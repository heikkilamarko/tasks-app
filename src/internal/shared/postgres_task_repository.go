package shared

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	// PostgreSQL driver
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresTaskRepositoryOptions struct {
	ConnectionString string
	Logger           *slog.Logger
}

type PostgresTaskRepository struct {
	Options PostgresTaskRepositoryOptions
	db      *sql.DB
}

func NewPostgresTaskRepository(ctx context.Context, options PostgresTaskRepositoryOptions) (*PostgresTaskRepository, error) {
	db, err := sql.Open("pgx", options.ConnectionString)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(10 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return &PostgresTaskRepository{options, db}, nil
}

func (repo *PostgresTaskRepository) Close() error {
	return repo.db.Close()
}

func (repo *PostgresTaskRepository) Create(ctx context.Context, task *Task) error {
	query := `
		INSERT INTO task (name, expires_at, expiring_info_at, expired_info_at, created_at, updated_at, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	err := repo.db.QueryRowContext(
		ctx,
		query,
		task.Name, task.ExpiresAt, task.ExpiringInfoAt, task.ExpiredInfoAt, task.CreatedAt, task.UpdatedAt, task.CompletedAt,
	).Scan(&task.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *PostgresTaskRepository) Update(ctx context.Context, task *Task) error {
	query := `
		UPDATE task
		SET name = $2, expires_at = $3, expiring_info_at = $4, expired_info_at = $5, updated_at = $6, completed_at = $7
		WHERE id = $1
	`

	_, err := repo.db.ExecContext(
		ctx,
		query,
		task.ID, task.Name, task.ExpiresAt, task.ExpiringInfoAt, task.ExpiredInfoAt, task.UpdatedAt, task.CompletedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *PostgresTaskRepository) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM task
		WHERE id = $1
	`

	_, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (repo *PostgresTaskRepository) GetByID(ctx context.Context, id int) (*Task, error) {
	query := `
		SELECT id, name, expires_at, expiring_info_at, expired_info_at, created_at, updated_at, completed_at
		FROM task
		WHERE id = $1
	`

	row := repo.db.QueryRowContext(ctx, query, id)

	task := &Task{}
	err := row.Scan(
		&task.ID, &task.Name, &task.ExpiresAt, &task.ExpiringInfoAt, &task.ExpiredInfoAt, &task.CreatedAt, &task.UpdatedAt, &task.CompletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return task, nil
}

func (repo *PostgresTaskRepository) GetActive(ctx context.Context) ([]*Task, error) {
	query := `
		SELECT id, name, expires_at, expiring_info_at, expired_info_at, created_at, updated_at, completed_at
		FROM task
		WHERE completed_at IS NULL
		ORDER BY created_at DESC
	`

	return repo.getTasks(ctx, query)
}

func (repo *PostgresTaskRepository) GetCompleted(ctx context.Context) ([]*Task, error) {
	query := `
		SELECT id, name, expires_at, expiring_info_at, expired_info_at, created_at, updated_at, completed_at
		FROM task
		WHERE completed_at IS NOT NULL
		ORDER BY completed_at DESC
	`

	return repo.getTasks(ctx, query)
}

func (repo *PostgresTaskRepository) GetExpiring(ctx context.Context, d time.Duration) ([]*Task, error) {
	query := `
		SELECT id, name, expires_at, expiring_info_at, expired_info_at, created_at, updated_at, completed_at
		FROM task
		WHERE completed_at IS NULL
		AND expires_at IS NOT NULL
		AND expiring_info_at IS NULL
		AND expires_at >= $1
		AND expires_at <= $2
		ORDER BY created_at ASC
	`

	t1 := time.Now().UTC()
	t2 := t1.Add(d)

	return repo.getTasks(ctx, query, t1, t2)
}

func (repo *PostgresTaskRepository) GetExpired(ctx context.Context) ([]*Task, error) {
	query := `
		SELECT id, name, expires_at, expiring_info_at, expired_info_at, created_at, updated_at, completed_at
		FROM task
		WHERE completed_at IS NULL
		AND expires_at IS NOT NULL
		AND expired_info_at IS NULL
		AND expires_at < $1
		ORDER BY created_at ASC
	`

	now := time.Now().UTC()

	return repo.getTasks(ctx, query, now)
}

func (repo *PostgresTaskRepository) DeleteCompleted(ctx context.Context, d time.Duration) (int64, error) {
	query := `
		DELETE FROM task
		WHERE completed_at IS NOT NULL
		AND completed_at < $1
	`

	t := time.Now().UTC().Add(-d)

	result, err := repo.db.ExecContext(ctx, query, t)
	if err != nil {
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *PostgresTaskRepository) getTasks(ctx context.Context, query string, args ...any) ([]*Task, error) {
	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}

		err := rows.Scan(
			&task.ID, &task.Name, &task.ExpiresAt, &task.ExpiringInfoAt, &task.ExpiredInfoAt, &task.CreatedAt, &task.UpdatedAt, &task.CompletedAt,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
