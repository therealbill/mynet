---
title: "Architecture"
description: "Two-pronged purpose: building AI features and maintaining agent quality"
weight: 1
---

# Architecture

The ai-development plugin serves two distinct audiences through two components: builders who add AI features to applications (ai-engineer) and plugin maintainers who keep agent definitions at a consistent quality bar (agent-modernizer). Both relate to "AI development," but they address different problems for different people.

## Two Purposes, One Plugin

Most plugins in the marketplace have a single focus. The ai-development plugin is unusual because it bundles an implementation agent with a quality-control skill. The rationale is that both components share a common concern -- making AI-powered systems work well in practice -- even though they operate at different levels.

The ai-engineer agent works at the application level: integrating LLMs, building recommendation engines, adding computer vision. The agent-modernizer skill works at the plugin ecosystem level: auditing agent definitions, enforcing formatting standards, producing consistent rewrites.

These are not in tension. A plugin maintainer using agent-modernizer benefits from understanding how AI agents should behave in practice. An AI feature builder using ai-engineer benefits from well-structured agent definitions that trigger reliably. The shared context justifies the shared plugin.

## ai-engineer: Practical AI Implementation

The ai-engineer agent favors pragmatic deployment over theoretical elegance. Its defaults reflect this: prefer pre-trained models over training from scratch, default to RAG over fine-tuning, use the smallest model that meets accuracy requirements. These are opinionated positions, not a survey of options.

This design decision matters because AI implementation has a well-known failure mode: over-engineering. Teams reach for custom model training when a managed API call would suffice, or fine-tune when RAG handles the use case. The ai-engineer agent steers toward the simplest solution that works, then scales up only with evidence that the simpler approach is insufficient.

The agent also enforces practical concerns that are easy to overlook: cost estimation before committing to an approach, fallback behavior for model downtime, graceful degradation when AI quality drops. These are the concerns that separate a prototype from a production system.

## agent-modernizer: Ecosystem Quality Control

As the marketplace grows, agent quality becomes a scaling problem. A single plugin author can maintain consistency across three agents. Across fifty plugins from different authors, definitions drift: some have verbose prompts that teach the model its own knowledge, others lack `<example>` blocks and trigger unreliably, others reference agents that do not exist.

The agent-modernizer skill addresses this by codifying audit criteria into a repeatable process. It checks frontmatter completeness, evaluates system prompt quality against known anti-patterns, and produces a findings table with severity levels. This transforms agent quality from a subjective judgment ("this prompt feels too long") into a structured evaluation ("34 bullet points listing concepts without decisions: Should fix").

## Why a Skill, Not an Agent

The agent-modernizer is a skill rather than an agent because it provides knowledge injection -- audit criteria, anti-pattern definitions, rewriting principles -- rather than autonomous multi-step behavior. It enhances whatever agent is handling the conversation by injecting the specific knowledge needed to evaluate and rewrite agent definitions.

An agent would be appropriate if the modernizer needed to autonomously discover agent files, run tests, and iterate on rewrites without user guidance. The current design assumes the user drives the process: they choose which agent to audit, review the findings, and decide whether to apply a rewrite. The skill provides the expertise; the user provides the judgment.

## Cross-Plugin Value

The agent-modernizer works on agents from any plugin in the marketplace, not just agents within ai-development. This makes it a meta-tool for marketplace quality. A plugin author working on the devops plugin can use agent-modernizer to audit their CI/CD agents. A maintainer reviewing a pull request that adds a new agent can use it to verify the definition meets current standards.

This cross-plugin applicability is deliberate. The audit criteria and rewriting principles are not specific to AI -- they apply to any agent definition regardless of domain. The skill lives in ai-development because agent definition quality is fundamentally an AI development concern: writing effective prompts, choosing the right model, and structuring triggers for reliable activation.

## Related

- [Agent Reference]({{< ref "/reference/agents" >}}) -- ai-engineer specification and capabilities
- [Skill Reference]({{< ref "/reference/skills" >}}) -- agent-modernizer audit criteria and severity levels
- research plugin -- Gathering domain knowledge before ai-engineer implementations
