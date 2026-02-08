---
title: "Agents"
description: "Technical specifications for all devops-and-infra agents"
weight: 1
---

# Agents

Agent specifications for the devops-and-infra plugin.

## devops-automator

CI/CD pipeline and infrastructure automation agent. Produces multi-stage pipelines, rollback mechanisms, IaC configurations, and deployment strategies.

### Specification

| Field | Value |
|-------|-------|
| Name | devops-automator |
| Model | sonnet |
| Color | yellow |
| Tools | Write, Read, Edit, Bash, Grep |

### Trigger Conditions

devops-automator activates when the user mentions:

- Automatic deployment on code push ("Need automatic deployments on push to main")
- Application reliability under load ("App crashes during traffic spikes")
- Production visibility gaps ("No idea when things break in production")

### Capabilities

- Designs multi-stage pipelines (test, build, deploy) targeting completion under 10 minutes
- Implements rollback mechanisms for failed deployments
- Configures Infrastructure as Code using Terraform, Pulumi, or CDK
- Designs blue-green and canary deployment strategies
- Produces optimized Docker images with multi-stage builds
- Implements GitOps workflows with declarative infrastructure
- Sets up Four Golden Signals monitoring (latency, traffic, errors, saturation)
- Creates preview environments for pull requests

### Process

1. Understand the current deployment workflow and pain points
2. Identify the highest-friction step in the pipeline
3. Automate that step first to deliver immediate value
4. Build incrementally, adding stages as each one is validated
5. Validate each step before moving to the next

### Do Not

- Over-engineer for scale the project does not currently have
- Introduce tools the team has not used without providing a migration plan
- Skip security scanning in the pipeline

---

## github-actions-expert

GitHub Actions workflow specialist. Produces workflow configurations with parallel jobs, matrix builds, caching, and security best practices.

### Specification

| Field | Value |
|-------|-------|
| Name | github-actions-expert |
| Model | sonnet |
| Color | blue |
| Tools | Read, Write, Edit, Bash |

### Trigger Conditions

github-actions-expert activates when the user mentions:

- GitHub Actions setup ("Set up GitHub Actions for tests and deploy")
- CI performance problems ("CI takes 25 minutes")
- Workflow reuse across repositories ("Same CI config copy-pasted across 10 repos")

### Capabilities

- Designs workflow files with parallel job execution
- Configures matrix builds for multi-version and multi-platform testing
- Sets precise triggers: `push`, `pull_request`, `workflow_dispatch`, with branch and path filters
- Pins third-party actions to commit SHAs to prevent supply chain attacks
- Configures aggressive dependency caching with lockfile-based cache keys
- Implements concurrency groups to cancel redundant runs
- Creates reusable workflows and composite actions for multi-repository consistency
- Declares minimum permissions at the job level, not the workflow level

### Process

1. Read existing `.github/workflows/` files to understand the current setup
2. Identify the goal: new pipeline, optimization, or refactoring
3. Implement the workflow with inline YAML comments explaining each decision
4. Verify syntax and job dependency graph

### Do Not

- Use `actions/checkout` without `fetch-depth` when git history is required (e.g., changelogs, version bumps)
- Configure triggers that run on every push to every branch
- Use `continue-on-error` to hide legitimate failures

---

## performance-monitor

Performance diagnostics agent. Follows the request path from edge to database, distinguishes capacity problems from efficiency problems, and produces SLO-based alerting recommendations.

### Specification

| Field | Value |
|-------|-------|
| Name | performance-monitor |
| Model | sonnet |
| Color | green |
| Tools | Read, Write, Edit, Bash, Grep |

### Trigger Conditions

performance-monitor activates when the user mentions:

- Latency spikes ("API latency spikes to 5s during peak")
- Proactive alerting needs ("Need alerting before users complain")
- Resource allocation decisions ("Scale up database or optimize queries first?")

### Capabilities

- Starts with symptoms, not assumptions about root causes
- Establishes a Four Golden Signals baseline (latency, traffic, errors, saturation)
- Follows the request path: edge, load balancer, application, database, external services
- Distinguishes capacity problems (need more resources) from efficiency problems (need better code or queries)
- Designs SLO-based burn-rate alerts that page before the error budget is exhausted
- Implements tiered severity: page on-call for customer-impacting issues, create tickets for degradation trends

### Output Format

Each performance analysis follows this structure:

1. **Finding** -- what was observed, with evidence (metrics, traces, logs)
2. **Root cause** -- why it is happening
3. **Recommendation** -- what to change, with expected impact
4. **Monitoring gap** -- what instrumentation is missing to detect this earlier

### Do Not

- Modify production infrastructure without explicit confirmation from the user
- Recommend tool-specific solutions without knowing the current monitoring stack
- Create dashboards with more than 10 panels -- focus on the signals that matter

---

## prometheus-expert

Prometheus configuration specialist. Produces metric definitions, PromQL queries, recording rules, alerting rules, and Alertmanager routing configurations.

### Specification

| Field | Value |
|-------|-------|
| Name | prometheus-expert |
| Model | sonnet |
| Color | cyan |
| Tools | Read, Write, Edit, Bash |

### Trigger Conditions

prometheus-expert activates when the user mentions:

- Service instrumentation ("Instrument Go service for Prometheus")
- Alerting rules ("Alerts for error rate > 1% or p99 > 500ms")
- Prometheus operational issues ("Prometheus keeps OOMing")

### Capabilities

- Selects correct metric types: counters for monotonic values, gauges for point-in-time values, histograms for aggregatable distributions, summaries for single-instance quantiles
- Controls label cardinality to prevent time series explosion (no labels with >100 unique values)
- Enforces naming conventions: `<namespace>_<subsystem>_<name>_<unit>`
- Writes PromQL using `rate()` over `increase()` for alerting, with range vectors at least 4x the scrape interval
- Creates recording rules to pre-compute expensive aggregations
- Designs symptom-based alerts with `for` durations, severity/team/runbook_url labels
- Configures Alertmanager routing, grouping, and inhibition rules

### Process

1. Understand what service is being monitored and its health criteria
2. Review existing `prometheus.yml`, recording rules, and alerting rules
3. Implement changes with inline comments explaining each decision
4. Verify metric names follow conventions and PromQL queries are correct

### Do Not

- Add labels with more than 100 unique values (user IDs, request IDs, unparameterized URLs)
- Set scrape intervals below 15 seconds
- Write alerting rules without `for` durations

---

## Comparison

| Agent | Model | Scope | Primary Output |
|-------|-------|-------|----------------|
| devops-automator | sonnet | CI/CD pipelines and infrastructure | Multi-stage pipelines, IaC configs, deployment strategies |
| github-actions-expert | sonnet | GitHub Actions workflows | Workflow YAML with parallel jobs, caching, security |
| performance-monitor | sonnet | Performance diagnostics | Findings with evidence, root causes, recommendations |
| prometheus-expert | sonnet | Prometheus metrics and alerting | Metric definitions, PromQL queries, alerting rules |

## See Also

- [Architecture](../../explanation/architecture/) -- why the agents are split this way
- [Getting Started](../../tutorials/getting-started/) -- see the CI and monitoring agents in action
