---
title: "Why Determinism Matters"
weight: 2
---

# Why Determinism Matters

Understanding determinism is crucial for writing correct Temporal workflows. This explains what determinism means in the Temporal context and why violations cause failures.

## The Replay Problem

Temporal workflows must produce the same results when replayed. Here's why.

### How Replay Works

When a worker picks up a workflow task:

```
History:
  1. WorkflowStarted
  2. ActivityScheduled (ValidateOrder)
  3. ActivityCompleted (result: valid)
  4. ActivityScheduled (ChargeCard)   ← Last recorded event

Worker execution:
  1. Start workflow code
  2. Hit ExecuteActivity(ValidateOrder) → Check history → Skip, return "valid"
  3. Hit ExecuteActivity(ChargeCard) → Check history → Skip (started but no result)
  4. → NEW EXECUTION STARTS HERE
  5. Continue with workflow logic
```

The worker **replays** the workflow code, comparing each action against recorded history.

### What Happens on Mismatch

```
History:
  Event 2: ActivityScheduled (ValidateOrder)

Replay execution:
  Hit ExecuteActivity(ProcessPayment)  ← MISMATCH!

Result: NonDeterministicError
```

The workflow code scheduled a different activity than what's in history. Temporal cannot know which version is correct, so it fails.

## Sources of Non-Determinism

### Time

**Wrong:**
```go
func MyWorkflow(ctx workflow.Context) error {
    if time.Now().Hour() < 12 {
        // Morning logic
    }
    // ...
}
```

On replay, `time.Now()` returns replay time (now), not original execution time. The branch taken might differ.

**Correct:**
```go
func MyWorkflow(ctx workflow.Context) error {
    if workflow.Now(ctx).Hour() < 12 {
        // Morning logic
    }
    // ...
}
```

`workflow.Now(ctx)` returns the workflow's logical time, consistent across replays.

### Random Values

**Wrong:**
```go
func MyWorkflow(ctx workflow.Context) error {
    id := uuid.New().String()
    // Use id...
}
```

Each replay generates a different UUID, causing activity inputs to differ.

**Correct:**
```go
func MyWorkflow(ctx workflow.Context) error {
    var id string
    workflow.SideEffect(ctx, func(ctx workflow.Context) interface{} {
        return uuid.New().String()
    }).Get(&id)
    // Use id...
}
```

`SideEffect` records the value on first execution, returns the recorded value on replay.

### Map Iteration

**Wrong:**
```go
func MyWorkflow(ctx workflow.Context) error {
    for key := range myMap {
        workflow.ExecuteActivity(ctx, Process, key).Get(ctx, nil)
    }
}
```

Go map iteration order is random. Activities might be scheduled in different order on replay.

**Correct:**
```go
func MyWorkflow(ctx workflow.Context) error {
    keys := make([]string, 0, len(myMap))
    for k := range myMap {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    for _, key := range keys {
        workflow.ExecuteActivity(ctx, Process, key).Get(ctx, nil)
    }
}
```

Sorting keys ensures consistent order across replays.

### Goroutines

**Wrong:**
```go
func MyWorkflow(ctx workflow.Context) error {
    go func() {
        // Background work
    }()
}
```

Go goroutines are not tracked by Temporal. On replay, they won't exist.

**Correct:**
```go
func MyWorkflow(ctx workflow.Context) error {
    workflow.Go(ctx, func(ctx workflow.Context) {
        // Background work (tracked by Temporal)
    })
}
```

`workflow.Go` creates a coroutine that Temporal tracks and replays correctly.

### External State

**Wrong:**
```go
var globalCounter int

func MyWorkflow(ctx workflow.Context) error {
    globalCounter++
    if globalCounter > 5 {
        // Different behavior
    }
}
```

Global state differs between original execution and replay (different process, different time).

**Correct:**
```go
func MyWorkflow(ctx workflow.Context) error {
    var counter int
    // Load from activity or pass as input
    workflow.ExecuteActivity(ctx, GetCounter).Get(ctx, &counter)

    counter++
    if counter > 5 {
        // Consistent behavior
    }
}
```

Use workflow-local state or fetch state through activities.

## The Determinism Contract

Workflow code must satisfy this contract:

> **Given the same history, workflow code must make the same decisions.**

"Decisions" include:

