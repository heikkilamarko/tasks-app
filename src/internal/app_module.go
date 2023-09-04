package internal

import "context"

type AppModule interface {
	Run(ctx context.Context) error
	Close() error
}
