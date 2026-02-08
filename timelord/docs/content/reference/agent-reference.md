---
title: "Agent Reference"
weight: 3
---

# Agent Reference

Specialized agents for Temporal.io development, operations, and debugging.

## Overview

Timelord includes three specialized agents, each focused on a specific domain. Agents are autonomous subprocesses that handle complex, multi-step tasks independently.

## temporal-dev

**Focus:** Application development with Temporal

**Color:** Blue

**Best for:**

- Designing workflow architectures
- Implementing activities and workers
- Writing deterministic workflow code
- Testing strategies and implementation
- SDK usage and patterns

**Triggers:**

- "How do I design this workflow?"
- "Help me implement an activity"
- "What's the best pattern for..."
- "How do I test my workflow?"

**Capabilities:**

| Task | Description |
|------|-------------|
| Workflow Design | Architecture patterns, determinism guidance |
| Activity Implementation | Timeouts, retries, idempotency |
| Worker Configuration | Task queues, concurrency, tuning |
| Testing | Unit tests, mocking, replay testing |
| SDK Patterns | Go SDK best practices |

**Example Interactions:**

```
User: "I need to process orders with payment, inventory, and shipping steps"
Agent: Analyzes requirements, suggests saga pattern with compensation...

User: "How do I handle long-running file processing?"
Agent: Recommends heartbeats, progress tracking, timeout configuration...

User: "Help me write tests for my OrderWorkflow"
Agent: Generates TestWorkflowEnvironment setup, mocks, assertions...
```

**Tools Available:**

- Read, Write, Edit (code generation)
- Grep, Glob (codebase analysis)
- Bash (running tests, CLI commands)

## temporal-ops

**Focus:** Cluster deployment and operations

**Color:** Green

**Best for:**

- Deploying Temporal clusters
- Kubernetes/Helm configuration
- Monitoring and alerting setup
- Security configuration
- Capacity planning and scaling

**Triggers:**

- "Help me deploy Temporal"
- "Set up monitoring for my cluster"
- "Configure mTLS"
- "How many shards do I need?"

**Capabilities:**

| Task | Description |
|------|-------------|
| Deployment | Helm charts, Kubernetes manifests |
| Monitoring | Prometheus, Grafana, alerts |
| Security | mTLS, authorization, network policies |
| Scaling | Shards, resources, capacity planning |
| Upgrades | Version migration, rolling updates |

**Example Interactions:**

```
User: "Deploy Temporal to my EKS cluster"
Agent: Generates Helm values, guides through deployment steps...

User: "Set up alerts for my production cluster"
Agent: Creates Prometheus rules, Grafana dashboards...

User: "We're expecting 100k concurrent workflows"
Agent: Calculates shard count, resource requirements...
```

**Tools Available:**

- Read, Write, Edit (configuration files)
- Grep, Glob (existing configs)
- Bash (kubectl, helm commands)
- WebFetch, WebSearch (documentation)

## temporal-debug

**Focus:** Troubleshooting and diagnosis

**Color:** Red

**Best for:**

- Diagnosing stuck workflows
- Analyzing event history
- Resolving errors
- Performance investigation
- Non-determinism detection

**Triggers:**

- "My workflow is stuck"
- "Why is this workflow failing?"
- "I'm getting non-deterministic errors"
- "Analyze this workflow's history"

**Capabilities:**

| Task | Description |
|------|-------------|
| Workflow Diagnosis | Status analysis, issue detection |
| Event History | Parse and explain event sequences |
| Error Resolution | Root cause analysis, fix suggestions |
| Performance | Latency investigation, bottlenecks |
| Recovery | Reset, terminate, retry guidance |

**Example Interactions:**

```
User: "My order workflow has been running for 2 hours"
Agent: Checks status, finds pending activity, identifies worker issue...

User: "Getting 'non-deterministic workflow definition' error"
Agent: Analyzes code for determinism violations, suggests fixes...

User: "Why did this workflow fail?"
Agent: Parses history, identifies failed activity, explains error...
```

**Diagnostic Commands:**

```bash
# Agent commonly uses these
timelord workflow diagnose <id> --json
timelord workflow describe <id> --json
timelord workflow history <id> --json
temporal task-queue describe --task-queue <queue>
```

**Tools Available:**

- Read, Write, Edit (code fixes)
- Grep, Glob (log analysis)
- Bash (CLI diagnostics)
- WebFetch, WebSearch (error lookup)

## Agent Selection Guide

| Scenario | Recommended Agent |
|----------|-------------------|
| "Design a workflow for X" | temporal-dev |
| "Deploy Temporal to production" | temporal-ops |
| "My workflow is stuck/failing" | temporal-debug |
| "Write tests for my workflow" | temporal-dev |
| "Set up monitoring" | temporal-ops |
| "Analyze this event history" | temporal-debug |
| "Implement activity with retries" | temporal-dev |
| "Configure mTLS" | temporal-ops |
| "Fix non-determinism error" | temporal-debug |
| "Scale my cluster" | temporal-ops |

## Invoking Agents

Agents are invoked automatically based on context. To explicitly request:

```
Use the temporal-debug agent to diagnose my stuck workflow
```

Or through the Task tool:

```
Task: Diagnose workflow order-12345
Agent: temporal-debug
```

## Agent Collaboration

Complex scenarios may involve multiple agents:

1. **Development → Debug:** Build workflow, then troubleshoot issues
2. **Ops → Debug:** Deploy cluster, then diagnose problems
3. **Debug → Dev:** Identify issue, then fix code

Example flow:

```
1. temporal-dev: Design and implement workflow
2. temporal-ops: Deploy to staging
3. temporal-debug: Diagnose test failures
4. temporal-dev: Fix identified issues
5. temporal-ops: Deploy to production
```

## Related

- [Skill Reference](/reference/skill-reference) - Knowledge modules
- [CLI Reference](/reference/cli-reference) - Command-line tools
