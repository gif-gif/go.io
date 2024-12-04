package gotime

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
)

var (
	TimeLayout string = "2006-01-02 15:04:05"
	DateLayout        = "20060102"
	TimeFormat        = map[string]string{
		"Y-m-d H:i:s": "2006-01-02 15:04:05",
		"Y-m-d":       "2006-01-02",
		"Ymd":         "20060102",
		"H:i:s":       "15:04:05",
		"Y":           "2006",
		"m":           "01",
		"d":           "02",
	}
)

// GetTimeNow 获取当前时间GetTimeNow()，用于测试时的时间修改
func GetTimeNow() time.Time {
	//redisConf := redis.RedisKeyConf{
	//	RedisConf: redis.RedisConf{
	//		Host: "122.228.113.235:17006",
	//		Type: "",
	//		Pass: "xiaozi527sport",
	//		TLS:  false,
	//	},
	//}
	//redisClient := redisConf.NewRedis()
	//val, err := redisClient.Get("test:now:gap:seconds")
	//if err != nil {
	//	fmt.Errorf("GetTimeNow test:now:gap error:%s", err.Error())
	//	return time.Now()
	//}
	//gap, _ := strconv.ParseInt(val, 10, 64)
	//return time.Now().Add(time.Second * time.Duration(gap))

	return time.Now()
}

func Ts2Time(t int64) time.Time {
	return time.Unix(t, 0)
}

func BeijingTimeLocation() *time.Location {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return loc
}

func GetChinaTomorrowAMSeconds(isBeijing bool) int64 {
	now := GetTimeNow()
	if isBeijing {
		loc, _ := time.LoadLocation("Asia/Shanghai")
		t, _ := time.ParseInLocation("2006-01-02", now.AddDate(0, 0, 1).Format("2006-01-02"), loc)
		secondsF := t.Sub(GetTimeNow()).Seconds()
		return int64(secondsF)
	} else {
		secondsF := now.Sub(GetTimeNow()).Seconds()
		return int64(secondsF)
	}
}

func GetLocalTomorrowAMSeconds() int64 {
	now := GetTimeNow()
	t, _ := time.ParseInLocation("2006-01-02", now.AddDate(0, 0, 1).Format("2006-01-02"), time.Local)
	secondsF := t.Sub(GetTimeNow()).Seconds()
	return int64(secondsF)
}

func GetTodayZero() time.Time {
	t := GetTimeNow()
	zero := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return zero
}

func GetZero(targetTime time.Time) time.Time {
	zero := time.Date(targetTime.Year(), targetTime.Month(), targetTime.Day(), 0, 0, 0, 0, targetTime.Location())
	return zero
}

// ParseTime  解析时间,"2021-03-17 00:00:00"
func ParseTime(timeStr string) (datetime time.Time) {
	datetime, _ = time.ParseInLocation(TimeLayout, timeStr, time.Local)
	return
}

//eg:20210812170000
func ParseTimeString(timeStr string) (datetime time.Time) {
	datetime, _ = time.ParseInLocation("20060102150405", timeStr, time.Local)
	return
}

func GetDateInterval(t1, t2 time.Time) int {
	t1 = time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	t2 = time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)

	interval := int(math.Abs(t1.Sub(t2).Hours())/24) + 1
	return interval
}

// SinceDays 获取过去的天数,dateString格式20060102
func SinceDays(dateString string) (int64, error) {
	targetTime, err := time.ParseInLocation("20060102", dateString, time.Local)
	if err != nil {
		return 0, err
	} else {
		days := int64(math.Ceil(time.Since(targetTime).Hours() / 24))
		if days == 0 {
			days = 1
		}
		return days, nil
	}
}

func IsSameDay(t1, t2 time.Time) bool {
	year1, month1, day1 := t1.Date()
	year2, month2, day2 := t2.Date()
	return day1 == day2 && month1 == month2 && year1 == year2
}

// LastHourStartAndEnd 上一个小时的开始和结束时间戳
func LastHourStartAndEnd(isBeijing bool) (int, int64, int64) {
	now1 := time.Now()
	var now time.Time
	if isBeijing {
		location, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			fmt.Println("Error loading location:", err)
			return 0, 0, 0
		}

		now = now1.In(location)
	} else {
		now = now1
	}

	startOfLastHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-1, 0, 0, 0, now.Location())
	endOfLastHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location()).Add(-time.Nanosecond)
	startOfLastHourUnix := startOfLastHour.Unix()
	endOfLastHourUnix := endOfLastHour.Unix()
	return startOfLastHour.Hour(), startOfLastHourUnix, endOfLastHourUnix
}

func CurrentHourStartAndEnd() (int, int64, int64) {
	now := time.Now()
	startOfHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
	endOfHour := startOfHour.Add(time.Hour - time.Second)

	startOfHourTimestamp := startOfHour.Unix()
	endOfHourTimestamp := endOfHour.Unix()
	return startOfHour.Hour(), startOfHourTimestamp, endOfHourTimestamp
}

