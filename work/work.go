package work

import (
	"time"
)

type Working interface {
	StartTime() time.Time
	EndTime(startTime time.Time) time.Time
}

type work struct {
	startTime string
	endTime   string
}

func NewWork(startTime, endTime string) (Working, error) {
	_, err := time.Parse("15:04", startTime)
	if err != nil {
		return nil, err
	}

	_, err = time.Parse("15:04", endTime)
	if err != nil {
		return nil, err
	}

	return &work{
		startTime: startTime,
		endTime:   endTime,
	}, nil
}

func (w *work) parseTime(t string) time.Time {
	s, _ := time.Parse("15:04", t)
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), s.Hour(), s.Minute(), 0, 0, now.Location())
}

func (w *work) StartTime() time.Time {
	return w.parseTime(w.startTime)
}

func (w *work) EndTime(startTime time.Time) time.Time {
	var d time.Duration
	if startTime.After(w.StartTime()) {
		d = startTime.Sub(w.StartTime())
	}

	e := w.parseTime(w.endTime)
	return e.Add(d)
}
