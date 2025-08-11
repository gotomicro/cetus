package xtime

import (
	"fmt"
	"time"
)

const (
	Day  = 24 * time.Hour
	Week = 7 * Day
	Year = 365 * Day
)

func Millisecond() int64 {
	return time.Now().UnixNano() / 1e6
}
\
func Microsecond() int64 {
	return time.Now().UnixNano() / 1e3
}

func Cost(t time.Time, label ...string) {
	tc := time.Since(t)
	if len(label) > 0 {
		fmt.Printf("label: %s, cost: %v\n", label[0], tc)
	} else {
		fmt.Printf("cost: %v\n", tc)
	}
}

// Timestamp2String 会格式化当前时间。
func Timestamp2String(timestamp int) string {
	tm := time.Unix(int64(timestamp), 0)
	return tm.Format("2006-01-02 15:04:05")
}

func Timestamp2String64(timestamp int64) string {
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02 15:04:05")
}

// String2Timestamp 会将当前时间转换为时间戳。
func String2Timestamp(str string) int64 {
	tm, _ := time.Parse("2006-01-02 15:04:05", str)
	return tm.Unix()
}

// GetTodayZeroPoint ..
func GetTodayZeroPoint() int64 {
	timeStr := time.Now().Format("2006-01-02")
	loc, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, loc)
	return t.Unix()
}

// GetTodayZeroPointTime ..
func GetTodayZeroPointTime() time.Time {
	timeStr := time.Now().Format("2006-01-02")
	loc, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, loc)
	return t
}

// GetYesterdayZeroPoint ..
func GetYesterdayZeroPoint() int64 {
	timeStr := time.Now().Format("2006-01-02")
	loc, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, loc)
	return t.Unix() - int64(60*60*24)
}
