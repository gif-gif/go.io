package gojob

import (
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"time"
)

type (
	CronsModel struct {
		cron gocron.Scheduler
	}
)

func New(options ...gocron.SchedulerOption) (*CronsModel, error) {
	c, err := gocron.NewScheduler(options...)
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

func (c *CronsModel) RemoveJob(jobID uuid.UUID) error {
	return c.cron.RemoveJob(jobID)
}

//job options start

// 在某一时刻运行
// s, _ := NewScheduler()
// defer func() { _ = s.Shutdown() }()
//
// start := time.Date(9999, 9, 9, 9, 9, 9, 9, time.UTC)
//
// j, _ := s.NewJob(
//
//	DurationJob(
//		time.Second,
//	),
//	NewTask(
//		func(one string, two int) {
//			fmt.Printf("%s, %d", one, two)
//		},
//		"one", 2,
//	),
//	WithStartAt(
//		WithStartDateTime(start),
//	),
//
// )
// s.Start()
//
// next, _ := j.NextRun()
// fmt.Println(next)
//
// _ = s.StopJobs()
//
// 定时执行启动 开始时间
func (c *CronsModel) WithStartAt(start time.Time) gocron.JobOption {
	//start := time.Date(9999, 9, 9, 9, 9, 9, 9, time.UTC)
	return gocron.WithStartAt(
		gocron.WithStartDateTime(start),
	)
}

// interval 月频, 0-6-->周日 周一 ... 周六, hours 具体执行时间列表
func (c *CronsModel) MonthlyJob(options *[]gocron.JobOption, interval uint, daysOfTheMonth []int, hours []uint, minute uint, fn any, parameters ...any) (gocron.Job, error) {
	if options == nil {
		options = &[]gocron.JobOption{}
	}
	hoursAtTime := []gocron.AtTime{}
	for _, hour := range hours {
		hoursAtTime = append(hoursAtTime, gocron.NewAtTime(hour, minute, 0))
	}

	return c.cron.NewJob(
		gocron.MonthlyJob(
			interval,
			gocron.NewDaysOfTheMonth(daysOfTheMonth[0], daysOfTheMonth[1:]...),
			gocron.NewAtTimes(
				gocron.NewAtTime(hours[0], minute, 0),
				hoursAtTime...,
			),
		),
		gocron.NewTask(
			fn,
			parameters,
		),
		*options...,
	)
}

// interval 周频, 0-6-->周日 周一 ... 周六, hours 具体执行时间列表
func (c *CronsModel) WeeklyJob(options *[]gocron.JobOption, interval uint, daysOfTheWeek []time.Weekday, hours []uint, minutes uint, fn any, parameters ...any) (gocron.Job, error) {
	if options == nil {
		options = &[]gocron.JobOption{}
	}
	hoursAtTime := []gocron.AtTime{}
	for _, hour := range hours {
		hoursAtTime = append(hoursAtTime, gocron.NewAtTime(hour, minutes, 0))
	}

	return c.cron.NewJob(
		gocron.WeeklyJob(
			interval,
			gocron.NewWeekdays(daysOfTheWeek[0], daysOfTheWeek[1:]...),
			gocron.NewAtTimes(
				gocron.NewAtTime(hours[0], minutes, 0),
				hoursAtTime...,
			),
		),
		gocron.NewTask(
			fn,
			parameters...,
		),
		*options...,
	)
}

// 当前时间 seconds 秒之后执行一次
func (c *CronsModel) OneTimeJobForSeconds(options *[]gocron.JobOption, seconds uint, fn any, parameters ...any) (gocron.Job, error) {
	if options == nil {
		options = &[]gocron.JobOption{}
	}
	return c.cron.NewJob(
		gocron.OneTimeJob(
			gocron.OneTimeJobStartDateTime(time.Now().Add(time.Duration(seconds)*time.Second)),
		),
		gocron.NewTask(
			fn,
			parameters...,
		),
		*options...,
	)
}

// 当前时间 minute 分钟之后执行一次
func (c *CronsModel) OneTimeJobForMinute(options *[]gocron.JobOption, minute uint, fn any, parameters ...any) (gocron.Job, error) {
	if options == nil {
		options = &[]gocron.JobOption{}
	}
	return c.cron.NewJob(
		gocron.OneTimeJob(
			gocron.OneTimeJobStartDateTime(time.Now().Add(time.Duration(minute)*time.Minute)),
		),
		gocron.NewTask(
			fn,
			parameters...,
		),
		*options...,
	)
}

// 每天定时执行
func (c *CronsModel) DailyJob(options *[]gocron.JobOption, interval uint, hours []uint, minute uint, fn any, parameters ...any) (gocron.Job, error) {
	if options == nil {
		options = &[]gocron.JobOption{}
	}

	hoursAtTime := []gocron.AtTime{}
	for _, hour := range hours {
		hoursAtTime = append(hoursAtTime, gocron.NewAtTime(hour, minute, 0))
	}

	return c.cron.NewJob(
		gocron.DailyJob(interval, gocron.NewAtTimes(
			gocron.NewAtTime(hours[0], minute, 0),
			hoursAtTime...,
		)),
		gocron.NewTask(
			fn,
			parameters...,
		),
		*options...,
	)
}

// 隔多少秒执行 , 可以通过 gocron.WithSingletonMode(gocron.LimitModeWait)来限制重叠执行
func (c *CronsModel) DurationJob(options *[]gocron.JobOption, seconds int, fn any, parameters ...any) (gocron.Job, error) {
	if options == nil {
		options = &[]gocron.JobOption{}
	}

	return c.cron.NewJob(
		gocron.DurationJob(
			time.Duration(seconds)*time.Second,
		),
		gocron.NewTask(
			fn,
			parameters...,
		),

		//gocron.WithSingletonMode(gocron.LimitModeReschedule),
		*options...,
	)
}

// DurationRandomJob 定义一个新作业，该作业以提供的最小和最大持续时间值之间的随机间隔运行
func (c *CronsModel) DurationRandomJob(options *[]gocron.JobOption, minDuration, maxDuration time.Duration, function any, parameters ...any) (gocron.Job, error) {
	if options == nil {
		options = &[]gocron.JobOption{}
	}
	return c.cron.NewJob(
		gocron.DurationRandomJob(
			minDuration, maxDuration,
		),
		gocron.NewTask(
			function,
			parameters...,
		),
		*options...,
	)
}

// spec is crontab pattern crontab 表达式
func (c *CronsModel) CronJob(spec string, options *[]gocron.JobOption, function any, parameters ...any) (gocron.Job, error) {
	if options == nil {
		options = &[]gocron.JobOption{}
	}
	return c.cron.NewJob(
		gocron.CronJob(
			// standard cron tab parsing
			spec,
			true,
		),
		gocron.NewTask(
			function,
			parameters...,
		),
		*options...,
	)
}

// crontab 每天0点0分0秒执行
func (c *CronsModel) Day(options *[]gocron.JobOption, fn any, parameters ...any) (gocron.Job, error) {
	return c.CronJob("0 0 0 * * *", options, fn, parameters...)
}

// crontab 每天x点0分0秒执行
func (c *CronsModel) DayHour(options *[]gocron.JobOption, hour int, fn any, parameters ...any) (gocron.Job, error) {
	return c.CronJob(fmt.Sprintf("0 0 %d * * *", hour), options, fn, parameters...)
}

// crontab 每天x点x分0秒执行
func (c *CronsModel) DayHourMinute(options *[]gocron.JobOption, hour, minute int, fn any, parameters ...any) (gocron.Job, error) {
	return c.CronJob(fmt.Sprintf("0 %d %d * * *", minute, hour), options, fn, parameters...)
}

// crontab 每小时执行
func (c *CronsModel) Hour(options *[]gocron.JobOption, fn any, parameters ...any) (gocron.Job, error) {
	return c.CronJob("0 0 */1 * * *", options, fn, parameters...)
}

// crontab 每隔x小时执行
func (c *CronsModel) HourX(options *[]gocron.JobOption, x int, fn any, parameters ...any) (gocron.Job, error) {
	return c.CronJob(fmt.Sprintf("0 0 */%d * * *", x), options, fn, parameters...)
}

// crontab 每分钟执行
func (c *CronsModel) Minute(options *[]gocron.JobOption, fn any, parameters ...any) (gocron.Job, error) {
	return c.CronJob("0 */1 * * * *", options, fn, parameters...)
}

// crontab 每隔x分钟执行
func (c *CronsModel) MinuteX(options *[]gocron.JobOption, x int, fn any, parameters ...any) (gocron.Job, error) {
	return c.CronJob(fmt.Sprintf("0 */%d * * * *", x), options, fn, parameters...)
}

// crontab 每秒钟执行
func (c *CronsModel) Second(options *[]gocron.JobOption, fn any, parameters ...any) (gocron.Job, error) {
	return c.CronJob("* * * * * *", options, fn, parameters...)
}

// crontab 每隔x秒执行
func (c *CronsModel) SecondX(options *[]gocron.JobOption, x int, fn any, parameters ...any) (gocron.Job, error) {
	return c.CronJob(fmt.Sprintf("*/%d * * * * *", x), options, fn, parameters...)
}
