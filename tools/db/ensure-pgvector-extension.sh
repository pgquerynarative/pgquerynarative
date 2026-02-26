#!/usr/bin/env bash
#
# Create the pgvector extension in the app database. Must be run as a PostgreSQL
# superuser (e.g. your macOS user with Homebrew Postgres, or postgres).
# Run once before 'make migrate' so migration 000007 can add the vector column.
#
set -e

DB="${PGDATABASE:-pgquerynarrative}"
HOST="${PGHOST:-localhost}"
PORT="${PGPORT:-5432}"

if [ -n "$SUPERVISOR_DB_URL" ]; then
  psql "$SUPERVISOR_DB_URL" -c 'CREATE EXTENSION IF NOT EXISTS vector;'
else
  psql -h "$HOST" -p "$PORT" -d "$DB" -c 'CREATE EXTENSION IF NOT EXISTS vector;'
fi

echo "Extension vector is installed. You can run: make migrate"
