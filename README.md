# Tasks App - Modular Monolith

A simple single binary task management app, built as a modular monolith. The app supports single-process and multi-process setups. In the multi-process setup, each instance is configured to run a subset of the modules.

Modules communicate using the [NATS](https://nats.io/) messaging system.

![components](doc/components.png)

See the `/examples` directory for some example setups.
