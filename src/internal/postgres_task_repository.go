package internal

import (
	"context"
	"database/sql"
	"time"
)

type PostgresTaskRepository struct {
	db *sql.DB
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
			return nil, nil // Return nil if no task with the given ID is found
		}
		return nil, err
	}

	return task, nil
}

func (repo *PostgresTaskRepository) GetAll(ctx context.Context) ([]*Task, error) {
	query := `
		SELECT id, name, expires_at, expiring_info_at, expired_info_at, created_at, updated_at, completed_at
		FROM task
		ORDER BY created_at ASC
	`

	rows, err := repo.db.QueryContext(ctx, query)
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

func (repo *PostgresTaskRepository) GetExpiring(ctx context.Context, expirationWindow time.Duration) ([]*Task, error) {
	now := time.Now().UTC()
	expirationTime := now.Add(expirationWindow)

	query := `
		SELECT id, name, expires_at, expiring_info_at, expired_info_at, created_at, updated_at, completed_at
		FROM task
		WHERE expires_at IS NOT NULL
		AND expiring_info_at IS NULL
		AND expires_at >= $1
		AND expires_at <= $2
		ORDER BY created_at ASC
    `

	rows, err := repo.db.QueryContext(ctx, query, now, expirationTime)
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

func (repo *PostgresTaskRepository) GetExpired(ctx context.Context) ([]*Task, error) {
	now := time.Now().UTC()

	query := `
		SELECT id, name, expires_at, expiring_info_at, expired_info_at, created_at, updated_at, completed_at
		FROM task
		WHERE expires_at IS NOT NULL
		AND expired_info_at IS NULL
		AND expires_at < $1
		ORDER BY created_at ASC
    `

	rows, err := repo.db.QueryContext(ctx, query, now)
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
