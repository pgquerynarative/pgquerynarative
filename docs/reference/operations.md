# Operations

Monitoring and runbooks for running PgQueryNarrative in production. For deployment steps, see [Deployment](deployment.md). For common issues, see [Troubleshooting](troubleshooting.md).

## Monitoring

### Health checks in use

The app does not expose dedicated `/health` or `/ready` endpoints. Deployments use an existing API endpoint for liveness and readiness:

| Environment | Probe type | Endpoint | Purpose |
|-------------|------------|----------|---------|
| Docker Compose | `healthcheck` | `GET http://localhost:8080/api/v1/queries/saved` | Container health |
| Kubernetes / Helm | `livenessProbe` | `GET /api/v1/queries/saved` (port 8080) | Restart if unresponsive |
| Kubernetes / Helm | `readinessProbe` | `GET /api/v1/queries/saved` (port 8080) | Remove from service if not ready |

That endpoint requires a working database connection (read-only pool). If the database is down, the probe fails and the orchestrator can act (e.g. restart pod, mark not ready).

**Limitation:** The probe does not check LLM availability. Report generation will fail if the LLM is unreachable, but the app will still be considered "healthy." For stricter readiness, you can add a dedicated readiness endpoint later that checks both DB and LLM.

### What to monitor

| Area | What to watch | How |
|------|----------------|-----|
| **Application** | Process up, HTTP responding | Existing liveness/readiness (above). Optionally: alert if probe fails for N consecutive checks. |
| **Database** | Connectivity, pool exhaustion | App logs (`database pool health check failed`, connection errors). Optionally: Postgres metrics (connections, replication lag). |
| **LLM** | Availability, latency, errors | App logs (report generation errors, timeouts). Optionally: scrape LLM provider metrics or track report success rate. |
| **Logs** | Errors, slow requests | Stdout/stderr; send to your log aggregator (e.g. Loki, CloudWatch). Look for `ErrPoolHealthCheckFailed`, report generation failures, 5xx. |
| **Resources** | Memory, CPU | Container/pod metrics (e.g. cAdvisor, Prometheus). Default limits: 256Mi memory, 500m CPU. |

### Optional: metrics and dashboards

The application does not expose Prometheus or other metrics endpoints. If you add them later, useful signals include:

- Request count and latency by endpoint (e.g. `/api/v1/queries/run`, `/api/v1/reports/generate`)
- Report generation success vs failure and latency
- Database pool usage (in use vs idle connections)
- Optional: LLM token usage or cost if provided by your LLM stack

You can already build dashboards from:

- **Orchestrator:** pod/container restarts, probe failures, resource usage (Kubernetes metrics).
- **Logs:** error rates, patterns (e.g. "failed to generate report").
- **Database:** Postgres metrics (connections, queries, replication) if you run your own Postgres or use a managed service with metrics.

---

## Runbooks

### Deploy (standard)

1. Build and push the image (see [Deployment](deployment.md)):
   ```bash
   docker build -f deploy/docker/Dockerfile -t your-registry/pgquerynarrative:<tag> .
   docker push your-registry/pgquerynarrative:<tag>
   ```
2. **Docker Compose:** `docker compose -f deploy/docker/docker-compose.yml up -d` (or with `image:` set to the new tag).
3. **Kubernetes:** Update `image` in `deploy/kubernetes/deployment.yaml` and `kubectl apply -f deploy/kubernetes/deployment.yaml`.
4. **Helm:** `helm upgrade pgqn ./deploy/helm/pgquerynarrative -n pgquerynarrative --set image.tag=<tag>` (and `image.repository` if needed).
5. Verify: hit the UI or `GET /api/v1/queries/saved` and confirm responses; check logs for errors.

### Rollback

1. **Docker Compose:** Change `image` (or Dockerfile) back to the previous tag and run `docker compose -f deploy/docker/docker-compose.yml up -d`.
2. **Kubernetes:** Revert `deployment.yaml` to the previous image and `kubectl apply -f deploy/kubernetes/deployment.yaml`, or use `kubectl rollout undo deployment/pgquerynarrative -n pgquerynarrative`.
3. **Helm:** `helm rollback pgqn -n pgquerynarrative` (or upgrade with the previous `image.tag`).
4. Confirm: same verification as deploy; check that reports and queries work again.

### Incident: Database unreachable

**Symptoms:** Probe failures, "connection refused" or "database pool health check failed" in logs, 503 or failed API calls.

**Actions:**

1. Confirm Postgres is running and reachable from the app (network, firewall, credentials).
2. If using Docker Compose: `docker compose -f deploy/docker/docker-compose.yml ps` and logs for the `postgres` service; restart if needed.
3. If using external Postgres: check the instance status, connectivity from the app network, and credentials in ConfigMap/Secret or env.
4. After DB is back, the app should pass readiness again; no app restart required unless you restarted Postgres and the app lost connections.

**See:** [Troubleshooting — Database](troubleshooting.md#database).

### Incident: LLM unreachable or report generation failing

**Symptoms:** Report generation fails or times out; errors in logs about LLM (e.g. connection refused, timeout, 5xx from provider).

**Actions:**

1. Check LLM configuration: `LLM_BASE_URL`, `LLM_PROVIDER`, `LLM_MODEL`, and `LLM_API_KEY` (if required). See [Configuration](../configuration.md) and [LLM setup](../getting-started/llm-setup.md).
2. If Ollama: ensure `ollama serve` is running and the model is pulled (`ollama pull <model>`). From Docker, use `http://host.docker.internal:11434` for host Ollama.
3. For cloud LLMs: verify API key, quota, and region/endpoint.
4. Retry a single report; if timeouts persist, consider a larger timeout or a smaller/faster model.

**Note:** Liveness/readiness do not depend on the LLM; the app stays "healthy" and only report generation is affected.

**See:** [Troubleshooting — Reports and LLM](troubleshooting.md#reports-and-llm).

### Incident: High memory or CPU

**Symptoms:** OOM kills, slow responses, or resource alerts on the app container/pod.

**Actions:**

1. Check current usage: `kubectl top pod -n pgquerynarrative` or Docker stats.
2. Correlate with load: many concurrent report generations (especially with large result sets or long narratives) can increase memory and CPU.
3. Short term: scale replicas (Helm: `replicaCount`) to spread load; increase memory/CPU limits if justified.
4. Long term: tune `DATABASE_MAX_CONNECTIONS`, consider rate limiting or queueing for report generation, and profile heavy endpoints.

**See:** [Docker resources](docker-resources.md) for default resource usage.

### Incident: Bad deployment (wrong config or image)

**Symptoms:** Wrong behavior after deploy (e.g. wrong DB, wrong LLM, broken UI).

**Actions:**

1. Roll back using the [Rollback](#rollback) steps above.
2. Fix ConfigMap/Secret or `values.yaml` (e.g. `DATABASE_HOST`, `LLM_BASE_URL`, secrets).
3. Redeploy with the correct image and config; verify again.

---

## See also

- [Deployment](deployment.md) — Build, Docker, Kubernetes, Helm
- [Troubleshooting](troubleshooting.md) — Common issues and solutions
- [Configuration](../configuration.md) — Environment variables
- [Docker resources](docker-resources.md) — Image and container resources
