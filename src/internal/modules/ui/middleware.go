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

func LoginMiddleware(auth *Auth) func(next http.Handler) http.Handler {
	natsJWT := &shared.NATSJWT{Config: auth.Config}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if auth.IsHubJWTCookieSet(r) {
				next.ServeHTTP(w, r)
				return
			}

			user, err := shared.GetUserContext(r.Context())
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			allowSub := fmt.Sprintf("task.%s.>", user.ID)

			jwt, err := natsJWT.CreateUserJWT(func(c *jwt.UserClaims) {
				c.BearerToken = true
				c.Sub.Allow.Add(allowSub)
			})
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			auth.SetHubJWTCookie(w, string(jwt))

			next.ServeHTTP(w, r)
		})
	}
}
