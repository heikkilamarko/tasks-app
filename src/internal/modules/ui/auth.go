package ui

import (
	"context"
	"net/http"
	"tasks-app/internal/shared"

	zhttp "github.com/zitadel/oidc/v3/pkg/http"
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
		zoidc.WithCodeFlow[*zoidc.UserInfoContext[*oidc.IDTokenClaims, *oidc.UserInfo], *oidc.IDTokenClaims, *oidc.UserInfo](
			zoidc.PKCEAuthentication(
				config.UI.AuthClientId,
				config.UI.AuthRedirectURI,
				[]string{
					oidc.ScopeOpenID,
					oidc.ScopeProfile,
					oidc.ScopeEmail,
				},
				zhttp.NewCookieHandler(
					[]byte(config.UI.AuthEncryptionKey),
					[]byte(config.UI.AuthEncryptionKey),
				),
			),
		),
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
		a.Authenticator.Callback(w, r)
	})
}

func (a *Auth) LogoutHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.Authenticator.Logout(w, r)
	})
}
