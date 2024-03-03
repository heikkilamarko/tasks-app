#!/bin/bash
set -e

cd "$(dirname "$0")"

caddy run --envfile caddy.env -c ../proxy/caddy/Caddyfile
