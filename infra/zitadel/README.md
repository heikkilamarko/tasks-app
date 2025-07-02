# ZITADEL

## Prepare Credentials

Generate the `zitadel-admin-sa.json` file from the `zitadel-admin-sa` Kubernetes secret.

## Caddy Proxy

> Only required if not using trusted certificates.

Add this entry to your hosts file to map `zitadel.test` to localhost:

```bash
127.0.0.1 zitadel.test
```

Start the Caddy proxy:

```bash
./caddy.sh
```

## Terraform

Initialize and apply the Terraform configuration:

```bash
terraform init
```

```bash
terraform apply -var-file=<(envsubst < vars.tfvars)
```

## Outputs

Use the following output values to configure the application:

```bash
terraform output -raw tasks_app_client_id
```

```bash
terraform output -raw email_notifier_token
```
