package gocrons

// 每秒钟执行
func Second(fn ...func()) *crontab {
	return __cron.Second(fn...)
}

// 每隔x秒执行
func SecondX(x int, fn ...func()) *crontab {
	return __cron.SecondX(x, fn...)
}
