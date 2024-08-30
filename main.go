package main

import (
	"log/slog"

	"go-home/os"
	"go-home/work"
)

func main() {
	worker, err := work.NewWork("10:00", "19:00")
	if err != nil {
		slog.Error("work error", "err", err)
		return
	}

	runner := Runner{
		OS:      os.RunnerOS,
		Working: worker,
	}

	err = runner.Run()
	if err != nil {
		slog.Error("runner error", "err", err)
	}
}
