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
| Trigger phrases | "modernize an agent", "audit agent definitions", "update agent format", "check agent quality", "rewrite agent prompts", "check agent routing" |

Also applies when reviewing agents for verbose content, missing frontmatter, weak triggering descriptions, or example blocks that inflate routing context.

### Files

| File | Contents |
|------|----------|
| `SKILL.md` | Process: run the audit script, confirm semantic findings by reading, report, rewrite, validate |
| `references/audit-criteria.md` | Frontmatter checklist, deterministic measurements, semantic judgment table, anti-pattern before/after examples, quality criteria, rewriting methodology, manual-answering guide, color guidelines |
| `references/question-catalog.json` | The semantic judgment questions, criteria, and thresholds the script sends to Jev; the rewrite-check questions used by `--compare` |
| `scripts/audit-agents.py` | Deterministic checks, Jev judgment calls, routing check, rewrite comparison |

### Script Usage

```
python3 scripts/audit-agents.py PATH [PATH ...]      # audit files and/or agent directories
python3 scripts/audit-agents.py --route DIR           # example self-routing check
python3 scripts/audit-agents.py --compare OLD NEW     # rewrite regression check
python3 scripts/audit-agents.py --no-jev PATH         # deterministic checks only
python3 scripts/audit-agents.py --json PATH           # emit machine-readable results
```

Paths are relative to the skill's directory.

| Flag | Effect |
|------|--------|
| (none) | Audit one agent file, or every `.md` file in a given directory |
| `--route DIR` | Check whether each agent's own `<example>` `user:` requests route to that agent (needs the key) |
| `--compare OLD NEW` | Run the rewrite-check questions against an original/rewritten file pair (needs the key) |
| `--no-jev` | Run deterministic checks only, skipping Jev entirely |
| `--json` | Emit results as JSON instead of the markdown findings tables |

**Exit codes:** `0` no finding is Must fix; `1` at least one finding is Must fix, or `--route`/`--compare` found a failure; `2` `TYPESAFE_API_KEY` is required for the requested operation and is not set.

### Environment

| Variable | Purpose |
|----------|---------|
| `TYPESAFE_API_KEY` | Bearer token sent to the TypeSafe Jev endpoint (`https://api.typesafe.ai/v1/systemone`). Required for semantic judgments, `--route`, and `--compare`. Deterministic checks run without it; the script prints a note to stderr and continues. |

Jev model: `jev-latest`.

State sent per agent: `{agent: {name, model, description, examples, body}, known_agents}`. `known_agents` is every agent name under `*/agents/*.md` across the marketplace when a `.claude-plugin/marketplace.json` exists two directories up from the agent; otherwise every `*.md` file in the agent's own plugin directory.

### Deterministic Checks

Body limits exclude example blocks; the body is the file content after frontmatter with `<example>...</example>` blocks removed.

**Frontmatter fields:**

| Field | Required | Valid values | Notes |
|-------|----------|-------------|-------|
| `name` | Yes | lowercase, hyphens, 3-50 chars, alphanumeric at both ends (`^[a-z0-9][a-z0-9-]{1,48}[a-z0-9]$`) | Missing or malformed is Must fix; a name that does not match the filename is Consider |
| `description` | Yes | 10-5,000 chars, no literal `\n` | Missing, under 10 chars, over 5,000 chars, or a literal `\n` is Must fix; over 120 words is Should fix |
| `model` | Yes | `inherit`, `sonnet`, `opus`, `haiku` | Missing is Must fix; a raw model ID (matches `claude-`, an 8-digit date, or an `N.N` version) or any other unknown alias is Must fix |
| `color` | Yes | `blue`, `cyan`, `green`, `yellow`, `magenta`, `red` | Missing or invalid is Must fix; a color shared by more than one agent in the same plugin directory is Consider |
| `tools` | Recommended | array of tool names | Missing is Consider (the agent gets every tool); present but not an array is Must fix |

**Examples:**

| Check | Condition | Severity |
|-------|-----------|----------|
| None present | Zero `<example>` blocks anywhere in the file | Must fix |
| In description | One or more `<example>` blocks appear before the frontmatter closes | Should fix |
| Too few | Fewer than 2 `<example>` blocks | Should fix |
| Too many | More than 5 `<example>` blocks | Consider |
| Incomplete | An example is missing `Context:`, `user:`, `assistant:`, or `<commentary>` | Must fix |

