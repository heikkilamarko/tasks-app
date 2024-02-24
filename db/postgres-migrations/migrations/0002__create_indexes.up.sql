CREATE INDEX idx_task_user_id ON task (user_id);
CREATE INDEX idx_task_expires_at ON task (expires_at);
CREATE INDEX idx_task_created_at ON task (created_at);
CREATE INDEX idx_task_completed_at ON task (completed_at);

CREATE INDEX idx_attachment_task_id ON attachment (task_id);
CREATE INDEX idx_attachment_file_name ON attachment (file_name);
