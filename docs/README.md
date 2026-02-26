# PgQueryNarrative documentation

PgQueryNarrative turns SQL query results into business narratives with AI. This documentation covers installation, configuration, the API, deployment, and development.

**New to the project?** Start with [Quick start](getting-started/quickstart.md) (Docker or local PostgreSQL), then [LLM setup](getting-started/llm-setup.md) for report generation. For deployment to Docker/Kubernetes/Helm, see [Deployment](reference/deployment.md).

**Where to find:** Run the app → [Quick start](getting-started/quickstart.md), [Installation](getting-started/installation.md). Configure env → [Configuration](configuration.md). Call the API → [API reference](api/README.md), [API examples](api/examples.md). Embed in Go → [Embedded integration](getting-started/embedded.md). Deploy → [Deployment](reference/deployment.md). Build/test → [Development setup](development/setup.md), [Testing](development/testing.md). Problems → [Troubleshooting](reference/troubleshooting.md).

---

## Getting started

| Document | Description |
|----------|-------------|
| [Installation](getting-started/installation.md) | Prerequisites, database setup, and running the application |
| [Quick start](getting-started/quickstart.md) | Minimal steps to run with Docker or local PostgreSQL |
| [LLM setup](getting-started/llm-setup.md) | Configure an LLM provider for report generation (Ollama, OpenAI, Claude, Gemini, Groq) and MCP |
| [Embedded integration](getting-started/embedded.md) | Use as a library or mount HTTP endpoints in Chi, Gin, or Echo |

---

## User guides

| Document | Description |
|----------|-------------|
| [Configuration](configuration.md) | Environment variables (server, database, LLM, metrics) |
| [CLI usage](usage/cli-usage.md) | Command-line interface for queries, saved queries, and reports |

---

## API

| Document | Description |
|----------|-------------|
| [API reference](api/README.md) | REST endpoints, request/response formats, error codes |
| [API examples](api/examples.md) | cURL examples for running queries, saving queries, and generating reports |

---

## Reference

| Document | Description |
|----------|-------------|
| [Deployment](reference/deployment.md) | Docker build/compose, Kubernetes manifests, Helm chart |
| [Operations](reference/operations.md) | Monitoring, health checks, and runbooks (deploy, rollback, incidents) |
| [Troubleshooting](reference/troubleshooting.md) | Common issues and solutions |
| [PostgreSQL extension](reference/postgres-extension.md) | Run queries and generate reports from SQL |

---

## Development

| Document | Description |
|----------|-------------|
| [Development setup](development/setup.md) | Build, test, code generation, and workflow |
| [Testing](development/testing.md) | Running unit and integration tests, QA checklist |

For contributing and security, see the [.github](https://github.com/pgquerynarrative/pgquerynarrative/tree/main/.github) directory (CONTRIBUTING.md, SECURITY.md).

---

## Changelog

Release history and unreleased changes: [CHANGELOG.md](../CHANGELOG.md).
