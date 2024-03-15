package helper

import (
	"strconv"
	"time"
)

func NanoTimestampStr() string {
	ts := time.Now().UnixNano()
	return strconv.FormatInt(ts, 10)
}

func UnixSecondToTime(second int64) time.Time {
	return time.Unix(second, 0)
}

func UnixSecondToStr(second int64) string {
	timeLayout := "2006-01-02 15:04:05"
	timeStr := time.Unix(second, 0).Format(timeLayout)
	return timeStr
}
func UnixSecondToInt64() int64 {
	return time.Now().Unix()
}
