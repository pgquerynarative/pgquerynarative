# Deployment

How to build, run, and deploy PgQueryNarrative with Docker, Kubernetes, or Helm. For a first-time run, see [Quick start](../getting-started/quickstart.md) or [Installation](../getting-started/installation.md).

**Note:** `make start-docker` (see [Quick start](../getting-started/quickstart.md)) uses the root `docker-compose.yml` and root `Dockerfile` for development. The sections below describe the **production** image (`deploy/docker/Dockerfile`), production compose (`deploy/docker/docker-compose.yml`), and Kubernetes/Helm.

## Docker

### Build the image

From the repository root (build context must be the repo root):

```bash
docker build -f deploy/docker/Dockerfile -t pgquerynarrative:latest .
```

The production Dockerfile is multi-stage: it builds the app and a migrate binary, then produces a minimal Alpine image with only the server binary, migrations, and entrypoint. No Go toolchain in the final image.

### Run with Docker Compose

Compose starts PostgreSQL and the app, runs migrations on startup, and optionally seeds demo data:

```bash
docker compose -f deploy/docker/docker-compose.yml up -d
```

Or build and run in one step:

```bash
docker compose -f deploy/docker/docker-compose.yml up -d --build
```

Set environment variables or use a `.env` file (when running Compose from the repo root, `.env` in the repo root is loaded). Important variables:

- `DATABASE_*` — Postgres connection (defaults point to the `postgres` service).
- `LLM_*` — LLM provider and model (default: Ollama at `http://host.docker.internal:11434`).
- `PGQUERYNARRATIVE_SEED=true` — seed demo data on first run (optional).

The app container waits for Postgres to be healthy, runs migrations, then starts the server. API: `http://localhost:8080`.

### Using a pre-built image

If you build the image elsewhere (e.g. CI) and push to a registry, point Compose at it:

```yaml
# In docker-compose override or env
# Build section can be replaced with:
services:
  app:
    image: your-registry/pgquerynarrative:1.0.0
```

---

## Kubernetes

Manifests are in `deploy/kubernetes/`. PostgreSQL is assumed to be external (managed instance or your own); the app connects via `DATABASE_HOST` and credentials from a Secret.

### Prerequisites

- Kubernetes cluster and `kubectl` configured.
- PostgreSQL reachable from the cluster (in-cluster service or external host). Create the database and roles (e.g. `pgquerynarrative_app`, `pgquerynarrative_readonly`) and run migrations once if the DB is empty.

### Apply order

1. Create the namespace (optional; you can use `default` and adjust resource names).
2. Create the Secret with real credentials (do not commit them).
3. Create the ConfigMap (or edit `configmap.yaml` to match your DB host and LLM).
4. Create Deployment and Service.
5. Optionally create the Ingress.

```bash
kubectl apply -f deploy/kubernetes/namespace.yaml
kubectl apply -f deploy/kubernetes/secret.yaml   # edit secret.yaml first
kubectl apply -f deploy/kubernetes/configmap.yaml
kubectl apply -f deploy/kubernetes/deployment.yaml
kubectl apply -f deploy/kubernetes/service.yaml
kubectl apply -f deploy/kubernetes/ingress.yaml   # optional
```

Set `image` in `deployment.yaml` to your image (e.g. `your-registry/pgquerynarrative:1.0.0`). Ensure `configmap.yaml` has the correct `DATABASE_HOST` (e.g. a K8s Service name or external hostname).

### Access

- Without Ingress: `kubectl port-forward -n pgquerynarrative svc/pgquerynarrative 8080:8080` then open `http://localhost:8080`.
- With Ingress: configure your ingress controller and DNS for the host in `ingress.yaml`.

---

## Helm

The chart is in `deploy/helm/pgquerynarrative/`. It deploys the same app with ConfigMap, Secret, Deployment, Service, and optional Ingress, driven by `values.yaml`.

### Install

From the repo root:

```bash
helm install pgqn ./deploy/helm/pgquerynarrative -n pgquerynarrative --create-namespace
```

Override values (e.g. image, database host, passwords):

```bash
helm install pgqn ./deploy/helm/pgquerynarrative -n pgquerynarrative --create-namespace \
  --set image.repository=your-registry/pgquerynarrative \
  --set image.tag=1.0.0 \
  --set database.host=your-postgres-host \
  --set secret.databasePassword=xxx \
  --set secret.databaseReadonlyPassword=xxx
```

Or use a custom values file:

```bash
helm install pgqn ./deploy/helm/pgquerynarrative -n pgquerynarrative --create-namespace -f my-values.yaml
```

### Upgrade / uninstall

```bash
helm upgrade pgqn ./deploy/helm/pgquerynarrative -n pgquerynarrative
helm uninstall pgqn -n pgquerynarrative
```

### Chart values

See `deploy/helm/pgquerynarrative/values.yaml`. Important keys:

- **image** — repository, tag, pullPolicy.
- **database** — host, port, name, user, sslMode, maxConnections, readonlyUser.
- **secret** — databasePassword, databaseReadonlyPassword (set via `--set` or a private values file).
- **llm** — provider, model, baseURL (and optional apiKey for cloud LLMs).
- **ingress.enabled** — set to `true` and set **ingress.host** to expose via Ingress.
- **seed** — set to `true` to run demo seed on startup (optional).

---

## Summary

| Method        | Path                          | Use case                    |
|---------------|-------------------------------|-----------------------------|
| Docker Compose| `deploy/docker/`              | Single host, dev or staging |
| Kubernetes    | `deploy/kubernetes/`           | Any cluster, raw manifests  |
| Helm          | `deploy/helm/pgquerynarrative/`| Parameterized K8s install   |

**Repository layout:** `deploy/docker/` (Dockerfile, entrypoint, production compose), `deploy/kubernetes/` (namespace, configmap, secret, deployment, service, ingress), `deploy/helm/pgquerynarrative/` (Chart.yaml, values.yaml, templates).

## See also

- [Configuration](../configuration.md) — Environment variables (database, LLM, server)
- [Operations](operations.md) — Monitoring, health checks, and runbooks
- [Docker resources](docker-resources.md) — Image and container resource usage
- [Quick start](../getting-started/quickstart.md) — Minimal run with Docker or local Postgres
- [Installation](../getting-started/installation.md) — Prerequisites and setup
- [API reference](../api/README.md) — REST endpoints
- [Troubleshooting](troubleshooting.md) — Common issues
- [Documentation index](../README.md)
