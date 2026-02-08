---
name: agent-modernizer
description: >
  This skill should be used when the user asks to "modernize an agent", "audit agent definitions",
  "update agent format", "check agent quality", "rewrite agent prompts", or wants to bring agent
  files up to current Claude Code plugin standards. Also applies when reviewing agents for verbose
  content, missing frontmatter, or weak triggering descriptions.
---

# Agent Modernizer

Audit and rewrite Claude Code plugin agent definitions to meet current standards. Transform verbose,
generic agent prompts into concise, opinionated definitions that trust the model's knowledge and
focus on decisions, boundaries, and priorities.

## When to Use

- Auditing existing agent `.md` files against the plugin agent spec
- Rewriting agents that have verbose bullet-point lists, missing frontmatter, or weak descriptions
- Creating new agents that follow current best practices from the start
- Reviewing a batch of agents in a plugin for consistency

## Process

### 1. Audit

Read the agent file and evaluate against the spec. Check every frontmatter field and assess the
system prompt quality. For the complete checklist, consult `references/audit-criteria.md`.

Produce a findings table:

| # | Field/Area | Issue | Severity |
|---|-----------|-------|----------|
| 1 | `color` | Missing — required field | Must fix |
| 2 | `description` | No `<example>` blocks | Must fix |
| ... | ... | ... | ... |

Severity levels:

- **Must fix** — Agent won't load or won't trigger correctly without this
- **Recommended** — Agent will work but is suboptimal
- **Minor** — Stylistic or organizational improvement

### 2. Assess System Prompt

Scan the body for common anti-patterns. For the full list, consult the "System Prompt Anti-Patterns"
section in `references/audit-criteria.md`. The key patterns to flag:

- **Topic lists without guidance** — two-word bullets that name concepts without providing decisions
- **Teaching the model its own knowledge** — listing stdlib functions, well-known patterns, or API methods that the chosen model already knows
- **Fictional content** — fake JSON progress trackers, fabricated metrics, delivery notification scripts
- **Phantom references** — mentions of agents, systems, or protocols that don't exist
- **Redundancy** — same topic covered multiple times, "don't" sections that invert the "do" section

Count total bullet points and flag if over ~20 — a sign the prompt is listing topics instead of providing guidance.

### 3. Rewrite

When rewriting, apply these principles:

**Trust the model.** Opus doesn't need to be told what `gofmt` is or that `errors.Is` exists. State
priorities and boundaries, not concept inventories.

**Decisions over topics.** Instead of "Branching strategies: Git Flow, GitHub Flow, trunk-based,"
write "Default to trunk-based with short-lived feature branches unless the project genuinely needs
release branches."

**Guard rails over checklists.** A "Do Not" section that prevents the 3-5 most common mistakes is
more valuable than a 20-item checklist of things to do.

**Concise role statement.** One sentence establishing who the agent is. No "years of experience"
padding.

**Concrete process.** Numbered steps for what happens when the agent is invoked. Each step should
describe an action, not a topic.

**Structured output.** Define how results should be reported — severity levels, file references,
format.

### 4. Validate

After rewriting, verify:

- All required frontmatter fields present (`name`, `description`, `model`, `color`)
- Description includes 2-4 `<example>` blocks with proper structure
- `tools` array follows least-privilege principle
- System prompt is under 3,000 characters (ideally 500-2,000)
- No anti-patterns remain
- Agent-specific domain knowledge is preserved (don't strip things the model genuinely can't infer)

### 5. Batch Audit

When auditing multiple agents in a directory, produce a summary table:

| Agent | Lines | Missing Fields | Anti-patterns | Action Needed |
|-------|-------|---------------|---------------|---------------|
| `code-reviewer` | 31 | `color`, examples | Generic checklist | Rewrite |
| `architect-reviewer` | 43 | `model`, `color`, `tools` | None significant | Fix frontmatter |

Prioritize rewrites by severity. Fix frontmatter-only issues first, then tackle full rewrites.

## Key Decisions for Rewrites

**Model selection:**

- `opus` — Complex judgment: architecture review, test diagnosis, code simplification, TUI design
- `sonnet` — Formulaic review: code review checklists, accessibility audits, Playwright tests
- `haiku` — Simple validation: format checks, linting
- `inherit` — When the parent's model is always appropriate

**When to merge agents:** If two agents share >70% of their system prompt or would need the same
context to operate, merge them. Add a section for the secondary concern rather than maintaining a
separate agent with overlapping triggers.

**When to split agents:** If an agent's system prompt exceeds 5,000 characters and covers genuinely
distinct domains (e.g., TUI development + database design), split into focused agents.

## Additional Resources

### Reference Files

For detailed audit criteria, anti-pattern examples, and rewriting methodology:

- **`references/audit-criteria.md`** — Complete frontmatter checklist, system prompt anti-patterns with before/after examples, quality criteria, color guidelines, and rewriting methodology
