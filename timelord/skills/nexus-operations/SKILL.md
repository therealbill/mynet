---
name: nexus-operations
description: Temporal Nexus implementation guidance for cross-namespace durable communication. Use when user asks about "nexus", "nexus operation", "nexus service", "nexus endpoint", "cross-namespace communication", "nexus caller", "nexus handler", "NexusClient", "ExecuteOperation", "WorkflowRunOperation", or "multi-namespace". Covers Go SDK (GA), with TypeScript/Python/Java in references. Do NOT use for architecture decisions about whether to adopt Nexus — use nexus-decision-guide instead.
version: 1.0.0
---

# Nexus Operations

Cross-namespace durable communication using Temporal Nexus (GA).

## Overview

Nexus enables workflows in separate namespaces to call each other through typed service contracts with Temporal's full durable execution guarantees.

| Feature | Nexus | Child Workflows | Activities | Signals |
|---------|-------|-----------------|-----------|---------|
| Cross-namespace | Yes | No (same NS only) | No | No |
| Durable execution | Yes | Yes | No | N/A |
| Request-response | Yes | Yes | Yes | No (one-way) |
| Arbitrary duration | Yes | Yes | No (< 10s) | N/A |
| Loose coupling | Yes | No | Yes | N/A |
| Service contract | Typed interface | Parent-child | Implicit | N/A |

**Architecture:**

```
┌──────────────────┐   Nexus Endpoint   ┌──────────────────┐
│  Caller Namespace │──────────────────>│ Handler Namespace │
│  ┌──────────────┐│   (routes to       │┌────────────────┐│
│  │Caller Workflow││   target NS +     ││Handler Worker   ││
│  │  NexusClient  ││   task queue)     ││ NexusService    ││
│  └──────────────┘│                    │└────────────────┘│
└──────────────────┘                    └──────────────────┘
```

## When to Use Nexus

**Use Nexus when:**

- Teams maintain separate namespaces and need to call each other's workflows
- You need durable execution guarantees across namespace boundaries
- You want clean, typed service contracts between teams/services
- You need both synchronous (< 10s) and asynchronous (long-running) cross-namespace calls

**Don't use Nexus when:**

- All workflows live in one namespace — use child workflows instead
- You need simple fire-and-forget notifications — use signals
- You need quick external API calls from workflows — use activities
- You need synchronous state mutation on a running workflow — use updates
- The added complexity (endpoints, handler workers, service contracts) isn't justified by the isolation benefits

## Choosing the Right Pattern

| Requirement | Best Option |
|-------------|-------------|
| Same namespace, simple decomposition | Child Workflows |
| External side effects, no durability needed | Activities |
| One-way notification to running workflow | Signals |
| Synchronous state mutation | Updates |
| Cross-namespace, durable, request-response | **Nexus** |
| Cross-namespace, short request (< 10s) | **Nexus sync operation** |
| Cross-namespace, long-running process | **Nexus async operation** |

## SDK Support

| SDK | Status | Notes |
|-----|--------|-------|
| Go | GA | Full production support |
| Java | GA | Full production support (1.28.0+) |
| Python | Public Preview | Recommended for production, API may improve |
| TypeScript | Experimental | APIs may have breaking changes |
| .NET | Experimental | APIs may have breaking changes (1.9.0+) |

## Key Concepts

**Nexus Endpoint:** A cluster-level routing resource that maps a name to a target namespace + task queue. Decouples callers from handler location. Created via CLI or Terraform.

**Nexus Service:** A named collection of Nexus Operations. Defines the microservice contract. Registered on handler workers.

**Nexus Operation:** An arbitrary-duration operation — either:
- **Synchronous**: RPC-style, completes immediately (< 10s), handler returns result directly
- **Asynchronous**: Starts a workflow as the backing operation, returns an operation token for tracking

**Nexus Machinery:** Temporal's built-in infrastructure providing at-least-once execution, automatic retries, circuit breaking, and rate limiting for Nexus calls.

## Endpoint Setup

Prerequisites: caller and handler namespaces must exist.

```bash
# Create namespaces (if not already present)
temporal operator namespace create --namespace handler-ns
temporal operator namespace create --namespace caller-ns

# Create Nexus endpoint
temporal operator nexus endpoint create \
  --name my-endpoint \
  --target-namespace handler-ns \
  --target-task-queue handler-tq

# List endpoints
temporal operator nexus endpoint list

# Describe endpoint
temporal operator nexus endpoint describe --name my-endpoint

# Update endpoint
temporal operator nexus endpoint update \
  --name my-endpoint \
  --target-task-queue new-handler-tq

# Delete endpoint
temporal operator nexus endpoint delete --name my-endpoint
```

## Handler Implementation (Go SDK)

### Synchronous Operation

For requests completing in < 10 seconds:

