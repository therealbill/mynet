---
title: "Scale Workers for High Throughput"
weight: 1
---

# Scale Workers for High Throughput

Configure and scale Temporal workers to handle high workflow and activity throughput.

## Problem

Your workflows are experiencing delays, task queues are backing up, or you need to handle increased load.

## Solution

Scale workers horizontally and tune concurrency settings.

## Prerequisites

- Running Temporal cluster
- Worker application deployed
- Monitoring configured (recommended)

## Steps

### 1. Identify Bottleneck

Check task queue latency:

```bash
# Via CLI
temporal task-queue describe --task-queue your-queue

# Or check Prometheus metric
# histogram_quantile(0.99, rate(temporal_schedule_to_start_latency_bucket[5m]))
```

High schedule-to-start latency indicates insufficient workers.

### 2. Configure Worker Concurrency

Tune worker options in your code:

```go
w := worker.New(c, "your-queue", worker.Options{
    // Workflow task concurrency
    MaxConcurrentWorkflowTaskExecutionSize: 1000,

    // Activity concurrency
    MaxConcurrentActivityExecutionSize: 200,

    // Local activity concurrency
    MaxConcurrentLocalActivityExecutionSize: 200,

    // Task queue pollers (network connections)
    MaxConcurrentWorkflowTaskPollers: 4,
    MaxConcurrentActivityTaskPollers: 4,
})
```

**Guidelines:**

| Setting | Default | Recommended Range |
|---------|---------|-------------------|
| WorkflowTaskExecutionSize | 1000 | 500-2000 |
| ActivityExecutionSize | 1000 | 100-500 |
| TaskPollers | 2 | 2-8 |

### 3. Scale Horizontally

Add more worker replicas in Kubernetes:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: temporal-worker
spec:
  replicas: 5  # Increase replicas
  template:
    spec:
      containers:
      - name: worker
        resources:
          requests:
            cpu: "500m"
            memory: "512Mi"
          limits:
            cpu: "2"
            memory: "2Gi"
```

Apply the change:

```bash
kubectl apply -f worker-deployment.yaml

# Or scale directly
kubectl scale deployment temporal-worker --replicas=10
```

### 4. Configure Horizontal Pod Autoscaler

Enable automatic scaling based on CPU:

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: temporal-worker-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: temporal-worker
  minReplicas: 3
  maxReplicas: 20
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

### 5. Use Multiple Task Queues

Distribute work across specialized queues:

```go
// High-priority queue
highPriorityWorker := worker.New(c, "high-priority", worker.Options{
    MaxConcurrentActivityExecutionSize: 50,
})

// Bulk processing queue
bulkWorker := worker.New(c, "bulk-processing", worker.Options{
    MaxConcurrentActivityExecutionSize: 200,
})

// Start both workers
go highPriorityWorker.Run(worker.InterruptCh())
bulkWorker.Run(worker.InterruptCh())
```

In workflows, route to appropriate queue:

```go
ao := workflow.ActivityOptions{
    TaskQueue: "high-priority",  // or "bulk-processing"
    StartToCloseTimeout: time.Minute,
}
```

### 6. Monitor After Scaling

Watch key metrics:

```promql
# Task queue latency (should decrease)
histogram_quantile(0.99, rate(temporal_schedule_to_start_latency_bucket[5m]))

# Worker task slots available
temporal_worker_task_slots_available

# Task dispatch rate
sum(rate(temporal_task_dispatch_total[5m])) by (task_queue)
```

## Verification

- [ ] Schedule-to-start latency < 1 second
- [ ] Task queue depth stable or decreasing
- [ ] Worker CPU utilization 50-70%
- [ ] No worker crashes or OOM events

## Troubleshooting

**Latency still high after scaling:**

- Check matching service replicas
- Verify pollers are connecting
- Look for activity timeouts

**Workers crashing:**

- Reduce concurrency settings
- Increase memory limits
- Check for activity panics

**Uneven load distribution:**

- Ensure all workers use same task queue
- Check for sticky execution settings
- Verify network connectivity

## Related

- `cluster-sizing` skill for server-side scaling
- `monitoring-setup` skill for metrics
- `worker-tuning` skill for advanced tuning