**Body:**

| Check | Condition | Severity |
|-------|-----------|----------|
| Empty | 0 characters | Must fix |
| Hard max | Over 5,000 characters | Should fix |
| Ideal max | Over 3,000 characters (and at or under 5,000) | Should fix |
| Ideal min | Under 500 characters | Consider |
| Bullets | Over 20 lines starting with `-`, `*`, or `+` (numbered steps counted separately, not limited) | Should fix |
| First person | Contains `I am`, `I will`, `I would`, or `I have` | Consider |
| Duplicate headings | The same numbered heading (for example `## 2.`) appears more than once | Consider |

### Semantic Judgment Catalog

Each row is one question in `question-catalog.json`. Noul questions return a probability (`p`); score questions return a probability-weighted score from 0 to 3 (fractional) plus a confidence.

| ID | Type | Detects | Flag rule |
|----|------|---------|-----------|
| `topic_lists` | Noul | Lists that name concepts without a decision | p >= 0.7 Should fix; 0.4 <= p < 0.7 reported as `Possible:`, Consider |
| `teaches_known` | Noul | Inventories of stdlib, APIs, or well-known practice with no project rule | p >= 0.7 Should fix; 0.4 <= p < 0.7 Consider |
| `fictional_content` | Noul | Fabricated metrics, fake progress objects, delivery scripts | p >= 0.7 Should fix; 0.4 <= p < 0.7 Consider |
| `phantom_refs` | Noul | Named collaborators or protocols absent from `known_agents` | p >= 0.7 Should fix; 0.4 <= p < 0.7 Consider |
| `overspecified_output` | Noul | JSON envelopes or protocols with no evidence of implementation | p >= 0.7 Should fix; 0.4 <= p < 0.7 Consider |
| `duplicate_coverage` | Noul | Same subject under two sections | p >= 0.7 Should fix; 0.4 <= p < 0.7 Consider |
| `mirrored_do_not` | Noul | A "Do Not" list that only inverts the positive list | p >= 0.7 Should fix; 0.4 <= p < 0.7 Consider |
| `distinct_domains` | Noul | Two unrelated domains in one agent; drives the Split action | p >= 0.7 (combined with body over 5,000 chars) |
| `has_role_statement` | Noul | Opening one-sentence role statement | p < 0.3 Should fix; 0.3 <= p < 0.6 Consider |
| `has_boundaries` | Noul | Explicit limits or stop conditions | p < 0.3 Should fix; 0.3 <= p < 0.6 Consider |
| `has_process` | Noul | Numbered action steps | p < 0.3 Should fix; 0.3 <= p < 0.6 Consider |
| `has_output_format` | Noul | Defined deliverable or report shape | p < 0.3 Should fix; 0.3 <= p < 0.6 Consider |
| `description_triggering` | Score 0-3 | Summary (0) to concrete triggers with a sibling handoff (3) | Score < 1.5 Should fix; 1.5 <= score < 2.5 Consider |
| `guidance_density` | Score 0-3 | Topic inventories (0) to decisions, defaults, and boundaries throughout (3) | Score < 1.5 Should fix |

Thresholds: `noul_flag = 0.7`, `noul_review = 0.4`, `routing_confidence = 0.7`.

### Severity Levels

| Severity | Deterministic triggers | Semantic triggers |
|----------|------------------------|-------------------|
| Must fix | Missing required field, bad name format, raw model ID, invalid color, literal `\n` in description, no examples anywhere, example missing a field | none |
| Should fix | Examples in description, description over 120 words, body over 3,000 chars, over 20 bullets | Anti-pattern Noul at or above 0.7, missing role/boundaries/process/output at or below 0.3, `description_triggering` below 1.5, `guidance_density` below 1.5 |
| Consider | No `tools` array, shared color in plugin, body under 500 chars, first person | Any Noul between 0.4 and 0.7 (`0.3`-`0.6` for role/boundaries/process/output), `description_triggering` below 2.5 |

### Action Rules

