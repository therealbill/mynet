---
title: "Understanding the orchestration model"
description: "How diataxis-orchestrator coordinates the six specialist agents, why coordination is separated from content creation, and how invocation order is determined."
weight: 3
doc_type: explanation
prerequisites:
  - "Familiarity with the diataxis-docs agent architecture (see [Architecture]({{< ref \"explanation/architecture\" >}}))"
  - "Understanding of the four Diataxis types (see [Diataxis in Practice]({{< ref \"explanation/diataxis-in-practice\" >}}))"
est_time: "12 minutes"
roles: ["developer", "architect"]
stability: stable
---

# Understanding the Orchestration Model

The diataxis-orchestrator agent coordinates six specialist agents through a structured workflow without writing any documentation content itself. This page explains how orchestration works, why coordination is separated from creation, and how the orchestrator determines which agents to invoke and in what order.

## The Problem

A documentation transformation involves multiple interdependent tasks: scanning existing content, classifying it, planning the restructuring, writing new content in four different types, and validating the result. These tasks have ordering constraints (you cannot validate content that has not been written) and context dependencies (writers need to know what reference pages exist before creating cross-links).

Without coordination, these tasks happen in an ad hoc order, leading to missing cross-links, redundant work, and content that does not fit together as a coherent documentation set.

## Overview

The orchestration model follows a pipeline pattern with three phases:

```
Phase 1: Analysis         Phase 2: Creation              Phase 3: Validation
+---------------+     +---------------------------+     +---------------------+
| doc-inventory  | --> | doc-reference-gen          | --> | doc-crosslink-      |
|                |     | doc-tutorial-writer        |     |   validator          |
|                |     | doc-howto-writer           |     |                     |
|                |     | doc-explanation-writer     |     |                     |
+---------------+     +---------------------------+     +---------------------+
```

The orchestrator manages transitions between phases, reviews output at each step, and adjusts the plan based on intermediate results.

## How It Works

### Phase 1: Analysis

Every orchestration begins with understanding what exists. The orchestrator gathers three pieces of context from the user:

- **Docs path**: where the existing documentation lives
- **Source code path**: where the project source code lives (for reference generation)
- **User goals**: restructure, create from scratch, audit, or fill specific gaps

For restructuring tasks, the orchestrator invokes **doc-inventory** first. The inventory agent scans every documentation file, classifies each by Diataxis type, identifies mixed-type pages, and reports coverage gaps. The resulting `inventory.json` becomes the foundation for the transformation plan.

For new documentation (no existing docs), the orchestrator skips inventory and moves directly to Phase 2.

### Phase 2: Creation

The orchestrator plans the content creation order based on the inventory results and a dependency principle: types that other types link to are created first.

**Reference comes first.** Reference pages are the most frequently linked-to documentation type. Tutorials link to reference for API details. How-to guides link to reference for parameter specifications. Explanations link to reference for technical facts. By generating reference first, the subsequent writer agents have concrete targets for their cross-links.

**Tutorials come second.** Tutorials establish the foundational learning path. How-to guides link to tutorials as prerequisites ("Complete the Getting Started tutorial first"). Creating tutorials before how-to guides ensures these prerequisite links have valid targets.

**How-to guides come third.** With reference and tutorials in place, how-to guides can link to both: reference for API details and tutorials for prerequisite knowledge. How-to guides also link to each other for related tasks, so they are written as a batch.

**Explanations come last.** Explanation pages link to all other types: reference for technical details, how-to guides for practical application, and tutorials for learning context. Writing explanations last means all link targets exist.

This ordering is not rigid. The orchestrator adjusts based on the specific transformation. If existing docs already have strong reference coverage, it may skip doc-reference-gen and start with tutorials. If the user only needs how-to guides, it invokes only doc-howto-writer.

### Between Each Agent

After each agent completes its work, the orchestrator:

1. **Reviews the output** to confirm it meets expectations and the agent's constraints
2. **Updates the plan** if the output reveals additional work needed (for example, a tutorial that references an API not yet in the reference pages)
3. **Passes context** to the next agent, including the paths of newly created pages so cross-links can be added
4. **Reports progress** to the user, summarizing what was created and what comes next

This review-between-steps pattern is what distinguishes orchestration from a simple sequential script. The orchestrator applies judgment to intermediate results and adapts.

### Phase 3: Validation

After all content is created, the orchestrator invokes **doc-crosslink-validator** to check the complete documentation set. Validation occurs after all content exists because:

- Cross-link validation requires all target pages to exist
- Type purity validation benefits from seeing the full documentation set in context
- Frontmatter validation catches issues across the entire set, not just individual pages

