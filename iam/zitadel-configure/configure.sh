#!/bin/bash

env_dir="../../config/dev"

terraform init
terraform apply -auto-approve

client_id=$(terraform output -raw tasks_app_client_id)

sed -i "" -e "s/APP_UI_AUTH_CLIENT_ID=.*/APP_UI_AUTH_CLIENT_ID=$client_id/" "$env_dir/tasks-app.env"