Computed per agent from its findings, metrics, and judgments. Rules 1, 2-4, and 5 are evaluated independently and their results are joined; only 2-4 are mutually exclusive.

1. Any Must fix finding: `Fix frontmatter` if every Must fix finding's area is something other than `body`; otherwise `Fix structure`.
2. If the body exceeds 5,000 characters and `distinct_domains` is at or above 0.7: `Split`.
3. Otherwise, if two or more of `topic_lists`, `teaches_known`, `fictional_content`, `phantom_refs`, `duplicate_coverage`, `mirrored_do_not`, `overspecified_output` are at or above 0.7, or `guidance_density` is below 1.5, or bullets exceed 20, or the body exceeds 5,000 characters: `Rewrite`.
4. Otherwise, if exactly one of those anti-pattern judgments is at or above 0.7, or the body exceeds 3,000 characters: `Trim`.
5. Independently of the above, if any example block sits inside the description: append `Move examples to body`.
6. If none of the above apply: `Keep`.

Multiple parts join with `; ` (for example `Fix frontmatter; Rewrite; Move examples to body`).

### Output Formats

**Findings table** (per agent):

`| # | Area | Issue | Severity |`

**Summary table** (printed when more than one agent is audited):

`| Agent | Body chars | Must fix | Should fix | Consider | Action |`

**Routing output** (`--route DIR`): one line per `<example>` `user:` request found across the directory's agents, in the form `result winner conf owner example request` (request text truncated to 70 characters). `result` is `ok` or `FAIL` for a non-proactive example, or `info` for a proactive example (its `Context:` or `<commentary>` contains "proactive"), which is not counted. `winner` and `conf` are Jev's chosen agent and confidence for that request; `owner` is the agent the example belongs to. A trailing line reports how many non-proactive requests failed to route to their own agent at or above the `routing_confidence` threshold, and notes that `info` rows are not counted.

**Compare output** (`--compare OLD NEW`): one line per rewrite-check question, in the form `flag question-id p=value`. `flag` is `FAIL` at p >= 0.7, `check` at 0.4 <= p < 0.7, otherwise `ok`. Exit is 1 if any question is `FAIL`.

### Conventions Enforced

- `name`, `description`, `model` (alias, never a raw ID), and `color` are required; `tools` is recommended
- `<example>` blocks live in the body, after the frontmatter and before the role statement; examples left in the description cost routing tokens on every turn
- The description stays under 120 words and states triggers, plus a "do not use for" handoff when sibling agents overlap
- Two to five examples, each with `Context:`, `user:`, `assistant:`, and `<commentary>`
- Body target 500-3,000 characters excluding examples; 5,000 is the rewrite-or-split line

### Model Selection Guide

| Model | Use case |
|-------|----------|
| `opus` | Judgment-heavy work: architecture review, test diagnosis, design |
| `sonnet` | Formulaic or playbook-driven work: checklists, accessibility audits, test generation, a fixed prototyping cycle, even when the domain is broad |
| `haiku` | Format and lint checks |
| `inherit` | When the parent's model is always right |

**Merge** two agents that share more than 70% of their prompt or need the same context; add a section for the secondary concern. **Split** an agent whose body exceeds 5,000 characters across genuinely distinct domains.

### Rewrite-Check Questions

Used by `--compare ORIGINAL REWRITTEN`, sent as `{original: full text, rewritten: full text}`. All three are Noul questions flagged `FAIL` at p >= 0.7 (the same `noul_flag` threshold).

| ID | Detects |
|----|---------|
| `drops_domain_knowledge` | `rewritten` omits a project-specific convention, platform gotcha, hard-won workaround, or concrete rule present in `original` that a model could not infer from general knowledge |
| `changes_role` | `rewritten` describes an agent with a different responsibility or scope than `original`, rather than the same agent stated more concisely |
| `adds_unsupported_claims` | `rewritten` introduces rules, defaults, or facts absent from `original` that are not standard, widely documented practice |

A `FAIL` on any of these blocks replacing the original until addressed.

## Related

- [Agent Reference](../../reference/agents/) -- ai-engineer specification
- [Getting Started](../../tutorials/getting-started/) -- Tutorial walking through a first audit
- [Audit Agent Definitions](../../howto/audit-agent-definitions/) -- Task-focused audit guide
