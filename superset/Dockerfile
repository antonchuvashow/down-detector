FROM apache/superset:6.1.0

USER root

RUN . /app/.venv/bin/activate && \
    uv pip install \
      "psycopg2-binary>=2.9.9" \
      "clickhouse-connect>=0.7.0"

USER superset