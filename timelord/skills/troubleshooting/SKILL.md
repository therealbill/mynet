---
name: troubleshooting
description: This skill should be used when the user asks about "workflow stuck", "workflow failing", "temporal error", "debug workflow", "diagnose temporal", "workflow not completing", "activity timeout", "non-deterministic error", or needs help resolving Temporal issues.
version: 1.0.0
---

# Temporal Troubleshooting

Guidance for diagnosing and resolving common Temporal workflow issues.

## Diagnostic Approach

1. **Identify symptoms** - What's the observable problem?
2. **Gather information** - Check status, events, logs
3. **Isolate cause** - Narrow down to specific component
4. **Apply fix** - Resolve with targeted solution
5. **Verify resolution** - Confirm problem is fixed

## Quick Diagnostics

### Check Workflow Status

```bash
# Describe workflow
temporal workflow describe --workflow-id <id>

# List running workflows
temporal workflow list --query "ExecutionStatus='Running'"

# List failed workflows
temporal workflow list --query "ExecutionStatus='Failed'"
```

### Check Event History

```bash
# Show all events
temporal workflow show --workflow-id <id>

# Get JSON for parsing
temporal workflow show --workflow-id <id> --output json > history.json

# Find specific events
temporal workflow show --workflow-id <id> --output json | \
  jq '.events[] | select(.eventType | contains("Failed"))'
```

### Check Task Queue

```bash
# Describe task queue
temporal task-queue describe --task-queue <queue>
```

## Common Issues

### Issue: Workflow Stuck (Not Progressing)

**Symptoms:**

- Workflow status is "Running" but nothing happening
- No recent events in history

**Diagnosis Tree:**

```
Workflow stuck
├── Last event is ActivityTaskScheduled?
│   ├── Workers running? → Check/start workers
│   ├── Correct task queue? → Fix queue name
│   └── Activity timed out? → Check timeout config
├── Last event is TimerStarted?
│   └── Timer still pending → Wait or reset
├── Last event is SignalExternalWorkflow?
│   └── Target workflow not responding → Check target
└── Last event is WorkflowTaskScheduled?
    ├── Workers running? → Check/start workers
    └── Workflow error? → Check worker logs
```

**Commands:**

```bash
# Check what workflow is waiting for
temporal workflow show --workflow-id <id> | tail -20

# Check if workers are connected
temporal task-queue describe --task-queue <queue>
```

**Solutions:**

| Waiting For | Solution |
|-------------|----------|
| Activity | Start workers, fix task queue |
| Timer | Wait or reset workflow |
| Signal | Send signal or reset |
| Child workflow | Check child status |
| Nexus operation | Check endpoint and handler workers |

### Issue: Non-Determinism Error

**Symptoms:**

- Error: `non-deterministic workflow definition`
- Error: `history mismatch`
- Workflow fails on replay

**Common Causes:**

1. **Time-based decisions**
   ```go
   // BAD
   if time.Now().Hour() < 12 { }

   // GOOD
   if workflow.Now(ctx).Hour() < 12 { }
   ```

2. **Random values**
   ```go
   // BAD
   id := uuid.New()

   // GOOD
   var id string
   workflow.SideEffect(ctx, func(ctx workflow.Context) interface{} {
       return uuid.New().String()
   }).Get(&id)
   ```

3. **Map iteration order**
   ```go
   // BAD - order varies
   for k, v := range myMap { }

   // GOOD - deterministic order
   keys := make([]string, 0, len(myMap))
   for k := range myMap {
       keys = append(keys, k)
   }
   sort.Strings(keys)
   for _, k := range keys {
       v := myMap[k]
   }
   ```

4. **Code changes without versioning**
   ```go
   // When changing workflow logic, use GetVersion
   v := workflow.GetVersion(ctx, "change-id", workflow.DefaultVersion, 1)
   if v == workflow.DefaultVersion {
       // Old logic
   } else {
       // New logic
   }
   ```

**Diagnosis:**

```bash
# Export history
temporal workflow show --workflow-id <id> --output json > history.json

# Create replay test
```

```go
func TestReplay(t *testing.T) {
    replayer := worker.NewWorkflowReplayer()
    replayer.RegisterWorkflow(YourWorkflow)
    err := replayer.ReplayWorkflowHistoryFromJSONFile(nil, "history.json")
    require.NoError(t, err)
}
```

### Issue: Activity Timeout

**Symptoms:**

- Activity fails with timeout error
- `ActivityTaskTimedOut` events in history

**Timeout Types:**

| Timeout | Meaning | Solution |
|---------|---------|----------|
| ScheduleToStart | Waiting for worker | Add workers, check queue |
| StartToClose | Execution too long | Increase timeout or optimize |
| ScheduleToClose | Total time exceeded | Increase or split activity |
| Heartbeat | No heartbeat received | Add heartbeat, check worker |

