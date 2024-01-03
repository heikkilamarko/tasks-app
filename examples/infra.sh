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
    docker compose -f infra.yml up --build -d
    sleep 5
    pushd zitadel && ./configure.sh && popd
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
