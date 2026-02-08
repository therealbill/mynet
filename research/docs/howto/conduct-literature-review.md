---
title: "Conduct a Literature Review"
description: "Steps to conduct an academic literature review using the academic-researcher agent with source evaluation and gap identification"
weight: 1
---

# Conduct a Literature Review

This guide walks through the steps to produce a structured academic literature review using the academic-researcher agent.

## Before you begin

Have a clear research question or topic. The more specific the question, the more focused the results. "Machine learning" is too broad; "transformer architectures for time series forecasting" gives the agent something to work with.

## Step 1: Dispatch the academic-researcher

Frame your request with a specific research question and any scope constraints:

```
Use the academic-researcher to conduct a literature review on
differential privacy techniques in federated learning.
Focus on papers from 2020 onward. I need to understand the
current state of the field and identify open problems.
```

The agent will:

- Start with recent survey and review papers to map the landscape
- Search ArXiv for preprints, Google Scholar for citation networks, and PubMed for biomedical topics
- Use multiple query variations (synonyms, related terms, author names from key papers)
- Track citation chains from foundational papers to recent work

## Step 2: Evaluate source quality with confidence markers

The academic-researcher assigns confidence markers to its findings. Review these carefully:

| Marker | Meaning |
|--------|---------|
| **Strong evidence** | Multiple replicated studies with consistent results |
| **Moderate evidence** | Few studies, but consistent findings |
| **Preliminary** | Single study or preprint, not yet replicated |
| **Contested** | Conflicting findings across studies |

When reviewing the output, pay attention to:

- **Peer review status** -- published in a peer-reviewed venue, or preprint only?
- **Citation count** -- highly cited papers carry more weight, but a high-citation paper with poor methodology is still weak evidence
- **Methodology quality** -- sample size, reproducibility, controls
- **Recency** -- older foundational work vs. cutting-edge findings

If the agent presents a finding as "strong evidence" but cites only one study, ask it to reassess.

## Step 3: Identify gaps and contradictions

A literature review is not a list of papers that agree. Ask the agent to explicitly surface:

- **Gaps**: areas where no studies exist, or where existing studies have methodological limitations
- **Contradictions**: cases where studies reach different conclusions, with an explanation of why (different populations, methodologies, definitions, or time periods)
- **Under-explored directions**: topics mentioned in papers' "future work" sections that have not been followed up

```
Identify the main gaps and contradictions in the literature you found.
Where do researchers disagree, and what hasn't been studied yet?
```

## Step 4: Produce the structured literature review

Ask the agent to organize its findings into a literature review structure:

```
Write a structured literature review from your findings.
Organize by theme, not chronologically. Include proper
(Author et al., Year) citations and identify the overall
state of the field.
```

The output should include:

- **Introduction** -- scope and significance of the research area
- **Thematic sections** -- findings grouped by sub-topic, each with citations and confidence levels
- **Contradictions section** -- where the field disagrees, with evidence for each position
- **Gaps and future directions** -- what remains unstudied or under-studied
- **References** -- complete citations in (Author et al., Year) format

## Step 5: Save and refine

Ask the agent to save the review to a file:

```
Write the literature review to research/lit-review-differential-privacy-fl.md
```

Review the output and iterate. Common refinements:

- Ask the agent to search for additional papers on a specific sub-topic that seems underrepresented
- Request deeper analysis of a contradiction between two key papers
- Ask for more detail on methodology assessment for the most-cited studies

## Combining with other agents

For a literature review that also evaluates software implementations:

1. Use academic-researcher for the scholarly literature
2. Use technical-researcher to evaluate any open-source implementations mentioned in the papers
3. Feed both outputs into research-synthesizer to merge academic and practical perspectives

See [Produce a Research Report]({{< ref "produce-research-report" >}}) for the full multi-agent workflow.
