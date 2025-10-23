package logic

import (
	"blue-bell_back/dao/mysql"
	"blue-bell_back/dao/redis"
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
	err = mysql.CreateCommunityPost(p)
	if err != nil {
		return err
	}

	//保存到redis
	err = redis.CreateCommunityPost(int64(id))
	return
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

// GetPostList 获取帖子列表逻辑
func GetPostList(page, size int64) (list []*models.ApiPostDetail, err error) {
	// 调用mysql.GetPostList获取帖子列表
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		// 如果获取帖子列表失败，记录错误日志并返回
		zap.L().Error("mysql.GetPostList failed.", zap.Error(err))
		return
	}
	// 初始化帖子详细信息列表
	list = make([]*models.ApiPostDetail, 0, len(posts))
	// 循环posts获取用户名和社区名称
	for _, post := range posts {
		// 1.根据作者id查询作者用户名
		author, err := mysql.GetAuthorNameById(uint64(post.AuthorID))
		if err != nil {
			// 如果获取作者用户名失败，记录错误日志并继续处理下一个帖子
			zap.L().Error("mysql.GetAuthorNameById(post.AuthorID) failed.",
				zap.Int64("authorID:", post.AuthorID),
				zap.Error(err))
			continue
		}
		// 2.根据社区id查询社区名称
		community, err := mysql.GetCommunityByID(post.CommunityID)
		if err != nil {
			// 如果获取社区名称失败，记录错误日志并继续处理下一个帖子
			zap.L().Error("mysql.GetAuthorNameById(post.AuthorID) failed.",
				zap.Int64("authorID:", post.AuthorID),
				zap.Error(err))
			continue
		}
		// 3.构建帖子详细信息对象
		apiPostDetail := &models.ApiPostDetail{
			AuthorName:      author.UserName,
			CommunityDetail: community,
			CommunityPost:   post,
		}
		// 4.将帖子详细信息添加到列表中并返回
		list = append(list, apiPostDetail)
	}
	return
}
