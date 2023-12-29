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
    docker compose up --build -d
    pushd zitadel && ./configure.sh && popd
    ;;
  "down")
    docker compose down -v
    git clean -dfX zitadel
    ;;
  *)
    usage
    ;;
esac

exit 0
