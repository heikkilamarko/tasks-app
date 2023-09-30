CREATE TABLE task (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    expires_at timestamptz,
    expiring_info_at timestamptz,
    expired_info_at timestamptz,
    created_at timestamptz NOT NULL,
    updated_at timestamptz,
    completed_at timestamptz
);

CREATE TABLE attachment (
    id SERIAL PRIMARY KEY,
    task_id INT REFERENCES task(id) ON DELETE CASCADE,
    file_name VARCHAR(200) NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz
);
