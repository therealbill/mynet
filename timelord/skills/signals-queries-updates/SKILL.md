---
name: signals-queries-updates
description: This skill should be used when the user asks about "workflow signal", "workflow query", "workflow update", "send signal", "query workflow state", "SignalWithStart", "workflow communication", or needs guidance on message passing patterns in Temporal.
version: 1.0.0
---

# Signals, Queries, and Updates

Patterns for communicating with running workflows in Temporal.

## Overview

| Feature | Direction | Mutation | Sync |
|---------|-----------|----------|------|
| Signal | External → Workflow | Yes | No |
| Query | External → Workflow | No | Yes |
| Update | External ↔ Workflow | Yes | Yes |

## Signals

Asynchronous messages sent to workflows. Fire-and-forget semantics.

### Defining Signal Handler

```go
func OrderWorkflow(ctx workflow.Context, order Order) error {
    var cancelRequested bool

    // Register signal handler
    cancelChan := workflow.GetSignalChannel(ctx, "cancel-order")

    // Handle signals in selector
    selector := workflow.NewSelector(ctx)
    selector.AddReceive(cancelChan, func(c workflow.ReceiveChannel, more bool) {
        var signal CancelSignal
        c.Receive(ctx, &signal)
        cancelRequested = true
    })

    // Process with signal checking
    for !cancelRequested {
        selector.Select(ctx)
        // Process order steps...
    }

    if cancelRequested {
        return compensate(ctx, order)
    }
    return nil
}
```

### Sending Signals

**From SDK:**

```go
// Get client
c, _ := client.Dial(client.Options{})

// Send signal
err := c.SignalWorkflow(ctx, workflowID, runID, "cancel-order", CancelSignal{
    Reason: "Customer requested",
})
```

**From CLI:**

```bash
temporal workflow signal \
  --workflow-id order-123 \
  --name cancel-order \
  --input '{"reason": "Customer requested"}'
```

### SignalWithStart

Start workflow or signal if already running:

```go
workflowOptions := client.StartWorkflowOptions{
    ID:        "order-123",
    TaskQueue: "orders",
}

we, err := c.SignalWithStartWorkflow(ctx,
    workflowID,
    "add-item",           // signal name
    ItemSignal{ID: "456"}, // signal data
    workflowOptions,
    OrderWorkflow,
    initialOrder,
)
```

### Signal Patterns

**Buffered signals with processing:**

```go
func ProcessingWorkflow(ctx workflow.Context) error {
    var items []Item
    itemChan := workflow.GetSignalChannel(ctx, "add-item")
    doneChan := workflow.GetSignalChannel(ctx, "done")

    done := false
    for !done {
        selector := workflow.NewSelector(ctx)

        selector.AddReceive(itemChan, func(c workflow.ReceiveChannel, more bool) {
            var item Item
            c.Receive(ctx, &item)
            items = append(items, item)
        })

        selector.AddReceive(doneChan, func(c workflow.ReceiveChannel, more bool) {
            c.Receive(ctx, nil)
            done = true
        })

        selector.Select(ctx)
    }

    return processItems(ctx, items)
}
```

## Queries

Synchronous read-only inspection of workflow state.

### Defining Query Handler

```go
func OrderWorkflow(ctx workflow.Context, order Order) error {
    var status string = "pending"
    var progress int = 0

    // Register query handler
    err := workflow.SetQueryHandler(ctx, "get-status", func() (OrderStatus, error) {
        return OrderStatus{
            Status:   status,
            Progress: progress,
        }, nil
    })
    if err != nil {
        return err
    }

    // Workflow logic updates status
    status = "processing"
    progress = 50
    // ...
    status = "completed"
    progress = 100

    return nil
}
```

### Executing Queries

**From SDK:**

```go
response, err := c.QueryWorkflow(ctx, workflowID, runID, "get-status")
if err != nil {
    return err
}

var status OrderStatus
err = response.Get(&status)
```

**From CLI:**

```bash
temporal workflow query \
  --workflow-id order-123 \
  --type get-status
```

### Query Best Practices

| Do | Don't |
|----|-------|
| Return computed state | Modify workflow state |
| Keep handlers fast | Perform I/O operations |
| Return serializable data | Return channels/contexts |
| Handle missing state | Panic on edge cases |

## Updates

Synchronous request-response mutations (Temporal 1.21+).

### Defining Update Handler

