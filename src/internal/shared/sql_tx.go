package shared

import (
	"context"
	"database/sql"
)

type sqlTxCtxKeyType string

const sqlTxCtxKey sqlTxCtxKeyType = "sqltx"

type SQLTx interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

type SQLTxManager struct {
	db *sql.DB
}

var _ TxManager = (*SQLTxManager)(nil)

func NewSQLTxManager(db *sql.DB) *SQLTxManager {
	return &SQLTxManager{db: db}
}

func (m *SQLTxManager) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fn(context.WithValue(ctx, sqlTxCtxKey, tx)); err != nil {
		return err
	}

	return tx.Commit()
}

func GetSQLTx(ctx context.Context) *sql.Tx {
	tx, _ := ctx.Value(sqlTxCtxKey).(*sql.Tx)
	return tx
}
