package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
)

// Go Web开发通用的脚手架
func main() {
	//1.	加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("Init settings failed, err :%v\n", err)
	}
	//2.	初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("Init logger failed, err :%v\n", err)
	}
	defer zap.L().Sync()
	//3.	初始化Mysql连接
	if err := mysql.Init(); err != nil {
		fmt.Printf("Init mysql failed, err :%v\n", err)
	}
	defer mysql.Close()
	//4.	初始化redis连接
	if err := redis.Init(); err != nil {
		fmt.Printf("Init redis failed, err :%v\n", err)
	}
	defer redis.Close()

	//5.	路由注册
	r := router.Setup()

	//6.	启动服务(优雅关机)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}

	go func() {
		//开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen %s\n", err)
		}
	}()

	//等待中断信号来优雅地关闭服务器，为关闭 服务器操作设置一个5s的超时
	quit := make(chan os.Signal, 1) //创建一个接收信号的通道
	//kill 发送syscall.SIGTERM信号
	//kill -2 发送syscall.SIGINT信号， 比如ctrl + c
	//kill -9 发送syscall.SIGKILL信号， 但是不能被捕获
	//signal.Notify把刚收到的 syscall.SIGINT或者syscall.SIGTERM信号发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //此处不会阻塞
	<-quit                                               //阻塞再次，当接受到上述两种信号才会往下执行
	zap.L().Info("Shutdown Server ...")
	//创建一个5s超时到context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//5s内优雅关闭服务(将为处理完的请求处理完再关闭服务)，超过5s就退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown:", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
