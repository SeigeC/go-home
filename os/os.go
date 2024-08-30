package os

import (
	"time"
)

type OS interface {
	GetUnlockTime() (time.Time, error)
	Notify(unlockTime, reminderTime time.Time) error
}

var RunnerOS OS
