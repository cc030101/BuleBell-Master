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
