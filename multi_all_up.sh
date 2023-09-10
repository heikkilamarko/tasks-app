#!/bin/bash

docker compose -f docker-compose.infra.yml -f docker-compose.multi.all.yml up --build -d
