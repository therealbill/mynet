---
title: "How Temporal Works"
weight: 1
---

# How Temporal Works

Understanding Temporal's architecture helps you build better workflows and troubleshoot issues effectively.

## The Core Idea

Temporal solves a fundamental problem: **how do you build long-running, reliable applications?**

Traditional approaches fail because:

- Processes crash and lose state
- Services go down during operations
- Network calls fail mid-transaction
- Deployments interrupt running work

Temporal's solution: **persist every decision the code makes, so it can resume from exactly where it left off.**

## Event Sourcing Foundation

At its heart, Temporal uses event sourcing. Instead of storing current state, it stores the sequence of events that produced that state.

### Traditional State Storage

```
┌─────────────────────────┐
│  Order #123             │
│  Status: Processing     │
│  Items: [A, B, C]       │
│  Total: $150            │
└─────────────────────────┘
```

If the system crashes, you know the current state but not how you got there.

### Temporal's Event History

```
Event 1: WorkflowStarted {orderId: 123, items: [A, B, C]}
Event 2: ActivityScheduled {activityType: ValidateOrder}
Event 3: ActivityCompleted {result: valid}
Event 4: ActivityScheduled {activityType: ChargeCard}
Event 5: ActivityStarted
... (system crashes here)
```

On restart, Temporal replays these events to reconstruct exactly where the workflow was.

## The Workflow Execution Model

### Workflow Tasks

Temporal executes workflows in incremental steps called "workflow tasks":

```
┌───────────────────────────────────────────────────────┐
│                    Temporal Server                     │
│  ┌─────────────────────────────────────────────────┐  │
│  │               Workflow History                   │  │
│  │  Event 1 → Event 2 → Event 3 → Event 4 → ...   │  │
│  └─────────────────────────────────────────────────┘  │
│                         │                              │
│              Workflow Task                             │
│                         ▼                              │
└───────────────────────────────────────────────────────┘
                          │
                          ▼
┌───────────────────────────────────────────────────────┐
│                      Worker                            │
│  1. Receive workflow task                              │
│  2. Replay history to reconstruct state               │
│  3. Execute new code until next async point           │
│  4. Return decisions (schedule activity, start timer) │
└───────────────────────────────────────────────────────┘
```

### The Replay Process

When a worker picks up a workflow task:

1. **Load history**: Get all events from the server
2. **Replay**: Execute workflow code, skipping completed operations
3. **Execute new code**: Continue from where history ends
4. **Return decisions**: Tell server what to do next

This is why **determinism matters**: the same code with the same history must produce the same results.

## Components Deep Dive

### Temporal Server

The server manages:

- **History Service**: Stores workflow event histories
- **Matching Service**: Matches tasks with workers
- **Frontend Service**: API gateway for clients and workers

```
                     ┌─────────────────┐
                     │    Frontend     │
                     │    Service      │
                     └────────┬────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
       ┌──────────┐    ┌──────────┐    ┌──────────┐
       │ History  │    │ Matching │    │ Worker   │
       │ Service  │    │ Service  │    │ Service  │
       └──────────┘    └──────────┘    └──────────┘
              │               │
              ▼               ▼
       ┌──────────┐    ┌──────────┐
       │ Database │    │  Task    │
       │(History) │    │  Queues  │
       └──────────┘    └──────────┘
```

### Workers

Workers are your applications that:

- **Poll task queues** for work
- **Execute workflow code** (replay + new execution)
- **Execute activity code** (actual work)
- **Report results** back to the server

Workers are stateless—any worker can pick up any task for workflows and activities it has registered.

### Task Queues

Task queues decouple workflow scheduling from execution:

```
┌─────────────────────────────────────────────────────┐
│                   Task Queue: "orders"              │
├─────────────────────────────────────────────────────┤
│  Workflow Task: order-123 (waiting)                 │
│  Activity Task: validate-order (waiting)            │
│  Activity Task: charge-card (in progress)           │
└─────────────────────────────────────────────────────┘
         │                    │                    │
         ▼                    ▼                    ▼
    ┌─────────┐         ┌─────────┐         ┌─────────┐
    │Worker 1 │         │Worker 2 │         │Worker 3 │
    └─────────┘         └─────────┘         └─────────┘
```

