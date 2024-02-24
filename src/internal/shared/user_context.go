package shared

import (
	"context"
	"errors"
)

type userCtxKeyType string

const userCtxKey userCtxKeyType = "user"

var ErrUserContextNotFound = errors.New("user context not found")

type UserContext struct {
	ID          string
	Name        string
	Email       string
	IDToken     string
	AccessToken string
}

func WithUserContext(ctx context.Context, user *UserContext) context.Context {
	return context.WithValue(ctx, userCtxKey, user)
}

func GetUserContext(ctx context.Context) (*UserContext, error) {
	user, ok := ctx.Value(userCtxKey).(*UserContext)
	if !ok {
		return nil, ErrUserContextNotFound
	}
	return user, nil
}
