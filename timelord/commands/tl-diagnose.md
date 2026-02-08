---
name: tl-diagnose
description: Diagnose workflow issues, analyze event history, and troubleshoot Temporal problems
arguments:
  - name: workflow-id
    description: Workflow ID to diagnose (optional - will prompt if not provided)
    required: false
  - name: issue
    description: Description of the issue you're experiencing
    required: false
---

# Temporal Workflow Diagnostics

Help diagnose and troubleshoot Temporal workflow issues.

## Diagnostic Process

When a user requests workflow diagnosis:

### 1. Gather Context

If no workflow ID provided, ask:
- What workflow is having issues?
- What symptoms are observed? (stuck, failing, slow, etc.)
- When did the issue start?

### 2. Run Diagnostics

Use the timelord-cli to gather information:

```bash
# If workflow ID is known
timelord workflow diagnose <workflow-id> --json

# Get workflow status
timelord workflow describe <workflow-id> --json

# Get event history
timelord workflow history <workflow-id> --json

# List recent workflows with issues
timelord workflow list --query "ExecutionStatus='Failed'" --json
timelord workflow list --query "ExecutionStatus='Running'" --json
```

### 3. Analyze Results

Based on the diagnostic output, check for:

**Stuck Workflows:**
- Pending activities waiting for workers
- Timers that haven't fired
- Waiting for signals
- Workflow task not being picked up

**Failed Workflows:**
- Activity failures and error messages
- Timeout issues
- Non-determinism errors
- Panic/crash in workflow code

**Performance Issues:**
- Long schedule-to-start latency
- Activity execution taking too long
- High retry counts

**Nexus Operation Issues:**
- NexusOperationTimedOut events → check handler and timeout config
- NexusOperationFailed events → check handler operation logic
- Endpoint not found → verify endpoint exists
- Handler not responding → check handler task queue workers

### 4. Provide Solutions

Use the troubleshooting skill for comprehensive guidance:

```
/skill troubleshooting
```

Common solutions to suggest:

| Issue | Solution |
|-------|----------|
| No workers | Start workers on correct task queue |
| Activity timeout | Increase timeout, add heartbeats |
| Non-determinism | Use GetVersion for code changes |
| Rate limiting | Add backoff, reduce concurrency |
| Nexus timeout | Increase scheduleToCloseTimeout, check handler |
| Nexus endpoint not found | Create endpoint with correct name |
| Nexus handler error | Debug handler operation/workflow |

### 5. Recovery Actions

If needed, suggest recovery commands:

```bash
# Reset workflow to retry
temporal workflow reset \
  --workflow-id <id> \
  --type LastWorkflowTask \
  --reason "Reset after fix"

# Terminate stuck workflow
temporal workflow terminate \
  --workflow-id <id> \
  --reason "Manual termination"

# Cancel workflow gracefully
temporal workflow cancel --workflow-id <id>
```

## Quick Diagnostics

For rapid triage, run:

```bash
# Check cluster health
timelord cluster status --json

# Check task queue
temporal task-queue describe --task-queue <queue>

# List problematic workflows
timelord workflow list --query "ExecutionStatus='Failed'" --limit 5 --json
```

## Related Skills

For deeper guidance, load these skills:

- **troubleshooting** - Complete error catalog and diagnosis trees
- **versioning-guide** - Help with non-determinism and workflow updates

## Related Agents

For complex debugging scenarios, use:

- **temporal-debug** - Specialized debugging agent for deep analysis

## Example Session

User: "My order workflow has been running for 2 hours"

1. Get workflow details:
   ```bash
   timelord workflow diagnose order-12345 --json
   ```

2. Analyze the diagnosis output for:
   - Pending activities
   - Last event type
   - Any failures or timeouts

3. Check if workers are running:
   ```bash
   temporal task-queue describe --task-queue order-processing
   ```

4. Provide targeted solution based on findings
