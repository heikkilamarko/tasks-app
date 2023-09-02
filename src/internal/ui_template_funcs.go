package internal

import (
	"time"
)

const (
	uiTimezone          = "Europe/Helsinki"
	uiTimeFormat        = "2006-01-02T15:04"
	uiDisplayTimeFormat = "02.01.2006 15.04"
)

func FormatUITime(from time.Time) string {
	l, err := time.LoadLocation(uiTimezone)
	if err != nil {
		return ""
	}

	return from.In(l).Format(uiTimeFormat)
}

func FormatUIDisplayTime(from time.Time) string {
	l, err := time.LoadLocation(uiTimezone)
	if err != nil {
		return ""
	}

	return from.In(l).Format(uiDisplayTimeFormat)
}

func ParseUITime(t string) (time.Time, error) {
	var empty time.Time

	l, err := time.LoadLocation(uiTimezone)
	if err != nil {
		return empty, err
	}

	pt, err := time.ParseInLocation(uiTimeFormat, t, l)
	if err != nil {
		return empty, err
	}

	return pt.UTC(), nil
}
