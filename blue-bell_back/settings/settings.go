package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

//Init 初始化配置信息
//返回： 可能发生的错误

func Init() (err error) {
	//指定配置文件 不需要带后缀
	viper.SetConfigName("config")
	//指定配置文件类型
	viper.SetConfigType("yaml")
	//指定查找路径 使用绝对路径/相对路径
	viper.AddConfigPath(".")
	//读取配置信息
	err = viper.ReadInConfig()

	//查看是否读取失败
	if err != nil {
		fmt.Printf("Viper read config failed, err%v\n", err)
		return
	}

	//监控配置文件的变化
	viper.WatchConfig()
	//当配置文件发生变化 回调函数启动
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
	})
	return
}
