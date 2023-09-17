#!/bin/bash

if [ ! -d "$1" ]; then
    echo "Example $1 does not exist."
    exit 1
fi

docker compose -f docker-compose.yml -f $1/docker-compose.yml down -v
