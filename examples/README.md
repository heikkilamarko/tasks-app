# Examples

## Single-Process Setup

In the single-process setup, all modules are enabled within a single process.

### Running

To start the single-process setup, execute the following command:

```bash
./example.sh up single
```

Open the app in your web browser: http://localhost:8000/ui

### Teardown

To stop the single-process setup, use the following command:

```bash
./example.sh down single
```

## Multi-Process Setup

A multi-process setup divides the application into two or more processes.

### Running

To start a multi-process setup, you can choose one of the following commands based on your needs:

```bash
./example.sh up multi_frontend_backend
```

or

```bash
./example.sh up multi_all
```

Once the setup is running, open the application in your web browser: http://localhost:8000/ui

### Teardown

To stop the multi-process setup, use one of the following commands, depending on your setup choice:

```bash
./example.sh down multi_frontend_backend
```

or

```bash
./example.sh down multi_all
```
