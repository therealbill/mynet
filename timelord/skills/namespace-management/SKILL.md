---
name: namespace-management
description: This skill should be used when the user asks about "temporal namespace", "create namespace", "namespace retention", "multi-tenant temporal", "namespace configuration", "namespace isolation", or needs guidance on organizing and managing Temporal namespaces.
version: 1.0.0
---

# Namespace Management

Guidance for creating, configuring, and managing Temporal namespaces.

## Namespace Concepts

Namespaces provide logical isolation for workflows:

| Aspect | Isolation Level |
|--------|-----------------|
| Workflow IDs | Per namespace (can reuse IDs) |
| Task queues | Per namespace |
| Search attributes | Per namespace |
| Retention | Configurable per namespace |
| Security | Independent authorization |
| Nexus Endpoints | Cross-namespace routing (endpoint → target NS + TQ) |

## Creating Namespaces

### Using CLI

```bash
# Basic creation
temporal operator namespace create --namespace orders

# With retention period
temporal operator namespace create \
  --namespace orders \
  --retention 168h

# With description
temporal operator namespace create \
  --namespace orders \
  --retention 168h \
  --description "Order processing workflows"
```

### Using timelord-cli

```bash
timelord namespace create orders --retention 168h --description "Order processing"
```

### Using SDK

```go
import (
    "go.temporal.io/api/operatorservice/v1"
    "go.temporal.io/sdk/client"
)

func createNamespace(c client.Client) error {
    ctx := context.Background()

    _, err := c.OperatorService().CreateNamespace(ctx, &operatorservice.CreateNamespaceRequest{
        Namespace: "orders",
        WorkflowExecutionRetentionPeriod: &durationpb.Duration{
            Seconds: 604800, // 7 days
        },
        Description: "Order processing workflows",
    })
    return err
}
```

## Namespace Configuration

### Retention Period

How long completed workflow history is retained:

| Environment | Recommended Retention |
|-------------|----------------------|
| Development | 1-3 days |
| Staging | 3-7 days |
| Production | 7-30 days |
| Compliance | Per requirements |

**Update retention:**

```bash
temporal operator namespace update \
  --namespace orders \
  --retention 336h  # 14 days
```

### Search Attributes

Custom searchable fields per namespace:

```bash
# Add search attribute
temporal operator search-attribute create \
  --namespace orders \
  --name CustomerId \
  --type Keyword

temporal operator search-attribute create \
  --namespace orders \
  --name OrderTotal \
  --type Double

temporal operator search-attribute create \
  --namespace orders \
  --name OrderDate \
  --type Datetime
```

**Available types:**

| Type | Description | Example |
|------|-------------|---------|
| Keyword | Exact match string | CustomerID, Status |
| Text | Full-text search | Description |
| Int | Integer values | Count, Attempts |
| Double | Floating point | Amount, Price |
| Bool | Boolean | IsActive, IsPriority |
| Datetime | Timestamps | CreatedAt, DueDate |
| KeywordList | Multiple keywords | Tags, Categories |

### Archival

Configure workflow archival:

```bash
# Enable archival
temporal operator namespace update \
  --namespace orders \
  --history-archival-state enabled \
  --history-archival-uri "s3://bucket/archival"
```

## Multi-Tenancy Patterns

### Namespace Per Team

```
├── team-a-dev
├── team-a-staging
├── team-a-prod
├── team-b-dev
├── team-b-staging
└── team-b-prod
```

**Benefits:**
- Clear ownership
- Independent configuration
- Separate quotas

### Namespace Per Environment

```
├── development
├── staging
├── production
└── production-canary
```

**Benefits:**
- Consistent naming
- Environment isolation
- Simplified promotion

### Namespace Per Service

```
├── orders
├── payments
├── inventory
├── shipping
└── notifications
```

**Benefits:**
- Service isolation
- Independent scaling
- Clear boundaries

### Connecting Services with Nexus

When using namespace-per-service, Nexus endpoints connect them:

```
orders ──nexus: payments-ep──> payments
orders ──nexus: inventory-ep──> inventory
shipping ──nexus: inventory-ep──> inventory
```

Each endpoint routes to the handler namespace's task queue. Teams own their handler services independently.

## Namespace Naming

### Conventions

**Pattern:** `{team}-{service}-{environment}`

```
orders-processing-prod
payments-api-staging
inventory-sync-dev
```

**Rules:**
- Use lowercase
- Use hyphens as separators
- Keep names descriptive but concise
- Include environment suffix

### Examples

| Service | Dev | Staging | Prod |
|---------|-----|---------|------|
| Orders | orders-dev | orders-staging | orders-prod |
| Payments | payments-dev | payments-staging | payments-prod |
| Shared | shared-dev | shared-staging | shared-prod |

