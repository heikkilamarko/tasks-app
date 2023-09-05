#!/bin/bash

docker compose -f docker-compose.infra.yml -f docker-compose.single.yml up --build -d
