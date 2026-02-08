---
title: "CLI Reference"
weight: 1
---

# Timelord CLI Reference

Complete reference for all timelord-cli commands.

## Global Options

| Option | Short | Description | Default |
|--------|-------|-------------|---------|
| `--json` | `-j` | Output in JSON format | false |
| `--help` | `-h` | Show help | - |

## scaffold

Generate Temporal project scaffolding.

### scaffold project

Create a new Temporal project with complete structure.

```bash
timelord scaffold project [name] [flags]
```

**Arguments:**

| Argument | Description |
|----------|-------------|
| `name` | Project name (required) |

**Flags:**

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--output` | `-o` | Output directory | `.` |

**Example:**

```bash
timelord scaffold project order-service --output ./services
```

**Output:**

```
order-service/
├── go.mod
├── main.go
├── workflows/
│   └── example_workflow.go
├── activities/
│   └── example_activity.go
├── worker/
│   └── main.go
├── Dockerfile
└── Makefile
```

### scaffold workflow

Generate a workflow file.

```bash
timelord scaffold workflow [name] [flags]
```

**Arguments:**

| Argument | Description |
|----------|-------------|
| `name` | Workflow name (required) |

**Flags:**

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--output` | `-o` | Output file path | `./workflows/<name>_workflow.go` |

**Example:**

```bash
timelord scaffold workflow OrderProcessing --output ./workflows/order.go
```

### scaffold activity

Generate an activity file.

```bash
timelord scaffold activity [name] [flags]
```

**Arguments:**

| Argument | Description |
|----------|-------------|
| `name` | Activity name (required) |

**Flags:**

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--output` | `-o` | Output file path | `./activities/<name>_activity.go` |

**Example:**

```bash
timelord scaffold activity SendEmail --output ./activities/email.go
```

### scaffold worker

Generate a worker file.

```bash
timelord scaffold worker [name] [flags]
```

**Arguments:**

| Argument | Description |
|----------|-------------|
| `name` | Worker name (required) |

**Flags:**

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--output` | `-o` | Output file path | `./worker/main.go` |
| `--task-queue` | `-q` | Task queue name | `default-queue` |

**Example:**

```bash
timelord scaffold worker orders --task-queue order-processing
```

## cluster

Manage and inspect Temporal clusters.

### cluster status

Check cluster health and connectivity.

```bash
timelord cluster status [flags]
```

**Flags:**

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--address` | `-a` | Temporal server address | `localhost:7233` |

**Example:**

```bash
timelord cluster status --address temporal.prod.example.com:7233
```

**Output (text):**

```
Cluster Status
==============
Status: SERVING
Address: localhost:7233
```

**Output (JSON):**

```json
{
  "status": "SERVING",
  "address": "localhost:7233"
}
```

### cluster info

Get cluster configuration details.

```bash
timelord cluster info [flags]
```

**Flags:**

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--address` | `-a` | Temporal server address | `localhost:7233` |

**Example:**

```bash
timelord cluster info --json
```

### cluster metrics

Display key cluster metrics summary.

```bash
timelord cluster metrics [flags]
```

**Flags:**

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--address` | `-a` | Temporal server address | `localhost:7233` |

**Example:**

```bash
timelord cluster metrics
```

## namespace

Manage Temporal namespaces.

### namespace list

List all namespaces.

```bash
timelord namespace list [flags]
```

**Flags:**

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--address` | `-a` | Temporal server address | `localhost:7233` |

**Example:**

```bash
timelord namespace list --json
```

**Output:**

```
Namespaces
----------
  - default
  - orders
  - payments
```

### namespace create

Create a new namespace.

```bash
timelord namespace create [name] [flags]
```

**Arguments:**

| Argument | Description |
|----------|-------------|
| `name` | Namespace name (required) |

**Flags:**

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--address` | `-a` | Temporal server address | `localhost:7233` |
| `--retention` | `-r` | Workflow execution retention period | `72h` |
| `--description` | `-d` | Namespace description | - |

**Example:**

```bash
timelord namespace create orders --retention 168h --description "Order processing workflows"
```

### namespace describe

Show namespace details.

```bash
timelord namespace describe [name] [flags]
```

**Arguments:**

| Argument | Description |
|----------|-------------|
| `name` | Namespace name (required) |

**Flags:**

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--address` | `-a` | Temporal server address | `localhost:7233` |

**Example:**

```bash
timelord namespace describe orders --json
```

## nexus

Manage Nexus endpoints for cross-namespace communication.

### Nexus Commands

| Command | Description |
|---------|-------------|
| `timelord nexus list` | List all Nexus endpoints |
| `timelord nexus create <name>` | Create a Nexus endpoint |
| `timelord nexus describe <name>` | Show endpoint details |
| `timelord nexus delete <name>` | Delete a Nexus endpoint |

