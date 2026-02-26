#!/bin/sh
# One-command extension install when Postgres runs in Docker (root docker-compose).
# Ensures DB is initialized, migrations are applied, then installs the extension.
# Usage: from repo root, run: make install-extension-docker

set -e

ROOT_DIR="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$ROOT_DIR"

# Check Docker Postgres is running
if ! docker compose exec -T postgres pg_isready -U postgres >/dev/null 2>&1; then
  echo "❌ Postgres container is not running."
  echo "   Start it with: docker compose up -d postgres"
  echo "   Then run: make db-init"
  echo "   Then run: make install-extension-docker"
  exit 1
fi

echo "📦 PgQueryNarrative extension (Docker)"
echo ""

# Apply migrations (idempotent: IF NOT EXISTS / CREATE OR REPLACE used)
echo "Applying migrations..."
for f in app/db/migrations/000001_create_schemas.up.sql \
         app/db/migrations/000002_create_app_tables.up.sql \
         app/db/migrations/000003_create_demo_schema.up.sql \
         app/db/migrations/000004_noop.up.sql; do
  docker compose exec -T postgres psql -U postgres -d pgquerynarrative -f - < "$f" || true
done

# Install extension from mounted path (if available) or via stdin
if docker compose exec -T postgres test -r /extension/pgquerynarrative--1.0.sql 2>/dev/null; then
  echo "Installing extension from /extension..."
  docker compose exec -T postgres psql -U postgres -d pgquerynarrative -f /extension/pgquerynarrative--1.0.sql
else
  echo "Installing extension (piping from host)..."
  docker compose exec -T postgres psql -U postgres -d pgquerynarrative -f - < infra/postgres-extension/pgquerynarrative--1.0.sql
fi

# Optional: HTTP-enabled version (only if http extension exists)
if docker compose exec -T postgres psql -U postgres -d pgquerynarrative -tAc "SELECT 1 FROM pg_extension WHERE extname = 'http'" 2>/dev/null | grep -q 1; then
  if docker compose exec -T postgres test -r /extension/pgquerynarrative--1.0--with-http.sql 2>/dev/null; then
    echo "Applying HTTP-enabled functions..."
    docker compose exec -T postgres psql -U postgres -d pgquerynarrative -f /extension/pgquerynarrative--1.0--with-http.sql
  else
    docker compose exec -T postgres psql -U postgres -d pgquerynarrative -f - < infra/postgres-extension/pgquerynarrative--1.0--with-http.sql
  fi
else
  echo "ℹ http extension not found; using base extension (functions return pending status until http is installed)."
fi

echo ""
echo "✅ Extension installed."
echo ""
echo "Quick test:"
echo "  docker compose exec postgres psql -U postgres -d pgquerynarrative -c \"SELECT pgquerynarrative_get_api_url();\""
echo "  docker compose exec postgres psql -U postgres -d pgquerynarrative -c \"SELECT pgquerynarrative_set_api_url('http://host.docker.internal:8080');\""
echo "  docker compose exec postgres psql -U postgres -d pgquerynarrative -c \"SELECT pgquerynarrative_run_query('SELECT 1', 1);\""
echo ""
