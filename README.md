# down-detector
Down-Detector is a service which checks whether the service is available.

```bash
docker compose up -d superset-postgres superset-redis clickhouse
docker compose --profile init run --rm superset-init
docker compose up -d superset
docker compose --profile tools up clickhouse-migrate postgres-migrate
```

# TODO

- [ ] Contexts please
- [ ] Workers for routes
- [ ] Errors should not be sent into the outer space from api
- [ ] Make Service Transactional Again
- [ ] Logging for each and every apchih
- [ ] Middlewares
- [ ] Swagger
- [ ] Default values handling inside configs 
- [ ] Loading configs from envs
- [ ] Delete Inspectors after route deletion