- Which activities to schedule
- What inputs to pass to activities
- Which timers to start
- Which signals to send
- When to complete

### What's Allowed

| Operation | Deterministic? | Notes |
|-----------|---------------|-------|
| `workflow.Now(ctx)` | ✓ | Returns workflow time |
| `workflow.SideEffect` | ✓ | Records and replays value |
| `workflow.ExecuteActivity` | ✓ | Tracked in history |
| `workflow.Sleep` | ✓ | Timer tracked in history |
| `workflow.Go` | ✓ | Coroutine tracked |
| Math operations | ✓ | Pure functions |
| String manipulation | ✓ | Pure functions |

### What's NOT Allowed

| Operation | Deterministic? | Alternative |
|-----------|---------------|-------------|
| `time.Now()` | ✗ | `workflow.Now(ctx)` |
| `rand.Int()` | ✗ | `workflow.SideEffect` |
| `uuid.New()` | ✗ | `workflow.SideEffect` |
| `os.Getenv()` | ✗ | Pass as input |
| Network calls | ✗ | Use activities |
| File I/O | ✗ | Use activities |
| Database queries | ✗ | Use activities |
| `go func()` | ✗ | `workflow.Go(ctx, ...)` |

## Code Changes and Versioning

When you change workflow code, running workflows might fail on replay.

### The Problem

```
Version 1:
  ExecuteActivity(A)
  ExecuteActivity(B)

Version 2:
  ExecuteActivity(A)
  ExecuteActivity(C)  ← Changed!
  ExecuteActivity(B)

Running workflow on V1, replays on V2:
  History: [A completed, B scheduled]
  Replay: Expects C after A
  Result: NonDeterministicError
```

### The Solution: GetVersion

```go
func MyWorkflow(ctx workflow.Context) error {
    workflow.ExecuteActivity(ctx, A).Get(ctx, nil)

    v := workflow.GetVersion(ctx, "add-activity-c", workflow.DefaultVersion, 1)
    if v >= 1 {
        // New workflows run C
        workflow.ExecuteActivity(ctx, C).Get(ctx, nil)
    }
    // Old workflows (DefaultVersion) skip C

    workflow.ExecuteActivity(ctx, B).Get(ctx, nil)
    return nil
}
```

`GetVersion` records a version marker in history:

- First execution: Records the current max version (1)
- Replay of old workflow: Returns DefaultVersion (-1)
- Replay of new workflow: Returns recorded version (1)

### Version Lifecycle

```
Day 1: Deploy with GetVersion
        Old workflows: Take old path (DefaultVersion)
        New workflows: Take new path (version 1)

Day N: All old workflows complete

Day N+1: Remove GetVersion, simplify to new code only
```

## Testing for Determinism

### Replay Testing

The gold standard for catching determinism issues:

```go
func TestReplayDeterminism(t *testing.T) {
    replayer := worker.NewWorkflowReplayer()
    replayer.RegisterWorkflow(MyWorkflow)

    // Replay production history
    err := replayer.ReplayWorkflowHistoryFromJSONFile(nil, "history.json")
    require.NoError(t, err)
}
```

Workflow changes that break determinism will fail this test.

### CI/CD Integration

1. Export histories from production periodically
2. Store representative histories in test fixtures
3. Run replay tests before deployment

```bash
# Export history
temporal workflow show --workflow-id prod-wf-123 --output json > fixtures/order-workflow.json

# Test
go test -run TestReplay
```

## Mental Model

Think of workflow code as a **script** that describes what should happen:

- The script runs multiple times (original + replays)
- Each run must produce the same sequence of commands
- Activities are the actors that do the real work
- The script shouldn't depend on what time it is when it runs

```
Workflow (Script):
  "First, validate the order"
  "Then, charge the card"
  "Finally, ship the items"

NOT:
  "If it's Tuesday, do X"
  "Pick a random item"
  "Check what's in this file"
```

## Summary

Determinism matters because:

1. **Temporal replays workflow code** to reconstruct state
2. **Non-deterministic code produces different results** on replay
3. **Mismatches cause workflow failures**

Ensure determinism by:

1. Using `workflow.*` APIs for time, randomness, concurrency
2. Keeping side effects in activities
3. Using `GetVersion` for code changes
4. Running replay tests before deployment

When in doubt, ask: "Will this code produce the same result if I run it again with the same inputs and history?"
