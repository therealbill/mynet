---
name: tl-test
description: Test and validate Temporal workflows for determinism and common issues
arguments:
  - name: path
    description: Path to workflow code or history file
    required: false
  - name: type
    description: Test type (validate, replay)
    required: false
---

# Temporal Workflow Testing

Help test and validate Temporal workflow implementations.

## Test Types

### 1. Code Validation

Scan workflow and activity code for common issues:

```bash
# Validate current directory
timelord test validate

# Validate specific path
timelord test validate ./workflows

# JSON output
timelord test validate --json
```

**Checks performed:**

| Check | Description |
|-------|-------------|
| Non-determinism | time.Now(), rand, uuid.New() in workflows |
| Missing context | Activities without context parameter |
| Error handling | Unhandled activity errors |
| Goroutine usage | go func() instead of workflow.Go() |

### 2. Replay Testing

Test workflow determinism by replaying history:

```bash
# Generate replay test from history
timelord test replay history.json --workflow MyWorkflow

# Export workflow history first
temporal workflow show --workflow-id my-wf --output json > history.json
```

**Replay process:**

1. Export workflow history from Temporal
2. Generate replay test code
3. Register workflow with replayer
4. Run test to verify determinism

## Testing Workflow

When user requests workflow testing:

### Step 1: Identify Test Type

**For code validation:**
```bash
timelord test validate ./path/to/code --json
```

**For replay testing:**
```bash
# First get the history
temporal workflow show --workflow-id <id> --output json > history.json

# Then generate replay test
timelord test replay history.json --workflow WorkflowName
```

### Step 2: Analyze Results

**Validation results:**

```json
{
  "valid": true,
  "warnings": [
    {
      "file": "workflow.go",
      "type": "NonDeterminism",
      "description": "time.Now() is non-deterministic",
      "suggestion": "Use workflow.Now(ctx) instead"
    }
  ]
}
```

**Replay results:**

```json
{
  "success": true,
  "eventCount": 42,
  "suggestions": ["Next steps..."]
}
```

### Step 3: Provide Guidance

Use the testing-strategies skill for comprehensive testing guidance:

```
/skill testing-strategies
```

## Common Validations

### Non-Determinism Checks

| Pattern | Issue | Fix |
|---------|-------|-----|
| `time.Now()` | Changes on replay | `workflow.Now(ctx)` |
| `rand.Int()` | Changes on replay | `workflow.SideEffect` |
| `uuid.New()` | Changes on replay | `workflow.SideEffect` |
| `os.Getenv()` | May change | Pass as input |
| `go func()` | Not tracked | `workflow.Go()` |
| `map` iteration | Order varies | Sort keys first |

### Activity Checks

| Pattern | Issue | Fix |
|---------|-------|-----|
| No context | Can't cancel/timeout | Add `context.Context` |
| No error return | Can't signal failure | Return `error` |
| No heartbeat | Long activity may timeout | Add heartbeating |

### Nexus Caller Checks

| Pattern | Issue | Fix |
|---------|-------|-----|
| Missing `scheduleToCloseTimeout` | Operation may hang indefinitely | Set `ScheduleToCloseTimeout` |
| Hardcoded endpoint name | Environment-specific | Use configuration/environment variable |
| No error handling on Nexus future | Lost failures | Handle `NexusOperationFailure` |

## Writing Good Tests

### Unit Test Template

```go
func TestMyWorkflow(t *testing.T) {
    testSuite := &testsuite.WorkflowTestSuite{}
    env := testSuite.NewTestWorkflowEnvironment()

    // Mock activities
    env.OnActivity(MyActivity, mock.Anything).Return("result", nil)

    // Execute
    env.ExecuteWorkflow(MyWorkflow, input)

    // Assert
    require.True(t, env.IsWorkflowCompleted())
    require.NoError(t, env.GetWorkflowError())

    var result string
    require.NoError(t, env.GetWorkflowResult(&result))
    require.Equal(t, "expected", result)
}
```

### Replay Test Template

```go
func TestReplay(t *testing.T) {
    replayer := worker.NewWorkflowReplayer()
    replayer.RegisterWorkflow(MyWorkflow)

    err := replayer.ReplayWorkflowHistoryFromJSONFile(nil, "history.json")
    require.NoError(t, err)
}
```

## Related Skills

For deeper testing guidance:

- **testing-strategies** - Comprehensive testing patterns
- **workflow-patterns** - Determinism requirements
- **troubleshooting** - Debug test failures

## Example Session

User: "Help me test my OrderWorkflow"

1. Run code validation:
   ```bash
   timelord test validate ./workflows --json
   ```

2. If issues found, explain fixes:
   ```
   Found: time.Now() usage
   Fix: Replace with workflow.Now(ctx)
   ```

3. Generate test code:
   ```go
   func TestOrderWorkflow(t *testing.T) {
       // ... generated test
   }
   ```

4. For replay testing, export history:
   ```bash
   temporal workflow show --workflow-id order-123 --output json > history.json
   timelord test replay history.json --workflow OrderWorkflow
   ```
