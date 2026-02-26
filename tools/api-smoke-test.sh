#!/usr/bin/env bash
#
# PgQueryNarrative API smoke test
#
# Exercises core and optional endpoints: queries, saved queries, suggestions
# (intent and similar), report generation, and report list. Requires the server
# to be running and the database migrated (including 000005 for embeddings).
# For similar-query and RAG: set EMBEDDING_BASE_URL and use Ollama with
# nomic-embed-text.
#
# Usage: ./tools/api-smoke-test.sh
#        PGQUERYNARRATIVE_URL=http://host:port ./tools/api-smoke-test.sh
#
set -e

BASE="${PGQUERYNARRATIVE_URL:-http://localhost:8080}"
API="${BASE}/api/v1"

if ! command -v jq >/dev/null 2>&1; then
  echo "Error: jq is required. Install jq and run again." >&2
  exit 1
fi

echo "PgQueryNarrative API smoke test"
echo "Base URL: $BASE"
echo ""

echo "1. Queries — run"
curl -sf -X POST "$API/queries/run" \
  -H "Content-Type: application/json" \
  -d '{"sql":"SELECT product_category, SUM(total_amount) AS total FROM demo.sales GROUP BY product_category ORDER BY total DESC","limit":100}' \
  | jq -r '.columns[].name, (.rows | length | "rows: \(.)")'

echo ""
echo "2. Queries — save (3 saved queries)"
for payload in \
  '{"name":"Sales by category","sql":"SELECT product_category, SUM(total_amount) AS total FROM demo.sales GROUP BY product_category ORDER BY total DESC","description":"Total revenue per product category"}' \
  '{"name":"Daily sales over time","sql":"SELECT date, SUM(total_amount) AS daily_total FROM demo.sales GROUP BY date ORDER BY date","description":"Daily revenue time series"}' \
  '{"name":"Revenue by region","sql":"SELECT region, SUM(total_amount) AS total FROM demo.sales GROUP BY region ORDER BY total DESC","description":"Sales totals by region"}'; do
  curl -sf -X POST "$API/queries/saved" -H "Content-Type: application/json" -d "$payload" \
    | jq -r '.id // .name // empty'
done

echo ""
echo "3. Suggestions — by intent (keyword match)"
curl -sf "$API/suggestions/queries?intent=sales&limit=5" \
  | jq -r '.suggestions[] | "  \(.source): \(.title)"'

echo ""
echo "4. Suggestions — similar (semantic; requires embeddings)"
curl -sf "$API/suggestions/similar?text=revenue%20by%20category&limit=5" \
  | jq -r 'if (.suggestions | length) > 0 then (.suggestions[] | "  \(.source): \(.title)") else "  (none — enable EMBEDDING_BASE_URL and re-save queries)" end'

echo ""
echo "5. Reports — generate (time-series, predictive summary; RAG if embeddings enabled)"
curl -sf -X POST "$API/reports/generate" \
  -H "Content-Type: application/json" \
  -d '{"sql":"SELECT date, SUM(total_amount) AS daily_total FROM demo.sales GROUP BY date ORDER BY date"}' \
  | jq -r '"  headline: " + (.narrative.headline // ""), (.metrics.time_series | to_entries[]? | "  \(.key): predictive_summary=\(.value.predictive_summary // "—")")'

echo ""
echo "6. Reports — list (latest 3)"
curl -sf "$API/reports?limit=3" \
  | jq -r '.items[]? | "  \(.id) \(.created_at)"'

echo ""
echo "Smoke test complete."
