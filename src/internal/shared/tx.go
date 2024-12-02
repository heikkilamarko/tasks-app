package shared

import (
	"context"
	"database/sql"
	"errors"
)

type DB interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type TxAdapters struct {
	TaskRepository TaskRepository
}

type TxProvider interface {
	Transact(txFunc func(adapters TxAdapters) error) error
}

type SQLTxProvider struct {
	db *sql.DB
}

func NewSQLTxProvider(db *sql.DB) *SQLTxProvider {
	return &SQLTxProvider{db}
}

func (p *SQLTxProvider) Transact(txFunc func(adapters TxAdapters) error) error {
	return runInTx(p.db, func(tx *sql.Tx) error {
		adapters := TxAdapters{
			TaskRepository: NewPostgresTaskRepository(tx),
		}

		return txFunc(adapters)
	})
}

func runInTx(db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err == nil {
		return tx.Commit()
	}

	rollbackErr := tx.Rollback()
	if rollbackErr != nil {
		return errors.Join(err, rollbackErr)
	}

	return err
}
