package timer

import (
	"log/slog"
	"time"
)

func Track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func Duration(msg string, start time.Time) {
	slog.Info(msg, "time", time.Since(start))
}
