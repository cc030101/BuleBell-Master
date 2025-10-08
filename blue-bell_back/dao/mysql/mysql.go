package mysql

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// db用于存储全局的数据库连接实例
var db *sqlx.DB

//Init初始化数据库连接
//该函数读取配置文件中国呢的数据库连接信息，构建DSN(数据源名称),并尝试建立数据库连接
//如果连接成果，会配置数据库连接池的最大打开连接数和最大空闲连接
//如果连接失败，会记录错误日志并返回错误

func Init() (err error) {
	//构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.dbname"),
	)
	//也可以用must connect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
		return
	}

	//设置数据库连接池的最大打开连接数
	db.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))
	//设置数据库连接吃的最大空闲连接数
	db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns"))
	return
}

// Close 关闭数据库连接
// 该函数用于释放数据库连接资源，应在程序退出前调用
func Close() {
	_ = db.Close()
}
