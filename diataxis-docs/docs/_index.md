---
title: "Diataxis Docs"
description: "Documentation transformation agents following the Diataxis framework"
weight: 10
---

# Diataxis Docs

Seven specialized agents that implement the Diataxis documentation framework: an orchestrator for coordinating transformations, four type-specific writers, an inventory scanner, and a quality validator.

## Components

| Type | Name | Description |
|------|------|-------------|
| Agent | diataxis-orchestrator | Master architect that coordinates full Diataxis transformations |
| Agent | doc-tutorial-writer | Creates 90-120 minute onboarding tutorials with checkpoints |
| Agent | doc-howto-writer | Writes task-focused how-to guides under 1800 words |
| Agent | doc-reference-gen | Generates reference documentation with the "no advice" rule |
| Agent | doc-explanation-writer | Creates conceptual and architectural explanation docs |
| Agent | doc-inventory | Scans and classifies docs by Diataxis type, identifies gaps |
| Agent | doc-crosslink-validator | Validates structure, frontmatter, type purity, and cross-links |

## Documentation

- [Getting Started](tutorials/getting-started/) — Restructure a documentation directory using the full Diataxis pipeline
- [Restructure Docs to Diataxis](howto/restructure-docs-to-diataxis/) — Transform existing docs to the four-type model
- [Write a Tutorial](howto/write-a-tutorial/) — Create tutorials with doc-tutorial-writer
- [Write Reference Docs](howto/write-reference-docs/) — Generate reference with doc-reference-gen
- [Validate Doc Quality](howto/validate-doc-quality/) — Check quality with doc-crosslink-validator
- [Agent Reference](reference/agents/) — All 7 agents specification
- [Architecture](explanation/architecture/) — The orchestrator pattern and 7 specialized agents
- [Diataxis in Practice](explanation/diataxis-in-practice/) — How the four types work together
- [Orchestration Model](explanation/orchestration-model/) — How the orchestrator coordinates specialist agents
