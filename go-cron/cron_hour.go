package gocrons

// 每小时执行
func Hour(fn ...func()) *crontab {
	return __cron.Hour()
}

// 每隔x小时执行
func HourX(x int, fn ...func()) *crontab {
	return __cron.HourX(x, fn...)
}
