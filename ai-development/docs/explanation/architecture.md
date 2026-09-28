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

As the marketplace grows, agent quality becomes a scaling problem. A single plugin author can maintain consistency across three agents. Across fifty plugins from different authors, definitions drift: some have verbose prompts that teach the model its own knowledge, others carry example blocks inside the description and inflate every turn's routing context, others describe themselves in summary form and trigger unreliably, others reference agents that do not exist.

The agent-modernizer skill addresses this by codifying audit criteria into a repeatable process. A script checks frontmatter completeness and measures lengths, a catalog of typed questions evaluates system prompt quality against known anti-patterns, and rules compose the two into a findings table with severity levels. This transforms agent quality from a subjective judgment ("this prompt feels too long") into a structured evaluation ("34 bullet points; topic lists without guidance at p=0.81: Should fix").

## Code Owns the Checks, Typed Questions Own the Judgments

Version 1.1.0 restructured the audit around a principle borrowed from TypeSafe's System One model: keep deterministic work in code and hand the model only narrow, typed judgments. Everything that can be measured is measured by `scripts/audit-agents.py` -- field presence, name format, model alias, example structure and placement, character and word counts, bullet counts. A model asked to count bullets or compare lengths gives unreliable answers; code gives the same answer every run.

What remains are judgments that need reading comprehension: does this list name topics without decisions, does this body reference an agent that does not exist, does this description read as selection criteria or as a summary. Each is a single question in `references/question-catalog.json` with explicit criteria for yes and no. Jev answers them as probabilities in one request per agent; when no key is present, the same questions are answered by reading the file. Either way the questions are the contract, so two audits of the same agent disagree only where the text is genuinely ambiguous.

Severity and the recommended action are then composed by rules in code, not by the model. That keeps the raw judgments reusable: a threshold can be changed without re-running inference, and a probability near 0.5 is surfaced as "possible, confirm by reading" rather than rounded to a verdict. The thresholds shipped with the skill were set by running the catalog over the 64 agents in this marketplace; the split of the original "redundancy" question into two narrower ones came from observing that the combined question sat near 0.5 for most agents.

The same discipline gives the skill two checks that a checklist could not: the routing check asks, for every example request in a plugin, which sibling description a router would pick, and the rewrite comparison asks whether a rewrite dropped knowledge the model could not infer. Both are single typed questions over state that code assembles.

## Why a Skill, Not an Agent

The agent-modernizer is a skill rather than an agent because it provides knowledge injection -- audit criteria, anti-pattern definitions, rewriting principles -- rather than autonomous multi-step behavior. It enhances whatever agent is handling the conversation by injecting the specific knowledge needed to evaluate and rewrite agent definitions.

An agent would be appropriate if the modernizer needed to autonomously discover agent files, run tests, and iterate on rewrites without user guidance. The current design assumes the user drives the process: they choose which agent to audit, review the findings, and decide whether to apply a rewrite. The skill provides the expertise; the user provides the judgment.

## Cross-Plugin Value

The agent-modernizer works on agents from any plugin in the marketplace, not just agents within ai-development. This makes it a meta-tool for marketplace quality. A plugin author working on the devops plugin can use agent-modernizer to audit their CI/CD agents. A maintainer reviewing a pull request that adds a new agent can use it to verify the definition meets current standards.

This cross-plugin applicability is deliberate. The audit criteria and rewriting principles are not specific to AI -- they apply to any agent definition regardless of domain. The skill lives in ai-development because agent definition quality is fundamentally an AI development concern: writing effective prompts, choosing the right model, and structuring triggers for reliable activation.

## Related

- [Agent Reference](../../reference/agents/) -- ai-engineer specification and capabilities
- [Skill Reference](../../reference/skills/) -- agent-modernizer audit criteria and severity levels
- research plugin -- Gathering domain knowledge before ai-engineer implementations
