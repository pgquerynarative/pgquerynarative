#!/bin/sh
# Production entrypoint: wait for Postgres, run migrations via binary, optional seed, then server.
set -e

DB_HOST="${DATABASE_HOST:-postgres}"
DB_PORT="${DATABASE_PORT:-5432}"
DB_NAME="${DATABASE_NAME:-pgquerynarrative}"
DB_USER="${DATABASE_USER:-pgquerynarrative_app}"
DB_PASSWORD="${DATABASE_PASSWORD:-pgquerynarrative_app}"

export PGPASSWORD="${DB_PASSWORD}"

attempts=30
while [ $attempts -gt 0 ]; do
  if pg_isready -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" >/dev/null 2>&1; then
    break
  fi
  attempts=$((attempts - 1))
  sleep 1
done

if [ $attempts -eq 0 ]; then
  echo "Postgres is not ready at ${DB_HOST}:${DB_PORT}"
  exit 1
fi

DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DATABASE_SSL_MODE:-disable}"
/app/bin/migrate -path /app/migrations -database "${DB_URL}" up

if [ "${PGQUERYNARRATIVE_SEED:-false}" = "true" ]; then
  psql "${DB_URL}" -f /app/tools/db/seed.sql
fi

exec /app/bin/server
