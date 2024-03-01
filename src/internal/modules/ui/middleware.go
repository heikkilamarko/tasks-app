package ui

import (
	"fmt"
	"log/slog"
	"net/http"
	"tasks-app/internal/shared"

	"github.com/nats-io/jwt/v2"
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

func LoginMiddleware(config *shared.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := r.Cookie(config.UI.HubJWTCookieName); err == nil {
				next.ServeHTTP(w, r)
				return
			}

			user, err := shared.GetUserContext(r.Context())
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			natsJWT := &shared.NATSJWT{
				Config: config,
				UserClaimsFunc: func(c *jwt.UserClaims) {
					c.Sub.Allow.Add(fmt.Sprintf("task.%s.>", user.ID))
				}}

			jwt, err := natsJWT.CreateUserJWT()
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     config.UI.HubJWTCookieName,
				Value:    string(jwt),
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			})

			next.ServeHTTP(w, r)
		})
	}
}

func LogoutMiddleware(config *shared.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.SetCookie(w, &http.Cookie{
				Name:     config.UI.HubJWTCookieName,
				Value:    "",
				MaxAge:   -1,
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			})

			next.ServeHTTP(w, r)
		})
	}
}
