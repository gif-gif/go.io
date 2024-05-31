package gojob

import (
	"fmt"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/go-co-op/gocron/v2"
)

type (
	CronsModel struct {
		cron gocron.Scheduler
	}
)

func New() (*CronsModel, error) {
	c, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}
	o := &CronsModel{
		cron: c,
	}
	return o, nil
}

func (o *CronsModel) Cron() gocron.Scheduler {
	return o.cron
}

func (c *CronsModel) Start() {
	c.cron.Start()
}

func (c *CronsModel) Stop() error {
	return c.cron.Shutdown()
}

func (c *CronsModel) Func(spec string, fn ...func()) {
	for _, f := range fn {
		_, err := c.cron.NewJob(
			gocron.CronJob(
				// standard cron tab parsing
				spec,
				true,
			),
			gocron.NewTask(
				f,
			),
		)
		if err != nil {
			golog.WithTag("gocrons").Error(err)
		}
	}
}

// 每天0点0分0秒执行
func (c *CronsModel) Day(fn ...func()) {
	c.Func("0 0 0 * * *", fn...)
}

// 每天x点0分0秒执行
func (c *CronsModel) DayHour(hour int, fn ...func()) {
	c.Func(fmt.Sprintf("0 0 %d * * *", hour), fn...)
}

// 每天x点x分0秒执行
func (c *CronsModel) DayHourMinute(hour, minute int, fn ...func()) {
	c.Func(fmt.Sprintf("0 %d %d * * *", minute, hour), fn...)
}

// 每小时执行
func (c *CronsModel) Hour(fn ...func()) {
	c.Func("0 0 */1 * * *", fn...)
}

// 每隔x小时执行
func (c *CronsModel) HourX(x int, fn ...func()) {
	c.Func(fmt.Sprintf("0 0 */%d * * *", x), fn...)
}

// 每分钟执行
func (c *CronsModel) Minute(fn ...func()) {
	c.Func("0 */1 * * * *", fn...)
}

// 每隔x分钟执行
func (c *CronsModel) MinuteX(x int, fn ...func()) {
	c.Func(fmt.Sprintf("0 */%d * * * *", x), fn...)
}

// 每秒钟执行
func (c *CronsModel) Second(fn ...func()) {
	c.Func("* * * * * *", fn...)
}

// 每隔x秒执行
func (c *CronsModel) SecondX(x int, fn ...func()) {
	c.Func(fmt.Sprintf("*/%d * * * * *", x), fn...)
}
