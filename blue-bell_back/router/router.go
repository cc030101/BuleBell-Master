package router

import ( //倒入自定义的日志哭，用于记录API请求的日志和恢复异常
	//gin框架，构建HTTP服务器
	"blue-bell_back/controller"
	"blue-bell_back/logger"
	"blue-bell_back/middlewares"
	"blue-bell_back/pkg/jwt"

	"net/http"
	"strings"

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

	v1 := r.Group("/api/v1")

	//注册路由业务
	// r.POST("/signup", controller.SignUpHandler)
	// r.POST("/login", controller.LoginHandler)

	// r.GET("/ping", JWTAuthMiddleware(), func(c *gin.Context) {
	// 	//如果是登陆用户，判断请求头中是否有有效的JWT
	// 	c.String(http.StatusOK, "pong")
	// })
	v1.POST("/signup", controller.SignUpHandler)
	v1.POST("/login", controller.LoginHandler)
	v1.Use(middlewares.JWTAuthMiddleware())
	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)
	}

	v1.POST("/community/post", controller.CreatePostHandler)

	//配置GET请求的路由，处理根路径的请求
	r.GET("/", func(c *gin.Context) {
		//响应客户端请求,返回HTTP状态码200和字符串"OK"
		c.String(http.StatusOK, "ok")
	})

	//返回初始化后的gin引擎
	return r
}

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// Authorization: Bearer xxxxxxx.xxx.xxxxxxx
		// 这里的具体实现方式要依据实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请求头中auth为空",
			})
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set("userID", mc.UserID)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}
