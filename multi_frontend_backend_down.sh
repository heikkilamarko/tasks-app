#!/bin/bash

docker compose -f docker-compose.infra.yml -f docker-compose.multi.frontend-backend.yml down -v
