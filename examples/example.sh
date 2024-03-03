#!/bin/sh
set -e

cd "$(dirname "$0")"

usage() {
  echo "Usage: $0 <action> <compose_file>"
  echo "  <action>        'up' to start, 'down' to stop."
  echo "  <compose_file>  Docker Compose file."
  exit 1
}

check_compose_file() {
  if [ ! -f "$1" ]; then
    echo "Error: Docker Compose file '$1' not found."
    exit 1
  fi
}

[ "$#" -lt 2 ] && usage

action="$1"
compose_file="$2"

check_compose_file "$compose_file"

case "$action" in
  "up")
    docker compose -f "$compose_file" up --build -d
    ;;
  "down")
    docker compose -f "$compose_file" down -v
    ;;
  *)
    usage
    ;;
esac

exit 0
