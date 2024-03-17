# Dev Setup

```bash
cd ../examples
```

```bash
./infra.sh up
```

```bash
docker service rm tasks-app-infra_caddy
```

```bash
cd ../src
```

```bash
./caddy.sh
```

```bash
./tasks-app.sh
```
