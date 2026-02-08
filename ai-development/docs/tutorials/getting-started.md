---
title: "Getting Started with AI Development"
description: "Use the agent-modernizer skill to audit an agent definition"
weight: 1
---

# Getting Started with AI Development

Audit an existing agent definition using the agent-modernizer skill, then apply the recommended rewrites to bring it up to current standards.

## What You'll Learn

- How to invoke the agent-modernizer skill
- How to read an audit findings table
- How to apply a rewrite and verify the result

## Prerequisites

- The ai-development plugin installed in your Claude Code environment
- An existing agent `.md` file to audit (from any plugin in the marketplace)

## Step 1: Choose an Agent to Audit

Pick any agent definition from another plugin. For this tutorial, assume you have an agent file at `devops/agents/ci-engineer.md` with some common issues: missing `color` field, no `<example>` blocks in the description, and a verbose system prompt full of topic lists.

Open the file or have its path ready:

```
devops/agents/ci-engineer.md
```

## Step 2: Request an Audit

Ask Claude to audit the agent definition. Any of these phrases will activate the agent-modernizer skill:

```
Audit this agent definition for quality
```

```
Check agent quality for devops/agents/ci-engineer.md
```

```
Modernize this agent
```

The skill activates automatically when it detects these trigger phrases. You do not need to reference the skill by name.

## Step 3: Observe the Audit Process

The agent-modernizer reads the agent file and evaluates it against the plugin agent specification. It checks every frontmatter field and assesses the system prompt for anti-patterns. You will see it work through:

1. **Frontmatter validation** -- Checks for required fields: `name`, `description`, `model`, `color`, `tools`
2. **Description assessment** -- Verifies `<example>` blocks exist with proper structure (Context, user, assistant, commentary)
3. **System prompt analysis** -- Scans for anti-patterns like topic lists without guidance, teaching the model its own knowledge, fictional progress tracking, and phantom agent references

## Step 4: Review the Findings Table

The audit produces a structured findings table with severity levels:

| # | Field/Area | Issue | Severity |
|---|-----------|-------|----------|
| 1 | `color` | Missing -- required field | Must fix |
| 2 | `description` | No `<example>` blocks -- agent won't trigger reliably | Must fix |
| 3 | System prompt | 34 bullet points listing concepts without decisions | Should fix |
| 4 | System prompt | References `security-auditor` agent that does not exist | Should fix |
| 5 | `tools` | Includes `WebSearch` but agent never needs web search | Consider |

Severity levels mean:

- **Must fix** -- The agent won't load or won't trigger correctly without this change
- **Should fix** -- The agent works but is suboptimal
- **Consider** -- A polish improvement, not urgent

## Step 5: Request a Rewrite

When multiple issues are found, ask for a full rewrite rather than fixing items individually:

```
Rewrite this agent based on the audit findings
```

The agent-modernizer applies its rewriting principles: trust the model's knowledge, provide decisions instead of topic lists, add guard rails instead of checklists, and write a concise role statement.

## Step 6: Review the Rewritten Agent

Compare the before and after. The rewritten agent should have:

- **Complete frontmatter** -- All required fields present (`name`, `description` with 2-4 `<example>` blocks, `model`, `color`, `tools`)
- **Concise system prompt** -- Under 3,000 characters, ideally 500-2,000. States the role in one sentence, provides decisions and boundaries, includes a numbered process
- **No anti-patterns** -- No topic lists without guidance, no phantom references, no fictional metrics

## Step 7: Verify the Result

Confirm the rewritten agent meets the specification:

- [ ] All required frontmatter fields are present
- [ ] Description includes `<example>` blocks with Context, user, assistant, and commentary
- [ ] `tools` array follows the principle of least privilege
- [ ] System prompt provides decisions, not concept inventories
- [ ] Domain-specific knowledge is preserved (only padding was removed)

## Summary

The agent-modernizer skill audits agent definitions against a structured checklist and produces a findings table with actionable severity levels. When multiple issues exist, it can rewrite the entire agent to current standards while preserving domain-specific knowledge that the model cannot infer on its own.

## Next Steps

- [Audit Agent Definitions]({{< ref "/howto/audit-agent-definitions" >}}) -- Task-focused guide for bringing agents up to standard
- [Skill Reference]({{< ref "/reference/skills" >}}) -- Full specification for agent-modernizer
- [Architecture]({{< ref "/explanation/architecture" >}}) -- Why this plugin serves both AI builders and plugin maintainers
