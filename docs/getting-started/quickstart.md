# Quick start

Minimal steps to run PgQueryNarrative.

## Prerequisites

Docker, or PostgreSQL 16+ and Go 1.24+.

## Docker

```bash
git clone https://github.com/your-org/pgquerynarrative.git
cd pgquerynarrative
make start-docker
```

Uses root `docker-compose.yml` (PostgreSQL + app). App: http://localhost:8080. Production image: [Deployment](../reference/deployment.md).

## Local PostgreSQL

```bash
git clone https://github.com/your-org/pgquerynarrative.git
cd pgquerynarrative
pg_isready   # ensure Postgres is running
make start-local
```

## Next steps

- **Web UI:** http://localhost:8080
- **CLI:** `make cli CMD='query "SELECT * FROM demo.sales LIMIT 5"'` — [CLI usage](../usage/cli-usage.md)
- **API:** [API examples](../api/examples.md)

## See also

- [Installation](installation.md) — Prerequisites and detailed setup
- [LLM setup](llm-setup.md) — Report generation and MCP
- [Configuration](../configuration.md) — Environment variables
- [Deployment](../reference/deployment.md) — Docker build, Kubernetes, Helm
- [Troubleshooting](../reference/troubleshooting.md) — Common issues
- [Documentation index](../README.md)
