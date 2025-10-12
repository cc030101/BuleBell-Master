package router

import ( //倒入自定义的日志哭，用于记录API请求的日志和恢复异常
	//gin框架，构建HTTP服务器
	"blue-bell_back/controller"
	"blue-bell_back/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Setup函数用于初始化并配置gin框架，设置中间件和路由
// 返回值: *gin.Engine，初始化后的gin引擎实例

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	//创建一个新的gin引擎实例
	r := gin.New()
	//使用自定义的日志记录🍺异常恢复中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	//注册路由业务
	r.POST("/signup", controller.SignUpHandler)
	r.POST("/login", controller.LoginHandler)
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "ping")
	})

	//配置GET请求的路由，处理根路径的请求
	r.GET("/", func(c *gin.Context) {
		//响应客户端请求,返回HTTP状态码200和字符串"OK"
		c.String(http.StatusOK, "ok")
	})

	//返回初始化后的gin引擎
	return r
}
