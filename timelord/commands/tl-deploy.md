---
name: tl-deploy
description: Get deployment guidance for Temporal clusters on Kubernetes
arguments:
  - name: environment
    description: "Target environment: local, development, staging, or production"
    required: false
---

# Deploy Temporal Cluster

Get guidance for deploying Temporal clusters to Kubernetes environments.

## Usage

```
/tl-deploy [environment]
```

## Environments

| Environment | Characteristics |
|-------------|-----------------|
| `local` | Docker Compose or minikube, minimal resources |
| `development` | Kubernetes cluster, 1 replica per service |
| `staging` | Production-like, reduced resources |
| `production` | Full HA, monitoring, security |

## Quick Local Setup

For local development with Docker Compose:

```bash
git clone https://github.com/temporalio/docker-compose.git
cd docker-compose
docker-compose up -d
```

Access:

- Server: `localhost:7233`
- Web UI: `http://localhost:8080`

## Kubernetes Deployment

### Prerequisites

1. Kubernetes cluster (1.24+)
2. Helm 3.x installed
3. kubectl configured
4. PostgreSQL database (can be provisioned by Helm)

### Add Helm Repository

```bash
helm repo add temporal https://go.temporal.io/helm-charts
helm repo update
```

### Deploy

**Development:**

```bash
helm install temporal temporal/temporal \
  --namespace temporal \
  --create-namespace \
  --set server.config.numHistoryShards=128 \
  --set cassandra.enabled=false \
  --set postgresql.enabled=true \
  --set elasticsearch.enabled=false
```

**Production:**

```bash
helm install temporal temporal/temporal \
  --namespace temporal \
  --create-namespace \
  -f values-production.yaml
```

## Deployment Checklist

### Pre-Deployment

- [ ] Kubernetes cluster provisioned and accessible
- [ ] PostgreSQL database created
- [ ] Database credentials stored in Kubernetes secret
- [ ] Helm repository added
- [ ] Resource requirements calculated
- [ ] History shard count determined (cannot change later)

### Deployment

- [ ] Create namespace
- [ ] Deploy Helm chart
- [ ] Verify all pods running
- [ ] Run schema migrations (if using external database)
- [ ] Create default namespace
- [ ] Verify cluster health

### Post-Deployment

- [ ] Configure monitoring (Prometheus/Grafana)
- [ ] Set up alerting
- [ ] Configure mTLS (production)
- [ ] Document access procedures
- [ ] Test workflow execution

## Verification Commands

```bash
# Check pod status
kubectl get pods -n temporal

# Check service endpoints
kubectl get svc -n temporal

# Port-forward for local access
kubectl port-forward svc/temporal-frontend 7233:7233 -n temporal

# Verify cluster health
temporal operator cluster health

# Create default namespace
temporal operator namespace create --namespace default --retention 3d
```

## Common Issues

**Pods stuck in Pending:**

- Check resource quotas
- Verify node capacity
- Check PVC provisioning

**Database connection errors:**

- Verify secret exists
- Check network connectivity
- Validate credentials

**Schema errors:**

- Run schema migrations manually
- Check admintools pod logs

## Next Steps

After deployment:

1. `/tl-status` - Verify cluster health
2. Set up monitoring - See `monitoring-setup` skill
3. Configure security - See `security-config` skill
4. Deploy workers - See worker configuration guide
