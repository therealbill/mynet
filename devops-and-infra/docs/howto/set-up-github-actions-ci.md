---
title: "Set Up GitHub Actions CI"
description: "Use github-actions-expert to create CI/CD workflows with proper caching and security"
weight: 1
---

# Set Up GitHub Actions CI

Create a CI/CD pipeline for your project with parallel jobs, dependency caching, security-pinned actions, and minimum permissions.

## Prerequisites

- Claude Code with the devops-and-infra plugin installed
- A project hosted on GitHub
- Knowledge of your project's language, test framework, and deployment target

## Steps

### 1. Describe Your Project

Tell github-actions-expert what needs to be built, tested, and deployed. Include:

- Programming language and version
- Test framework and how tests are run
- Build output (Docker image, static files, compiled binary)
- Deployment target (AWS, GCP, Vercel, self-hosted)
- Branch strategy (which branches trigger which stages)

Example prompts that trigger github-actions-expert:

- "Set up GitHub Actions for tests and deploy"
- "CI takes 25 minutes, help me speed it up"
- "Same CI config copy-pasted across 10 repos"

The more context you provide, the more precise the generated workflow.

### 2. Trigger github-actions-expert

github-actions-expert activates when it detects GitHub Actions, CI/CD pipeline, or workflow optimization context in your request. If it does not activate automatically, be explicit:

```
Create a GitHub Actions workflow for a Python FastAPI app. Tests run with
pytest, the app deploys as a Docker image to Google Cloud Run, and we need
matrix builds for Python 3.11 and 3.12.
```

The agent reads any existing `.github/workflows/` files first to understand your current setup before proposing changes.

### 3. Review Workflow Structure

github-actions-expert produces a workflow with these characteristics:

- **Parallel jobs:** Independent stages (lint, test, type-check) run simultaneously. Dependent stages (build, deploy) use `needs:` to wait
- **Concurrency groups:** Configured to cancel in-progress runs when a new commit is pushed to the same branch, preventing resource waste
- **Matrix builds:** Used when testing across multiple language versions or operating systems. The matrix is defined explicitly, not as a combinatorial explosion
- **Triggers:** Precise event filters. Pull requests run tests. Pushes to `main` run the full pipeline including deploy. Manual triggers via `workflow_dispatch` are included for emergency deploys

Verify that the job dependency graph matches your expected pipeline flow. Jobs that can run in parallel should not have unnecessary `needs:` constraints.

### 4. Review Security Practices

github-actions-expert applies several security measures by default:

- **Pinned actions:** Third-party actions are pinned to commit SHAs (`uses: actions/checkout@a5ac7e...`) rather than version tags. Version tags can be reassigned; commit SHAs cannot
- **Minimum permissions:** Each job declares only the permissions it needs. The test job gets `contents: read`. The deploy job gets `contents: read` and `id-token: write` for OIDC. No job gets `write-all`
- **Secrets in environment:** Secrets are referenced via `${{ secrets.NAME }}`, never echoed to logs or passed as command-line arguments
- **Dependency review:** For pull requests, the workflow can include `actions/dependency-review-action` to flag new dependencies with known vulnerabilities

Check that no secrets are exposed in build logs and that the permissions for each job are the minimum required for that job's tasks.

### 5. Commit the Workflow and Verify

Commit the `.github/workflows/` files and push to trigger the pipeline:

```
git add .github/workflows/
git commit -m "Add CI/CD workflow"
git push
```

Open the Actions tab in your GitHub repository and verify that:

- The workflow triggers on the correct events
- Parallel jobs run simultaneously, not sequentially
- Caching reduces dependency installation time on the second run
- All stages complete successfully

## Verification

A properly configured pipeline meets these criteria:

- [ ] Total pipeline time is under 10 minutes for a standard test-build-deploy cycle
- [ ] All stages pass on the first run (no missing secrets or misconfigured permissions)
- [ ] Cache hit on the second run reduces dependency installation to under 10 seconds
- [ ] Pull requests run tests but do not deploy
- [ ] Pushes to `main` run the full pipeline including deploy

## Troubleshooting

**Pipeline is slow despite parallel jobs:**

Check the caching configuration. The most common cause of slow pipelines is re-downloading dependencies on every run. Verify that the cache key includes the lockfile hash (`package-lock.json`, `go.sum`, `requirements.txt`) and that the cache path matches where your package manager stores downloads.

**Flaky tests in CI but not locally:**

Matrix builds can expose environment-dependent failures. Check whether the test relies on a specific timezone, locale, or file system behavior. github-actions-expert can add `services:` blocks for databases or other dependencies that tests need.

**Same workflow needed across multiple repositories:**

Ask github-actions-expert about reusable workflows and composite actions. Reusable workflows live in a central repository and are called from other repositories with `uses: org/repo/.github/workflows/ci.yml@main`. Composite actions bundle multiple steps into a single action that can be versioned and shared.

## See Also

- [Getting Started](../../tutorials/getting-started/) -- tutorial covering CI setup and monitoring together
- [Architecture](../../explanation/architecture/) -- why github-actions-expert is separate from devops-automator
- [Agent Reference](../../reference/agents/) -- github-actions-expert specification
