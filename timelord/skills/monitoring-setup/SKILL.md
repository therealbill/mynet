---
name: monitoring-setup
description: This skill should be used when the user asks about "Temporal monitoring", "Prometheus Temporal", "Grafana Temporal", "Temporal metrics", "Temporal alerts", "monitor workflows", "Temporal observability", or needs guidance on setting up monitoring for Temporal clusters.
version: 1.0.0
---

# Temporal Monitoring Setup

Guidance for configuring monitoring, alerting, and observability for Temporal clusters.

## Monitoring Stack

Recommended stack:

- **Prometheus**: Metrics collection
- **Grafana**: Visualization and dashboards
- **Alertmanager**: Alert routing and notifications

## Metrics Overview

Temporal exposes Prometheus metrics on port 9090 for each service.

### Key Metric Categories

| Category | Prefix | Purpose |
|----------|--------|---------|
| Service | `temporal_*` | Service health, latency |
| Persistence | `temporal_persistence_*` | Database operations |
| Workflow | `temporal_workflow_*` | Workflow execution |
| Activity | `temporal_activity_*` | Activity execution |
| Task Queue | `temporal_task_*` | Task dispatch |

## Prometheus Configuration

### Scrape Configuration

```yaml
# prometheus.yml
scrape_configs:
  - job_name: 'temporal-frontend'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_label_app_kubernetes_io_component]
        regex: frontend
        action: keep
      - source_labels: [__meta_kubernetes_pod_container_port_number]
        regex: "9090"
        action: keep

  - job_name: 'temporal-history'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_label_app_kubernetes_io_component]
        regex: history
        action: keep

  - job_name: 'temporal-matching'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_label_app_kubernetes_io_component]
        regex: matching
        action: keep

  - job_name: 'temporal-worker'
    kubernetes_sd_configs:
      - role: pod
    relabel_configs:
      - source_labels: [__meta_kubernetes_pod_label_app_kubernetes_io_component]
        regex: worker
        action: keep
```

### With Helm

Enable Prometheus in Helm values:

```yaml
prometheus:
  enabled: true
  nodeExporter:
    enabled: false

server:
  metrics:
    prometheus:
      timerType: histogram
```

## Essential Metrics

### Service Health

```promql
# Frontend request rate
sum(rate(temporal_frontend_requests_total[5m])) by (operation)

# Frontend errors
sum(rate(temporal_frontend_errors_total[5m])) by (operation, error_type)

# Frontend latency p99
histogram_quantile(0.99,
  sum(rate(temporal_frontend_request_latency_bucket[5m])) by (le, operation)
)
```

### Persistence (Database)

```promql
# Persistence request rate
sum(rate(temporal_persistence_requests_total[5m])) by (operation)

# Persistence latency p99
histogram_quantile(0.99,
  sum(rate(temporal_persistence_latency_bucket[5m])) by (le, operation)
)

# Persistence errors
sum(rate(temporal_persistence_errors_total[5m])) by (operation, error_type)
```

### Workflow Execution

```promql
# Workflow start rate
sum(rate(temporal_workflow_started_total[5m])) by (namespace, workflow_type)

# Workflow completion rate
sum(rate(temporal_workflow_completed_total[5m])) by (namespace, workflow_type)

# Workflow failure rate
sum(rate(temporal_workflow_failed_total[5m])) by (namespace, workflow_type)

# Workflow execution latency
histogram_quantile(0.99,
  sum(rate(temporal_workflow_endtoend_latency_bucket[5m])) by (le, workflow_type)
)
```

### Task Queue Health

```promql
# Schedule-to-start latency (task wait time)
histogram_quantile(0.99,
  sum(rate(temporal_schedule_to_start_latency_bucket[5m])) by (le, task_queue)
)

# Task dispatch rate
sum(rate(temporal_task_dispatch_total[5m])) by (task_queue, task_type)

# Task backlog
temporal_task_queue_depth
```

## Grafana Dashboards

### Import Official Dashboards

Temporal provides official Grafana dashboards:

```bash
# Dashboard IDs for import
# Server Overview: 10270
# SDK Metrics: 10271
```

### Custom Dashboard Panels

**Service Health Panel:**

```json
{
  "title": "Frontend Request Rate",
  "type": "graph",
  "targets": [{
    "expr": "sum(rate(temporal_frontend_requests_total[5m])) by (operation)",
    "legendFormat": "{{operation}}"
  }]
}
```

**Task Queue Latency Panel:**

```json
{
  "title": "Schedule-to-Start Latency (p99)",
  "type": "graph",
  "targets": [{
    "expr": "histogram_quantile(0.99, sum(rate(temporal_schedule_to_start_latency_bucket[5m])) by (le, task_queue))",
    "legendFormat": "{{task_queue}}"
  }]
}
```

