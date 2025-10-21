package controller

import (
	"blue-bell_back/dao/mysql"
	"blue-bell_back/logic"
	"blue-bell_back/models"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignUpHandler 处理用户注册请求
// 参数: c *gin.Context 提供了请求的上下文，用于处理HTTP请求和响应
func SignUpHandler(c *gin.Context) {
	//1. 获取参数和参数校验
	p := new(models.ParamSignUp)
	if err := c.ShouldBindBodyWithJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("SignUp with invalid param", zap.Error(err))
		// 判断err是不是validator.ValidationErrors 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidPassword)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidPassword, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 手动对请求参数进行详细的业务规则校验
	//if len(p.UserName) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.Password != p.RePassword {
	//	// 请求参数有误，直接返回响应
	//	zap.L().Error("SignUp with invalid param")
	//	c.JSON(http.StatusOK, gin.H{
	//		"msg": "请求参数有误",
	//	})
	//	return
	//}

	fmt.Println(p)
	//2. 业务处理
	if err := logic.SignUp(p); err != nil {
		zap.L().Error("login.SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorInvalidPassword) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	//3. 返回响应
	ResponseSueccess(c, nil)
}
func LoginHandler(c *gin.Context) {
	// 1.获取请求参数及参数校验
	p := new(models.ParamLogin)
	if err := c.ShouldBindJSON(p); err != nil {
		// 请求参数有误，直接返回响应
		zap.L().Error("Login with invalid param", zap.Error(err))
		// 判断err是不是validator.ValidationErrors 类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}

	// 2.业务逻辑处理
	if err := logic.Login(p); err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}
	// 3.返回响应
	ResponseSueccess(c, nil)
}
