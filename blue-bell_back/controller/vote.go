package controller

import (
	"blue-bell_back/logic"
	"blue-bell_back/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

//CommunityVote 处理社区投票的函数
// 验证投票参数的合法性，然后获取当前的用户ID

func CommunityVote(c *gin.Context) {
	//参数校验
	p := new(models.ParamCommunityVote)
	if err := c.ShouldBindJSON(p); err != nil {
		//类型断言
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		errData := removeTopStruct(errs.Translate(trans))
		ResponseErrorWithMsg(c, CodeInvalidParam, errData)
	}

	//获取用户id
	userID, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}

	//具体的投票业务逻辑
	err = logic.CommunityVote(userID, p)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}

	ResponseSuccess(c, nil)
}
