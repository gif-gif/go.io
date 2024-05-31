package gocrons

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

var (
	__cron = &crontab{c: cron.New(cron.WithSeconds())}
)

type crontab struct {
	c *cron.Cron
}

func Cron() *cron.Cron {
	return __cron.c
}

func (c *crontab) Start() {
	c.c.Start()
}

func (c *crontab) Stop() {
	c.c.Stop()
}

func (c *crontab) Func(spec string, fn ...func()) *crontab {
	for _, f := range fn {
		c.c.AddFunc(spec, f)
	}
	return c
}

func (c *crontab) Job(spec string, job ...cron.Job) *crontab {
	for _, j := range job {
		c.c.AddJob(spec, j)
	}
	return c
}

// 每天0点0分0秒执行
func (c *crontab) Day(fn ...func()) *crontab {
	return __cron.Func("0 0 0 * * *", fn...)
}

// 每天x点0分0秒执行
func (c *crontab) DayHour(hour int, fn ...func()) *crontab {
	return __cron.Func(fmt.Sprintf("0 0 %d * * *", hour), fn...)
}

// 每天x点x分0秒执行
func (c *crontab) DayHourMinute(hour, minute int, fn ...func()) *crontab {
	return __cron.Func(fmt.Sprintf("0 %d %d * * *", minute, hour), fn...)
}

// 每小时执行
func (c *crontab) Hour(fn ...func()) *crontab {
	return __cron.Func("0 0 */1 * * *", fn...)
}

// 每隔x小时执行
func (c *crontab) HourX(x int, fn ...func()) *crontab {
	return __cron.Func(fmt.Sprintf("0 0 */%d * * *", x), fn...)
}

// 每分钟执行
func (c *crontab) Minute(fn ...func()) *crontab {
	return __cron.Func("0 */1 * * * *", fn...)
}

// 每隔x分钟执行
func (c *crontab) MinuteX(x int, fn ...func()) *crontab {
	return __cron.Func(fmt.Sprintf("0 */%d * * * *", x), fn...)
}

// 每秒钟执行
func (c *crontab) Second(fn ...func()) *crontab {
	return __cron.Func("* * * * * *", fn...)
}

// 每隔x秒执行
func (c *crontab) SecondX(x int, fn ...func()) *crontab {
	return __cron.Func(fmt.Sprintf("*/%d * * * * *", x), fn...)
}
