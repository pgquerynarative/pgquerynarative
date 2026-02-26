#!/usr/bin/env bash
#
# Install pgvector on macOS (Homebrew) for use with local PostgreSQL.
# pgvector does NOT require shared_preload_libraries; CREATE EXTENSION is enough.
#
# Prerequisites: PostgreSQL and Postgres app data already set up (e.g. make db-init).
# After this script, run as a Postgres superuser: ./tools/db/ensure-pgvector-extension.sh  then: make migrate
#
set -e

echo "Installing pgvector via Homebrew..."
brew install pgvector

echo ""
echo "Checking PostgreSQL version..."
if command -v psql >/dev/null 2>&1; then
  PSQL_VERSION=$(psql --version 2>/dev/null | sed -n 's/.* \([0-9]*\)\.*/\1/p' | head -1)
  echo "  psql reports version: $PSQL_VERSION"
else
  echo "  psql not in PATH; ensure PostgreSQL is installed (e.g. brew install postgresql@16 or postgresql@17)."
fi

echo ""
echo "Homebrew pgvector is built for PostgreSQL 17/18. If you use Postgres 16, either:"
echo "  - Upgrade: brew install postgresql@17 and migrate your data, or"
echo "  - Use Docker with an image that includes pgvector (e.g. ankane/pgvector)."
echo ""
echo "Next steps:"
echo "  1. Restart PostgreSQL if it was running during install: brew services restart postgresql@17  (or your version)"
echo "  2. Create the vector extension as superuser: ./tools/db/ensure-pgvector-extension.sh"
echo "  3. Run migrations: make migrate"
echo ""
echo "Done. No shared_preload_libraries change is required for pgvector."
