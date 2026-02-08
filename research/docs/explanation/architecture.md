---
title: "Architecture"
description: "The research pipeline architecture: why five specialized agents, how they connect, and the role each plays"
weight: 1
---

# Research Pipeline Architecture

The research plugin uses five agents organized into a three-stage pipeline: gather, synthesize, report. This page explains why the plugin is structured this way and what role each agent plays.

## Why five agents instead of one

A single monolithic "research everything and write a report" agent would be simpler to invoke but worse at every individual task. Research, synthesis, and writing are different cognitive activities that benefit from different system prompts, tool sets, and even model selections.

**Different source types require different evaluation criteria.** Academic literature demands citation chain tracking, peer review status checks, and methodology assessment. Code repository evaluation requires examining commit history, issue tracker health, and API surface area. A general investigation needs to balance academic, industry, government, and news sources with explicit conflict-of-interest awareness. One set of evaluation instructions cannot serve all three well.

**Different tools for different jobs.** The academic-researcher and comprehensive-researcher need WebSearch and WebFetch to find and retrieve sources across the internet. The technical-researcher additionally needs Bash and Grep to inspect code repositories and run analysis commands. The research-synthesizer and report-generator need only Read and Write -- they work with material that has already been gathered. Giving every agent every tool would broaden the attack surface of each prompt and dilute focus.

**Separation enables parallel dispatch.** When a research question spans multiple domains, you can dispatch the academic-researcher, technical-researcher, and comprehensive-researcher simultaneously. Each runs with its own specialized prompt and tool set. A monolithic agent would have to handle these sequentially.

**Synthesis is a distinct skill.** Merging findings from multiple sources into a unified thematic analysis is fundamentally different from gathering those findings in the first place. The research-synthesizer uses the opus model precisely because integration -- finding cross-cutting themes, surfacing contradictions that span sources, assessing evidence quality holistically -- benefits from stronger reasoning capabilities.

**Writing is not researching.** The report-generator does not search for information. It takes findings that have already been gathered and synthesized, then structures them into a document appropriate for a specific audience. Separating the writer from the researcher prevents a common failure mode: an agent that fills gaps in its research by quietly generating plausible-sounding content rather than acknowledging the gaps.

## The three-stage pipeline

```
Stage 1: GATHER          Stage 2: SYNTHESIZE       Stage 3: REPORT
-------------------------+--------------------------+-----------------------
academic-researcher    --+                          |
comprehensive-researcher-+-> research-synthesizer --+-> report-generator
technical-researcher   --+                          |
```

### Stage 1: Gather

One or more researcher agents investigate the topic. Each agent produces findings in its domain with citations and confidence markers. The agents write their outputs to files.

- **academic-researcher** -- scholarly rigor, peer-reviewed sources, evidence evaluation
- **comprehensive-researcher** -- breadth across source types, topic decomposition, triangulation
- **technical-researcher** -- code and repository analysis, library comparison, project health

You choose which agents to dispatch based on the research question. A purely academic question might only need the academic-researcher. A technology decision might involve all three.

### Stage 2: Synthesize

The research-synthesizer reads all gathered findings and produces a unified analysis. It organizes by theme rather than by source, so the output reads as "the evidence on topic X shows..." rather than "agent A said X, agent B said Y."

The synthesizer preserves disagreements. When sources contradict each other, the contradiction appears in the output with evidence for each side. The synthesizer also identifies gaps -- questions that none of the gathering agents answered.

This stage runs on the opus model because integration across disparate sources demands stronger reasoning than individual source gathering.

### Stage 3: Report

The report-generator transforms the synthesis into a polished document. It adapts structure, tone, and terminology to the target audience (technical, executive, policy, or general). It uses sequential citations [1], [2] and produces a references section.

The report-generator does not add findings. If the synthesis has gaps, those gaps appear as acknowledged limitations in the report rather than being filled with generated content.

## Agent model and tool selections

| Agent | Model | Tools | Rationale |
|-------|-------|-------|-----------|
| academic-researcher | sonnet | WebSearch, WebFetch, Read, Write | Web access for source discovery; sonnet for efficient search and evaluation |
| comprehensive-researcher | sonnet | WebSearch, WebFetch, Read, Write | Same tool needs as academic; sonnet handles structured decomposition well |
| technical-researcher | sonnet | WebSearch, WebFetch, Read, Write, Bash, Grep | Additional tools for code inspection; sonnet for systematic technical analysis |
| research-synthesizer | opus | Read, Write | No web access (works with gathered material only); opus for complex integration reasoning |
| report-generator | sonnet | Read, Write | No web access (works with synthesized material only); sonnet for fluent writing |

The research-synthesizer is the only agent on the opus model. This reflects that cross-source integration -- finding themes that span multiple research streams, identifying subtle contradictions, and assessing evidence quality holistically -- is the most reasoning-intensive task in the pipeline.

## Design constraints

**Each agent has a single responsibility.** Gathering agents gather. The synthesizer synthesizes. The reporter writes. This makes each agent's behavior predictable and its output verifiable.

**No agent researches and writes.** This is the most important constraint. Combining research and writing in one agent creates an incentive to generate plausible content when research comes up short. Separating the two forces gaps to be surfaced rather than papered over.

**The synthesizer does not summarize.** Summarization loses information. Synthesis integrates information -- finding patterns and conflicts that only become visible when looking across all inputs together. The output of synthesis may be longer than any single input.
