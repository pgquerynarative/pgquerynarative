# Operations

Monitoring and runbooks for production. Deployment: [Deployment](deployment.md). Common issues: [Troubleshooting](troubleshooting.md).

## Monitoring

### Health checks

The app does not expose `/health` or `/ready`. Deployments use an API endpoint for liveness/readiness:

| Environment | Probe | Endpoint | Purpose |
|-------------|--------|----------|---------|
| Docker Compose | healthcheck | `GET http://localhost:8080/api/v1/queries/saved` | Container health |
| Kubernetes / Helm | livenessProbe, readinessProbe | `GET /api/v1/queries/saved` (port 8080) | Restart if unresponsive; remove from service if not ready |

That endpoint requires a working DB connection (read-only pool). If the DB is down, the probe fails. **Limitation:** Probe does not check LLM; report generation can fail while the app is "healthy."

### What to monitor

| Area | Watch | How |
|------|--------|-----|
| Application | Process up, HTTP responding | Liveness/readiness above; alert on N consecutive failures. |
| Database | Connectivity, pool exhaustion | App logs (connection errors, pool health). Postgres metrics if available. |
| LLM | Availability, latency, errors | App logs (report errors, timeouts). Report success rate. |
| Logs | Errors, slow requests | Stdout/stderr; aggregator (Loki, CloudWatch). Look for pool/report failures, 5xx. |
| Resources | Memory, CPU | Container/pod metrics (cAdvisor, Prometheus). |

## Runbooks

### Deploy (standard)

1. Build and push: `docker build -f deploy/docker/Dockerfile -t your-registry/pgquerynarrative:<tag> .` then `docker push ...`.
2. **Compose:** `docker compose -f deploy/docker/docker-compose.yml up -d` (or set `image:` to new tag).
3. **Kubernetes:** Update `image` in `deployment.yaml`; `kubectl apply -f deploy/kubernetes/deployment.yaml`.
4. **Helm:** `helm upgrade pgqn ./deploy/helm/pgquerynarrative -n pgquerynarrative --set image.tag=<tag>`.
5. Verify: UI or `GET /api/v1/queries/saved`; check logs.

### Rollback

1. **Compose:** Revert image; `docker compose -f deploy/docker/docker-compose.yml up -d`.
2. **Kubernetes:** Revert `deployment.yaml` and apply, or `kubectl rollout undo deployment/pgquerynarrative -n pgquerynarrative`.
3. **Helm:** `helm rollback pgqn -n pgquerynarrative` or upgrade with previous `image.tag`.
4. Confirm: same verification as deploy.

### Incident: Database unreachable

**Symptoms:** Probe failures, "connection refused" or pool health errors in logs, 503 or failed API calls.

**Actions:** Confirm Postgres is running and reachable (network, firewall, credentials). Compose: check `postgres` service and logs; restart if needed. External Postgres: check instance, connectivity, credentials. After DB is back, app should pass readiness. See [Troubleshooting – Database](troubleshooting.md#database).

### Incident: LLM unreachable or report generation failing

**Symptoms:** Report generation fails or times out; LLM errors in logs.

**Actions:** Check `LLM_BASE_URL`, `LLM_PROVIDER`, `LLM_MODEL`, `LLM_API_KEY`. [Configuration](../configuration.md), [LLM setup](../getting-started/llm-setup.md). Ollama: ensure `ollama serve` and model pulled; Docker: `http://host.docker.internal:11434`. Cloud: verify API key, quota, region. See [Troubleshooting – Reports and LLM](troubleshooting.md#reports-and-llm).

### Incident: High memory or CPU

**Symptoms:** OOM, slow responses, resource alerts.

**Actions:** Check usage (e.g. `kubectl top pod`). Correlate with load (concurrent reports). Short term: scale replicas or increase limits. Long term: tune `DATABASE_MAX_CONNECTIONS`, consider rate limiting or queueing for reports.

### Incident: Bad deployment (wrong config or image)

**Symptoms:** Wrong behavior after deploy.

**Actions:** Roll back (see Rollback above). Fix ConfigMap/Secret or values; redeploy with correct image and config; verify.

## See also

- [Deployment](deployment.md) — Build, Docker, Kubernetes, Helm
- [Troubleshooting](troubleshooting.md) — Common issues
- [Configuration](../configuration.md) — Environment variables
- [Documentation index](../README.md)
