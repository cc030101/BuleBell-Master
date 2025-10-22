package logic

import (
	"blue-bell_back/dao/mysql"
	"blue-bell_back/models"
	"blue-bell_back/pkg/snowflake"

	"go.uber.org/zap"
)

//	社区相关
//	GetCommunityList 处理获取社区列表

func GetCommunityList() ([]*models.Community, error) {
	return mysql.GetCommunityList()
}

// GetCommunityDetail 处理获取社区详情
func GetCommunityDetail(id int64) (detail *models.CommunityDetail, err error) {
	return mysql.GetCommunityByID(id)
}

// CreateCommunityPost 创建帖子
func CreateCommunityPost(p *models.CommunityPost) (err error) {
	//1.生成id
	var id int64
	id = snowflake.GenID()
	p.ID = id
	//2.保存到数据库
	return mysql.CreateCommunityPost(p)
}

// GetPostDetail 获取帖子详情
func GetPostDetail(id uint64) (detail *models.ApiPostDetail, err error) {
	post, err := mysql.GetPostDetailByID(id)
	if err != nil {
		zap.L().Error("mysql.GetPostDetailByID(id) failed.",
			zap.Uint64("authorID:", id),
			zap.Error(err))
		return
	}
	// 1.根据作者id查询作者用户名
	user, err := mysql.GetAuthorNameById(uint64(post.AuthorID))
	if err != nil {
		zap.L().Error("mysql.GetAuthorNameById(post.AuthorID) failed.",
			zap.Int64("authorID:", post.AuthorID),
			zap.Error(err))
		return
	}
	// 2.根据社区id查询社区名称
	communityDetail, err := mysql.GetCommunityByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityNameById(post.CommunityID) failed.",
			zap.Int64("authorID:", post.CommunityID),
			zap.Error(err))
		return
	}
	// 3.组装帖子详细信息并返回
	detail = &models.ApiPostDetail{
		AuthorName:      user.UserName,
		CommunityDetail: communityDetail,
		CommunityPost:   post,
	}
	return
}
