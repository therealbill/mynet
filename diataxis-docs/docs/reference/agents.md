---
title: "Agent reference"
description: "Complete specification of all seven diataxis-docs agents: orchestrator, four type-specific writers, inventory scanner, and cross-link validator."
weight: 1
doc_type: reference
stability: stable
version: "1.0.0"
---

# Agent Reference

Complete specification of all seven agents in the diataxis-docs plugin.

## diataxis-orchestrator

| Property | Value |
|----------|-------|
| Model | inherit |
| Color | blue |
| Tools | Read, Write, Edit, Grep, Glob, Bash |

### Purpose

Coordinates Diataxis transformation workflows. Assesses current documentation state, plans the restructuring, and delegates to the six specialist agents.

### Capabilities

- Assess documentation state and clarify transformation goals (restructure, new docs, audit, fill gaps)
- Gather context: docs path, source code path, top user tasks, available introspection tools
- Plan delegation workflow and agent invocation order
- Review each agent's output before invoking the next
- Maintain partial progress if any agent encounters issues

### Coordination Workflows

| Workflow | Agent Sequence |
|----------|----------------|
| Restructure existing docs | doc-inventory, doc-reference-gen, doc-tutorial-writer, doc-howto-writer, doc-explanation-writer, doc-crosslink-validator |
| New documentation from scratch | doc-reference-gen, doc-tutorial-writer, doc-howto-writer, doc-explanation-writer, doc-crosslink-validator |
| Audit only | doc-inventory, doc-crosslink-validator |

### Constraints

- Does not write documentation content directly
- Does not skip the assessment phase
- Does not invoke agents without passing docs path and source path
- Does not proceed without user confirmation on structural changes

### Target Structure

```
docs/
  start-here.md
  tutorials/00-onboarding.md
  how-to/*.md
  reference/api/*.md
  explanation/*.md
  _reports/inventory.json, quality.json
```

### Trigger Conditions

- "Restructure our docs to follow Diataxis"
- "We need docs for this new service"
- "Can you check if our docs follow Diataxis properly?"
- "Audit our documentation coverage"

---

## doc-tutorial-writer

| Property | Value |
|----------|-------|
| Model | inherit |
| Color | yellow |
| Tools | Read, Write, Edit, Grep, Glob, Bash |

### Purpose

Creates learning-oriented onboarding tutorials that guide complete beginners through building a working project.

### Capabilities

- Design progressive learning journeys from installation to working outcome
- Write step-by-step instructions with complete code and expected output
- Add checkpoints every 3-5 steps with verification commands
- Include troubleshooting sections covering common setup issues
- Cross-link to reference, how-to guides, and explanations

### Constraints

- Target duration: 90-120 minutes for onboarding tutorials
- One path only: no alternatives or choices presented
- Progressive skill-building: each step introduces one new concept
- Fading scaffolding: early steps are explicit, later steps assume learned patterns
- No explanations of why things work (link to explanation pages instead)
- No best practices or optimization (link to how-to guides instead)

### Output Structure

```
frontmatter (title, summary, prerequisites, est_time, roles, stability)
# Title (states what learner will build)
## Prerequisites
## What you'll build
## Steps 1-N (action title, instructions, code, expected output)
### Checkpoint (every 3-5 steps)
## What you built
## Next steps
## Troubleshooting
```

### Trigger Conditions

- "Create an onboarding tutorial for our CLI tool"
- "Our getting started guide is broken"
- "We need tutorials for beginners and intermediate users"

---

## doc-howto-writer

| Property | Value |
|----------|-------|
| Model | inherit |
| Color | magenta |
| Tools | Read, Write, Edit, Grep, Glob, Bash |

### Purpose

Creates task-focused how-to guides that solve specific problems with numbered steps.

### Capabilities

- Write single-goal guides with numbered, sequential steps
- Include complete code and commands at each step
- Add verification sections proving the goal was achieved
- Write troubleshooting sections covering 3-5 common failure modes
- Cross-link to reference, tutorials, explanations, and related how-to guides

### Constraints

- Maximum length: 1800 words
- One guide per goal: multi-goal content must be split
- Assumes existing knowledge: links to tutorials for basics
- Prescriptive: picks the recommended approach, does not list alternatives
- No teaching from scratch (link to tutorials instead)
- No explanations of why (link to explanation pages instead)
- No API specifications (link to reference pages instead)

### Output Structure

```
frontmatter (title, summary, prerequisites, est_time, roles, stability)
# How to [Goal]
Goal statement, time estimate
## Prerequisites
## Steps (numbered: 1. Action title, code, expected result)
## Verify it works
## Troubleshooting (3-5 problem/symptom/cause/solution blocks)
## Next steps / See also
```

