---
title: "Skills"
description: "AI development skill specifications"
weight: 2
---

# Skills

Specialized skills provided by the ai-development plugin.

## agent-modernizer

| Field | Value |
|-------|-------|
| Name | agent-modernizer |
| Trigger phrases | "modernize an agent", "audit agent definitions", "update agent format", "check agent quality", "rewrite agent prompts" |

**Process:** Reads agent `.md` file, evaluates frontmatter fields, assesses system prompt quality, produces findings table with severity levels, optionally rewrites the agent.

### Audit Criteria

**Frontmatter fields checked:**

| Field | Required | Valid Values | Notes |
|-------|----------|-------------|-------|
| `name` | Yes | lowercase, hyphens, 3-50 chars | Must start and end alphanumeric |
| `description` | Yes | 10-5,000 chars with `<example>` blocks | Controls agent triggering |
| `model` | Yes | `inherit`, `sonnet`, `opus`, `haiku` | Raw model IDs are invalid |
| `color` | Yes | `blue`, `cyan`, `green`, `yellow`, `magenta`, `red` | Distinct per plugin |
| `tools` | Recommended | Array of tool names | Principle of least privilege |

**Description requirements:**

- Summary of when to use the agent
- 2-4 `<example>` blocks, each containing: `Context:`, `user:`, `assistant:`, `<commentary>`

### System Prompt Anti-Patterns

| Anti-Pattern | Description |
|-------------|-------------|
| Topic lists without guidance | Two-word bullets naming concepts without decisions |
| Teaching model knowledge | Listing stdlib functions or well-known patterns the model already knows |
| Fictional content | Fake JSON progress trackers, fabricated metrics |
| Phantom references | Mentions of agents or systems that do not exist |
| Redundancy | Same topic covered multiple times; "Don't" section inverting the "Do" section |
| Over-specification | Detailed JSON schemas for protocols that are not implemented |

### Severity Levels

| Level | Meaning | Examples |
|-------|---------|---------|
| Must fix | Agent won't load or trigger correctly | Missing `color`, no `<example>` blocks, raw model ID |
| Should fix | Agent works but is suboptimal | Verbose system prompt, phantom references, topic lists |
| Consider | Polish improvement | Tool array includes unnecessary tools, minor redundancy |

### Rewriting Principles

- **Trust the model** -- State priorities and boundaries, not concept inventories
- **Decisions over topics** -- Specify what to prefer and when to choose which approach
- **Guard rails over checklists** -- A "Do Not" section preventing common mistakes outweighs a long checklist
- **Concise role statement** -- One sentence establishing who the agent is
- **Concrete process** -- Numbered steps describing actions, not topics
- **Target length** -- System prompt under 3,000 characters, ideally 500-2,000

### Batch Audit Output

When auditing a directory of agents, produces a summary table:

| Agent | Lines | Missing Fields | Anti-patterns | Action Needed |
|-------|-------|---------------|---------------|---------------|
| `code-reviewer` | 31 | `color`, examples | Generic checklist | Rewrite |
| `architect-reviewer` | 43 | `model`, `color`, `tools` | None significant | Fix frontmatter |

### Model Selection Guide

| Model | Use Case |
|-------|----------|
| `opus` | Complex judgment: architecture review, test diagnosis, code simplification |
| `sonnet` | Formulaic review: code review checklists, accessibility audits |
| `haiku` | Simple validation: format checks, linting |
| `inherit` | When the parent's model is always appropriate |

## Related

- [Agent Reference](../../reference/agents/) -- ai-engineer specification
- [Getting Started](../../tutorials/getting-started/) -- Tutorial walking through a first audit
- [Audit Agent Definitions](../../howto/audit-agent-definitions/) -- Task-focused audit guide
