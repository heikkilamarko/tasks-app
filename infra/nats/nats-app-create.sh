#!/bin/sh
set -euo pipefail

cd "$(dirname "$0")"

export NATS_CONTEXT=tasks-app-admin

nats stream add --config "streams/tasks.json"
nats stream add --config "streams/tasks_dlq.json"

nats consumer add tasks --config "consumers/tasks.json"
