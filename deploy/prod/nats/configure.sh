#!/bin/bash
set -e

cd "$(dirname "$0")"

export NKEYS_PATH=nsc/keys

config_dir="../config"

operator_name="tasks_app_operator"
account_name="tasks_app_account"
admin_user_name="admin_user"
app_user_name="app_user"

nats_auth_file="../../../infra/nats/auth.conf"

rm -rf nsc keys

nsc env -s nsc/stores

nsc add operator --generate-signing-key --sys -n $operator_name
nsc edit operator --require-signing-keys

nsc add account -n $account_name
nsc edit account -n $account_name \
    --sk generate \
    --js-disk-storage -1 \
    --js-mem-storage -1

nsc add user -a $account_name -n $admin_user_name
nsc add user -a $account_name -n $app_user_name

nsc generate config --mem-resolver --sys-account SYS > $nats_auth_file

nsc export keys --accounts --account $account_name --dir keys

account_public_key=$(nsc describe account --json --field sub | tr -d '"')

signing_key_file=keys/$(nsc describe account --json --field nats.signing_keys.0 | tr -d '"').nk
account_seed=$(cat "$signing_key_file")

sed -i "" -e "s/APP_SHARED_NATS_ACCOUNT_PUBLIC_KEY=.*/APP_SHARED_NATS_ACCOUNT_PUBLIC_KEY=$account_public_key/" "$config_dir/tasks-app.env"
sed -i "" -e "s/APP_SHARED_NATS_ACCOUNT_SEED=.*/APP_SHARED_NATS_ACCOUNT_SEED=$account_seed/" "$config_dir/tasks-app.env"