## Access Control

### Per-Namespace Authorization

Configure authorization per namespace:

```yaml
# Namespace-specific permissions
namespaces:
  orders-prod:
    allowed_principals:
      - team-orders@example.com
      - platform-team@example.com
    permissions:
      - READ
      - WRITE
      - ADMIN
  orders-staging:
    allowed_principals:
      - team-orders@example.com
    permissions:
      - READ
      - WRITE
```

### Certificate-Based Access

Map certificates to namespaces:

```go
type ClaimMapper struct{}

func (c *ClaimMapper) GetClaims(authInfo *authorization.AuthInfo) (*authorization.Claims, error) {
    // Extract namespace from certificate CN
    cn := authInfo.TLSSubject.CommonName
    // CN format: service.namespace.example.com

    namespace := extractNamespace(cn)

    return &authorization.Claims{
        Namespaces: map[string]authorization.NamespaceClaims{
            namespace: {
                Permissions: []authorization.Permission{
                    authorization.PermissionRead,
                    authorization.PermissionWrite,
                },
            },
        },
    }, nil
}
```

## Nexus Endpoint Management

Nexus endpoints route cross-namespace calls to handler task queues.

### CLI Commands

```bash
# Create endpoint
temporal operator nexus endpoint create \
  --name payments-endpoint \
  --target-namespace payments-ns \
  --target-task-queue payments-tq

# List endpoints
temporal operator nexus endpoint list

# Describe endpoint
temporal operator nexus endpoint describe --name payments-endpoint

# Update endpoint
temporal operator nexus endpoint update \
  --name payments-endpoint \
  --target-task-queue new-payments-tq

# Delete endpoint
temporal operator nexus endpoint delete --name payments-endpoint
```

### Endpoint Naming

Follow the namespace naming pattern: `{service}-{target}-ep`

| Caller | Handler | Endpoint Name |
|--------|---------|---------------|
| orders | payments | payments-ep |
| orders | inventory | inventory-ep |
| shipping | notifications | notifications-ep |

### Topology

```
┌──────────────┐   payments-ep    ┌──────────────┐
│  orders-ns   │────────────────>│ payments-ns  │
│              │   inventory-ep   ┌──────────────┐
│              │────────────────>│ inventory-ns │
└──────────────┘                 └──────────────┘
```

## Operations

### List Namespaces

```bash
# CLI
temporal operator namespace list

# timelord
timelord namespace list --json
```

### Describe Namespace

```bash
# CLI
temporal operator namespace describe --namespace orders

# timelord
timelord namespace describe orders --json
```

### Update Namespace

```bash
temporal operator namespace update \
  --namespace orders \
  --retention 336h \
  --description "Updated description"
```

### Delete Namespace

```bash
# Requires force flag - destructive!
temporal operator namespace delete --namespace orders-test
```

## Migration Between Namespaces

### Workflow Migration Strategy

1. **Stop new workflows** in source namespace
2. **Wait for completion** of running workflows
3. **Start new workflows** in target namespace
4. **Verify functionality** in target
5. **Archive/delete** source namespace

### Data Migration

Export and replay approach:

```bash
# Export histories
temporal workflow list --namespace source-ns --query "ExecutionStatus='Completed'" | \
  while read wf; do
    temporal workflow show --namespace source-ns --workflow-id $wf --output json > histories/$wf.json
  done

# Replay in target (for verification)
# Run replay tests against new namespace
```

## Monitoring

### Namespace Metrics

```promql
# Workflows per namespace
sum by (namespace) (temporal_workflow_active_count)

# Task queue depth per namespace
sum by (namespace, task_queue) (temporal_task_queue_depth)

# Failures per namespace
rate(temporal_workflow_failed_total[5m]) by (namespace)
```

### Alerts

```yaml
groups:
  - name: namespace-alerts
    rules:
      - alert: NamespaceHighFailureRate
        expr: rate(temporal_workflow_failed_total[5m]) > 0.1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High failure rate in namespace {{ $labels.namespace }}"
```

## Best Practices

### Configuration

- Set appropriate retention (not too long, not too short)
- Define search attributes upfront
- Document namespace purpose
- Configure archival for compliance

### Organization

- Use consistent naming conventions
- Document namespace ownership
- Plan namespace structure before creation
- Separate environments clearly

### Security

- Apply least-privilege access
- Use namespace-level authorization
- Audit namespace access
- Rotate credentials regularly

### Operations

- Monitor each namespace independently
- Set up alerts per namespace
- Plan capacity per namespace
- Document runbooks per namespace

## Additional Resources

### Reference Files

For detailed namespace patterns, consult:

- **`references/namespace-patterns.md`** - Advanced organization patterns
- **`references/migration-procedures.md`** - Namespace migration guides
