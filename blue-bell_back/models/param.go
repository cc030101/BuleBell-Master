package models

const (
	Page       = 10
	Size       = 1
	OrderTime  = "time"
	OrderScore = "score"
)

// 注册请求参数结构体
type ParamSignUp struct {
	UserName   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// ParamLogin 登陆请求参数
type ParamLogin struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ParamOrderList 获取帖子列表
type ParamOrderList struct {
	Page  int64  `json:"page" form:"page"`
	Size  int64  `json:"size" form:"size"`
	Order string `json:"order" form:"order"`
}

// ParamCommunityPostList 社区下帖子列表的接口
type ParamCommunityPostList struct {
	*ParamOrderList
	CommunityID int64 `json:"community_id" form:"community_id"`
}
