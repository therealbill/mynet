---
name: worker-tuning
description: This skill should be used when the user asks about "worker concurrency", "task queue", "worker performance", "worker scaling", "MaxConcurrentActivityExecutionSize", "MaxConcurrentWorkflowTaskExecutionSize", "poller configuration", or needs guidance on optimizing Temporal worker performance.
version: 1.0.0
---

# Worker Tuning

Guidance for configuring and optimizing Temporal worker performance.

## Worker Options Overview

Key configuration parameters:

| Option | Default | Purpose |
|--------|---------|---------|
| `MaxConcurrentActivityExecutionSize` | 1000 | Max parallel activities |
| `MaxConcurrentWorkflowTaskExecutionSize` | 1000 | Max parallel workflow tasks |
| `MaxConcurrentLocalActivityExecutionSize` | 1000 | Max parallel local activities |
| `WorkerActivitiesPerSecond` | 100000 | Activity rate limit |
| `MaxConcurrentActivityTaskPollers` | 2 | Activity poll goroutines |
| `MaxConcurrentWorkflowTaskPollers` | 2 | Workflow poll goroutines |

## Basic Configuration

```go
import (
    "go.temporal.io/sdk/worker"
)

func createWorker(c client.Client) worker.Worker {
    options := worker.Options{
        // Concurrency limits
        MaxConcurrentActivityExecutionSize:         100,
        MaxConcurrentWorkflowTaskExecutionSize:     100,
        MaxConcurrentLocalActivityExecutionSize:    100,

        // Rate limiting
        WorkerActivitiesPerSecond: 1000,

        // Pollers
        MaxConcurrentActivityTaskPollers:  5,
        MaxConcurrentWorkflowTaskPollers:  5,
    }

    return worker.New(c, "my-task-queue", options)
}
```

## Tuning Guidelines

### Activity Concurrency

**Factors to consider:**

- Resource usage per activity (CPU, memory, network)
- External service rate limits
- Database connection pools
- Available system resources

**Calculation:**

```
MaxConcurrentActivityExecutionSize = min(
    Available CPU cores × activities per core,
    External service rate limit,
    Database connection pool size,
    Available memory / memory per activity
)
```

**Example configurations:**

| Activity Type | Recommended Concurrency |
|--------------|------------------------|
| CPU-bound | 2 × CPU cores |
| I/O-bound (fast) | 50-200 |
| I/O-bound (slow) | 10-50 |
| External API calls | API rate limit / workers |

### Workflow Task Concurrency

Workflow tasks are typically fast (< 1s). Higher concurrency is usually safe:

```go
MaxConcurrentWorkflowTaskExecutionSize: 1000,  // Default is good
```

Reduce if workflows have expensive queries or local activities.

### Poller Configuration

**Pollers per task queue:**

```go
// For high-throughput queues
MaxConcurrentActivityTaskPollers:  10,
MaxConcurrentWorkflowTaskPollers:  10,

// For low-throughput queues
MaxConcurrentActivityTaskPollers:  2,
MaxConcurrentWorkflowTaskPollers:  2,
```

**Rule of thumb:** More pollers help with latency-sensitive workloads but consume more connections.

## Task Queue Strategies

### Dedicated Task Queues

Separate task queues for different activity types:

```go
// CPU-intensive activities
cpuWorker := worker.New(c, "cpu-intensive", worker.Options{
    MaxConcurrentActivityExecutionSize: 4,  // Limited to CPU cores
})
cpuWorker.RegisterActivity(RenderVideo)
cpuWorker.RegisterActivity(ProcessImage)

// I/O-bound activities
ioWorker := worker.New(c, "io-bound", worker.Options{
    MaxConcurrentActivityExecutionSize: 100,  // Higher concurrency
})
ioWorker.RegisterActivity(FetchData)
ioWorker.RegisterActivity(SendEmail)
```

### Priority Task Queues

```go
// High priority - more workers
highPriorityWorker := worker.New(c, "orders-high-priority", worker.Options{
    MaxConcurrentActivityExecutionSize: 50,
})

// Normal priority
normalWorker := worker.New(c, "orders-normal", worker.Options{
    MaxConcurrentActivityExecutionSize: 20,
})
```

### Sticky Execution

Workflows prefer returning to the same worker (sticky execution). Configure cache:

```go
worker.Options{
    StickyScheduleToStartTimeout: 5 * time.Second,  // How long to wait for sticky worker
}
```

## Resource Management

### Memory Considerations

```go
// Limit concurrent executions to control memory
worker.Options{
    MaxConcurrentActivityExecutionSize:         50,   // Each activity uses memory
    MaxConcurrentWorkflowTaskExecutionSize:     100,  // Workflow state in memory
}
```

