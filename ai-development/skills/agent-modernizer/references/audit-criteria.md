# Agent Audit Criteria

## Frontmatter Requirements

Every agent must have all required fields. The script checks each; this table is the contract.

| Field | Required | Valid Values | Notes |
|-------|----------|-------------|-------|
| `name` | Yes | lowercase, hyphens, 3-50 chars | Must start and end alphanumeric; should match the filename |
| `description` | Yes | 10-5,000 chars, under 120 words | What the router reads; no `<example>` blocks here |
| `model` | Yes | `inherit`, `sonnet`, `opus`, `haiku` | Raw model IDs such as `claude-sonnet-4-20250514` are invalid |
| `color` | Yes | `blue`, `cyan`, `green`, `yellow`, `magenta`, `red` | Distinct per plugin |
| `tools` | Recommended | Array of tool names | Principle of least privilege; omitting it grants every tool |

### Description and Example Placement

The description is loaded into the routing context for every agent on every turn. Examples in the
description cost 250-480 tokens per agent; moved to the body they cost nothing until the agent is
invoked. The standard layout is:

```markdown
---
name: cli-developer
description: >
  Use when the user asks to build, restructure, or polish a command-line tool: command trees,
  flags, config layering, help text, shell completions. Do not use for interactive terminal
  screens (use go-tui-developer) or palette and symbol decisions (use cli-ui-designer).
model: opus
color: blue
tools: ["Read", "Write", "Edit", "Bash"]
---

<example>
Context: User is starting a new Go CLI tool from scratch
user: "I need to build a CLI for managing database migrations"
assistant: "I'll use the cli-developer agent to design the command tree, flag conventions, and output formats."
<commentary>
New CLI requiring decisions about command structure, configuration layering, and output conventions.
</commentary>
</example>

You are a senior CLI developer who ...
```

