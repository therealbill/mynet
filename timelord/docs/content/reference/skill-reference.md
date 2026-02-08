---
title: "Skill Reference"
weight: 2
---

# Skill Reference

Catalog of all Timelord skills with descriptions and usage guidance.

## Overview

Skills are specialized knowledge modules that Claude automatically activates based on context. Each skill provides focused expertise for specific Temporal.io topics.

## Development Skills

Skills focused on building Temporal applications.

### workflow-patterns

**Triggers:** "workflow pattern", "saga pattern", "state machine workflow", "Continue-As-New", "deterministic workflow", "workflow design"

Guidance for designing workflow implementations including:

- **Determinism Rules** - What operations are safe in workflows
- **Saga Pattern** - Compensating transactions for distributed workflows
- **State Machine** - Event-driven workflow design
- **Continue-As-New** - Handling long-running workflows
- **Child Workflows** - Composition and modularity

**Key Topics:**

| Pattern | Use Case |
|---------|----------|
| Sequential | Simple ordered steps |
| Parallel | Independent concurrent activities |
| Saga | Multi-step with compensation |
| State Machine | Event-driven transitions |
| Long-running | Workflows spanning days/weeks |

### activity-design

**Triggers:** "activity timeout", "retry policy", "heartbeat", "idempotent activity", "activity options", "activity design"

Best practices for implementing activities:

- **Timeout Configuration** - StartToClose, ScheduleToClose, Heartbeat
- **Retry Policies** - Backoff, max attempts, non-retryable errors
- **Heartbeats** - Progress reporting for long activities
- **Idempotency** - Safe retry handling

**Timeout Quick Reference:**

| Timeout | Purpose | Typical Value |
|---------|---------|---------------|
| StartToCloseTimeout | Max execution time | 1-30 minutes |
| ScheduleToCloseTimeout | Total allowed time | 1-24 hours |
| HeartbeatTimeout | Liveness check | 30-60 seconds |

### testing-strategies

**Triggers:** "test workflow", "TestWorkflowEnvironment", "mock activity", "replay test", "workflow testing"

Comprehensive testing approaches:

- **Unit Testing** - TestWorkflowEnvironment usage
- **Mocking** - Activity and child workflow mocks
- **Replay Testing** - Verify determinism with history
- **Integration Testing** - End-to-end with real cluster

**Test Environment Setup:**

```go
func TestWorkflow(t *testing.T) {
    testSuite := &testsuite.WorkflowTestSuite{}
    env := testSuite.NewTestWorkflowEnvironment()

    // Mock activities
    env.OnActivity(MyActivity, mock.Anything).Return(result, nil)

    // Execute workflow
    env.ExecuteWorkflow(MyWorkflow, input)

    // Assert results
    require.True(t, env.IsWorkflowCompleted())
    require.NoError(t, env.GetWorkflowError())
}
```

### signals-queries-updates

**Triggers:** "workflow signal", "workflow query", "workflow update", "send signal", "query workflow state"

Message passing patterns:

- **Signals** - Async messages to workflows
- **Queries** - Read-only state inspection
- **Updates** - Sync request-response mutations

### versioning-guide

**Triggers:** "workflow versioning", "GetVersion", "workflow migration", "update running workflow", "workflow reset", "backward compatible workflow"

Safe workflow evolution:

- **GetVersion API** - Branch code for compatibility
- **Version Numbers** - DefaultVersion and increments
- **Migration Strategies** - Rolling updates
- **Workflow Reset** - Recovery from stuck states

**GetVersion Pattern:**

```go
v := workflow.GetVersion(ctx, "change-id", workflow.DefaultVersion, 1)
if v == workflow.DefaultVersion {
    // Old logic
} else {
    // New logic (v >= 1)
}
```

### nexus-operations

**Triggers:** "nexus operation", "cross-namespace", "nexus service", "nexus handler", "nexus caller", "sync operation", "async operation"

Comprehensive Nexus implementation patterns:

- **Sync Operations** - Quick request-response calls (< 10s)
- **Async Operations** - Workflow-backed long-running operations
- **Service Registration** - Handler worker setup with Nexus services
- **Caller Patterns** - Invoking operations from caller workflows
- **Error Handling** - NexusOperationFailure and retry strategies

**Operation Type Quick Reference:**

| Type | Duration | Handler Pattern | Use Case |
|------|----------|----------------|----------|
| Sync | < 10 seconds | `nexus.NewSyncOperation` | Quick lookups, validations |
| Async | Unlimited | `temporalnexus.NewWorkflowRunOperation` | Long-running processes |

### nexus-decision-guide

**Triggers:** "should I use nexus", "nexus vs child workflow", "cross-namespace architecture", "nexus evaluation", "namespace communication"

Architecture evaluation for Nexus adoption:

