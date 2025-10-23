package mysql

import (
	"blue-bell_back/models"
	"database/sql"

	"go.uber.org/zap"
)

// GetCommunityList查询所有社区
//此函数从数据库中查询所有社区的信息，并返回社区列表
// 如果查询结果为空， 则不会返回错误，而是将错误置空

func GetCommunityList() (communityList []*models.Community, err error) {
	sqlStr := "select community_id, community_name from community"
	if err := db.Select(&communityList, sqlStr); err != nil {
		//如果查询为空
		if err == sql.ErrNoRows {
			zap.L().Warn("no result from community table")
			err = nil
		}
	}
	return
}

// GetCommunityByID 根据社区ID查询社区详情
// 此函数根据给定的社区ID查询社区的详细信息。
// 如果ID有效且找到对应的社区，则返回社区详情。
// 如果ID无效或没有找到对应的社区，则返回ErrorInvalidID错误。
func GetCommunityByID(id int64) (communityDetail *models.CommunityDetail, err error) {
	//申请内存
	communityDetail = new(models.CommunityDetail)
	sqlStr := "select community_id, community_name, introducion, create_time from community where id = ?"

	if err := db.Get(communityDetail, sqlStr, id); err != nil {
		//判断id是否有效
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
			return nil, err
		}
	}
	return
}

// CreateCommunityPost 创建社区的帖子
// 参数： post：指向包含帖子信息的CommunityPost结构体的指针，包含帖子ID、标题、作者ID、社区ID和内容
// 返回值：如果插入操作成功，则返回nil；否则返回错误信息
func CreateCommunityPost(post *models.CommunityPost) (err error) {
	//定义sql语句来插入帖子信息到数据库
	sqlStr := "insert into post(post_id, title, author_id, community_id, content) value(?,?,?,?,?)"
	//执行sql语句，插入帖子信息，并检查是否有错误发生
	_, err = db.Exec(sqlStr, post.ID, post.Title, post.AuthorID, post.CommunityID, post.Content)

	if err != nil {
		return err
	}
	return
}

// GetAuthorNameById 根据用户id查询用户名称
func GetAuthorNameById(userID uint64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := "select user_id, username from user where user_id = ?"
	if err := db.Get(user, sqlStr, userID); err != nil {
		//判断id是否有效
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
			return nil, err
		}
	}
	return
}

// GetPostDetailByID 根据帖子ID查询帖子详情
func GetPostDetailByID(postId uint64) (postDetail *models.CommunityPost, err error) {
	postDetail = new(models.CommunityPost)
	sqlStr := "select post_id, title, content, author_id, community_id, status, create_time from post where post_id = ?"
	//执行sql查询，并将结果存储到postDetail中
	if err := db.Get(postDetail, sqlStr, postId); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID
			return nil, err
		}
	}
	return
}

//GetPostList查询帖子列表

func GetPostList(page, size int64) (list []*models.CommunityPost, err error) {
	list = make([]*models.CommunityPost, 0, 2)
	sqlStr := "select post_id, title, content, author_id, community_id, status, create_time from post order by create_time desc limit ?,?"
	err = db.Select(&list, sqlStr, (page-1)*size, size)
	return
}

//CheckPostExist 检查帖子是否存在

func CheckPostExist(id string) (exist bool, err error) {
	var count int

	sqlStr := "select count(community_id) from post where post_id = ?"
	if err := db.Get(&count, sqlStr, id); err != nil {
		return exist, err
	}
	zap.L().Debug("select count from post", zap.Int("count:", count))
	if count > 0 {
		return true, nil
	}

	return exist, nil
}
