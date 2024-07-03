package gojob

import (
	"fmt"
	golog "github.com/gif-gif/go.io/go-log"
	"github.com/go-co-op/gocron/v2"
	"time"
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

// interval 月频, 0-6-->周日 周一 ... 周六, hours 具体执行时间列表
func (c *CronsModel) MonthlyJob(interval uint, daysOfTheMonth []int, hours []uint, minutes uint, fn ...func()) {
	hoursAtTime := []gocron.AtTime{}
	for _, hour := range hours {
		hoursAtTime = append(hoursAtTime, gocron.NewAtTime(hour, minutes, 0))
	}

	for _, f := range fn {
		_, err := c.cron.NewJob(
			gocron.MonthlyJob(
				interval,
				gocron.NewDaysOfTheMonth(daysOfTheMonth[0], daysOfTheMonth[1:]...),
				gocron.NewAtTimes(
					gocron.NewAtTime(hours[0], minutes, 0),
					hoursAtTime...,
				),
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

// interval 周频, 0-6-->周日 周一 ... 周六, hours 具体执行时间列表
func (c *CronsModel) WeeklyJob(interval uint, daysOfTheWeek []time.Weekday, hours []uint, minutes uint, fn ...func()) {
	hoursAtTime := []gocron.AtTime{}
	for _, hour := range hours {
		hoursAtTime = append(hoursAtTime, gocron.NewAtTime(hour, minutes, 0))
	}

	for _, f := range fn {
		_, err := c.cron.NewJob(
			gocron.WeeklyJob(
				interval,
				gocron.NewWeekdays(daysOfTheWeek[0], daysOfTheWeek[1:]...),
				gocron.NewAtTimes(
					gocron.NewAtTime(hours[0], minutes, 0),
					hoursAtTime...,
				),
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

// 当前时间 seconds 秒之后执行一次
func (c *CronsModel) OneTimeJobForSeconds(seconds uint, fn ...func()) {
	for _, f := range fn {
		_, err := c.cron.NewJob(
			gocron.OneTimeJob(
				gocron.OneTimeJobStartDateTime(time.Now().Add(time.Duration(seconds)*time.Second)),
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

// 当前时间 minute 分钟之后执行一次
func (c *CronsModel) OneTimeJobForMinute(minute uint, fn ...func()) {
	for _, f := range fn {
		_, err := c.cron.NewJob(
			gocron.OneTimeJob(
				gocron.OneTimeJobStartDateTime(time.Now().Add(time.Duration(minute)*time.Minute)),
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

// 每天定时执行
func (c *CronsModel) DailyJob(interval uint, hours []uint, minute uint, fn ...func()) {
	hoursAtTime := []gocron.AtTime{}
	for _, hour := range hours {
		hoursAtTime = append(hoursAtTime, gocron.NewAtTime(hour, minute, 0))
	}

	for _, f := range fn {
		_, err := c.cron.NewJob(
			gocron.DailyJob(interval, gocron.NewAtTimes(
				gocron.NewAtTime(hours[0], minute, 0),
				hoursAtTime...,
			)),
			gocron.NewTask(
				f,
			),
		)
		if err != nil {
			golog.WithTag("gocrons").Error(err)
		}
	}
}

// 隔多少秒执行
func (c *CronsModel) DurationJob(seconds int, fn ...func()) {
	for _, f := range fn {
		_, err := c.cron.NewJob(
			gocron.DurationJob(
				time.Duration(seconds)*time.Second,
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

// DurationRandomJob 定义一个新作业，该作业以提供的最小和最大持续时间值之间的随机间隔运行
func (c *CronsModel) DurationRandomJob(minDuration, maxDuration time.Duration, fn ...func()) {
	for _, f := range fn {
		_, err := c.cron.NewJob(
			gocron.DurationRandomJob(
				minDuration, maxDuration,
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

// spec is crontab pattern crontab 表达式
func (c *CronsModel) CronJob(spec string, fn ...func()) {
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

// crontab 每天0点0分0秒执行
func (c *CronsModel) Day(fn ...func()) {
	c.CronJob("0 0 0 * * *", fn...)
}

// crontab 每天x点0分0秒执行
func (c *CronsModel) DayHour(hour int, fn ...func()) {
	c.CronJob(fmt.Sprintf("0 0 %d * * *", hour), fn...)
}

// crontab 每天x点x分0秒执行
func (c *CronsModel) DayHourMinute(hour, minute int, fn ...func()) {
	c.CronJob(fmt.Sprintf("0 %d %d * * *", minute, hour), fn...)
}

// crontab 每小时执行
func (c *CronsModel) Hour(fn ...func()) {
	c.CronJob("0 0 */1 * * *", fn...)
}

// crontab 每隔x小时执行
func (c *CronsModel) HourX(x int, fn ...func()) {
	c.CronJob(fmt.Sprintf("0 0 */%d * * *", x), fn...)
}

// crontab 每分钟执行
func (c *CronsModel) Minute(fn ...func()) {
	c.CronJob("0 */1 * * * *", fn...)
}

// crontab 每隔x分钟执行
func (c *CronsModel) MinuteX(x int, fn ...func()) {
	c.CronJob(fmt.Sprintf("0 */%d * * * *", x), fn...)
}

// crontab 每秒钟执行
func (c *CronsModel) Second(fn ...func()) {
	c.CronJob("* * * * * *", fn...)
}

// crontab 每隔x秒执行
func (c *CronsModel) SecondX(x int, fn ...func()) {
	c.CronJob(fmt.Sprintf("*/%d * * * * *", x), fn...)
}
