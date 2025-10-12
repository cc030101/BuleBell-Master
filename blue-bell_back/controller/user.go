package controller

import (
	"blue-bell_back/logic"
	"blue-bell_back/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignUpHandler 处理用户注册请求
// 参数: c *gin.Context 提供了请求的上下文，用于处理HTTP请求和响应
func SignUpHandler(c *gin.Context) {
	//1.获取参数和参数校验

	p := new(models.ParamSignUp)
	if err := c.ShouldBindBodyWithJSON(p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))

		//判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(), //翻译错误
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": removeTopStruct(errs.Translate(trans)),
		})
	}
	return

	// 手动对请求参数进行详细的业务规则校验
	//if len(p.UserName) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.Password != p.RePassword {
	//    // 请求参数有误，直接返回响应
	//    zap.L().Error("SignUp with invalid param")
	//    c.JSON(http.StatusOK, gin.H{
	//        "msg": "请求参数有误",
	//    })
	//    return
	//}

	fmt.Println(p)

	//2.业务处理

	if err := logic.SignUp(p); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "注册失败",
		})
		return
	}

	//3.返回响应
	c.JSONP(http.StatusOK, gin.H{
		"msg": "success",
	})
}

func LoginHandler(c *gin.Context) {
	//1.获取请求参数及参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		//请求参数有误，直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))
		//判断err是不是validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(), //翻译错误
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": removeTopStruct(errs.Translate(trans)),
		})

		return
	}

	//2.业务逻辑处理
	if err := logic.Login(p); err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.UserName), zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"msg": "用户名或密码错误",
		})
		return
	}

	//3.返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "登录成功",
	})

}