**Estimate memory:**
- Workflow state: 1-10 KB per workflow
- Activity execution: Depends on payload
- Pollers: ~1 MB per poller

### CPU Considerations

```go
// For CPU-bound activities
runtime.GOMAXPROCS(4)  // Limit Go runtime to 4 cores

worker.Options{
    MaxConcurrentActivityExecutionSize: 4,  // Match GOMAXPROCS
}
```

### Connection Pooling

Each poller maintains connections. Balance pollers vs. connections:

```go
// Total connections ≈ pollers × 2
worker.Options{
    MaxConcurrentActivityTaskPollers:  5,   // ~10 connections
    MaxConcurrentWorkflowTaskPollers:  5,   // ~10 connections
}
```

## Scaling Patterns

### Horizontal Scaling

Scale workers across nodes:

```yaml
# Kubernetes deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: temporal-worker
spec:
  replicas: 10  # Scale horizontally
  template:
    spec:
      containers:
      - name: worker
        resources:
          requests:
            cpu: "2"
            memory: "2Gi"
          limits:
            cpu: "4"
            memory: "4Gi"
```

### Auto-Scaling

Scale based on task queue backlog:

```yaml
# HPA based on custom metrics
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: worker-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: temporal-worker
  minReplicas: 2
  maxReplicas: 20
  metrics:
  - type: External
    external:
      metric:
        name: temporal_task_queue_depth
        selector:
          matchLabels:
            task_queue: "my-queue"
      target:
        type: Value
        value: "100"
```

### Vertical Scaling

Increase resources and concurrency:

```go
// Larger instance, higher concurrency
worker.Options{
    MaxConcurrentActivityExecutionSize:        200,
    MaxConcurrentWorkflowTaskExecutionSize:    500,
    MaxConcurrentActivityTaskPollers:          20,
    MaxConcurrentWorkflowTaskPollers:          20,
}
```

## Monitoring Worker Performance

### Key Metrics

```promql
# Activity execution rate
rate(temporal_worker_activity_task_execution_total[5m])

# Schedule-to-start latency (waiting for worker)
histogram_quantile(0.99, rate(temporal_activity_schedule_to_start_latency_bucket[5m]))

# Execution latency
histogram_quantile(0.99, rate(temporal_activity_execution_latency_bucket[5m]))

# Task queue backlog
temporal_task_queue_depth{task_queue="my-queue"}
```

### Performance Indicators

| Metric | Healthy | Action if unhealthy |
|--------|---------|---------------------|
| Schedule-to-start p99 | < 100ms | Add workers or pollers |
| Execution latency p99 | < timeout | Optimize activity |
| Task queue depth | Near 0 | Add workers |
| Worker CPU | < 80% | Increase concurrency |

## Troubleshooting

### High Schedule-to-Start Latency

Causes:
- Not enough workers
- Not enough pollers
- Task queue name mismatch

Solutions:
- Scale up worker replicas
- Increase `MaxConcurrentActivityTaskPollers`
- Verify task queue names match

### Activities Timing Out

Causes:
- Concurrency too high for resources
- External service slow
- Activity too large

Solutions:
- Reduce `MaxConcurrentActivityExecutionSize`
- Add activity heartbeats
- Break into smaller activities

### Worker Memory Issues

Causes:
- Too many concurrent executions
- Large payloads in memory
- Memory leaks

Solutions:
- Reduce concurrency limits
- Use smaller payloads
- Profile memory usage

## Configuration Examples

### Low-Latency Worker

```go
worker.Options{
    MaxConcurrentActivityExecutionSize:    200,
    MaxConcurrentWorkflowTaskExecutionSize: 500,
    MaxConcurrentActivityTaskPollers:      20,
    MaxConcurrentWorkflowTaskPollers:      20,
}
```

### Resource-Constrained Worker

```go
worker.Options{
    MaxConcurrentActivityExecutionSize:    10,
    MaxConcurrentWorkflowTaskExecutionSize: 50,
    MaxConcurrentActivityTaskPollers:      2,
    MaxConcurrentWorkflowTaskPollers:      2,
}
```

### Batch Processing Worker

```go
worker.Options{
    MaxConcurrentActivityExecutionSize:    50,
    WorkerActivitiesPerSecond:             100,  // Rate limit
    MaxConcurrentActivityTaskPollers:      5,
}
```

## Additional Resources

### Reference Files

For detailed tuning patterns, consult:

- **`references/scaling-patterns.md`** - Advanced scaling strategies
- **`references/performance-testing.md`** - Load testing approaches
