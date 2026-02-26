# Configuration

Environment variables only. Sensible defaults for local use.

## Logging

| Variable | Default | Description |
|---------|---------|-------------|
| `LOG_DEBUG` | (empty) | `1` or `true` = verbose logging |

## Server

| Variable | Default | Description |
|---------|---------|-------------|
| `PGQUERYNARRATIVE_HOST` | `0.0.0.0` | Bind address |
| `PGQUERYNARRATIVE_PORT` | `8080` | Server port |
| `PGQUERYNARRATIVE_READ_TIMEOUT` | `15s` | Request read timeout |
| `PGQUERYNARRATIVE_WRITE_TIMEOUT` | `60s` | Response write timeout |

## Database

| Variable | Default | Description |
|---------|---------|-------------|
| `POSTGRES_IMAGE` | `postgres:18-alpine` | Docker Postgres image |
| `DATABASE_HOST` | `localhost` | Database host |
| `DATABASE_PORT` | `5432` | Database port |
| `DATABASE_NAME` | `pgquerynarrative` | Database name |
| `DATABASE_USER` | `pgquerynarrative_app` | Application user |
| `DATABASE_PASSWORD` | `pgquerynarrative_app` | Application password |
| `DATABASE_READONLY_USER` | `pgquerynarrative_readonly` | Read-only user |
| `DATABASE_READONLY_PASSWORD` | `pgquerynarrative_readonly` | Read-only password |
| `DATABASE_SSL_MODE` | `disable` | SSL mode (disable/require/verify-full) |
| `DATABASE_MAX_CONNECTIONS` | `10` | Max connection pool size |
| `QUERY_TIMEOUT` | `30s` | Query execution timeout |

## LLM {#llm}

Required for report generation. See [LLM setup](getting-started/llm-setup.md).

| Variable | Default | Description |
|---------|---------|-------------|
| `LLM_PROVIDER` | `ollama` | ollama \| gemini \| claude \| openai \| groq |
| `LLM_MODEL` | `llama3.2` | Model name |
| `LLM_BASE_URL` | `http://localhost:11434` | LLM API base URL. Docker: `http://host.docker.internal:11434` |
| `LLM_API_KEY` | (empty) | API key (cloud providers) |

## Embeddings (optional) {#embeddings}

Used for similar-query retrieval (`GET /suggestions/similar`) and RAG in report generation. When not set, those features are disabled.

| Variable | Default | Description |
|---------|---------|-------------|
| `EMBEDDING_BASE_URL` | (empty) | Embedding API URL. If empty and `LLM_PROVIDER=ollama`, defaults to `LLM_BASE_URL`. |
| `EMBEDDING_MODEL` | `nomic-embed-text` | Embedding model (e.g. Ollama `nomic-embed-text`). |

Ollama: `ollama pull nomic-embed-text`. See [Semantic search (pgvector)](reference/semantic-search-pgvector.md).

## MCP (Claude desktop / Cursor) {#mcp-claude-desktop--cursor}

1. Build: `make build-mcp` → `bin/mcp-server`.
2. Edit MCP config:
   - **Claude:** macOS `~/Library/Application Support/Claude/claude_desktop_config.json`; Windows `%APPDATA%\Claude\`; Linux `~/.config/Claude/`.
   - **Cursor:** Settings → MCP or the MCP config file.
3. Add under `mcpServers` (replace path):
   ```json
   "pgquerynarrative": {
     "command": "/FULL/PATH/TO/pgquerynarrative/bin/mcp-server"
   }
   ```
   If app is not at http://localhost:8080: `"env": { "PGQUERYNARRATIVE_URL": "http://localhost:8080" }`. See `config/mcp-example.json`.
4. Restart client. Tools: `run_query`, `generate_report`, `list_saved_queries`, `get_report`, `list_reports`.

## Metrics

| Variable | Default | Description |
|---------|---------|-------------|
| `PERIOD_TREND_THRESHOLD_PERCENT` | `0.5` | Min % change to label trend "up"/"down"; below = "flat". |

## Security

| Variable | Default | Description |
|---------|---------|-------------|
| `SECURITY_AUTH_ENABLED` | `false` | Enable auth (future) |

## Loading config

- **Env:** `export PGQUERYNARRATIVE_PORT=8081` then start.
- **.env:** Create `.env` in project root (gitignored); `export $(cat .env | xargs)` before starting. Do not commit secrets.
- **Docker Compose:** Set `environment` under `app` in `docker-compose.yml`.

## Production

Change default passwords; use SSL for DB (`DATABASE_SSL_MODE=require`); use secrets management. Recommended: `QUERY_TIMEOUT=60s`, `DATABASE_MAX_CONNECTIONS=50`.

Invalid config causes clear startup errors. See [Troubleshooting](reference/troubleshooting.md) for other issues.

## See also

- [Installation](getting-started/installation.md) — Prerequisites and run methods
- [Quick start](getting-started/quickstart.md) — Minimal run
- [LLM setup](getting-started/llm-setup.md) — Providers and MCP
- [Deployment](reference/deployment.md) — Docker, Kubernetes, Helm
- [API reference](api/README.md) — REST endpoints
- [Troubleshooting](reference/troubleshooting.md) — Common issues
- [Documentation index](README.md)