The description must therefore carry the triggering conditions itself: concrete requests or
situations a user would type, plus a "do not use for" handoff when a sibling agent overlaps. A
description that only summarizes the agent ("Expert in Playwright testing for modern web
applications") scores low on triggering strength because it gives the router nothing to match.

**Common failures:**

- No examples anywhere: the agent triggers unreliably (Must fix)
- Examples still inside the description (Should fix; move them, do not delete them)
- Examples escaped as `\n` literals instead of real newlines, so the parser cannot read them (Must fix)
- An example missing `Context:`, `user:`, `assistant:`, or `<commentary>` (Must fix)
- A proactive example whose commentary does not say so; mark it `Proactive trigger:` so the routing
  check treats it as informational

## Deterministic Measurements

The script reports these; use its numbers rather than estimating.

| Measurement | Where | Limits |
|-------------|-------|--------|
| Body characters | Everything after the frontmatter, example blocks removed | 500-3,000 target; over 3,000 Should fix; over 5,000 rewrite or split |
| Bullet points | Lines starting with `-`, `*`, or `+` in the body | Over 20 is Should fix; numbered steps are counted separately and not limited |
| Description words | The folded `description` value | Over 120 is Should fix |
| Example count | `<example>` blocks anywhere in the file | 2-5; fewer is Should fix, more is Consider |
| First person | `I am`, `I will`, `I would`, `I have` in the body | Consider; prompts address the agent as "You" |
| Duplicate section numbers | Headings such as `## 2.` appearing twice | Consider |

## Semantic Judgments

Each judgment is one question in `question-catalog.json` with its own criteria. The script sends them
all in one request per agent with the state `{agent: {name, model, description, examples, body},
known_agents}`; `known_agents` is every agent name in the marketplace, or in the plugin when there is
no marketplace manifest. Answering manually means answering each question against that same state.

| Question | Type | Detects | Flag rule |
|----------|------|---------|-----------|
| `topic_lists` | Noul | Lists that name concepts without a decision | p at or above 0.7 |
| `teaches_known` | Noul | Inventories of stdlib, APIs, or well-known practice with no project rule | p at or above 0.7 |
| `fictional_content` | Noul | Fabricated metrics, fake progress objects, delivery scripts | p at or above 0.7 |
| `phantom_refs` | Noul | Named collaborators or protocols absent from `known_agents` | p at or above 0.7 |
| `overspecified_output` | Noul | JSON envelopes or protocols with no evidence of implementation | p at or above 0.7 |
| `duplicate_coverage` | Noul | Same subject under two sections | p at or above 0.7 |
| `mirrored_do_not` | Noul | A "Do Not" list that only inverts the positive list | p at or above 0.7 |
| `distinct_domains` | Noul | Two unrelated domains in one agent; drives Split | p at or above 0.7 with body over 5,000 chars |
| `has_role_statement` | Noul | Opening one-sentence role | p at or below 0.3 is Should fix |
| `has_boundaries` | Noul | Explicit limits or stop conditions | p at or below 0.3 is Should fix |
| `has_process` | Noul | Numbered action steps | p at or below 0.3 is Should fix |
| `has_output_format` | Noul | Defined deliverable or report shape | p at or below 0.3 is Should fix |
| `description_triggering` | Score 0-3 | Summary (0) to concrete triggers with handoff (3) | Below 1.5 Should fix; below 2.5 Consider |
| `guidance_density` | Score 0-3 | Topic inventories (0) to decisions throughout (3) | Below 1.5 Should fix |

Any Noul between 0.4 and 0.7 is reported as `Possible:` and must be confirmed or dropped by reading.
Thresholds came from running the catalog over 64 marketplace agents; the split of the original
"redundant" question into `duplicate_coverage` and `mirrored_do_not` happened because the combined
question sat near 0.5 for most agents, which is the model saying the question was ambiguous.

### Answering the catalog manually

1. Load the agent the way the script does: description, examples, and body separated, example blocks
   removed from the body.
2. Take one question at a time. Read its `criteria.true` and `criteria.false` before scanning the body.
3. Record yes, no, or uncertain and the lines that decide it. Uncertain maps to the 0.4-0.7 band.
4. For the two Score questions, pick the level whose description matches; do not average.
5. Apply the severity table in `SKILL.md`. Do not invent findings the questions do not cover; add a
   question to the catalog instead.

## System Prompt Anti-Patterns

### 1. Topic Lists Without Guidance

```markdown
# BAD
Branching strategies:
- Git Flow implementation
- GitHub Flow setup
- Trunk-based development
```

```markdown
# GOOD
**Branching strategy** — Default to trunk-based with short-lived feature branches unless the
project ships multiple supported release lines.
```

The model already knows these concepts exist. The prompt should say *when to choose which* and *what
to prefer*.

### 2. Fictional Progress Tracking

```markdown
# BAD
Progress tracking:
{ "merge_conflicts_reduced": "67%", "team_satisfaction": "4.5/5" }
```

Not real data. Remove entirely.

### 3. Phantom Agent References

```markdown
# BAD
- Collaborate with devops-engineer on CI/CD
- Support release-manager on versioning
```

Remove unless those agents exist in the marketplace. Cross-plugin references to real agents are fine.

### 4. Redundant Sections

Signs: the same subject under two headings, a "Do Not" list that inverts the "Do" list item for item,
duplicate section numbers. A Defaults section and a Do Not section that mention the same technology
while stating different rules are not redundant.

### 5. Teaching the Model Its Own Knowledge

```markdown
# BAD for an Opus-class agent
- Use `errors.Is` / `errors.As` over string comparison
- Prefer `io.Reader`/`io.Writer` interfaces
```

```markdown
# GOOD
Apply Effective Go and Go Code Review Comments conventions. Favor clarity over brevity.
```

If the information is in the model's training data and is not project-specific, it is padding.

### 6. Over-Specification of Output

```json
{ "requesting_agent": "cli-developer", "request_type": "get_cli_context", "payload": {} }
```

Unless this protocol is implemented and consumed, remove it. A human-readable `**Output:**` line is
the right level of specification.

## System Prompt Quality Criteria

1. **States the role** in one sentence
2. **Provides decisions, not topics**: what to prefer, what to default to, when to choose which
3. **Defines boundaries**: what not to do, what to leave alone, where to stop
4. **Specifies process**: numbered steps naming actions
5. **Defines output**: how results are reported
6. **Handles edge cases** the model would not anticipate (generated code, vendor directories)
7. **Stays under 3,000 characters** excluding examples

## Rewriting Methodology

1. **Identify the core purpose**: strip padding to find what the agent actually does
2. **Choose the model**: Opus for judgment, Sonnet for formulaic procedure, Haiku for simple checks
3. **Write for that model**: Opus needs boundaries and priorities; Sonnet needs explicit procedure
4. **Keep domain knowledge the model lacks**: project conventions, platform gotchas such as the
   termenv/AppleScript macOS workaround, numeric limits
5. **Add guard rails**: a "Do Not" section naming distinct mistakes
6. **Write the description as selection criteria**: concrete requests, a handoff to overlapping siblings
7. **Run `--compare original.md rewritten.md`** before replacing the file; a FAIL on
   `drops_domain_knowledge`, `changes_role`, or `adds_unsupported_claims` means the rewrite went too far

## Color Guidelines

- **blue/cyan**: analysis, exploration, development
- **green**: validation, accessibility, success-oriented work
- **yellow**: review, caution, architectural decisions
- **red**: security, critical issues
- **magenta**: creative, generation

Use distinct colors within a plugin.
