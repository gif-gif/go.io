package main

import (
	gocontext "github.com/gif-gif/go.io/go-context"
	gojob "github.com/gif-gif/go.io/go-job"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/go-co-op/gocron/v2"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"time"
)

func main() {
	cron, _ := gojob.New()
	cron.Start()
	n := 0
	job, _ := cron.SecondX(&[]gocron.JobOption{
		gocron.WithSingletonMode(gocron.LimitModeWait)}, 1, func() {
		golog.Info("Starting-" + gconv.String(n))
		time.Sleep(5 * time.Second)
		n++
	})

	time.Sleep(time.Second * 10)
	err := cron.RemoveJob(job.ID())
	if err != nil {
		glog.Error(err)
	}
	glog.Info("Done1")

	job, _ = cron.SecondX(nil, 1, func() {
		golog.Info("Starting1")
	})

	time.Sleep(time.Second * 10)
	err = cron.RemoveJob(job.ID())
	if err != nil {
		glog.Error(err)
	}
	glog.Info("Done2")

	<-gocontext.WithCancel().Done()
}
