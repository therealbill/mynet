---
title: "Configure Prometheus Monitoring"
description: "Use prometheus-expert to set up metrics collection and alerting"
weight: 2
---

# Configure Prometheus Monitoring

Add Prometheus monitoring to your service with correctly typed metrics, controlled label cardinality, and symptom-based alerting that pages for real problems.

## Prerequisites

- Claude Code with the devops-and-infra plugin installed
- A running service that you want to monitor
- Prometheus deployed (or planned) in your infrastructure
- Knowledge of what "healthy" means for your service (latency targets, acceptable error rates)

## Steps

### 1. Describe Your Service and Health Criteria

Tell prometheus-expert what the service does and what matters for its health. Include:

- Service type (HTTP API, gRPC service, background worker, message consumer)
- Current pain points (latency spikes, unknown error rates, no alerting)
- SLO targets if you have them (e.g., 99.9% of requests under 200ms)
- Who should be paged and when

Example prompts that trigger prometheus-expert:

- "Instrument Go service for Prometheus"
- "Alerts for error rate > 1% or p99 > 500ms"
- "Prometheus keeps OOMing"

### 2. Trigger prometheus-expert for Instrumentation Guidance

prometheus-expert activates when it detects Prometheus, metrics instrumentation, or PromQL context in your request. If it does not activate automatically, be explicit:

```
Set up Prometheus metrics for our Go API. We need request duration histograms,
error rate counters, and alerts that page when the error rate exceeds 1% over
5 minutes or p99 latency exceeds 500ms for 10 minutes.
```

The agent reviews any existing `prometheus.yml`, alerting rules, and application instrumentation before proposing changes.

### 3. Review Metric Types and Naming Conventions

prometheus-expert selects the correct metric type for each measurement:

- **Counters** for values that only increase: total requests, total errors, bytes processed
- **Gauges** for values that go up and down: active connections, queue depth, temperature
- **Histograms** for distributions you need to aggregate: request duration, response size. Histograms are preferred over summaries because they are aggregatable across instances
- **Summaries** only when you need accurate quantiles on a single instance and aggregation across instances is not required

All metric names follow the `<namespace>_<subsystem>_<name>_<unit>` convention:

- `myapp_http_request_duration_seconds` (not `request_time` or `http_latency`)
- `myapp_http_requests_total` (not `num_requests` or `request_count`)
- `myapp_db_connections_active` (not `db_conns` or `active_connections`)

Verify that every metric name ends with the unit (`_seconds`, `_bytes`, `_total`) and that the namespace matches your service name.

### 4. Review Alerting Rules

prometheus-expert designs alerts based on symptoms, not causes:

- **Symptom-based:** "Error rate exceeds 1%" alerts on what users experience. "Pod restarted" alerts on infrastructure noise that may or may not affect users
- **`for` durations:** Every alert includes a `for` clause to prevent single-scrape spikes from paging. Error rate alerts typically use 5 minutes. Latency alerts typically use 10 minutes
- **Severity labels:** Each alert carries `severity` (critical, warning, info), `team`, and `runbook_url` labels. Critical pages on-call. Warning creates a ticket. Info is dashboard-only
- **PromQL correctness:** Rate calculations use `rate()` over a range vector at least 4x the scrape interval. Recording rules pre-compute expensive aggregations

Review the `for` durations and severity levels. Noisy alerts -- those that fire frequently without indicating real problems -- erode trust in the alerting system.

### 5. Configure Alertmanager Routing

prometheus-expert provides Alertmanager configuration that routes alerts based on their labels:

- **Critical alerts** with `severity: critical` route to PagerDuty or Opsgenie for immediate on-call notification
- **Warning alerts** with `severity: warning` route to Slack or email for next-business-day response
- **Grouping:** Alerts are grouped by `alertname` and `service` to prevent notification storms when multiple instances of the same service trigger the same alert
- **Inhibition rules:** Critical alerts for a service suppress warning alerts for the same service to reduce noise during incidents

Verify that the routing tree matches your team's on-call structure and escalation policy.

## Verification

A properly configured monitoring setup meets these criteria:

- [ ] All defined metrics appear in Prometheus (check `/targets` for scrape status)
- [ ] Metric names follow the `<namespace>_<subsystem>_<name>_<unit>` convention
- [ ] No label has more than 100 unique values (check with `count by (__name__)({__name__=~"myapp_.*"})`)
- [ ] Alerts fire when you simulate failure conditions (inject errors, add artificial latency)
- [ ] Alertmanager routes test alerts to the correct channel

## Troubleshooting

**Prometheus OOM (out of memory):**

High label cardinality is the most common cause. Labels like `user_id`, `request_id`, or unparameterized URL paths create a unique time series for every unique value. prometheus-expert collapses high-cardinality path parameters (e.g., `/users/123` becomes `/users/{id}`) and avoids unbounded label values entirely. If Prometheus is already OOMing, check `topk(10, count by (__name__)({__name__=~".+"}))` to find the metrics with the most series.

**Metrics not appearing in Prometheus:**

Verify the scrape target is configured in `prometheus.yml` and that the `/metrics` endpoint is reachable from the Prometheus server. Check the Prometheus `/targets` page for scrape errors. Common causes: firewall rules blocking the metrics port, incorrect service discovery configuration, or the application not exposing the `/metrics` endpoint.

**Alerts fire too often (noisy):**

Increase the `for` duration or adjust the threshold. If a latency alert fires every day during a predictable traffic spike, the threshold may be set too low for your actual traffic pattern. prometheus-expert can help recalibrate thresholds based on historical data from `rate()` and `histogram_quantile()` queries over longer time ranges.

## See Also

- [Getting Started]({{< ref "tutorials/getting-started" >}}) -- tutorial covering CI setup and monitoring together
- [Architecture]({{< ref "explanation/architecture" >}}) -- why prometheus-expert and performance-monitor are separate agents
- [Agent Reference]({{< ref "reference/agents" >}}) -- prometheus-expert specification