## Alerting Rules

### Critical Alerts

```yaml
# prometheus-rules.yaml
groups:
  - name: temporal-critical
    rules:
      - alert: TemporalServiceDown
        expr: up{job=~"temporal-.*"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Temporal service {{ $labels.job }} is down"

      - alert: TemporalPersistenceHighLatency
        expr: |
          histogram_quantile(0.99,
            sum(rate(temporal_persistence_latency_bucket[5m])) by (le, operation)
          ) > 0.5
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Database latency p99 > 500ms for {{ $labels.operation }}"

      - alert: TemporalPersistenceErrors
        expr: |
          sum(rate(temporal_persistence_errors_total[5m])) by (operation) > 0.1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Persistence errors detected for {{ $labels.operation }}"
```

### Warning Alerts

```yaml
  - name: temporal-warning
    rules:
      - alert: TemporalHighScheduleToStartLatency
        expr: |
          histogram_quantile(0.99,
            sum(rate(temporal_schedule_to_start_latency_bucket[5m])) by (le, task_queue)
          ) > 10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Task queue {{ $labels.task_queue }} has high latency"
          description: "Tasks are waiting >10s to start. Consider scaling workers."

      - alert: TemporalHighFrontendLatency
        expr: |
          histogram_quantile(0.99,
            sum(rate(temporal_frontend_request_latency_bucket[5m])) by (le, operation)
          ) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Frontend latency high for {{ $labels.operation }}"

      - alert: TemporalWorkflowFailureRate
        expr: |
          sum(rate(temporal_workflow_failed_total[5m])) /
          sum(rate(temporal_workflow_completed_total[5m])) > 0.05
        for: 10m
        labels:
          severity: warning
        annotations:
          summary: "Workflow failure rate > 5%"

      - alert: TemporalTaskQueueBacklog
        expr: temporal_task_queue_depth > 1000
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Task queue {{ $labels.task_queue }} backlog > 1000"
```

## Worker Metrics

Configure SDK metrics in your workers:

```go
import (
    "go.temporal.io/sdk/client"
    "go.temporal.io/sdk/contrib/opentelemetry"
    "go.temporal.io/sdk/worker"

    "github.com/uber-go/tally/v4"
    "github.com/uber-go/tally/v4/prometheus"
)

func main() {
    // Create Prometheus reporter
    reporter := prometheus.NewReporter(prometheus.Options{})
    scope, closer := tally.NewRootScope(tally.ScopeOptions{
        Tags:           map[string]string{"service": "my-worker"},
        CachedReporter: reporter,
        Separator:      prometheus.DefaultSeparator,
    }, time.Second)
    defer closer.Close()

    // Create client with metrics
    c, _ := client.Dial(client.Options{
        MetricsHandler: sdktally.NewMetricsHandler(scope),
    })
    defer c.Close()

    // Create worker
    w := worker.New(c, "task-queue", worker.Options{})
    // ...
}
```

### Key SDK Metrics

```promql
# Activity execution latency
temporal_activity_execution_latency_bucket

# Workflow task latency
temporal_workflow_task_execution_latency_bucket

# Worker task slots
temporal_worker_task_slots_available
```

## Logging Configuration

Configure structured logging for correlation:

```yaml
# Temporal server config
log:
  stdout: true
  level: info
  outputFile: ""
```

Include workflow/run IDs in worker logs:

```go
func YourActivity(ctx context.Context, input Input) error {
    logger := activity.GetLogger(ctx)
    logger.Info("Processing",
        "workflowID", activity.GetInfo(ctx).WorkflowExecution.ID,
        "runID", activity.GetInfo(ctx).WorkflowExecution.RunID,
        "activityID", activity.GetInfo(ctx).ActivityID,
    )
    // ...
}
```

## Health Checks

### Liveness Probe

```yaml
livenessProbe:
  httpGet:
    path: /health
    port: 7233
  initialDelaySeconds: 30
  periodSeconds: 10
```

### Readiness Probe

```yaml
readinessProbe:
  httpGet:
    path: /health
    port: 7233
  initialDelaySeconds: 10
  periodSeconds: 5
```

## Troubleshooting with Metrics

| Symptom | Metrics to Check | Likely Cause |
|---------|------------------|--------------|
| Slow workflows | schedule_to_start_latency | Need more workers |
| Workflow failures | workflow_failed_total | Check activity errors |
| API errors | frontend_errors_total | Auth, rate limits |
| DB issues | persistence_latency | Database performance |

## Additional Resources

### Reference Files

For complete metric reference, consult:

- **`references/metrics-reference.md`** - Complete metric documentation
- **`references/dashboard-json.md`** - Grafana dashboard definitions
