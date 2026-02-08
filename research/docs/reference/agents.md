---
title: "Agent Reference"
description: "Technical specifications for all five research plugin agents including models, tools, capabilities, and trigger examples"
weight: 1
---

# Agent Reference

The research plugin provides five agents. This page documents their specifications.

---

## academic-researcher

| Field | Value |
|-------|-------|
| **Model** | sonnet |
| **Color** | cyan |
| **Tools** | WebSearch, WebFetch, Read, Write |

### Description

Searches and evaluates scholarly literature from academic databases including ArXiv, Google Scholar, and PubMed. Assesses source quality by peer review status, citation count, journal reputation, methodology quality, sample size, and reproducibility. Uses multiple query variations including synonyms, related terms, and author names. Tracks citation chains from foundational papers to recent work.

### Confidence markers

- **Strong evidence** -- multiple replicated studies with consistent results
- **Moderate evidence** -- few studies, but consistent findings
- **Preliminary** -- single study or preprint, not yet replicated
- **Contested** -- conflicting findings across studies

### Citation format

(Author et al., Year) with enough detail to locate the source.

### Trigger examples

- "What does the research say about transformer architectures for time series forecasting?"
- "Is there solid evidence that retrieval-augmented generation reduces hallucination in LLMs?"
- "I need a literature review on differential privacy in federated learning"

### Key capabilities

- Systematic literature search with multiple query strategies
- Source quality evaluation (peer review, methodology, citation analysis)
- Evidence landscape mapping with confidence levels
- Gap identification: unstudied areas, weak methodology, unreplicated results
- Conflict analysis: explains why studies disagree (methodology, population, definitions, time period)

---

## comprehensive-researcher

| Field | Value |
|-------|-------|
| **Model** | sonnet |
| **Color** | blue |
| **Tools** | WebSearch, WebFetch, Read, Write |

### Description

Conducts deep multi-source investigations by decomposing topics into 5-8 specific sub-questions covering different angles (technical, economic, social, historical, forward-looking). Sources diversely from academic papers, government reports, industry analysis, expert commentary, and news reporting. Triangulates claims across multiple independent sources.

### Source prioritization (highest to lowest weight)

1. Peer-reviewed research and systematic reviews
2. Government and institutional reports with transparent methodology
3. Industry analysis from reputable firms (notes potential conflicts of interest)
4. Expert commentary and opinion (labeled as such)
5. News reporting (for timeliness, not as primary evidence)

### Output structure

- Executive summary: 3-5 bullet points
- Main body organized by theme or sub-question, every claim cited
- Confidence markers throughout
- Contradictions and debates presented with evidence for each side
- Knowledge gaps identified
- Full bibliography

### Trigger examples

- "Give me a comprehensive briefing on the current state of nuclear fusion energy"
- "What are the real trade-offs of adopting WebAssembly for our backend services?"
- "Research the current state of AI code generation tools and their impact on developer productivity"

### Key capabilities

- Topic decomposition into 5-8 answerable sub-questions
- Multi-source type coverage (academic, government, industry, expert, news)
- Claim triangulation across independent sources
- Evidence quality assessment with confidence markers
- Structured output with executive summary and themed findings

---

## technical-researcher

| Field | Value |
|-------|-------|
| **Model** | sonnet |
| **Color** | green |
| **Tools** | WebSearch, WebFetch, Read, Write, Bash, Grep |

### Description

Evaluates code repositories, libraries, frameworks, and implementation approaches. Assesses projects using a structured framework covering code quality, project health, community, API design, and performance. Sources include GitHub/GitLab repositories, package registries, technical documentation, developer forums, and benchmarks.

### Evaluation framework

- **Code quality** -- architecture patterns, test coverage, error handling, documentation quality, type safety
- **Project health** -- last commit date, open vs. closed issue ratio, PR review cadence, bus factor (active maintainer count)
- **Community** -- stars/forks (rough popularity, not quality), Stack Overflow activity, forum engagement, corporate backing
- **API design** -- ergonomics, consistency, breaking change history, migration path quality
- **Performance** -- published benchmarks (assessed with skepticism), known bottlenecks, scaling characteristics

