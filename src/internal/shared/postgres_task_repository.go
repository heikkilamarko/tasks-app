package shared

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresTaskRepository struct {
	db *sql.DB
}

func NewPostgresTaskRepository(ctx context.Context, config *Config) (*PostgresTaskRepository, error) {
	db, err := sql.Open("pgx", config.Shared.PostgresConnectionString)
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

	return &PostgresTaskRepository{db}, nil
}

func (repo *PostgresTaskRepository) Close() error {
	return repo.db.Close()
}

func (repo *PostgresTaskRepository) Create(ctx context.Context, task *Task) error {
	query := `
		INSERT INTO task
			(name, expires_at, expiring_info_at, expired_info_at, created_at, updated_at, completed_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	return repo.db.QueryRowContext(
		ctx,
		query,
		task.Name, task.ExpiresAt, task.ExpiringInfoAt, task.ExpiredInfoAt, task.CreatedAt, task.UpdatedAt, task.CompletedAt,
	).Scan(&task.ID)
}

func (repo *PostgresTaskRepository) Update(ctx context.Context, task *Task) error {
	query := `
		UPDATE task
		SET 
			name = $2,
			expires_at = $3,
			expiring_info_at = $4,
			expired_info_at = $5,
			updated_at = $6,
			completed_at = $7
		WHERE
			id = $1
	`

	_, err := repo.db.ExecContext(ctx, query, task.ID, task.Name, task.ExpiresAt, task.ExpiringInfoAt, task.ExpiredInfoAt, task.UpdatedAt, task.CompletedAt)
	return err
}

func (repo *PostgresTaskRepository) UpdateAttachments(ctx context.Context, taskID int, inserted []string, deleted map[int]string) error {
	now := time.Now().UTC()

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO attachment
			(task_id, file_name, created_at)
		VALUES
			($1, $2, $3)
	`

	for _, name := range inserted {
		if _, err := tx.ExecContext(ctx, query, taskID, name, now); err != nil {
			return err
		}
	}

	query = `
		DELETE FROM attachment
		WHERE id = $1
	`

	for id := range deleted {
		if _, err := tx.ExecContext(ctx, query, id); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (repo *PostgresTaskRepository) Delete(ctx context.Context, id int) error {
	query := `
		DELETE FROM task
		WHERE id = $1
	`

	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *PostgresTaskRepository) GetByID(ctx context.Context, id int) (*Task, error) {
	query := `
		WHERE t.id = $1
	`

	tasks, err := repo.getTasks(ctx, query, id)
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, nil
	}

	return tasks[0], nil
}

func (repo *PostgresTaskRepository) GetActive(ctx context.Context, offset int, limit int) ([]*Task, error) {
	query := `
		WHERE t.completed_at IS NULL
		ORDER BY t.created_at DESC
		LIMIT $1 OFFSET $2
	`

	return repo.getTasks(ctx, query, limit, offset)
}

func (repo *PostgresTaskRepository) GetCompleted(ctx context.Context, offset int, limit int) ([]*Task, error) {
	query := `
		WHERE t.completed_at IS NOT NULL
		ORDER BY t.completed_at DESC
		LIMIT $1 OFFSET $2
	`

	return repo.getTasks(ctx, query, limit, offset)
}

func (repo *PostgresTaskRepository) GetExpiring(ctx context.Context, d time.Duration) ([]*Task, error) {
	query := `
		WHERE t.completed_at IS NULL
		AND t.expiring_info_at IS NULL
		AND t.expires_at IS NOT NULL
		AND t.expires_at >= $1
		AND t.expires_at <= $2
		ORDER BY t.created_at ASC
	`

	t1 := time.Now().UTC()
	t2 := t1.Add(d)

	return repo.getTasks(ctx, query, t1, t2)
}

func (repo *PostgresTaskRepository) GetExpired(ctx context.Context) ([]*Task, error) {
	query := `
		WHERE t.completed_at IS NULL
		AND t.expired_info_at IS NULL
		AND t.expires_at IS NOT NULL
		AND t.expires_at < $1
		ORDER BY t.created_at ASC
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
	var tasks []*Task

	query = `
		SELECT
			t.id,
			t.name,
			t.expires_at,
			t.expiring_info_at,
			t.expired_info_at,
			t.created_at,
			t.updated_at,
			t.completed_at,
			a.id,
			a.task_id,
			a.file_name,
			a.created_at,
			a.updated_at
		FROM
			task t
		LEFT JOIN
			attachment a ON t.id = a.task_id
	` + query

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasksMap := make(map[int]*Task)

	for rows.Next() {
		var t Task
		var a struct {
			ID        *int
			TaskID    *int
			FileName  *string
			CreatedAt *time.Time
			UpdatedAt *time.Time
		}

		if err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.ExpiresAt,
			&t.ExpiringInfoAt,
			&t.ExpiredInfoAt,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.CompletedAt,
			&a.ID,
			&a.TaskID,
			&a.FileName,
			&a.CreatedAt,
			&a.UpdatedAt,
		); err != nil {
			return nil, err
		}

		task, ok := tasksMap[t.ID]
		if !ok {
			task = &t
			tasks = append(tasks, task)
			tasksMap[task.ID] = task
		}

		if a.ID != nil {
			task.Attachments = append(task.Attachments, &Attachment{*a.ID, *a.TaskID, *a.FileName, *a.CreatedAt, a.UpdatedAt})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
