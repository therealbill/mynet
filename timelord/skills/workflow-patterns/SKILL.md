---
name: workflow-patterns
description: This skill should be used when the user asks about "workflow patterns", "saga pattern", "state machine workflow", "Continue-As-New", "deterministic workflow", "long-running workflow", "workflow design", or needs guidance on structuring Temporal workflows for specific business scenarios.
version: 1.0.0
---

# Temporal Workflow Patterns

Guidance for designing and implementing Temporal workflow patterns in Go and Python.

## Core Principle: Determinism

All Temporal workflows must be deterministic. The same inputs must always produce the same sequence of commands.

**Determinism Rules:**

- Never use `time.Now()` - use `workflow.Now(ctx)`
- Never use `rand` - use `workflow.SideEffect()` for random values
- Never use UUIDs directly - generate via `workflow.SideEffect()`
- No direct I/O, network calls, or database operations
- No goroutines - use `workflow.Go()` for concurrency
- No global mutable state

**Allowed Operations:**

- Activity execution via `workflow.ExecuteActivity()`
- Child workflows via `workflow.ExecuteChildWorkflow()`
- Timers via `workflow.Sleep()` and `workflow.NewTimer()`
- Signals via `workflow.GetSignalChannel()`
- Queries via `workflow.SetQueryHandler()`
- Local activities via `workflow.ExecuteLocalActivity()`
- Nexus operations via `workflow.NewNexusClient()` and `client.ExecuteOperation()`

## Pattern: Saga (Compensating Transactions)

Implement distributed transactions with compensating actions for rollback.

**When to Use:**

- Multi-step operations across services
- Need to undo completed steps on failure
- Eventual consistency is acceptable

**Go Implementation:**

```go
func OrderSagaWorkflow(ctx workflow.Context, order Order) error {
    var compensations []func(context.Context) error

    // Step 1: Reserve inventory
    err := workflow.ExecuteActivity(ctx, ReserveInventory, order).Get(ctx, nil)
    if err != nil {
        return err
    }
    compensations = append(compensations, func(ctx context.Context) error {
        return ReleaseInventory(ctx, order)
    })

    // Step 2: Charge payment
    err = workflow.ExecuteActivity(ctx, ChargePayment, order).Get(ctx, nil)
    if err != nil {
        return compensate(ctx, compensations)
    }
    compensations = append(compensations, func(ctx context.Context) error {
        return RefundPayment(ctx, order)
    })

    // Step 3: Ship order
    err = workflow.ExecuteActivity(ctx, ShipOrder, order).Get(ctx, nil)
    if err != nil {
        return compensate(ctx, compensations)
    }

    return nil
}

func compensate(ctx workflow.Context, compensations []func(context.Context) error) error {
    // Execute compensations in reverse order
    for i := len(compensations) - 1; i >= 0; i-- {
        ao := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
        actCtx := workflow.WithActivityOptions(ctx, ao)
        workflow.ExecuteActivity(actCtx, compensations[i]).Get(ctx, nil)
    }
    return errors.New("saga rolled back")
}
```

## Pattern: State Machine

Explicit state transitions controlled by signals.

**When to Use:**

- Well-defined states and transitions
- External events drive state changes
- Need audit trail of state changes

**Go Implementation:**

```go
type OrderState string

const (
    StatePending   OrderState = "pending"
    StateApproved  OrderState = "approved"
    StateShipped   OrderState = "shipped"
    StateCancelled OrderState = "cancelled"
)

func OrderStateMachineWorkflow(ctx workflow.Context, orderID string) error {
    state := StatePending

    workflow.SetQueryHandler(ctx, "getState", func() (OrderState, error) {
        return state, nil
    })

    approvalCh := workflow.GetSignalChannel(ctx, "approve")
    shipCh := workflow.GetSignalChannel(ctx, "ship")
    cancelCh := workflow.GetSignalChannel(ctx, "cancel")

    for state != StateShipped && state != StateCancelled {
        selector := workflow.NewSelector(ctx)

        selector.AddReceive(approvalCh, func(c workflow.ReceiveChannel, more bool) {
            if state == StatePending {
                state = StateApproved
            }
        })

        selector.AddReceive(shipCh, func(c workflow.ReceiveChannel, more bool) {
            if state == StateApproved {
                state = StateShipped
            }
        })

        selector.AddReceive(cancelCh, func(c workflow.ReceiveChannel, more bool) {
            if state != StateShipped {
                state = StateCancelled
            }
        })

        selector.Select(ctx)
    }

    return nil
}
```

## Pattern: Entity (Long-Lived)

Workflow representing a business entity with long lifecycle.

**When to Use:**

- Entity with lifecycle spanning days/months/years
- Accumulates state over time
- Responds to many events

**Key Consideration:** Use Continue-As-New to prevent history growth.

```go
func CustomerEntityWorkflow(ctx workflow.Context, customer Customer) error {
    // Set up query handlers
    workflow.SetQueryHandler(ctx, "getCustomer", func() (Customer, error) {
        return customer, nil
    })

    eventCh := workflow.GetSignalChannel(ctx, "event")
    eventCount := 0
    maxEvents := 1000 // Trigger Continue-As-New

    for {
        selector := workflow.NewSelector(ctx)

        selector.AddReceive(eventCh, func(c workflow.ReceiveChannel, more bool) {
            var event CustomerEvent
            c.Receive(ctx, &event)
            customer = processEvent(customer, event)
            eventCount++
        })

        selector.Select(ctx)

        // Continue-As-New to manage history size
        if eventCount >= maxEvents {
            return workflow.NewContinueAsNewError(ctx, CustomerEntityWorkflow, customer)
        }
    }
}
```

