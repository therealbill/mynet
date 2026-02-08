---
title: "Architecture"
description: "The DevOps pipeline — build, deploy, monitor — and agent roles"
weight: 1
---

# Architecture

The devops-and-infra plugin covers the full DevOps pipeline: building code, deploying it, and monitoring it in production. This page explains how four agents divide that pipeline, why each agent exists as a separate specialist, and how they relate to each other and to other plugins.

## The DevOps Pipeline

The plugin maps to three stages of the DevOps lifecycle:

1. **Build:** github-actions-expert handles CI/CD pipeline configuration -- the YAML workflows that test, build, and package code on every push
2. **Deploy:** devops-automator handles infrastructure and deployment strategy -- the Terraform configs, Docker images, blue-green rollouts, and GitOps workflows that get code into production
3. **Monitor:** performance-monitor and prometheus-expert handle observability -- the metrics, alerts, and diagnostic analysis that tell you whether production is healthy

Each stage has different inputs, outputs, and expertise requirements. A GitHub Actions workflow is YAML with a specific schema, job dependency graph, and caching model. A Terraform configuration is HCL with a state model, dependency graph, and provider ecosystem. A Prometheus alerting rule is PromQL with a time series model, range vector semantics, and label cardinality constraints. These are fundamentally different domains that share a DevOps umbrella.

## The Build-Deploy Split

github-actions-expert and devops-automator operate on different levels of abstraction.

github-actions-expert works at the CI/CD pipeline level. Its domain is GitHub Actions: workflow files, job definitions, matrix builds, caching strategies, action pinning, and concurrency groups. The output is `.github/workflows/*.yml` files that define how code moves from a commit to a tested, built artifact.

devops-automator works at the infrastructure and deployment level. Its domain is everything after the artifact is built: Docker image optimization, Terraform or Pulumi configurations, blue-green vs. canary deployment strategies, rollback mechanisms, and preview environments. The output is infrastructure code and deployment configurations.

The split exists because these two domains have different optimization targets. A CI/CD pipeline optimizes for speed and feedback quality -- run tests fast, cache aggressively, fail early. A deployment strategy optimizes for reliability and safety -- roll out gradually, verify health, roll back automatically. The person who knows that `actions/cache` needs a lockfile-based key is not necessarily the person who knows when to choose canary over blue-green.

In practice, the boundary is the built artifact. github-actions-expert produces it; devops-automator deploys it.

## The Monitoring Pair

prometheus-expert and performance-monitor both deal with production observability, but they approach it from different directions.

prometheus-expert is a tool specialist. Its domain is Prometheus: metric types (counters, gauges, histograms, summaries), the `<namespace>_<subsystem>_<name>_<unit>` naming convention, PromQL query construction, recording rules, alerting rules with `for` durations, label cardinality management, and Alertmanager routing. When you need to write a Prometheus alerting rule or figure out why Prometheus is running out of memory, prometheus-expert has the specific technical knowledge.

performance-monitor is a diagnostician. Its domain is performance analysis: following a request from edge to database, distinguishing capacity problems from efficiency problems, establishing Four Golden Signals baselines, designing SLO-based burn-rate alerts, and producing structured findings with evidence, root causes, and recommendations. performance-monitor is not tied to any specific monitoring tool. It reasons about performance in general terms and recommends what to measure, not how to configure a specific tool.

The split mirrors a real-world distinction. The person who knows that `rate()` range vectors should be at least 4x the scrape interval is a Prometheus specialist. The person who knows that a latency spike is caused by connection pool exhaustion rather than slow queries is a performance diagnostician. These are different skills. prometheus-expert handles the former. performance-monitor handles the latter.

When both are needed, the typical flow is: performance-monitor identifies the problem and recommends what to monitor, then prometheus-expert implements the specific Prometheus configuration.

## Why All Agents Use Sonnet

All four agents run on sonnet. DevOps tasks are well-defined technical patterns rather than open-ended architectural reasoning.

A GitHub Actions workflow has a specific schema. A Prometheus alerting rule has specific syntax. A deployment strategy has a finite set of options (blue-green, canary, rolling, recreate) with known tradeoffs. A performance diagnosis follows a structured methodology (symptoms, request path, capacity vs. efficiency, root cause).

These tasks require domain expertise -- knowing the right pattern for a given situation -- more than they require the broad, open-ended reasoning that opus provides. The CI/CD pipeline that github-actions-expert produces is not a novel design. It is the correct application of known best practices (cache aggressively, pin actions, minimize permissions, parallelize jobs) to a specific project. Sonnet handles this pattern-matching efficiently.

## Cross-Plugin Relationships

The devops-and-infra plugin sits downstream from development plugins and upstream from production:

- **backend-development** produces the services that devops-and-infra deploys and monitors. backend-architect designs the system; devops-automator deploys it. The boundary is the built artifact and its infrastructure requirements.
- **code-quality** provides quality gates that run inside the CI/CD pipeline. Linting, testing, and security scanning happen within the GitHub Actions workflow that github-actions-expert creates. code-quality agents define what checks to run; github-actions-expert defines where and how they run in the pipeline.
- **web-development** and **cli-development** produce frontend and CLI applications that also need CI/CD pipelines and monitoring. github-actions-expert handles their pipeline configuration with the same principles (caching, parallelism, security) adapted to different build tools and deployment targets.

These plugins do not import each other. They collaborate through Claude's dispatch layer. backend-architect may specify that a service needs a health check endpoint. devops-automator references that endpoint in a deployment health probe. The boundary is clean -- each plugin owns its domain.

## See Also

- [Agent Reference](../../reference/agents/) -- full specification of all four agents
- [Getting Started](../../tutorials/getting-started/) -- see the CI and monitoring agents in practice
- [Set Up GitHub Actions CI](../../howto/set-up-github-actions-ci/) -- practical guide for the build stage
- [Configure Prometheus Monitoring](../../howto/configure-prometheus-monitoring/) -- practical guide for the monitoring stage
