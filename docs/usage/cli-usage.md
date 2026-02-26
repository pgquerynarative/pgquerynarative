# CLI usage

Command-line access to the API. The app must be running (`make start-docker` or `make start-local`). CLI runs in a container (Docker) or on the host (local binary).

## Commands

| Command | Description |
|---------|-------------|
| `make cli CMD='query "SQL"'` | Run query (optional limit: `query "SQL" 10`) |
| `make cli CMD='list'` | List saved queries |
| `make cli CMD='get "uuid"'` | Get saved query by ID |
| `make cli CMD='save "Name" "SQL"'` | Save query (optional: `"tags,a,b"`) |
| `make cli CMD='report "SQL"'` | Generate report (requires LLM) |

**Interactive:** `make cli-shell` then e.g. `pgquerynarrative query "SELECT * FROM demo.sales LIMIT 5"` or alias `pqn list`.

## Examples

```bash
make cli CMD='query "SELECT product_category, SUM(total_amount) FROM demo.sales GROUP BY product_category"'
make cli CMD='save "Top Products" "SELECT product_category, SUM(total_amount) FROM demo.sales GROUP BY product_category"'
make cli CMD='report "SELECT product_category, SUM(total_amount) FROM demo.sales GROUP BY product_category"'
```

## Environment (CLI)

| Variable | Default | Description |
|----------|---------|-------------|
| `PGQUERYNARRATIVE_API_URL` | `http://app:8080` | API base URL (use `http://localhost:8080` when CLI on host) |
| `PGQUERYNARRATIVE_FORMAT` | `table` | Output: `table` or `json` |

JSON output: `make cli CMD='query "SELECT 1"'` with `PGQUERYNARRATIVE_FORMAT=json` (when running CLI on host).

**Quoting:** Quote SQL in the outer command: `make cli CMD='query "SELECT * FROM demo.sales"'`. For single quotes in SQL: `'\''`.

## See also

- [API examples](../api/examples.md) — cURL equivalents
- [API reference](../api/README.md) — REST endpoints
- [Quick start](../getting-started/quickstart.md) — Start the app
- [Troubleshooting](../reference/troubleshooting.md) — Common issues
- [Documentation index](../README.md)
