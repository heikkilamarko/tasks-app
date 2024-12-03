package shared

import (
	"context"
	"fmt"
	"time"
)

type PostgresTaskRepository struct {
	db DB
}

var _ TaskRepository = (*PostgresTaskRepository)(nil)

func NewPostgresTaskRepository(db DB) *PostgresTaskRepository {
	return &PostgresTaskRepository{db}
}

func (repo *PostgresTaskRepository) Create(ctx context.Context, task *Task) error {
	user, err := GetUserContext(ctx)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO task
			(user_id, name, expires_at, expiring_info_at, expired_info_at, created_at, updated_at, completed_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	return repo.db.QueryRowContext(
		ctx,
		query,
		user.ID, task.Name, task.ExpiresAt, task.ExpiringInfoAt, task.ExpiredInfoAt, task.CreatedAt, task.UpdatedAt, task.CompletedAt,
	).Scan(&task.ID)
}

func (repo *PostgresTaskRepository) Update(ctx context.Context, task *Task) error {
	user, _ := GetUserContext(ctx)

	query := `
		UPDATE task
		SET
			name = $1,
			expires_at = $2,
			expiring_info_at = $3,
			expired_info_at = $4,
			updated_at = $5,
			completed_at = $6
		WHERE
			id = $7
    `
	args := []any{task.Name, task.ExpiresAt, task.ExpiringInfoAt, task.ExpiredInfoAt, task.UpdatedAt, task.CompletedAt, task.ID}

	if user != nil {
		query += "AND user_id = $8"
		args = append(args, user.ID)
	}

	_, err := repo.db.ExecContext(ctx, query, args...)
	return err
}

func (repo *PostgresTaskRepository) UpdateAttachments(ctx context.Context, taskID int, inserted []string, deleted map[int]string) error {
	now := UTCNow()

	query := `
		INSERT INTO attachment
			(task_id, file_name, created_at)
		VALUES
			($1, $2, $3)
	`

	for _, name := range inserted {
		if _, err := repo.db.ExecContext(ctx, query, taskID, name, now); err != nil {
			return err
		}
	}

	query = `
		DELETE FROM attachment
		WHERE id = $1
	`

	for id := range deleted {
		if _, err := repo.db.ExecContext(ctx, query, id); err != nil {
			return err
		}
	}

	return nil
}

func (repo *PostgresTaskRepository) Delete(ctx context.Context, id int) error {
	user, _ := GetUserContext(ctx)

	query := `
		DELETE FROM task
		WHERE id = $1
	`
	args := []any{id}

	if user != nil {
		query += "AND user_id = $2"
		args = append(args, user.ID)
	}

	_, err := repo.db.ExecContext(ctx, query, args...)
	return err
}

func (repo *PostgresTaskRepository) GetByID(ctx context.Context, id int) (*Task, error) {
	user, _ := GetUserContext(ctx)

	where := `
		WHERE t.id = $1
	`
	args := []any{id}

	if user != nil {
		where += "AND t.user_id = $2"
		args = append(args, user.ID)
	}

	tasks, err := repo.getTasks(ctx, where, "", args...)
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, nil
	}

	return tasks[0], nil
}

func (repo *PostgresTaskRepository) GetActive(ctx context.Context, offset int, limit int) ([]*Task, error) {
	user, _ := GetUserContext(ctx)

	where := `
		WHERE t.completed_at IS NULL
	`
	args := []any{}

	if user != nil {
		where += "AND t.user_id = $1"
		args = append(args, user.ID)
	}

	orderBy := `
		ORDER BY t.created_at DESC
		LIMIT $2 OFFSET $3
	`
	args = append(args, limit, offset)

	return repo.getTasks(ctx, where, orderBy, args...)
}

func (repo *PostgresTaskRepository) GetCompleted(ctx context.Context, offset int, limit int) ([]*Task, error) {
	user, _ := GetUserContext(ctx)

	where := `
		WHERE t.completed_at IS NOT NULL
	`
	args := []any{}

	if user != nil {
		where += "AND t.user_id = $1"
		args = append(args, user.ID)
	}

	orderBy := `
		ORDER BY t.completed_at DESC
		LIMIT $2 OFFSET $3
	`
	args = append(args, limit, offset)

	return repo.getTasks(ctx, where, orderBy, args...)
}

func (repo *PostgresTaskRepository) GetExpiring(ctx context.Context, d time.Duration) ([]*Task, error) {
	user, _ := GetUserContext(ctx)

	t1 := UTCNow()
	t2 := t1.Add(d)

	where := `
		WHERE t.completed_at IS NULL
		AND t.expiring_info_at IS NULL
		AND t.expires_at IS NOT NULL
		AND t.expires_at >= $1
		AND t.expires_at <= $2
	`
	args := []any{t1, t2}

	if user != nil {
		where += "AND t.user_id = $3"
		args = append(args, user.ID)
	}

	orderBy := "ORDER BY t.created_at ASC"

	return repo.getTasks(ctx, where, orderBy, args...)
}

func (repo *PostgresTaskRepository) GetExpired(ctx context.Context) ([]*Task, error) {
	user, _ := GetUserContext(ctx)

	now := UTCNow()

	where := `
		WHERE t.completed_at IS NULL
		AND t.expired_info_at IS NULL
		AND t.expires_at IS NOT NULL
		AND t.expires_at < $1
	`
	args := []any{now}

	if user != nil {
		where += "AND t.user_id = $2"
		args = append(args, user.ID)
	}

	orderBy := "ORDER BY t.created_at ASC"

	return repo.getTasks(ctx, where, orderBy, args...)
}

func (repo *PostgresTaskRepository) DeleteCompleted(ctx context.Context, d time.Duration) (int64, error) {
	user, _ := GetUserContext(ctx)

	t := UTCNow().Add(-d)

	query := `
		DELETE FROM task
		WHERE completed_at IS NOT NULL
		AND completed_at < $1
	`
	args := []any{t}

	if user != nil {
		query += "AND user_id = $2"
		args = append(args, user.ID)
	}

	result, err := repo.db.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *PostgresTaskRepository) getTasks(ctx context.Context, where string, orderBy string, args ...any) ([]*Task, error) {
	var tasks []*Task

	query := fmt.Sprintf(`
		SELECT
			t.id,
			t.user_id,
			t.name,
			t.expires_at,
			t.expiring_info_at,
			t.expired_info_at,
			t.created_at,
			t.updated_at,
			t.completed_at,
			COALESCE(jsonb_agg(a) FILTER (WHERE a.task_id IS NOT NULL), '[]') AS attachments
		FROM
			task t
		LEFT JOIN
			attachment a ON t.id = a.task_id
		%s
		GROUP BY t.id
		%s
	`, where, orderBy)

	rows, err := repo.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		t := &Task{}

		if err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Name,
			&t.ExpiresAt,
			&t.ExpiringInfoAt,
			&t.ExpiredInfoAt,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.CompletedAt,
			&t.Attachments,
		); err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
