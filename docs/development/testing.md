# Testing

Unit and integration tests. Unit tests: `test/unit/`; run by package or by test name.

## Running tests

**All unit tests:**

```bash
make test-unit
```

Or: `go test ./test/unit/... ./cmd/server/... ./pkg/narrative/... -v`

**By package:**

```bash
go test ./test/unit/app/catalog/... -v
go test ./test/unit/app/charts/... -v
go test ./test/unit/app/metrics/... -v
go test ./test/unit/app/llm/... -v
go test ./test/unit/app/queryrunner/... -v
go test ./test/unit/app/service/... -v
go test ./test/unit/app/story/... -v
go test ./test/unit/app/suggestions/... -v
go test ./test/unit/web/... -v
go test ./cmd/server/... -v
go test ./pkg/narrative/... -v
```

**Single test:** `go test ./test/unit/app/service/... -run TestBuildPerfSuggestions_LimitApplied -v`

**Integration tests** (require Docker): `go test ./test/integration/... -v`

**E2E tests:** `go test ./test/e2e/... -v`

## Test layout

| Package | What is tested |
|---------|-----------------|
| `test/unit/app/catalog` | Schema/catalog loader |
| `test/unit/app/charts` | Chart suggestions |
| `test/unit/app/metrics` | Period comparison, trend, anomalies, data quality |
| `test/unit/app/llm` | Narrative prompt builder |
| `test/unit/app/queryrunner` | SQL validation (schema, SELECT-only, disallowed keywords) |
| `test/unit/app/story` | Narrative sanitizer |
| `test/unit/app/service` | Perf suggestions, metrics-to-API conversion |
| `test/unit/app/suggestions` | Query suggestions (curated, limit) |
| `test/unit/web` | Report HTML |
| `cmd/server` | Request logging middleware |
| `pkg/narrative` | Client, run-query options, validation |

**Integration** (`test/integration/...`): query runner; schema and suggestions against real Postgres; reports List/Get. **E2E** (`test/e2e/...`): full HTTP API against real Postgres (queries, schema, suggestions, reports).

## QA checklist

- Chart suggestions; period comparison and time-series in reports; data quality and perf suggestions; narrative content (no spurious "previous period" for single-period queries).
- API: `POST /reports/generate`, `GET /reports/{id}`, `GET /reports`; metrics structure; errors (invalid SQL → 400, not found → 404, LLM failure → 500).

Example API checks: [API examples](../api/examples.md).

## See also

- [Development setup](setup.md) — Build, commands, workflow
- [API examples](../api/examples.md) — Example API calls for QA
- [Troubleshooting](../reference/troubleshooting.md) — Common issues
- [Documentation index](../README.md)
