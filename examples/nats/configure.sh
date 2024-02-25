#!/bin/bash
set -e

cd "$(dirname "$0")"

export NKEYS_PATH=nsc/keys

operator_name="tasks_app_operator"
account_name="tasks_app_account"
admin_user_name="admin_user"
app_user_name="app_user"
ui_user_name="ui_user"

nats_auth_file="../../messaging/nats/auth.conf"

rm -rf nsc

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
