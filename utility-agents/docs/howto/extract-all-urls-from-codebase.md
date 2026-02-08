---
title: "Extract All URLs from Codebase"
description: "Use url-link-extractor to build a URL inventory for migration or audit"
weight: 2
---

# Extract All URLs from Codebase

Build a comprehensive URL inventory from your codebase for domain migration, security audit, or dependency review.

## Problem

You need to know every URL in your project. Maybe you are migrating to a new domain and need to update all internal references. Maybe you are auditing external dependencies before a security review. Maybe you suspect hardcoded development URLs are still in production code. In all cases, you need a complete, categorized list with source locations.

## Solution

Use the url-link-extractor agent to scan the codebase and produce a structured URL inventory grouped by category, with file paths and line numbers for every occurrence.

## Prerequisites

- The utility-agents plugin installed
- A codebase to scan (any language or framework)

## Steps

### 1. Trigger the Extraction

Tell the agent what you need:

```
Find every URL in this project for our domain migration.
```

Or for a security-focused extraction:

```
Are there any hardcoded localhost or dev URLs in the code?
```

The url-link-extractor agent activates. It uses Glob to identify relevant file types -- HTML, JavaScript, TypeScript, CSS, markdown, JSON, YAML, and configuration files. It then runs Grep with URL-matching patterns across each file type to extract every URL reference. Dependency directories like `node_modules` and `vendor` are excluded automatically.

### 2. Review the Categories

The inventory is organized into four categories:

- **Navigation and content** -- href attributes, markdown links, router paths, internal page references. For a domain migration, these are the URLs that need updating.
- **Assets and resources** -- image sources, script tags, stylesheet references, CDN URLs. These may point to your current CDN or asset domain.
- **APIs and services** -- fetch URLs, API base paths, webhook endpoints, environment variable URL references. These are the backend connections your application depends on.
- **External references** -- third-party documentation, social media links, mailto and tel links. These typically do not need updating during a migration but are useful for dependency auditing.

Each entry lists the URL, the file path, and the line number.

### 3. Filter for Your Specific Need

The full inventory serves different purposes depending on your goal:

**For domain migration** -- focus on Navigation and Assets categories. These contain the URLs that reference your current domain. External references and API endpoints on third-party domains do not need updating.

**For security audit** -- focus on APIs and flagged items. Look for hardcoded localhost URLs, HTTP (non-HTTPS) endpoints, and URLs containing credentials or tokens in query parameters.

**For dependency review** -- focus on External references and Assets. These show which third-party services and CDNs your project depends on.

### 4. Export the Inventory

The agent outputs results as a structured report. Use this output as a checklist for your migration or audit. For large projects, you may want to save the output to a file:

```
Extract all URLs and write the inventory to url-inventory.md.
```

The agent uses Write to save the categorized report to the specified file.

### 5. Act on the Results

For a domain migration, systematically update each URL from the old domain to the new one. The file paths and line numbers in the inventory point you directly to each location.

For a security audit, investigate each flagged item. Hardcoded localhost URLs in production code are a potential misconfiguration. HTTP endpoints may need upgrading to HTTPS.

## Verification

Cross-reference the inventory against URLs you already know about. If your project has a documented list of API endpoints or CDN URLs, confirm they all appear in the extraction results. Missing entries may indicate URLs constructed dynamically at runtime, which static extraction cannot detect.

After completing these steps you should have:

- [ ] A categorized URL inventory with file paths and line numbers
- [ ] Flagged items reviewed and triaged
- [ ] A filtered list relevant to your specific goal (migration, security, or audit)
- [ ] Known URLs cross-referenced against the inventory for completeness

## Troubleshooting

**False positives from comments** -- code comments sometimes contain example URLs or pseudo-URLs that are not functional references. The agent attempts to distinguish real URL references from illustrative examples, but some false positives may appear. Review flagged items in context before acting on them.

**URLs in minified code** -- minified JavaScript and CSS files contain URLs that are harder to parse. The agent extracts URLs from minified files, but line numbers may be less useful when the entire file is a single line. Consider running extraction against source files rather than build output.

**Dynamically constructed URLs** -- URLs built at runtime by concatenating strings (e.g., `baseUrl + '/api/users'`) may not appear as complete URLs in the extraction. The agent captures the partial components it can find, but the full URL may only be visible at runtime. Check for base URL variables and template strings to catch these cases.
