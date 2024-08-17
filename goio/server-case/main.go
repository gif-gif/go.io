package main

import (
	"context"
	"github.com/gif-gif/go.io/goio"
	"github.com/gif-gif/go.io/goio/server-case/router"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	startSever()
}

func startSever() {
	goio.Init(goio.DEVELOPMENT)
	s := goio.NewServer(
		goio.ServerNameOption("serverName"),
		goio.EnvOption(goio.Env),
		//goio.EnableEncryptionOption(),
		goio.PProfEnableOption(false),
		goio.NoLogPathsOption("/captcha/get"),
	)

	// 路由
	router.Routes(s.Group("/"))
	// 启动
	s.Run(":1000")
}

// 简介的web服务器启动
func simpleServer() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// 静态页面返回
	//r.LoadHTMLGlob("templates/*")
	//r.GET("/", func(c *gin.Context) {
	//	c.HTML(200, "index.tmpl", gin.H{
	//		"title": "Parse-video",
	//	})
	//})

	srv := &http.Server{
		Addr:    ":1000",
		Handler: r,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器 (设置 5 秒的超时时间)
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
