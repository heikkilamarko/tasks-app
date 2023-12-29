#!/bin/bash

export $(cat tasks-app.env | xargs)

go run cmd/tasks-app/main.go
