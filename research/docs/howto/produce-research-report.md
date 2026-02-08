---
title: "Produce a Research Report"
description: "Steps to gather findings from researcher agents, synthesize them, and produce a polished report for any audience"
weight: 2
---

# Produce a Research Report

This guide covers the steps to go from raw research questions to a polished, audience-appropriate report using the research plugin's synthesis and reporting agents.

## Before you begin

Decide what you are researching and for whom. The research question determines which agents to dispatch. The audience determines how the report-generator structures the final output.

## Step 1: Gather findings from researcher agents

Choose the researcher agents that match your topic. You can dispatch multiple agents in parallel for broader coverage.

| Agent | Best for |
|-------|----------|
| academic-researcher | Peer-reviewed literature, citation analysis, evidence evaluation |
| comprehensive-researcher | Broad multi-source investigation across academic, industry, news, government |
| technical-researcher | Code repositories, library comparison, API evaluation, project health |

For a single-source investigation, dispatch one agent:

```
Use the comprehensive-researcher to investigate the trade-offs
of adopting event-driven architecture for our payment processing system.
Write findings to research/eda-findings.md
```

For broader coverage, dispatch multiple agents and have each write to a separate file:

```
Use the academic-researcher to research event-driven architecture
patterns and their reliability characteristics. Write findings to
research/eda-academic.md

Use the technical-researcher to evaluate Kafka, Pulsar, and NATS
for our event streaming needs. Write findings to research/eda-technical.md

Use the comprehensive-researcher to investigate production adoption
of event-driven architectures in financial services. Write findings to
research/eda-comprehensive.md
```

## Step 2: Synthesize multi-source findings

Once all researcher agents have completed, use the research-synthesizer to merge their outputs. The synthesizer organizes by theme rather than by source, surfaces contradictions, and identifies evidence gaps.

```
Use the research-synthesizer to consolidate the findings in
research/eda-academic.md, research/eda-technical.md, and
research/eda-comprehensive.md. Write the synthesis to
research/eda-synthesis.md
```

The synthesizer produces:

- **Key findings** with confidence levels
- **Thematic analysis** organized by strength of evidence
- **Unique insights** from individual sources
- **Contradictions** where sources disagree
- **Evidence gaps** where more investigation is needed

If you dispatched only a single researcher agent, you can skip synthesis and feed its output directly to the report-generator. Synthesis is most valuable when merging findings from multiple agents or sources.

## Step 3: Generate the report

Use the report-generator to transform the synthesis (or raw findings) into a polished document. Specify the target audience so the agent adapts tone, terminology, and structure.

```
Use the report-generator to create a report from research/eda-synthesis.md.
The audience is our VP of Engineering who needs to decide whether to
invest in event-driven architecture for payment processing.
Write the report to research/eda-report.md
```

### Audience adaptation

The report-generator adjusts its output based on the audience you specify:

- **Technical audience** -- precise terminology, methodology details, code examples
- **Executive audience** -- lead with implications and decisions, minimize jargon, emphasize business impact
- **Policy audience** -- actionable recommendations, regulatory context, stakeholder impact
- **General audience** -- define terms on first use, use analogies, prioritize accessibility

### Citation format

The report-generator uses sequential citations [1], [2], [3] linked to a references section at the end. Every claim in the report maps to a numbered source.

## Step 4: Review and iterate

Check the final report for:

- **Citation coverage** -- every substantive claim should have a [n] citation
- **Audience match** -- language and structure appropriate for the specified audience
- **Preserved contradictions** -- disagreements from the synthesis should appear in the report, not be smoothed away
- **Proportional length** -- the report should match its content, neither padding thin findings nor compressing rich ones

Common refinements:

- Ask the report-generator to expand the executive summary if it is too terse
- Request a comparison table for technical alternatives
- Ask for a "recommended next steps" section tailored to the audience

## Quick reference: the full pipeline

```
Step 1: Dispatch researcher agent(s) --> findings files
Step 2: research-synthesizer          --> synthesis file
Step 3: report-generator              --> final report
```

For simple investigations, steps 1 and 3 may suffice. For complex topics with multiple source types, the three-stage pipeline produces the most thorough and well-organized result.

## Related guides

- [Conduct a Literature Review](../conduct-literature-review/) -- academic-focused workflow with source evaluation
- [Architecture](../../explanation/architecture/) -- why the pipeline is structured this way
- [Agent Reference](../../reference/agents/) -- specifications for all five agents
