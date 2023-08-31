package internal

import "context"

type MessagingClient interface {
	SendMsg(ctx context.Context, subject string, data any) error
	SendPersistentMsg(ctx context.Context, subject string, data any) error
}
