#!/bin/sh

migrate -path /migrations/postgres -database $POSTGRES_POSTGRES_CONNECTIONSTRING up
migrate -path /migrations/tasks_app -database $POSTGRES_TASKS_APP_CONNECTIONSTRING up
