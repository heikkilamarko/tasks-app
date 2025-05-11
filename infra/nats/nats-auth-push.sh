#!/bin/bash
set -euo pipefail

nats auth account push \
    --context azure-k3s-demo-system \
    --operator example \
    tasks-app
