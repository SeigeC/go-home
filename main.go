package main

import (
	"log/slog"
	"time"

	"go-home/os"
)

type Working interface {
	StartTime() time.Time
	EndTime(startTime time.Time) time.Time
}

var runner Runner

func main() {
	runner := Runner{
		OS: os.RunnerOS,
	}

	unlockTime, err := runner.OS.GetUnlockTime()
	if err != nil {
		slog.Error("获取解锁时间失败：", slog.Any("err", err))
		return
	}

	slog.Info("解锁时间", slog.String("time", unlockTime.Format(time.DateTime)))

	// 计算 9 小时后的时间
	reminderTime := unlockTime.Add(9 * time.Hour)
	slog.Info("提醒时间", slog.String("time", reminderTime.Format(time.DateTime)))

	<-time.After(time.Until(reminderTime))

	// 调用系统通知
	err = runner.OS.Notify(unlockTime, reminderTime)
	if err != nil {
		slog.Error("Error setting reminder:", slog.Any("err", err))
		return
	}
}
