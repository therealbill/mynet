---
name: upgrade-guide
description: This skill should be used when the user asks about "upgrade temporal", "temporal version upgrade", "migrate temporal", "update temporal server", "upgrade cluster", "schema migration", "SDK compatibility", or needs help planning safe Temporal upgrades.
version: 1.1.0
---

# Temporal Upgrade Guide

Safe procedures for upgrading Temporal server, SDKs, and schema versions.

## When to Use

- Planning a Temporal server version upgrade
- Checking SDK compatibility for a server upgrade
- Running schema migrations during upgrade
- Rolling back a failed upgrade
- Upgrading from one major version to another

## Upgrade Planning

### Pre-Upgrade Checklist

1. **Review release notes** for the target version
2. **Check SDK compatibility matrix** - ensure your application SDK versions are compatible
3. **Identify schema migrations** - check if database schema changes are required
4. **Test in staging** - always upgrade a non-production environment first
5. **Plan maintenance window** - coordinate with stakeholders
6. **Prepare rollback procedure** - document how to revert if needed

### Version Compatibility

Follow sequential upgrade paths for major versions. Do not skip major versions:

```
1.21.x -> 1.22.x -> 1.23.x -> 1.24.x
```

Check SDK compatibility before upgrading:

- Go SDK and Python SDK versions must match the server version range
- Review the compatibility matrix in the Temporal release notes
- Update application dependencies before deploying workers against the new server

## Upgrade Process

### 1. Backup Database

Always backup both the main and visibility databases before upgrading:

```bash
pg_dump -h <host> -U temporal temporal > temporal_backup_$(date +%Y%m%d).sql
pg_dump -h <host> -U temporal temporal_visibility > temporal_visibility_backup_$(date +%Y%m%d).sql
```

### 2. Scale Down Application Workers

```bash
kubectl scale deployment <worker-deployment> --replicas=0 -n <app-namespace>
```

Wait for workers to drain active tasks before proceeding.

### 3. Upgrade Temporal Server

```bash
helm repo update
helm upgrade temporal temporal/temporal \
  --namespace temporal \
  -f values.yaml \
  --version <target_version>
```

### 4. Run Schema Migrations (If Required)

Check release notes for migration requirements. Use the `temporal-sql-tool` from the admintools pod to apply versioned schema changes.

### 5. Verify Cluster Health

```bash
temporal operator cluster health
temporal operator cluster describe
kubectl get pods -n temporal
```

### 6. Update SDK and Deploy Workers

Update your application's SDK dependency to a compatible version, then deploy:

```bash
kubectl apply -f worker-deployment.yaml
kubectl rollout status deployment/<worker-deployment> -n <app-namespace>
```

### 7. Post-Upgrade Verification

- All Temporal pods running
- Cluster health check passes
- Workers connected and processing tasks
- Test workflow completes successfully
- Monitoring dashboards show normal metrics

## Rollback Procedure

If issues occur after upgrading:

1. **Rollback Helm release**: `helm rollback temporal <previous-revision> -n temporal`
2. **Restore database** (if schema changed): restore from backup
3. **Redeploy previous workers**: `kubectl rollout undo deployment/<worker> -n <app-namespace>`

## Common Upgrade Issues

- **Pod startup failures**: Check init container logs, verify schema compatibility, review resource limits
- **Schema migration errors**: Verify database connectivity, check admintools logs
- **SDK incompatibility**: Update SDK to a version compatible with the new server, check for deprecated API usage
- **Non-determinism after upgrade**: New server versions may change default behavior; review workflow replay compatibility

## Related

- Use `/tl-upgrade` command for interactive upgrade planning
- Use the `temporal-ops` agent for operational guidance during upgrades
- See `versioning-guide` skill for workflow versioning during SDK updates
