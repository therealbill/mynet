---
title: "Review Code Before Committing"
description: "Use code-reviewer for pre-commit code review"
weight: 1
---

# Review Code Before Committing

Catch bugs, security issues, and maintainability problems before they reach your commit history using the code-reviewer agent.

## Prerequisites

- Claude Code with the code-quality plugin installed
- A project with staged or unstaged code changes

## Steps

### 1. Stage Your Changes

Stage the files you intend to commit. code-reviewer examines both staged and unstaged diffs, but staging first ensures you and the agent are looking at the same set of changes:

```bash
git add -p
```

For large changesets, consider reviewing in batches. Stage one logical unit of work at a time rather than the entire diff.

### 2. Ask code-reviewer to Review

Use any of these phrases to activate the code-reviewer agent:

```
Review the changes I just made
```

```
Does this look right to you?
```

```
Check my code before I push
```

code-reviewer runs `git diff` (staged and unstaged), reads the full context of every modified file, then reviews the changes against its priority list: correctness, security, maintainability, and observability.

### 3. Read the Severity-Grouped Output

code-reviewer organizes findings into three levels:

- **Critical** -- must fix before merging. Logic errors, nil dereference, race conditions, injection vulnerabilities, exposed secrets, broken trace propagation.
- **Warnings** -- should fix. Confusing names, duplicated logic, functions doing too many things, missing error context, unbounded spans in loops.
- **Suggestions** -- consider fixing. Improvements that are not urgent but would benefit the codebase.

Each finding includes the file path, line number, description of the problem, and a concrete fix.

### 4. Fix Critical Issues

Address every Critical finding before committing. These are issues that will cause bugs in production or create security vulnerabilities. Common Critical findings include:

- Unchecked error returns leading to nil dereference
- SQL injection from string concatenation in queries
- Race conditions on shared state without synchronization
- Using `context.Background()` where a propagated context is available, severing distributed traces

### 5. Address Warnings

Fix Warnings in the same commit when practical. These issues will not break production immediately but will make the code harder to maintain. If a Warning requires a larger refactor that is out of scope, note it for a follow-up.

### 6. Consider Suggestions

Review Suggestions and apply the ones that improve clarity without adding scope. It is acceptable to defer Suggestions to a separate commit or ignore them if they conflict with your immediate goals.

### 7. Verification

After fixing the reported issues, run code-reviewer again to confirm the fixes are clean:

```
Review my changes again
```

A clean review returns "No issues found." If new findings appear, they are typically secondary issues that were masked by the original problems.

Once the review is clean, commit with confidence:

```bash
git commit
```

## Troubleshooting

**"No issues found" when you expected feedback:**

code-reviewer only flags real problems. If the code is clean, it says so. If you suspect it missed something specific, ask directly: "Check if the database query in user.go is vulnerable to SQL injection." Targeted questions prompt deeper inspection of a specific concern.

**Large diffs produce too many findings:**

Break the review into smaller units. Stage one file or one logical change at a time and ask for a focused review: "Review only the changes in auth/middleware.go." This produces more actionable output than reviewing 20 files at once.

**Findings about code you did not change:**

code-reviewer reads the full context of modified files, not just the diff. If it flags an existing issue in a file you touched, that is intentional -- your change may interact with the existing problem. Decide whether to fix it in this commit or defer it.

## See Also

- [Getting Started]({{< ref "tutorials/getting-started" >}}) -- tutorial walkthrough of the review-test-fix cycle
- [Agent Reference]({{< ref "reference/agents" >}}) -- code-reviewer specification
- [Architecture]({{< ref "explanation/architecture" >}}) -- why code-reviewer is a separate specialized agent
