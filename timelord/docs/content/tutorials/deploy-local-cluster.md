---
title: "Deploy a Local Temporal Cluster"
weight: 2
---

# Deploy a Local Temporal Cluster

Set up a local Temporal development environment using Docker Compose or Kubernetes.

## What You'll Learn

- Run Temporal locally with Docker Compose
- Deploy Temporal to a local Kubernetes cluster
- Access the Temporal Web UI
- Configure the Temporal CLI

## Prerequisites

- Docker and Docker Compose
- OR: minikube/kind and Helm 3.x
- Temporal CLI (`temporal`)

## Option 1: Docker Compose (Recommended for Development)

The fastest way to get started.

### Step 1: Clone the Repository

```bash
git clone https://github.com/temporalio/docker-compose.git
cd docker-compose
```

### Step 2: Start Temporal

```bash
docker-compose up -d
```

This starts:

- Temporal Server (frontend, history, matching, worker)
- PostgreSQL database
- Elasticsearch (for visibility)
- Temporal Web UI

### Step 3: Verify Services

```bash
docker-compose ps
```

All services should show `Up` status.

### Step 4: Access Temporal

- **Web UI**: http://localhost:8080
- **Server gRPC**: localhost:7233

### Step 5: Configure CLI

```bash
export TEMPORAL_ADDRESS=localhost:7233

# Verify connection
temporal operator cluster health
```

### Step 6: Create a Namespace

```bash
temporal operator namespace create \
  --namespace default \
  --retention 3d
```

### Checkpoint

You should now have:

- [x] Temporal server running
- [x] Web UI accessible at localhost:8080
- [x] CLI connected and working
- [x] Default namespace created

## Option 2: Local Kubernetes (minikube)

For testing Kubernetes deployments locally.

### Step 1: Start minikube

```bash
minikube start --memory 4096 --cpus 2
```

### Step 2: Add Helm Repository

```bash
helm repo add temporal https://go.temporal.io/helm-charts
helm repo update
```

### Step 3: Deploy Temporal

```bash
helm install temporal temporal/temporal \
  --namespace temporal \
  --create-namespace \
  --set server.config.numHistoryShards=128 \
  --set cassandra.enabled=false \
  --set postgresql.enabled=true \
  --set elasticsearch.enabled=false \
  --set prometheus.enabled=false \
  --set grafana.enabled=false \
  --wait
```

This takes a few minutes. Watch the progress:

```bash
kubectl get pods -n temporal -w
```

### Step 4: Port Forward Services

```bash
# Frontend service (for CLI and workers)
kubectl port-forward svc/temporal-frontend 7233:7233 -n temporal &

# Web UI
kubectl port-forward svc/temporal-web 8080:8080 -n temporal &
```

### Step 5: Verify and Create Namespace

```bash
export TEMPORAL_ADDRESS=localhost:7233

temporal operator cluster health
temporal operator namespace create --namespace default --retention 3d
```

### Checkpoint

You should now have:

- [x] minikube cluster running
- [x] Temporal deployed via Helm
- [x] Port forwarding configured
- [x] CLI connected
- [x] Default namespace created

## Testing Your Setup

### Start a Test Workflow

Create a simple test using the CLI:

```bash
# Start a workflow
temporal workflow start \
  --task-queue test-queue \
  --type TestWorkflow \
  --workflow-id test-1

# This will fail (no worker) but confirms the server works
```

### View in Web UI

1. Open http://localhost:8080
2. Navigate to the default namespace
3. See your test workflow (it will be pending)

## Stopping the Environment

### Docker Compose

```bash
cd docker-compose
docker-compose down

# To remove volumes (clean slate)
docker-compose down -v
```

### minikube

```bash
# Stop port forwards
pkill -f "port-forward"

# Delete Temporal
helm uninstall temporal -n temporal

# Stop minikube
minikube stop
```

## Troubleshooting

**Docker Compose won't start:**

- Check Docker is running
- Verify ports 7233, 8080 are available
- Check logs: `docker-compose logs temporal`

**Helm install times out:**

- Increase timeout: `--timeout 10m`
- Check pod status: `kubectl describe pod -n temporal`
- Check events: `kubectl get events -n temporal`

**CLI can't connect:**

- Verify TEMPORAL_ADDRESS is set
- Check port forwarding is active
- Test network: `nc -zv localhost 7233`

## Next Steps

With your local cluster running:

1. [Build your first workflow](/tutorials/first-workflow/)
2. Configure monitoring (see `monitoring-setup` skill)
3. Explore the Web UI

## Clean Up

When finished:

**Docker Compose:**

```bash
docker-compose down -v
```

**minikube:**

```bash
minikube delete
```
