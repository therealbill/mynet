---
title: "Set Up Nexus Endpoints"
weight: 5
---

# Set Up Nexus Endpoints

Configure Nexus endpoints for cross-namespace communication between Temporal services.

## Problem

You have workflows in separate namespaces that need to communicate with each other using durable execution guarantees.

## Solution

Create Nexus endpoints that route calls from caller namespaces to handler namespaces, then implement handler services and caller workflows.

## Prerequisites

- Temporal server running (local dev or production)
- At least two namespaces (caller and handler)
- Access to `temporal` CLI

## Steps

### 1. Create Namespaces

```bash
temporal operator namespace create --namespace caller-ns
temporal operator namespace create --namespace handler-ns
```

### 2. Create a Nexus Endpoint

```bash
temporal operator nexus endpoint create \
  --name my-endpoint \
  --target-namespace handler-ns \
  --target-task-queue handler-tq
```

### 3. Verify the Endpoint

```bash
temporal operator nexus endpoint list
temporal operator nexus endpoint describe --name my-endpoint
```

### 4. Implement the Handler Service (Go)

```go
import (
    "context"
    "github.com/nexus-rpc/sdk-go/nexus"
    temporalnexus "go.temporal.io/sdk/temporalnexus"
    "go.temporal.io/sdk/client"
    "go.temporal.io/sdk/worker"
    "go.temporal.io/sdk/workflow"
)

// Sync operation for quick requests (< 10s)
var EchoOp = nexus.NewSyncOperation("echo",
    func(ctx context.Context, input EchoInput, opts nexus.StartOperationOptions) (EchoOutput, error) {
        return EchoOutput{Message: input.Message}, nil
    })

// Async operation backed by a workflow
var ProcessOp = temporalnexus.NewWorkflowRunOperation("process",
    ProcessWorkflow,
    func(ctx context.Context, input ProcessInput, opts nexus.StartOperationOptions) (client.StartWorkflowOptions, error) {
        return client.StartWorkflowOptions{ID: opts.RequestID}, nil
    })

// Register on handler worker
func main() {
    c, _ := client.Dial(client.Options{Namespace: "handler-ns"})
    w := worker.New(c, "handler-tq", worker.Options{})

    service := nexus.NewService("my-service")
    service.Register(EchoOp)
    service.Register(ProcessOp)
    w.RegisterNexusService(service)
    w.RegisterWorkflow(ProcessWorkflow)

    w.Run(worker.InterruptCh())
}
```

### 5. Implement the Caller Workflow (Go)

```go
func CallerWorkflow(ctx workflow.Context, input CallerInput) (*CallerOutput, error) {
    nexusClient := workflow.NewNexusClient("my-endpoint", "my-service")

    future := nexusClient.ExecuteOperation(ctx, ProcessOp, ProcessInput{
        ID: input.ID,
    }, workflow.NexusOperationOptions{
        ScheduleToCloseTimeout: 10 * time.Minute,
    })

    var result ProcessOutput
    if err := future.Get(ctx, &result); err != nil {
        return nil, fmt.Errorf("nexus operation failed: %w", err)
    }

    return &CallerOutput{Result: result}, nil
}
```

### 6. Start Both Workers

Start the handler worker first (connects to handler-ns), then the caller worker (connects to caller-ns).

## Verification

- [ ] Endpoint created and visible in `temporal operator nexus endpoint list`
- [ ] Handler worker running and registered with Nexus service
- [ ] Caller workflow can execute Nexus operations successfully
- [ ] Cross-namespace communication works end-to-end

## Troubleshooting

**Endpoint not found:** Verify endpoint name matches exactly in caller workflow's `NewNexusClient` call.

**Handler not responding:** Check handler workers are running on the correct task queue in the handler namespace.

**Operation timeout:** Increase `ScheduleToCloseTimeout` or check handler workflow for issues.

## Related

- `nexus-operations` skill for detailed implementation patterns
- `nexus-decision-guide` skill for architecture evaluation
- [Handle Workflow Failures](../handle-workflow-failures/) for recovery strategies
