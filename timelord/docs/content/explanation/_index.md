---
title: "Explanation"
weight: 3
bookCollapseSection: true
---

# Explanation

Conceptual documentation explaining how Temporal works and why it's designed the way it is.

## Contents

### [How Temporal Works]({{< relref "how-temporal-works" >}})

Understanding Temporal's core architecture:

- Event sourcing and workflow state
- History replay and determinism
- Task queues and workers
- Persistence and durability

### [Why Determinism Matters]({{< relref "why-determinism-matters" >}})

Deep dive into workflow determinism:

- What makes code non-deterministic
- How replay reconstructs state
- Common determinism pitfalls
- Versioning strategies

### [Why Nexus]({{< relref "why-nexus" >}})

Understanding when and why to use Temporal Nexus:

- The problem Nexus solves for cross-namespace communication
- What Nexus provides over activity-based workarounds
- Nexus vs alternatives comparison
- When Nexus is overkill and when it shines

## Who This Is For

This section is for readers who want to understand:

- **The "why"** behind Temporal's design decisions
- **Mental models** for thinking about workflow orchestration
- **Trade-offs** between different approaches
- **Background context** for making architectural decisions

## Contrast with Other Sections

| Section | Purpose |
|---------|---------|
| **Tutorials** | Learn by building step-by-step |
| **How-to** | Solve specific problems |
| **Reference** | Look up exact specifications |
| **Explanation** | Understand concepts deeply |
