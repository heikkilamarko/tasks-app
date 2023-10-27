#!/bin/bash

action="$1"
example_directory="$2"

usage() {
    echo "Usage: $0 [up|down] <example_directory>"
    exit 1
}

if [ "$#" -ne 2 ]; then
    usage
fi

if [ ! -d "$example_directory" ]; then
    echo "Example directory '$example_directory' does not exist."
    exit 1
fi

case "$action" in
    up)
        echo "Bringing up containers..."
        docker compose -f docker-compose.yml -f "$example_directory/docker-compose.yml" up --quiet-pull --build -d
        ;;
    down)
        echo "Bringing down containers..."
        docker compose -f docker-compose.yml -f "$example_directory/docker-compose.yml" down -v
        ;;
    *)
        usage
        ;;
esac
