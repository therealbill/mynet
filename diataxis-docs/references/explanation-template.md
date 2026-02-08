# Explanation Template

Use this template when creating explanation/conceptual documentation. Read this before writing.

## Frontmatter

```yaml
---
title: "Understanding [concept/system]"
summary: "How [system] works internally, why it was designed this way, and the trade-offs involved."
doc_type: explanation
prerequisites:
  - "Familiarity with [concept] (see [tutorial link])"
est_time: "20 minutes"
roles: ["developer", "architect"]
stability: stable
---
```

## Page Structure

```markdown
# Understanding [concept/system]

[Opening paragraph: what this document covers and why a reader would care.
Frame the problem this concept/system addresses.]

## The problem

[What challenge or need motivated this design? What would happen without it?
Ground the reader in the real-world context before introducing the solution.]

## Overview

[High-level mental model. Use an analogy if it genuinely clarifies.
Include a diagram showing the key components and their relationships.]

\`\`\`
┌──────────┐     ┌──────────┐     ┌──────────┐
│ Component│────>│ Component│────>│ Component│
│    A     │     │    B     │     │    C     │
└──────────┘     └──────────┘     └──────────┘
     ▲                                  │
     └──────────────────────────────────┘
\`\`\`

Think of [system] as [mental model framing]. [One or two sentences
expanding the model.]

## How it works

### [Aspect 1]

[Explain how this aspect works at a conceptual level.
Use code examples for illustration, not as recipes.]

\`\`\`python
# Illustrative — shows the concept, not a complete implementation
event = Event(type="user.created", data=user)
bus.publish(event)  # Async: handlers process independently
\`\`\`

[Explain what this demonstrates about the design.]

### [Aspect 2]

[Continue with the next key aspect...]

### [Aspect 3]

...

## Why it's designed this way

[Explain the design rationale. What alternatives were considered?
What constraints shaped the decision?]

### Alternatives considered

| Approach | Pros | Cons | Why not chosen |
|----------|------|------|----------------|
| [Alternative A] | [benefits] | [drawbacks] | [reason] |
| [Alternative B] | [benefits] | [drawbacks] | [reason] |
| [Current approach] | [benefits] | [drawbacks] | **Chosen because:** [reason] |

## Trade-offs

[Every design gains something and sacrifices something. Be explicit.]

| Gained | Sacrificed |
|--------|-----------|
| [Benefit 1] | [Cost 1] |
| [Benefit 2] | [Cost 2] |
| [Benefit 3] | [Cost 3] |

## Common misconceptions

### "[Misconception 1]"

[Why people think this, and what's actually true.]

### "[Misconception 2]"

[Why people think this, and what's actually true.]

## Related

- **Do it:** [How to configure/use this](../how-to/relevant-task.md)
- **Look it up:** [API reference for this system](../reference/api.md)
- **Related concept:** [Understanding related-concept](./related-concept.md)
```

## Writing Style

### Discursive, not imperative

Explanation uses analytical language. It discusses, compares, contextualizes.

Good: "The system uses event sourcing because it needs a complete audit trail. This trades storage efficiency for auditability."

Bad: "Use event sourcing for audit trails. First, create an event store..."

### Mental models over component lists

Good: "Think of the permission system as a series of gates. Each request passes through authorization, then role checking, then resource-level permissions. A request must pass all gates."

Bad: "The permission system has three components: AuthorizationService, RoleChecker, and ResourcePermissionManager."

### Trade-offs are mandatory

Every explanation of a design decision must include what was gained AND what was sacrificed. If you can't identify a trade-off, you haven't understood the design deeply enough.

### Code is illustrative, not executable

Code examples in explanations show concepts, not complete implementations. They may be simplified, incomplete, or pseudocode. They should never be copy-pasteable recipes — that's what how-to guides are for.

### Diagrams for architecture

ASCII diagrams, Mermaid, or descriptions of visual relationships are essential for architectural explanations. If you're explaining how components interact without a visual, the explanation is incomplete.
