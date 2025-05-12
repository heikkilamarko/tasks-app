# Tasks App - Modular Monolith

A simple single binary task management app, built as a modular monolith. The app supports single-process and multi-process setups. In the multi-process setup, each instance is configured to run a subset of the modules.

## Tech Stack

| TECHNOLOGY                                 | DESCRIPTION                                   |
| ------------------------------------------ | --------------------------------------------- |
| [NATS](https://nats.io/)                   | Messaging, WebSockets, KV Store, Object Store |
| [PostgreSQL](https://www.postgresql.org/)  | Database                                      |
| [ZITADEL](https://zitadel.com/)            | Identity and Access Management (IAM)          |
| [Terraform](https://www.terraform.io/)     | Infrastructure Automation                     |
| [Bash](https://www.gnu.org/software/bash/) | Scripting                                     |
| [Go](https://go.dev/)                      | Programming Language                          |
| [htmx](https://htmx.org/)                  | Web Technology                                |
| [Hyperscript](https://hyperscript.org/)    | Frontend Scripting Language                   |

## Infrastructure Setup

The infrastructure is built on top of the Azure K3s setup provided in the following GitHub repository:

https://github.com/heikkilamarko/azure-k3s-demo

### NATS Configuration

Follow the instructions in [infra/nats](infra/nats) to configure NATS.

### PostgreSQL Configuration

Follow the instructions in [infra/postgresql](infra/postgresql) to configure PostgreSQL.

### ZITADEL Configuration

Follow the instructions in [infra/zitadel](infra/zitadel) to configure ZITADEL.

## Application Deployment

### Build the Application

```bash
./ci.sh <docker_image_tag>
```

### Deploy as a Single Binary

```bash
./cd.sh <docker_image_tag> single
```

```bash
./cleanup.sh <docker_image_tag> single
```

### Deploy as Microservices

```bash
./cd.sh <docker_image_tag> micro
```

```bash
./cleanup.sh <docker_image_tag> micro
```
