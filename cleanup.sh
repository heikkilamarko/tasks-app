#!/bin/bash
set -euo pipefail

export IMAGE_TAG="$1"

kubectl delete secret nats-app-cred --namespace=tasks-app

envsubst < k8s/tasks-app.yaml | kubectl delete -f -
