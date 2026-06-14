# down-detector
Down-Detector is a service which checks whether the service is available.

```bash
docker compose up -d superset-postgres superset-redis clickhouse
docker compose --profile init run --rm superset-init
docker compose up -d superset
docker compose --profile tools up clickhouse-migrate postgres-migrate
```

## Environment variables

### Application

| Variable | Default | Description |
| --- | --- | --- |
| `POSTGRES_ADDR` | `localhost:5432` | Detector Postgres host and port. |
| `POSTGRES_DB` | `detector` | Detector Postgres database. |
| `POSTGRES_USER` | `detector` | Detector Postgres user. |
| `POSTGRES_PASSWORD` | `detector` | Detector Postgres password. |
| `POSTGRES_SSL_MODE` | `disable` | Detector Postgres SSL mode. |
| `CLICKHOUSE_ADDR` | `localhost:9000` | ClickHouse native protocol host and port. |
| `CLICKHOUSE_DB` | `analytics` | ClickHouse database. |
| `CLICKHOUSE_USER` | `dev` | ClickHouse user. |
| `CLICKHOUSE_PASSWORD` | `password` | ClickHouse password. |
| `SERVER_PORT` | `5436` | HTTP server port. |
| `GIN_MODE` | `test` | Gin mode. |
| `SCHEDULER_CRON` | `*/10 * * * * *` | Inspector scheduler cron expression with seconds. |
| `REPORT_SUBMITTER_SOURCE` | `inspector` | Source written to generated reports. |
| `REPORT_SUBMITTER_LATITUDE` | `55.160023` | Latitude written to generated reports. |
| `REPORT_SUBMITTER_LONGITUDE` | `61.401998` | Longitude written to generated reports. |
| `REPORT_SUBMITTER_IP` | empty | IP written to generated reports. |
| `REPORT_SUBMITTER_PLATFORM` | `ios` | Platform written to generated reports. |
| `SUPERSET_BASE_URL` | `http://localhost:8088` | Superset URL used by the application client. |
| `SUPERSET_ADMIN_USER` | `admin` | Superset admin username used by the application client and init container. |
| `SUPERSET_ADMIN_PASSWORD` | `admin` | Superset admin password used by the application client and init container. |
| `SUPERSET_GUEST_USERNAME` | `guest` | Guest username for embedded dashboard tokens. |
| `SUPERSET_GUEST_FIRSTNAME` | `guest` | Guest first name for embedded dashboard tokens. |
| `SUPERSET_GUEST_LASTNAME` | `guest` | Guest last name for embedded dashboard tokens. |
| `SUPERSET_DASHBOARDS` | empty | Comma-separated dashboards in `name:id` format. Overrides `SUPERSET_DASHBOARD_NAME` and `SUPERSET_DASHBOARD_ID`. |
| `SUPERSET_DASHBOARD_NAME` | `Main` | Default dashboard display name. |
| `SUPERSET_DASHBOARD_ID` | `8f9771f2-2c8a-4e7e-8f10-a5ecd7d39255` | Default dashboard ID. |

### Docker compose and Superset

| Variable | Default | Description |
| --- | --- | --- |
| `COMPOSE_ENV` | `local` | Compose environment label for ClickHouse. |
| `CLICKHOUSE_DEFAULT_ACCESS_MANAGEMENT` | `1` | Enables ClickHouse SQL-driven access management in compose. |
| `SUPERSET_CONFIG_PATH` | `/app/pythonpath/superset_config.py` | Superset config file path inside the container. |
| `SUPERSET_SECRET_KEY` | required | Superset Flask secret key. |
| `CORS_OPTIONS_ORIGINS` | `*` | Comma-separated allowed origins for Superset CORS. |
| `SQLALCHEMY_DATABASE_URI` | set in compose | Superset metadata database URI. |
| `REDIS_HOST` | `redis` in config, `superset-redis` in compose | Redis host for Superset caches. |
| `REDIS_PORT` | `6379` | Redis port for Superset caches. |
| `CACHE_REDIS_DB` | `1` | Redis DB for Superset cache. |
| `RATELIMIT_STORAGE_URI` | `redis://redis:6379/4` in config, `redis://superset-redis:6379/4` in compose | Superset rate limit storage URI. |
| `SUPERSET_ADMIN_FIRST_NAME` | `Admin` | Superset admin first name for init. |
| `SUPERSET_ADMIN_LAST_NAME` | `User` | Superset admin last name for init. |
| `SUPERSET_ADMIN_EMAIL` | `admin@example.com` | Superset admin email for init. |
| `NODE_ENV` | `production` | chouse-ui runtime environment. |
| `JWT_SECRET` | required | chouse-ui JWT secret. |
| `RBAC_ENCRYPTION_KEY` | required | chouse-ui RBAC encryption key. |
| `RBAC_ENCRYPTION_SALT` | required | chouse-ui RBAC encryption salt. |

# TODO

- [ ] Contexts please
- [ ] Workers for routes
- [ ] Errors should not be sent into the outer space from api
- [ ] Make Service Transactional Again
- [ ] Logging for each and every apchih
- [ ] Middlewares
- [ ] Swagger
- [ ] Default values handling inside configs 
- [x] Loading configs from envs
- [ ] Delete Inspectors after route deletion
