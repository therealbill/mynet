---
title: "Getting Started with Code Quality"
description: "Learn the review-test-fix cycle using code-reviewer and test-writer-fixer"
weight: 1
---

# Getting Started with Code Quality

Learn the core review-test-fix workflow by making a code change, reviewing it for bugs and security issues, then writing tests to lock in correctness.

## What You'll Build

By the end of this tutorial, you will have:

- Made a code change and triggered code-reviewer to analyze it
- Read severity-grouped findings (Critical, Warnings, Suggestions)
- Triggered test-writer-fixer to create tests for the changed code
- Seen the test-write-fix cycle run end to end

This takes about 15 minutes.

## Prerequisites

- Claude Code CLI installed and authenticated (run `claude --version` to verify)
- The code-quality plugin installed in your project's `.claude/settings.json`
- A project with an existing codebase and recent code changes (or the willingness to make a small modification)

## Step 1: Make a Code Change

Open your project and make a small, meaningful modification. For example, add a new function, fix a bug, or change how an API endpoint handles errors. The change should touch at least one file with existing logic -- this gives the reviewer real context to work with.

If you do not have a change ready, modify an existing function to handle a new edge case:

```
Add input validation to the createUser function — reject empty email addresses
```

Stage or leave the changes unstaged. code-reviewer examines both staged and unstaged diffs.

## Step 2: Trigger code-reviewer

Ask Claude Code to review your changes. Any of these phrases activate the code-reviewer agent:

```
Review the changes I just made
```

```
Does this look right to you?
```

```
Check my code before I push
```

Claude Code matches your request to the code-reviewer agent. The agent runs `git diff` to see what changed, then reads the full context of each modified file -- not just the diff hunks. This context is essential for catching issues that span multiple functions or depend on surrounding code.

## Step 3: Review the Findings

code-reviewer groups its output by severity:

- **Critical** -- issues that will cause bugs or security vulnerabilities. These must be fixed before merging. Examples: nil dereference on an unchecked return value, SQL injection from unsanitized input, a race condition on shared state.
- **Warnings** -- issues that will cause maintenance problems. These should be fixed. Examples: confusing variable names, duplicated logic across functions, missing context in error messages.
- **Suggestions** -- improvements that are not urgent. Examples: a function doing slightly too many things, an opportunity to extract shared logic.

For each finding, code-reviewer provides the file path, line number, what is wrong, and how to fix it. If the code is clean, it says "No issues found" rather than manufacturing feedback.

### Checkpoint

At this point you should have:

- Triggered code-reviewer with one of the activation phrases
- Received a severity-grouped report of findings
- Understood which issues are critical versus advisory

If code-reviewer did not activate, verify the code-quality plugin is listed in your `.claude/settings.json` and restart Claude Code.

## Step 4: Trigger test-writer-fixer

Now ask Claude Code to write tests for the code you changed:

```
Write tests for the changes I just made
```

Claude Code matches this to the test-writer-fixer agent. The agent identifies the modified files, examines the existing test patterns in your project, and writes new tests that follow those patterns.

test-writer-fixer produces tests that:

- Test behavior, not implementation -- assertions target observable outcomes
- Use descriptive names that read as specifications (e.g., `"returns error when email is empty"`)
- Cover the happy path, error conditions, and edge cases in that priority order
- Follow your project's existing test framework and conventions

## Step 5: See the Test-Write-Fix Cycle

After test-writer-fixer creates the tests, it runs them automatically. If any test fails, the agent classifies the failure:

- **Real bug** -- the test caught a genuine problem in your code. The agent reports the bug rather than modifying the test.
- **Stale expectation** -- your code behavior legitimately changed, but an existing test still asserts the old behavior. The agent updates the assertion to match the new behavior.
- **Brittle test** -- a test depends on implementation details, timing, or ordering. The agent refactors it for resilience while preserving the original intent.
- **Environment issue** -- missing dependencies, configuration, or state. The agent fixes the setup.

Watch this cycle: write tests, run them, classify any failures, fix accordingly. This is the core loop you will use throughout your development workflow.

### Checkpoint

At this point you should have:

- New test files or test cases covering your code changes
- All tests passing after the agent's write-fix cycle
- An understanding of how the agent classifies and handles test failures

## What You Learned

In this tutorial, you:

- **Triggered code-reviewer** by asking for a review of your changes -- the agent analyzed your diff for correctness, security, and maintainability issues
- **Read severity-grouped findings** and understood the difference between Critical issues that block merging and Suggestions that can wait
- **Triggered test-writer-fixer** to create tests that follow your project's conventions and cover the changed code
- **Observed the test-write-fix cycle** where the agent writes tests, runs them, classifies failures, and fixes them without weakening test intent

## Next Steps

- [Review Code Before Committing]({{< ref "howto/review-code-before-committing" >}}) -- integrate code review into your pre-commit workflow
- [Write and Fix Tests]({{< ref "howto/write-and-fix-tests" >}}) -- advanced test creation and repair techniques
- [Audit Accessibility]({{< ref "howto/audit-accessibility" >}}) -- check web content for WCAG compliance
- [Choosing the Right Agent]({{< ref "explanation/choosing-the-right-agent" >}}) -- when to use each of the five code quality agents
- For language-specific code quality expertise, combine these agents with the [programming-languages](/programming-languages/) plugin
