#!/bin/bash
set -euo pipefail

export IMAGE="crk3sdemo.azurecr.io/tasks-app:$1"
export IMAGE_MIGRATIONS="crk3sdemo.azurecr.io/tasks-app-migrations:$1"

export ZITADEL_IP="$(terraform -chdir=../azure-k3s-demo/infra output -raw vm_public_ip)"
export ZITADEL_CLIENT_ID="$(terraform -chdir=infra/zitadel output -raw tasks_app_client_id)"
export ZITADEL_EMAIL_NOTIFIER_PAT="$(terraform -chdir=infra/zitadel output -raw email_notifier_token)"

export NATS_ACCOUNT_PUBLIC_KEY="__TODO__"
export NATS_ACCOUNT_SEED="__TODO__"

envsubst < k8s/tasks-app.yaml | kubectl apply -f -

kubectl create secret generic nats-app-cred --from-file=infra/nats/app.cred --namespace=tasks-app
