#!/bin/bash

docker compose -f docker-compose.infra.yml -f docker-compose.single.yml down -v
