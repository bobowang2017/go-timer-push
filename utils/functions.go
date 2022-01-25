package utils

import (
	"time"
)

//计算当前日期是第几周
func WeekCountNow() int {
	_, thisWeek := time.Now().ISOWeek()
	return thisWeek
}
