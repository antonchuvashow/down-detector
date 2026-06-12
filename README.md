# down-detector
Down-Detector is a service which checks whether the service is available.

```bash
docker compose up -d postgres redis clickhouse
docker compose --profile init run --rm superset-init
docker compose up -d superset
docker compose --profile tools up migrate
```