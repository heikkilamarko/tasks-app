package ui

import (
	"net/http"
	"slices"
	"time"
)

const (
	TimezoneDefault = "Europe/Helsinki"
	TimeFormat      = "02.01.2006 15.04"
	TimeFormatISO   = "2006-01-02T15:04"
)

const CookieNameTimezone = "timezone"

var SupportedTimezones = []string{
	"Europe/Helsinki",
	"Europe/London",
	"America/New_York",
}

func IsValidTimezone(timezone string) bool {
	return slices.Contains(SupportedTimezones, timezone)
}

func SetTimezoneCookie(w http.ResponseWriter, timezone string) {
	cookie := &http.Cookie{
		Name:     CookieNameTimezone,
		Value:    timezone,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   365 * 24 * 60 * 60, // One year in seconds
	}
	http.SetCookie(w, cookie)
}

func GetTimezone(r *http.Request) string {
	cookie, err := r.Cookie(CookieNameTimezone)
	if err != nil {
		return TimezoneDefault
	}

	if IsValidTimezone(cookie.Value) {
		return cookie.Value
	}

	return TimezoneDefault
}

func GetLocation(r *http.Request) *time.Location {
	l, _ := time.LoadLocation(GetTimezone(r))
	return l
}
