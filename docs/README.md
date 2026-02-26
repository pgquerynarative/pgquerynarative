# PgQueryNarrative documentation

PgQueryNarrative turns SQL results into business narratives using an LLM. Run read-only SQL against PostgreSQL, get metrics and chart suggestions, and generate narrative reports.

**Start here:** [Quick start](getting-started/quickstart.md) → [LLM setup](getting-started/llm-setup.md) (for reports) → [Configuration](configuration.md).

---

## Documentation index

### Getting started

| Doc | Description |
|-----|-------------|
| [Installation](getting-started/installation.md) | Prerequisites, database setup, run methods |
| [Quick start](getting-started/quickstart.md) | Run with Docker or local PostgreSQL |
| [LLM setup](getting-started/llm-setup.md) | Configure LLM for report generation (Ollama, OpenAI, Claude, Gemini, Groq) and MCP |
| [Embedded integration](getting-started/embedded.md) | Use as a Go library or mount HTTP in Chi, Gin, Echo |

### User guides

| Doc | Description |
|-----|-------------|
| [Configuration](configuration.md) | Environment variables (server, database, LLM, embeddings) |
| [CLI usage](usage/cli-usage.md) | Command-line interface for queries, saved queries, reports |

### API

| Doc | Description |
|-----|-------------|
| [API reference](api/README.md) | REST endpoints, request/response, error codes |
| [API examples](api/examples.md) | cURL examples for run, save, reports |

### Reference

| Doc | Description |
|-----|-------------|
| [Deployment](reference/deployment.md) | Docker, Kubernetes, Helm |
| [Operations](reference/operations.md) | Monitoring, health checks, runbooks |
| [Troubleshooting](reference/troubleshooting.md) | Common issues and fixes |
| [PostgreSQL extension](reference/postgres-extension.md) | Call the API from SQL via `CREATE EXTENSION pgquerynarrative` |
| [Semantic search (pgvector)](reference/semantic-search-pgvector.md) | Embeddings, similar-query search, RAG |

### Development

| Doc | Description |
|-----|-------------|
| [Development setup](development/setup.md) | Build, test, codegen, workflow |
| [Testing](development/testing.md) | Unit, integration, E2E tests |

**Contributing & security:** [.github/CONTRIBUTING.md](../.github/CONTRIBUTING.md), [.github/SECURITY.md](../.github/SECURITY.md).

**Changelog:** [CHANGELOG.md](../CHANGELOG.md).