func GetNowDateForLocation(isBeijing bool) string {
	tm := time.Now()
	var tmInLocation time.Time
	if isBeijing {
		// 加载时区
		location, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			fmt.Println("Error loading location:", err)
			return ""
		}
		// 将时间转换为指定时区
		tmInLocation = tm.In(location)
	} else {
		tmInLocation = tm
	}

	// 将时间格式化为日期字符串
	dateStr := tmInLocation.Format(time.DateOnly)
	return dateStr
}

func ToDate(ts int64, isBeijing bool) string {
	tm := time.Unix(ts, 0)
	var tmInLocation time.Time
	if isBeijing {
		// 加载时区
		location, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			fmt.Println("Error loading location:", err)
			return ""
		}

		// 将时间转换为指定时区
		tmInLocation = tm.In(location)
	} else {
		tmInLocation = tm
	}

	// 将时间格式化为日期字符串
	dateStr := tmInLocation.Format(time.DateOnly)
	return dateStr
}

func ToDateTime(ts int64, isBeijing bool) string {
	tm := time.Unix(ts, 0)
	var tmInLocation time.Time

	if isBeijing {
		// 加载时区
		location, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			fmt.Println("Error loading location:", err)
			return ""
		}

		// 将时间转换为指定时区
		tmInLocation = tm.In(location)
	} else {
		tmInLocation = tm
	}

	// 将时间格式化为日期字符串
	dateStr := tmInLocation.Format(time.DateTime)
	return dateStr
}

// 返回相差天数
func TimeRangeDay(stTime int64, endTIme int64) int {
	startTime := time.Unix(stTime, 0)
	endTime := time.Unix(endTIme, 0)
	durationDays := int(endTime.Sub(startTime).Hours()/24 + 1)
	return durationDays
}

func TimeRangeDates(startDate string, endDate string) []string {
	//startDate := "2021-01-01"
	//endDate := "2021-01-10"
	dates := []string{}
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		fmt.Println(d.Format("2006-01-02"))
		dates = append(dates, d.Format("2006-01-02"))
	}

	return dates
}

func DatesForRangeTs(startTimestamp int64, endTimestamp int64, format string) []string {
	if format == "" {
		format = "2006-01-02"
	}
	// 转换为time.Time对象
	startTime := time.Unix(startTimestamp, 0)
	endTime := time.Unix(endTimestamp, 0)

	// 创建一个当前时间的副本
	currentTime := startTime
	dates := []string{}
	// 循环直到达到结束时间
	for currentTime.Before(endTime) {
		// 格式化输出日期
		dates = append(dates, currentTime.Format(format))
		// 增加一天
		currentTime = currentTime.AddDate(0, 0, 1)
	}

	// 输出最后一天（如果最后一天的时间戳在当天内）
	if currentTime.Format(format) == endTime.Format(format) {
		dates = append(dates, currentTime.Format(format))
	}

	return dates
}

// convertTo24HourFormat 将时间字符串转换为 0-23 小时制
// timeStr: 时间字符串，例如 "3:04 PM"
// 3:04 PM-->15:04 ,必须是这两个格式之一
func ConvertAmPmHourTo24HourFormat(timeStr string, layout24Hour string) (string, error) {
	// 定义时间解析格式
	const layout12Hour = "3:04 PM"
	if layout24Hour == "" {
		layout24Hour = "15:04"
	}

	if !strings.HasPrefix(layout24Hour, "15") {
		return "", errors.New("layout24Hour must start with 15 or 15:04")
	}

	// 解析时间字符串
	t, err := time.Parse(layout12Hour, timeStr)
	if err != nil {
		return "", fmt.Errorf("failed to parse time: %v", err)
	}

	// 格式化为 24 小时制
	return t.Format(layout24Hour), nil
}

// 格式为： 2021-08-06T07:00:00+0000 To time.Time
func ConvertToGMTTime(gmtTime string) (time.Time, error) {
	t, err := time.Parse("2006-01-02T15:04:05-0700", gmtTime)
	return t, err
}

// 判断是不是 2024-11-22 4 格式的时间串
func IsValidDateTime(str string) bool {
	// 定义日期时间格式的正则表达式，忽略多余空格
	re := regexp.MustCompile(`^\d{4}-\d{1,2}-\d{1,2}\s+\d{1,2}$`)
	// 使用正则表达式匹配字符串
	isValidFormat := re.MatchString(strings.TrimSpace(str))
	return isValidFormat
}

func IsValidDate(str string) bool {
	// 定义日期时间格式的正则表达式，忽略多余空格
	re := regexp.MustCompile(`^\d{4}-\d{1,2}-\d{1,2}$`)
	// 使用正则表达式匹配字符串
	isValidFormat := re.MatchString(strings.TrimSpace(str))
	return isValidFormat
}
