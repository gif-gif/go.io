package goutils

import "time"

func DateTime(format string) string {
	return time.Now().Format(format)
}

func Today() string {
	return DateTime(time.DateOnly)
}

func Now() string {
	return DateTime(time.DateTime)
}

func NextDate(d int) string {
	return time.Now().AddDate(0, 0, d).Format(time.DateOnly)
}

func Ts2Date(ts int64) string {
	return time.Unix(ts, 0).Format(time.DateOnly)
}

// 本地化日期
func Ts2DateLocal(ts int64, timeOffsetSec int) string {
	t := time.Unix(ts, 0)
	location := time.FixedZone("OffsetZone", timeOffsetSec)
	// 将时间对象转换为特定时区的时间
	localTime := t.In(location)
	return localTime.Format(time.DateOnly)
}

func Ts2DateTime(ts int64) string {
	return time.Unix(ts, 0).Format(time.DateTime)
}

// 本地化日期时间
func Ts2DateTimeLocal(ts int64, timeOffsetSec int) string {
	t := time.Unix(ts, 0)
	location := time.FixedZone("OffsetZone", timeOffsetSec)
	// 将时间对象转换为特定时区的时间
	localTime := t.In(location)
	return localTime.Format(time.DateTime)
}

func Date2Ts(date string) int64 {
	ti, _ := time.ParseInLocation(time.DateOnly, date, time.Local)
	return ti.Unix()
}

func DateTime2Ts(dateTime string) int64 {
	ti, _ := time.ParseInLocation(time.DateTime, dateTime, time.Local)
	return ti.Unix()
}
