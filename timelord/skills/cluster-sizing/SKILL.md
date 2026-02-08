---
name: cluster-sizing
description: This skill should be used when the user asks about "Temporal sizing", "history shards", "cluster capacity", "Temporal resources", "scale Temporal", "Temporal performance", "how many shards", or needs guidance on capacity planning for Temporal clusters.
version: 1.0.0
---

# Temporal Cluster Sizing

Guidance for sizing Temporal clusters based on workload requirements.

## Key Sizing Factors

| Factor | Impact | Cannot Change |
|--------|--------|---------------|
| History Shards | Workflow parallelism | Yes (set at creation) |
| History Replicas | Throughput, availability | No |
| Matching Replicas | Task dispatch rate | No |
| Frontend Replicas | API request rate | No |
| Database Size | History storage | No |

## History Shards

**Critical**: History shards cannot be changed after cluster creation.

Shards determine maximum workflow parallelism. Each workflow belongs to one shard.

### Sizing Guidelines

| Concurrent Workflows | Recommended Shards |
|---------------------|-------------------|
| < 10,000 | 128 |
| 10,000 - 100,000 | 256 |
| 100,000 - 500,000 | 512 |
| 500,000 - 2,000,000 | 1024 |
| > 2,000,000 | 2048 or 4096 |

### Calculation Formula

```
shards = ceil(max_concurrent_workflows / 1000) * safety_factor

# Round up to nearest power of 2
# safety_factor = 2-4x for growth
```

**Example**: Expecting 50,000 concurrent workflows with 3x growth:

```
base = 50,000 / 1000 = 50
with_growth = 50 * 3 = 150
nearest_power_of_2 = 256 shards
```

### Shard Distribution

Shards distribute across history service replicas:

```
shards_per_replica = total_shards / history_replicas

# Example: 512 shards, 4 replicas = 128 shards/replica
```

More replicas = better distribution = higher throughput.

## Service Sizing

### Frontend Service

Handles API requests, authentication, rate limiting.

| Load Level | Replicas | CPU | Memory |
|------------|----------|-----|--------|
| Low (<100 rps) | 1-2 | 500m | 1Gi |
| Medium (100-1000 rps) | 3 | 1 | 2Gi |
| High (1000-5000 rps) | 5 | 2 | 4Gi |
| Very High (>5000 rps) | 10+ | 4 | 8Gi |

### History Service

Manages workflow state and event history.

| Shards | Replicas | CPU/replica | Memory/replica |
|--------|----------|-------------|----------------|
| 128 | 2 | 1 | 2Gi |
| 256 | 3 | 2 | 4Gi |
| 512 | 4-6 | 2 | 4Gi |
| 1024 | 8-12 | 4 | 8Gi |
| 2048 | 16-24 | 4 | 8Gi |

### Matching Service

Dispatches tasks to workers.

| Task Rate | Replicas | CPU | Memory |
|-----------|----------|-----|--------|
| Low (<1000/s) | 2 | 500m | 1Gi |
| Medium (1000-10000/s) | 3 | 1 | 2Gi |
| High (>10000/s) | 5+ | 2 | 4Gi |

### Worker Service (Internal)

Handles internal system workflows. Scale with cluster size:

| Cluster Size | Replicas | CPU | Memory |
|--------------|----------|-----|--------|
| Small | 1 | 200m | 256Mi |
| Medium | 1 | 500m | 512Mi |
| Large | 2 | 1 | 1Gi |

## Database Sizing

### PostgreSQL Recommendations

| Workflow Volume | CPU | Memory | Storage | IOPS |
|-----------------|-----|--------|---------|------|
| < 100K workflows | 2 | 8GB | 100GB | 3000 |
| 100K-1M workflows | 4 | 16GB | 500GB | 6000 |
| 1M-10M workflows | 8 | 32GB | 1TB | 12000 |
| > 10M workflows | 16+ | 64GB+ | 2TB+ | 20000+ |

### Storage Calculation

```
storage_per_workflow = avg_history_events * event_size
                     = 100 events * 1KB = 100KB

total_storage = workflows * storage_per_workflow * retention_multiplier
              = 1,000,000 * 100KB * 1.5 = 150GB
```

**Retention**: Configure appropriate workflow retention to manage storage.

## Elasticsearch Sizing

For visibility queries (optional but recommended):

| Indexed Workflows | Nodes | CPU/node | Memory/node | Storage/node |
|-------------------|-------|----------|-------------|--------------|
| < 1M | 3 | 1 | 2Gi | 50Gi |
| 1M-10M | 3 | 2 | 4Gi | 200Gi |
| > 10M | 5+ | 4 | 8Gi | 500Gi |

## Configuration Templates

### Small Cluster (Dev/Test)

```yaml
server:
  config:
    numHistoryShards: 128
  replicaCount:
    frontend: 1
    history: 1
    matching: 1
    worker: 1
  resources:
    frontend:
      requests: {cpu: "250m", memory: "512Mi"}
    history:
      requests: {cpu: "500m", memory: "1Gi"}
    matching:
      requests: {cpu: "250m", memory: "512Mi"}
```

### Medium Cluster (Production Start)

```yaml
server:
  config:
    numHistoryShards: 256
  replicaCount:
    frontend: 3
    history: 3
    matching: 3
    worker: 1
  resources:
    frontend:
      requests: {cpu: "500m", memory: "1Gi"}
      limits: {cpu: "2", memory: "4Gi"}
    history:
      requests: {cpu: "1", memory: "2Gi"}
      limits: {cpu: "4", memory: "8Gi"}
    matching:
      requests: {cpu: "500m", memory: "1Gi"}
      limits: {cpu: "2", memory: "4Gi"}
```

### Large Cluster (High Volume)

```yaml
server:
  config:
    numHistoryShards: 1024
  replicaCount:
    frontend: 5
    history: 10
    matching: 5
    worker: 2
  resources:
    frontend:
      requests: {cpu: "2", memory: "4Gi"}
      limits: {cpu: "4", memory: "8Gi"}
    history:
      requests: {cpu: "4", memory: "8Gi"}
      limits: {cpu: "8", memory: "16Gi"}
    matching:
      requests: {cpu: "2", memory: "4Gi"}
      limits: {cpu: "4", memory: "8Gi"}
```

## Scaling Guidelines

### Horizontal Scaling

Scale replicas when:

- CPU utilization > 70% sustained
- Memory utilization > 80%
- Request latency p99 > SLA
- Task backlog growing

### Vertical Scaling

Increase resources when:

- Replica count at practical limit
- Database connection pooling maxed
- GC pressure affecting latency

## Monitoring for Sizing Decisions

Key metrics to watch:

```promql
# History service load
sum(rate(temporal_persistence_requests_total[5m])) by (operation)

# Task latency (indicates matching capacity)
histogram_quantile(0.99, rate(temporal_schedule_to_start_latency_bucket[5m]))

# Workflow throughput
sum(rate(temporal_workflow_completed_total[5m]))

# Shard distribution
temporal_history_shard_count
```

## Common Sizing Mistakes

| Mistake | Impact | Solution |
|---------|--------|----------|
| Too few shards | Cannot scale later | Start with more shards |
| Undersized history | Latency spikes | Increase memory, replicas |
| Single frontend | Single point of failure | Minimum 2 for HA |
| No Elasticsearch | Slow visibility queries | Enable for production |

## Additional Resources

### Reference Files

For detailed sizing calculations, consult:

- **`references/sizing-calculator.md`** - Detailed sizing formulas
- **`references/benchmark-results.md`** - Performance benchmark data
