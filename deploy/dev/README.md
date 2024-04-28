# Examples

## Prerequisites

Before running the examples, please complete the following steps:

### 1. Hosts

Add the following entries to your `/etc/hosts` file:

```
127.0.0.1    www.tasks-app.com
127.0.0.1    auth.tasks-app.com
127.0.0.1    smtp.tasks-app.com
```

### 2. Certificates

We use Caddy's internal CA to generate certificates for the example setup. The root certificate is located at:

```
./certs/root.crt
```

Ensure that the certificate is added to the trust store on your development machine.

### 3. Infra Services

To start and configure the infra services (caddy, zitadel, postgres, and nats), execute the following command:

```bash
./infra.sh up
```

After running the examples, use the following command to tear down the infra services:

```bash
./infra.sh down
```

## ZITADEL Console (UI)

https://auth.tasks-app.com

| Username        | Initial Password |
| --------------- | ---------------- |
| `zitadel-admin` | `S3c_r3t!`       |

## smtp4dev UI

https://smtp.tasks-app.com

| Username | Password   |
| -------- | ---------- |
| `admin`  | `S3c_r3t!` |

## Example Users

When running the examples, use the following example credentials to log in to the application:

| Username  | Initial Password |
| --------- | ---------------- |
| `johndoe` | `S3c_r3t!`       |
| `janedoe` | `S3c_r3t!`       |

## Examples

### Single-Process Setup

In the single-process setup, all modules are enabled within a single process.

#### Start

```bash
./example.sh up example_single.yml
```

App URL: http://www.tasks-app.com/ui

#### Stop

```bash
./example.sh down example_single.yml
```

### Multi-Process Setup 1

A multi-process setup that divides the application into frontend and backend processes.

#### Start

```bash
./example.sh up example_multi_frontend_backend.yml
```

App URL: http://www.tasks-app.com/ui

#### Stop

```bash
./example.sh down example_multi_frontend_backend.yml
```

### Multi-Process Setup 2

A multi-process setup where each module operates within its own dedicated process.

#### Start

```bash
./example.sh up example_multi_all.yml
```

App URL: http://www.tasks-app.com/ui

#### Stop

```bash
./example.sh down example_multi_all.yml
```
