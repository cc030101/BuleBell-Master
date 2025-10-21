package middlewares

import (
	"blue-bell_back/controller"
	"blue-bell_back/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationKey = "Authorization"
	BearerKey        = "Bearer"
	NullKey          = ""
	SpaceKey         = " "
)

//JWTAuthMiddleware 基于JWT的认证中间件

func JWTAuthMiddleware() func(c *gin.Context) {

	//token存放在header中的authorization，并使用bearer开头
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get(AuthorizationKey)
		if authHeader == NullKey {
			controller.ResponseError(c, controller.CodeNeedLogin)
			c.Abort()
			return
		}
		//按空格分割
		parts := strings.SplitN(authHeader, SpaceKey, 2)
		if !(len(parts) == 2) && parts[0] == BearerKey {
			controller.ResponseError(c, controller.CodeInvalidParam)
			c.Abort()
			return
		}

		//parts[1]获取到的是token string，使用之前定义好的解析函数
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			controller.ResponseError(c, controller.CodeInvalidParam)
			c.Abort()
			return
		}

		//将当前请求的userID信息保存在请求的上下文c中
		c.Set(controller.CtxUserIDKey, mc.UserID)
		c.Next() //后续可以通过c.get(CtxUserIDKey)获取用户信息
	}
}
