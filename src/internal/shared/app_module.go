package shared

import "context"

type AppModule interface {
	Name() string
	Run(ctx context.Context) error
	Close() error
}
