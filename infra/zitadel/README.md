# ZITADEL

## Update the Hosts File

Add the following entry to your system’s hosts file to map `zitadel.local` to your local machine:

```bash
127.0.0.1 zitadel.local
```

## Prepare ZITADEL Credentials

Create the `zitadel-admin-sa.json` file from the Kubernetes secret `zitadel-admin-sa`.

## Caddy Proxy

```bash
./caddy.sh
```

## Terraform

```bash
terraform init
```

```bash
terraform apply
```

## Outputs

Use the output values below to configure the application:

```bash
terraform output -raw tasks_app_client_id
```

```bash
terraform output -raw email_notifier_token
```
