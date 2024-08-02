package ui

import (
	"net/http"
	"net/url"
	"strings"
)

const DefaultRedirectURL = "/ui"

func GetRedirectURL(r *http.Request) string {
	currentURL := r.Header.Get("HX-Current-URL")
	if currentURL == "" {
		currentURL = r.Header.Get("Referer")
		if currentURL == "" {
			return DefaultRedirectURL
		}
	}

	redirectURL, err := url.Parse(currentURL)
	if err != nil {
		return DefaultRedirectURL
	}

	if !strings.HasPrefix(redirectURL.Path, DefaultRedirectURL) {
		return DefaultRedirectURL
	}

	return redirectURL.Path
}
