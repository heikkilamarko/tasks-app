package ui

import (
	"context"
	"crypto/tls"
	"net/http"
	"tasks-app/internal/shared"

	"github.com/go-chi/chi/v5"
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

type AuthInfo struct {
	UserName    string
	UserEmail   string
	AccessToken string
}

func NewAuth(ctx context.Context, config *shared.Config) (*Auth, error) {
	zhttp.DefaultHTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

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

func (a *Auth) RegisterRoutes(router chi.Router) {
	router.Handle(a.Config.UI.AuthPath+"/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.Authenticator.Authenticate(w, r, "/ui")
	}))

	router.Handle(a.Config.UI.AuthPath+"/callback", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.Authenticator.Callback(w, r)
	}))

	router.Handle(a.Config.UI.AuthPath+"/logout", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.Authenticator.Logout(w, r)
	}))
}

func (a *Auth) GetAuthInfo(r *http.Request) *AuthInfo {
	if ctx := a.Middleware.Context(r.Context()); ctx != nil {
		return &AuthInfo{
			ctx.UserInfo.Name,
			ctx.UserInfo.Email,
			ctx.Tokens.AccessToken,
		}
	}
	return nil
}
