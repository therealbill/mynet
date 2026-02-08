# Changelog

All notable changes to Timelord will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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

## [1.0.0] - 2024-01-15

### Added

#### Agents
- **temporal-dev**: Development agent for workflow design, activities, and testing
- **temporal-ops**: Operations agent for cluster deployment, monitoring, and security
- **temporal-debug**: Debugging agent for troubleshooting and event history analysis

#### Skills (14 total)

Development skills:
- **workflow-patterns**: Determinism, saga, state machine, Continue-As-New patterns
- **activity-design**: Timeouts, retries, heartbeats, idempotency
- **testing-strategies**: TestWorkflowEnvironment, mocking, replay testing
- **worker-tuning**: Concurrency, task queues, performance optimization
- **versioning-guide**: GetVersion API, workflow migration, reset commands
- **signals-queries-updates**: Signals, queries, and updates message passing

Operations skills:
- **cluster-deployment**: Helm deployment, local k8s, production configuration
- **cluster-sizing**: History shards, resource planning, capacity
- **monitoring-setup**: Prometheus metrics, Grafana dashboards, alerting
- **security-config**: mTLS, authorization, network policies
- **namespace-management**: Multi-tenancy, namespace organization
- **visibility-search**: Search attributes, query syntax, Elasticsearch

Troubleshooting skills:
- **troubleshooting**: Common errors, diagnosis workflows, recovery

#### Commands
- **/tl-scaffold**: Generate projects, workflows, activities, and workers
- **/tl-status**: Check cluster health and workflow status
- **/tl-deploy**: Deployment guidance with Kubernetes context
- **/tl-upgrade**: Upgrade planning and migration guidance
- **/tl-diagnose**: Workflow troubleshooting and analysis
- **/tl-test**: Workflow validation and replay testing

#### CLI Tool (timelord-cli)

Scaffold commands:
- `timelord scaffold project` - Full project structure
- `timelord scaffold workflow` - Workflow boilerplate
- `timelord scaffold activity` - Activity boilerplate
- `timelord scaffold worker` - Worker boilerplate

Cluster commands:
- `timelord cluster status` - Health check
- `timelord cluster info` - Configuration details
- `timelord cluster metrics` - Metrics summary

Namespace commands:
- `timelord namespace list` - List namespaces
- `timelord namespace create` - Create namespace
- `timelord namespace describe` - Namespace details

Workflow commands:
- `timelord workflow list` - List with filters
- `timelord workflow describe` - Execution details
- `timelord workflow history` - Event history
- `timelord workflow diagnose` - Issue detection

Test commands:
- `timelord test validate` - Code validation
- `timelord test replay` - Replay test generation

#### Documentation

Tutorials:
- First Workflow: Build your first Temporal workflow
- Deploy Local Cluster: Docker Compose and minikube setup
- Production Setup: Production deployment guide

How-to guides:
- Scale Workers: Horizontal and vertical scaling
- Configure mTLS: Certificate generation and setup
- Setup Monitoring: Prometheus and Grafana
- Handle Workflow Failures: Recovery strategies

Explanations:
- How Temporal Works: Architecture deep dive
- Why Determinism Matters: Understanding replay

Reference:
- CLI Reference: Complete command documentation
- Skill Reference: All skills with triggers
- Agent Reference: Agent capabilities

#### Templates
- Go workflow template
- Go activity template
- Go worker template
- Full project scaffold with Dockerfile and Makefile

### Technical Details
- Primary language support: Go
- Secondary language support: Python (patterns and examples)
- Deployment target: Kubernetes (EKS production, local development)
- Database: PostgreSQL
- CLI built with Cobra framework
- JSON output support for all CLI commands
