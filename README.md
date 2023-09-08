# Tasks App - Modular Monolith

A simple single binary task management app, built as a modular monolith. The app supports single-process and multi-process setups. In the multi-process setup, each instance is configured to run a subset of the modules.

Modules communicate using the [NATS](https://nats.io/) messaging system.

## Single-process setup

In the single-process setup, all modules are enabled within a single process.

### Run

To start the single-process setup, run the following commands:

```bash
./single_up.sh
```

```bash
open http://localhost:8000/ui
```

### Teardown

To stop the single-process setup, use the following command:

```bash
./single_down.sh
```

## Multi-process setup

The example multi-process setup divides the app into two processes: one for the UI module and another for the backend modules.

Note. Users have the flexibility to run any combination of the modules to tailor the setup to their specific requirements.

### Run

To start the multi-process setup, run the following commands:

```bash
./multi_up.sh
```

```bash
open http://localhost:8000/ui
```

### Teardown

To stop the multi-process setup, use the following command:

```bash
./multi_down.sh
```
