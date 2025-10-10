package mysql

import (
	"blue-bell_back/settings"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// db用于存储全局的数据库连接实例
var db *sqlx.DB

// 参数: cfg: 包含MySQL配置信息的结构体指针
// 返回值: 可能发生的错误
func Init(cfg *settings.MySQLConfig) (err error) {
	//构建DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		// viper.GetString("mysql.user"),
		// viper.GetString("mysql.password"),
		// viper.GetString("mysql.host"),
		// viper.GetInt("mysql.port"),
		// viper.GetString("mysql.dbname"),
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)

	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect DB failed", zap.Error(err))
		return
	}

	//设置数据库连接池的最大打开连接数
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	//设置数据库连接吃的最大空闲连接数
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return
}

func Close() {
	_ = db.Close()
}
