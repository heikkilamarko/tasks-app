#!/bin/sh
set -euo pipefail

cd "$(dirname "$0")"

nats stream add \
    --context tasks-app-admin \
    --config streams/tasks.json

nats stream add \
    --context tasks-app-admin \
    --config streams/tasks_dlq.json

nats consumer add tasks \
    --context tasks-app-admin \
    --config consumers/tasks.json