```go
import (
    "context"
    "github.com/nexus-rpc/sdk-go/nexus"
)

// Synchronous operation — returns result directly
var EchoOp = nexus.NewSyncOperation("echo",
    func(ctx context.Context, input EchoInput, opts nexus.StartOperationOptions) (EchoOutput, error) {
        return EchoOutput{Message: input.Message}, nil
    })
```

### Asynchronous (Workflow-Backed) Operation

For long-running processes:

```go
import (
    "context"
    "go.temporal.io/sdk/client"
    temporalnexus "go.temporal.io/sdk/temporalnexus"
)

// Async operation — starts a workflow, returns when workflow completes
var ProcessOrderOp = temporalnexus.NewWorkflowRunOperation("process-order",
    ProcessOrderWorkflow,
    func(ctx context.Context, input OrderInput, opts nexus.StartOperationOptions) (client.StartWorkflowOptions, error) {
        return client.StartWorkflowOptions{
            // Use RequestID as workflow ID to deduplicate retries
            ID: opts.RequestID,
        }, nil
    })
```

### Service Definition and Worker Registration

```go
import (
    "github.com/nexus-rpc/sdk-go/nexus"
    "go.temporal.io/sdk/worker"
)

// Create and register service
func registerNexusService(w worker.Worker) {
    service := nexus.NewService("order-service")
    service.Register(EchoOp)
    service.Register(ProcessOrderOp)
    w.RegisterNexusService(service)
}
```

## Other SDK Implementations

For TypeScript (experimental), Python (public preview), and Java (GA) handler and caller examples, see `references/nexus-multi-sdk.md`.

## Caller Workflow Patterns (Go SDK)

```go
func CallerWorkflow(ctx workflow.Context, input CallerInput) (*CallerOutput, error) {
    // Create Nexus client — endpoint name + service name
    nexusClient := workflow.NewNexusClient("my-endpoint", "order-service")

    // Execute operation (blocks until result)
    future := nexusClient.ExecuteOperation(ctx, ProcessOrderOp, OrderInput{
        ID: input.OrderID,
    }, workflow.NexusOperationOptions{
        ScheduleToCloseTimeout: 10 * time.Minute,
    })

    var result OrderOutput
    if err := future.Get(ctx, &result); err != nil {
        return nil, fmt.Errorf("nexus operation failed: %w", err)
    }

    return &CallerOutput{OrderResult: result}, nil
}
```

## Error Handling

### Error Types

| Error Type | Where | Description | Retryable? |
|------------|-------|-------------|------------|
| `OperationError` | Handler | Application-level failure | No |
| `HandlerError` | Handler | Framework/infrastructure error | Depends on type |
| `NexusOperationFailure` | Caller | Wraps handler error for caller workflow | No |

### Handling Errors in Callers (Go)

```go
var result OrderOutput
if err := future.Get(ctx, &result); err != nil {
    var nexusErr *temporal.NexusOperationFailure
    if errors.As(err, &nexusErr) {
        // Application-level failure from handler
        logger.Error("Nexus operation failed", "error", nexusErr)
        return nil, err
    }
    // Other errors (timeout, cancellation)
    return nil, err
}
```

### Returning Errors from Handlers (Go)

```go
var ProcessOp = nexus.NewSyncOperation("process",
    func(ctx context.Context, input Input, opts nexus.StartOperationOptions) (Output, error) {
        if input.ID == "" {
            // Return application error — caller gets NexusOperationFailure
            return Output{}, nexus.HandlerErrorf(nexus.HandlerErrorTypeBadRequest, "ID is required")
        }
        return Output{Result: "ok"}, nil
    })
```

## Event History

Nexus-specific events in caller workflow history:

| Event | Meaning |
|-------|---------|
| `NexusOperationScheduled` | Caller scheduled a Nexus operation |
| `NexusOperationStarted` | Async operation accepted by handler (returned operation token) |
| `NexusOperationCompleted` | Operation finished successfully |
| `NexusOperationFailed` | Operation returned an error |
| `NexusOperationCanceled` | Operation was canceled |
| `NexusOperationTimedOut` | `scheduleToCloseTimeout` exceeded |

## Best Practices

- **Use sync operations for < 10s requests** — simpler, lower overhead
- **Use async (workflow-backed) for long-running work** — gets full workflow durability
- **Deduplicate with RequestID** — use `opts.RequestID` as workflow ID in async handlers to safely handle retries
- **Implement idempotent handlers** — Nexus provides at-least-once semantics, so handlers may be called more than once
- **Always set `scheduleToCloseTimeout`** on caller side — without it, operations may hang indefinitely
- **Treat Nexus service contracts as stable APIs** — they're microservice boundaries; version them carefully
- **Use separate namespaces** for caller and handler to get true isolation benefits
- **Don't use Nexus within a single namespace** — child workflows are simpler and equally durable

## Additional Resources

### Reference Files

For detailed patterns, consult:

- **`references/nexus-patterns.md`** - Advanced Nexus usage patterns
- **`references/nexus-multi-sdk.md`** - Multi-SDK implementation examples