### Defaults

- Compares at least 2-3 alternatives including the "do nothing" or "use stdlib" option
- Checks issue trackers for dealbreaker bugs and maintainer responsiveness
- Reviews breaking changes between recent major versions
- Notes license and commercial restrictions

### Trigger examples

- "I need to pick a Go HTTP router -- compare chi, gorilla/mux, and the standard library"
- "How does SQLite's WAL mode actually work under the hood?"
- "Is Drizzle ORM mature enough for production use?"

### Key capabilities

- Repository analysis (code, issues, PRs, release notes, commit history)
- Package registry data (download stats, version history, dependency count)
- Multi-alternative comparison with trade-offs
- Production readiness assessment
- Technical deep-dive on implementation internals

---

## research-synthesizer

| Field | Value |
|-------|-------|
| **Model** | opus |
| **Color** | yellow |
| **Tools** | Read, Write |

### Description

Consolidates findings from multiple sources or parallel research agents into a unified analysis. Organizes output by theme, not by source. Preserves disagreements rather than resolving them. Surfaces patterns, contradictions, and gaps that only become visible when looking across all inputs together. Works exclusively with provided input -- does not conduct its own research.

### Integration principles

- **Merge, not stack** -- output organized by theme, not by source
- **Preserve disagreement** -- contradictions are valuable information, not noise to be resolved
- **Attribute proportionally** -- single-source findings noted as such, convergent findings noted as corroborated
- **Surface gaps** -- unanswered questions across all research
- **Maintain confidence markers** -- strong consensus, emerging agreement, contested, insufficient evidence

### Output structure

- **Key findings** -- 3-5 most important conclusions with confidence levels
- **Thematic analysis** -- major themes by strength of evidence, with attribution
- **Unique insights** -- important single-source findings
- **Contradictions and tensions** -- where sources disagree, with evidence for each side
- **Evidence gaps** -- what was not covered, what needs further investigation
- **Source assessment** -- quality and coverage of the inputs received

### Trigger examples

- "I dispatched academic, technical, and market researchers on quantum computing -- synthesize their findings"
- "I have notes from 8 different sources on this topic -- help me make sense of it all"
- "We've done the research -- now tell me what it all means together"

### Key capabilities

- Cross-source thematic integration
- Contradiction and tension identification
- Evidence quality assessment across combined inputs
- Gap analysis spanning all source material
- Confidence level assignment based on corroboration

---

## report-generator

| Field | Value |
|-------|-------|
| **Model** | sonnet |
| **Color** | magenta |
| **Tools** | Read, Write |

### Description

Transforms research findings into polished, well-structured documents. Adapts structure, tone, and terminology to the target audience. Uses sequential citation numbering [1], [2] linked to a references section. Works exclusively with provided input -- a writer, not a researcher.

### Report structure

1. **Executive Summary** -- 3-5 bullet points (written last, placed first)
2. **Introduction** -- context, research scope, significance
3. **Key Findings** -- organized by theme or importance, every claim cited [n]
4. **Analysis** -- implications of findings for the reader
5. **Contradictions and Open Questions** -- conflicting evidence presented fairly
6. **Conclusion** -- takeaways, implications, suggested next steps
7. **References** -- sequential numbering, consistent format

### Audience adaptation

| Audience | Characteristics |
|----------|----------------|
| **Technical** | Precise terminology, methodology details, code examples |
| **Executive** | Lead with implications and decisions, minimal jargon, business impact |
| **Policy** | Actionable recommendations, regulatory context, stakeholder impact |
| **General** | Terms defined on first use, analogies, accessibility over precision |

### Citation format

Sequential numbering [1], [2], [3] with a references section at the end of the document.

### Trigger examples

- "Turn these research findings into a report for our engineering leadership"
- "Create a technical comparison report from this library evaluation data"
- "Package up all our AI safety research into a comprehensive report with proper citations"

### Key capabilities

- Audience-adaptive writing (technical, executive, policy, general)
- Sequential citation management
- Structured report formatting with Markdown
- Comparison tables for side-by-side analysis
- Length calibration: matches report length to content depth