- **Decision Criteria** - When Nexus is the right choice
- **Alternatives Comparison** - Nexus vs activities, child workflows, signals
- **Migration Path** - Moving from activity-based cross-namespace calls to Nexus
- **Operational Overhead** - What Nexus requires to run effectively

**Decision Quick Reference:**

| Scenario | Recommendation |
|----------|---------------|
| Same namespace | Child workflows |
| Cross-namespace, short-lived | Nexus sync operation |
| Cross-namespace, long-running | Nexus async operation |
| Fire-and-forget notification | Signals |
| External HTTP API call | Regular activity |

## Operations Skills

Skills focused on deploying and managing Temporal clusters.

### cluster-deployment

**Triggers:** "deploy temporal", "install temporal", "helm temporal", "kubernetes temporal", "temporal cluster setup"

Deployment guidance:

- **Local Development** - Docker Compose, minikube
- **Kubernetes** - Helm chart configuration
- **Production** - HA setup, resource allocation
- **Cloud Providers** - EKS, GKE, AKS specifics

**Quick Start:**

```bash
# Add Helm repo
helm repo add temporal https://charts.temporal.io

# Install with PostgreSQL
helm install temporal temporal/temporal \
  --set server.replicaCount=3 \
  --set cassandra.enabled=false \
  --set postgresql.enabled=true
```

### cluster-sizing

**Triggers:** "history shards", "temporal sizing", "cluster capacity", "shard configuration", "temporal resources"

Capacity planning:

- **History Shards** - Critical one-time decision
- **Resource Allocation** - CPU, memory per component
- **Scaling Guidelines** - When and how to scale
- **Database Sizing** - PostgreSQL/Cassandra requirements

**Shard Guidelines:**

| Concurrent Workflows | Recommended Shards |
|---------------------|-------------------|
| < 10,000 | 512 |
| 10,000 - 100,000 | 1024 |
| 100,000 - 500,000 | 2048 |
| > 500,000 | 4096 |

### monitoring-setup

**Triggers:** "temporal metrics", "prometheus temporal", "grafana temporal", "temporal alerts", "monitor temporal"

Observability configuration:

- **Prometheus Metrics** - Key metrics to collect
- **Grafana Dashboards** - Visualization setup
- **Alerting Rules** - Critical alert definitions
- **SLO Configuration** - Service level objectives

**Critical Metrics:**

| Metric | Alert Threshold |
|--------|-----------------|
| `schedule_to_start_latency_p99` | > 5s |
| `workflow_task_failed_total` | > 10/min |
| `persistence_latency_p99` | > 100ms |

### security-config

**Triggers:** "temporal mTLS", "temporal authorization", "temporal security", "secure temporal", "TLS temporal"

Security implementation:

- **mTLS Setup** - Certificate generation and configuration
- **Authorization** - Role-based access control
- **Namespace Isolation** - Multi-tenant security
- **Network Policies** - Kubernetes network security

### namespace-management

**Triggers:** "temporal namespace", "create namespace", "namespace retention", "multi-tenant temporal"

Namespace operations:

- **Creation** - New namespace setup
- **Configuration** - Retention, search attributes
- **Multi-tenancy** - Isolation strategies
- **Migration** - Moving workflows between namespaces

### visibility-search

**Triggers:** "search attributes", "temporal search", "list workflows", "visibility query", "elasticsearch temporal"

Workflow discovery:

- **Search Attributes** - Custom indexed fields
- **Query Syntax** - Filter expressions
- **Elasticsearch Setup** - Advanced visibility
- **Performance** - Query optimization

## Debugging Skills

Skills focused on troubleshooting and diagnosis.

### troubleshooting

**Triggers:** "workflow stuck", "workflow failing", "temporal error", "debug workflow", "diagnose temporal", "workflow not completing", "activity timeout", "non-deterministic error"

Issue resolution:

- **Diagnostic Approach** - Systematic troubleshooting
- **Common Issues** - Stuck, failed, slow workflows
- **Error Catalog** - Known errors and fixes
- **Recovery Actions** - Reset, terminate, cancel

**Diagnosis Tree:**

```
Issue Detected
├── Workflow Stuck?
│   ├── Check pending activities
│   ├── Check task queue workers
│   └── Check pending timers/signals
├── Workflow Failed?
│   ├── Check activity errors
│   ├── Check non-determinism
│   └── Check panic/crash
└── Performance Issues?
    ├── Check schedule-to-start latency
    └── Check activity execution time
```

## Skill Usage

Skills activate automatically based on conversation context. To explicitly request a skill:

```
Tell me about workflow patterns for saga transactions
```

Or reference directly:

```
Using the testing-strategies skill, help me write tests for my workflow
```

## Related

- [Agent Reference](/reference/agent-reference) - Specialized agents
- [CLI Reference](/reference/cli-reference) - Command-line tools
