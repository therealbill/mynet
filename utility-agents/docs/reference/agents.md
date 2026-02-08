---
title: "Agents"
description: "Technical specifications for all utility-agents"
weight: 1
---

# Agents

Technical specifications for the three agents in the utility-agents plugin.

## episode-orchestrator

| Field | Value |
|-------|-------|
| Name | episode-orchestrator |
| Model | sonnet |
| Color | blue |
| Tools | Read, Write |

**Trigger conditions:**

- Processing episode data through a multi-step pipeline
- Validating episode payloads for completeness
- Coordinating sequential workflow steps with output forwarding
- Handling batch episode processing with per-episode validation

**Trigger patterns:**

- "Process this episode: Title: Pilot, Duration: 42min"
- "I need to process the season 2 premiere"
- "Run these three episodes through the pipeline"

**Completeness detection:**

An episode payload is complete when it contains a title and at least one of:

- Duration
- Air date
- Episode number

If the payload is incomplete, the agent asks exactly one clarifying question.

**Process:**

1. Parse the incoming request to identify episode data
2. Validate that required fields are present and have reasonable values
3. If incomplete, ask one focused clarifying question
4. Execute pipeline steps in their configured sequence
5. Pass outputs forward so later steps can use earlier results
6. Return a consolidated summary of all step results

**Constraints:**

- Does not invent episode data that was not provided
- Does not skip validation to proceed faster
- Does not silently drop errors from pipeline steps
- Does not reorder the configured processing sequence

---

## url-context-validator

| Field | Value |
|-------|-------|
| Name | url-context-validator |
| Model | sonnet |
| Color | cyan |
| Tools | Read, Write, WebFetch, WebSearch |

**Trigger conditions:**

- Checking whether links in a document work and point to appropriate content
- Validating references in articles, READMEs, or documentation pages
- Finding dead, misaligned, or outdated links
- Assessing content freshness behind working URLs

**Trigger patterns:**

- "Check links in my README"
- "Validate the references in this article"
- "Are any links dead or pointing to outdated docs?"

**Check types:**

| Check | Description |
|-------|-------------|
| Reachability | Uses WebFetch to confirm the URL loads successfully |
| Redirects | Notes when a URL redirects; records both original and final destination |
| Anchor text alignment | Compares link text against destination content |
| Content freshness | Assesses whether linked content is outdated even if the URL works |
| Better alternatives | Uses WebSearch to find more authoritative or current sources for stale links |
| Security concerns | Flags HTTP links in HTTPS contexts or suspicious domains |

**Severity classification:**

| Severity | Meaning |
|----------|---------|
| Broken | URL returns an error (404, 500, connection refused) |
| Misaligned | URL works but destination content does not match link text or surrounding context |
| Stale | URL works and is contextually reasonable but content is outdated |
| Fine | URL works, matches context, points to current content |

**Report fields per link:**

- Status (working, broken, redirect)
- Whether destination content matches the surrounding context
- Recommended action (keep, update URL, replace with better source, remove)
- Suggested replacement URL when applicable

**Process:**

1. Extract all URLs from the provided content using Read
2. Fetch each URL with WebFetch and record the result
3. For working links, compare destination content against surrounding context and anchor text
4. Report findings grouped by severity: broken, misaligned, stale, fine

**Constraints:**

- Does not claim to measure response times -- WebFetch does not expose this
- Does not claim to detect regional access restrictions
- Does not skip contextual analysis for working links
- Does not report only technical status without contextual assessment

---

## url-link-extractor

| Field | Value |
|-------|-------|
| Name | url-link-extractor |
| Model | sonnet |
| Color | green |
| Tools | Read, Write, Grep, Glob, Bash |

**Trigger conditions:**

- Building a comprehensive URL inventory from a codebase
- Preparing for domain migration by finding all internal URLs
- Auditing external dependencies referenced in code
- Finding hardcoded development or localhost URLs in production code

**Trigger patterns:**

- "Find every URL for domain migration"
- "List all external URLs"
- "Any hardcoded localhost URLs in the code?"

**URL categories:**

| Category | Examples |
|----------|----------|
| Navigation and content | href attributes, markdown links, router paths, internal page references |
| Assets and resources | Image sources, script tags, stylesheets, fonts, media files, CDN references |
| APIs and services | fetch/axios URLs, API base URLs, webhook endpoints, environment variable URLs |
| External references | Third-party links, social media URLs, mailto/tel links, documentation references |

**Flagged patterns:**

- Hardcoded localhost or 127.0.0.1 URLs
- Duplicate URLs across multiple files
- Inconsistent base paths (same API at different versions)
- HTTP URLs where HTTPS is expected

**Output fields per URL:**

- The URL
- Category (navigation, asset, API, external)
- File path where the URL appears
- Line number within the file
- Count of unique vs. total URLs
- Flagged items requiring attention

**Process:**

1. Use Glob to identify relevant file types (HTML, JS, TS, CSS, MD, JSON, YAML, config files)
2. Use Grep with URL-matching patterns to extract URLs from each file type
3. Categorize each URL by type (internal vs. external) and purpose (navigation, asset, API, external)
4. Flag problematic patterns: hardcoded localhost, duplicates, inconsistent base paths
5. Report findings as a structured inventory with file paths and line numbers

**Constraints:**

- Does not validate whether URLs are reachable -- that is url-context-validator's responsibility
- Does not report URLs found in node_modules, vendor, or other dependency directories
- Does not omit the source file and line number -- location context is required for every entry

---

## Agent Comparison

| Agent | Model | Scope | Primary Output |
|-------|-------|-------|----------------|
| episode-orchestrator | sonnet | Workflow orchestration with validation | Consolidated pipeline execution summary |
| url-context-validator | sonnet | URL validation (functional + contextual) | Severity-grouped link health report |
| url-link-extractor | sonnet | URL extraction and cataloging | Categorized URL inventory with source locations |

## Pipeline Relationship

The url-link-extractor and url-context-validator agents form a natural pipeline:

1. **url-link-extractor** scans the codebase and produces a categorized URL inventory with file locations
2. **url-context-validator** takes URLs and checks each one for reachability, contextual alignment, and freshness

Extraction uses Grep and Glob to scan files. Validation uses WebFetch and WebSearch to check destinations. The agents use different toolsets because they solve different problems.
