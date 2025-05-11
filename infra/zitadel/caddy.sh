#!/bin/sh
set -euo pipefail

cd "$(dirname "$0")"

zitadel_address="https://$(terraform -chdir=../../../azure-k3s-demo/infra output -raw vm_public_ip)"

sudo ZITADEL_ADDRESS=$zitadel_address caddy run
