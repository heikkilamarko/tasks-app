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

type Timezone struct {
	Name     string
	Timezone string
}

var SupportedTimezones = []Timezone{
	{Name: "Helsinki", Timezone: "Europe/Helsinki"},
	{Name: "London", Timezone: "Europe/London"},
	{Name: "New York", Timezone: "America/New_York"},
}

func IsValidTimezone(timezone string) bool {
	return slices.ContainsFunc(SupportedTimezones, func(tz Timezone) bool {
		return tz.Timezone == timezone
	})
}

func SetTimezoneCookie(w http.ResponseWriter, timezone string) {
	cookie := &http.Cookie{
		Name:     CookieNameTimezone,
		Value:    timezone,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   365 * 24 * 60 * 60, // One year in seconds
	}
	http.SetCookie(w, cookie)
}

func GetLocation(r *http.Request) *time.Location {
	l, _ := time.LoadLocation(GetTimezone(r))
	return l
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
