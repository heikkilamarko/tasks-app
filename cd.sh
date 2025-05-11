#!/bin/bash
set -euo pipefail

export IMAGE_TAG="$1"

envsubst < k8s/tasks-app.yaml | kubectl apply -f -

kubectl create secret generic nats-app-cred --from-file=infra/nats/app.cred --namespace=tasks-app
