package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"go-home/os"
	"go-home/work"
)

type Runner struct {
	OS      os.OS
	Working work.Working
}

func (r Runner) Run() (err error) {
	for {
		isWorkingDay := true
		// isWorkingDay, err := r.IsWorkingDay()
		// if err != nil {
		// 	return err
		// }

		slog.Info(fmt.Sprintf("%s is working day", time.Now().Format(time.DateOnly)))

		if isWorkingDay {
			err = r.StartOneDay()
			if err != nil {
				return err
			}
		}

		d := time.Until(NextDay(time.Now()))
		slog.Info(fmt.Sprintf("等待 %s 后开始统计下一天", d))
		time.Sleep(d)
	}
}

func (r Runner) StartOneDay() error {
	// get unlock time
	unlockTime, err := r.GetUnlockTime()
	if err != nil {
		if errors.Is(err, os.ErrNotFoundUnlockTime) {
			slog.Warn("今天上班期间没找到解锁时间，等待重新执行")
			return nil
		}
		return err
	}

	slog.Info(fmt.Sprintf("开始执行，解锁时间为: %s", unlockTime.Format(time.TimeOnly)))

	endTime := r.Working.EndTime(unlockTime)
	slog.Info(fmt.Sprintf("预计下班时间为: %s", endTime.Format(time.TimeOnly)))

	slog.Info(fmt.Sprintf("剩余时间为: %s", time.Until(endTime)))
	time.Sleep(time.Until(endTime))
	return r.OS.Notify(unlockTime, endTime)
}

func (r Runner) GetUnlockTime() (time.Time, error) {
	// 没到上班时间等待到上班时间
	if !time.Now().Before(r.Working.StartTime()) {
		time.Sleep(time.Until(r.Working.StartTime()))
	}

	endTime := r.Working.EndTime(r.Working.StartTime())
	for {
		unlockTime, err := r.OS.GetUnlockTime()
		if err != nil {
			if errors.Is(err, os.ErrNotFoundUnlockTime) {
				// 到下班时间后不再寻找
				if time.Now().After(endTime) {
					return time.Time{}, os.ErrNotFoundUnlockTime
				}
				time.Sleep(5 * time.Minute)
				continue
			}
			return time.Time{}, err
		}
		return unlockTime, nil
	}
}

func (r Runner) IsWorkingDay() (bool, error) {
	year := time.Now().Year()
	resp, err := http.Get(fmt.Sprintf("https://api.jiejiariapi.com/v1/holidays/%d", year))
	if err != nil {
		return false, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	var holidays map[string]struct {
		Date     string `json:"date"`
		Name     string `json:"name"`
		IsOffDay bool   `json:"isOffDay"`
	}

	err = json.Unmarshal(body, &holidays)
	if err != nil {
		return false, err
	}

	if holiday, ok := holidays[time.Now().Format(time.DateOnly)]; ok {
		return !holiday.IsOffDay, nil
	}

	weekDay := time.Now().Weekday()
	return weekDay != time.Saturday && weekDay != time.Sunday, nil
}
