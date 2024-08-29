package os

import (
	"os"
	"time"
)

func init() {
	// 设置环境变量 TZ 为东八区（Asia/Shanghai）
	os.Setenv("TZ", "Asia/Shanghai")
}

type OS interface {
	GetUnlockTime() (time.Time, error)
	Notify(unlockTime, reminderTime time.Time) error
}

var RunnerOS OS
