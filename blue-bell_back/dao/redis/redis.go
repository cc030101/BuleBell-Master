package redis

import (
	"blue-bell_back/settings"
	"fmt"

	"github.com/go-redis/redis"
)

var rdb *redis.Client //rdb是全局的Redis客户端实例

// Init 初始化 Redis 客户端
// 参数 cfg 是 Redis 服务器的配置信息
// 返回可能的错误
func Init(cfg *settings.RedisConfig) (err error) {
	//根据配置信息创建Redis客户端实例
	rdb = redis.NewClient(&redis.Options{
		//使用viper读取配置信息，格式转化为host:port形式
		Addr: fmt.Sprintf("%s:%d",
			// viper.GetString("redis.host"),
			// viper.GetInt("redis.port"),
			cfg.Host,
			cfg.Port,
		),
		// //读取redis密码配置
		// Password: viper.GetString("redis.password"),
		// //读取redis数据库编号配置
		// DB: viper.GetInt("redis.db"),
		// //读取连接池大小配置
		// PoolSize: viper.GetInt("redis.pool_size"),
		// //读取最小空闲连接数配置
		// MinIdleConns: viper.GetInt("redis.mid_idle_conns"),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	//测试链接是否成果
	_, err = rdb.Ping().Result()
	return
}

// Close 关闭 Redis 客户端连接
func Close() {
	_ = rdb.Close()
}
