package gotime

import "time"

// 当前时区相关日期函数
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

func Ts2DateTime(ts int64) string {
	return time.Unix(ts, 0).Format(time.DateTime)
}

func Date2Ts(date string) int64 {
	ti, _ := time.ParseInLocation(time.DateOnly, date, time.Local)
	return ti.Unix()
}

func DateTime2Ts(dateTime string) int64 {
	return DateTime2TsLocalFormat(dateTime, time.Local, time.DateTime)
}

//--------------本地化时间--------------
//没有format 时日志格式为  time.DateOnly 和 time.DateTime

// 当前时间通用 格式用法
func DateTimeLocal(format string, timeOffsetSec int) string {
	return Ts2DateTimeLocalFormat(time.Now().Unix(), timeOffsetSec, format)
}

// NowDate
func TodayLocal(timeOffsetSec int) string {
	return Ts2DateLocal(time.Now().Unix(), timeOffsetSec)
}

// NowDateTime
func NowLocal(timeOffsetSec int) string {
	return Ts2DateTimeLocal(time.Now().Unix(), timeOffsetSec)
}

// NextDate
func NextDateLocal(d int, timeOffsetSec int) string {
	location := time.FixedZone("OffsetZone", timeOffsetSec)
	t := time.Unix(time.Now().Unix(), 0)
	t.In(location)
	return Ts2DateTimeLocalFormat(t.AddDate(0, 0, d).Unix(), timeOffsetSec, time.DateOnly)
}

// 本地化日期 给定时间戳（以秒为单位）
func Ts2DateLocal(ts int64, timeOffsetSec int) string {
	return Ts2DateTimeLocalFormat(ts, timeOffsetSec, time.DateOnly)
}

// 本地化日期时间 给定时间戳（以秒为单位）
func Ts2DateTimeLocal(ts int64, timeOffsetSec int) string {
	return Ts2DateTimeLocalFormat(ts, timeOffsetSec, time.DateTime)
}

func Date2TsLocal(date string, location *time.Location) int64 {
	return DateTime2TsLocalFormat(date, location, time.DateOnly)
}

func DateTime2TsLocal(dateTime string, location *time.Location) int64 {
	return DateTime2TsLocalFormat(dateTime, location, time.DateTime)
}

func DateTime2TsLocalFormat(dateTime string, location *time.Location, format string) int64 {
	ti, _ := time.ParseInLocation(format, dateTime, location)
	return ti.Unix()
}

// 将时间对象转换为特定时区的时间
func Ts2DateTimeLocalFormat(ts int64, timeOffsetSec int, format string) string {
	t := time.Unix(ts, 0)
	location := time.FixedZone("OffsetZone", timeOffsetSec)
	// 将时间对象转换为特定时区的时间
	localTime := t.In(location)
	return localTime.Format(format)
}
