---
title: "Handle Workflow Failures"
weight: 4
---

# Handle Workflow Failures

Diagnose and recover from workflow failures in production.

## Problem

Workflows are failing and you need to identify the cause and recover.

## Solution

Use systematic diagnosis to identify failure causes, then apply appropriate recovery strategies.

## Prerequisites

- Access to Temporal CLI or UI
- Workflow ID of failed execution
- Access to worker logs

## Steps

### 1. Identify the Failure

**List failed workflows:**

```bash
timelord workflow list --query "ExecutionStatus='Failed'" --json
```

**Get workflow details:**

```bash
timelord workflow describe <workflow-id> --json
```

**Run diagnosis:**

```bash
timelord workflow diagnose <workflow-id> --json
```

### 2. Analyze Event History

**Get full history:**

```bash
timelord workflow history <workflow-id> --json > history.json
```

**Find failure events:**

```bash
# Find failed activities
cat history.json | jq '.events[] | select(.eventType | contains("Failed"))'

# Find the last events
cat history.json | jq '.events[-5:]'
```

### 3. Identify Failure Type

**Activity Failure:**

```json
{
  "eventType": "ActivityTaskFailed",
  "activityTaskFailedEventAttributes": {
    "failure": {
      "message": "connection refused"
    }
  }
}
```

**Causes:**
- External service unavailable
- Timeout exceeded
- Code error/panic
- Resource exhaustion

**Workflow Failure:**

```json
{
  "eventType": "WorkflowExecutionFailed",
  "workflowExecutionFailedEventAttributes": {
    "failure": {
      "message": "workflow logic error"
    }
  }
}
```

**Causes:**
- Unhandled error in workflow
- Non-determinism
- Panic in workflow code

### 4. Apply Recovery Strategy

#### For Activity Failures

**Option A: Let retries handle it**

Check retry configuration:

```go
ao := workflow.ActivityOptions{
    RetryPolicy: &temporal.RetryPolicy{
        InitialInterval:    time.Second,
        BackoffCoefficient: 2.0,
        MaximumInterval:    time.Minute,
        MaximumAttempts:    5,
    },
}
```

**Option B: Reset workflow**

Reset to retry from specific point:

```bash
# Reset to last workflow task
temporal workflow reset \
  --workflow-id <id> \
  --type LastWorkflowTask \
  --reason "Retry after external service recovered"
```

**Option C: Send signal to skip**

If workflow supports it:

```bash
temporal workflow signal \
  --workflow-id <id> \
  --name skip-activity \
  --input '{"activityId": "problematic-activity"}'
```

#### For Workflow Failures

**Option A: Fix code and reset**

1. Deploy fixed code
2. Reset workflow:

```bash
temporal workflow reset \
  --workflow-id <id> \
  --type FirstWorkflowTask \
  --reason "Retry with fixed code"
```

**Option B: Terminate and restart**

```bash
# Terminate failed workflow
temporal workflow terminate \
  --workflow-id <id> \
  --reason "Replacing with new execution"

# Start new workflow
temporal workflow start \
  --workflow-id <new-id> \
  --type <WorkflowType> \
  --task-queue <queue> \
  --input '<input>'
```

#### For Non-Determinism Errors

1. Add versioning to code:

```go
v := workflow.GetVersion(ctx, "fix-issue", workflow.DefaultVersion, 1)
if v == workflow.DefaultVersion {
    // Old code path
} else {
    // Fixed code path
}
```

2. Deploy updated code
3. Reset workflow if needed

### 5. Batch Recovery

For multiple failed workflows:

```bash
# Reset all failed workflows of a type
temporal workflow reset-batch \
  --query "WorkflowType='OrderWorkflow' AND ExecutionStatus='Failed'" \
  --type LastWorkflowTask \
  --reason "Batch recovery after service fix"
```

### 6. Prevent Future Failures

**Add better error handling:**

```go
err := workflow.ExecuteActivity(ctx, MyActivity, input).Get(ctx, &result)
if err != nil {
    var appErr *temporal.ApplicationError
    if errors.As(err, &appErr) {
        // Handle specific application errors
        if appErr.Type() == "ValidationError" {
            return handleValidationError(ctx, appErr)
        }
    }
    return err
}
```

**Add circuit breakers:**

```go
// Track consecutive failures
if consecutiveFailures > 3 {
    workflow.Sleep(ctx, time.Minute*5)
    consecutiveFailures = 0
}
```

**Add monitoring:**

```go
workflow.UpsertSearchAttributes(ctx, map[string]interface{}{
    "FailureCount": failureCount,
    "LastError":    lastError.Error(),
})
```

## Verification

- [ ] Workflow recovered successfully
- [ ] Root cause identified
- [ ] Fix deployed to prevent recurrence
- [ ] Monitoring in place

## Recovery Decision Tree

```
Workflow Failed
├── Activity Failed?
│   ├── Retries exhausted → Reset or fix activity
│   ├── Non-retryable error → Fix code, reset
│   └── External service down → Wait, reset
├── Workflow Error?
│   ├── Non-determinism → Add versioning, reset
│   ├── Panic → Fix code, reset
│   └── Logic error → Fix code, reset
└── Timeout?
    ├── Workflow timeout → Increase timeout, restart
    └── Activity timeout → Increase timeout, reset
```

## Common Failure Patterns

| Pattern | Symptom | Solution |
|---------|---------|----------|
| Service unavailable | Connection refused | Add retries, circuit breaker |
| Rate limiting | 429 errors | Add backoff, reduce concurrency |
| Database deadlock | Transaction failed | Add retry, review queries |
| Memory exhaustion | OOM kills | Reduce batch size, add pagination |
| Non-determinism | History mismatch | Add GetVersion, review code |

## Troubleshooting

**Reset fails:**
- Check workflow is in failed/terminated state
- Verify event-id exists in history
- Check you have permission

**Workflow restarts but fails again:**
- Root cause not fixed
- Check worker logs for errors
- Verify external dependencies

**Batch reset hangs:**
- Query returns too many workflows
- Add time bounds to query
- Process in smaller batches

## Related

- `troubleshooting` skill for detailed diagnosis
- `versioning-guide` skill for non-determinism fixes
- Workflow diagnose command
