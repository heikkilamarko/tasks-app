package ui

import (
	"strconv"
	"time"
)

func ValidateID(value string) (int, bool) {
	v, err := strconv.Atoi(value)
	return v, err == nil
}

func ValidateName(value string) (string, bool) {
	if value == "" {
		return value, false
	}

	return value, true
}

func ValidateExpiresAt(value string) (*time.Time, bool) {
	if value == "" {
		return nil, true
	}

	v, err := ParseTime(value)
	return v, err == nil
}
