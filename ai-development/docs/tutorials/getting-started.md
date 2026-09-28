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

Pick any agent definition from another plugin. For this tutorial, assume you have an agent file at `devops/agents/ci-engineer.md` with some common issues: a missing `color` field, example blocks sitting inside the `description` instead of the body (costing routing tokens on every turn), a summary-style description with no triggering phrases, and a verbose system prompt full of topic lists.

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

From `ai-development/skills/agent-modernizer/`, the skill runs the audit script against your file:

```bash
python3 scripts/audit-agents.py path/to/devops/agents/ci-engineer.md
```

Pointed at a directory instead of a single file, the same script audits every agent in it and adds a summary table; `--route plugin/agents/` additionally checks whether each agent's own examples route to it rather than to a sibling.

If `TYPESAFE_API_KEY` is set, the semantic judgments -- whether the prompt is a topic list, whether it defines an output format, and so on -- come from TypeSafe Jev and carry a probability. Without the key, the script still runs every deterministic check (fields, lengths, example placement), and Claude answers the semantic judgments by reading the file directly.

## Step 4: Review the Findings Table

The script prints a findings table sorted by severity. A findings table from a live run looks like this:

| # | Field/Area | Issue | Severity |
|---|-----------|-------|----------|
| 1 | `description` | 225 words; every agent description is loaded for routing on every turn, keep it under 120 | Should fix |
| 2 | `examples` | 3 example block(s) inside `description`; move them to the body | Should fix |
| 3 | `body` | No defined output format (p=0.03) | Should fix |
| 4 | `description` | No 'do not use for' handoff to sibling agents (score 2.0/3) | Consider |

Severity levels mean:

- **Must fix** -- The agent won't load or won't trigger correctly without this change
- **Should fix** -- The agent works but is suboptimal
- **Consider** -- A polish improvement, not urgent

## Step 5: Confirm a Finding by Reading

Before acting on the table, check one finding against the source. Open `devops/agents/ci-engineer.md` and look at the `description` field: the `<example>` blocks are sitting inside it, above the closing `---`, which is what row 2 flagged. This matters most for findings that carry a probability rather than a hard count -- a Jev judgment is a signal, not a verdict, and dropping a finding the text doesn't support is a normal part of the audit.

## Step 6: Request a Rewrite

When multiple issues are found, ask for a full rewrite rather than fixing items individually:

```
Rewrite this agent based on the audit findings
```

The agent-modernizer applies its rewriting principles: trust the model's knowledge, provide decisions instead of topic lists, add guard rails instead of checklists, and write a concise role statement.

## Step 7: Review the Rewritten Agent

Compare the before and after. The rewritten agent should have:

- **Complete frontmatter** -- All required fields present (`name`, `description`, `model`, `color`, `tools`)
- **Examples in the body** -- 2-5 `<example>` blocks after the frontmatter, each with `Context:`, `user:`, `assistant:`, and `<commentary>`
- **A tight, triggering description** -- Under 120 words, stating concrete requests plus a "do not use for" handoff to overlapping siblings
- **Concise system prompt** -- 500-3,000 characters excluding examples. States the role in one sentence, provides decisions and boundaries, includes a numbered process
- **No anti-patterns** -- No topic lists without guidance, no phantom references, no fictional metrics

## Step 8: Verify the Result

Confirm the rewrite before replacing the original:

- [ ] Re-run the script against the rewrite; it exits 0
- [ ] Run `python3 scripts/audit-agents.py --compare original.md rewritten.md`
- [ ] Examples live in the body, not the description
- [ ] `description` is under 120 words and includes a handoff to sibling agents
- [ ] Body is under 3,000 characters

`--compare` needs `TYPESAFE_API_KEY` and checks whether the rewrite dropped non-inferable knowledge, changed the role, or added claims the original never made. A FAIL on any of those blocks replacing the original until addressed.

## Summary

The agent-modernizer skill runs a deterministic-plus-semantic audit script against agent definitions and produces a findings table with actionable severity levels. When multiple issues exist, it can rewrite the entire agent to current standards, and `--compare` checks the rewrite against the original before you replace it.

## Next Steps

- [Audit Agent Definitions](../../howto/audit-agent-definitions/) -- Task-focused guide for bringing agents up to standard
- [Skill Reference](../../reference/skills/) -- Full specification for agent-modernizer
- [Architecture](../../explanation/architecture/) -- Why this plugin serves both AI builders and plugin maintainers
