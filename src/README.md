# Dev Setup

```bash
cd ../examples
```

```bash
./infra.sh up
```

```bash
docker compose -f infra.yml stop caddy
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
