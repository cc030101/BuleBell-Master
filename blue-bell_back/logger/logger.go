package logger

import (
	"blue-bell_back/settings"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack" //用于日志轮转（按大小切割日志文件）

	//配置管理库，从配置文件读取日志设置。
	"go.uber.org/zap" //高性能结构化日志库。
	"go.uber.org/zap/zapcore"
)

// Init// Init 初始化日志系统
// cfg: 日志配置信息，包含日志文件名、最大大小、备份文件最大数量和最大保留天数
// func Init() (err error) {
func Init(cfg *settings.LogConfig) (err error) {
	writeSyncer := getLogWriter( //调用 getLogWriter() 创建一个支持日志轮转的写入器。
		// viper.GetString("log.filename"), //日志文件路径
		// viper.GetInt("log.max_size"),    //单个日志文件最大 MB 数、保留备份个数、保留天数。
		// viper.GetInt("log.max_backups"),
		// viper.GetInt("log.max_age"),
		cfg.FileName,
		cfg.MaxSize,
		cfg.MaxBackups,
		cfg.MaxAge,
	)
	encoder := getEncoder() //调用 getEncoder() 定义日志输出格式（JSON 格式
	var l = new(zapcore.Level)
	//err = l.UnmarshalText([]byte(viper.GetString("log.level")))
	//将配置中的日志级别文本转换为zapcore.Level类型
	err = l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return
	}

	//创建zapcore.Core实例
	core := zapcore.NewCore(encoder, writeSyncer, l)

	//创建新的Logger实例
	lg := zap.New(core, zap.AddCaller())
	// 替换zap库中全局的logger
	zap.ReplaceGlobals(lg) //将这个 logger 设为全局默认 logger，后续使用 zap.L() 即可获取。
	return
}

// getEncoder 创建一个日志编码器
// 返回值: zapcore.Encoder类型的日志编码器
func getEncoder() zapcore.Encoder {
	// 创建一个适用于生产环境的编码配置。
	encoderConfig := zap.NewProductionEncoderConfig()
	// 设置时间戳的编码方式为 ISO8601 格式。
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder //2025-04-05T12:34:56.789+0800
	// 将时间戳字段的键名设置为 "time"。
	encoderConfig.TimeKey = "time"
	// 设置日志级别字段的编码方式为大写形式。
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder //日志级别大写显示（如 "INFO", "ERROR"）
	// 设置持续时间字段的编码方式为秒。
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder //耗时以秒为单位
	// 设置调用者信息字段的编码方式为短格式。
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder //显示简短调用者信息（如 main.go:123）
	// 根据上述配置返回一个新的 JSON 编码器。
	return zapcore.NewJSONEncoder(encoderConfig)

	//返回一个 JSON 编码器，最终日志看起来像
	// {"time":"2025-04-05T12:34:56.789+0800","level":"INFO","caller":"main.go:50","msg":"/api/users","status":200,"method":"GET",...}
}

// getLogWriter 创建一个日志写入器
// 参数: filename: 日志文件名 maxSize: 单个日志文件最大大小（MB）
// maxBackup: 最大备份文件数量 maxAge: 最大保留天数
// 返回值: zapcore.WriteSyncer类型的日志写入器
func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	// 创建一个Lumberjack日志记录器实例，并配置其参数
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	//lumberjack.Logger 实现日志自动切割
	// 将Lumberjack日志记录器包装为zapcore.WriteSyncer接口类型并返回。
	// 这样做使得Lumberjack可以与Zap日志库无缝集成。
	return zapcore.AddSync(lumberJackLogger)
}

// GinLogger 创建一个gin框架使用的日志记录中间件
// 返回值: gin.HandlerFunc类型的中间件
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()
		// 获取请求路径
		path := c.Request.URL.Path
		// 获取查询字符串
		query := c.Request.URL.RawQuery
		// 调用c.Next()以继续执行链中的其他中间件和处理函数
		c.Next()

		// 计算请求处理时间
		cost := time.Since(start)
		// 使用zap记录请求信息
		zap.L().Info(path,
			// 记录HTTP状态码
			zap.Int("status", c.Writer.Status()),
			// 记录HTTP方法
			zap.String("method", c.Request.Method),
			// 记录请求路径
			zap.String("path", path),
			// 记录查询字符串
			zap.String("query", query),
			// 记录客户端IP
			zap.String("ip", c.ClientIP()),
			// 记录用户代理
			zap.String("user-agent", c.Request.UserAgent()),
			// 记录私有错误信息
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			// 记录请求处理时间
			zap.Duration("cost", cost),
			//请求路径、查询参数、状态码、HTTP 方法、客户端 IP、User-Agent、错误信息、耗时。
			//{"time":"...","level":"INFO","caller":"...","msg":"/login","status":200,"method":"POST","path":"/login","query":"","ip":"192.168.1.100","user-agent":"Mozilla/...","errors":"","cost":0.015}
		)
	}
}

// GinRecovery 是一个中间件，用于恢复出现的panic，防止服务崩溃。
// 参数stack表示是否在日志中包含调用栈信息。
// 返回一个gin.HandlerFunc处理函数。
func GinRecovery(stack bool) gin.HandlerFunc { //Gin Panic 恢复中间件
	return func(c *gin.Context) {
		// 使用defer来捕获执行过程中的panic。
		defer func() { //在 defer 中 recover 捕获 panic。
			if err := recover(); err != nil {
				// 初始化断开连接错误标志。
				var brokenPipe bool //特殊处理 broken pipe 错误
				if ne, ok := err.(*net.OpError); ok {
					// 判断内部错误是否为os.SyscallError类型。
					if se, ok := ne.Err.(*os.SyscallError); ok {
						//if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
						// 	brokenPipe = true
						// }
						if strings.Contains(strings.ToLower(se.Error()),
							"broken pipe") || strings.Contains(strings.ToLower(se.Error()),
							"connection reset by peer") {
							// 设置断开连接错误标志。
							brokenPipe = true
						}
					}
				}

				// 尝试获取完整的HTTP请求信息。
				httpRequest, _ := httputil.DumpRequest(c.Request, false) //记录请求头等信息
				// 如果是客户端断开连接错误，记录错误日志并终止处理。
				if brokenPipe {
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // 错误检查
					c.Abort()
					return
				}
				// 如果不是客户端断开连接错误，根据配置决定是否记录调用栈信息。
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
				// 终止处理并返回500内部服务器错误状态码。
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		// 继续执行请求的下一个处理函数。
		c.Next()
	}
}
