---
title: "Timelord Documentation"
type: docs
---

# Timelord

A Claude Code plugin for Temporal.io expertise.

Timelord helps platform engineers and application developers work effectively with Temporal.io, providing guidance for workflow development, cluster operations, and troubleshooting.

## Quick Start

### Installation

```bash
# Install from marketplace
claude plugin install timelord

# Or install from local directory
claude plugin install ./timelord
```

### Scaffold a Project

```
/tl-scaffold project myapp
```

### Get Workflow Guidance

Ask Claude about workflow patterns:

> "I need to implement an order processing workflow with payment and shipping"

The `temporal-dev` agent will help design your solution with proper patterns.

## Documentation

### Tutorials

Step-by-step guides for learning Temporal.io:

- [Your First Workflow](tutorials/first-workflow/) - Build a complete workflow from scratch

### How-To Guides

Task-oriented guides for common operations:

- Scale workers for high throughput
- Configure mTLS for production
- Set up monitoring with Prometheus

### Reference

Technical specifications:

- CLI Reference
- Skill Reference
- Agent Capabilities

### Explanation

Conceptual documentation:

- Temporal Architecture
- Workflow Determinism
- Activity Design Principles

## Components

### Agents

| Agent | Focus |
|-------|-------|
| `temporal-dev` | Workflow development, activities, testing |
| `temporal-ops` | Cluster operations, monitoring, scaling |
| `temporal-debug` | Troubleshooting, diagnosis, event history |

### Commands

| Command | Description |
|---------|-------------|
| `/tl-scaffold` | Generate project components |
| `/tl-status` | Check cluster and workflow status |
| `/tl-deploy` | Deployment guidance |
| `/tl-diagnose` | Workflow troubleshooting |
| `/tl-test` | Run workflow tests |
| `/tl-upgrade` | Upgrade planning |

### Skills

**Development:**

- `workflow-patterns` - Determinism, saga, state machine
- `activity-design` - Idempotency, retries, heartbeats
- `testing-strategies` - Unit tests, replay testing

**Operations:**

- `cluster-deployment` - Helm, Kubernetes setup
- `cluster-sizing` - Capacity planning
- `monitoring-setup` - Prometheus, Grafana

**Troubleshooting:**

- `troubleshooting` - Common errors, diagnosis
- `versioning-guide` - Safe workflow updates

## Technology Stack

- **Languages:** Go (primary), Python (secondary)
- **Deployment:** Kubernetes (EKS, local)
- **Database:** PostgreSQL
- **SDK:** Temporal Go SDK, Temporal Python SDK
