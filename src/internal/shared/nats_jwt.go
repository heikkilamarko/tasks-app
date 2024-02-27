package shared

import (
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nkeys"
)

type NATSJWT struct {
	Config         *Config
	UserClaimsFunc func(c *jwt.UserClaims)
}

func (g *NATSJWT) CreateUserJWT() (string, error) {
	accountKP, err := nkeys.FromSeed([]byte(g.Config.Shared.NATSAccountSeed))
	if err != nil {
		return "", fmt.Errorf("get account key pair: %w", err)
	}

	accountPub := g.Config.Shared.NATSAccountPublicKey

	userKP, err := nkeys.CreateUser()
	if err != nil {
		return "", fmt.Errorf("create user key pair: %w", err)
	}

	userPub, err := userKP.PublicKey()
	if err != nil {
		return "", fmt.Errorf("get user public key: %w", err)
	}

	userClaims := jwt.NewUserClaims(userPub)
	userClaims.IssuerAccount = accountPub
	userClaims.Expires = time.Now().Add(time.Hour).Unix()
	userClaims.BearerToken = true
	userClaims.NatsLimits.Subs = 1000
	userClaims.NatsLimits.Payload = 1_000_000
	if g.UserClaimsFunc != nil {
		g.UserClaimsFunc(userClaims)
	}

	vr := jwt.ValidationResults{}
	userClaims.Validate(&vr)
	if vr.IsBlocking(true) {
		return "", errors.New("validate user claims")
	}

	userJWT, err := userClaims.Encode(accountKP)
	if err != nil {
		return "", fmt.Errorf("encode user claims: %w", err)
	}

	return userJWT, nil
}
