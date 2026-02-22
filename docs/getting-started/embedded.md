# Embedded integration

Use PgQueryNarrative inside your own Go service by creating a `narrative.Client` and either calling it directly (library mode) or mounting its HTTP endpoints with the provided middleware (Chi, Gin, Echo).

## Library usage

Create a client from config, run queries and generate reports in code, then close the client when done.

```go
import (
    "github.com/pgquerynarrative/pgquerynarrative/pkg/narrative"
    "github.com/pgquerynarrative/pgquerynarrative/app/config"
)

cfg := narrative.FromAppConfig(config.Load())
client, err := narrative.NewClient(ctx, cfg)
if err != nil { ... }
defer client.Close()

result, err := client.RunQuery(ctx, "SELECT ... FROM demo.sales", 100)
report, err := client.GenerateReport(ctx, sql)
schema, err := client.GetSchema(ctx)
```

See `examples/library-usage/basic.go`.

## HTTP middleware (Chi, Gin, Echo)

Mount narrative endpoints on your existing router so you get query/run, report/generate, schema, and suggestions without writing handlers.

**Chi:**

```go
import narrativemw "github.com/pgquerynarrative/pgquerynarrative/pkg/narrative/middleware"

r := chi.NewRouter()
narrativemw.MountChi(r, client, "/api")
```

**Gin:**

```go
r := gin.Default()
narrativemw.MountGin(r, client, "/api")
```

**Echo:**

```go
e := echo.New()
narrativemw.MountEcho(e, client, "/api")
```

Mounted routes (with prefix `/api`):

| Method | Path | Description |
|--------|------|-------------|
| POST | /api/query/run | Body: `{"sql":"...", "limit": N}`. Run read-only SQL. |
| POST | /api/report/generate | Body: `{"sql":"..."}`. Generate narrative report (requires LLM). |
| GET | /api/schema | Allowed schemas, tables, columns. |
| GET | /api/suggestions/queries | Query: `intent`, `limit`. Suggested SQL (curated + saved). |

Use an empty prefix (`""`) to mount at the root (e.g. `/query/run`).

## Examples

- `examples/library-usage/basic.go` – client only, no HTTP
- `examples/chi-integration/main.go` – Chi + middleware
- `examples/gin-integration/main.go` – Gin + middleware
- `examples/echo-integration/main.go` – Echo + middleware

Build and run (set `DATABASE_*` and `LLM_*` as needed):

```bash
go build -o bin/example-chi ./examples/chi-integration
./bin/example-chi
```

Config is the same as the standalone server: `DATABASE_*`, `LLM_*`, etc. See [Configuration](../configuration.md).

## See also

- [Configuration](../configuration.md) — Database and LLM environment variables
- [API reference](../api/README.md) — REST endpoints (aligned with middleware routes)
- [Deployment](../reference/deployment.md) — Running in Docker or Kubernetes
- [Documentation index](../README.md)
