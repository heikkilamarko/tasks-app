#!/bin/bash
set -e

cd "$(dirname "$0")"

export TF_VAR_initial_password="S3c_r3t!"

config_dir="../config"
src_dir="../../../src"

terraform init
terraform apply -auto-approve

tasks_app_client_id=$(terraform output -raw tasks_app_client_id)
email_notifier_token=$(terraform output -raw email_notifier_token)

sed -i "" -e "s/APP_UI_AUTH_CLIENT_ID=.*/APP_UI_AUTH_CLIENT_ID=$tasks_app_client_id/" "$config_dir/tasks-app.env"
sed -i "" -e "s/APP_UI_AUTH_CLIENT_ID=.*/APP_UI_AUTH_CLIENT_ID=$tasks_app_client_id/" "$src_dir/tasks-app.env"

sed -i "" -e "s/APP_EMAIL_NOTIFIER_ZITADEL_PAT=.*/APP_EMAIL_NOTIFIER_ZITADEL_PAT=$email_notifier_token/" "$config_dir/tasks-app.env"
sed -i "" -e "s/APP_EMAIL_NOTIFIER_ZITADEL_PAT=.*/APP_EMAIL_NOTIFIER_ZITADEL_PAT=$email_notifier_token/" "$src_dir/tasks-app.env"
