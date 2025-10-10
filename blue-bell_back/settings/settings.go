package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf是一个全局变量，用于
var Conf = new(AppConfig)

// 定义了应用程序配置的结构体，包含了应用程序的基本配置信息
type AppConfig struct {
	// Name         string                 `mapstructure:"name"` //应用程序名称
	// Mode         string                 `mapstructure:"mode"` //程序运行模式
	// Port         int                    `mapstructure:"port"` //程序运行端口

	Name         string                 `mapstructure:"name"`       // 应用程序名称
	Mode         string                 `mapstructure:"mode"`       // 应用程序运行模式
	Port         int                    `mapstructure:"port"`       // 应用程序运行端口
	StartTime    string                 `mapstructure:"start_time"` // 程序开始时间
	MachineID    int64                  `mapstructure:"machine_id"` //机器ID
	*LogConfig   `mapstructure:"log"`   //日志配置信息
	*MySQLConfig `mapstructure:"mysql"` //MySQL数据库配置信息
	*RedisConfig `mapsturcture:"redis"` //redis数据库配置信息
}

// LogConfig定义了日志配置的结构体，包含了日志的相关配置信息
type LogConfig struct {
	Level      string `mapstructure:"level"`      //日志级别
	FileName   string `mapstructure:"filename"`   //日志文件名
	MaxSize    int    `mapstructure:"max_size"`   //日志文件最大的尺寸（MB）
	MaxAge     int    `mapstructure:"max_age"`    //日志文件最大保留的天数
	MaxBackups int    `mapstructure:max_backups"` //最多保留的日志文件个数

}

//MySQLConfig定义了MySQL数据库配置的结构体，包含了连接MySQL数据库所需的信息

type MySQLConfig struct {
	Host         string `mapstructure:"host"`           //数据库主机地址
	Port         int    `mapstructure:"port"`           //数据库端口
	User         string `mapstructure:"user"`           //数据库用户名
	Password     string `mapstructure:"password"`       //
	DbName       string `mapstructure::"dbname"`        //数据库
	MaxOpenConns int    `mapstructure:"max_open_conns"` //最大打开的连接数
	MaxIdleConns int    `mapstructure:max_idle_conns"`  //最大空闲连接数
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`      //Redis主机地址
	Port     int    `mapstructure:"port"`      //端口
	Password string `mapstructure:"password"`  //
	DB       int    `mapstructure:"db"`        //数据编号
	PoolSize int    `mapstructure:"pool_size"` //连接池大小
}

func Init() (err error) {
	// //指定配置文件 不需要带后缀
	// viper.SetConfigName("config")
	// //指定配置文件类型
	// viper.SetConfigType("yaml")
	// //指定查找路径 使用绝对路径/相对路径
	// viper.AddConfigPath(".")
	// //读取配置信息
	// err = viper.ReadInConfig()

	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath(".")      //指定查找路径 使用绝对路径/相对路径
	viper.AddConfigPath("./conf") //可以设置多个查找配置文件路径

	//读取配置信息
	err = viper.ReadInConfig()
	//查看是否读取失败
	if err != nil {
		// fmt.Printf("Viper read config failed, err%v\n", err)
		fmt.Printf("viper.ReadInconfig failed, err:%v\n", err)
		return
	}

	//把读取到的配置信息反序列到Conf全局变量中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshall failed, err:%v\n", err)
	}

	//监控配置文件的变化
	viper.WatchConfig()
	//当配置文件发生变化 回调函数启动
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshall failed, err:%v\n", err)
		}
	})
	return
}
