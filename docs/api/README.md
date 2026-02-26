# API reference

REST API base: `http://localhost:8080/api/v1`. No authentication in current version.

## Queries

| Method | Path | Description |
|--------|------|-------------|
| POST | `/queries/run` | Body: `{"sql":"...", "limit": 100}`. Run read-only SQL. Returns optional `period_comparison`, `chart_suggestions`. |
| POST | `/queries/saved` | Body: `{"name","sql","tags"}`. Save query. |
| GET | `/queries/saved` | Query: `limit`, `offset`, `tags`. List saved. |
| GET | `/queries/saved/{id}` | Get saved query. |
| DELETE | `/queries/saved/{id}` | Delete saved query. |

## Schema

| Method | Path | Description |
|--------|------|-------------|
| GET | `/schema` | Allowed schemas, tables, columns. Used by MCP `get_schema`. |

## Suggestions

| Method | Path | Description |
|--------|------|-------------|
| GET | `/suggestions/queries` | Query: `intent`, `limit` (default 5). Suggested SQL (curated + saved). |
| GET | `/suggestions/similar` | Query: `text` (required), `limit` (default 5). Saved queries semantically similar to text (embedding-based). Requires [embeddings](../reference/semantic-search-pgvector.md) enabled. |

## Reports

| Method | Path | Description |
|--------|------|-------------|
| POST | `/reports/generate` | Body: `{"sql":"...", "saved_query_id": "uuid"}`. Generate report (requires [LLM](../getting-started/llm-setup.md)). |
| GET | `/reports/{id}` | Get report. Metrics: `time_series`, `data_quality`, `perf_suggestions`, `predictive_summary`, etc. |
| GET | `/reports` | Query: `limit`, `offset`, `saved_query_id`. List reports. |

## Errors

JSON: `{"name","message","code"}`. Codes: `VALIDATION_ERROR`, `TIMEOUT_ERROR`, `NOT_FOUND`, `LLM_ERROR`. No rate limiting in current version.

**OpenAPI:** `api/gen/http/openapi.json`, `openapi.yaml` (in repo).

## See also

- [API examples](examples.md) — cURL for queries, reports, saved queries
- [Configuration](../configuration.md) — Environment variables
- [Deployment](../reference/deployment.md) — Running the API
- [Embedded integration](../getting-started/embedded.md) — Library and middleware
- [Documentation index](../README.md)
