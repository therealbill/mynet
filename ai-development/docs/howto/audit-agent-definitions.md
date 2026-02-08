---
title: "Audit Agent Definitions"
description: "Modernize agent definitions for current plugin standards"
weight: 1
---

# Audit Agent Definitions

Bring an existing agent definition up to current Claude Code plugin standards using the agent-modernizer skill.

## Problem

An agent definition has stale formatting, missing frontmatter fields, verbose system prompts, or lacks proper `<example>` blocks. It may load but triggers unreliably or provides unfocused assistance.

## Prerequisites

- The ai-development plugin installed
- Path to the agent `.md` file you want to modernize

## Steps

### 1. Open or Reference the Agent File

Identify the agent `.md` file to audit. You can either open it directly or reference it by path:

```
devops/agents/ci-engineer.md
```

If auditing multiple agents in a directory, reference the directory instead. The skill supports batch audits.

### 2. Request the Audit

Use any of these phrases to activate the agent-modernizer skill:

```
Modernize this agent
```

```
Audit this agent definition
```

```
Check agent quality
```

```
Update agent format
```

The skill reads the file, evaluates every frontmatter field against the spec, and scans the system prompt for anti-patterns.

### 3. Review the Findings Table

Focus on **Must fix** items first. These prevent the agent from loading or triggering correctly:

- Missing required frontmatter fields (`name`, `description`, `model`, `color`)
- No `<example>` blocks in the description
- Raw model IDs (e.g., `claude-sonnet-4-20250514`) instead of aliases (`sonnet`)

Then address **Should fix** items:

- Topic lists without decisions or guidance
- Phantom references to agents that do not exist
- Redundant sections (a "Don't" list that inverts the "Do" list)
- Teaching the model information it already knows

**Consider** items are polish. Address them if you are already rewriting.

### 4. Request a Rewrite

If multiple issues are found, request a full rewrite rather than patching individual fields:

```
Rewrite this agent based on the findings
```

For single-field fixes (e.g., just adding a missing `color`), apply the change directly.

### 5. Review the Rewritten Agent

Verify the rewritten agent includes:

- **Frontmatter fields** -- `name`, `description` (with 2-4 `<example>` blocks), `model` (alias, not raw ID), `color`, `tools` (least-privilege array)
- **Concise system prompt** -- One-sentence role statement, decisions and boundaries instead of topic inventories, numbered process steps, a focused "Do Not" section
- **No anti-patterns** -- No fictional metrics, no phantom agent references, no concept lists the model already knows

### 6. Verify the Agent Loads

After applying the rewritten definition:

1. Confirm the file parses correctly (valid YAML frontmatter, no escaped newline literals in the description)
2. Test that trigger phrases from the `<example>` blocks activate the agent
3. Verify the `tools` array includes only tools the agent genuinely needs

## Batch Audits

To audit all agents in a plugin at once, reference the agents directory:

```
Audit all agents in devops/agents/
```

The skill produces a summary table:

| Agent | Lines | Missing Fields | Anti-patterns | Action Needed |
|-------|-------|---------------|---------------|---------------|
| `ci-engineer` | 47 | `color`, examples | Topic lists | Rewrite |
| `deploy-manager` | 28 | None | Phantom refs | Fix prompts |

Prioritize by severity. Fix frontmatter-only issues first, then full rewrites.

## Related

- [Getting Started](../../tutorials/getting-started/) -- Step-by-step tutorial walking through a first audit
- [Skill Reference](../../reference/skills/) -- agent-modernizer specification and audit criteria
- research plugin -- Use for gathering domain knowledge before writing agent prompts from scratch
