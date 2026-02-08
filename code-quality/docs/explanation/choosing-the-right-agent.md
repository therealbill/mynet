---
title: "Choosing the Right Agent"
description: "When to use each code quality agent"
weight: 2
---

# Choosing the Right Agent

The code-quality plugin provides five agents, each specialized for a different aspect of code quality. This page explains when to use each one and how they work together.

## Decision Matrix

| What you are doing | Agent to use |
|--------------------|--------------|
| Just finished writing or modifying code | code-reviewer |
| Added a new service, module, or changed inter-layer communication | architect-reviewer |
| Need tests for new or changed code | test-writer-fixer |
| Tests are failing after a refactor | test-writer-fixer |
| Need end-to-end browser tests for user flows | playwright-expert |
| Playwright tests are flaky or timing out | playwright-expert |
| Setting up Playwright in CI | playwright-expert |
| Need to audit web content for accessibility | web-accessibility-checker |
| Preparing for WCAG AA compliance | web-accessibility-checker |
| Checking form accessibility | web-accessibility-checker |

When the situation is ambiguous, consider what kind of feedback you need:

- **"Is this code correct?"** -- code-reviewer
- **"Is this code well-structured?"** -- architect-reviewer
- **"Is this code tested?"** -- test-writer-fixer
- **"Does this work in a browser?"** -- playwright-expert
- **"Can everyone use this?"** -- web-accessibility-checker

## Combining Agents

The most thorough quality workflow uses multiple agents in sequence. The order matters because each agent builds on assumptions that previous agents have validated.

**For a backend change:**

1. code-reviewer -- catch bugs and security issues
2. architect-reviewer -- verify structural consistency (if the change is structural)
3. test-writer-fixer -- write or update tests

**For a frontend/UI change:**

1. code-reviewer -- catch bugs in component logic
2. test-writer-fixer -- write unit tests for component behavior
3. playwright-expert -- write E2E tests for user flows
4. web-accessibility-checker -- audit accessibility compliance

**For a refactor:**

1. architect-reviewer -- verify the refactor respects boundaries and patterns
2. code-reviewer -- catch bugs introduced by the refactor
3. test-writer-fixer -- fix any tests broken by the refactor

Not every change needs every agent. A one-line bug fix needs code-reviewer and possibly test-writer-fixer. A full-stack feature might use all five.

## What These Agents Do Not Do

The code-quality agents are reviewers and testers. They analyze and validate code but do not write production application code.

- **They do not implement features.** If you need code written, describe the feature to Claude Code directly or use a development-focused plugin (web-development, backend-development, cli-development).
- **They do not deploy code.** Deployment, infrastructure, and CI/CD pipeline creation are outside their scope.
- **They do not manage dependencies.** Package updates, version conflicts, and dependency audits are not part of their review process.
- **They do not replace human judgment.** The agents flag issues and provide recommendations. The decision to fix, defer, or dismiss a finding remains with the developer.

## When NOT to Use Them

Some changes do not benefit from code quality review:

- **Simple typo fixes.** Correcting a variable name or fixing a comment does not need a full code review pass.
- **Documentation-only changes.** Changes to markdown files, comments, or README content are outside the scope of these agents.
- **Configuration changes.** Updating environment variables, CI configuration, or deployment manifests does not trigger useful feedback from code-reviewer or architect-reviewer.
- **Generated code.** Files produced by code generators, schema compilers, or build tools should not be reviewed as if they were hand-written. The generator's configuration is what matters, not its output.

In these cases, proceed directly to committing without running the quality agents. Using the agents on changes that do not benefit from them adds time without adding value.

## See Also

- [Agent Reference](../../reference/agents/) -- technical specifications for all five agents
- [Architecture](../../explanation/architecture/) -- why five specialized agents instead of one
- [Getting Started](../../tutorials/getting-started/) -- tutorial walkthrough of the review-test-fix cycle
