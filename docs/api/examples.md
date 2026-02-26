# API examples

Base URL: `http://localhost:8080` (or set `PGQUERYNARRATIVE_PORT`). All examples use `application/json`.

## Run query

```bash
curl -X POST http://localhost:8080/api/v1/queries/run \
  -H "Content-Type: application/json" \
  -d '{
    "sql": "SELECT product_category, SUM(total_amount) AS total FROM demo.sales GROUP BY product_category",
    "limit": 10
  }'
```

Response: `columns`, `rows`, `row_count`, `execution_time_ms`, optional `chart_suggestions`, `period_comparison` (when result has a time column and measures).

## Save query

```bash
curl -X POST http://localhost:8080/api/v1/queries/saved \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Top Products",
    "sql": "SELECT product_category, SUM(total_amount) FROM demo.sales GROUP BY product_category",
    "tags": ["sales", "top"]
  }'
```

## Generate report

Requires a configured [LLM](../getting-started/llm-setup.md). See [Configuration – LLM](../configuration.md#llm).

```bash
curl -X POST http://localhost:8080/api/v1/reports/generate \
  -H "Content-Type: application/json" \
  -d '{
    "sql": "SELECT product_category, SUM(total_amount) AS total FROM demo.sales GROUP BY product_category"
  }'
```

Response: `narrative`, `metrics` (aggregates, time_series, data_quality, perf_suggestions, etc.).

### Single-period query

When the result has no time-series comparison, the narrative will not mention "previous period":

```bash
curl -s -X POST http://localhost:8080/api/v1/reports/generate \
  -H "Content-Type: application/json" \
  -d '{"sql": "SELECT SUM(total_amount) AS total_fares, COUNT(*) AS trips FROM demo.sales"}' \
  | jq '.narrative'
```

### Get / list reports

```bash
curl -s http://localhost:8080/api/v1/reports/REPORT_UUID | jq .
curl -s "http://localhost:8080/api/v1/reports?limit=5&offset=0" | jq '.items[] | {id, created_at}'
```

## See also

- [API reference](README.md) — Full endpoint list and error codes
- [LLM setup](../getting-started/llm-setup.md) — Required for report generation
- [Configuration](../configuration.md) — Base URL and LLM
- [Quick start](../getting-started/quickstart.md) — Get the app running
- [Documentation index](../README.md)
