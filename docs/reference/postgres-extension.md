# PostgreSQL extension

Call PgQueryNarrative from SQL: run queries, generate reports, and list saved queries. Requires the PgQueryNarrative service to be running and PostgreSQL 16+. Extension files live in `infra/postgres-extension/`.

## Easy install (Docker)

From the repo root:

```bash
make setup-extension-docker
```

Starts Postgres, inits the DB, runs migrations, installs the extension, and seeds demo data. If Postgres is already running, use `make install-extension-docker` to install or reinstall the extension only.

## Install (summary)

- **Docker:** `make setup-extension-docker` (full setup) or `make install-extension-docker` (extension only).
- **Local:** `make install-extension` (set `PGHOST`, `PGPORT`, `PGQUERYNARRATIVE_DB`, `PGQUERYNARRATIVE_USER`), or run `pgquerynarrative--1.0.sql` and optionally `pgquerynarrative--1.0--with-http.sql` (requires the [http](https://github.com/pramsey/pgsql-http) extension) with `psql`.

## Configuration

```sql
SELECT pgquerynarrative_set_api_url('http://localhost:8080');
SELECT pgquerynarrative_get_api_url();
```

## Functions

| Function | Description |
|----------|-------------|
| `pgquerynarrative_run_query(query_sql, row_limit)` | Run a read-only query; returns JSON. Default limit 100. |
| `pgquerynarrative_generate_report(query_sql)` | Generate a narrative report (JSON). |
| `pgquerynarrative_list_saved(limit, offset)` | List saved queries (JSON). |

## Example

```sql
SELECT pgquerynarrative_run_query(
  'SELECT product_category, SUM(total_amount) FROM demo.sales GROUP BY product_category',
  10
);
SELECT pgquerynarrative_generate_report(
  'SELECT product_category, SUM(total_amount) FROM demo.sales GROUP BY product_category'
);
```

## Troubleshooting

- **Placeholder or empty result:** Install the HTTP extension and apply `pgquerynarrative--1.0--with-http.sql`.
- **Connection errors:** Verify the PgQueryNarrative service is running and the API URL is correct (use `http://host.docker.internal:8080` when Postgres is in Docker and the app is on the host).
- **Role/database issues:** See [Troubleshooting](troubleshooting.md).

## See also

- [API reference](../api/README.md) — REST endpoints called by the extension
- [Configuration](../configuration.md) — Server and database config
- [Deployment](deployment.md) — Running the service (Docker, Kubernetes, Helm)
- [Troubleshooting](troubleshooting.md) — Common issues
- [Documentation index](../README.md)