### Title Convention

Titles start with "How to" followed by a verb phrase: "How to configure OAuth," "How to deploy to production."

### Trigger Conditions

- "Write a guide for setting up OAuth authentication"
- "Our deployment guide is 4000 words and covers too much"
- "Document the 5 most common things our users need to do"

---

## doc-reference-gen

| Property | Value |
|----------|-------|
| Model | inherit |
| Color | green |
| Tools | Read, Write, Edit, Grep, Glob, Bash |

### Purpose

Generates complete, accurate API/CLI/SDK reference documentation from source code or introspection tools.

### Capabilities

- Extract public API surfaces from source code using introspection tools
- Generate structured reference pages with consistent formatting
- Document parameters, return values, error codes, and limits
- Support multiple reference types: REST API, CLI, SDK, configuration
- Cross-link to how-to guides and explanations

### The "No Advice" Rule

Reference pages state facts. They never advise. The following transformations apply:

| Violation | Replacement |
|-----------|-------------|
| "You should validate inputs" | "Input validation is the caller's responsibility" |
| "We recommend the async version" | "An async version is available as `sendEmailAsync`" |
| "It is best to use connection pooling" | "Connection pooling is supported via the `pool` parameter" |

Any content containing opinions, advice, or "you should" is moved to a how-to guide.

### Supported Introspection Tools

| Language/Platform | Tools |
|-------------------|-------|
| TypeScript/JavaScript | typedoc, JSDoc |
| Python | pydoc, type hints |
| Go | go doc |
| Rust | cargo doc |
| REST API | OpenAPI/Swagger, GraphQL schema |
| CLI | --help output, man pages |

### Constraints

- Every public API surface must be documented
- Consistent structure across all entries of the same kind
- Parameters include: type, required/optional, default, constraints
- Error codes include: condition and HTTP status
- Examples show syntax only, not workflows or use cases
- Version and stability markers on every page

### Output Structure (API)

```
frontmatter (title, summary, stability, version)
# Name
## Signature
## Parameters (table: name, type, required, default, description)
## Returns
## Errors (table: code, condition, status)
## Limits
## Examples (syntax only)
## See Also
```

### Output Structure (CLI)

```
frontmatter (title, summary, stability, version)
# Command Name
## Synopsis
## Arguments
## Options (table: flag, type, default, description)
## Exit Codes
## Examples
## See Also
```

### Trigger Conditions

- "Generate reference docs for our REST API"
- "Document all CLI commands and flags"
- "Our API docs have too many 'you should' statements"

---

## doc-explanation-writer

| Property | Value |
|----------|-------|
| Model | inherit |
| Color | red |
| Tools | Read, Write, Edit, Grep, Glob, Bash |

### Purpose

Creates understanding-oriented conceptual and architectural documentation that explains why systems work the way they do.

### Capabilities

- Explain design decisions with rationale, alternatives considered, and trade-offs
- Build mental models using analogies, diagrams, and progressive explanation
- Discuss trade-offs explicitly with benefit-vs-cost analysis
- Clarify common misconceptions
- Cross-link to how-to guides for practical application and reference for details

### Constraints

- Understanding-oriented: answers "why" and "how it works internally"
- No step-by-step instructions (link to how-to guides instead)
- No bare API specifications (link to reference pages instead)
- Not prescriptive: descriptive language only
- Trade-offs must be discussed honestly (every design gains and sacrifices something)
- Code examples are for illustration only, not as recipes

### Output Structure

```
frontmatter (title, summary, prerequisites, est_time, roles, stability)
# Understanding [Concept]
Opening: what this covers and why it matters
## The problem
## Overview (high-level mental model, diagram)
## How it works
## Why it's designed this way
## Trade-offs (table: benefit vs cost)
## Common misconceptions
## Related
```

### Trigger Conditions

- "Document why we chose event-driven architecture"
- "Users don't understand how our permission system works"
- "Our config guide spends 2000 words explaining the config philosophy"

---

## doc-crosslink-validator

| Property | Value |
|----------|-------|
| Model | inherit |
| Color | blue |
| Tools | Read, Write, Edit, Grep, Glob, Bash |

### Purpose

Validates documentation structure, frontmatter, Diataxis type separation, and cross-linking. Generates quality reports.

### Capabilities

- Validate frontmatter completeness and correctness
- Check type classification and type purity
- Verify all internal links resolve
- Analyze cross-linking patterns and identify gaps
- Find orphaned pages
- Generate machine-readable and human-readable quality reports

### Validation Checks

