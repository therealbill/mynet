# Temporal Nexus Support for Timelord Plugin — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.
>
> **FIRST:** Copy this plan to `docs/plans/2026-02-07-nexus-support.md` in the timelord project (`/Users/bill/Projects/claude-plugins/timelord/docs/plans/2026-02-07-nexus-support.md`) before starting implementation. Create the `docs/plans/` directory if it doesn't exist.

**Goal:** Add comprehensive Temporal Nexus (GA) support to the timelord plugin — a new skill, a decision-guidance skill, updates to 4 existing skills, 3 agents, 4 commands, CLI tool extensions, Go templates, docs, and top-level metadata.

**Architecture:** Nexus is Temporal's first-class cross-namespace communication primitive (GA in Go & Java, preview in Python, experimental in TS/.NET). The plugin will treat it as a peer to workflows, activities, signals, queries, and updates — with strong emphasis on helping users decide *when* Nexus is the right tool and *when* it isn't.

**Tech Stack:** Markdown (skills/agents/commands), Go (CLI tool + templates)

**Skill quality standard:** All new and modified skills must follow [Anthropic's Complete Guide to Building Skills for Claude](https://resources.anthropic.com/hubfs/The-Complete-Guide-to-Building-Skill-for-Claude.pdf). Key requirements: description under 1024 chars with `[What] + [When] + [Key capabilities]` structure; SKILL.md under 5,000 words using progressive disclosure (heavy content in `references/`); no XML tags in frontmatter; include error handling and examples; use negative triggers to prevent over-triggering between related skills.

**Brain files:** Do NOT directly modify `.claude/brain/` files. These are managed by the brain skills plugin and will auto-update on the next scan.

---

## Task 1: Create the `nexus-operations` skill

**Executor:** Use `plugin-dev:skill-development` skill for structure guidance, then write the files.

The core knowledge artifact. All other changes reference this.

**Files:**

- Create: `skills/nexus-operations/SKILL.md`
- Create: `skills/nexus-operations/references/nexus-multi-sdk.md`

**Pattern to follow:** `skills/signals-queries-updates/SKILL.md` — YAML frontmatter with `name`, `description`, `version`, then markdown. Keep SKILL.md under 5,000 words per Anthropic skill guide; move TypeScript/Python/Java examples to `references/nexus-multi-sdk.md` (progressive disclosure level 3).

**Step 1: Create skill directory and SKILL.md**

````markdown
---
name: nexus-operations
description: Temporal Nexus implementation guidance for cross-namespace durable communication. Use when user asks about "nexus", "nexus operation", "nexus service", "nexus endpoint", "cross-namespace communication", "nexus caller", "nexus handler", "NexusClient", "ExecuteOperation", "WorkflowRunOperation", or "multi-namespace". Covers Go SDK (GA), with TypeScript/Python/Java in references. Do NOT use for architecture decisions about whether to adopt Nexus — use nexus-decision-guide instead.
version: 1.0.0
---

# Nexus Operations

Cross-namespace durable communication using Temporal Nexus (GA).

## Overview

Nexus enables workflows in separate namespaces to call each other through typed service contracts with Temporal's full durable execution guarantees.

| Feature | Nexus | Child Workflows | Activities | Signals |
|---------|-------|-----------------|-----------|---------|
| Cross-namespace | Yes | No (same NS only) | No | No |
| Durable execution | Yes | Yes | No | N/A |
| Request-response | Yes | Yes | Yes | No (one-way) |
| Arbitrary duration | Yes | Yes | No (< 10s) | N/A |
| Loose coupling | Yes | No | Yes | N/A |
| Service contract | Typed interface | Parent-child | Implicit | N/A |

**Architecture:**

```
┌──────────────────┐   Nexus Endpoint   ┌──────────────────┐
│  Caller Namespace │──────────────────>│ Handler Namespace │
│  ┌──────────────┐│   (routes to       │┌────────────────┐│
│  │Caller Workflow││   target NS +     ││Handler Worker   ││
│  │  NexusClient  ││   task queue)     ││ NexusService    ││
│  └──────────────┘│                    │└────────────────┘│
└──────────────────┘                    └──────────────────┘
```

## When to Use Nexus

**Use Nexus when:**

- Teams maintain separate namespaces and need to call each other's workflows
- You need durable execution guarantees across namespace boundaries
- You want clean, typed service contracts between teams/services
- You need both synchronous (< 10s) and asynchronous (long-running) cross-namespace calls

**Don't use Nexus when:**

- All workflows live in one namespace — use child workflows instead
- You need simple fire-and-forget notifications — use signals
- You need quick external API calls from workflows — use activities
- You need synchronous state mutation on a running workflow — use updates
- The added complexity (endpoints, handler workers, service contracts) isn't justified by the isolation benefits

## Choosing the Right Pattern

| Requirement | Best Option |
|-------------|-------------|
| Same namespace, simple decomposition | Child Workflows |
| External side effects, no durability needed | Activities |
| One-way notification to running workflow | Signals |
| Synchronous state mutation | Updates |
| Cross-namespace, durable, request-response | **Nexus** |
| Cross-namespace, short request (< 10s) | **Nexus sync operation** |
| Cross-namespace, long-running process | **Nexus async operation** |

## SDK Support

| SDK | Status | Notes |
|-----|--------|-------|
| Go | GA | Full production support |
| Java | GA | Full production support (1.28.0+) |
| Python | Public Preview | Recommended for production, API may improve |
| TypeScript | Experimental | APIs may have breaking changes |
| .NET | Experimental | APIs may have breaking changes (1.9.0+) |

## Key Concepts

**Nexus Endpoint:** A cluster-level routing resource that maps a name to a target namespace + task queue. Decouples callers from handler location. Created via CLI or Terraform.

**Nexus Service:** A named collection of Nexus Operations. Defines the microservice contract. Registered on handler workers.

**Nexus Operation:** An arbitrary-duration operation — either:
- **Synchronous**: RPC-style, completes immediately (< 10s), handler returns result directly
- **Asynchronous**: Starts a workflow as the backing operation, returns an operation token for tracking

**Nexus Machinery:** Temporal's built-in infrastructure providing at-least-once execution, automatic retries, circuit breaking, and rate limiting for Nexus calls.

## Endpoint Setup

Prerequisites: caller and handler namespaces must exist.

```bash
# Create namespaces (if not already present)
temporal operator namespace create --namespace handler-ns
temporal operator namespace create --namespace caller-ns

# Create Nexus endpoint
temporal operator nexus endpoint create \
  --name my-endpoint \
  --target-namespace handler-ns \
  --target-task-queue handler-tq

# List endpoints
temporal operator nexus endpoint list

# Describe endpoint
temporal operator nexus endpoint describe --name my-endpoint

# Update endpoint
temporal operator nexus endpoint update \
  --name my-endpoint \
  --target-task-queue new-handler-tq

# Delete endpoint
temporal operator nexus endpoint delete --name my-endpoint
```

## Handler Implementation (Go SDK)

### Synchronous Operation

For requests completing in < 10 seconds:

```go
import (
    "context"
    "github.com/nexus-rpc/sdk-go/nexus"
)

// Synchronous operation — returns result directly
var EchoOp = nexus.NewSyncOperation("echo",
    func(ctx context.Context, input EchoInput, opts nexus.StartOperationOptions) (EchoOutput, error) {
        return EchoOutput{Message: input.Message}, nil
    })
```

### Asynchronous (Workflow-Backed) Operation

For long-running processes:

```go
import (
    "context"
    "go.temporal.io/sdk/client"
    temporalnexus "go.temporal.io/sdk/temporalnexus"
)

// Async operation — starts a workflow, returns when workflow completes
var ProcessOrderOp = temporalnexus.NewWorkflowRunOperation("process-order",
    ProcessOrderWorkflow,
    func(ctx context.Context, input OrderInput, opts nexus.StartOperationOptions) (client.StartWorkflowOptions, error) {
        return client.StartWorkflowOptions{
            // Use RequestID as workflow ID to deduplicate retries
            ID: opts.RequestID,
        }, nil
    })
```

### Service Definition and Worker Registration

```go
import (
    "github.com/nexus-rpc/sdk-go/nexus"
    "go.temporal.io/sdk/worker"
)

// Create and register service
func registerNexusService(w worker.Worker) {
    service := nexus.NewService("order-service")
    service.Register(EchoOp)
    service.Register(ProcessOrderOp)
    w.RegisterNexusService(service)
}
```

## Other SDK Implementations

For TypeScript (experimental), Python (public preview), and Java (GA) handler and caller examples, see `references/nexus-multi-sdk.md`.

## Caller Workflow Patterns (Go SDK)

```go
func CallerWorkflow(ctx workflow.Context, input CallerInput) (*CallerOutput, error) {
    // Create Nexus client — endpoint name + service name
    nexusClient := workflow.NewNexusClient("my-endpoint", "order-service")

    // Execute operation (blocks until result)
    future := nexusClient.ExecuteOperation(ctx, ProcessOrderOp, OrderInput{
        ID: input.OrderID,
    }, workflow.NexusOperationOptions{
        ScheduleToCloseTimeout: 10 * time.Minute,
    })

    var result OrderOutput
    if err := future.Get(ctx, &result); err != nil {
        return nil, fmt.Errorf("nexus operation failed: %w", err)
    }

    return &CallerOutput{OrderResult: result}, nil
}
```

## Error Handling

### Error Types

| Error Type | Where | Description | Retryable? |
|------------|-------|-------------|------------|
| `OperationError` | Handler | Application-level failure | No |
| `HandlerError` | Handler | Framework/infrastructure error | Depends on type |
| `NexusOperationFailure` | Caller | Wraps handler error for caller workflow | No |

### Handling Errors in Callers (Go)

```go
var result OrderOutput
if err := future.Get(ctx, &result); err != nil {
    var nexusErr *temporal.NexusOperationFailure
    if errors.As(err, &nexusErr) {
        // Application-level failure from handler
        logger.Error("Nexus operation failed", "error", nexusErr)
        return nil, err
    }
    // Other errors (timeout, cancellation)
    return nil, err
}
```

### Returning Errors from Handlers (Go)

```go
var ProcessOp = nexus.NewSyncOperation("process",
    func(ctx context.Context, input Input, opts nexus.StartOperationOptions) (Output, error) {
        if input.ID == "" {
            // Return application error — caller gets NexusOperationFailure
            return Output{}, nexus.HandlerErrorf(nexus.HandlerErrorTypeBadRequest, "ID is required")
        }
        return Output{Result: "ok"}, nil
    })
```

## Event History

Nexus-specific events in caller workflow history:

| Event | Meaning |
|-------|---------|
| `NexusOperationScheduled` | Caller scheduled a Nexus operation |
| `NexusOperationStarted` | Async operation accepted by handler (returned operation token) |
| `NexusOperationCompleted` | Operation finished successfully |
| `NexusOperationFailed` | Operation returned an error |
| `NexusOperationCanceled` | Operation was canceled |
| `NexusOperationTimedOut` | `scheduleToCloseTimeout` exceeded |

## Best Practices

- **Use sync operations for < 10s requests** — simpler, lower overhead
- **Use async (workflow-backed) for long-running work** — gets full workflow durability
- **Deduplicate with RequestID** — use `opts.RequestID` as workflow ID in async handlers to safely handle retries
- **Implement idempotent handlers** — Nexus provides at-least-once semantics, so handlers may be called more than once
- **Always set `scheduleToCloseTimeout`** on caller side — without it, operations may hang indefinitely
- **Treat Nexus service contracts as stable APIs** — they're microservice boundaries; version them carefully
- **Use separate namespaces** for caller and handler to get true isolation benefits
- **Don't use Nexus within a single namespace** — child workflows are simpler and equally durable

## Additional Resources

### Reference Files

For detailed patterns, consult:

- **`references/nexus-patterns.md`** - Advanced Nexus usage patterns
- **`references/nexus-multi-sdk.md`** - Multi-SDK implementation examples
````

**Step 2: Create `skills/nexus-operations/references/nexus-multi-sdk.md`**

This file contains the TypeScript, Python, and Java SDK examples (progressive disclosure level 3). Include:

- TypeScript handler: service contract definition (`nexus.service()`), sync handler, async handler (`WorkflowRunOperationHandler`)
- TypeScript caller: `wf.createNexusClient()` + `executeOperation()`
- Python handler: `@nexusrpc.service`, `@nexusrpc.handler.sync_operation`, `@temporal_nexus.workflow_run_operation`
- Java handler: `@ServiceImpl`, `@OperationImpl`, `OperationHandler.sync()`, `WorkflowRunOperationHandler.fromWorkflowMethod()`

Use the exact code examples from the original plan (preserved below for reference):

**TypeScript handler service contract:**
```typescript
import * as nexus from 'nexus-rpc';
export const orderService = nexus.service('order-service', {
    echo: nexus.operation<EchoInput, EchoOutput>(),
    processOrder: nexus.operation<OrderInput, OrderOutput>(),
});
```

**TypeScript handler implementations:**
```typescript
import * as temporalNexus from '@temporalio/nexus';
const handlers = {
    echo: async (ctx, input: EchoInput): Promise<EchoOutput> => {
        return { message: input.message };
    },
    processOrder: new temporalNexus.WorkflowRunOperationHandler<OrderInput, OrderOutput>(
        async (ctx, input: OrderInput) => {
            return await temporalNexus.startWorkflow(ctx, processOrderWorkflow, {
                args: [input],
                workflowId: ctx.requestId ?? randomUUID(),
            });
        }
    ),
};
```

**TypeScript caller:**
```typescript
const nexusClient = wf.createNexusClient({ service: orderService, endpoint: 'my-endpoint' });
const result = await nexusClient.executeOperation('processOrder', { id: input.orderId }, { scheduleToCloseTimeout: '10m' });
```

**Python handler:**
```python
@nexusrpc.service
class OrderService:
    process_order: nexusrpc.Operation[OrderInput, OrderOutput]
    echo: nexusrpc.Operation[EchoInput, EchoOutput]

@nexusrpc.handler.service_handler(service=OrderService)
class OrderServiceHandler:
    @nexusrpc.handler.sync_operation
    async def echo(self, ctx, input: EchoInput) -> EchoOutput:
        return EchoOutput(message=input.message)

    @temporal_nexus.workflow_run_operation
    async def process_order(self, ctx, input: OrderInput):
        return await temporal_nexus.start_workflow(ctx, ProcessOrderWorkflow, args=[input])
```

**Java handler:**
```java
@ServiceImpl(service = OrderNexusService.class)
public class OrderNexusServiceImpl {
    @OperationImpl
    public OperationHandler<EchoInput, EchoOutput> echo() {
        return OperationHandler.sync((ctx, input) -> new EchoOutput(input.getMessage()));
    }
    @OperationImpl
    public OperationHandler<OrderInput, OrderOutput> processOrder() {
        return WorkflowRunOperationHandler.fromWorkflowMethod(
            (ctx, input) -> Workflow.newWorkflowStub(ProcessOrderWorkflow.class)::processOrder);
    }
}
```

**Step 3: Commit**

```bash
git add skills/nexus-operations/
git commit -m "feat: add nexus-operations skill for cross-namespace communication"
```

---

## Task 2: Create the `nexus-decision-guide` skill

**Executor:** Use `plugin-dev:skill-development` skill for structure guidance, then write the files.

A decision-focused skill that helps engineers evaluate whether Nexus is right for their architecture.

**Files:**

- Create: `skills/nexus-decision-guide/SKILL.md`

**Step 1: Create skill directory and SKILL.md**

````markdown
---
name: nexus-decision-guide
description: Architecture decision framework for evaluating whether Temporal Nexus is right for your use case. Use when user asks "should I use nexus", "nexus vs child workflows", "nexus vs activities", "when to use nexus", "cross-namespace pattern", "multi-namespace architecture", "nexus tradeoffs", or "nexus benefits". Covers tradeoff analysis, complexity scoring, and migration paths. Do NOT use for Nexus implementation details or code examples — use nexus-operations instead.
version: 1.0.0
---

# Nexus Decision Guide

Guidance for deciding whether Temporal Nexus is the right communication pattern for your architecture.

## Decision Framework

### Step 1: Do You Need Cross-Namespace Communication?

If all your workflows run in a single namespace, **stop here — you don't need Nexus.**

- Same-namespace decomposition → **Child Workflows**
- Same-namespace messaging → **Signals, Queries, or Updates**
- Same-namespace external calls → **Activities**

### Step 2: What Drives the Namespace Separation?

| Reason | Nexus Fit |
|--------|-----------|
| Team ownership boundaries | Strong — Nexus provides clean service contracts |
| Security/compliance isolation | Strong — namespaces + Nexus maintain isolation |
| Independent scaling | Strong — separate workers, separate scaling |
| Environment separation (dev/staging/prod) | Not applicable — use same namespace structure per env |
| Just organizational preference | Weak — consider if the operational overhead is worth it |

### Step 3: What Are the Communication Requirements?

| Need | Best Pattern |
|------|-------------|
| Fire-and-forget notification | Signals (if same NS) or Activities calling Signal API |
| Quick data lookup (< 10s) | Nexus sync operation |
| Long-running cross-team process | Nexus async (workflow-backed) operation |
| Cross-namespace, no durability needed | Activity with Temporal client (simpler, less durable) |
| Cross-namespace with full durability | **Nexus** (the primary use case) |

### Step 4: Evaluate the Tradeoffs

**Nexus adds:**

- Nexus Endpoints (cluster-level routing resources to create and manage)
- Handler workers (separate workers for Nexus service registration)
- Service contracts (typed interfaces to define and maintain)
- Cross-namespace debugging complexity (tracing across caller + handler)

**Nexus provides:**

- At-least-once execution with automatic retries and circuit breaking
- Clean API boundaries between teams
- Independent deployment and scaling per namespace
- Full durable execution across namespace boundaries

**The question:** Does the isolation and durability benefit outweigh the operational complexity for your case?

## Architecture Patterns

### Pattern: Service Mesh via Nexus

Multiple teams each owning a namespace, connected via Nexus endpoints:

```
┌──────────────┐   Nexus: payments-ep   ┌──────────────┐
│  orders-ns   │──────────────────────>│ payments-ns  │
│  (Team A)    │                        │ (Team B)     │
│              │   Nexus: inventory-ep  ┌──────────────┐
│              │──────────────────────>│ inventory-ns │
└──────────────┘                       │ (Team C)     │
                                       └──────────────┘
```

**Good when:** Teams need independence, clear ownership, separate deployment cycles.

### Pattern: Shared Platform Service

Common platform service (e.g., notifications, payments) exposed via Nexus to all consuming namespaces:

```
┌────────────┐
│ service-a  │──┐
└────────────┘  │   Nexus: notify-ep   ┌────────────────┐
                ├────────────────────>│ notifications  │
┌────────────┐  │                     │ (platform team)│
│ service-b  │──┘                     └────────────────┘
└────────────┘
```

**Good when:** A central capability is consumed by many teams.

### Anti-Pattern: Nexus Within Single Namespace

Don't use Nexus if caller and handler are in the same namespace. Child workflows give you the same decomposition with less overhead.

### Anti-Pattern: Nexus for Simple External Calls

If you just need to call an HTTP API and don't need cross-namespace durability, use a regular Activity. Nexus is for Temporal-to-Temporal communication.

## Complexity Assessment

Rate your scenario:

| Factor | Low (1) | Medium (2) | High (3) |
|--------|---------|------------|----------|
| Number of consuming namespaces | 1 | 2-3 | 4+ |
| Handler operation complexity | Simple sync | Mix of sync/async | Complex chains |
| Team independence requirement | Nice to have | Important | Critical |
| Existing namespace topology | Single NS | Few NS, ad-hoc | Multi-NS established |

**Score 4-6:** Consider if simpler patterns (child workflows, activities) suffice.
**Score 7-9:** Nexus is likely a good fit.
**Score 10-12:** Nexus is strongly recommended.

## Migration Path

If you're currently using ad-hoc cross-namespace communication (activities calling Temporal client, etc.):

1. Define Nexus service contracts for existing cross-namespace calls
2. Create Nexus endpoints pointing to handler namespaces
3. Implement handler services wrapping existing workflows
4. Migrate callers one at a time to use NexusClient
5. Remove old activity-based cross-namespace glue code

## Additional Resources

- **nexus-operations** skill — detailed implementation guidance
- **namespace-management** skill — namespace setup and Nexus endpoint management
````

**Step 2: Commit**

```bash
git add skills/nexus-decision-guide/SKILL.md
git commit -m "feat: add nexus-decision-guide skill for architecture evaluation"
```

---

## Task 3: Update `namespace-management` skill with Nexus endpoint management

**Executor:** Use `plugin-dev:skill-development` skill for guidance on modifications.

**Files:**

- Modify: `skills/namespace-management/SKILL.md`

**Step 1: Add Nexus row to Namespace Concepts table** (after line 21)

Add to the table:

```markdown
| Nexus Endpoints | Cross-namespace routing (endpoint → target NS + TQ) |
```

**Step 2: Add new section before "## Operations"** (before line 263)

Insert a new section "## Nexus Endpoint Management" containing:

- CLI commands: create, list, describe, update, delete endpoints
- Endpoint naming conventions (following existing `{team}-{service}-{environment}` pattern)
- Topology diagram showing Nexus connecting the multi-tenancy patterns already in the skill

**Step 3: Add "Namespace Per Service with Nexus" subsection** under "## Multi-Tenancy Patterns" (after the existing three patterns around line 183)

Show how the "Namespace Per Service" pattern pairs with Nexus endpoints:

````markdown
### Connecting Services with Nexus

When using namespace-per-service, Nexus endpoints connect them:

```
orders ──nexus: payments-ep──> payments
orders ──nexus: inventory-ep──> inventory
shipping ──nexus: inventory-ep──> inventory
```

Each endpoint routes to the handler namespace's task queue. Teams own their handler services independently.
````

**Step 4: Commit**

```bash
git add skills/namespace-management/SKILL.md
git commit -m "feat: add Nexus endpoint management to namespace-management skill"
```

---

## Task 4: Update `workflow-patterns` skill with Nexus pattern

**Executor:** Use `plugin-dev:skill-development` skill for guidance on modifications.

**Files:**

- Modify: `skills/workflow-patterns/SKILL.md`

**Step 1: Add Nexus to allowed operations list** (line 31, after local activities)

```markdown
- Nexus operations via `workflow.NewNexusClient()` and `client.ExecuteOperation()`
```

**Step 2: Add new pattern section** after "## Pattern: Cron/Scheduled" (after line 299)

````markdown
## Pattern: Cross-Namespace Communication (Nexus)

Invoke workflows in other namespaces through typed service contracts with durable execution guarantees.

**When to Use:**

- Workflows in different namespaces need to communicate
- Teams require namespace isolation but need to call each other's services
- You need durability guarantees across namespace boundaries

**When NOT to Use:**

- All workflows in the same namespace — use child workflows instead
- Simple external API calls — use activities
- Fire-and-forget messaging — use signals

**Go Implementation (Caller):**

```go
func OrderWorkflow(ctx workflow.Context, order Order) (*OrderResult, error) {
    nexusClient := workflow.NewNexusClient("payments-endpoint", "payment-service")

    future := nexusClient.ExecuteOperation(ctx, "charge", ChargeInput{
        OrderID: order.ID,
        Amount:  order.Total,
    }, workflow.NexusOperationOptions{
        ScheduleToCloseTimeout: 5 * time.Minute,
    })

    var chargeResult ChargeOutput
    if err := future.Get(ctx, &chargeResult); err != nil {
        return nil, fmt.Errorf("payment failed: %w", err)
    }

    return &OrderResult{PaymentID: chargeResult.TransactionID}, nil
}
```
````

**Step 3: Add row to pattern selection table** (line 338)

```markdown
| Cross-namespace durable communication | Nexus Operations |
```

**Step 4: Commit**

```bash
git add skills/workflow-patterns/SKILL.md
git commit -m "feat: add Nexus cross-namespace pattern to workflow-patterns skill"
```

---

## Task 5: Update `testing-strategies` skill with Nexus testing

**Executor:** Use `plugin-dev:skill-development` skill for guidance on modifications.

**Files:**

- Modify: `skills/testing-strategies/SKILL.md`

**Step 1: Add new section** after "## Testing Child Workflows" (after line 298)

````markdown
## Testing Nexus Operations

### Testing Sync Operation Handlers

Sync operations are plain functions — test them directly:

```go
func TestEchoOp(t *testing.T) {
    result, err := echoHandler(context.Background(), EchoInput{Message: "hello"}, nexus.StartOperationOptions{})
    require.NoError(t, err)
    require.Equal(t, "hello", result.Message)
}
```

### Testing Async Operation Handlers

Async operations start workflows — test the backing workflow as a normal workflow test, and test the operation handler setup separately.

### Testing Caller Workflows with Nexus Mocks

Mock Nexus operations in the test environment:

```go
func (s *WorkflowTestSuite) TestCallerWorkflow_NexusSuccess() {
    env := s.NewTestWorkflowEnvironment()

    // Register Nexus operation handler for testing
    env.RegisterNexusService(testNexusService)

    env.ExecuteWorkflow(CallerWorkflow, CallerInput{OrderID: "order-1"})

    s.True(env.IsWorkflowCompleted())
    s.NoError(env.GetWorkflowError())
}
```

### Integration Testing Nexus

Test against a real dev server with multiple namespaces:

```bash
# Start dev server
temporal server start-dev

# Create test namespaces and endpoint
temporal operator namespace create --namespace test-caller
temporal operator namespace create --namespace test-handler
temporal operator nexus endpoint create \
  --name test-endpoint \
  --target-namespace test-handler \
  --target-task-queue test-handler-tq
```

Then run both handler and caller workers, execute the caller workflow, and verify the end-to-end result.
````

**Step 2: Commit**

```bash
git add skills/testing-strategies/SKILL.md
git commit -m "feat: add Nexus testing patterns to testing-strategies skill"
```

---

## Task 6: Update `troubleshooting` skill with Nexus issues

**Executor:** Use `plugin-dev:skill-development` skill for guidance on modifications.

**Files:**

- Modify: `skills/troubleshooting/SKILL.md`

**Step 1: Add new issue section** after "### Issue: High Latency" (after line 266)

````markdown
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
````

**Step 2: Add to "Stuck Workflows" solutions table** (line 98)

```markdown
| Nexus operation | Check endpoint and handler workers |
```

**Step 3: Commit**

```bash
git add skills/troubleshooting/SKILL.md
git commit -m "feat: add Nexus troubleshooting to troubleshooting skill"
```

---

## Task 7: Update `temporal-dev` agent with Nexus development guidance

**Executor:** Use `plugin-dev:agent-development` skill for agent modification guidance.

**Files:**

- Modify: `agents/temporal-dev.md`

**Step 1: Add Nexus example to frontmatter** (after the versioning example, before `model: inherit` at line 41)

```markdown
  <example>
  Context: User wants to implement cross-namespace communication
  user: "I need my orders service to call the payments service in another namespace"
  assistant: "I'll use the temporal-dev agent to help you design a Nexus service for cross-namespace communication between your orders and payments services."
  <commentary>
  Nexus operations enable durable cross-namespace communication, which is the recommended pattern for service-to-service calls across namespace boundaries.
  </commentary>
  </example>
```

**Step 2: Add responsibility 6** (after line 54)

```markdown
6. **Nexus Operations**: Design cross-namespace services using Nexus endpoints, sync/async operations, and typed service contracts
```

**Step 3: Add Nexus SDK patterns section** (after the Testing Approach section, around line 142)

````markdown
**Nexus SDK Patterns:**

Handler service definition (Go):
```go
service := nexus.NewService("my-service")
service.Register(mySyncOp)
service.Register(myAsyncOp)
w.RegisterNexusService(service)
```

Caller workflow pattern (Go):
```go
nexusClient := workflow.NewNexusClient("endpoint-name", "service-name")
future := nexusClient.ExecuteOperation(ctx, operation, input, workflow.NexusOperationOptions{
    ScheduleToCloseTimeout: 10 * time.Minute,
})
```

Key design decisions:
- Use sync operations for requests completing < 10 seconds
- Use async (workflow-backed) operations for long-running processes
- Use `opts.RequestID` as workflow ID in async handlers to deduplicate retries
- Don't use Nexus within a single namespace — use child workflows instead
````

**Step 4: Add step 7 to Analysis Process** (after line 151)

```markdown
7. Evaluate cross-namespace requirements (Nexus vs child workflows vs activities)
```

**Step 5: Add to Common Patterns** (after line 168)

```markdown
- **Nexus Operations**: Cross-namespace durable communication with typed service contracts
```

**Step 6: Commit**

```bash
git add agents/temporal-dev.md
git commit -m "feat: add Nexus development guidance to temporal-dev agent"
```

---

## Task 8: Update `temporal-ops` agent with Nexus endpoint management

**Executor:** Use `plugin-dev:agent-development` skill for agent modification guidance.

**Files:**

- Modify: `agents/temporal-ops.md`

**Step 1: Add Nexus example to frontmatter** (after the security example, before `model: inherit`)

```markdown
  <example>
  Context: User needs to set up Nexus endpoints for cross-namespace communication
  user: "How do I configure Nexus endpoints between our orders and payments namespaces?"
  assistant: "I'll use the temporal-ops agent to guide you through creating Nexus endpoints with proper namespace configuration."
  <commentary>
  Nexus endpoint management is an operational concern involving namespace configuration, routing, and access control.
  </commentary>
  </example>
```

**Step 2: Add responsibility 7** (after line 55)

```markdown
7. **Nexus Endpoints**: Create and manage Nexus endpoints for cross-namespace routing
```

**Step 3: Add Nexus section** after the Namespace Management section (after line 249)

````markdown
**Nexus Endpoint Operations:**

```bash
# Create endpoint
temporal operator nexus endpoint create \
  --name payments-endpoint \
  --target-namespace payments-ns \
  --target-task-queue payments-tq

# List endpoints
temporal operator nexus endpoint list

# Update endpoint
temporal operator nexus endpoint update \
  --name payments-endpoint \
  --target-task-queue new-payments-tq

# Delete endpoint
temporal operator nexus endpoint delete --name payments-endpoint
```

**Nexus Topology Planning:**

```
┌──────────────┐    Nexus Endpoint    ┌──────────────┐
│  orders-ns   │───── payments ──────>│ payments-ns  │
│  (caller)    │───── inventory ─────>│ inventory-ns │
└──────────────┘                      └──────────────┘
```
````

**Step 4: Commit**

```bash
git add agents/temporal-ops.md
git commit -m "feat: add Nexus endpoint management to temporal-ops agent"
```

---

## Task 9: Update `temporal-debug` agent with Nexus debugging

**Executor:** Use `plugin-dev:agent-development` skill for agent modification guidance.

**Files:**

- Modify: `agents/temporal-debug.md`

**Step 1: Add Nexus example to frontmatter** (after the activity timeout example, before `model: inherit`)

```markdown
  <example>
  Context: User has a Nexus operation that isn't completing
  user: "My cross-namespace Nexus call to the payments service is timing out"
  assistant: "I'll use the temporal-debug agent to investigate the Nexus operation failure, checking both the caller and handler workflows, the endpoint configuration, and handler worker connectivity."
  <commentary>
  Nexus debugging requires analyzing both the caller and handler sides, plus the endpoint routing configuration.
  </commentary>
  </example>
```

**Step 2: Add Nexus events to Event History Analysis table** (after line 207)

```markdown
| `NexusOperationScheduled` | Nexus operation queued by caller |
| `NexusOperationStarted` | Handler accepted the operation |
| `NexusOperationCompleted` | Operation finished successfully |
| `NexusOperationFailed` | Operation returned error |
| `NexusOperationCanceled` | Operation was canceled |
| `NexusOperationTimedOut` | Operation exceeded timeout |
```

**Step 3: Add Nexus to the event flow example** (after line 221)

````markdown
**Nexus Event Flow:**

```
WorkflowTaskCompleted
  └─> NexusOperationScheduled
      └─> NexusOperationStarted (async only)
          └─> NexusOperationCompleted
              └─> WorkflowTaskScheduled
```
````

**Step 4: Add Nexus debugging commands** (after Debugging Commands section around line 243)

````markdown
**Nexus-Specific Debugging:**

```bash
# Find Nexus events in caller workflow
temporal workflow show --workflow-id <caller-id> --output json | \
  jq '.events[] | select(.eventType | contains("Nexus"))'

# Check endpoint configuration
temporal operator nexus endpoint describe --name <endpoint-name>

# Check handler workers
temporal task-queue describe --task-queue <handler-tq> --namespace <handler-ns>

# If async operation: find handler workflow
# The handler workflow ID typically matches the caller's operation request ID
```
````

**Step 5: Commit**

```bash
git add agents/temporal-debug.md
git commit -m "feat: add Nexus debugging to temporal-debug agent"
```

---

## Task 10: Update `tl-scaffold` command with Nexus scaffolding

**Executor:** Use `plugin-dev:command-development` skill for command modification guidance.

**Files:**

- Modify: `commands/tl-scaffold.md`

**Step 1: Add nexus types to Component Types table** (after line 31)

```markdown
| `nexus-service` | Nexus service handler with sync and async operation definitions |
| `nexus-caller` | Caller workflow with NexusClient setup for cross-namespace calls |
```

**Step 2: Add examples** (after existing examples around line 51)

````markdown
**Create a Nexus service (handler side):**
```
/tl-scaffold nexus-service PaymentService
```

**Create a Nexus caller workflow:**
```
/tl-scaffold nexus-caller OrderPayment
```
````

**Step 3: Add output descriptions** (after line 76)

```markdown
- **nexus-service**: Creates `<name>_nexus_service.go` with service definition, sync/async operations, and worker registration
- **nexus-caller**: Creates `<name>_nexus_caller.go` with NexusClient setup, ExecuteOperation pattern, and error handling
```

**Step 4: Commit**

```bash
git add commands/tl-scaffold.md
git commit -m "feat: add Nexus scaffolding to tl-scaffold command"
```

---

## Task 11: Update `tl-status` command with Nexus endpoint status

**Executor:** Use `plugin-dev:command-development` skill for command modification guidance.

**Files:**

- Modify: `commands/tl-status.md`

**Step 1: Add nexus target to table** (after line 29)

```markdown
| `nexus` | List Nexus endpoints or show endpoint details |
```

**Step 2: Add examples** (after existing examples around line 52)

````markdown
**List Nexus endpoints:**
```
/tl-status nexus
```

**Show Nexus endpoint details:**
```
/tl-status nexus my-endpoint
```
````

**Step 3: Add execution commands** (after line 76)

````markdown
**Nexus endpoint operations:**
```bash
# List all endpoints
temporal operator nexus endpoint list

# Describe specific endpoint
temporal operator nexus endpoint describe --name {{id}}
```
````

**Step 4: Commit**

```bash
git add commands/tl-status.md
git commit -m "feat: add Nexus endpoint status to tl-status command"
```

---

## Task 12: Update `tl-diagnose` and `tl-test` commands

**Executor:** Use `plugin-dev:command-development` skill for command modification guidance.

**Files:**

- Modify: `commands/tl-diagnose.md`
- Modify: `commands/tl-test.md`

**Step 1: Add Nexus issues to tl-diagnose** (after line 66, in the "Analyze Results" section)

```markdown
**Nexus Operation Issues:**
- NexusOperationTimedOut events → check handler and timeout config
- NexusOperationFailed events → check handler operation logic
- Endpoint not found → verify endpoint exists
- Handler not responding → check handler task queue workers
```

And add to the solutions table (after line 83):

```markdown
| Nexus timeout | Increase scheduleToCloseTimeout, check handler |
| Nexus endpoint not found | Create endpoint with correct name |
| Nexus handler error | Debug handler operation/workflow |
```

**Step 2: Add Nexus checks to tl-test** (after the Activity Checks table at line 138)

```markdown
### Nexus Caller Checks

| Pattern | Issue | Fix |
|---------|-------|-----|
| Missing `scheduleToCloseTimeout` | Operation may hang indefinitely | Set `ScheduleToCloseTimeout` |
| Hardcoded endpoint name | Environment-specific | Use configuration/environment variable |
| No error handling on Nexus future | Lost failures | Handle `NexusOperationFailure` |
```

**Step 3: Commit**

```bash
git add commands/tl-diagnose.md commands/tl-test.md
git commit -m "feat: add Nexus diagnostics to tl-diagnose and tl-test commands"
```

---

## Task 13: Add CLI `nexus.go` command file

**Executor:** Use `go-architect` agent (subagent_type=go-architect) for Go implementation.

**Files:**

- Create: `tools/timelord-cli/cmd/nexus.go`

**Pattern to follow:** `tools/timelord-cli/cmd/namespace.go` — uses `runTemporalCLI()` helper from `cluster.go:89`, cobra command structure, JSON/text output modes.

**Step 1: Create `nexus.go`**

Implement:

- `nexusCmd` — parent command `timelord nexus`
- `nexusListCmd` — `timelord nexus list` → delegates to `temporal operator nexus endpoint list`
- `nexusCreateCmd` — `timelord nexus create <name> --target-namespace <ns> --target-task-queue <tq>` → delegates to `temporal operator nexus endpoint create`
- `nexusDescribeCmd` — `timelord nexus describe <name>` → delegates to `temporal operator nexus endpoint describe`
- `nexusDeleteCmd` — `timelord nexus delete <name>` → delegates to `temporal operator nexus endpoint delete`

Types:

```go
type NexusEndpointInfo struct {
    Name            string `json:"name"`
    TargetNamespace string `json:"targetNamespace"`
    TargetTaskQueue string `json:"targetTaskQueue"`
}

type NexusEndpointListResult struct {
    Endpoints []NexusEndpointInfo `json:"endpoints"`
    Count     int                 `json:"count"`
}
```

Follow the same patterns as `namespace.go`: `init()` registers with `rootCmd`, flags for `--target-namespace` and `--target-task-queue`, `runTemporalCLI()` for delegation, `jsonOutput` toggle for output format.

**Step 2: Commit**

```bash
git add tools/timelord-cli/cmd/nexus.go
git commit -m "feat: add timelord nexus CLI commands for endpoint management"
```

---

## Task 14: Update CLI `scaffold.go` with Nexus templates

**Executor:** Use `go-architect` agent (subagent_type=go-architect) for Go implementation.

**Files:**

- Modify: `tools/timelord-cli/cmd/scaffold.go`
- Create: `templates/nexus-service.go.tmpl`
- Create: `templates/nexus-caller.go.tmpl`

**Step 1: Add `NexusTemplateData` struct and template constants**

Extend `TemplateData` with `Endpoint` and `ServiceName` fields. Add `nexusServiceTemplate` and `nexusCallerTemplate` constants following the pattern of `workflowTemplate`, `activityTemplate`, and `workerTemplate`.

The `nexusServiceTemplate` should generate a Go file with:

- Input/output types
- A sync operation (`nexus.NewSyncOperation`)
- An async operation (`temporalnexus.NewWorkflowRunOperation`)
- A backing workflow function
- A `NewXxxService()` constructor that creates the service and registers operations

The `nexusCallerTemplate` should generate a Go file with:

- Input/output types
- A caller workflow function using `workflow.NewNexusClient()` and `ExecuteOperation()`
- Proper error handling
- `scheduleToCloseTimeout` configured

**Step 2: Add cobra commands**

```go
var scaffoldNexusServiceCmd = &cobra.Command{
    Use:   "nexus-service [name]",
    Short: "Scaffold a Nexus service handler",
    Args:  cobra.ExactArgs(1),
    RunE:  runScaffoldNexusService,
}

var scaffoldNexusCallerCmd = &cobra.Command{
    Use:   "nexus-caller [name]",
    Short: "Scaffold a Nexus caller workflow",
    Args:  cobra.ExactArgs(1),
    RunE:  runScaffoldNexusCaller,
}
```

Add flags: `--endpoint` and `--service-name` for the caller template.

Register in `init()`:

```go
scaffoldCmd.AddCommand(scaffoldNexusServiceCmd)
scaffoldCmd.AddCommand(scaffoldNexusCallerCmd)
```

Update `executeTemplate()` to include `"nexus-service"` and `"nexus-caller"` in the templates map.

**Step 3: Create `.go.tmpl` files** (duplicating the template content for external template loading)

**Step 4: Build and verify**

```bash
cd tools/timelord-cli && go build -o bin/timelord . && ./bin/timelord scaffold --help
```

Verify `nexus-service` and `nexus-caller` appear in the scaffold subcommands.

**Step 5: Commit**

```bash
git add tools/timelord-cli/cmd/scaffold.go templates/nexus-service.go.tmpl templates/nexus-caller.go.tmpl
git commit -m "feat: add Nexus service and caller scaffolding templates"
```

---

## Task 15: Update top-level metadata files

**Executor:** Direct file edits — no specialized agent needed.

**Files:**

- Modify: `README.md`
- Modify: `CHANGELOG.md`
- Modify: `.claude-plugin/plugin.json`

**Step 1: Update README.md**

- Change "14 specialized skills" to "16 specialized skills" (line 9)
- Add `nexus-operations` and `nexus-decision-guide` to the Development skills list

**Step 2: Update CHANGELOG.md**

Add a new version entry at the top:

```markdown
## [1.1.0] - 2026-02-07

### Added

#### Nexus Support (GA)

- **nexus-operations** skill: Complete Nexus implementation guidance with multi-SDK examples (Go, TypeScript, Python, Java)
- **nexus-decision-guide** skill: Architecture decision framework for evaluating Nexus vs alternatives
- Nexus endpoint management in namespace-management skill
- Nexus cross-namespace pattern in workflow-patterns skill
- Nexus testing patterns in testing-strategies skill
- Nexus troubleshooting in troubleshooting skill
- Nexus development guidance in temporal-dev agent
- Nexus endpoint management in temporal-ops agent
- Nexus debugging in temporal-debug agent
- `/tl-scaffold nexus-service` — Scaffold Nexus handler services
- `/tl-scaffold nexus-caller` — Scaffold Nexus caller workflows
- `/tl-status nexus` — Check Nexus endpoint status
- `timelord nexus list|create|describe|delete` CLI commands
- Nexus service and caller Go templates
- How-to guide: Set up Nexus endpoints
- How-to guide: Call across namespaces with Nexus
- Reference: Nexus CLI commands and API surface
- Explanation: Why Nexus for cross-namespace communication
```

**Step 3: Update plugin.json**

- Add `"nexus"` to `keywords` array
- Update `version` to `"1.1.0"`

**Step 4: Commit**

```bash
git add README.md CHANGELOG.md .claude-plugin/plugin.json
git commit -m "chore: bump version to 1.1.0 with Nexus GA support"
```

---

## Task 16: Add Nexus documentation to Diataxis model

**Executor:** Use the `diataxis-orchestrator` agent to coordinate, with `doc-howto-writer`, `doc-reference-gen`, `doc-explanation-writer`, and `doc-crosslink-validator` agents for each document type.

Nexus needs proper documentation integrated with the existing Diataxis structure at `docs/content/`. The existing docs follow Hugo-style frontmatter (`title`, `weight`) and are organized into `tutorials/`, `how-to/`, `reference/`, and `explanation/`.

**Files:**

- Create: `docs/content/how-to/setup-nexus-endpoints.md`
- Create: `docs/content/how-to/call-across-namespaces.md`
- Create: `docs/content/reference/nexus-reference.md`
- Create: `docs/content/explanation/why-nexus.md`
- Modify: `docs/content/_index.md`
- Modify: `docs/content/reference/cli-reference.md`

**Step 1: Create how-to guide — Set Up Nexus Endpoints** (use `doc-howto-writer` agent)

`docs/content/how-to/setup-nexus-endpoints.md`

Goal-oriented recipe for creating Nexus endpoints between two namespaces. Follow the existing how-to pattern (see `handle-workflow-failures.md`: frontmatter, Problem, Solution, Prerequisites, Steps, Verification).

```markdown
---
title: "Set Up Nexus Endpoints"
weight: 5
---

# Set Up Nexus Endpoints

Connect two namespaces for cross-namespace workflow communication using Temporal Nexus.

## Problem

You have workflows in separate namespaces that need to call each other with durable execution guarantees.

## Solution

Create Nexus endpoints to route calls from a caller namespace to a handler namespace's task queue.

## Prerequisites

- Temporal CLI installed
- Both caller and handler namespaces exist
- Handler workers running on the target task queue

## Steps

### 1. Create the Nexus Endpoint
...CLI commands for create, with explanation of each flag...

### 2. Register a Nexus Service on the Handler Worker
...Go code snippet for service registration...

### 3. Create a Caller Workflow Using NexusClient
...Go code snippet for caller workflow...

### 4. Verify the Endpoint
...CLI commands for describe, list...

## Verification

...How to confirm end-to-end connectivity...
```

**Step 2: Create how-to guide — Call Across Namespaces** (use `doc-howto-writer` agent)

`docs/content/how-to/call-across-namespaces.md`

Focused on the caller-side patterns: creating NexusClient, executing operations, handling errors, setting timeouts.

```markdown
---
title: "Call Across Namespaces with Nexus"
weight: 6
---
```

Cover: NexusClient creation, sync vs async operations, error handling for `NexusOperationFailure`, timeout configuration.

**Step 3: Create reference — Nexus Reference** (use `doc-reference-gen` agent)

`docs/content/reference/nexus-reference.md`

Pure technical specification — no advice (per Diataxis reference rules). Include:

- Nexus endpoint CLI commands (create, list, describe, update, delete) with all flags
- Nexus event types table (`NexusOperationScheduled`, `Started`, `Completed`, `Failed`, `Canceled`, `TimedOut`)
- Nexus error types (`OperationError`, `HandlerError`, `NexusOperationFailure`)
- Nexus Go SDK API surface (`nexus.NewSyncOperation`, `temporalnexus.NewWorkflowRunOperation`, `workflow.NewNexusClient`, `NexusOperationOptions`)
- SDK support status table (Go GA, Java GA, Python preview, TS/NET experimental)

```markdown
---
title: "Nexus Reference"
weight: 4
---
```

**Step 4: Update CLI Reference** (use `doc-reference-gen` agent)

Add `timelord nexus` subcommands to `docs/content/reference/cli-reference.md`:

- `timelord nexus list` — List all Nexus endpoints
- `timelord nexus create <name>` — Create endpoint (flags: `--target-namespace`, `--target-task-queue`)
- `timelord nexus describe <name>` — Show endpoint details
- `timelord nexus delete <name>` — Delete endpoint

**Step 5: Create explanation — Why Nexus** (use `doc-explanation-writer` agent)

`docs/content/explanation/why-nexus.md`

Understanding-oriented document explaining *why* Nexus exists and *when* it's the right choice. Not a how-to — explains the architectural motivation:

- The namespace isolation problem
- Why activities calling Temporal client across namespaces is fragile
- How Nexus solves this with first-class durable cross-namespace communication
- Nexus vs child workflows vs activities vs signals — conceptual comparison
- The service contract model and why it matters for team independence

```markdown
---
title: "Why Nexus for Cross-Namespace Communication"
weight: 3
---
```

**Step 6: Update documentation index**

Add Nexus entries to `docs/content/_index.md`:

- Under "How-To Guides": add "Set up Nexus endpoints" and "Call across namespaces with Nexus"
- Under "Reference": add "Nexus Reference"
- Under "Explanation": add "Why Nexus"
- Under "Components > Agents" table: note Nexus in temporal-dev/ops/debug descriptions
- Under "Components > Skills > Development": add `nexus-operations` and `nexus-decision-guide`

**Step 7: Cross-link validation** (use `doc-crosslink-validator` agent)

Run the crosslink validator to ensure:

- All new docs have correct frontmatter
- Cross-references between how-to, reference, and explanation are valid
- New docs are properly linked from the index
- No orphaned pages

**Step 8: Commit**

```bash
git add docs/content/
git commit -m "docs: add Nexus documentation across Diataxis model (how-to, reference, explanation)"
```

---

## Verification

**Executor:** Use `superpowers:verification-before-completion` skill.

After all tasks, verify the complete integration:

1. **CLI builds:** `cd tools/timelord-cli && go build -o bin/timelord .`
2. **Scaffold works:** `./bin/timelord scaffold nexus-service TestService --output /tmp/nexus-test && cat /tmp/nexus-test/test_service_nexus_service.go`
3. **Nexus commands exist:** `./bin/timelord nexus --help`
4. **Skills are discoverable:** Check that the new skills' `description` fields contain all relevant trigger phrases and negative triggers
5. **Skill size check:** Verify `skills/nexus-operations/SKILL.md` is under 5,000 words (`wc -w`)
6. **Progressive disclosure:** Verify `skills/nexus-operations/references/nexus-multi-sdk.md` exists with TS/Python/Java examples
7. **Agent examples trigger correctly:** Verify the new `<example>` blocks in agents would match user queries about Nexus
8. **No broken cross-references:** Grep for references to skills/agents and verify they point to correct names
9. **Plugin metadata:** Verify `plugin.json` version is `1.1.0` and keywords include `"nexus"`
10. **No brain file modifications:** Confirm no changes to `.claude/brain/` — these auto-update via brain skills
11. **Diataxis docs complete:** Verify all 4 new doc files exist under `docs/content/`, have correct frontmatter, and are linked from `_index.md`
12. **Diataxis crosslinks valid:** Run `doc-crosslink-validator` agent to confirm no orphaned or broken references