Benefits:

- **Scalability**: Add workers to handle more load
- **Isolation**: Different task queues for different concerns
- **Routing**: Direct specific work to specific workers

## Durability Guarantees

### Workflow Durability

Workflows are durable because:

1. **Every decision is persisted** before it takes effect
2. **Events are immutable** once written
3. **Replay reconstructs state** exactly

Even if:

- Workers crash → Another worker picks up the task
- Server crashes → Resumes from persisted history
- Network fails → Retries with exactly-once semantics

### Activity Execution

Activities have different guarantees:

- **At-least-once execution**: May run multiple times on failure
- **Heartbeats**: Long activities report progress
- **Retries**: Configurable retry policies

This is why activities must be **idempotent** or handle duplicate execution.

## History Shards

For scalability, Temporal partitions workflows across shards:

```
┌─────────────────────────────────────────────────────┐
│                  History Service                     │
├─────────────────────────────────────────────────────┤
│  Shard 1    │  Shard 2    │  Shard 3    │  ...     │
│  ├ wf-001   │  ├ wf-002   │  ├ wf-003   │          │
│  ├ wf-004   │  ├ wf-005   │  ├ wf-006   │          │
│  └ wf-007   │  └ wf-008   │  └ wf-009   │          │
└─────────────────────────────────────────────────────┘
```

Key points:

- **Shard count is fixed** at cluster creation
- **Workflows hash to shards** by workflow ID
- **More shards = more parallelism** but more overhead

This is why [cluster sizing](/reference/skill-reference/#cluster-sizing) is important.

## Putting It Together

A typical workflow execution:

```
1. Client: StartWorkflow(OrderWorkflow, order)
   └─→ Server: Create history, schedule workflow task

2. Worker: Poll workflow task queue
   └─→ Server: Return task with history (empty)

3. Worker: Execute OrderWorkflow until first activity
   └─→ Server: Return command: ScheduleActivity(ValidateOrder)

4. Server: Record ActivityScheduled event
   └─→ Queue activity task

5. Worker: Poll activity task queue
   └─→ Server: Return ValidateOrder task

6. Worker: Execute ValidateOrder activity
   └─→ Server: Record ActivityCompleted, schedule workflow task

7. Worker: Poll workflow task queue
   └─→ Server: Return task with history (4 events)

8. Worker: Replay history (skip ValidateOrder), continue execution
   └─→ Server: Return command: ScheduleActivity(ChargeCard)

... continues until workflow completes
```

## Why This Architecture?

### Versus Message Queues

Message queues (RabbitMQ, SQS) handle individual messages. Temporal handles entire workflows:

- **State management**: Temporal tracks workflow progress
- **Retries**: Built-in, with history preservation
- **Visibility**: Query and search running workflows
- **Time**: Timers and sleep built-in

### Versus Databases + Cron

Rolling your own with databases and cron:

- **State machines are complex**: Temporal handles this
- **Failure handling**: Automatic retries and recovery
- **Scalability**: Distributed, partitioned architecture
- **Testing**: Replay testing catches issues early

### Versus Serverless Orchestration

AWS Step Functions, Azure Durable Functions:

- **Code-first**: Write workflows in your language
- **Testability**: Unit test with TestWorkflowEnvironment
- **Portability**: Self-hosted, cloud-agnostic
- **Flexibility**: Full language features available

## Summary

Temporal's architecture provides:

1. **Durability** through event sourcing
2. **Scalability** through sharding and task queues
3. **Reliability** through deterministic replay
4. **Simplicity** through code-first workflows

Understanding these concepts helps you:

- Design workflows that leverage Temporal's strengths
- Troubleshoot issues by understanding the execution model
- Make informed decisions about configuration and scaling
