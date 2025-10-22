package controller

import (
	"blue-bell_back/logic"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 社区相关
// CommunityHandler 处理获取社区列表的函数
// 该函数查询所有社区的信息，并返回给客户端

func CommunityHandler(c *gin.Context) {
	//1.查询到所有社区的信息(community_id, community_name)
	list, err := logic.GetCommunityList()
	if err != nil {
		zap.L().Error("service.GetCommunityList failed.", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, list)
}

// CommunityDetailHandler 处理获取社区详情的函数
// 该函数根据社区ID查询社区详情，并返回给客户端

func CommunityDetailHandler(c *gin.Context) {
	//1.拿到id
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	//根据社区id查询社区详情
	detail, err := logic.GetCommunityDetail(id)
	if err != nil {
		zap.L().Error("service.GetCommunityList failed.", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, detail)
}
