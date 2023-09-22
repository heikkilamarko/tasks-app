package ui

import "time"

const (
	timezone      = "Europe/Helsinki"
	timeFormat    = "02.01.2006 15.04"
	timeFormatISO = "2006-01-02T15:04"
)

func ParseTime(t string) (*time.Time, error) {
	l, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, err
	}

	pt, err := time.ParseInLocation(timeFormatISO, t, l)
	if err != nil {
		return nil, err
	}

	pt = pt.UTC()

	return &pt, nil
}
