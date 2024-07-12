package main

import (
	"fmt"
	gojob "github.com/gif-gif/go.io/go-job"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/goio"
	"github.com/go-co-op/gocron/v2"
	"github.com/gogf/gf/util/gconv"
	"github.com/google/uuid"
	"time"
)

func main() {
	goio.Init(goio.DEVELOPMENT)
	//testJob()
	simpleUseGoJob()
}

func testJob() {
	n := 0
	cron, err := gojob.New()
	if err != nil {
		golog.WithTag("gojob").Error(err)
	}
	defer cron.Stop()
	cron.Start()

	job, err := cron.DurationJob(&[]gocron.JobOption{
		gocron.WithLimitedRuns(20),                        //最大执行次数
		gocron.WithSingletonMode(gocron.LimitModeWait),    // 限制重叠执行
		gocron.WithStartAt(gocron.WithStartImmediately()), //马上开始
		gocron.WithEventListeners(
			gocron.AfterJobRunsWithError(
				func(jobID uuid.UUID, jobName string, err error) {
					golog.WithTag("AfterJobRunsWithError-gojob").Error(jobID, jobName, err.Error())
					cron.Stop()
				},
			),
			gocron.AfterJobRunsWithPanic(
				func(jobID uuid.UUID, jobName string, err any) {
					golog.WithTag("AfterJobRunsWithPanic-gojob").Error(jobID, jobName, err)
					cron.Stop()
				},
			),
			gocron.AfterLockError(func(jobID uuid.UUID, jobName string, err error) {
				golog.WithTag("AfterLockError-gojob").Error(jobID, jobName, err.Error())
				cron.Stop()
			}),
		),
	}, 1, func(nn int) error {
		golog.WithTag("gojobStart").Info("testing->" + gconv.String(nn))
		time.Sleep(time.Second * 5)
		golog.WithTag("gojobEnd").Info("testing->" + gconv.String(nn))
		a := 1 / nn                                            //test for panic
		return fmt.Errorf("gojobEnd failed" + gconv.String(a)) //test for error
	}, n)

	if err != nil {
		golog.WithTag("gojob").Error(err)
	} else {
		golog.WithTag("gojob").Info("job.ID:" + job.ID().String())
	}

	time.Sleep(time.Second * 500)
	golog.InfoF("end of gojob")
}

func simpleUseGoJob() {
	n := 0
	cron, err := gojob.New()
	if err != nil {
		golog.WithTag("gojob").Error(err)
	}
	defer cron.Stop()
	cron.Start()

	job, err := cron.SecondX(nil, 1, func(nn int) error {
		golog.WithTag("gojob").Info("testing->" + gconv.String(nn))
		return nil
	}, n)

	if err != nil {
		golog.WithTag("gojob").Error(err)
	} else {
		golog.WithTag("gojob").Info("job.ID:" + job.ID().String())
	}

	time.Sleep(time.Second * 500)
	golog.InfoF("end of gojob")
}
