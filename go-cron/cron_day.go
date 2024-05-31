package gocrons

// 每天0点0分0秒执行
func Day(fn ...func()) *crontab {
	return __cron.Day(fn...)
}

// 每天x点0分0秒执行
func DayHour(hour int, fn ...func()) *crontab {
	return __cron.DayHour(hour, fn...)
}

// 每天x点x分0秒执行
func DayHourMinute(hour, minute int, fn ...func()) *crontab {
	return __cron.DayHourMinute(hour, minute, fn...)
}
