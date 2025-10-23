package logic

import (
	"blue-bell_back/dao/mysql"
	"blue-bell_back/models"

	"go.uber.org/zap"
)

/*
	case1 direction = 1
		1.1	用户没有投过票 投了赞成
		1.2	用户投过反对票 改投赞成
	case2 direction = 0
		2.1	用户投过赞成票 取消了
		2.2 用过投过反对票 取消了
	case3 direction = -1
		3.1 用户没有投过票 投了赞成票
		3.2 用户投过赞成票 改投反对票

	投票功能的限制
	每个帖子自发布之日起7日内允许投票 超过该时间不允许再投票
	1.到期之后将redis中保存的帖子对应的赞成和反对票存入mysql
	2.到期之后删除 KeyPostVoteZSetPreFix

*/

// CommunityVote 帖子投票和功能逻辑处理
// 参数：p 包含帖子ID和投票方向的参数结构体
func CommunityVote(userID int64, p *models.ParamCommunityVote) (err error) {
	// 判断post id对应的帖子是否存在
	exist, err := mysql.CheckPostExist(uint64(p.PostID))
	if err != nil {
		zap.L().Error("mysql.CheckPostExist(post.CommunityID) failed.",
			zap.Uint64("post_id:", p.PostID),
			zap.Error(err))
		return
	}
	if !exist {
		zap.L().Error("post id not exist.",
			zap.Uint64("post_id:", p.PostID),
			zap.Error(err))
		return
	}
	// 调用Redis投票功能
	err = redis.VoteForCommunity(userID, int64(p.PostID), float64(p.Direction))
	if err != nil {
		zap.L().Error("redis.VoteForCommunity(userID, p.PostID, float64(p.Direction)).", zap.Error(err))
		return
	}
	return
}
