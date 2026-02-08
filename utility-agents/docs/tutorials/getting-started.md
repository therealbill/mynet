---
title: "Getting Started with Utility Agents"
description: "Extract and validate URLs in a project using url-link-extractor and url-context-validator"
weight: 1
---

# Getting Started with Utility Agents

Walk through extracting every URL from a project and then validating them for correctness and relevance using the url-link-extractor and url-context-validator agents.

## What You'll Learn

- Use url-link-extractor to build a categorized URL inventory from a project
- Read the extraction output and understand its categories
- Identify flagged problems in the inventory
- Feed extracted URLs to url-context-validator for validation
- Interpret validation results grouped by severity

## Prerequisites

- The utility-agents plugin installed in your Claude Code environment
- A project that contains markdown files, HTML pages, or web application code with links -- any project with at least a handful of URLs will work

## Step 1: Extract All URLs

Start by asking the url-link-extractor agent to scan your project:

```
Find every URL in this project.
```

The agent activates automatically. It uses Glob to identify relevant file types -- markdown, HTML, JavaScript, TypeScript, CSS, JSON, YAML, and configuration files. It then runs Grep with URL-matching patterns across those files to extract every URL reference.

The agent skips `node_modules`, `vendor`, and other dependency directories automatically. It only reports URLs in your own source files.

### Checkpoint

You should see the agent scanning file types in sequence, building up a list of URLs with their source file paths and line numbers.

## Step 2: Review the Categorized Output

The extraction results are organized into four categories:

- **Navigation and content** -- href attributes in HTML, markdown links, router paths, and internal page references. These are the links your users click on.
- **Assets and resources** -- image sources, script tags, stylesheet links, font references, and CDN URLs. These are the files your application loads.
- **APIs and services** -- fetch or axios URLs, API base URLs, webhook endpoints, and URLs stored in environment variable references. These are the services your application talks to.
- **External references** -- third-party documentation links, social media URLs, mailto and tel links. These point outside your project.

Each entry includes the file path and line number where the URL appears, so you can locate it immediately.

### Checkpoint

Review the category breakdown. URLs should appear in the category that matches their purpose. A documentation link to an external API guide should be under "External references," not under "APIs and services."

## Step 3: Note Flagged Problems

The url-link-extractor flags patterns that commonly indicate issues:

- **Hardcoded localhost URLs** -- `http://localhost:3000` or `http://127.0.0.1:8080` left in production code
- **Duplicate URLs** -- the same URL referenced from many locations, which creates maintenance burden during migrations
- **Inconsistent base paths** -- the same API referenced as both `https://api.example.com/v1` and `https://api.example.com/v2` in different files

Review the flagged items. These are not necessarily errors -- a development configuration file legitimately contains localhost URLs -- but they are worth inspecting.

### Checkpoint

You should see a summary of flagged items at the end of the extraction report. If nothing is flagged, your URL hygiene is already good.

## Step 4: Validate the Extracted URLs

Now hand the content to url-context-validator to check whether the URLs actually work and point to appropriate content:

```
Check the links in this project for broken or outdated references.
```

The url-context-validator agent takes over. For each URL, it performs two levels of checking:

1. **Technical validation** -- it uses WebFetch to confirm the URL loads successfully. It notes redirects, recording both the original URL and the final destination.
2. **Contextual validation** -- for working links, it compares the destination content against the surrounding text. A link labeled "React documentation" that points to a Vue.js page is technically working but contextually wrong.

The agent also uses WebSearch to look for better alternatives when a linked resource appears stale -- for example, if a tutorial link points to a version 2 guide when version 4 is current.

### Checkpoint

The agent should be fetching each URL and reporting results. This step takes longer than extraction because it makes network requests.

## Step 5: Review Validation Results

The validation report groups findings by severity:

- **Broken** -- URLs that return errors (404, 500, connection refused). These need immediate attention: replace with a working URL or remove the link entirely.
- **Misaligned** -- URLs that work but point to content that does not match the anchor text or surrounding context. Update the anchor text or find a more appropriate link target.
- **Stale** -- URLs that work and are contextually correct but point to outdated content. The agent may suggest a more current alternative found via WebSearch.
- **Fine** -- URLs that work, match their context, and point to current content. No action needed.

For each finding, the agent provides a recommended action: keep, update URL, replace with a better source, or remove.

### Checkpoint

Review the broken and misaligned findings first -- these are the highest-priority items. Stale links are lower urgency but worth addressing before they become broken.

## Summary

You used two agents in sequence to build a complete URL health assessment:

- **url-link-extractor** scanned the codebase and produced a categorized inventory with file paths and line numbers. It flagged problematic patterns like hardcoded localhost URLs. It did not check whether any URL was reachable -- that is not its job.
- **url-context-validator** took the URLs and verified each one for both technical functionality and contextual appropriateness. It checked not just "does this link work?" but "does this link point to what the reader expects?"

This extract-then-validate pipeline separates concerns cleanly: extraction needs Grep and Glob to scan files, while validation needs WebFetch and WebSearch to check destinations. Different tools, different agents.

## Next Steps

- [Validate URLs in Documentation]({{< relref "../howto/validate-urls-in-documentation" >}}) -- focused procedure for documentation link checks
- [Extract All URLs from Codebase]({{< relref "../howto/extract-all-urls-from-codebase" >}}) -- targeted extraction for migration or audit
- [Agent Reference]({{< relref "../reference/agents" >}}) -- full specifications for all three agents
- [Architecture]({{< relref "../explanation/architecture" >}}) -- why two URL agents instead of one
