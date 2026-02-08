# Timelord

> Expert guidance for Temporal.io workflow orchestration—from first workflow to production cluster.

A Claude Code plugin providing comprehensive expertise for deploying, managing, and developing with Temporal.io (self-hosted).

## Features

- **16 specialized skills** covering workflows, activities, testing, deployment, Nexus, and troubleshooting
- **3 expert agents** for development, operations, and debugging
- **CLI tool** for scaffolding projects and inspecting clusters
- **6 slash commands** for common tasks
- **Complete documentation** with tutorials, how-to guides, and conceptual explanations

## Overview

Timelord helps platform engineers and application developers work effectively with Temporal.io:

- **Workflow Development**: Design patterns, activity implementation, worker configuration
- **Cluster Operations**: Deployment, monitoring, scaling, security
- **Troubleshooting**: Diagnosis, event history analysis, common error resolution

## Target Users

| Role | Focus |
|------|-------|
| Application Developers | Workflow patterns, activities, testing, versioning |
| Platform Engineers | Cluster deployment, monitoring, scaling, security |

## Technology Stack

- **Languages**: Go (primary), Python (secondary)
- **Deployment**: Kubernetes (EKS production, local k8s development)
- **Database**: PostgreSQL

## Components

### Agents

- `temporal-dev` - Workflow development guidance
- `temporal-ops` - Cluster operations expertise
- `temporal-debug` - Troubleshooting specialist

### Commands

| Command | Description |
|---------|-------------|
| `/tl-scaffold` | Generate project, workflow, activity, or worker code |
| `/tl-status` | Check cluster and workflow status |
| `/tl-deploy` | Deployment guidance with Kubernetes context |
| `/tl-diagnose` | Workflow troubleshooting |
| `/tl-test` | Run workflow tests |
| `/tl-upgrade` | Upgrade planning |

### Skills

**Development**:

- `workflow-patterns` - Determinism, saga, state machine, Continue-As-New
- `activity-design` - Idempotency, retries, timeouts, heartbeats
- `testing-strategies` - TestWorkflowEnvironment, mocking, replay
- `worker-tuning` - Concurrency, task queues, resource management
- `versioning-guide` - GetVersion, workflow reset, safe deployments
- `signals-queries-updates` - Message passing patterns
- `nexus-operations` - Cross-namespace durable communication with Nexus
- `nexus-decision-guide` - When to use Nexus vs alternatives

**Operations**:

- `cluster-deployment` - Helm, local k8s, production configs
- `cluster-sizing` - Shards, resources, capacity planning
- `monitoring-setup` - Prometheus, Grafana, alerting
- `security-config` - mTLS, authorization
- `namespace-management` - Multi-tenancy patterns
- `visibility-search` - Search attributes, Elasticsearch

**Troubleshooting**:

- `troubleshooting` - Common errors, diagnosis workflows
- `upgrade-guide` - Safe upgrade procedures

### CLI Tool

The `timelord-cli` provides scaffolding and cluster interaction:

```bash
# Scaffold a new project
timelord scaffold project myapp

# Generate a workflow
timelord scaffold workflow OrderProcessing

# Check cluster status
timelord cluster status

# Describe a workflow execution
timelord workflow describe <workflow-id>
```

## Installation

```bash
# From marketplace (when published)
claude plugin install timelord

# From local directory
claude plugin install ./timelord
```

## Quick Start

1. **Scaffold a project**:
   ```
   /tl-scaffold project myapp
   ```

2. **Get workflow guidance**:
   Ask about workflow patterns and the `temporal-dev` agent will help design your solution.

3. **Deploy to Kubernetes**:
   ```
   /tl-deploy
   ```

## Documentation

Full documentation available at `docs/` including:

- Tutorials for getting started
- How-to guides for common tasks
- Reference documentation
- Conceptual explanations

## License

MIT
