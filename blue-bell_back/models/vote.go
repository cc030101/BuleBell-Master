package models

type ParamCommunityVote struct {
	PostID    uint64 `json:"post_id,string" binding:"required"`                 // 帖子id
	Direction int8   `json:"direction,string" bingding:"required,oneof=1 0 -1"` //赞成1反对-1取消0
}
