# API reference

REST API base: `http://localhost:8080/api/v1`. No authentication in current version.

## Queries

| Method | Path | Description |
|--------|------|-------------|
| POST | `/queries/run` | Body: `{"sql":"...", "limit": 100}`. Run read-only SQL. Returns optional `period_comparison` and `chart_suggestions`. |
| POST | `/queries/saved` | Body: `{"name","sql","tags"}`. Save query. |
| GET | `/queries/saved` | Query: `limit`, `offset`, `tags`. List saved. |
| GET | `/queries/saved/{id}` | Get saved query. |
| DELETE | `/queries/saved/{id}` | Delete saved query. |

## Schema

| Method | Path | Description |
|--------|------|-------------|
| GET | `/schema` | Returns allowed schemas with tables and columns (from `information_schema`). Used by MCP `get_schema` and for discovery. |

## Suggestions

| Method | Path | Description |
|--------|------|-------------|
| GET | `/suggestions/queries` | Query: `intent` (optional), `limit` (default 5). Returns suggested SQL: curated examples plus saved queries matching intent (name/description/sql). Used by MCP `suggest_queries`. |
| GET | `/suggestions/similar` | Query: `text` (required), `limit` (default 5). Returns saved queries semantically similar to the given text (embedding-based). Requires embeddings to be enabled (`EMBEDDING_BASE_URL`). |

## Reports

| Method | Path | Description |
|--------|------|-------------|
| POST | `/reports/generate` | Body: `{"sql":"...", "saved_query_id": "uuid"}`. Generate report (requires LLM). |
| GET | `/reports/{id}` | Get report. Metrics include `time_series` (with optional `periods`, `moving_average`, `anomalies`, `trend_summary`, `next_period_forecast`, `predictive_summary`), `data_quality`, `perf_suggestions`, `period_current_label`, `period_previous_label`. |
| GET | `/reports` | Query: `limit`, `offset`, `saved_query_id`. List reports. |

## Errors

JSON: `{"name","message","code"}`. Codes: `VALIDATION_ERROR`, `TIMEOUT_ERROR`, `NOT_FOUND`, `LLM_ERROR`. No rate limiting in current version.

**OpenAPI:** `api/gen/http/openapi.json`, `openapi.yaml`, `openapi3.json`, `openapi3.yaml` (in repo).

## See also

- [API examples](examples.md) — cURL examples for queries, reports, saved queries
- [Configuration](../configuration.md) — Environment variables (base URL, auth when available)
- [Deployment](../reference/deployment.md) — Running the API (Docker, Kubernetes, Helm)
- [Embedded integration](../getting-started/embedded.md) — Use the API via library or middleware
- [Documentation index](../README.md)
