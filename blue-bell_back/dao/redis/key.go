package redis

//redis key

const (
	KeyPreFix             = "Forum:"
	KeyPostTimeZSet       = "post:time"   //帖子及发帖时间
	KeyPostScoreZSet      = "post:score"  //帖子及投票的分数
	KeyPostVoteZSetPreFix = "post:voted:" //记录用户的投票类型
)

func getRedisKey(key string) string {
	return KeyPreFix + key
}
