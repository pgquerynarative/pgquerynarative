# Troubleshooting

Common issues and fixes.

## Environment and dependencies

| Issue | Solution |
|-------|----------|
| **Docker not found** | Install [Docker Desktop](https://www.docker.com/products/docker-desktop). Verify: `docker info`. |
| **make: command not found** | Install Make (e.g. `brew install make` on macOS). |
| **Port 8080 in use** | Use another port: `PGQUERYNARRATIVE_PORT=8081 make start-docker` or `make start-local`. |

## Database {#database}

| Issue | Solution |
|-------|----------|
| **PostgreSQL connection refused** | **Docker:** `make start-docker`. **Local:** start Postgres (e.g. `brew services start postgresql@18`), then `make start-local`. |
| **Role does not exist / permission denied** | Run `make db-init` then `make migrate`, or `make start-local` once. If `demo.sales` denied: grant SELECT to `pgquerynarrative_readonly` on the table. |

## Reports and LLM {#reports-and-llm}

| Issue | Solution |
|-------|----------|
| **Failed to parse narrative JSON** | LLM output may be truncated. Ensure Ollama is running and model pulled (`ollama serve`, `ollama pull llama3.2`). Try a larger model or restart Ollama. |
| **Report generation fails or times out** | See [LLM setup](../getting-started/llm-setup.md). App in Docker, Ollama on host: `LLM_BASE_URL=http://host.docker.internal:11434`. |
| **Narrative shows wrong number scale** | Ensure prompt and data use correct magnitude; check app version. |

## Extension (PostgreSQL)

| Issue | Solution |
|-------|----------|
| **CREATE EXTENSION pgquerynarrative fails** | Copy extension files first: [PostgreSQL extension](postgres-extension.md). Local: `make install-extension`. Docker: `make install-extension-docker`. Then run `CREATE EXTENSION pgquerynarrative;` in psql. |

## See also

- [Configuration](../configuration.md) — Environment variables
- [Installation](../getting-started/installation.md) — Prerequisites and setup
- [Quick start](../getting-started/quickstart.md) — Minimal run
- [LLM setup](../getting-started/llm-setup.md) — Report generation and MCP
- [Deployment](deployment.md) — Docker, Kubernetes, Helm
- [Operations](operations.md) — Monitoring and runbooks
- [Documentation index](../README.md)
