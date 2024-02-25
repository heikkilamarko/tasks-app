#!/bin/sh
set -e

usage() {
  echo "Usage: $0 <action>"
  echo "  <action>  'up' to start, 'down' to stop."
  exit 1
}

[ "$#" -lt 1 ] && usage

case "$1" in
  "up")
    nats/configure.sh
    docker compose -f infra.yml up --build -d
    sleep 10
    zitadel/configure.sh
    ;;
  "down")
    docker compose -f infra.yml down -v
    git clean -dfX zitadel
    ;;
  *)
    usage
    ;;
esac

exit 0
