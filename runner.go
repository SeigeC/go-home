package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"go-home/os"
)

type Runner struct {
	OS      os.OS
	Working Working
}

type runtime struct {
	StartTime time.Time
	EndTime   time.Time
}

func newRuntime() *runtime {
	return &runtime{}
}

func (r Runner) Run() (err error) {
	var runtime *runtime

start:
	for {
		if runtime == nil {
			isWorkingDay, err := r.IsWorkingDay()
			if err != nil {
				return err
			}

			if !isWorkingDay {
				time.Sleep(time.Until(NextDay(time.Now())))
				goto start
			}

			runtime = newRuntime()
		}

		unlockTime, err := r.GetUnlockTime()
		if err != nil {
			if errors.Is(err, os.ErrNotFoundUnlockTime) {
				runtime = nil
				time.Sleep(time.Until(NextDay(time.Now())))
				goto start
			}
		}

		runtime.StartTime = unlockTime
		runtime.EndTime = r.Working.EndTime(runtime.StartTime)
	}
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
