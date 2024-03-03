#!/bin/bash
set -e

cd "$(dirname "$0")"

export $(cat tasks-app.env | xargs)

go run cmd/tasks-app/main.go
