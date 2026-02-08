# How-To Guide Template

Use this template when creating how-to guides. Read this before writing.

## Frontmatter

```yaml
---
title: "How to [accomplish specific goal]"
summary: "Configure/Set up/Deploy [thing] to [achieve outcome]."
doc_type: how-to
prerequisites:
  - "Completed [tutorial name] or equivalent experience"
  - "[Tool/service] configured and running"
est_time: "15 minutes"
roles: ["developer", "operator"]
stability: stable
---
```

## Page Structure

```markdown
# How to [accomplish specific goal]

[One sentence: what this guide helps you do and when you'd need it.]

**Time:** ~15 minutes

## Prerequisites

- [Specific version/tool/config required]
- [Link to tutorial if foundational knowledge needed]

## Step 1: [Verb phrase — Configure, Create, Install]

[Minimal context — just enough to know what you're doing]

\`\`\`bash
[Complete command or code block]
\`\`\`

Expected result:

\`\`\`
[Output confirming success]
\`\`\`

## Step 2: [Next verb phrase]

...

## Step N: [Final verb phrase]

...

## Verify it works

[Command or action that proves the goal was achieved]

\`\`\`bash
[test/verification command]
\`\`\`

You should see:

\`\`\`
[Expected output proving success]
\`\`\`

## Troubleshooting

### [Error message or symptom]

**Cause:** [Brief explanation]
**Fix:** [Specific command or action]

### [Another error]

**Cause:** ...
**Fix:** ...

### [Permission/config issue]

**Cause:** ...
**Fix:** ...

## See also

- [Related how-to guide] — for [related task]
- [API reference] — for [detailed specifications]
- [Explanation] — for [understanding why this works]
```

## Key Constraints

### Word count: under 1800 words

If you're over 1800 words:

- **Move conceptual content** to an explanation page and link to it
- **Move API details** to reference and link to it
- **Split into multiple guides** if you're covering multiple goals

### One guide = one goal

Test: can you complete the title "How to ___" with a single verb phrase? If you need "and" in the title, you have two guides.

Good: "How to configure OAuth with GitHub"
Bad: "How to configure OAuth and set up role-based access control"

### Prescriptive, not descriptive

Don't present alternatives. Pick the recommended approach. If there are genuinely different approaches for different situations, write separate how-to guides for each.

Bad: "You can use either `docker compose` or `kubectl`. If you're using Docker..."
Good: Two separate guides — "How to deploy with Docker Compose" and "How to deploy with Kubernetes"

### Prerequisites are links, not lessons

If the reader needs foundational knowledge, link to the tutorial. Don't re-teach it.

Bad: "First, let's understand how OAuth works. OAuth is a protocol that..."
Good: "Prerequisites: Completed the [Authentication tutorial](../tutorials/auth.md)"