| Check | Severity | Description |
|-------|----------|-------------|
| Missing required frontmatter | Error | title, summary, doc_type, stability must be present |
| Invalid doc_type | Error | Must be one of: tutorial, how-to, reference, explanation |
| Invalid stability | Error | Must be one of: stable, beta, experimental, deprecated |
| Type purity violation | Error | Content mixing Diataxis types |
| Broken internal link | Error | Link target does not exist |
| Missing prerequisites field | Warning | Tutorials and how-tos need prerequisites |
| Missing est_time field | Warning | Tutorials and how-tos need time estimates |
| Orphaned page | Warning | No other page links to this page |
| Insufficient cross-links | Warning | Below minimum cross-link threshold |
| Cross-link opportunity | Info | Suggested link between related pages |

### Cross-link Minimum Thresholds

| Doc Type | Minimum Links | Link Targets |
|----------|---------------|--------------|
| Tutorial | 3 | Reference, how-to guides, explanations |
| How-to | 5 | Prerequisites, reference, related how-tos, explanations |
| Reference | 2 | How-to guides showing usage |
| Explanation | 4 | How-to guides, reference, related concepts |

### Type Purity Rules

| Doc Type | Red Flags |
|----------|-----------|
| Tutorial | API specs, conceptual discussion, no project walkthrough |
| How-to | Teaching basics, explaining why |
| Reference | Step-by-step instructions, "we recommend" |
| Explanation | Step-by-step instructions, bare API specs |

### Severity Levels

| Level | Meaning | Action |
|-------|---------|--------|
| Error | Must fix | Blocks publishing |
| Warning | Should fix | Degrades quality |
| Info | Nice to have | Improves discoverability |

### Output

- `docs/_reports/quality.json` -- machine-readable report with issue counts, per-file issues, link graph analysis
- Markdown summary -- human-readable overview with prioritized recommendations

### Trigger Conditions

- "Check our docs for quality issues before we publish"
- "Do any of our reference pages have 'you should' statements?"
- "Our docs feel disconnected"

---

## doc-inventory

| Property | Value |
|----------|-------|
| Model | inherit |
| Color | cyan |
| Tools | Read, Write, Edit, Grep, Glob, Bash |

### Purpose

Scans existing documentation and classifies each page by Diataxis type. Identifies gaps and generates an inventory report.

### Capabilities

- Scan all `.md`, `.mdx`, and `.rst` files in a documentation directory
- Classify each page as Tutorial, How-to, Reference, Explanation, or Mixed
- Rate quality match: strong, moderate, weak, mixed
- Identify documentation gaps by Diataxis type
- Recommend reorganization actions (move, split, merge)
- Generate structured inventory report

### Classification Signals

| Type | Signals |
|------|---------|
| Tutorial | "Getting started," "quickstart," "your first...," walks through building something, has checkpoints, imperative mood |
| How-to | "How to...," "configuring...," recipe-style, numbered steps, single goal, assumes knowledge |
| Reference | API docs, CLI reference, config options, consistent structure, no advice, often auto-generated |
| Explanation | "Architecture," "concepts," "understanding...," discusses "why," design decisions, trade-offs, diagrams |
| Mixed | Combines multiple types; classified with note on which types are present |

### Output

`docs/_reports/inventory.json` with structure:

```json
{
  "scanned_at": "ISO timestamp",
  "docs_path": "path scanned",
  "total_pages": 0,
  "by_type": {
    "tutorial": 0,
    "how-to": 0,
    "reference": 0,
    "explanation": 0,
    "mixed": 0
  },
  "pages": [
    {
      "path": "relative path",
      "title": "page title",
      "type": "classified type",
      "quality": "strong|moderate|weak|mixed",
      "issues": ["list of issues"],
      "word_count": 0
    }
  ],
  "gaps": [
    {
      "type": "missing type",
      "severity": "critical|high|medium|low",
      "description": "what is missing",
      "suggestion": "what to create"
    }
  ],
  "reorganization": [
    {
      "action": "move|split|merge",
      "from": "source path",
      "to": "target path",
      "reason": "why"
    }
  ],
  "recommendations": []
}
```

### Constraints

- Scans every documentation file without skipping any
- Reads content before classifying (does not guess from filenames)
- Does not create documentation, only inventories and recommends
- Mixed pages include specific guidance on which types are present and how to split

### Trigger Conditions

- "Audit our docs and tell me what we have"
- "We have docs but I think we're missing important pieces"
- "Our docs are a mess -- some pages mix tutorials with reference material"

---

## See Also

- {{< ref "tutorials/getting-started" >}} -- Tutorial using all seven agents
- {{< ref "howto/restructure-docs-to-diataxis" >}} -- Restructuring workflow using these agents
- {{< ref "explanation/architecture" >}} -- Why seven specialized agents instead of one
- {{< ref "explanation/orchestration-model" >}} -- How the orchestrator coordinates the specialists
