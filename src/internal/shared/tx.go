package shared

import "context"

type TxManager interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}