The validator generates a quality report. If the report contains errors, the orchestrator can re-invoke the appropriate writer agent to fix specific issues, then re-validate.

## Why Coordination Is Separated from Creation

The orchestrator deliberately does not write documentation content. This separation exists for three reasons.

### Constraint Enforcement

Each content-layer agent has type-specific constraints: the tutorial writer follows one path only, the reference generator enforces "no advice," the how-to writer stays under 1800 words. If the orchestrator wrote content directly, it would need to enforce all four sets of constraints simultaneously, which erodes constraint effectiveness as context grows.

By delegating to specialists, each agent operates with a focused system prompt containing only the constraints relevant to its type. The constraints are structurally enforced, not merely suggested.

### Separation of Concerns

Coordination logic (what to do when, in what order, with what context) is fundamentally different from content creation logic (how to write a good tutorial, how to enforce the "no advice" rule). Mixing these concerns in one agent would make both harder to improve.

Improving the orchestrator's sequencing logic does not risk degrading tutorial quality. Improving the tutorial writer's checkpoint design does not risk breaking the coordination workflow.

### Auditability

When the orchestrator delegates to a specialist, the delegation is visible: the orchestrator invoked doc-reference-gen with specific parameters and received specific output. This makes the transformation auditable. If a reference page contains advice, the issue traces to doc-reference-gen, not to an opaque "orchestrator wrote everything" process.

## How the Orchestrator Decides What to Invoke

The orchestrator uses three inputs to determine its workflow:

### 1. User Goal

The user's stated goal selects the workflow template:

| User Goal | Workflow |
|-----------|----------|
| "Restructure our docs" | Full pipeline: inventory, all four writers, validator |
| "Create docs from scratch" | Skip inventory, start with reference generation |
| "Audit our docs" | Inventory and validator only |
| "Fill gaps" | Inventory to identify gaps, then only the writers needed for missing types |
| "Write a tutorial" | Tutorial writer only |
| "Generate reference" | Reference generator only |

### 2. Inventory Results

When doc-inventory runs, its output refines the plan:

- **No existing tutorials** (critical gap): doc-tutorial-writer is prioritized
- **Mixed-type pages found**: the relevant writers are invoked to create split pages
- **Strong reference coverage**: doc-reference-gen may be skipped or limited to updates
- **No existing docs at all**: the inventory agent confirms this and the orchestrator switches to the "from scratch" workflow

### 3. Intermediate Results

The orchestrator adjusts after each agent's output:

- If doc-reference-gen creates pages that reveal previously unknown API surfaces, the tutorial plan may expand to cover them
- If doc-tutorial-writer creates a tutorial referencing a how-to that does not exist, that how-to is added to the plan for doc-howto-writer
- If doc-crosslink-validator reports type purity errors, the orchestrator re-invokes the appropriate writer to fix the specific page

## Trade-offs

| Benefit | Cost |
|---------|------|
| Clear separation of coordination and creation | More agent invocations per transformation |
| Each agent operates with focused constraints | Context must be passed between agents explicitly |
| Intermediate review catches issues early | Pipeline takes longer than a single-agent approach |
| Workflow adapts to inventory findings | Orchestrator logic is more complex than a simple script |
| Partial progress preserved on failure | Partial results may need manual cleanup |
| Auditable delegation chain | Users must trust the orchestrator's judgment on sequencing |

## Common Misconceptions

**"The orchestrator is just running agents in a fixed order."**
The order is adapted based on user goals, inventory results, and intermediate output. While the default sequence (inventory, reference, tutorial, how-to, explanation, validation) is common, the orchestrator skips, reorders, or re-invokes agents as needed.

**"Skipping the orchestrator and invoking agents directly is faster."**
For single-agent tasks (writing one how-to guide), invoking doc-howto-writer directly is appropriate and the orchestrator is unnecessary. For multi-agent tasks (restructuring a documentation directory), the orchestrator's sequencing, context passing, and intermediate review prevent issues that take longer to fix after the fact.

**"The orchestrator makes all decisions without user input."**
The orchestrator asks for user confirmation before major structural changes. It presents the transformation plan and waits for approval. The user retains control over what gets restructured and how.

## Related

- {{< ref "reference/agents" >}} -- Specifications for the orchestrator and all six specialist agents
- {{< ref "explanation/architecture" >}} -- Why seven specialized agents instead of one
- {{< ref "explanation/diataxis-in-practice" >}} -- The Diataxis types that drive the creation order
- {{< ref "howto/restructure-docs-to-diataxis" >}} -- Practical guide using the orchestration workflow
- {{< ref "tutorials/getting-started" >}} -- Tutorial walking through the full orchestrated pipeline
