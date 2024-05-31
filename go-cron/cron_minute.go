package gocrons

// 每分钟执行
func Minute(fn ...func()) *crontab {
	return __cron.Minute(fn...)
}

// 每隔x分钟执行
func MinuteX(x int, fn ...func()) *crontab {
	return __cron.MinuteX(x, fn...)
}
