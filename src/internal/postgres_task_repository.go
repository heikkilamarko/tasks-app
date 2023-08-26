package internal

import "database/sql"

type PostgresTaskRepository struct {
	db *sql.DB
}

func (repo *PostgresTaskRepository) Create(task *Task) error {
	query := `
        INSERT INTO task (name, expires_at, created_at, updated_at, completed_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `

	err := repo.db.QueryRow(
		query,
		task.Name, task.ExpiresAt, task.CreatedAt, task.UpdatedAt, task.CompletedAt,
	).Scan(&task.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *PostgresTaskRepository) Update(task *Task) error {
	query := `
        UPDATE task
        SET name = $2, expires_at = $3, updated_at = $4, completed_at = $5
        WHERE id = $1
    `

	_, err := repo.db.Exec(
		query,
		task.ID, task.Name, task.ExpiresAt, task.UpdatedAt, task.CompletedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repo *PostgresTaskRepository) Delete(id int) error {
	query := `
        DELETE FROM task
        WHERE id = $1
    `

	_, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

func (repo *PostgresTaskRepository) GetByID(id int) (*Task, error) {
	query := `
        SELECT id, name, expires_at, created_at, updated_at, completed_at
        FROM task
        WHERE id = $1
    `

	row := repo.db.QueryRow(query, id)

	task := &Task{}
	err := row.Scan(
		&task.ID, &task.Name, &task.ExpiresAt, &task.CreatedAt, &task.UpdatedAt, &task.CompletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil if no task with the given ID is found
		}
		return nil, err
	}

	return task, nil
}

func (repo *PostgresTaskRepository) GetAll() ([]*Task, error) {
	rows, err := repo.db.Query(`
        SELECT id, name, expires_at, created_at, updated_at, completed_at
        FROM task
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}

		err := rows.Scan(
			&task.ID, &task.Name, &task.ExpiresAt, &task.CreatedAt, &task.UpdatedAt, &task.CompletedAt,
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
