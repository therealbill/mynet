---
name: versioning-guide
description: This skill should be used when the user asks about "workflow versioning", "GetVersion", "workflow migration", "update running workflow", "workflow reset", "deploy workflow changes", "backward compatible workflow", or needs guidance on safely updating Temporal workflows.
version: 1.0.0
---

# Temporal Workflow Versioning

Guidance for safely updating workflows while maintaining compatibility with running executions.

## Why Versioning Matters

Temporal replays workflow history to reconstruct state. If workflow code changes, replay may fail with non-determinism errors. Versioning ensures:

- Running workflows complete successfully
- New workflows use updated logic
- Safe rollback if issues occur

## GetVersion API

The primary tool for workflow versioning.

### Basic Usage

```go
func MyWorkflow(ctx workflow.Context, input Input) error {
    // Introduce a version branch
    v := workflow.GetVersion(ctx, "change-id", workflow.DefaultVersion, 1)

    if v == workflow.DefaultVersion {
        // Original logic (for workflows started before change)
        err := oldBehavior(ctx)
    } else {
        // New logic (v == 1, for new workflows)
        err := newBehavior(ctx)
    }

    return err
}
```

### Parameters

| Parameter | Description |
|-----------|-------------|
| `ctx` | Workflow context |
| `changeID` | Unique identifier for this change point |
| `minSupported` | Minimum version to support (usually DefaultVersion) |
| `maxSupported` | Current version number |

### Version Numbers

- `workflow.DefaultVersion` (-1): Original code before any versioning
- `1, 2, 3, ...`: Incremental version numbers

## Common Scenarios

### Adding New Activity

```go
func OrderWorkflow(ctx workflow.Context, order Order) error {
    // Original activities
    err := workflow.ExecuteActivity(ctx, ValidateOrder, order).Get(ctx, nil)
    if err != nil {
        return err
    }

    // Version gate for new activity
    v := workflow.GetVersion(ctx, "add-fraud-check", workflow.DefaultVersion, 1)
    if v >= 1 {
        // New: Add fraud check for new workflows
        err = workflow.ExecuteActivity(ctx, CheckFraud, order).Get(ctx, nil)
        if err != nil {
            return err
        }
    }

    // Continue with rest of workflow
    err = workflow.ExecuteActivity(ctx, ProcessPayment, order).Get(ctx, nil)
    return err
}
```

### Removing Activity

```go
func OrderWorkflow(ctx workflow.Context, order Order) error {
    v := workflow.GetVersion(ctx, "remove-legacy-check", workflow.DefaultVersion, 1)

    if v == workflow.DefaultVersion {
        // Old: Keep legacy activity for existing workflows
        workflow.ExecuteActivity(ctx, LegacyCheck, order).Get(ctx, nil)
    }
    // New workflows (v >= 1) skip the legacy check

    // Continue with workflow
    return workflow.ExecuteActivity(ctx, ProcessOrder, order).Get(ctx, nil)
}
```

### Changing Activity Parameters

```go
func OrderWorkflow(ctx workflow.Context, order Order) error {
    v := workflow.GetVersion(ctx, "new-shipping-params", workflow.DefaultVersion, 1)

    var shippingResult ShippingResult
    if v == workflow.DefaultVersion {
        // Old parameter format
        err := workflow.ExecuteActivity(ctx, ShipOrder, order.ID).Get(ctx, &shippingResult)
    } else {
        // New parameter format with more data
        input := ShippingInput{
            OrderID:  order.ID,
            Priority: order.Priority,
            Address:  order.ShippingAddress,
        }
        err := workflow.ExecuteActivity(ctx, ShipOrderV2, input).Get(ctx, &shippingResult)
    }

    return err
}
```

### Multiple Versions

```go
func OrderWorkflow(ctx workflow.Context, order Order) error {
    v := workflow.GetVersion(ctx, "payment-flow", workflow.DefaultVersion, 2)

    switch v {
    case workflow.DefaultVersion:
        // Original: single payment call
        return processPaymentV1(ctx, order)
    case 1:
        // V1: Added retry logic
        return processPaymentV2(ctx, order)
    default:
        // V2 (current): Added fraud check
        return processPaymentV3(ctx, order)
    }
}
```

## Deployment Strategy

### Safe Deployment Process

1. **Add version gates** before deployment
2. **Deploy workers** with new code
3. **Monitor** for non-determinism errors
4. **Let old workflows complete** naturally
5. **Clean up** version gates after all old workflows finish

### Deployment Timeline

