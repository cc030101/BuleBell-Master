package controller

import (
	"blue-bell_back/logic"
	"blue-bell_back/models"
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

// CreatePostHandler 创建帖子函数
func CreatePostHandler(c *gin.Context) {
	//1.获取参数
	post := new(models.CommunityPost)
	if err := c.ShouldBindJSON(post); err != nil {
		//如果参数一场就记录日志并返回错误
		zap.L().Debug("c.ShouldBindJSON(post) failed.", zap.Any("err", err))
		zap.L().Error("create community failed.", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	//2.获取用户id
	userId, err := getCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	post.AuthorID = int64(userId)
	//3. 存储数据
	if err := logic.CreateCommunityPost(post); err != nil {
		//创建失败 返回错误信息
		zap.L().Error("service.CreateCommunityPost failed.", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	//4. 返回
	ResponseSuccess(c, post)
}

// PostDetailHandler 帖子详情函数
func PostDetailHandler(c *gin.Context) {
	//1.拿到postID
	id := c.Param("id")
	postId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		zap.L().Error("get post detail failed. invalid param", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}

	//2.根据帖子id查询帖子详情
	detail, err := logic.GetPostDetail(uint64(postId))
	if err != nil {
		zap.L().Error("service.GetPostDetail failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSuccess(c, detail)
}

// GetPostListHandler 获取帖子列表函数
func GetPostListHandler(c *gin.Context) {
	page, size, err := GetPageInfo(c)

	if err != nil {
		page = 1
		size = 10
	}

	//2.获取数据
	list, err := logic.GetPostList(page, size)
	if err != nil {
		zap.L().Error("service.GetPostList failed.", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}

	ResponseSuccess(c, list)
}
