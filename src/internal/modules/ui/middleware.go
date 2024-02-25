package ui

import (
	"log/slog"
	"net/http"
	"os"
	"tasks-app/internal/shared"
)

type Middleware func(http.Handler) http.Handler

func HandleWithMiddleware(mux *http.ServeMux, pattern string, handler http.Handler, middlewares ...Middleware) {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	mux.Handle(pattern, handler)
}

func ErrorRecoveryMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					if err == http.ErrAbortHandler {
						panic(err)
					}
					logger.Error("error recovery", "panic", err)
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func UserContextMiddleware(auth *Auth) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(shared.WithUserContext(r.Context(), auth.GetUserContext(r)))
			next.ServeHTTP(w, r)
		})
	}
}

func NATSMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			jwt, err := os.ReadFile("/ui_user.jwt")
			if err == nil {
				cookie := &http.Cookie{
					Name:     "nats_jwt",
					Value:    string(jwt),
					Path:     "/",
					Secure:   true,
					HttpOnly: true,
					SameSite: http.SameSiteStrictMode,
				}
				http.SetCookie(w, cookie)
			}
			next.ServeHTTP(w, r)
		})
	}
}
