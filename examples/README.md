# Examples

## Prerequisites

Before running the examples, please complete the following steps:

### 1. Hosts

Add the following entries to your `/etc/hosts` file:

```
127.0.0.1    www.tasks-app.com
127.0.0.1    auth.tasks-app.com
```

### 2. Certificates

We use Caddy's internal CA to generate certificates for the example setup. The root certificate is located at:

```
/proxy/caddy/certs/root.crt
```

Ensure that the certificate is added to the trust store on your development machine.

### 3. Infra Services

To start the infra services (caddy, zitadel, postgres, and nats), execute the following command:

```bash
docker compose up --build -d
```

### 4. ZITADEL Configuration

To configure ZITADEL resources, execute the following commands:

1. Navigate to the `zitadel-configure` directory:

   ```bash
   cd zitadel-configure
   ```

2. Run the configuration script:

   ```bash
   ./configure.sh
   ```

#### Users

When running the examples, use the following example credentials to log in to the application:

| Username | Initial Password |
| -------- | ---------------- |
| `editor` | `S3c_r3t!`       |
| `viewer` | `S3c_r3t!`       |

#### ZITADEL Console (UI)

https://auth.tasks-app.com/ui/console

| Username        | Initial Password |
| --------------- | ---------------- |
| `zitadel-admin` | `S3c_r3t!`       |

## Single-Process Setup

In the single-process setup, all modules are enabled within a single process.

### Running

To start the single-process setup, execute the following command:

```bash
docker compose -f compose_single.yml up --build -d
```

Open the app in your web browser: http://www.tasks-app.com/ui

### Teardown

To stop the single-process setup and tear down the services, use the following command:

```bash
docker compose -f compose_single.yml down -v
```

## Multi-Process Setup

A multi-process setup divides the application into two or more processes.

### Running

To start a multi-process setup, you can choose one of the following commands based on your needs:

```bash
docker compose -f compose_multi_frontend_backend.yml up --build -d
```

or

```bash
docker compose -f compose_multi_all.yml up --build -d
```

Once the setup is running, open the application in your web browser: http://www.tasks-app.com/ui

### Teardown

To stop the multi-process setup, use one of the following commands, depending on your setup choice:

```bash
docker compose -f compose_multi_frontend_backend.yml down -v
```

or

```bash
docker compose -f compose_multi_all.yml down -v
```

## Teardown Infra Services

To tear down the infra services, execute the following command:

```bash
docker compose down -v
```
