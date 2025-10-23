package redis

import "blue-bell_back/models"

func GetPostListByID(p *models.ParamOrderList) ([]string, error) {
	//1.从redis获取id
	//根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZSet)

	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}

	//2.确定查询的索引起始点
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	//3.zrevrange查询
	return rdb.ZRevRange(key, start, end).Result()
}
