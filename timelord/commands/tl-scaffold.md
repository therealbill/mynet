---
name: tl-scaffold
description: Scaffold Temporal.io project components using timelord-cli
arguments:
  - name: type
    description: "Component type: project, workflow, activity, or worker"
    required: true
  - name: name
    description: Name for the component
    required: true
---

# Scaffold Temporal Components

Generate boilerplate code for Temporal.io applications.

## Usage

```
/tl-scaffold <type> <name>
```

## Component Types

| Type | Description |
|------|-------------|
| `project` | Full project structure with go.mod, Makefile, Dockerfile |
| `workflow` | Workflow file with input/output types and activity configuration |
| `activity` | Activity struct with heartbeat and cancellation handling |
| `worker` | Worker main file with registration and configuration |
| `nexus-service` | Nexus service handler with sync and async operation definitions |
| `nexus-caller` | Caller workflow with NexusClient setup for cross-namespace calls |

## Examples

**Create a new project:**
```
/tl-scaffold project myapp
```

**Create a workflow:**
```
/tl-scaffold workflow OrderProcessing
```

**Create an activity:**
```
/tl-scaffold activity Payment
```

**Create a worker:**
```
/tl-scaffold worker orders
```

**Create a Nexus service (handler side):**
```
/tl-scaffold nexus-service PaymentService
```

**Create a Nexus caller workflow:**
```
/tl-scaffold nexus-caller OrderPayment
```

## Execution

Run the timelord-cli scaffold command:

```bash
${CLAUDE_PLUGIN_ROOT}/tools/timelord-cli/bin/timelord scaffold {{type}} {{name}} --json
```

If the CLI is not built, build it first:

```bash
cd ${CLAUDE_PLUGIN_ROOT}/tools/timelord-cli && go build -o bin/timelord .
```

## Output

The command generates files in the current directory (or specified output directory):

- **project**: Creates a complete directory structure with all necessary files
- **workflow**: Creates `<name>_workflow.go` with workflow definition
- **activity**: Creates `<name>_activities.go` with activity struct and methods
- **worker**: Creates `<name>_worker.go` with worker setup
- **nexus-service**: Creates `<name>_nexus_service.go` with service definition, sync/async operations, and worker registration
- **nexus-caller**: Creates `<name>_nexus_caller.go` with NexusClient setup, ExecuteOperation pattern, and error handling

## Next Steps After Scaffolding

**For projects:**
1. `cd <project-name>`
2. `go mod tidy`
3. Implement your workflow logic
4. `make run-worker` to start the worker

**For individual components:**
1. Review and customize the generated code
2. Register workflows/activities with the worker
3. Configure task queues as needed
