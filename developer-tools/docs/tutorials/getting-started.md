---
title: "Getting Started with Developer Tools"
description: "Set up a Git workflow with git-workflow-manager and generate docs with documentation-engineer"
weight: 1
---

# Getting Started with Developer Tools

Learn to establish a professional Git workflow and generate API documentation for a project that currently has neither, using two developer-tools agents working in sequence.

## What You'll Build

By the end of this tutorial, you will have:

- Described your project context to git-workflow-manager and received a complete branching strategy
- Reviewed branch naming, protection rules, and merge policy tailored to your team
- Set up automation including PR templates, commit hooks, and CI triggers
- Switched to documentation-engineer and generated API reference documentation from your source code
- Understood how these two agents address different parts of the developer workflow

This takes about 20 minutes.

## Prerequisites

- Claude Code CLI installed and authenticated (run `claude --version` to verify)
- The developer-tools plugin installed in your project's `.claude/settings.json`
- A project with source code that has no formal Git workflow (no branch protection, no commit conventions) and no structured documentation

If your project already has a Git workflow but needs documentation, skip to Step 4.

## Step 1: Describe Your Team and Release Cadence

Start by telling Claude Code about your development context. git-workflow-manager needs to understand your constraints before recommending a strategy. Any of these phrases activate the agent:

```
Set up our Git workflow
```

```
Main branch history is a mess
```

When the agent activates, provide concrete details about your situation:

```
We're a team of 4. We deploy to production twice a week from main.
We don't have branch protection or commit conventions.
PRs sometimes sit for days without review.
```

git-workflow-manager uses these constraints to select an appropriate branching model. A small team deploying frequently gets a different recommendation than a large team shipping quarterly releases. The agent reads your repository structure, existing Git configuration, and CI setup to ground its recommendations in your actual project state.

## Step 2: Review the Branching Strategy

The agent produces a branching strategy document covering:

- **Branch naming** -- conventions like `feature/`, `fix/`, `release/` prefixes with descriptive slugs
- **Protection rules** -- which branches require PR review, status checks, and linear history
- **Merge policy** -- squash-merge for feature branches (clean history), real merge commits for releases (preserve context)

For a team of 4 deploying twice weekly, the agent typically recommends trunk-based development or GitHub Flow rather than full Git Flow. Read the output carefully. If the strategy feels too heavyweight for your team, say so -- the agent adjusts based on feedback.

### Checkpoint

At this point you should have:

- A branching model recommendation matched to your team size and release cadence
- Branch naming conventions documented
- Protection rules and merge policy specified

If git-workflow-manager did not activate, verify the developer-tools plugin is listed in your `.claude/settings.json` and restart Claude Code.

## Step 3: Review the Automation

After establishing the branching strategy, git-workflow-manager configures automation to enforce it:

- **PR templates** -- structured templates that prompt authors for context, testing steps, and reviewer notes
- **Commit hooks** -- local hooks that validate commit message format against Conventional Commits (`feat:`, `fix:`, `docs:`, `chore:`)
- **CI triggers** -- workflows that run on PR creation and push to protected branches
- **Labels and auto-merge** -- label-driven automation for routine dependency updates or documentation changes

Review each automation artifact. The agent creates files in your repository (`.github/pull_request_template.md`, `.husky/` hooks, workflow files) that you can inspect before committing. Test a commit hook by making a small change and verifying the hook validates your message format.

## Step 4: Switch to documentation-engineer

Now shift from workflow to documentation. Tell Claude Code about your undocumented API:

```
We need docs for our API — it has 12 endpoints and zero documentation
```

```
Nobody reads our docs because they don't exist yet
```

documentation-engineer activates and begins by auditing your codebase. It reads route definitions, request/response schemas, authentication patterns, and middleware to understand what your API actually does before generating anything.

## Step 5: Review the Generated API Reference

documentation-engineer produces structured reference documentation organized by resource:

- **Endpoint listing** -- each route with its HTTP method, path, and purpose
- **Request parameters** -- path params, query params, and request body schemas with types and validation rules
- **Response formats** -- success and error response shapes with status codes
- **Authentication requirements** -- which endpoints require auth and what token format they expect

The generated docs live in your repository as Markdown files, typically under `docs/` or a location the agent selects based on your project structure. They are code artifacts -- they go through PR review, get versioned alongside your source, and can be validated in CI.

Review the output against your actual API behavior. If an endpoint is missing, the agent likely could not detect it from static analysis. Point it out and the agent adds it.

### Checkpoint

At this point you should have:

- API reference documentation covering all detected endpoints
- Documentation files committed to your repository
- An understanding of the docs-as-code workflow

## What You Learned

In this tutorial, you:

- **Described your project context** to git-workflow-manager and received a branching strategy matched to your team size and release cadence
- **Reviewed automation artifacts** including PR templates, commit hooks, and CI triggers that enforce the workflow
- **Triggered documentation-engineer** to audit your codebase and generate API reference documentation
- **Saw docs-as-code in action** -- documentation generated from source, stored in the repo, and reviewed like any other code change

## Next Steps

- [Set Up Git Workflow]({{< ref "howto/set-up-git-workflow" >}}) -- detailed guide for branching strategy and automation configuration
- [Generate API Documentation]({{< ref "howto/generate-api-documentation" >}}) -- advanced documentation generation with CI validation
- [Architecture]({{< ref "explanation/architecture" >}}) -- understand why developer-tools uses three specialized agents
- For Diataxis-structured documentation beyond API reference, see the [diataxis-docs](/diataxis-docs/) plugin
- For code review integration in your new workflow, see the [code-quality](/code-quality/) plugin