**Create endpoint:**

```bash
timelord nexus create my-endpoint \
  --target-namespace handler-ns \
  --target-task-queue handler-tq
```

**Scaffold Nexus service:**

```bash
timelord scaffold nexus-service PaymentService --output ./handlers
```

**Scaffold Nexus caller:**

```bash
timelord scaffold nexus-caller OrderPayment \
  --endpoint payments-endpoint \
  --service-name payment-service \
  --output ./workflows
```

## workflow

Manage and inspect workflow executions.

### workflow list

List workflow executions.

```bash
timelord workflow list [flags]
```

**Flags:**

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--namespace` | `-n` | Temporal namespace | `default` |
| `--query` | `-q` | List filter query | - |
| `--limit` | `-l` | Maximum workflows to return | `10` |

**Example:**

```bash
# List running workflows
timelord workflow list --query "ExecutionStatus='Running'" --limit 20

# List failed workflows
timelord workflow list --query "ExecutionStatus='Failed'" --json

# List by workflow type
timelord workflow list --query "WorkflowType='OrderWorkflow'"
```

**Query Syntax:**

Common query patterns:

| Query | Description |
|-------|-------------|
| `ExecutionStatus='Running'` | Running workflows |
| `ExecutionStatus='Failed'` | Failed workflows |
| `ExecutionStatus='Completed'` | Completed workflows |
| `WorkflowType='Name'` | By workflow type |
| `StartTime > '2024-01-01'` | Started after date |

### workflow describe

Show workflow execution details.

```bash
timelord workflow describe [workflow-id] [flags]
```

**Arguments:**

| Argument | Description |
|----------|-------------|
| `workflow-id` | Workflow ID (required) |

**Flags:**

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--namespace` | `-n` | Temporal namespace | `default` |
| `--run-id` | - | Run ID (optional) | - |

**Example:**

```bash
timelord workflow describe order-12345 --json
timelord workflow describe order-12345 --run-id abc-123
```

**Output:**

```
Workflow: order-12345
----------------------------
Status: Running
Type: OrderWorkflow
Start Time: 2024-01-15T10:30:00Z
Task Queue: order-processing
History Length: 42 events
```

### workflow history

Display workflow event history.

```bash
timelord workflow history [workflow-id] [flags]
```

**Arguments:**

| Argument | Description |
|----------|-------------|
| `workflow-id` | Workflow ID (required) |

**Flags:**

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--namespace` | `-n` | Temporal namespace | `default` |
| `--run-id` | - | Run ID (optional) | - |

**Example:**

```bash
timelord workflow history order-12345
timelord workflow history order-12345 --json
```

### workflow diagnose

Analyze workflow execution for issues.

```bash
timelord workflow diagnose [workflow-id] [flags]
```

**Arguments:**

| Argument | Description |
|----------|-------------|
| `workflow-id` | Workflow ID (required) |

**Flags:**

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--namespace` | `-n` | Temporal namespace | `default` |
| `--run-id` | - | Run ID (optional) | - |

**Example:**

```bash
timelord workflow diagnose order-12345 --json
```

**Output (text):**

```
Diagnosis for Workflow: order-12345
================================
Status: Running
Last Event: ActivityTaskScheduled

Pending Items:
  - 2 pending activity(ies)

Issues Found:
  [WARNING] ActivityTimeout: 1 activity(ies) timed out

Suggestions:
  - Consider increasing activity timeout
  - Add heartbeats for long-running activities
  - Check worker connectivity
```

**Output (JSON):**

```json
{
  "workflowId": "order-12345",
  "status": "Running",
  "issues": [
    {
      "severity": "WARNING",
      "type": "ActivityTimeout",
      "description": "1 activity(ies) timed out"
    }
  ],
  "suggestions": [
    "Consider increasing activity timeout",
    "Add heartbeats for long-running activities",
    "Check worker connectivity"
  ],
  "lastEvent": "ActivityTaskScheduled",
  "pendingItems": [
    "2 pending activity(ies)"
  ]
}
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `TEMPORAL_ADDRESS` | Default Temporal server address | `localhost:7233` |
| `TEMPORAL_NAMESPACE` | Default namespace | `default` |
| `TEMPORAL_TLS_CERT` | Path to TLS certificate | - |
| `TEMPORAL_TLS_KEY` | Path to TLS key | - |
| `TEMPORAL_TLS_CA` | Path to CA certificate | - |

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Invalid arguments |
| 3 | Connection error |
| 4 | Not found |

## See Also

- [Temporal CLI Documentation](https://docs.temporal.io/cli)
- [Troubleshooting Guide](/reference/troubleshooting)
