---
title: "Architecture"
description: "Utility vs domain-specific agents and the URL analysis pipeline"
weight: 1
---

# Architecture

How the utility-agents plugin fits into the broader plugin ecosystem, why URL analysis is split into two agents, and how episode-orchestrator demonstrates a reusable workflow pattern.

## Utility Agents vs. Domain-Specific Agents

Most plugins in this marketplace are domain-specific. The web-development plugin contains agents that understand React, Next.js, and CSS. The backend-development plugin has agents for API design and SQL optimization. These agents carry deep knowledge of a particular technology stack and are useful only within that domain.

Utility agents are different. They serve any domain because they operate on universal patterns rather than technology-specific knowledge. URLs exist in every web project regardless of framework. Workflow orchestration applies to any multi-step process regardless of what the steps do. The utility-agents plugin provides capabilities that other plugins can leverage but should not duplicate.

This distinction matters for plugin design. A domain-specific agent like the web-development plugin's accessibility-auditor knows how to check ARIA attributes in React components -- it has domain knowledge that makes it effective in one context but irrelevant in others. The url-context-validator checks whether links work and make sense -- a task that applies equally to a React app's README, a Go project's API documentation, and a marketing site's landing pages.

## The URL Analysis Pipeline

URL analysis in this plugin follows a two-stage pipeline: extract, then validate.

**Stage 1: Extraction** -- url-link-extractor scans a codebase and produces a categorized inventory of every URL it finds. It reports file paths, line numbers, URL categories, and flagged patterns. It does not check whether any URL is reachable.

**Stage 2: Validation** -- url-context-validator takes URLs and checks each one for technical reachability and contextual appropriateness. It fetches destinations, compares content against link text, searches for better alternatives, and reports findings grouped by severity. It does not scan codebases for URLs.

The pipeline can be used end-to-end (extract everything, then validate everything) or each stage can be used independently. You might use url-link-extractor alone when preparing a migration inventory where you do not care about reachability -- you just need to know what to update. You might use url-context-validator alone when you already have a specific document whose links you want to check.

## Why Two URL Agents Instead of One

The split is driven by toolset requirements.

URL extraction is a file-scanning problem. The url-link-extractor needs Grep to match URL patterns across files, Glob to identify which file types to scan, and Bash for any custom filtering. It never needs to make network requests. Its job is entirely local: read files, find patterns, categorize results.

URL validation is a network and comparison problem. The url-context-validator needs WebFetch to load destination pages and WebSearch to find alternative sources. It never needs to scan files for URL patterns. Its job is entirely external: fetch pages, compare content, assess quality.

Combining these into a single agent would produce an agent with six tools (Read, Write, Grep, Glob, Bash, WebFetch, WebSearch) and a system prompt covering both file scanning patterns and content assessment criteria. The tool surface area grows, the prompt loses focus, and the agent becomes mediocre at both tasks instead of good at one.

The separation also enables independent use. A user who needs a URL inventory for migration planning calls the extractor alone. A user who wants to check five specific links calls the validator alone. Neither user loads capabilities they do not need.

## Episode-Orchestrator as a Workflow Pattern

The episode-orchestrator agent serves a different purpose than the URL agents. Where the URL agents operate as a pipeline (output of one feeds input of the other), the episode-orchestrator demonstrates a self-contained multi-step workflow pattern.

Its design embodies several principles that apply to any orchestration agent:

**Completeness detection before execution.** The agent defines a minimum viable payload -- a title plus at least one of duration, air date, or episode number. It checks for completeness before running any pipeline step. This prevents partial data from propagating through a workflow and producing confusing failures downstream.

**Single-question clarification.** When data is incomplete, the agent asks exactly one clarifying question rather than presenting a list of all missing fields. This keeps the interaction focused and avoids overwhelming the user with a form-like experience in a conversational interface.

**Sequential execution with output forwarding.** Pipeline steps run in order, and each step receives the outputs of all previous steps. This is a simple but effective orchestration pattern: no parallel execution complexity, no dependency graph resolution, just a sequence where context accumulates.

**Explicit failure reporting.** When a step fails, the agent reports which step failed, what went wrong, and whether remaining steps were skipped or continued. This is the minimum viable error reporting for any multi-step workflow -- the user always knows where the process stopped and why.

## Cross-Plugin Relationships

The utility-agents plugin intersects with other plugins at specific points:

**developer-tools** -- the documentation-engineer agent in the developer-tools plugin works at the document level, improving structure, clarity, and completeness. The url-context-validator works at the link level, checking whether references within documents are valid. These are complementary: documentation-engineer ensures the document is well-written, url-context-validator ensures the links within it are accurate.

**diataxis-docs** -- the diataxis-docs plugin validates documentation against the Diataxis framework (tutorials, how-to guides, explanation, reference). It assesses whether a document is in the right category and follows the right structure. The url-context-validator can then check that the links within a properly structured document point to appropriate destinations.

These are not dependencies -- no plugin requires another to function. They are points of natural composition where using multiple plugins together produces better results than using any one alone.
