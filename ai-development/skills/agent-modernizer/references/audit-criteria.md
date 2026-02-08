# Agent Audit Criteria

## Frontmatter Requirements

Every agent must have all required fields. Check each:

| Field | Required | Valid Values | Notes |
|-------|----------|-------------|-------|
| `name` | Yes | lowercase, hyphens, 3-50 chars | Must start/end alphanumeric |
| `description` | Yes | 10-5,000 chars with `<example>` blocks | This is what triggers the agent |
| `model` | Yes | `inherit`, `sonnet`, `opus`, `haiku` | Raw model IDs are invalid |
| `color` | Yes | `blue`, `cyan`, `green`, `yellow`, `magenta`, `red` | Distinct per plugin |
| `tools` | Recommended | Array of tool names | Principle of least privilege |

### Description Format

The description must include:

1. A summary of when to use the agent ("Use this agent when..." or equivalent)
2. 2-4 `<example>` blocks, each containing:
   - `Context:` — situation description
   - `user:` — what the user says (in quotes)
   - `assistant:` — how Claude should respond (in quotes)
   - `<commentary>` — why this agent is appropriate

**Common description failures:**
- Missing examples entirely → agent won't trigger reliably
- Examples escaped as `\n` literals instead of actual newlines → parser can't read them
- No triggering conditions — reads like a summary instead of selection criteria
- Raw model IDs like `claude-sonnet-4-20250514` instead of `sonnet`

## System Prompt Anti-Patterns

### 1. Topic Lists Without Guidance

**Problem:** Two-word bullet lists that name concepts without providing decisions.

```markdown
# BAD
Branching strategies:
- Git Flow implementation
- GitHub Flow setup
- Trunk-based development
- Feature branch workflow
```

```markdown
# GOOD
**Branching strategy** — Choose the simplest model that works. Default to
trunk-based with short-lived feature branches unless the project genuinely
needs release branches.
```

**Why it matters:** The model already knows these concepts exist. The prompt should tell it *when to choose which* and *what to prefer*.

### 2. Fictional Progress Tracking

**Problem:** Fake JSON progress objects and delivery notifications with fabricated metrics.

```markdown
# BAD
Progress tracking:
{
  "merge_conflicts_reduced": "67%",
  "team_satisfaction": "4.5/5"
}

Delivery notification:
"Achieved 4.8/5 developer satisfaction rating."
```

These are not real data. Remove entirely.

### 3. Phantom Agent References

**Problem:** References to agents that don't exist in the plugin.

```markdown
# BAD
Integration with other agents:
- Collaborate with devops-engineer on CI/CD
- Support release-manager on versioning
- Work with security-auditor on policies
```

Remove unless these agents actually exist in the same plugin.

### 4. Redundant Sections

**Problem:** The same topic listed multiple times with slightly different groupings.

Signs of redundancy:
- A "don't" section that's just the inverse of the "do" section
- Topic categories that overlap (e.g., "Branching strategies" listed under both "Workflow patterns" and "Branching best practices")
- Duplicate numbering (two sections numbered `2.`)

### 5. Teaching the Model Its Own Knowledge

**Problem:** Listing every API method, stdlib function, or well-known pattern.

```markdown
# BAD for an Opus-class agent
- Use `errors.Is` / `errors.As` over string comparison
- Use `sync.Once`, `sync.Pool`, `sync.Map` where appropriate
- Prefer `io.Reader`/`io.Writer` interfaces
```

```markdown
# GOOD — behavioral boundary instead
Apply Effective Go and Go Code Review Comments conventions.
Favor clarity over brevity.
```

**Rule of thumb:** If the information is in the model's training data and isn't project-specific, it's padding. Focus on *priorities, decisions, and boundaries* the model can't infer.

### 6. Over-Specification of Output

**Problem:** Prescribing detailed JSON schemas, report structures, or communication protocols for inter-agent messages that don't exist.

```json
{
  "requesting_agent": "cli-developer",
  "request_type": "get_cli_context",
  "payload": { "query": "..." }
}
```

Unless this protocol is actually implemented and consumed, remove it.

## System Prompt Quality Criteria

A good system prompt:

1. **States the role** in one sentence — who is this agent?
2. **Provides decisions, not topics** — what to prefer, what to default to, when to choose which approach
3. **Defines boundaries** — what NOT to do, what to leave alone, where to stop
4. **Specifies process** — numbered steps for what the agent does when invoked
5. **Defines output format** — how to report results (severity levels, structure)
6. **Handles edge cases** — specific scenarios to watch for (generated code, vendor directories, etc.)
7. **Stays under 3,000 characters** for the body — ideally 500-2,000

## Rewriting Methodology

When rewriting an agent:

1. **Identify the core purpose** — what does this agent actually do? Strip away the padding to find it.
2. **Choose the right model** — Opus for complex judgment (architecture review, test repair), Sonnet for formulaic tasks (code review checklists, accessibility audits), Haiku for simple checks.
3. **Write for the model you chose** — Opus needs boundaries and priorities. Sonnet needs more explicit procedures. Don't teach Opus things it already knows.
4. **Add domain knowledge the model lacks** — project-specific conventions, hard-won lessons (like the termenv/AppleScript macOS workaround), platform-specific gotchas.
5. **Include guard rails** — a "Do Not" section prevents the specific mistakes that matter.
6. **Verify examples trigger correctly** — examples should cover the main invocation patterns with diverse phrasing.

## Color Guidelines

Choose colors that convey the agent's role:

- **blue/cyan**: Analysis, exploration, development
- **green**: Success-oriented, validation, accessibility
- **yellow**: Review, caution, architectural decisions
- **red**: Security, critical issues
- **magenta**: Creative, generation

Use distinct colors for agents in the same plugin to aid visual identification.
