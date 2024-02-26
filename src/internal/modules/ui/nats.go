package ui

import (
	"errors"
	"fmt"
	"tasks-app/internal/shared"
	"time"

	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nkeys"
)

func GenerateUserJWT(userID string, config *shared.Config) (string, error) {
	accountPublicKey := config.Shared.NATSAccountPublicKey

	accountSigningKey, err := getAccountSigningKey(config.Shared.NATSAccountSeed)
	if err != nil {
		return "", fmt.Errorf("get account signing key: %w", err)
	}

	userPublicKey, err := generateUserKey()
	if err != nil {
		return "", fmt.Errorf("generate user key: %w", err)
	}

	userJWT, err := generateUserJWT(userID, userPublicKey, accountPublicKey, accountSigningKey)
	if err != nil {
		return "", fmt.Errorf("generate user jwt: %w", err)
	}

	return userJWT, nil
}

func getAccountSigningKey(seed string) (nkeys.KeyPair, error) {
	kp, err := nkeys.ParseDecoratedNKey([]byte(seed))
	if err != nil {
		return nil, err
	}
	return kp, nil
}

func generateUserKey() (string, error) {
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

func generateUserJWT(userID string, userPublicKey, accountPublicKey string, accountSigningKey nkeys.KeyPair) (string, error) {
	uc := jwt.NewUserClaims(userPublicKey)

	uc.IssuerAccount = accountPublicKey
	uc.Expires = time.Now().Add(time.Hour).Unix()
	uc.BearerToken = true
	uc.NatsLimits.Subs = 1000
	uc.NatsLimits.Payload = 1_000_000
	uc.Sub.Allow.Add(fmt.Sprintf("tasks.ui.%s.>", userID))

	vr := jwt.ValidationResults{}
	uc.Validate(&vr)
	if vr.IsBlocking(true) {
		return "", errors.New("generated user user claims are invalid")
	}

	userJWT, err := uc.Encode(accountSigningKey)
	if err != nil {
		return "", err
	}

	return userJWT, nil
}
