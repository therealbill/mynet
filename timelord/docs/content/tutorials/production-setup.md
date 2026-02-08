---
title: "Production Temporal Setup"
weight: 3
---

# Production Temporal Setup

Deploy a production-ready Temporal cluster on Kubernetes with PostgreSQL, monitoring, and security.

## What You'll Learn

- Size your cluster appropriately
- Configure external PostgreSQL
- Deploy with Helm to EKS
- Set up monitoring
- Enable mTLS

## Prerequisites

- Kubernetes cluster (EKS, GKE, or similar)
- kubectl and Helm 3.x configured
- PostgreSQL database (RDS, Cloud SQL, or self-managed)
- Domain name for ingress (optional)

## Step 1: Plan Your Cluster

### Determine History Shards

History shards **cannot be changed** after creation. Size appropriately:

| Expected Concurrent Workflows | Shards |
|------------------------------|--------|
| < 100,000 | 256 |
| 100,000 - 500,000 | 512 |
| > 500,000 | 1024 |

### Calculate Resources

For a moderate production workload:

```yaml
# Frontend: 3 replicas, 2Gi each
# History: 4 replicas, 4Gi each
# Matching: 3 replicas, 2Gi each
```

## Step 2: Prepare PostgreSQL

### Create Databases

```sql
CREATE DATABASE temporal;
CREATE DATABASE temporal_visibility;
CREATE USER temporal WITH ENCRYPTED PASSWORD 'your-secure-password';
GRANT ALL PRIVILEGES ON DATABASE temporal TO temporal;
GRANT ALL PRIVILEGES ON DATABASE temporal_visibility TO temporal;
```

### Create Kubernetes Secret

```bash
kubectl create namespace temporal

kubectl create secret generic temporal-db-credentials \
  --namespace temporal \
  --from-literal=password='your-secure-password'
```

## Step 3: Create Helm Values

Create `values-production.yaml`:

```yaml
server:
  replicaCount:
    frontend: 3
    history: 4
    matching: 3
    worker: 1

  config:
    numHistoryShards: 512

    persistence:
      default:
        driver: sql
        sql:
          driver: postgres
          host: your-rds-endpoint.region.rds.amazonaws.com
          port: 5432
          database: temporal
          user: temporal
          existingSecret: temporal-db-credentials
          maxConns: 20
          maxIdleConns: 20
          maxConnLifetime: "1h"

      visibility:
        driver: sql
        sql:
          driver: postgres
          host: your-rds-endpoint.region.rds.amazonaws.com
          port: 5432
          database: temporal_visibility
          user: temporal
          existingSecret: temporal-db-credentials
          maxConns: 10
          maxIdleConns: 10

  resources:
    frontend:
      requests:
        cpu: "500m"
        memory: "1Gi"
      limits:
        cpu: "2"
        memory: "4Gi"
    history:
      requests:
        cpu: "1"
        memory: "2Gi"
      limits:
        cpu: "4"
        memory: "8Gi"
    matching:
      requests:
        cpu: "500m"
        memory: "1Gi"
      limits:
        cpu: "2"
        memory: "4Gi"
    worker:
      requests:
        cpu: "200m"
        memory: "256Mi"
      limits:
        cpu: "1"
        memory: "1Gi"

# Disable embedded databases
cassandra:
  enabled: false

mysql:
  enabled: false

postgresql:
  enabled: false

# Enable Elasticsearch for advanced visibility
elasticsearch:
  enabled: true
  replicas: 3
  minimumMasterNodes: 2
  resources:
    requests:
      cpu: "1"
      memory: "2Gi"
    limits:
      cpu: "2"
      memory: "4Gi"
  persistence:
    enabled: true
    size: 100Gi

# Web UI
web:
  enabled: true
  replicaCount: 2

# Monitoring
prometheus:
  enabled: true

grafana:
  enabled: true
```

## Step 4: Deploy Temporal

```bash
# Add Helm repo
helm repo add temporal https://go.temporal.io/helm-charts
helm repo update

# Deploy
helm install temporal temporal/temporal \
  --namespace temporal \
  -f values-production.yaml \
  --wait --timeout 15m
```

Monitor deployment:

```bash
kubectl get pods -n temporal -w
```

## Step 5: Initialize Schema

If using external PostgreSQL, run schema setup:

```bash
# Get admintools pod
kubectl exec -it temporal-admintools-0 -n temporal -- /bin/bash

# Inside the pod, run:
temporal-sql-tool \
  --plugin postgres \
  --ep your-rds-endpoint.region.rds.amazonaws.com \
  -p 5432 \
  -u temporal \
  --pw 'your-secure-password' \
  --db temporal \
  setup-schema -v 0.0

temporal-sql-tool \
  --plugin postgres \
  --ep your-rds-endpoint.region.rds.amazonaws.com \
  -p 5432 \
  -u temporal \
  --pw 'your-secure-password' \
  --db temporal_visibility \
  setup-schema -v 0.0

# Exit pod
exit
```

## Step 6: Create Namespaces

```bash
# Port forward for CLI access
kubectl port-forward svc/temporal-frontend 7233:7233 -n temporal &

# Create production namespace
temporal operator namespace create \
  --namespace production \
  --retention 30d \
  --description "Production workflows"

# Create staging namespace
temporal operator namespace create \
  --namespace staging \
  --retention 7d \
  --description "Staging workflows"
```

## Step 7: Configure Monitoring

### Access Grafana

```bash
kubectl port-forward svc/temporal-grafana 3000:80 -n temporal
```

Open http://localhost:3000 (default: admin/admin)

### Import Dashboards

1. Go to Dashboards → Import
2. Import dashboard ID `10270` (Temporal Server)
3. Import dashboard ID `10271` (Temporal SDK)

### Configure Alerts

Add to your Prometheus alerts:

```yaml
groups:
  - name: temporal-production
    rules:
      - alert: TemporalServiceDown
        expr: up{job=~"temporal-.*"} == 0
        for: 1m
        labels:
          severity: critical

      - alert: TemporalHighLatency
        expr: |
          histogram_quantile(0.99,
            rate(temporal_persistence_latency_bucket[5m])
          ) > 0.5
        for: 5m
        labels:
          severity: warning
```

## Step 8: Verify Deployment

```bash
# Check cluster health
temporal operator cluster health

# Check cluster info
temporal operator cluster describe

# List namespaces
temporal operator namespace list
```

## Checkpoint

Production deployment complete:

- [x] PostgreSQL configured with proper credentials
- [x] Helm deployment successful
- [x] All pods running and healthy
- [x] Schema initialized
- [x] Production namespace created
- [x] Monitoring accessible

## Next Steps

1. **Security**: Enable mTLS (see How-To guide)
2. **Ingress**: Configure ALB/nginx ingress
3. **Workers**: Deploy your application workers
4. **Backup**: Set up database backup schedule
5. **Runbooks**: Document operational procedures

## Maintenance Checklist

### Daily

- [ ] Check Grafana dashboards
- [ ] Review alert status
- [ ] Monitor task queue latency

### Weekly

- [ ] Review error rates
- [ ] Check database size growth
- [ ] Validate backup integrity

### Monthly

- [ ] Review capacity utilization
- [ ] Plan scaling if needed
- [ ] Review and apply security patches
