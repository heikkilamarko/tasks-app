#!/bin/sh
set -e

cd "$(dirname "$0")"

usage() {
  echo "Usage: $0 <action>"
  echo "  <action>  'up' to start, 'down' to stop."
  exit 1
}

[ "$#" -lt 1 ] && usage

case "$1" in
  "up")
    nats/configure.sh
    mkdir -p zitadel/machinekey
    docker compose -f infra.prod.yml build
    docker stack deploy -c infra.prod.yml tasks-app-infra
    echo "Waiting for services to start..."
    sleep 60
    # scp -r raspi:/home/raspi/zitadel/machinekey ./zitadel/
    zitadel/configure.sh
    ;;
  "down")
    docker stack rm tasks-app-infra
    sleep 10
    docker volume prune -a -f
    git clean -dfX nats
    git clean -dfX ../messaging/nats
    git clean -dfX zitadel
    ;;
  *)
    usage
    ;;
esac

exit 0
