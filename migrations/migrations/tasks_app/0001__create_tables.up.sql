CREATE TABLE task (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id VARCHAR(200) NOT NULL,
    name VARCHAR(200) NOT NULL,
    expires_at TIMESTAMPTZ,
    expiring_info_at TIMESTAMPTZ,
    expired_info_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ
);

CREATE TABLE attachment (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    task_id BIGINT REFERENCES task(id) ON DELETE CASCADE,
    file_name VARCHAR(200) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ
);
