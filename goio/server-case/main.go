package main

import (
	"context"
	"flag"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/go-utils/prometheusx"
	"github.com/gif-gif/go.io/goio"
	"github.com/gif-gif/go.io/goio/server"
	conf "github.com/gif-gif/go.io/goio/server-case/config"
	"github.com/gif-gif/go.io/goio/server-case/router"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	yamlFile = flag.String("yaml", "etc/api-local.yaml", "yaml file")
)

func main() {
	startSever()
}

func startSever() {
	// 初始化命令行参数
	goio.FlagInit()

	// yaml文件必须
	if *yamlFile == "" {
		flag.PrintDefaults()
		golog.Fatal("yaml file not exist.")
	}

	// 加载配置文件
	confs := &conf.Config{}
	err := goserver.LoadYamlConfig(*yamlFile, confs)
	if err != nil {
		golog.WithTag("main").Error(err)
		return
	}
	goio.Init(confs.Env)
	goio.SetupLogDefault()
	goio.Setup("")
	prometheusx.Init(confs.Prometheus)
	prometheusx.AlertErr(confs.Server.Name, "main start")

	s := goserver.NewServer(
		goserver.ServerNameOption("serverName"),
		goserver.EnvOption(goio.Env),
		//goserver.EnableEncryptionOption("1a3295a2408d553a8458085e7435898e", "119f54388848cb4306f6d2067a4713fce4193504ca368d648196c840ba87da65"),
		goserver.PProfEnableOption(false),
		goserver.NoLogPathsOption("/captcha/get"),
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
