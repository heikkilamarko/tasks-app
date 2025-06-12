# NATS

## NATS TLS Certificate Setup

Before proceeding, ensure you have the NATS TLS CA certificate available at `certs/ca.crt`.
For details on setting up the certificate, refer to the [NATS example setup documentation](https://github.com/heikkilamarko/azure-k3s-demo/tree/main/examples/nats).

## Create Auth Configuration

Generate NATS auth configuration:

```bash
./nats-auth-create.sh
```

## Deploy Auth Configuration

Deploy the generated auth configuration to the NATS cluster:

```bash
./nats-auth-push.sh
```

## Create App Configuration

Create application-specific NATS configuration:

```bash
./nats-app-create.sh
```
