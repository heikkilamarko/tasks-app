package ui

import (
	"errors"
	"strings"
	"time"
)

func Dict(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dictionary call")
	}

	root := make(map[string]any)

	for i := 0; i < len(values); i += 2 {
		dict := root
		var key string
		switch v := values[i].(type) {
		case string:
			key = v
		case []string:
			for i := 0; i < len(v)-1; i++ {
				key = v[i]
				var m map[string]any
				v, found := dict[key]
				if found {
					m = v.(map[string]any)
				} else {
					m = make(map[string]any)
					dict[key] = m
				}
				dict = m
			}
			key = v[len(v)-1]
		default:
			return nil, errors.New("invalid dictionary key")
		}
		dict[key] = values[i+1]
	}

	return root, nil
}

func FormatTime(l *time.Location, t time.Time) string {
	return t.In(l).Format(TimeFormat)
}

func FormatISOTime(l *time.Location, t time.Time) string {
	return t.In(l).Format(TimeFormatISO)
}

func FormatTimezone(tz string) string {
	parts := strings.Split(tz, "/")
	if len(parts) != 2 {
		return tz
	}
	return strings.ReplaceAll(parts[1], "_", " ")
}
