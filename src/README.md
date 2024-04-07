# Dev Setup

```bash
cd ../deploy/dev
```

```bash
./infra.sh up
```

```bash
docker service rm tasks-app-infra_caddy
```

```bash
cd ../../src
```

```bash
caddy run
```

```bash
./tasks-app.sh
```
