package ui

import (
	"net/http"
	"slices"
)

const (
	ThemeDark    = "dark"
	ThemeLight   = "light"
	ThemeDefault = ThemeDark
)

const CookieNameTheme = "theme"

var SupportedThemes = []string{
	ThemeDark,
	ThemeLight,
}

func IsValidTheme(theme string) bool {
	return slices.Contains(SupportedThemes, theme)
}

func SetThemeCookie(w http.ResponseWriter, theme string) {
	cookie := &http.Cookie{
		Name:     CookieNameTheme,
		Value:    theme,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   365 * 24 * 60 * 60, // One year in seconds
	}
	http.SetCookie(w, cookie)
}

func GetTheme(r *http.Request) string {
	cookie, err := r.Cookie(CookieNameTheme)
	if err != nil {
		return ThemeDefault
	}

	if IsValidTheme(cookie.Value) {
		return cookie.Value
	}

	return ThemeDefault
}
