#!/bin/bash
set -e

cd "$(dirname "$0")"

export NKEYS_PATH=nsc/keys

env_dir="../../config/dev"
dev_dir="../../src"

operator_name="tasks_app_operator"
account_name="tasks_app_account"
admin_user_name="admin_user"
app_user_name="app_user"
ui_user_name="ui_user"

nats_auth_file="../../messaging/nats/auth.conf"

rm -rf nsc keys

nsc env -s nsc/stores

nsc add operator -n $operator_name -s

nsc add account -n $account_name
nsc edit account -n $account_name \
    --js-disk-storage -1 \
    --js-mem-storage -1

nsc add user -a $account_name -n $admin_user_name
nsc add user -a $account_name -n $app_user_name

nsc add user -a $account_name -n $ui_user_name \
    --bearer \
    --allow-sub "tasks.ui.>" \
    --deny-pub ">"
nsc edit user -a $account_name -n $ui_user_name \
    --subs 1000 \
    --payload 1MB

nsc generate config --mem-resolver > $nats_auth_file

nsc export keys --accounts --account $account_name --dir keys

file_path=$(find keys -maxdepth 1 -type f -name "*.nk")
account_public_key=$(basename "$file_path" .nk)
account_seed=$(cat "$file_path")

sed -i "" -e "s/APP_SHARED_NATS_ACCOUNT_PUBLIC_KEY=.*/APP_SHARED_NATS_ACCOUNT_PUBLIC_KEY=$account_public_key/" "$env_dir/tasks-app.env"
sed -i "" -e "s/APP_SHARED_NATS_ACCOUNT_PUBLIC_KEY=.*/APP_SHARED_NATS_ACCOUNT_PUBLIC_KEY=$account_public_key/" "$dev_dir/tasks-app.env"

sed -i "" -e "s/APP_SHARED_NATS_ACCOUNT_SEED=.*/APP_SHARED_NATS_ACCOUNT_SEED=$account_seed/" "$env_dir/tasks-app.env"
sed -i "" -e "s/APP_SHARED_NATS_ACCOUNT_SEED=.*/APP_SHARED_NATS_ACCOUNT_SEED=$account_seed/" "$dev_dir/tasks-app.env"
