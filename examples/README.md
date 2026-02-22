# Examples

- **Library:** `library-usage/` — use PgQueryNarrative as a Go library (run query, optional report). Requires PostgreSQL and optionally LLM.
- **Embedded:** `gin-integration/`, `echo-integration/`, `chi-integration/` — minimal HTTP servers with narrative endpoints. Build: `go build -o bin/example-gin ./examples/gin-integration` (same for echo, chi).
- **Business scenario:** `business-scenario/` — run `bash run-business-scenario.sh` (server must be running; see [Quick start](../docs/getting-started/quickstart.md)).

See [docs/README.md](../docs/README.md) and [API examples](../docs/api/examples.md).
