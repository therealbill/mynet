# Tutorial Template

Use this template when creating onboarding tutorials. Read this before writing.

## Frontmatter

```yaml
---
title: "Build your first [thing] with [tool]"
summary: "A hands-on tutorial that takes you from zero to a working [thing] in about [N] minutes."
doc_type: tutorial
prerequisites:
  - "[Tool] installed (version X+)"
  - "Basic familiarity with [language/concept]"
est_time: "90 minutes"
roles: ["developer"]
stability: stable
---
```

## Page Structure

```markdown
# Build your first [thing] with [tool]

In this tutorial, you'll build a [concrete outcome]. By the end, you'll have
[working result that demonstrates core value].

**What you'll build:** [screenshot or diagram of end result]

**Time:** ~90 minutes

## Prerequisites

- [Tool] version X or later ([installation link])
- A text editor
- [Any other requirement with install link]

**Verify your setup:**

\`\`\`bash
[tool] --version
# Expected: [tool] vX.Y.Z
\`\`\`

## Step 1: [Action that produces visible result]

[One paragraph of context — what this step achieves, not why it works]

\`\`\`bash
[Complete command — never use "..." or "[your-value]"]
\`\`\`

You should see:

\`\`\`
[Exact expected output]
\`\`\`

## Step 2: [Next action]

[Instructions...]

## Checkpoint: Verify [milestone]

At this point you should have:

- [ ] [Concrete, verifiable condition]
- [ ] [Another verifiable condition]

Run this to confirm:

\`\`\`bash
[verification command]
\`\`\`

Expected output:

\`\`\`
[exact output]
\`\`\`

> **Stuck?** See [Troubleshooting](#troubleshooting) below.

## Steps 3-N...

[Continue pattern: action title, context sentence, complete code, expected output]

[Checkpoint every 3-5 steps]

## What you built

You now have a working [thing] that:

- [Capability 1]
- [Capability 2]
- [Capability 3]

## Next steps

- **Do more:** [How-to guide for common next task] →
- **Go deeper:** [Explanation of how it works internally] →
- **Look up details:** [Reference for the API/CLI you used] →

## Troubleshooting

### [Symptom 1: exact error message or behavior]

**Cause:** [Why this happens]

**Fix:**

\`\`\`bash
[Exact fix command]
\`\`\`

### [Symptom 2]

**Cause:** ...
**Fix:** ...
```

## Key Decisions

### Tutorial length

- **Onboarding tutorials:** 90-120 minutes, 15-25 steps
- **Feature tutorials:** 30-45 minutes, 8-12 steps
- **Quick-start tutorials:** 10-15 minutes, 5-7 steps

### Scaffolding fade

Early steps (1-5): Very explicit — show every keystroke, every file path, every expected output.

Middle steps (6-15): Slightly less hand-holding — "Create a file `routes.js` with:" (no longer explaining what a file is).

Late steps (15+): Assume patterns from earlier — "Add another route following the same pattern" (but still show the code).

### What NOT to include

- **Choices** — Never say "you could use X or Y." Pick X. Period.
- **Explanations** — Never explain how the internals work. Link to an explanation page.
- **Best practices** — Don't teach the "right" way. Teach the "working" way.
- **Error handling** — Add it only if the tutorial won't work without it. Save proper error handling for how-to guides.
