---
title: "Getting Started with Mynet"
description: "Install the code-quality plugin from the Mynet marketplace, make a code change, and run the code-reviewer agent to catch issues before you commit."
weight: 1
---

# Getting Started with Mynet

Install a plugin from the Mynet marketplace, trigger an agent on real code, and learn how plugins extend Claude Code with specialized capabilities.

## What You'll Build

By the end of this tutorial, you will have:

- Cloned the Mynet plugin marketplace
- Installed the **code-quality** plugin into a project
- Made a deliberate code change with a bug
- Triggered the **code-reviewer** agent to catch the bug
- Understood the agent's structured output

This takes about 15 minutes.

## Prerequisites

- Claude Code CLI installed and authenticated (run `claude --version` to verify)
- Git installed
- A project directory with at least one source file (any language)

## Step 1: Clone the Mynet Marketplace

Open your terminal and clone the Mynet repository:

```bash
git clone https://github.com/your-org/claude-plugins.git ~/mynet
```

This downloads the full marketplace. Each top-level directory is an independent plugin.

Run this command to see the available plugins:

```bash
ls -d ~/mynet/*/
```

You should see directories like `code-quality/`, `web-development/`, `backend-development/`, `devops-and-infra/`, and others. Each one is a self-contained plugin you can install into any project.

## Step 2: Explore the Marketplace Manifest

Open the marketplace manifest to see what plugins are available:

```bash
cat ~/mynet/.claude-plugin/marketplace.json
```

You should see a JSON file listing every plugin with its name, description, version, and keywords. Look for the `code-quality` entry:

```json
{
  "name": "code-quality",
  "source": "./code-quality",
  "description": "Code review, testing, accessibility, and architectural quality agents",
  "version": "1.0.0",
  "keywords": ["code-review", "testing", "accessibility", "architecture", "quality"]
}
```

This tells you the plugin provides agents for code review, testing, accessibility checking, and architecture review.

## Step 3: Examine the Plugin Structure

Look inside the code-quality plugin:

```bash
ls ~/mynet/code-quality/agents/
```

You should see:

```
architect-review.md
code-reviewer.md
playwright-expert.md
test-writer-fixer.md
web-accessibility-checker.md
```

Each `.md` file defines one agent. The **code-reviewer** agent is the one you will use in this tutorial. It reviews recent code changes for correctness, security, and maintainability.

### Checkpoint

At this point you should have:

- The Mynet repository cloned to `~/mynet`
- Confirmed that the `code-quality` plugin exists with five agents
- Identified the `code-reviewer.md` agent file

If `ls ~/mynet/code-quality/agents/` does not show five `.md` files, re-run the clone command from Step 1.

## Step 4: Install the Plugin in Your Project

Navigate to your project directory (any project with source code will work). Create or edit your project's `.claude/settings.json` to add the plugin path:

```bash
mkdir -p .claude
cat > .claude/settings.json << 'EOF'
{
  "plugins": [
    "~/mynet/code-quality"
  ]
}
EOF
```

This tells Claude Code to load the code-quality plugin whenever you work in this project. The agents defined in the plugin become available in your Claude Code sessions.

## Step 5: Create a File with a Bug

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

## Step 6: Stage the Change

Add the file to git so it appears in `git diff`:

```bash
git add example.py
```

The code-reviewer agent works by examining `git diff` output. Staging the file makes it visible to the agent.

### Checkpoint

At this point you should have:

- A `.claude/settings.json` file pointing to the code-quality plugin
- An `example.py` file staged in git with two known issues
- Run `git diff --cached` to verify the staged file appears in the diff

## Step 7: Launch Claude Code and Trigger the Agent

Start Claude Code in your project directory:

```bash
claude
```

Once Claude Code starts, type this prompt:

```
Review the changes I just made
```

Claude Code recognizes that you are asking for a code review and routes your request to the **code-reviewer** agent from the code-quality plugin.

## Step 8: Read the Agent's Output

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

If Claude Code did not route to the code-reviewer agent, verify that `.claude/settings.json` contains the correct plugin path and restart Claude Code.

## Step 9: Fix the Issues and Re-review

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

- **Cloned the Mynet marketplace** and explored its structure -- a collection of independent plugins, each with its own agents, skills, and commands
- **Installed a plugin** by adding its path to your project's `.claude/settings.json`
- **Triggered an agent** by asking Claude Code a natural language question that matched the agent's description
- **Interpreted structured output** from the code-reviewer agent, understanding its severity categories and fix recommendations
- **Iterated on feedback** by fixing issues and re-running the review

## Next Steps

Now that you have a working plugin installation, explore further:

- **Try other agents in code-quality** -- the `test-writer-fixer` agent generates tests for your code, and the `architect-review` agent evaluates your project's architecture
- **Browse the marketplace** -- look at `web-development`, `backend-development`, `devops-and-infra`, and other plugins for agents that match your workflow
- **Install multiple plugins** -- add more entries to the `plugins` array in `.claude/settings.json` to combine capabilities from different plugins

## Troubleshooting

**Claude Code does not recognize the plugin**

Verify the path in `.claude/settings.json` is correct and points to the plugin directory (the one containing `.claude-plugin/plugin.json`). Restart Claude Code after changing settings.

**The code-reviewer agent does not trigger**

The agent activates on prompts related to code review. Use phrases like "review the changes I just made," "check my code before I push," or "does this look right to you?" If `git diff` shows no changes, the agent has nothing to review -- make sure your file is staged.

**git diff shows no output**

Run `git diff --cached` to check staged changes. If you forgot to run `git add`, the changes are unstaged and `git diff` (without `--cached`) shows them. The code-reviewer agent checks both staged and unstaged changes, but the file must be tracked by git.

**The review output is empty or says "No issues found"**

The agent only reports real problems. If your code is clean, a "No issues found" response is correct. To verify the agent is working, introduce a deliberate bug (like the division-by-zero example above) and re-run the review.
