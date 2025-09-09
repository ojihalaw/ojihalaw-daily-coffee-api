package utils

import (
	"time"
)

// TodayString -> return string "YYYYMMDD"
func TodayString() string {
	return time.Now().Format("20060102")
}
