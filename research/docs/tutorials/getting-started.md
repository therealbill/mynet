---
title: "Getting Started with Research"
description: "Walk through a complete research investigation using the three-stage pipeline: gather, synthesize, report"
weight: 1
---

# Getting Started with Research

In this tutorial you will run a complete research investigation using the research plugin's three-stage pipeline. By the end, you will have gathered findings from multiple sources, consolidated them into a unified analysis, and produced a polished report ready for your audience.

## What you will build

A research report on a topic of your choosing, produced through three stages:

1. **Gather** -- dispatch the comprehensive-researcher to investigate broadly
2. **Synthesize** -- use the research-synthesizer to consolidate findings
3. **Report** -- use the report-generator to produce a polished deliverable

## Prerequisites

- The research plugin installed and available in your Claude Code session
- A topic you want to investigate (this tutorial uses "the current state of WebAssembly for server-side applications" as an example)

## Stage 1: Gather findings with comprehensive-researcher

Start by dispatching the comprehensive-researcher agent. This agent decomposes your topic into 5-8 sub-questions, searches across academic, industry, government, and news sources, and triangulates claims across independent sources.

Ask Claude to use the comprehensive-researcher agent:

```
Research the current state of WebAssembly for server-side applications.
Cover performance characteristics, ecosystem maturity, production adoption,
and comparison with native compiled languages.
```

The comprehensive-researcher will:

- Break your topic into sub-questions covering different angles (technical feasibility, performance benchmarks, ecosystem maturity, production case studies, limitations)
- Search across multiple source types for each sub-question
- Rate evidence quality using confidence markers: "strong evidence", "moderate evidence", "preliminary", "contested"
- Produce a structured report with an executive summary, themed findings, and bibliography

**Checkpoint:** Before moving on, verify that the output includes:

- An executive summary with 3-5 bullet points
- Findings organized by theme with citations
- Confidence markers on key claims
- A bibliography section at the end

If any of these are missing, ask the agent to expand on the incomplete sections before proceeding.

### Saving findings to a file

Ask the researcher to write its findings to a file so subsequent agents can read them:

```
Write your findings to research/wasm-server-findings.md
```

This gives the synthesizer and report-generator a concrete file to work from.

## Stage 2: Consolidate with research-synthesizer

Now bring in the research-synthesizer agent. This agent reads the gathered findings and produces a unified analysis organized by theme rather than by source. Its job is integration, not summarization.

```
Use the research-synthesizer to consolidate the findings in
research/wasm-server-findings.md. Identify cross-cutting themes,
contradictions, and evidence gaps.
```

The synthesizer will produce:

- **Key findings** -- the 3-5 most important conclusions with confidence levels
- **Thematic analysis** -- major themes organized by strength of evidence
- **Unique insights** -- findings that add important perspective but come from a single source
- **Contradictions and tensions** -- where sources disagree
- **Evidence gaps** -- what remains unanswered

**Checkpoint:** Verify the synthesizer output contains all five sections above. Pay special attention to the contradictions section -- if all findings appear to agree perfectly, the synthesis may have smoothed over genuine disagreements. Ask the synthesizer to look again if the output seems artificially harmonious.

Ask the synthesizer to save its output:

```
Write the synthesis to research/wasm-server-synthesis.md
```

## Stage 3: Generate the report with report-generator

The final stage transforms the synthesized findings into a polished document tailored to your audience. The report-generator is a writer, not a researcher -- it works exclusively with the material you provide.

Specify your audience when invoking the agent:

```
Use the report-generator to create a technical report from
research/wasm-server-synthesis.md. The audience is engineering leadership
evaluating whether to adopt WebAssembly for backend services.
```

The report-generator will:

- Write an executive summary (written last, placed first)
- Provide an introduction with context and scope
- Present key findings organized by theme with sequential citations [1], [2]
- Connect findings to implications for the reader
- Present contradictions and open questions honestly
- Include a conclusion with takeaways and next steps
- Append a numbered references section

**Checkpoint:** Verify the final report:

- Every claim has a citation marker [n] that maps to the references section
- The tone matches your specified audience (technical but leadership-oriented in this example)
- Contradictions from the synthesis are preserved, not hidden
- The report length matches the content -- it should not pad thin findings or compress rich ones

## The complete pipeline at a glance

```
comprehensive-researcher  -->  research-synthesizer  -->  report-generator
      (gather)                    (synthesize)              (report)
```

Each stage has a distinct role:

- **Gathering** produces raw, cited findings from diverse sources
- **Synthesis** finds patterns, contradictions, and gaps across those findings
- **Reporting** transforms the synthesis into a readable document for a specific audience

## Extending the pipeline

For deeper investigations, you can dispatch multiple researcher agents in parallel before synthesis:

- Use **academic-researcher** for peer-reviewed literature
- Use **technical-researcher** for code repository evaluation
- Use **comprehensive-researcher** for broad multi-source coverage

Feed all their outputs into the research-synthesizer together. The synthesizer is designed to merge findings from multiple agents, organizing by theme rather than by source. See the [research pipeline explanation](../../explanation/research-pipeline/) for details on parallel dispatch.

## Next steps

- [Conduct a Literature Review](../../howto/conduct-literature-review/) -- use the academic-researcher for scholarly sources
- [Produce a Research Report](../../howto/produce-research-report/) -- detailed steps for the synthesis-to-report workflow
- [Agent Reference](../../reference/agents/) -- specifications for all five agents
