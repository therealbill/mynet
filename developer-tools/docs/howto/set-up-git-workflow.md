---
title: "Set Up Git Workflow"
description: "Use git-workflow-manager to establish branching strategy and automation"
weight: 1
---

# Set Up Git Workflow

Establish a branching strategy with commit conventions, merge policy, and automation using git-workflow-manager.

## Goal

By the end of this guide, your project will have:

- A branching model matched to your team's constraints
- Commit message conventions enforced by hooks
- PR templates and merge policy configured
- Automation for releases and routine maintenance

## Step 1: Describe Your Constraints

git-workflow-manager needs context before it can recommend a strategy. Provide these details when activating the agent:

- **Team size** -- solo developer, small team (2-6), or large team (7+)
- **Release cadence** -- continuous deployment, weekly releases, or versioned releases on a schedule
- **CI integration** -- what CI system you use (GitHub Actions, GitLab CI, Jenkins) and whether it currently runs on PRs
- **Current pain points** -- messy history, long-lived branches, merge conflicts, unclear release process

```
Set up our Git workflow — we're 3 developers, deploy weekly,
use GitHub Actions, and our main branch history is hard to follow
```

The agent reads your repository's current Git configuration, branch structure, and CI files to ground its recommendations in your actual setup rather than starting from a blank slate.

## Step 2: Review the Branching Model

git-workflow-manager selects from established branching strategies:

- **Trunk-based development** -- all work lands on `main` through short-lived feature branches. Best for small teams with continuous deployment and strong CI coverage.
- **GitHub Flow** -- feature branches merge to `main` via PR. Slightly more structured than trunk-based, with explicit review gates. Good for teams that deploy after PR approval.
- **Git Flow** -- `develop`, `release`, and `hotfix` branches alongside `main`. Suited for teams shipping versioned releases on a schedule with long stabilization periods.

The agent picks the simplest model that satisfies your constraints. If it recommends Git Flow and you deploy continuously, push back -- the agent adjusts.

## Step 3: Review Commit Conventions and Merge Policy

The agent configures Conventional Commits with these prefixes:

- `feat:` -- new functionality
- `fix:` -- bug repair
- `docs:` -- documentation changes
- `chore:` -- maintenance, dependencies, tooling
- `refactor:` -- code restructuring without behavior change
- `test:` -- test additions or modifications

Merge policy varies by branch type:

- **Feature branches** -- squash-merge to keep `main` history clean, one commit per feature
- **Release branches** -- real merge commits to preserve the full context of a release
- **Hotfix branches** -- merge commit with clear provenance back to `main`

## Step 4: Review Automation

The agent creates or modifies these files in your repository:

- **PR template** (`.github/pull_request_template.md`) -- structured sections for description, testing steps, and reviewer checklist
- **Commit hooks** (`.husky/commit-msg` or similar) -- validates commit message format on every commit
- **CI workflow** -- runs linting, tests, and status checks on PR branches
- **Release workflow** -- tag-driven releases with changelog generation from Conventional Commits
- **Labels** -- GitHub labels for PR categorization (`feature`, `fix`, `docs`, `breaking`)

If your project uses a monorepo or has special CI requirements, the agent adapts the automation accordingly.

## Step 5: Review Release Management

For projects that ship versioned releases, the agent configures:

- **Semantic versioning** -- `feat:` bumps minor, `fix:` bumps patch, `BREAKING CHANGE:` bumps major
- **Changelog generation** -- automated from commit messages, grouped by type
- **Tag-driven deployment** -- pushing a version tag triggers the release pipeline

For continuous deployment projects, this step is lighter -- every merge to `main` is a release, and the agent configures that pipeline instead.

## Verification

Test the workflow end to end:

1. Create a feature branch following the naming convention
2. Make a commit -- verify the hook validates your message format
3. Push and open a PR -- verify the template appears and CI runs
4. Merge the PR -- verify the merge policy (squash vs merge commit) matches configuration

If all four steps work as expected, the workflow is operational.

## Troubleshooting

**Trunk-based is not enough.** If your team needs release stabilization periods or maintains multiple versions simultaneously, tell the agent. It will recommend GitHub Flow or Git Flow with the appropriate branch structure.

**Conflicts with existing workflow.** If your team already has conventions that differ from the agent's recommendations, describe them explicitly. git-workflow-manager adapts to existing practices rather than overwriting them. The agent implements minimal changes to solve the stated problem.

**Hooks not running.** Verify that your hook manager (Husky, lefthook, or similar) is installed and that `npm install` or equivalent has been run. Check that the hook files are executable (`chmod +x`).

**Agent-parallel workflows.** For teams running multiple Claude Code agents concurrently, git-workflow-manager can configure Git worktrees so each agent works in an isolated directory without branch conflicts. Mention this need explicitly when describing your constraints.
