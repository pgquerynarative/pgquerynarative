# Installation

Prerequisites and how to run PgQueryNarrative (Docker or from source).

## Prerequisites

- **Go 1.24+** (for building from source)
- **PostgreSQL 16+** or Docker

For report generation: [LLM setup](llm-setup.md).

## Docker (recommended)

```bash
git clone https://github.com/your-org/pgquerynarrative.git
cd pgquerynarrative
make start-docker
```

Starts PostgreSQL, runs migrations and seed, then the app. API: http://localhost:8080.

## Local (from source)

1. **Install Go and PostgreSQL** (e.g. macOS: `brew install go postgresql@18`).

2. **Clone and setup:**
   ```bash
   git clone https://github.com/your-org/pgquerynarrative.git
   cd pgquerynarrative
   make setup
   make generate
   make build
   ```

3. **Database:** With Postgres running:
   ```bash
   make db-init
   make migrate
   make seed
   ```

4. **Run:** `make run` or `./bin/server`

## Verify

```bash
curl http://localhost:8080/api/v1/queries/saved
```

## PostgreSQL versions

Supported: 16, 17, 18. Docker default: `postgres:18-alpine`. Override: `POSTGRES_IMAGE=postgres:17-alpine make start-docker`.

## See also

- [Quick start](quickstart.md) — Minimal run steps
- [Configuration](../configuration.md) — Environment variables
- [Deployment](../reference/deployment.md) — Docker, Kubernetes, Helm
- [Troubleshooting](../reference/troubleshooting.md) — Common issues
- [Documentation index](../README.md)
