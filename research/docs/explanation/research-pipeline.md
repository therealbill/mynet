---
title: "Research Pipeline"
description: "Why the research plugin uses specialized agents, how parallel dispatch works, and how synthesis and reporting transform raw findings into deliverables"
weight: 2
---

# The Research Pipeline in Depth

This page explores the reasoning behind the research plugin's pipeline design: why agents specialize, how to dispatch them in parallel, what synthesis actually does, and how audience adaptation works in report generation.

## Why specialization matters

Different source types demand different evaluation criteria, and a single set of instructions cannot serve all of them well.

**Academic sources** require tracking citation chains, assessing peer review status, evaluating methodology (sample size, controls, reproducibility), and distinguishing established consensus from preliminary findings. The academic-researcher is tuned for this: it starts with survey papers to map a field, then drills into primary sources, and applies confidence markers based on replication status.

**Code repositories and libraries** require examining commit history, issue tracker health, PR review cadence, breaking change history, and API design quality. The technical-researcher has Bash and Grep in its tool set specifically so it can inspect repositories directly. Its evaluation framework covers code quality, project health, community signals, and performance characteristics. It defaults to comparing at least 2-3 alternatives including the "do nothing" option.

**Broad investigations** require balancing source types with awareness of their relative strengths and limitations. The comprehensive-researcher decomposes topics into sub-questions and sources diversely -- academic papers carry the most weight, but government reports, industry analysis, expert commentary, and news reporting each contribute different perspectives. It explicitly labels the source type of each finding and notes potential conflicts of interest.

These three evaluation frameworks are different enough that combining them into one prompt would either make the prompt unwieldy or force compromises that weaken all three.

## Parallel dispatch

The most powerful usage pattern is dispatching multiple researcher agents simultaneously on the same topic. Each agent explores the topic through its own lens:

```
academic-researcher   -->  scholarly evidence with confidence markers
technical-researcher  -->  repository analysis with health metrics
comprehensive-researcher -->  multi-source investigation with triangulation
```

These agents run independently -- they do not read each other's output. Each writes its findings to a separate file. This independence is a feature, not a limitation: it means each agent applies its full evaluation framework without being influenced by findings from a different methodology.

The parallel dispatch pattern works because the agents cover genuinely different source types. The academic-researcher might find peer-reviewed studies showing that a technology works well in controlled experiments. The technical-researcher might discover that the reference implementation has significant unresolved bugs. The comprehensive-researcher might find industry reports showing mixed adoption results. These perspectives are complementary, and their occasional contradictions are informative.

## Synthesis as the critical step

The research-synthesizer is the integration layer of the pipeline. It takes the separate outputs from gathering agents and produces a unified analysis organized by theme, not by source. This distinction matters.

**Source-organized output** reads: "The academic researcher found X. The technical researcher found Y. The comprehensive researcher found Z." This is concatenation, not synthesis.

**Theme-organized output** reads: "On the question of reliability, academic studies show strong evidence of X [from academic-researcher], though the reference implementation has open issues affecting Y [from technical-researcher], and production adopters report Z [from comprehensive-researcher]."

The synthesizer identifies four categories of findings:

- **Corroborated findings** -- themes that appear across multiple inputs, strengthened by convergence
- **Unique insights** -- findings from a single source that add important perspective but lack independent corroboration
- **Contradictions** -- cases where sources disagree, presented with evidence for each side and possible explanations for the disagreement
- **Gaps** -- questions that none of the inputs addressed, either because they were not asked or because evidence is unavailable

The synthesizer runs on the opus model because this kind of cross-source integration -- holding multiple research streams in context, finding non-obvious connections, and identifying where apparent agreement masks a subtle disagreement -- requires stronger reasoning than individual source gathering.

## Report generation: audience adaptation and citation management

The report-generator is the final stage. It is a writer, not a researcher. This constraint is fundamental: the report-generator does not search for new information, and if the synthesis has gaps, those gaps appear as acknowledged limitations in the report.

### Audience adaptation

The same synthesis can produce very different reports depending on the target audience:

- A **technical audience** gets precise terminology, methodology details, and code examples where relevant
- An **executive audience** gets implications and decisions first, minimal jargon, and emphasis on business impact
- A **policy audience** gets actionable recommendations, regulatory context, and stakeholder impact analysis
- A **general audience** gets terms defined on first use, analogies for complex concepts, and accessibility prioritized over precision

The audience does not change the findings -- it changes how those findings are presented. A contradiction in the evidence appears in all four versions; it is explained differently depending on the reader.

### Citation management

The report-generator uses sequential citation numbering [1], [2], [3] linked to a references section. Every substantive claim in the report maps to a numbered source. This is a stricter standard than the gathering agents use (which cite inline in various formats) because the report is the deliverable that the reader will share, reference, and evaluate.

## Connecting to other plugins

Research findings often feed into downstream work. The output of the research pipeline -- a report file in Markdown with structured findings and citations -- can serve as input to other plugin workflows. For example, technical research findings about library choices can inform architecture decisions, and academic literature reviews can ground design documentation in established evidence.

The research plugin produces artifacts (Markdown files with structured content and citations) that are designed to be consumed by both humans and other agents. The file-based handoff between pipeline stages (gather to files, synthesize from files, report from files) means any stage's output is inspectable and reusable.