```
Day 0: Deploy code with GetVersion
        - New workflows: use new logic
        - Old workflows: use old logic

Day 1-N: Old workflows complete naturally
         Monitor for issues

Day N+1: Once all old workflows done:
         - Remove version gates
         - Simplify to new logic only
```

### Monitoring During Rollout

```promql
# Watch for non-determinism errors
sum(rate(temporal_workflow_task_execution_failed_total{failure_reason="NonDeterminismError"}[5m]))

# Track workflow completions by version
# (requires custom search attributes)
```

## Workflow Reset

Alternative to versioning for fixing stuck workflows.

### When to Use Reset

- Workflow stuck due to bug
- Need to retry from specific point
- Cannot wait for natural completion

### Reset Commands

```bash
# Reset to specific event
temporal workflow reset \
  --workflow-id <id> \
  --event-id 10 \
  --reason "Reset after bug fix"

# Reset to last workflow task
temporal workflow reset \
  --workflow-id <id> \
  --type LastWorkflowTask \
  --reason "Retry last decision"

# Reset to first workflow task
temporal workflow reset \
  --workflow-id <id> \
  --type FirstWorkflowTask \
  --reason "Start over"

# Batch reset multiple workflows
temporal workflow reset-batch \
  --query "WorkflowType='OrderWorkflow' AND ExecutionStatus='Running'" \
  --type LastWorkflowTask \
  --reason "Batch reset for bug fix"
```

### Reset Types

| Type | Description | Use Case |
|------|-------------|----------|
| `FirstWorkflowTask` | Restart from beginning | Complete redo |
| `LastWorkflowTask` | Redo last decision | Recent bug |
| `--event-id N` | Reset to specific event | Targeted retry |

## Best Practices

### Version Gate Naming

Use descriptive, unique change IDs:

```go
// GOOD: Descriptive and unique
workflow.GetVersion(ctx, "add-fraud-check-2024-01", ...)
workflow.GetVersion(ctx, "payment-retry-logic", ...)
workflow.GetVersion(ctx, "shipping-v2-params", ...)

// BAD: Generic or reused
workflow.GetVersion(ctx, "change1", ...)
workflow.GetVersion(ctx, "fix", ...)
```

### Cleanup Strategy

After all old workflows complete:

```go
// Before cleanup (versioned)
func OrderWorkflow(ctx workflow.Context, order Order) error {
    v := workflow.GetVersion(ctx, "add-validation", workflow.DefaultVersion, 1)
    if v >= 1 {
        validate(ctx, order)
    }
    return process(ctx, order)
}

// After cleanup (simplified)
func OrderWorkflow(ctx workflow.Context, order Order) error {
    validate(ctx, order)  // Now always runs
    return process(ctx, order)
}
```

### Testing Versions

Test both old and new paths:

```go
func TestOrderWorkflow_NewVersion(t *testing.T) {
    env := s.NewTestWorkflowEnvironment()
    // Test new logic
    env.ExecuteWorkflow(OrderWorkflow, testOrder)
    // Assert new behavior
}

func TestOrderWorkflow_ReplayOldHistory(t *testing.T) {
    replayer := worker.NewWorkflowReplayer()
    replayer.RegisterWorkflow(OrderWorkflow)
    // Replay with old history
    err := replayer.ReplayWorkflowHistoryFromJSONFile(nil, "old_history.json")
    require.NoError(t, err)
}
```

## Anti-Patterns

### Don't Change Version IDs

```go
// WRONG: Changing ID breaks replay
// v1 code
workflow.GetVersion(ctx, "my-change", ...)
// v2 code - BROKEN
workflow.GetVersion(ctx, "my-change-v2", ...)  // Different ID!

// RIGHT: Keep same ID, increment version
// v1 code
workflow.GetVersion(ctx, "my-change", DefaultVersion, 1)
// v2 code
workflow.GetVersion(ctx, "my-change", DefaultVersion, 2)  // Same ID
```

### Don't Lower maxSupported

```go
// WRONG: Lowering version breaks existing workflows
// v1 code
workflow.GetVersion(ctx, "change", DefaultVersion, 2)
// v2 code - BROKEN
workflow.GetVersion(ctx, "change", DefaultVersion, 1)  // Lowered!

// RIGHT: Only increase versions
workflow.GetVersion(ctx, "change", DefaultVersion, 3)
```

## Additional Resources

### Reference Files

For detailed versioning patterns, consult:

- **`references/version-patterns.md`** - Advanced versioning scenarios
- **`references/migration-guide.md`** - Large-scale migration strategies
