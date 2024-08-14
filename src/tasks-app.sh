#!/bin/bash
set -e

cd "$(dirname "$0")"

pushd internal/modules/ui/web
npm install
npm run build:dev
popd

export $(cat tasks-app.env | xargs)

go run cmd/tasks-app/main.go
