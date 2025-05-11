#!/bin/bash
set -euo pipefail

cd "$(dirname "$0")"

nats_server="$(terraform -chdir=../../../azure-k3s-demo/infra output -raw vm_public_ip):4222"

nats auth account add \
    --operator example \
    --bearer \
    --jetstream \
    --js-disk 1GB \
    --js-memory 1GB \
    --payload 2MB \
    --defaults \
    tasks-app

nats auth user add \
    --operator example \
    --bearer \
    --payload=-1 \
    --defaults \
    admin tasks-app

nats auth user add \
    --operator example \
    --bearer \
    --payload=-1 \
    --defaults \
    app tasks-app

nats auth user credential \
    --operator example \
    admin.cred admin tasks-app

nats ctx add tasks-app-admin \
    --server $nats_server \
    --creds "$PWD/admin.cred"

nats auth user credential \
    --operator example \
    app.cred app tasks-app

nats ctx add tasks-app-app \
    --server $nats_server \
    --creds "$PWD/app.cred"
