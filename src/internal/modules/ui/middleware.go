package ui

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const (
	cookieMaxAge    = 60 * 60 // 1h
	cookieSessionID = "tasks_app_session_id"
)

type ctxKeySessionID struct{}

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sessionID string
		c, err := r.Cookie(cookieSessionID)
		if err == http.ErrNoCookie {
			u, _ := uuid.NewRandom()
			sessionID = u.String()
			http.SetCookie(w, &http.Cookie{
				Name:     cookieSessionID,
				Value:    sessionID,
				MaxAge:   cookieMaxAge,
				HttpOnly: true,
			})
		} else if err != nil {
			return
		} else {
			sessionID = c.Value
		}
		ctx := context.WithValue(r.Context(), ctxKeySessionID{}, sessionID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func SessionID(r *http.Request) string {
	v := r.Context().Value(ctxKeySessionID{})
	if v != nil {
		return v.(string)
	}
	return ""
}
