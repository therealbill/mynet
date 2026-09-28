---
title: "Audit Agent Definitions"
description: "Modernize agent definitions for current plugin standards"
weight: 1
---

# Audit Agent Definitions

Bring an existing agent definition up to current Claude Code plugin standards using the agent-modernizer skill.

## Problem

An agent definition has stale formatting, missing frontmatter fields, verbose system prompts, or example blocks in the wrong place. It may load but trigger unreliably or provide unfocused assistance.

## Prerequisites

- The ai-development plugin installed
- Python 3 available on the machine running Claude Code (the audit script needs nothing else)
- Path to the agent `.md` file or agents directory you want to audit
- Optional: a `TYPESAFE_API_KEY` environment variable, if you want the harder judgments sent to TypeSafe Jev instead of answered manually

## Steps

### 1. Run the Audit Script

The skill starts with the deterministic script, run from the skill directory (`ai-development/skills/agent-modernizer/`), with paths relative to it:

```bash
python3 scripts/audit-agents.py path/to/agent.md          # one agent
python3 scripts/audit-agents.py plugin/agents/             # batch, with summary table
python3 scripts/audit-agents.py --route plugin/agents/     # example self-routing check
python3 scripts/audit-agents.py --compare old.md new.md    # regression check after a rewrite
python3 scripts/audit-agents.py --no-jev path/to/agent.md  # deterministic checks only
python3 scripts/audit-agents.py --json path/to/agent.md    # machine-readable output
```

The script checks everything code can check by itself: required fields, name format, model alias, color, example structure and placement, body and description lengths, and bullet counts.

If `TYPESAFE_API_KEY` is set, it also sends the harder judgments — topic lists, phantom agent references, guidance density, triggering strength, and more — to TypeSafe Jev and reports each as a probability. Without the key, it prints the deterministic findings only, and you answer `references/question-catalog.json` yourself, one question at a time, against the same state the script would have sent (name, description, examples, body, known agent names).

### 2. Confirm Semantic Findings by Reading

Treat Jev's output as a signal, not a verdict. For every semantic finding, whether it shows a Noul `p=` or a Score such as `1.4/3`, and for every `Possible:` line, open the file and check the lines that decide it. Drop findings the text does not support.

A probability near 0.5 means the question was hard to call — that is uncertainty, not a "moderate" version of the problem. Never report a probability as a finding without reading the lines behind it.

### 3. Review the Findings

The report is a table per agent, sorted by severity, with a computed action:

| # | Area | Issue | Severity |
|---|------|-------|----------|
| 1 | `description` | 225 words; keep under 120 | Should fix |
| 2 | `examples` | 3 blocks inside `description`; move to body | Should fix |
| 3 | `body` | No defined output format (p=0.03) | Should fix |

Focus on **Must fix** items first — these block the agent from loading or triggering correctly: a missing required frontmatter field, a bad name format, a raw model ID instead of an alias, an invalid color, literal `\n` escapes in the description, no `<example>` blocks anywhere in the file, or an example missing a required field.

Then address **Should fix** items: examples still sitting inside the description (they belong in the body, after the frontmatter and before the role statement — moving them is the fix, not deleting them), a description over 120 words, a body over 3,000 characters, more than 20 bullets, or a confirmed anti-pattern such as a topic list, phantom agent reference, or a missing role, boundaries, process, or output.

**Consider** items are polish: no `tools` array, a shared color within the plugin, a body under 500 characters, or first-person phrasing. The full severity and action rules live in the skill's `SKILL.md`; the script applies them for you.

For a directory audit, also check the results from `--route`. A routing failure on a proactive example is informational — a proactive example doesn't need to beat every sibling. A failure on a plain request example means the description, or a sibling's description, needs a "do not use for" handoff sentence.

### 4. Request a Rewrite

If multiple issues are found, request a full rewrite rather than patching individual fields:

```
Rewrite this agent based on the findings
```

The rewrite favors decisions over topic inventories, guard rails over checklists, and a one-sentence role with numbered process steps and a defined output. It keeps whatever the model cannot infer on its own: project conventions, platform gotchas, version-specific workarounds, and numeric limits.

For single-field fixes (for example, just adding a missing `color`), apply the change directly instead.

### 5. Review the Rewritten Agent

Verify the rewritten agent includes:

- **Frontmatter** — `name`, `model` (alias, not a raw ID), `color`, and `tools` scoped to what the agent needs
- **A description under 120 words** written as selection criteria: concrete requests or situations a user would type, plus a "do not use for (use sibling-agent)" handoff where a sibling agent overlaps
- **Example blocks in the body**, after the frontmatter and before the role statement — two to five of them, each with `Context:`, `user:`, `assistant:`, and `<commentary>`
- **A concise system prompt** — one-sentence role, decisions instead of topic lists, numbered process steps, and a focused "Do Not" section that doesn't just invert the positive instructions
- **No anti-patterns** — no fictional metrics, no phantom agent references, no concept lists the model already knows

### 6. Verify the Agent Loads

After applying the rewritten definition:

```bash
python3 scripts/audit-agents.py rewritten.md
python3 scripts/audit-agents.py --compare original.md rewritten.md   # needs TYPESAFE_API_KEY
python3 scripts/audit-agents.py --route plugin/agents/
```

Re-running the script against the rewritten file should show no Must-fix findings and a clean example count. `--compare` asks whether the rewrite dropped non-inferable knowledge, changed the agent's role, or added claims the original never made — a FAIL on any of these blocks replacing the original until addressed. Finally, confirm that trigger phrases from the `<example>` blocks activate the agent as expected.

## Batch Audits

To audit all agents in a plugin at once, point the script at the directory:

```bash
python3 scripts/audit-agents.py devops/agents/
```

The script produces a summary table with per-severity counts and a computed action per agent:

| Agent | Body chars | Must fix | Should fix | Consider | Action |
|-------|-----------:|---------:|-----------:|---------:|--------|
| `ci-engineer` | 4120 | 0 | 2 | 1 | Rewrite |
| `deploy-manager` | 1870 | 0 | 0 | 1 | Keep |

Fix any Must-fix items first, then work through the rewrites in severity order.

## Related

- [Getting Started](../../tutorials/getting-started/) -- Step-by-step tutorial walking through a first audit
- [Skill Reference](../../reference/skills/) -- agent-modernizer specification and audit criteria
- research plugin -- Use for gathering domain knowledge before writing agent prompts from scratch
