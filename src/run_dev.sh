#!/bin/bash

export $(cat .env.development | xargs)

go run cmd/tasks-app/main.go
