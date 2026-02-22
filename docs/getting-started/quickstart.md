# Quick start

## Prerequisites

- Docker, or PostgreSQL 16+ and Go 1.24+

## Docker

```bash
git clone https://github.com/your-org/pgquerynarrative.git
cd pgquerynarrative
make start-docker
```

This uses the root `docker-compose.yml` (PostgreSQL + app from root `Dockerfile`). App: `http://localhost:8080`. For a production-style image and compose, see [Deployment](../reference/deployment.md).

## Local PostgreSQL

```bash
git clone https://github.com/your-org/pgquerynarrative.git
cd pgquerynarrative
pg_isready   # ensure Postgres is running
make start-local
```

## First steps

- **Web:** http://localhost:8080
- **CLI:** `make cli CMD='query "SELECT * FROM demo.sales LIMIT 5"'` ([CLI usage](../usage/cli-usage.md))
- **API:** [API examples](../api/examples.md)

## See also

- [Installation](installation.md) — Prerequisites and detailed setup
- [LLM setup](llm-setup.md) — Report generation and MCP
- [Configuration](../configuration.md) — Environment variables
- [Deployment](../reference/deployment.md) — Docker build/compose, Kubernetes, Helm
- [API examples](../api/examples.md) — cURL for queries and reports
- [Troubleshooting](../reference/troubleshooting.md) — Common issues
- [Documentation index](../README.md)
