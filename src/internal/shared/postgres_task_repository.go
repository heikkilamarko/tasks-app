package shared

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

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

func (repo *PostgresTaskRepository) UpdateAttachments(ctx context.Context, taskID int, inserted []string, deleted map[int]string) error {
	now := time.Now().UTC()

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, name := range inserted {
		_, err := tx.ExecContext(ctx,
			"INSERT INTO attachment (task_id, file_name, created_at) VALUES ($1, $2, $3)",
			taskID, name, now,
		)
		if err != nil {
			return err
		}
	}

	for id := range deleted {
		_, err := tx.ExecContext(ctx,
			"DELETE FROM attachment WHERE id = $1",
			id,
		)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
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

func (repo *PostgresTaskRepository) GetActive(ctx context.Context) ([]*Task, error) {
	query := `
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
		WHERE t.completed_at IS NULL
		ORDER BY t.created_at DESC
	`

	return repo.getTasks(ctx, query)
}

func (repo *PostgresTaskRepository) GetCompleted(ctx context.Context) ([]*Task, error) {
	query := `
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
		WHERE t.completed_at IS NOT NULL
		ORDER BY t.completed_at DESC
	`

	return repo.getTasks(ctx, query)
}

func (repo *PostgresTaskRepository) GetExpiring(ctx context.Context, d time.Duration) ([]*Task, error) {
	query := `
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
	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasksMap := make(map[int]*Task)
	var tasks []*Task

	for rows.Next() {
		task := &Task{}
		attachment := &EmptyAttachment{}

		err := rows.Scan(
			&task.ID,
			&task.Name,
			&task.ExpiresAt,
			&task.ExpiringInfoAt,
			&task.ExpiredInfoAt,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.CompletedAt,
			&attachment.ID,
			&attachment.TaskID,
			&attachment.FileName,
			&attachment.CreatedAt,
			&attachment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		existingTask, ok := tasksMap[task.ID]
		if !ok {
			task.Attachments = []*Attachment{}
			tasksMap[task.ID] = task
			tasks = append(tasks, task)
			existingTask = task
		}

		if attachment.ID != nil {
			existingTask.Attachments = append(existingTask.Attachments, &Attachment{
				ID:        *attachment.ID,
				TaskID:    *attachment.TaskID,
				FileName:  *attachment.FileName,
				CreatedAt: *attachment.CreatedAt,
				UpdatedAt: attachment.UpdatedAt,
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
