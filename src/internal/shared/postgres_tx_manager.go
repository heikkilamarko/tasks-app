package shared

import "database/sql"

type PostgresTxManager struct {
	db *sql.DB
}

var _ TxManager = (*PostgresTxManager)(nil)

func NewPostgresTxManager(db *sql.DB) *PostgresTxManager {
	return &PostgresTxManager{db}
}

func (m *PostgresTxManager) RunInTx(fn func(txc TxContext) error) error {
	return runInTx(m.db, func(tx *sql.Tx) error {
		return fn(TxContext{
			TaskRepository: NewPostgresTaskRepository(tx),
		})
	})
}