```go
func OrderWorkflow(ctx workflow.Context, order Order) error {
    var items []Item

    // Register update handler with validator
    err := workflow.SetUpdateHandlerWithOptions(ctx, "add-item",
        func(ctx workflow.Context, item Item) (AddItemResult, error) {
            items = append(items, item)
            return AddItemResult{
                ItemCount: len(items),
                Total:     calculateTotal(items),
            }, nil
        },
        workflow.UpdateHandlerOptions{
            Validator: func(ctx workflow.Context, item Item) error {
                if item.Quantity <= 0 {
                    return errors.New("quantity must be positive")
                }
                return nil
            },
        },
    )
    if err != nil {
        return err
    }

    // Wait for completion signal
    workflow.GetSignalChannel(ctx, "complete").Receive(ctx, nil)

    return processOrder(ctx, items)
}
```

### Executing Updates

**From SDK:**

```go
handle, err := c.UpdateWorkflow(ctx, client.UpdateWorkflowOptions{
    WorkflowID:   workflowID,
    UpdateName:   "add-item",
    Args:         []interface{}{item},
    WaitForStage: client.WorkflowUpdateStageCompleted,
})
if err != nil {
    return err
}

var result AddItemResult
err = handle.Get(ctx, &result)
```

**From CLI:**

```bash
temporal workflow update \
  --workflow-id order-123 \
  --name add-item \
  --input '{"id": "item-456", "quantity": 2}'
```

### Update vs Signal

| Use Update When | Use Signal When |
|-----------------|-----------------|
| Need confirmation | Fire-and-forget OK |
| Need return value | No response needed |
| Validation required | Simple notification |
| Sync semantics needed | Async acceptable |

## Common Patterns

### Approval Workflow

```go
func ApprovalWorkflow(ctx workflow.Context, request Request) error {
    var approved bool
    var decision Decision

    approveChan := workflow.GetSignalChannel(ctx, "approve")
    rejectChan := workflow.GetSignalChannel(ctx, "reject")

    // Set query for status
    workflow.SetQueryHandler(ctx, "status", func() string {
        if approved {
            return "approved"
        }
        return "pending"
    })

    // Wait for decision
    selector := workflow.NewSelector(ctx)
    selector.AddReceive(approveChan, func(c workflow.ReceiveChannel, _ bool) {
        c.Receive(ctx, &decision)
        approved = true
    })
    selector.AddReceive(rejectChan, func(c workflow.ReceiveChannel, _ bool) {
        c.Receive(ctx, &decision)
        approved = false
    })
    selector.Select(ctx)

    if approved {
        return executeRequest(ctx, request)
    }
    return workflow.NewApplicationError("Request rejected", "REJECTED")
}
```

### Long-Polling Pattern

```go
func LongRunningWorkflow(ctx workflow.Context) error {
    var results []Result

    // Query handler returns current results
    workflow.SetQueryHandler(ctx, "get-results", func() []Result {
        return results
    })

    // Process items, updating results
    for i := 0; i < 100; i++ {
        result := processItem(ctx, i)
        results = append(results, result)
    }

    return nil
}
```

## Error Handling

### Signal Errors

Signals are fire-and-forget. Handle errors in workflow:

```go
selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
    var signal MySignal
    c.Receive(ctx, &signal)

    if err := validateSignal(signal); err != nil {
        // Log error, don't return - signal is already received
        workflow.GetLogger(ctx).Error("Invalid signal", "error", err)
        return
    }
    // Process valid signal
})
```

### Query Errors

Return errors from query handlers:

```go
workflow.SetQueryHandler(ctx, "get-data", func(key string) (Data, error) {
    data, ok := dataMap[key]
    if !ok {
        return Data{}, fmt.Errorf("key not found: %s", key)
    }
    return data, nil
})
```

### Update Validation

Use validators to reject invalid updates:

```go
workflow.SetUpdateHandlerWithOptions(ctx, "update-quantity",
    func(ctx workflow.Context, qty int) error {
        // Handler runs after validation passes
        quantity = qty
        return nil
    },
    workflow.UpdateHandlerOptions{
        Validator: func(ctx workflow.Context, qty int) error {
            if qty < 0 {
                return errors.New("quantity cannot be negative")
            }
            if qty > maxQuantity {
                return fmt.Errorf("quantity exceeds max: %d", maxQuantity)
            }
            return nil
        },
    },
)
```

## Additional Resources

### Reference Files

For detailed patterns, consult:

- **`references/signal-patterns.md`** - Advanced signal handling
- **`references/update-patterns.md`** - Update implementation examples
