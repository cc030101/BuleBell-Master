package logger

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack" //用于日志轮转（按大小切割日志文件）
	"github.com/spf13/viper"          //配置管理库，从配置文件读取日志设置。
	"go.uber.org/zap"                 //高性能结构化日志库。
	"go.uber.org/zap/zapcore"
)

// Init 初始化lg
// 该函数读取配置文件中的日志设置，初始化一个zap.Logger，并将其设置为全局默认使用
func Init() (err error) {
	writeSyncer := getLogWriter( //调用 getLogWriter() 创建一个支持日志轮转的写入器。
		viper.GetString("log.filename"), //日志文件路径
		viper.GetInt("log.max_size"),    //单个日志文件最大 MB 数、保留备份个数、保留天数。
		viper.GetInt("log.max_backups"),
		viper.GetInt("log.max_age"),
	)
	encoder := getEncoder() //调用 getEncoder() 定义日志输出格式（JSON 格式
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(viper.GetString("log.level")))
	if err != nil {
		return
	}
	core := zapcore.NewCore(encoder, writeSyncer, l)

	lg := zap.New(core, zap.AddCaller())
	// 替换zap库中全局的logger
	zap.ReplaceGlobals(lg) //将这个 logger 设为全局默认 logger，后续使用 zap.L() 即可获取。
	return
}

// getEncoder 返回一个zapcore.Encoder，用于配置日志的编码格式
// 该函数设置了日志的时间格式、级别格式、持续时间格式和调用者格式
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder //2025-04-05T12:34:56.789+0800
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder       //日志级别大写显示（如 "INFO", "ERROR"）
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder //耗时以秒为单位
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder       //显示简短调用者信息（如 main.go:123）
	return zapcore.NewJSONEncoder(encoderConfig)

	//返回一个 JSON 编码器，最终日志看起来像
	// {"time":"2025-04-05T12:34:56.789+0800","level":"INFO","caller":"main.go:50","msg":"/api/users","status":200,"method":"GET",...}
}

// getLogWriter 返回一个zapcore.WriteSyncer，用于配置日志的写入方式
// 该函数设置了日志文件的名称、最大大小、备份文件数量和最大保存天数
func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	//lumberjack.Logger 实现日志自动切割
	return zapcore.AddSync(lumberJackLogger)
}

// GinLogger 接收gin框架默认的日志
// 该函数记录了每个HTTP请求的路径、查询参数、状态码、方法、IP、用户代理、错误和耗时
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		zap.L().Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
			//请求路径、查询参数、状态码、HTTP 方法、客户端 IP、User-Agent、错误信息、耗时。
			//{"time":"...","level":"INFO","caller":"...","msg":"/login","status":200,"method":"POST","path":"/login","query":"","ip":"192.168.1.100","user-agent":"Mozilla/...","errors":"","cost":0.015}
		)
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
// 该函数捕获panic异常，记录错误日志，并中止请求处理
func GinRecovery(stack bool) gin.HandlerFunc { //Gin Panic 恢复中间件
	return func(c *gin.Context) {
		defer func() { //在 defer 中 recover 捕获 panic。
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool //特殊处理 broken pipe 错误
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false) //记录请求头等信息
				if brokenPipe {
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: err check
					c.Abort()
					return
				}

				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
