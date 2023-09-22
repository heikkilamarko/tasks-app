package ui

import (
	"os"
	"time"
)

func RenderEnv(key string) string {
	return os.Getenv(key)
}

func RenderTime(from time.Time) string {
	l, err := time.LoadLocation(timezone)
	if err != nil {
		return ""
	}

	return from.In(l).Format(timeFormat)
}

func RenderISOTime(from time.Time) string {
	l, err := time.LoadLocation(timezone)
	if err != nil {
		return ""
	}

	return from.In(l).Format(timeFormatISO)
}
