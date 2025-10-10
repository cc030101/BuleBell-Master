package router

import (
	//倒入自定义的日志哭，用于记录API请求的日志和恢复异常
	"blue-bell_back/logger"
	//net/http库，处理HTTP请求
	"net/http"
	//gin框架，构建HTTP服务器
	"github.com/gin-gonic/gin"
)

// 返回值为初始化后的gin引擎实例
func Setup() *gin.Engine {
	//创建一个新的gin引擎实例
	r := gin.New()

	//使用自定义的日志记录和异常恢复中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	//配置GET请求的路由，处理根路径的请求
	r.GET("/", func(c *gin.Context) {
		//相应客户端请求，返回HTTP状态码200和字符串
		c.String(http.StatusOK, "ok")
	})
	//返回初始化后的gin实例
	return r
}
