package ui

import (
	"context"
	"net/http"
	"tasks-app/internal/shared"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/authentication"
	zoidc "github.com/zitadel/zitadel-go/v3/pkg/authentication/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

type Auth struct {
	Authenticator *authentication.Authenticator[*zoidc.UserInfoContext[*oidc.IDTokenClaims, *oidc.UserInfo]]
	Middleware    *authentication.Interceptor[*zoidc.UserInfoContext[*oidc.IDTokenClaims, *oidc.UserInfo]]
	Config        *shared.Config
}

func NewAuth(ctx context.Context, config *shared.Config) (*Auth, error) {
	authenticator, err := authentication.New(
		ctx,
		zitadel.New(config.UI.AuthDomain),
		config.UI.AuthEncryptionKey,
		zoidc.DefaultAuthentication(config.UI.AuthClientId, config.UI.AuthRedirectURI, config.UI.AuthEncryptionKey),
	)
	if err != nil {
		return nil, err
	}

	middleware := authentication.Middleware(authenticator)

	return &Auth{authenticator, middleware, config}, nil
}

func (a *Auth) GetUserContext(r *http.Request) *shared.UserContext {
	if ctx := a.Middleware.Context(r.Context()); ctx != nil {
		return &shared.UserContext{
			ID:          ctx.UserInfo.Subject,
			Name:        ctx.UserInfo.Name,
			Email:       ctx.UserInfo.Email,
			IDToken:     ctx.Tokens.IDToken,
			AccessToken: ctx.Tokens.AccessToken,
		}
	}
	return nil
}

func (a *Auth) LoginHandler(requestedURI string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.Authenticator.Authenticate(w, r, requestedURI)
	})
}

func (a *Auth) CallbackHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.DeleteHubJWTCookie(w)
		a.Authenticator.Callback(w, r)
	})
}

func (a *Auth) LogoutHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.DeleteHubJWTCookie(w)
		a.Authenticator.Logout(w, r)
	})
}

func (a *Auth) IsHubJWTCookieSet(r *http.Request) bool {
	_, err := r.Cookie(a.Config.UI.HubJWTCookieName)
	return err == nil
}

func (a *Auth) SetHubJWTCookie(w http.ResponseWriter, jwt string) {
	http.SetCookie(w, &http.Cookie{
		Name:     a.Config.UI.HubJWTCookieName,
		Value:    jwt,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

func (a *Auth) DeleteHubJWTCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     a.Config.UI.HubJWTCookieName,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}
