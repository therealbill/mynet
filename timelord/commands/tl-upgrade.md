---
name: tl-upgrade
description: Plan and execute Temporal cluster upgrades safely
arguments:
  - name: target_version
    description: Target Temporal version to upgrade to
    required: false
---

# Upgrade Temporal Cluster

Plan and execute safe upgrades for Temporal clusters.

## Usage

```
/tl-upgrade [target_version]
```

## Upgrade Process Overview

1. Review release notes for breaking changes
2. Check compatibility matrix
3. Backup database
4. Scale down application workers
5. Upgrade Temporal server
6. Run schema migrations (if required)
7. Verify cluster health
8. Upgrade SDK versions
9. Deploy updated workers
10. Monitor for issues

## Pre-Upgrade Checklist

### Review Changes

- [ ] Read release notes for target version
- [ ] Identify breaking changes
- [ ] Check schema migration requirements
- [ ] Review SDK compatibility

### Prepare Environment

- [ ] Backup PostgreSQL database
- [ ] Document current cluster configuration
- [ ] Verify rollback procedure
- [ ] Schedule maintenance window
- [ ] Notify stakeholders

### Application Preparation

- [ ] Identify SDK version requirements
- [ ] Test application with new SDK locally
- [ ] Prepare updated worker deployments

## Upgrade Steps

### Step 1: Backup Database

```bash
# PostgreSQL backup
pg_dump -h <host> -U temporal temporal > temporal_backup_$(date +%Y%m%d).sql
pg_dump -h <host> -U temporal temporal_visibility > temporal_visibility_backup_$(date +%Y%m%d).sql
```

### Step 2: Scale Down Application Workers

```bash
# Scale down your application workers
kubectl scale deployment my-worker --replicas=0 -n my-app

# Wait for workers to drain
kubectl get pods -n my-app -w
```

### Step 3: Upgrade Temporal Server

```bash
# Update Helm repository
helm repo update

# Check available versions
helm search repo temporal/temporal --versions

# Upgrade with Helm
helm upgrade temporal temporal/temporal \
  --namespace temporal \
  -f values.yaml \
  --version <target_version>
```

### Step 4: Run Schema Migrations (If Required)

Check release notes for migration requirements:

```bash
# Connect to admintools pod
kubectl exec -it temporal-admintools-0 -n temporal -- /bin/bash

# Run migrations
temporal-sql-tool \
  --plugin postgres \
  --ep <postgresql-host> \
  -p 5432 \
  -u temporal \
  --pw '<password>' \
  --db temporal \
  update-schema -d /etc/temporal/schema/postgresql/v96/temporal/versioned
```

### Step 5: Verify Cluster Health

```bash
# Check pod status
kubectl get pods -n temporal

# Verify cluster health
temporal operator cluster health

# Check cluster info
temporal operator cluster describe
```

### Step 6: Upgrade SDK and Deploy Workers

Update SDK in go.mod:

```go
require (
    go.temporal.io/sdk v1.29.1  // Update to compatible version
)
```

Deploy updated workers:

```bash
kubectl apply -f worker-deployment.yaml
kubectl rollout status deployment/my-worker -n my-app
```

### Step 7: Verify Application

```bash
# Check worker logs
kubectl logs -f deployment/my-worker -n my-app

# Start a test workflow
temporal workflow start \
  --task-queue my-queue \
  --type TestWorkflow \
  --workflow-id test-upgrade-$(date +%s)
```

## Version Compatibility

### SDK Compatibility Matrix

| Server Version | Go SDK | Python SDK |
|----------------|--------|------------|
| 1.24.x | 1.27.x - 1.29.x | 1.5.x - 1.7.x |
| 1.23.x | 1.25.x - 1.28.x | 1.4.x - 1.6.x |
| 1.22.x | 1.23.x - 1.27.x | 1.3.x - 1.5.x |

### Upgrade Path

For major version upgrades, follow sequential path:

```
1.21.x → 1.22.x → 1.23.x → 1.24.x
```

Do not skip major versions.

## Rollback Procedure

If issues occur after upgrade:

### Rollback Helm Release

```bash
# List release history
helm history temporal -n temporal

# Rollback to previous version
helm rollback temporal <revision> -n temporal
```

### Restore Database (If Schema Changed)

```bash
# Restore from backup
psql -h <host> -U temporal temporal < temporal_backup_YYYYMMDD.sql
```

### Redeploy Previous Worker Version

```bash
kubectl rollout undo deployment/my-worker -n my-app
```

## Common Upgrade Issues

**Pod startup failures:**

- Check init container logs
- Verify schema compatibility
- Review resource limits

**Schema migration errors:**

- Verify database connectivity
- Check admintools logs
- Manual migration may be required

**SDK compatibility:**

- Update SDK to compatible version
- Check for deprecated API usage
- Review migration guide

## Post-Upgrade Verification

- [ ] All Temporal pods running
- [ ] Cluster health check passes
- [ ] Application workers connected
- [ ] Test workflow completes successfully
- [ ] Monitoring dashboards show normal metrics
- [ ] No errors in logs
