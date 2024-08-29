package main

import (
	"time"
)

// 第二天 0 点
func NextDay(now time.Time) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
}
