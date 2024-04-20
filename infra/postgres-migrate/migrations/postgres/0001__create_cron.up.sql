CREATE EXTENSION IF NOT EXISTS pg_cron;

SELECT cron.schedule_in_database('tasks-app-vacuum', '0 3 * * *', 'vacuum full', 'tasks_app');
