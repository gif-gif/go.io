package main

import (
	gocontext "github.com/gif-gif/go.io/go-context"
	goredis "github.com/gif-gif/go.io/go-db/go-redis"
	gojob "github.com/gif-gif/go.io/go-job"
	golog "github.com/gif-gif/go.io/go-log"
	goasynq "github.com/gif-gif/go.io/go-mq/go-asynq"
	"github.com/gif-gif/go.io/go-mq/go-asynq/test/tasks"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/hibiken/asynq"
	"log"
	"time"
)

const redisAddr = "127.0.0.1:6379"

func main() {
	goutils.AsyncFunc(func() {
		serverTest()
	})
	clientAsynq()
	<-gocontext.Cancel().Done()
}

func clientAsynq() {
	goasynq.InitClient(goasynq.ClientConfig{
		Config: goredis.Config{
			Addr: redisAddr,
		},
	})
	client := goasynq.DefaultClient()
	defer goasynq.DefaultClient().Close()

	info, err := client.Enqueue(tasks.TypeEmailDelivery, tasks.EmailDeliveryPayload{UserID: 42, TemplateID: "some:template:id"}, asynq.ProcessIn(24*time.Second))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	info, err = client.Enqueue(tasks.TypeImageResize, tasks.ImageResizePayload{SourceURL: "https://example.com/myassets/image.jpg"}, asynq.MaxRetry(10), asynq.Timeout(3*time.Second))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	cron, err := gojob.New()
	if err != nil {
		golog.WithTag("gojob").Error(err)
	}

	cron.Start()
	n := 0
	cron.SecondX(nil, 1, func() { //for test time
		n++
		golog.WithTag("SecondX").Info(n)
	})
}

func serverTest() {
	goasynq.InitServer(goasynq.ServerConfig{
		Config: goredis.Config{
			Addr: redisAddr,
		},
	})
	server := goasynq.DefaultServer()
	server.HandleFunc(tasks.TypeEmailDelivery, tasks.HandleEmailDeliveryTask)
	server.Handle(tasks.TypeImageResize, tasks.NewImageProcessor())
	golog.WithTag("wwww").Info("server running")
}
