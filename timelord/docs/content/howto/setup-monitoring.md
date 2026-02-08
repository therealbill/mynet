---
title: "Set Up Monitoring and Alerting"
weight: 3
---

# Set Up Monitoring and Alerting

Configure Prometheus and Grafana for Temporal cluster observability.

## Problem

You need visibility into your Temporal cluster's health, performance, and issues.

## Solution

Deploy Prometheus for metrics collection and Grafana for visualization with pre-built dashboards.

## Prerequisites

- Running Temporal cluster
- Kubernetes cluster with Helm
- kubectl access

## Steps

### 1. Deploy Prometheus

**Add Prometheus Helm repo:**

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
```

**Create values file:**

```yaml
# prometheus-values.yaml
server:
  persistentVolume:
    size: 50Gi

alertmanager:
  enabled: true

serverFiles:
  prometheus.yml:
    scrape_configs:
      - job_name: 'temporal'
        kubernetes_sd_configs:
          - role: pod
            namespaces:
              names:
                - temporal
        relabel_configs:
          - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
            action: keep
            regex: true
          - source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
            action: replace
            target_label: __metrics_path__
            regex: (.+)
          - source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]
            action: replace
            regex: ([^:]+)(?::\d+)?;(\d+)
            replacement: $1:$2
            target_label: __address__
```

**Install Prometheus:**

```bash
helm install prometheus prometheus-community/prometheus \
  --namespace monitoring \
  --create-namespace \
  -f prometheus-values.yaml
```

### 2. Deploy Grafana

**Add Grafana Helm repo:**

```bash
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
```

**Create values file:**

```yaml
# grafana-values.yaml
persistence:
  enabled: true
  size: 10Gi

adminPassword: "your-secure-password"

datasources:
  datasources.yaml:
    apiVersion: 1
    datasources:
      - name: Prometheus
        type: prometheus
        url: http://prometheus-server.monitoring.svc.cluster.local
        access: proxy
        isDefault: true

dashboardProviders:
  dashboardproviders.yaml:
    apiVersion: 1
    providers:
      - name: 'temporal'
        orgId: 1
        folder: 'Temporal'
        type: file
        disableDeletion: false
        editable: true
        options:
          path: /var/lib/grafana/dashboards/temporal
```

**Install Grafana:**

```bash
helm install grafana grafana/grafana \
  --namespace monitoring \
  -f grafana-values.yaml
```

### 3. Configure Temporal Metrics

Ensure Temporal is exposing metrics. Update Helm values:

```yaml
# temporal-values.yaml
server:
  metrics:
    prometheus:
      timerType: histogram

  podAnnotations:
    prometheus.io/scrape: "true"
    prometheus.io/port: "9090"
```

Upgrade Temporal:

```bash
helm upgrade temporal temporal/temporal \
  --namespace temporal \
  -f temporal-values.yaml
```

### 4. Import Dashboards

**Temporal Server Dashboard:**

Create `temporal-server-dashboard.json`:

```json
{
  "title": "Temporal Server",
  "panels": [
    {
      "title": "Workflow Starts",
      "type": "graph",
      "targets": [
        {
          "expr": "rate(temporal_workflow_started_total[5m])",
          "legendFormat": "{{namespace}}"
        }
      ]
    },
    {
      "title": "Schedule-to-Start Latency",
      "type": "graph",
      "targets": [
        {
          "expr": "histogram_quantile(0.99, rate(temporal_schedule_to_start_latency_bucket[5m]))",
          "legendFormat": "p99"
        },
        {
          "expr": "histogram_quantile(0.95, rate(temporal_schedule_to_start_latency_bucket[5m]))",
          "legendFormat": "p95"
        }
      ]
    },
    {
      "title": "Task Queue Depth",
      "type": "graph",
      "targets": [
        {
          "expr": "temporal_task_queue_depth",
          "legendFormat": "{{task_queue}}"
        }
      ]
    }
  ]
}
```

Import via Grafana UI or API.

### 5. Set Up Alerts

**Create Prometheus alerting rules:**

```yaml
# temporal-alerts.yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: temporal-alerts
  namespace: monitoring
spec:
  groups:
    - name: temporal
      rules:
        - alert: TemporalHighScheduleToStartLatency
          expr: histogram_quantile(0.99, rate(temporal_schedule_to_start_latency_bucket[5m])) > 5
          for: 5m
          labels:
            severity: warning
          annotations:
            summary: "High schedule-to-start latency"
            description: "P99 latency is {{ $value }}s"

        - alert: TemporalWorkflowFailures
          expr: rate(temporal_workflow_failed_total[5m]) > 0.1
          for: 5m
          labels:
            severity: warning
          annotations:
            summary: "Elevated workflow failures"
            description: "Failure rate is {{ $value }}/s"

        - alert: TemporalPersistenceLatency
          expr: histogram_quantile(0.99, rate(temporal_persistence_latency_bucket[5m])) > 0.1
          for: 5m
          labels:
            severity: warning
          annotations:
            summary: "High persistence latency"
            description: "P99 latency is {{ $value }}s"
```

Apply rules:

```bash
kubectl apply -f temporal-alerts.yaml
```

### 6. Configure AlertManager

**Set up notification channels:**

```yaml
# alertmanager-config.yaml
global:
  slack_api_url: 'https://hooks.slack.com/services/xxx'

route:
  receiver: 'slack-notifications'
  group_by: ['alertname', 'namespace']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h

receivers:
  - name: 'slack-notifications'
    slack_configs:
      - channel: '#temporal-alerts'
        send_resolved: true
        title: '{{ .Status | toUpper }}: {{ .CommonAnnotations.summary }}'
        text: '{{ .CommonAnnotations.description }}'
```

## Verification

- [ ] Prometheus is scraping Temporal metrics
- [ ] Grafana dashboards show data
- [ ] Test alerts fire correctly
- [ ] Notifications reach configured channels

## Key Metrics to Monitor

| Metric | Description | Alert Threshold |
|--------|-------------|-----------------|
| `temporal_workflow_started_total` | Workflow start rate | Depends on baseline |
| `temporal_workflow_failed_total` | Failure rate | > 1% of starts |
| `temporal_schedule_to_start_latency` | Worker pickup delay | p99 > 5s |
| `temporal_persistence_latency` | Database latency | p99 > 100ms |
| `temporal_task_queue_depth` | Queue backlog | > 1000 |

## Troubleshooting

**No metrics in Prometheus:**
- Check pod annotations
- Verify network policies allow scraping
- Check Prometheus targets page

**Missing dashboards:**
- Verify dashboard JSON is valid
- Check Grafana provisioner logs
- Import manually if needed

**Alerts not firing:**
- Check AlertManager configuration
- Verify PrometheusRule is loaded
- Test notification channel

## Related

- `monitoring-setup` skill for detailed guidance
- Production setup tutorial
