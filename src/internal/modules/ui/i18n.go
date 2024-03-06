package ui

import (
	"net/http"
	"slices"
)

const (
	LanguageDefault    = "en"
	CookieNameLanguage = "language"
)

var SupportedLanguages = []string{
	"en",
	"fi",
}

func IsValidLanguage(lang string) bool {
	return slices.Contains(SupportedLanguages, lang)
}

func SetLanguageCookie(w http.ResponseWriter, lang string) {
	cookie := &http.Cookie{
		Name:     CookieNameLanguage,
		Value:    lang,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   365 * 24 * 60 * 60, // One year in seconds
	}
	http.SetCookie(w, cookie)
}

func GetLanguage(r *http.Request) string {
	cookie, err := r.Cookie(CookieNameLanguage)
	if err != nil {
		return LanguageDefault
	}

	if IsValidLanguage(cookie.Value) {
		return cookie.Value
	}

	return LanguageDefault
}
