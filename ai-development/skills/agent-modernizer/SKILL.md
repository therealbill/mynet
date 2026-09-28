---
name: agent-modernizer
description: >
  This skill should be used when the user asks to "modernize an agent", "audit agent definitions",
  "update agent format", "check agent quality", "rewrite agent prompts", "check agent routing", or
  wants to bring agent files up to current Claude Code plugin standards. Also applies when reviewing
  agents for verbose content, missing frontmatter, weak triggering descriptions, or example blocks
  that inflate routing context.
---

# Agent Modernizer

Audit and rewrite Claude Code plugin agent definitions. Code owns every check code can make: field
presence, name format, model alias, color, example structure and placement, lengths, and bullet counts.
The remaining judgments (topic lists, phantom references, guidance density, triggering strength) are
typed questions with explicit criteria in `references/question-catalog.json`. TypeSafe Jev answers
them when `TYPESAFE_API_KEY` is set; you answer them otherwise. Severity and action come from the rule
table below, not from impression.

## When to Use

- Auditing one agent `.md` file or a whole `agents/` directory
- Checking that each agent's own examples route to it and not to a sibling
- Rewriting agents with verbose prompts, missing frontmatter, or summary-style descriptions
- Regression-checking a rewrite against the original before replacing it

## Conventions Enforced

- `name`, `description`, `model` (alias, never a raw ID), and `color` are required; `tools` is recommended
- `<example>` blocks live in the **body**, after the frontmatter and before the role statement. Every
  description is loaded for routing on every turn; examples there cost hundreds of tokens per agent
- The description stays under 120 words and states triggers, plus a "do not use for" handoff when
  sibling agents overlap
- Two to five examples, each with `Context:`, `user:`, `assistant:`, and `<commentary>`
- Body target 500-3,000 characters excluding examples; 5,000 is the rewrite-or-split line

## Process

### 1. Run the audit script

```bash
python3 scripts/audit-agents.py path/to/agent.md          # one agent
python3 scripts/audit-agents.py plugin/agents/            # batch, with summary table
python3 scripts/audit-agents.py --route plugin/agents/    # example self-routing check
python3 scripts/audit-agents.py --no-jev path/to/agent.md # deterministic checks only
```

Paths are relative to this skill's directory. The script needs only Python 3; it calls Jev when the
key is set and says so when it is not. Use its numbers for characters, words, bullets, and example
counts. Do not estimate them by eye.

### 2. Confirm semantic findings by reading

Jev output is a signal, not a verdict. For every semantic finding, whether it shows a Noul `p=` or a
Score such as `1.4/3`, and for every `Possible:` line, open the file and cite the lines that decide it. Drop findings the text does
not support. A Noul near 0.5 means the question was hard to call, not that the problem is moderate.

Without a key, answer `references/question-catalog.json` yourself: one question at a time, against the
same state the script would send (name, description, examples, body, known agent names), recording
yes, no, or uncertain with the deciding lines. Do not skip questions or merge them.

### 3. Report

Per agent, a findings table sorted by severity, then the action:

| # | Area | Issue | Severity |
|---|------|-------|----------|
| 1 | `description` | 225 words; keep under 120 | Should fix |
| 2 | `examples` | 3 blocks inside `description`; move to body | Should fix |
| 3 | `body` | No defined output format (p=0.03) | Should fix |

For a directory, add the script's summary table and the routing results. A routing failure on a
proactive example is informational; on a plain request it means the description or a sibling's
description needs a handoff sentence.

### 4. Rewrite

**Trust the model.** State priorities and boundaries, not concept inventories. Opus does not need to
be told what `gofmt` is.

**Decisions over topics.** Replace "Branching strategies: Git Flow, GitHub Flow, trunk-based" with
"Default to trunk-based with short-lived feature branches unless the project ships supported releases."

**Guard rails over checklists.** A "Do Not" section naming the three to five mistakes that matter beats
a twenty-item list of things to do, and must not just invert the positive instructions.

**One-sentence role, numbered actions, defined output.** Each process step names an action. An
`**Output:**` line says what the agent hands back and in what shape.

**Keep what the model cannot infer.** Project conventions, platform gotchas, version-specific
workarounds, and numeric limits survive the rewrite, reworded if needed.

### 5. Validate

```bash
python3 scripts/audit-agents.py rewritten.md
python3 scripts/audit-agents.py --compare original.md rewritten.md   # needs the key
python3 scripts/audit-agents.py --route plugin/agents/
```

`--compare` asks whether the rewrite drops non-inferable knowledge, changes the role, or adds claims
the original never made. A FAIL on any of these blocks replacing the original until addressed.

## Severity and Action Rules

| Severity | Deterministic triggers | Semantic triggers |
|----------|------------------------|-------------------|
| Must fix | Missing required field, bad name format, raw model ID, invalid color, literal `\n` in description, no examples anywhere, example missing a field | none |
| Should fix | Examples in description, description over 120 words, body over 3,000 chars, over 20 bullets | Anti-pattern Noul at or above 0.7, missing role, boundaries, process, or output at or below 0.3, triggering score below 1.5, guidance density below 1.5 |
| Consider | No `tools` array, shared color in plugin, body under 500 chars, first person | Any Noul between 0.4 and 0.7 (confirm by reading), triggering score below 2.5 |

"Anti-pattern" in the action rules means one of the seven catalog questions `topic_lists`,
`teaches_known`, `fictional_content`, `phantom_refs`, `overspecified_output`, `duplicate_coverage`,
and `mirrored_do_not`. A missing role, boundaries, process, or output is a completeness gap that
produces a finding but does not by itself drive Rewrite or Trim.

Action: **Fix frontmatter** when any Must fix exists; **Split** when the body exceeds 5,000 chars and
covers distinct domains; **Rewrite** for two or more confirmed anti-patterns, density below 1.5, over
20 bullets, or a body over 5,000 chars; **Trim** for one anti-pattern or a body over 3,000 chars;
**Move examples to body** whenever they sit in the description; otherwise **Keep**. Thresholds were
set by running the catalog over this marketplace's 64 agents; treat them as starting points when the
target repository differs.

## Key Decisions for Rewrites

**Model selection:** `opus` for judgment-heavy work (architecture review, test diagnosis, design);
`sonnet` for formulaic or playbook-driven work (checklists, accessibility audits, test generation, a
fixed prototyping cycle) even when the domain is broad; `haiku` for format and lint checks; `inherit`
when the parent's model is always right.

**Merge** two agents that share more than 70% of their prompt or need the same context; add a section
for the secondary concern. **Split** an agent whose body exceeds 5,000 characters across genuinely
distinct domains.

## Do Not

- Flag body-placed examples as "missing from description"; that placement is the standard here
- Report a Jev probability as a finding without reading the lines behind it
- Rewrite on the strength of the script alone; the script ranks work, the rewrite is yours
- Strip a fact because it looks like padding before checking it is inferable from general knowledge
- Add timeline or effort estimates to findings

## Reference Files

- `references/audit-criteria.md` — frontmatter checklist, anti-patterns with before/after examples,
  quality criteria, color guidelines, rewriting methodology, manual answering guide
- `references/question-catalog.json` — the exact judgments, criteria, and thresholds the script sends
- `scripts/audit-agents.py` — deterministic checks, Jev judgments, routing check, rewrite comparison
