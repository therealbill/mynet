---
title: "Why Nexus"
weight: 3
---

# Why Nexus

Understanding when and why to use Temporal Nexus for cross-namespace communication.

## The Problem Nexus Solves

As organizations scale their Temporal usage, they often need multiple namespaces:

- **Team isolation**: Each team owns their namespace with independent workflows
- **Security boundaries**: Sensitive workflows run in restricted namespaces
- **Independent scaling**: Different services scale independently

But these namespaces need to communicate. Before Nexus, teams had limited options:

### The Activity Workaround

```go
// Anti-pattern: Using an activity to call across namespaces
func CrossNamespaceActivity(ctx context.Context, input Input) (Output, error) {
    client, _ := client.Dial(client.Options{Namespace: "other-ns"})
    run, _ := client.ExecuteWorkflow(ctx, client.StartWorkflowOptions{...}, "OtherWorkflow", input)
    var result Output
    run.Get(ctx, &result)
    return result, nil
}
```

Problems with this approach:

- **No durability**: If the activity worker crashes, the connection to the remote workflow is lost
- **Timeout constraints**: Activities have execution timeouts, but the remote workflow may run for hours
- **No built-in retries**: You must implement your own retry logic for the cross-namespace call
- **Tight coupling**: The caller must know the remote namespace, task queue, and workflow type

## What Nexus Provides

Nexus solves these problems by making cross-namespace communication a first-class Temporal concept:

### Durable Execution Across Boundaries

When a caller workflow invokes a Nexus operation:

1. Temporal records `NexusOperationScheduled` in the caller's history
2. The operation is routed through the Nexus endpoint to the handler
3. For async operations, the handler starts a workflow and returns a token
4. Temporal tracks the operation's lifecycle with dedicated events
5. When the handler completes, the caller is notified via `NexusOperationCompleted`

If anything crashes at any point, Temporal's replay mechanism ensures the operation continues from where it left off.

### Typed Service Contracts

Nexus services define explicit contracts between teams:

```go
// Handler defines the service contract
service := nexus.NewService("payment-service")
service.Register(ChargeOp)    // sync: quick charges
service.Register(RefundOp)    // async: long-running refund process
```

This makes cross-namespace dependencies explicit, versioned, and documented.

### Infrastructure-Level Routing

Nexus endpoints decouple callers from handler locations:

```
Caller knows: endpoint name + service name
Endpoint maps to: target namespace + task queue
```

If the handler team changes their task queue or reorganizes their namespace, only the endpoint configuration changes -- caller code stays the same.

## Nexus vs Alternatives

| Approach | Durability | Duration | Coupling | Complexity |
|----------|-----------|----------|----------|------------|
| Activity with client | None | Limited by activity timeout | High | Low |
| Nexus sync operation | Full | < 10 seconds | Low | Medium |
| Nexus async operation | Full | Unlimited | Low | Medium |
| Child workflow | Full | Unlimited | N/A (same NS only) | Low |

### When Nexus Is Overkill

Nexus adds operational overhead (endpoints, handler workers, service contracts). Don't use it when:

- All workflows live in one namespace -- child workflows are simpler
- You need fire-and-forget notifications -- signals work fine
- You're calling external HTTP APIs -- regular activities suffice
- The cross-namespace need is a one-off -- consider consolidating namespaces instead

### When Nexus Shines

- Multiple teams maintaining independent namespaces
- Platform services consumed by many namespaces (notifications, payments, audit)
- Regulated environments requiring namespace-level isolation
- Long-running cross-team processes needing durability guarantees

## The Architecture

```
┌────────────────┐                          ┌────────────────┐
│  Caller NS     │   Nexus Endpoint         │  Handler NS    │
│  ┌──────────┐  │   (cluster-level         │  ┌──────────┐  │
│  │ Workflow  │──│──  routing config  ─────│──│  Handler  │  │
│  │ NexusClnt │  │   maps name to          │  │  Worker   │  │
│  └──────────┘  │   NS + task queue)       │  │  NexusSvc │  │
└────────────────┘                          │  └──────────┘  │
                                            └────────────────┘
```

Nexus leverages Temporal's existing infrastructure -- task queues, workers, history -- while adding the routing and lifecycle management needed for cross-namespace communication.

## Summary

Nexus exists because namespace boundaries shouldn't mean giving up Temporal's durability guarantees. It provides the same reliability for cross-namespace calls that child workflows provide within a single namespace, while maintaining the loose coupling and team independence that separate namespaces are meant to provide.
