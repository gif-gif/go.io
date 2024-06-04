package gocrons

import (
	"fmt"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/gif-gif/go.io/goio"
	"github.com/robfig/cron/v3"
)

type (
	CronsModel struct {
		cron *cron.Cron
	}
)

func New() *CronsModel {
	o := &CronsModel{
		cron: cron.New(cron.WithSeconds()),
	}
	return o
}

func (o *CronsModel) Cron() *cron.Cron {
	return o.cron
}

func (c *CronsModel) Start() {
	c.cron.Start()
}

func (c *CronsModel) Stop() {
	c.cron.Stop()
}

func (c *CronsModel) Func(spec string, fn ...func()) {
	for _, f := range fn {
		_, err := c.cron.AddFunc(spec, f)
		if err != nil {
			if goio.Env == goio.DEVELOPMENT {

			}
			golog.WithTag("gocrons").Error(err)
		}
	}
}

func (c *CronsModel) Job(spec string, job ...cron.Job) {
	for _, j := range job {
		_, err := c.cron.AddJob(spec, j)
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
