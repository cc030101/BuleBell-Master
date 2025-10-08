package redis

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
)

var rdb *redis.Client //rdb是一个全局的Redis客户端实例

//Init初始化连接
//该函数读取配置信息，创建并验证与Redis服务器的连接
//返回值：如果连接成果，则err为nil； 否则， err中包含错误信息

func Init() (err error) {
	//根据配置信息创建Redis客户端实例
	rdb = redis.NewClient(&redis.Options{
		//使用viper读取配置信息，格式转化为host:port形式
		Addr: fmt.Sprintf("%s:%d",
			viper.GetString("redis.host"),
			viper.GetInt("redis.port"),
		),
		//读取redis密码配置
		Password: viper.GetString("redis.password"),
		//读取redis数据库编号配置
		DB: viper.GetInt("redis.db"),
		//读取连接池大小配置
		PoolSize: viper.GetInt("redis.pool_size"),
		//读取最小空闲连接数配置
		MinIdleConns: viper.GetInt("redis.mid_idle_conns"),
	})

	//测试链接是否成果
	_, err = rdb.Ping().Result()
	return
}

//Close 关闭Redis连接
//释放Redis客户端实例占用的资源

func Close() {
	_ = rdb.Close()
}
