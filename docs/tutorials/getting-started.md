---
title: "Getting Started with Mynet"
description: "Install the code-quality plugin from the Mynet marketplace, make a code change, and run the code-reviewer agent to catch issues before you commit."
weight: 1
---

# Getting Started with Mynet

Install a plugin from the Mynet marketplace, trigger an agent on real code, and learn how plugins extend Claude Code with specialized capabilities.

## What You'll Build

By the end of this tutorial, you will have:

- Added the Mynet marketplace and installed a plugin using Claude Code's plugin manager
- Installed the **code-quality** plugin into a project
- Made a deliberate code change with a bug
- Triggered the **code-reviewer** agent to catch the bug
- Understood the agent's structured output

This takes about 15 minutes.

## Prerequisites

- Claude Code CLI installed and authenticated (run `claude --version` to verify)
- A project directory with at least one source file (any language)

## Step 1: Add the Mynet Marketplace

Inside a Claude Code session, register the Mynet marketplace so its plugins appear in the Discover tab:

```
/plugin marketplace add therealbill/mynet
```

This is a one-time setup step. Once added, the marketplace persists across sessions -- you do not need to add it again.

### Checkpoint

Run this command to confirm the marketplace was registered:

```
/plugin marketplace list
```

You should see `mynet` in the output. If it does not appear, re-run the add command from above.

## Step 2: Browse and Install the code-quality Plugin

You can install the plugin interactively or with a direct command.

**Option A: Interactive** -- Run `/plugin`, navigate to the Discover tab, find `code-quality`, and install it.

**Option B: Direct command:**

```
/plugin install code-quality@mynet --scope project
```

The `--scope project` flag makes the plugin available to anyone who works in this project directory. Without it, the plugin installs for your user only.

### Checkpoint

Run `/plugin` and check the Installed tab. You should see **code-quality** listed with its agents:

- architect-review
- code-reviewer
- playwright-expert
- test-writer-fixer
- web-accessibility-checker

If the plugin does not appear, verify the marketplace was added in Step 1 and try installing again.

## Step 3: Create a File with a Bug

Create a small source file with a deliberate bug. This gives the code-reviewer agent something to find.

Create a file called `example.py` in your project:

```python
def calculate_average(numbers):
    total = 0
    for n in numbers:
        total += n
    return total / len(numbers)

def process_user_input(user_data):
    query = f"SELECT * FROM users WHERE name = '{user_data['name']}'"
    return query
```

This file has two problems that a code reviewer should catch:

- `calculate_average` will crash with a `ZeroDivisionError` if given an empty list
- `process_user_input` has a SQL injection vulnerability

## Step 4: Stage the Change

Add the file to git so it appears in `git diff`:

```bash
git add example.py
```

The code-reviewer agent works by examining `git diff` output. Staging the file makes it visible to the agent.

### Checkpoint

At this point you should have:

- The code-quality plugin installed and visible in the `/plugin` Installed tab
- An `example.py` file staged in git with two known issues
- Run `git diff --cached` to verify the staged file appears in the diff

## Step 5: Trigger the Code-Reviewer Agent

In your Claude Code session, type this prompt:

```
Review the changes I just made
```

Claude Code recognizes that you are asking for a code review and routes your request to the **code-reviewer** agent from the code-quality plugin.

## Step 6: Read the Agent's Output

The code-reviewer agent runs through a structured review process:

1. It runs `git diff` to see your staged changes
2. It reads the full file for context
3. It analyzes the code against its review priorities
4. It reports findings grouped by severity

You should see output structured like this:

```
**Critical** (must fix before merging):

- example.py, line 5: `calculate_average` divides by `len(numbers)` without
  checking for an empty list. This will raise a ZeroDivisionError.
  Fix: Add a guard clause — `if not numbers: return 0`

- example.py, line 8: `process_user_input` builds a SQL query using string
  interpolation with unsanitized user input. This is a SQL injection
  vulnerability.
  Fix: Use parameterized queries instead of string formatting.
```

The agent categorizes issues into three levels:

- **Critical** -- bugs or security vulnerabilities that must be fixed before merging
- **Warnings** -- issues that will cause maintenance problems over time
- **Suggestions** -- improvements that are not urgent

The agent skips any category that has no findings. If the code is clean, it reports "No issues found."

> **What just happened?** The code-reviewer agent read your git diff, analyzed the code against its review priorities (correctness, security, maintainability, observability), and produced a structured report. It operates as a senior reviewer -- it catches real problems and avoids nitpicking style.

### Checkpoint

At this point you should have:

- Received a structured code review from the code-reviewer agent
- Seen at least two Critical findings (the division-by-zero and SQL injection)
- Understood the severity categories in the output

If Claude Code did not route to the code-reviewer agent, verify that the plugin is enabled in the `/plugin` Installed tab.

## Step 7: Fix the Issues and Re-review

Update `example.py` to fix both problems:

```python
def calculate_average(numbers):
    if not numbers:
        return 0
    total = 0
    for n in numbers:
        total += n
    return total / len(numbers)

def process_user_input(user_data, db_connection):
    query = "SELECT * FROM users WHERE name = %s"
    return db_connection.execute(query, (user_data['name'],))
```

Stage the updated file:

```bash
git add example.py
```

Then ask Claude Code again:

```
Review the changes I just made
```

You should see a clean report with no Critical findings. The agent confirms the fixes address the original issues.

## What You Learned

In this tutorial, you:

- **Added the Mynet marketplace** to Claude Code -- a one-time registration that makes all marketplace plugins discoverable
- **Installed a plugin** using the `/plugin` command with project scope
- **Triggered an agent** by asking Claude Code a natural language question that matched the agent's description
- **Interpreted structured output** from the code-reviewer agent, understanding its severity categories and fix recommendations
- **Iterated on feedback** by fixing issues and re-running the review

## Next Steps

Now that you have a working plugin installation, explore further:

- **Try other agents in code-quality** -- the `test-writer-fixer` agent generates tests for your code, and the `architect-review` agent evaluates your project's architecture
- **Browse the marketplace** -- the [Plugin Catalog](../../reference/plugin-catalog/) lists all available plugins with their component counts
- **Install multiple plugins** -- see [How to Use Plugins Together](../../howto/use-plugins-together/) for combining agents across plugins
- **Find the right plugin** -- use [How to Find the Right Plugin](../../howto/find-right-plugin/) to match plugins to your workflow

## Troubleshooting

**Plugin not showing in Discover tab**

Verify the marketplace was added with `/plugin marketplace list`. If it is not listed, re-run `/plugin marketplace add therealbill/mynet`.

**Agent not triggering after install**

Verify the plugin is enabled in the `/plugin` Installed tab. The agent activates on prompts related to code review. Use phrases like "review the changes I just made," "check my code before I push," or "does this look right to you?"

**git diff shows no output**

Run `git diff --cached` to check staged changes. If you forgot to run `git add`, the changes are unstaged and `git diff` (without `--cached`) shows them. The code-reviewer agent checks both staged and unstaged changes, but the file must be tracked by git.

**The review output is empty or says "No issues found"**

The agent only reports real problems. If your code is clean, a "No issues found" response is correct. To verify the agent is working, introduce a deliberate bug (like the division-by-zero example above) and re-run the review.
