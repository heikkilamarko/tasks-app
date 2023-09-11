#!/bin/bash

docker compose -f docker-compose.yml -f $1/docker-compose.yml down -v
