---
title: "Call Across Namespaces"
weight: 6
---

# Call Across Namespaces

Use Nexus operations to invoke workflows in other namespaces from within a caller workflow.

## Problem

Your workflow needs to trigger and wait for a long-running process managed by another team in a different namespace.

## Solution

Use `workflow.NewNexusClient` to create a client bound to a Nexus endpoint and service, then call `ExecuteOperation` to invoke the remote operation.

## Prerequisites

- Nexus endpoint configured (see [Set Up Nexus Endpoints](/how-to/setup-nexus-endpoints/))
- Handler service running in the target namespace
- Understanding of sync vs async Nexus operations

## Steps

### 1. Choose Operation Type

| Need | Type | Handler Pattern |
|------|------|----------------|
| Quick lookup (< 10s) | Sync | `nexus.NewSyncOperation` |
| Long-running process | Async | `temporalnexus.NewWorkflowRunOperation` |

### 2. Call a Sync Operation

```go
func CallerWorkflow(ctx workflow.Context, input Input) (*Output, error) {
    nexusClient := workflow.NewNexusClient("my-endpoint", "my-service")

    future := nexusClient.ExecuteOperation(ctx, "echo", EchoInput{
        Message: "hello",
    }, workflow.NexusOperationOptions{
        ScheduleToCloseTimeout: 30 * time.Second,
    })

    var result EchoOutput
    if err := future.Get(ctx, &result); err != nil {
        return nil, err
    }
    return &Output{Message: result.Message}, nil
}
```

### 3. Call an Async Operation

```go
func CallerWorkflow(ctx workflow.Context, input Input) (*Output, error) {
    nexusClient := workflow.NewNexusClient("my-endpoint", "my-service")

    // This starts a workflow in the handler namespace
    // and blocks until that workflow completes
    future := nexusClient.ExecuteOperation(ctx, "process-order", OrderInput{
        OrderID: input.OrderID,
    }, workflow.NexusOperationOptions{
        ScheduleToCloseTimeout: 1 * time.Hour,
    })

    var result OrderOutput
    if err := future.Get(ctx, &result); err != nil {
        return nil, fmt.Errorf("order processing failed: %w", err)
    }
    return &Output{OrderResult: result}, nil
}
```

### 4. Handle Errors

```go
var result Output
if err := future.Get(ctx, &result); err != nil {
    var nexusErr *temporal.NexusOperationFailure
    if errors.As(err, &nexusErr) {
        // Application-level failure from handler
        logger.Error("Nexus operation failed", "error", nexusErr)
    }
    return nil, err
}
```

### 5. Call Multiple Operations

```go
func OrderWorkflow(ctx workflow.Context, order Order) (*OrderResult, error) {
    paymentsClient := workflow.NewNexusClient("payments-ep", "payment-service")
    inventoryClient := workflow.NewNexusClient("inventory-ep", "inventory-service")

    // Execute in sequence
    payFuture := paymentsClient.ExecuteOperation(ctx, "charge", ChargeInput{
        Amount: order.Total,
    }, workflow.NexusOperationOptions{
        ScheduleToCloseTimeout: 5 * time.Minute,
    })

    var payResult ChargeOutput
    if err := payFuture.Get(ctx, &payResult); err != nil {
        return nil, err
    }

    reserveFuture := inventoryClient.ExecuteOperation(ctx, "reserve", ReserveInput{
        Items: order.Items,
    }, workflow.NexusOperationOptions{
        ScheduleToCloseTimeout: 5 * time.Minute,
    })

    var reserveResult ReserveOutput
    if err := reserveFuture.Get(ctx, &reserveResult); err != nil {
        return nil, err
    }

    return &OrderResult{PaymentID: payResult.ID, ReservationID: reserveResult.ID}, nil
}
```

## Verification

- [ ] Nexus operations complete successfully
- [ ] Error handling works for failure cases
- [ ] Timeouts are set appropriately for each operation

## Related

- [Set Up Nexus Endpoints](/how-to/setup-nexus-endpoints/) for initial configuration
- `nexus-operations` skill for comprehensive patterns
