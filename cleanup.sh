#!/bin/bash
set -eu

cd "$(dirname "$0")"

export IMAGE="crk3sdemo.azurecr.io/tasks-app:$1"
export IMAGE_MIGRATIONS="crk3sdemo.azurecr.io/tasks-app-migrations:$1"

export ZITADEL_IP="$(terraform -chdir=../azure-k3s-demo/infra output -raw vm_public_ip)"
export ZITADEL_CLIENT_ID="$(terraform -chdir=infra/zitadel output -raw tasks_app_client_id)"
export ZITADEL_EMAIL_NOTIFIER_PAT="$(terraform -chdir=infra/zitadel output -raw email_notifier_token)"

nsc env -s ~/.local/share/nats/nsc/stores
nsc export keys --accounts --account tasks-app --dir keys
export NATS_ACCOUNT_PUBLIC_KEY=$(nsc describe account --name tasks-app --json | jq -r .sub)
export NATS_ACCOUNT_SEED="$(cat keys/$NATS_ACCOUNT_PUBLIC_KEY.nk)"
rm -rf keys

kubectl delete secret nats-app-cred --namespace=tasks-app || true

envsubst < k8s/tasks-app-$2.yaml | kubectl delete -f - || true
