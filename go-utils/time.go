package goutils

import "time"

func DateTime(format string) string {
	return time.Now().Format(format)
}

func Today() string {
	return DateTime("2006-01-02")
}

func Now() string {
	return DateTime("2006-01-02 15:04:05")
}

func NextDate(d int) string {
	return time.Now().AddDate(0, 0, d).Format("2006-01-02")
}

func Ts2Date(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02")
}

func Ts2DateTime(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05")
}

func Date2Ts(date string) int64 {
	ti, _ := time.ParseInLocation("2006-01-02", date, time.Local)
	return ti.Unix()
}

func DateTime2Ts(dateTime string) int64 {
	ti, _ := time.ParseInLocation("2006-01-02 15:04:05", dateTime, time.Local)
	return ti.Unix()
}
