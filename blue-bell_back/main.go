package main

import (
	"blue-bell_back/controller"
	"blue-bell_back/dao/mysql"
	"blue-bell_back/dao/redis"
	"blue-bell_back/logger"
	"blue-bell_back/pkg/snowflake"
	"blue-bell_back/router"
	"blue-bell_back/settings"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Go Web开发较通用的脚手架模板

func main() {
	// 1.加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed,err:%v\n", err)
		return
	}
	// 2.初始化日志
	//if err := logger.Init(); err != nil {
	if err := logger.Init(settings.Conf.LogConfig); err != nil {
		fmt.Printf("init logger failed,err:%v\n", err)
		return
	}
	//将缓冲区的日志写入磁盘（防止程序结束前日志丢失）。
	//必须调用 .Sync()，否则最后几条日志可能不会落盘。
	defer zap.L().Sync()
	// 3.初始化MySQL连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed,err:%v\n", err)
		return
	}
	defer mysql.Close()
	// 4.初始化Redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed,err:%v\n", err)
		return
	}
	defer redis.Close()

	//5.初始化snowflake
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
	}

	// 6.初始化gin框架内置的校验器使用的翻译器
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Println("init validator trans err:%v", err)
		return
	}

	//7.注册路由
	r := router.Setup(settings.Conf.Mode)

	// 8.启动服务(优雅关机)
	srv := &http.Server{
		//Addr:    fmt.Sprintf(":#{settings.Conf.Port}"),  // 错误 TODO:
		Addr:    fmt.Sprintf(":%d", viper.GetInt("port")), //Go 中要用 fmt.Sprintf 或直接拼接
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { //是主动调用 Shutdown() 导致的关闭，属于正常行为，不打印错误。
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown:", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
