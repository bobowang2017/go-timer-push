package utils

import (
	"go-timer-push/logger"
	"time"
)

//对于协程内部运行的函数，如果发生panic会导致整个程序崩溃，故需要手动recover
func SafeGo(do func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logger.Logger.Error(err)
			}
		}()
		do()
	}()
}

//计算当前日期是第几周
func WeekCountNow() int {
	_, thisWeek := time.Now().ISOWeek()
	return thisWeek
}