## Pattern: Continue-As-New

Reset workflow history while preserving state.

**When to Use:**

- Long-running workflows
- Event history approaching limits (50,000 events default)
- Periodic batch processing

**History Size Guidelines:**

| Scenario | Recommended Threshold |
|----------|----------------------|
| Simple workflows | 10,000 events |
| Complex workflows | 5,000 events |
| High-frequency signals | 1,000 events |

```go
func BatchProcessorWorkflow(ctx workflow.Context, cursor string) error {
    for i := 0; i < 100; i++ {
        var result BatchResult
        err := workflow.ExecuteActivity(ctx, ProcessBatch, cursor).Get(ctx, &result)
        if err != nil {
            return err
        }
        cursor = result.NextCursor

        if result.Done {
            return nil
        }
    }

    // Continue with new history, preserving cursor
    return workflow.NewContinueAsNewError(ctx, BatchProcessorWorkflow, cursor)
}
```

## Pattern: Child Workflows

Decompose complex workflows into smaller, manageable units.

**When to Use:**

- Need independent failure handling
- Separate retry policies required
- Logical grouping of operations
- Parallel execution of workflow segments

```go
func ParentWorkflow(ctx workflow.Context, orders []Order) error {
    var futures []workflow.ChildWorkflowFuture

    // Launch child workflows in parallel
    for _, order := range orders {
        cwo := workflow.ChildWorkflowOptions{
            WorkflowID: "order-" + order.ID,
        }
        childCtx := workflow.WithChildOptions(ctx, cwo)
        future := workflow.ExecuteChildWorkflow(childCtx, ProcessOrderWorkflow, order)
        futures = append(futures, future)
    }

    // Wait for all children
    for _, future := range futures {
        if err := future.Get(ctx, nil); err != nil {
            // Handle individual order failure
            workflow.GetLogger(ctx).Error("Order failed", "error", err)
        }
    }

    return nil
}
```

## Pattern: Cron/Scheduled

Recurring workflow executions on a schedule.

**Configuration:**

```go
func main() {
    c, _ := client.Dial(client.Options{})

    options := client.StartWorkflowOptions{
        ID:           "daily-report",
        TaskQueue:    "reports",
        CronSchedule: "0 9 * * *", // 9 AM daily
    }

    c.ExecuteWorkflow(context.Background(), options, DailyReportWorkflow)
}
```

**Cron Workflow Best Practice:**

```go
func DailyReportWorkflow(ctx workflow.Context) error {
    // Each cron run gets workflow.Now() as start time
    reportDate := workflow.Now(ctx)

    // Process report for this execution
    err := workflow.ExecuteActivity(ctx, GenerateReport, reportDate).Get(ctx, nil)

    // Cron workflows should complete, not loop
    return err
}
```

## Pattern: Cross-Namespace Communication (Nexus)

Invoke workflows in other namespaces through typed service contracts with durable execution guarantees.

**When to Use:**

- Workflows in different namespaces need to communicate
- Teams require namespace isolation but need to call each other's services
- You need durability guarantees across namespace boundaries

**When NOT to Use:**

- All workflows in the same namespace — use child workflows instead
- Simple external API calls — use activities
- Fire-and-forget messaging — use signals

**Go Implementation (Caller):**

```go
func OrderWorkflow(ctx workflow.Context, order Order) (*OrderResult, error) {
    nexusClient := workflow.NewNexusClient("payments-endpoint", "payment-service")

    future := nexusClient.ExecuteOperation(ctx, "charge", ChargeInput{
        OrderID: order.ID,
        Amount:  order.Total,
    }, workflow.NexusOperationOptions{
        ScheduleToCloseTimeout: 5 * time.Minute,
    })

    var chargeResult ChargeOutput
    if err := future.Get(ctx, &chargeResult); err != nil {
        return nil, fmt.Errorf("payment failed: %w", err)
    }

    return &OrderResult{PaymentID: chargeResult.TransactionID}, nil
}
```

## Anti-Patterns to Avoid

**Never do these in workflow code:**

```go
// WRONG: Non-deterministic time
deadline := time.Now().Add(time.Hour)

// CORRECT: Use workflow time
deadline := workflow.Now(ctx).Add(time.Hour)

// WRONG: Direct random
id := uuid.New().String()

// CORRECT: SideEffect for randomness
var id string
workflow.SideEffect(ctx, func(ctx workflow.Context) interface{} {
    return uuid.New().String()
}).Get(&id)

// WRONG: Goroutines
go processItem(item)

// CORRECT: Workflow goroutines
workflow.Go(ctx, func(ctx workflow.Context) {
    processItem(ctx, item)
})
```

## Choosing the Right Pattern

| Requirement | Pattern |
|-------------|---------|
| Multi-service transaction with rollback | Saga |
| Discrete states with signal-driven transitions | State Machine |
| Long-lived business entity | Entity + Continue-As-New |
| Processing large datasets | Batch + Continue-As-New |
| Recurring scheduled tasks | Cron |
| Complex workflow decomposition | Child Workflows |
| Cross-namespace durable communication | Nexus Operations |

## Additional Resources

### Reference Files

For detailed examples and advanced patterns, consult:

- **`references/saga-examples.md`** - Complete saga implementations
- **`references/state-machine-examples.md`** - State machine variations

### Examples

Working code examples in `examples/`:

- **`order-saga.go`** - E-commerce saga workflow
- **`approval-state-machine.go`** - Document approval workflow
