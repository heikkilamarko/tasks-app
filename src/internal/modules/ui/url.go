package ui

import (
	"errors"
	"net/http"
	"net/url"
	"slices"
)

const DefaultRedirectURL = "/ui"

var ErrBadReferer = errors.New("referer invalid")

func GetRedirectURL(r *http.Request, trustedHosts []string) (string, error) {
	referer, err := url.Parse(r.Referer())
	if err != nil || referer.String() == "" {
		return "", ErrBadReferer
	}

	if !slices.Contains(trustedHosts, referer.Host) {
		return "", ErrBadReferer
	}

	return referer.Path, nil
}
