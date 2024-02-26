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
	accountPublicKey := g.Config.Shared.NATSAccountPublicKey

	accountSigningKey, err := g.createAccountSigningKey(g.Config.Shared.NATSAccountSeed)
	if err != nil {
		return "", fmt.Errorf("create account signing key: %w", err)
	}

	userPublicKey, err := g.createUserPublicKey()
	if err != nil {
		return "", fmt.Errorf("create user public key: %w", err)
	}

	userJWT, err := g.createUserJWT(userPublicKey, accountPublicKey, accountSigningKey)
	if err != nil {
		return "", fmt.Errorf("create user jwt: %w", err)
	}

	return userJWT, nil
}

func (g *NATSJWT) createAccountSigningKey(seed string) (nkeys.KeyPair, error) {
	kp, err := nkeys.ParseDecoratedNKey([]byte(seed))
	if err != nil {
		return nil, err
	}
	return kp, nil
}

func (g *NATSJWT) createUserPublicKey() (string, error) {
	kp, err := nkeys.CreateUser()
	if err != nil {
		return "", err
	}
	publicKey, err := kp.PublicKey()
	if err != nil {
		return "", err
	}
	return publicKey, nil
}

func (g *NATSJWT) createUserJWT(userPublicKey, accountPublicKey string, accountSigningKey nkeys.KeyPair) (string, error) {
	uc := jwt.NewUserClaims(userPublicKey)

	uc.IssuerAccount = accountPublicKey
	uc.Expires = time.Now().Add(time.Hour).Unix()
	uc.BearerToken = true
	uc.NatsLimits.Subs = 1000
	uc.NatsLimits.Payload = 1_000_000

	if g.UserClaimsFunc != nil {
		g.UserClaimsFunc(uc)
	}

	vr := jwt.ValidationResults{}
	uc.Validate(&vr)
	if vr.IsBlocking(true) {
		return "", errors.New("user claims are invalid")
	}

	userJWT, err := uc.Encode(accountSigningKey)
	if err != nil {
		return "", err
	}

	return userJWT, nil
}
