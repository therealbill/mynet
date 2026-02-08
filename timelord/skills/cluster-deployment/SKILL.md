---
name: cluster-deployment
description: This skill should be used when the user asks about "deploy Temporal", "Temporal Helm chart", "install Temporal", "Kubernetes Temporal", "EKS Temporal", "local Temporal cluster", "temporal-server deployment", or needs guidance on deploying self-hosted Temporal clusters.
version: 1.0.0
---

# Temporal Cluster Deployment

Guidance for deploying self-hosted Temporal clusters using Helm on Kubernetes.

## Deployment Options

| Environment | Method | Database | Use Case |
|-------------|--------|----------|----------|
| Local Dev | docker-compose | SQLite/PostgreSQL | Development, testing |
| Local K8s | Helm | PostgreSQL | Integration testing |
| Production | Helm | PostgreSQL + Elasticsearch | Production workloads |

## Local Development with Docker Compose

Quick setup for local development:

```bash
git clone https://github.com/temporalio/docker-compose.git
cd docker-compose
docker-compose up -d
```

Access points:

- Temporal Server: `localhost:7233`
- Web UI: `http://localhost:8080`

## Kubernetes Deployment with Helm

### Prerequisites

- Kubernetes cluster (1.24+)
- Helm 3.x
- kubectl configured
- PostgreSQL database (or provision with Helm)
- Optional: Elasticsearch for advanced visibility

### Add Helm Repository

```bash
helm repo add temporal https://go.temporal.io/helm-charts
helm repo update
```

### Development Configuration

Minimal configuration for development/testing:

```yaml
# values-dev.yaml
server:
  replicaCount:
    frontend: 1
    history: 1
    matching: 1
    worker: 1

  config:
    numHistoryShards: 128

cassandra:
  enabled: false

mysql:
  enabled: false

postgresql:
  enabled: true

elasticsearch:
  enabled: false

prometheus:
  enabled: false

grafana:
  enabled: false
```

Deploy:

```bash
helm install temporal temporal/temporal \
  -f values-dev.yaml \
  --namespace temporal \
  --create-namespace
```

### Production Configuration

Full production configuration with external PostgreSQL:

```yaml
# values-production.yaml
server:
  replicaCount:
    frontend: 3
    history: 3
    matching: 3
    worker: 1

  config:
    numHistoryShards: 512

    persistence:
      default:
        driver: sql
        sql:
          driver: postgres
          host: your-postgresql-host
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
          host: your-postgresql-host
          port: 5432
          database: temporal_visibility
          user: temporal
          existingSecret: temporal-db-credentials

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

cassandra:
  enabled: false

mysql:
  enabled: false

postgresql:
  enabled: false  # Using external PostgreSQL

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

web:
  enabled: true
  replicaCount: 2

prometheus:
  enabled: true

grafana:
  enabled: true
```

### Database Secret

Create the database credentials secret:

```bash
kubectl create secret generic temporal-db-credentials \
  --namespace temporal \
  --from-literal=password='your-db-password'
```

### Deploy

```bash
helm install temporal temporal/temporal \
  -f values-production.yaml \
  --namespace temporal \
  --create-namespace \
  --wait
```

## Schema Management

Initialize database schemas before first deployment:

```bash
# Run schema setup job
kubectl exec -it temporal-admintools-0 -n temporal -- \
  temporal-sql-tool \
    --plugin postgres \
    --ep your-postgresql-host \
    -p 5432 \
    -u temporal \
    --pw 'your-password' \
    --db temporal \
    setup-schema -v 0.0

# Run visibility schema
kubectl exec -it temporal-admintools-0 -n temporal -- \
  temporal-sql-tool \
    --plugin postgres \
    --ep your-postgresql-host \
    -p 5432 \
    -u temporal \
    --pw 'your-password' \
    --db temporal_visibility \
    setup-schema -v 0.0
```

## EKS-Specific Configuration

For Amazon EKS deployments:

```yaml
# values-eks.yaml
server:
  config:
    persistence:
      default:
        sql:
          host: temporal.cluster-xxxxx.region.rds.amazonaws.com
          # Use IAM authentication
          connectAttributes:
            aws_region: us-east-1

# Use ALB for ingress
web:
  ingress:
    enabled: true
    annotations:
      kubernetes.io/ingress.class: alb
      alb.ingress.kubernetes.io/scheme: internal
      alb.ingress.kubernetes.io/target-type: ip
    hosts:
      - temporal.internal.example.com

# Service account for IAM roles
serviceAccount:
  create: true
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::ACCOUNT:role/temporal-role
```

## Verification

Check deployment status:

```bash
# Check pods
kubectl get pods -n temporal

# Check services
kubectl get svc -n temporal

# Port-forward to test locally
kubectl port-forward svc/temporal-frontend 7233:7233 -n temporal
kubectl port-forward svc/temporal-web 8080:8080 -n temporal

# Verify cluster health
temporal operator cluster health
```

## Post-Deployment

After deployment:

1. Create default namespace
2. Configure monitoring alerts
3. Set up mTLS for production
4. Configure worker deployments
5. Set up backup procedures

```bash
# Create default namespace
temporal operator namespace create --namespace default --retention 3d
```

## Troubleshooting

**Pods not starting:**

```bash
kubectl describe pod <pod-name> -n temporal
kubectl logs <pod-name> -n temporal
```

**Database connection issues:**

- Verify network connectivity
- Check secret credentials
- Verify database exists and user has permissions

**Schema errors:**

- Run schema migrations manually
- Check admintools logs

## Additional Resources

### Reference Files

For advanced configurations, consult:

- **`references/helm-values-reference.md`** - Complete Helm values documentation
- **`references/eks-deployment.md`** - EKS-specific deployment guide
