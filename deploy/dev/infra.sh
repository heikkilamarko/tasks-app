#!/bin/sh
set -e

cd "$(dirname "$0")"

usage() {
  echo "Usage: $0 <action>"
  echo "  <action>  'up' to start, 'down' to stop."
  exit 1
}

is_zitadel_ready() {
  curl -sf -o /dev/null https://auth.tasks-app.com/debug/ready
}

[ "$#" -lt 1 ] && usage

case "$1" in
  "up")
    nats/configure.sh
    mkdir -p zitadel/machinekey
    docker compose -f infra.yml build
    docker stack deploy --detach=true -c infra.yml tasks-app-infra
    while ! is_zitadel_ready; do
      echo "Waiting for ZITADEL to start..."
      sleep 5
    done
    zitadel/configure.sh
    ;;
  "down")
    docker stack rm --detach=false tasks-app-infra
    sleep 10
    docker volume prune -a -f
    git clean -dfX nats
    git clean -dfX ../../infra/nats
    git clean -dfX zitadel
    ;;
  *)
    usage
    ;;
esac

exit 0
