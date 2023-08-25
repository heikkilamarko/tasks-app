INSERT INTO task (name, expires_at, created_at, updated_at, completed_at)
VALUES
    ('Task 1', '2023-08-30 15:00:00', '2023-08-25 10:00:00', NULL, NULL),
    ('Task 2', '2023-09-05 12:00:00', '2023-08-25 11:30:00', NULL, NULL),
    ('Task 3', '2023-09-10 18:00:00', '2023-08-25 14:15:00', NULL, NULL),
    ('Task 4', NULL, '2023-08-25 16:30:00', NULL, NULL),
    ('Task 5', '2023-08-28 09:00:00', '2023-08-25 19:45:00', NULL, NULL);

INSERT INTO notification (task_id, name, created_at)
VALUES
    (1, 'expiring', '2023-08-29 15:00:00'),
    (2, 'expired', '2023-09-06 12:15:00'),
    (3, 'expiring', '2023-09-08 16:30:00'),
    (3, 'expired', '2023-09-11 18:30:00'),
    (5, 'expiring', '2023-08-27 08:30:00');
