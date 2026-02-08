---
name: tl-status
description: Check Temporal cluster and workflow status
arguments:
  - name: target
    description: "What to check: cluster, workflow, or namespace (default: cluster)"
    required: false
  - name: id
    description: Workflow ID or namespace name (when checking specific items)
    required: false
---

# Check Temporal Status

Check the status of Temporal clusters, namespaces, and workflows.

## Usage

```
/tl-status [target] [id]
```

## Targets

| Target | Description |
|--------|-------------|
| `cluster` | Check cluster health and configuration (default) |
| `namespace` | List namespaces or show namespace details |
| `workflow` | Show workflow execution status |
| `nexus` | List Nexus endpoints or show endpoint details |

## Examples

**Check cluster status:**
```
/tl-status
/tl-status cluster
```

**List namespaces:**
```
/tl-status namespace
```

**Show namespace details:**
```
/tl-status namespace production
```

**Check workflow status:**
```
/tl-status workflow my-workflow-id
```

**List Nexus endpoints:**
```
/tl-status nexus
```

**Show Nexus endpoint details:**
```
/tl-status nexus my-endpoint
```

## Execution

Use the Temporal CLI (`temporal`) to check status:

**Cluster status:**
```bash
temporal operator cluster health
temporal operator cluster describe
```

**Namespace operations:**
```bash
# List namespaces
temporal operator namespace list

# Describe namespace
temporal operator namespace describe --namespace {{id}}
```

**Workflow status:**
```bash
temporal workflow describe --workflow-id {{id}}
```

**Nexus endpoint operations:**
```bash
# List all endpoints
temporal operator nexus endpoint list

# Describe specific endpoint
temporal operator nexus endpoint describe --name {{id}}
```

## Prerequisites

Ensure the Temporal CLI is installed and configured:

```bash
# Check if temporal CLI is available
which temporal

# Configure server address if needed
export TEMPORAL_ADDRESS=localhost:7233
```

## Status Indicators

**Cluster Health:**

- `SERVING` - Cluster is healthy and accepting requests
- `NOT_SERVING` - Cluster is not accepting requests

**Workflow States:**

- `Running` - Workflow is actively executing
- `Completed` - Workflow finished successfully
- `Failed` - Workflow terminated with an error
- `Canceled` - Workflow was canceled
- `Terminated` - Workflow was forcefully terminated
- `TimedOut` - Workflow exceeded its execution timeout
- `ContinuedAsNew` - Workflow continued with new history

## Troubleshooting

If status checks fail:

1. Verify Temporal server is running
2. Check `TEMPORAL_ADDRESS` environment variable
3. Ensure network connectivity to the cluster
4. Verify namespace exists and you have access
