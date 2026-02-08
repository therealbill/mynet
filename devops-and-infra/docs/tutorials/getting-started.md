---
title: "Getting Started with DevOps and Infrastructure"
description: "Set up GitHub Actions CI and add Prometheus monitoring"
weight: 1
---

# Getting Started with DevOps and Infrastructure

Set up a CI/CD pipeline with GitHub Actions and add Prometheus monitoring to your service using the devops-and-infra plugin's specialized agents.

## What You'll Build

By the end of this tutorial, you will have:

- Triggered github-actions-expert with a project description and received a complete CI/CD workflow
- Reviewed the workflow structure including parallel jobs, dependency caching, and security practices
- Switched to prometheus-expert and received instrumentation guidance for your service
- Reviewed metric types, naming conventions, and alerting rules
- Understood how the CI/CD and monitoring agents complement each other in the DevOps pipeline

This takes about 30 minutes.

## Prerequisites

- Claude Code CLI installed and authenticated (run `claude --version` to verify)
- The devops-and-infra plugin installed in your project's `.claude/settings.json`
- A project hosted on GitHub (this tutorial uses a Node.js API as the example, but any language works)
- Basic familiarity with CI/CD concepts (pipelines, stages, triggers)

## Step 1: Describe Your Project to github-actions-expert

Open Claude Code in your project directory and describe your CI/CD needs:

```
Set up GitHub Actions for our Node.js API. We use Jest for tests, deploy to
AWS ECS, and want the pipeline to run on pull requests and pushes to main.
```

Claude Code matches your request to the github-actions-expert agent based on the GitHub Actions context. The agent activates and begins designing your workflow.

## Step 2: Review the Generated Workflow

github-actions-expert produces a complete `.github/workflows/` configuration. For a Node.js API deploying to ECS, expect:

- **Triggers:** `pull_request` for test runs, `push` to `main` for deploy -- not every push to every branch
- **Jobs:** Separate jobs for lint, test, build, and deploy -- structured for parallelism where possible
- **Caching:** `actions/cache` configured for `node_modules` with a hash of `package-lock.json` as the cache key
- **Secrets management:** AWS credentials referenced via `${{ secrets.AWS_ACCESS_KEY_ID }}` -- never hardcoded
- **Pinned actions:** Third-party actions pinned to commit SHAs rather than version tags to prevent supply chain attacks
- **Permissions:** Minimum permissions declared at the job level, not the workflow level

Review the YAML carefully. github-actions-expert includes inline comments explaining why each configuration choice was made -- why a specific cache key was chosen, why jobs are structured in a particular order, why certain permissions are set.

## Step 3: See How Jobs Run in Parallel

github-actions-expert structures the workflow so independent jobs run concurrently. In the generated workflow, notice:

- **Lint and test run in parallel** -- they have no dependency on each other, so they execute simultaneously
- **Build depends on both lint and test** -- it uses `needs: [lint, test]` to wait for both to pass
- **Deploy depends on build** -- it only runs when the build succeeds and the trigger is a push to `main`

This parallel structure is intentional. A sequential pipeline (lint, then test, then build, then deploy) takes the sum of all job times. A parallel pipeline takes only as long as the longest parallel path. github-actions-expert optimizes for total pipeline time under 10 minutes.

### Checkpoint

At this point you should have:

- A `.github/workflows/` directory with a complete CI/CD configuration
- Parallel jobs for lint, test, build, and deploy
- Proper caching, secrets management, and pinned actions

If github-actions-expert did not activate, verify the devops-and-infra plugin is listed in your `.claude/settings.json` and restart Claude Code.

## Step 4: Switch to prometheus-expert for Monitoring

With the pipeline in place, add monitoring to the service that gets deployed. Ask:

```
Instrument our Node.js API for Prometheus. We need to track request latency,
error rates, and active connections. Alert if error rate exceeds 1% or p99
latency exceeds 500ms.
```

Claude Code routes this to prometheus-expert. The agent begins by understanding what your service does and what "healthy" looks like, then produces instrumentation guidance.

## Step 5: Review the Instrumentation

prometheus-expert produces Prometheus configuration covering several areas:

- **Metric types:** Histograms for request latency (not summaries -- histograms are aggregatable), counters for total requests and errors, gauges for active connections
- **Naming conventions:** Metrics follow the `<namespace>_<subsystem>_<name>_<unit>` pattern -- for example, `myapi_http_request_duration_seconds` rather than `request_time` or `latency`
- **Label design:** Labels like `method`, `path`, and `status_code` -- but with cardinality control. The agent collapses high-cardinality path parameters (e.g., `/users/:id` becomes `/users/{id}`) to prevent label explosion
- **Alerting rules:** Symptom-based alerts with `for` durations to avoid false positives. The error rate alert uses `rate()` over a 5-minute window, not instantaneous spikes. The latency alert uses `histogram_quantile(0.99, ...)` with a 10-minute `for` duration
- **Recording rules:** Pre-computed aggregations for dashboard queries that would otherwise be expensive at query time

Review the metric names, label cardinality, and alert thresholds. prometheus-expert includes comments explaining each decision -- why a histogram bucket distribution was chosen, why a specific `for` duration prevents false positives, why certain labels were excluded to control cardinality.

## What You Learned

In this tutorial, you:

- **Triggered github-actions-expert** by describing a CI/CD requirement -- the agent produced a complete GitHub Actions workflow with parallel jobs, caching, and security best practices
- **Reviewed pipeline structure** including job dependencies, trigger configuration, and secrets management
- **Triggered prometheus-expert** by describing monitoring requirements -- the agent produced metric definitions, naming conventions, and alerting rules
- **Reviewed instrumentation design** including metric type selection, label cardinality control, and symptom-based alerting

The core workflow is: **github-actions-expert builds the pipeline, prometheus-expert monitors what it deploys.** The CI/CD pipeline ensures code reaches production safely. Prometheus monitoring ensures you know when production is unhealthy.

## Next Steps

- [Set Up GitHub Actions CI]({{< ref "howto/set-up-github-actions-ci" >}}) -- detailed guide for CI/CD workflow configuration
- [Configure Prometheus Monitoring]({{< ref "howto/configure-prometheus-monitoring" >}}) -- detailed guide for metrics and alerting
- [Architecture]({{< ref "explanation/architecture" >}}) -- understand how the four agents divide the DevOps pipeline
- [Agent Reference]({{< ref "reference/agents" >}}) -- full specification of all four agents

For the services being deployed, see the [backend-development](/backend-development/) plugin. For pre-commit quality checks in your pipeline, see the [code-quality](/code-quality/) plugin.