**Diagnosis:**

```bash
# Find timeout events
temporal workflow show --workflow-id <id> --output json | \
  jq '.events[] | select(.eventType == "ActivityTaskTimedOut")'
```

**Solutions:**

```go
ao := workflow.ActivityOptions{
    // Increase timeouts if needed
    StartToCloseTimeout: 30 * time.Minute,

    // Add heartbeat for long activities
    HeartbeatTimeout: 30 * time.Second,

    // Configure retries
    RetryPolicy: &temporal.RetryPolicy{
        MaximumAttempts: 5,
    },
}
```

### Issue: Activity Failures

**Symptoms:**

- Activities failing repeatedly
- `ActivityTaskFailed` events

**Diagnosis:**

```bash
# Find failure details
temporal workflow show --workflow-id <id> --output json | \
  jq '.events[] | select(.eventType == "ActivityTaskFailed") | .activityTaskFailedEventAttributes'
```

**Common Causes:**

| Error Type | Cause | Solution |
|------------|-------|----------|
| Connection error | Network issue | Add retries, check connectivity |
| Panic | Code bug | Fix activity code |
| Application error | Business logic | Check error handling |
| Resource exhausted | Rate limiting | Add backoff, reduce load |

### Issue: High Latency

**Symptoms:**

- Workflows slow to progress
- Long time between events

**Check Points:**

```promql
# Schedule-to-start latency (Prometheus)
histogram_quantile(0.99, rate(temporal_schedule_to_start_latency_bucket[5m]))

# Persistence latency
histogram_quantile(0.99, rate(temporal_persistence_latency_bucket[5m]))
```

**Solutions:**

| Latency Location | Solution |
|-----------------|----------|
| Schedule-to-start | Add workers |
| Activity execution | Optimize activity |
| Database | Scale database |
| Network | Check connectivity |

### Issue: Nexus Operation Failures

**Symptoms:**

- `NexusOperationFailed` or `NexusOperationTimedOut` events in caller workflow history
- Caller workflow stuck waiting for Nexus operation

**Diagnosis Tree:**

```
Nexus operation issue
├── NexusOperationTimedOut?
│   ├── scheduleToCloseTimeout too short → Increase timeout
│   └── Handler workflow stuck → Debug handler workflow in handler namespace
├── NexusOperationFailed?
│   ├── OperationError → Check handler operation logic
│   ├── HandlerError → Check handler worker logs/infrastructure
│   └── Endpoint misconfigured → Verify endpoint config
└── NexusOperationScheduled but never started?
    ├── Endpoint exists? → temporal operator nexus endpoint list
    ├── Handler workers running? → Check handler task queue
    └── Target namespace accessible? → Verify namespace exists
```

**Diagnostic Commands:**

```bash
# Check Nexus events in caller workflow
temporal workflow show --workflow-id <caller-wf-id> --output json | \
  jq '.events[] | select(.eventType | contains("Nexus"))'

# List Nexus endpoints
temporal operator nexus endpoint list

# Describe specific endpoint
temporal operator nexus endpoint describe --name <endpoint-name>

# Check handler task queue in handler namespace
temporal task-queue describe --task-queue <handler-tq> --namespace <handler-ns>
```

**Common Nexus Issues:**

| Issue | Cause | Solution |
|-------|-------|----------|
| Endpoint not found | Endpoint not created or wrong name | Create/verify endpoint |
| Handler not responding | No workers on handler task queue | Start handler workers |
| Operation timeout | `scheduleToCloseTimeout` too short | Increase caller timeout |
| Handler error | Bug in handler operation code | Fix handler code |
| Cross-namespace auth | Permissions not configured | Configure namespace access |

## Recovery Actions

### Reset Workflow

Reset to retry from a specific point:

```bash
# Reset to specific event
temporal workflow reset \
  --workflow-id <id> \
  --event-id <event-id> \
  --reason "Reset for retry"

# Reset to last workflow task
temporal workflow reset \
  --workflow-id <id> \
  --type LastWorkflowTask \
  --reason "Reset after fix"
```

### Terminate Workflow

Force stop a stuck workflow:

```bash
temporal workflow terminate \
  --workflow-id <id> \
  --reason "Manual termination - issue description"
```

### Cancel Workflow

Request graceful cancellation:

```bash
temporal workflow cancel --workflow-id <id>
```

## Prevention

- Use replay tests in CI/CD
- Monitor key metrics
- Set appropriate timeouts
- Add heartbeats to long activities
- Use GetVersion for workflow changes
- Test with the production history

## Additional Resources

### Reference Files

For detailed error catalogs, consult:

- **`references/error-catalog.md`** - Complete error reference
- **`references/diagnostic-queries.md`** - Prometheus queries
