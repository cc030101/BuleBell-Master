package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// CtxUserIDKey 是用于在 Gin 上下文中存储用户 ID 的键。
// 它用于在后续处理程序中检索用户 ID。
const CtxUserIDKey = "userID"

//ErrorUserNotLogin 是一个错误，表示当前用户未登录
//当尝试获取当前用户失败时，将返回此错误

var ErrorUserNotLogin = errors.New("当前用户位登陆")

// getCurrentUser 尝试从 Gin 上下文中获取当前用户的 ID。
// 如果用户未登录或用户 ID 类型不正确，将返回 ErrorUserNotLogin 错误。
// 参数:c - Gin 上下文指针，用于从中提取用户 ID。
// 返回值:userID - 用户 ID，如果成功获取。
// err - 错误，如果获取失败。

func getCurrentUser(c *gin.Context) (userID uint64, err error) {
	//从上下文获取ID
	uid, ok := c.Get(CtxUserIDKey)
	if !ok {
		//未找到则返回错误
		err = ErrorUserNotLogin
		return
	}
	//将用户ID转换为uint64类型
	userID, ok = uid.(uint64)
	if !ok {
		// 转换失败，返回错误
		err = ErrorUserNotLogin
		return
	}
	return
}
