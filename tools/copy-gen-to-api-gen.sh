#!/bin/sh
# Copy generated code from gen/ to api/gen/ and fix import paths.
# Goa generates into gen/; the app uses only api/gen. Single tree, no gen/ in imports.
set -e
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

[ -d gen ] || exit 0

cp gen/queries/*.go api/gen/queries/
cp gen/reports/*.go api/gen/reports/
cp -r gen/http/* api/gen/http/

[ -d gen/schema ] && mkdir -p api/gen/schema && cp gen/schema/*.go api/gen/schema/
[ -d gen/suggestions ] && mkdir -p api/gen/suggestions && cp gen/suggestions/*.go api/gen/suggestions/

# Fix imports in api/gen: gen/ -> api/gen/
for f in api/gen/http/queries/server/*.go api/gen/http/queries/client/*.go \
         api/gen/http/reports/server/*.go api/gen/http/reports/client/*.go; do
	[ -f "$f" ] || continue
	sed -i.bak 's|github.com/pgquerynarrative/pgquerynarrative/gen/queries|github.com/pgquerynarrative/pgquerynarrative/api/gen/queries|g' "$f"
	sed -i.bak 's|github.com/pgquerynarrative/pgquerynarrative/gen/reports|github.com/pgquerynarrative/pgquerynarrative/api/gen/reports|g' "$f"
	rm -f "$f.bak"
done
for f in api/gen/http/schema/server/*.go api/gen/http/schema/client/*.go \
         api/gen/http/suggestions/server/*.go api/gen/http/suggestions/client/*.go; do
	[ -f "$f" ] || continue
	sed -i.bak 's|github.com/pgquerynarrative/pgquerynarrative/gen/schema|github.com/pgquerynarrative/pgquerynarrative/api/gen/schema|g' "$f"
	sed -i.bak 's|github.com/pgquerynarrative/pgquerynarrative/gen/suggestions|github.com/pgquerynarrative/pgquerynarrative/api/gen/suggestions|g' "$f"
	rm -f "$f.bak"
done
for f in api/gen/http/cli/pgquerynarrative/cli.go; do
	[ -f "$f" ] || continue
	sed -i.bak 's|github.com/pgquerynarrative/pgquerynarrative/gen/http/|github.com/pgquerynarrative/pgquerynarrative/api/gen/http/|g' "$f"
	rm -f "$f.bak"
done

echo "Synced gen/ to api/gen/ (imports fixed)"
