CREATE TABLE task (
    id SERIAL PRIMARY KEY,
    name VARCHAR(2000) NOT NULL,
    expires_at timestamptz,
    created_at timestamptz NOT NULL,
    updated_at timestamptz,
    completed_at timestamptz
);

CREATE TABLE notification (
    id SERIAL PRIMARY KEY,
    task_id INT REFERENCES task(id) ON DELETE CASCADE,
    name TEXT NOT NULL CHECK (name IN ('expiring', 'expired')),
    created_at timestamptz NOT NULL
);
