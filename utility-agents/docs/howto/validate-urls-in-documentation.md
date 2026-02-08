---
title: "Validate URLs in Documentation"
description: "Use url-context-validator to check links for freshness and relevance"
weight: 1
---

# Validate URLs in Documentation

Verify that all links in your documentation are working, contextually appropriate, and pointing to current content.

## Problem

Documentation accumulates links over time. Pages get moved, APIs get versioned, tutorials get rewritten. A link that worked six months ago may now 404, redirect to a different page, or point to an outdated version of a resource. Simple HTTP checks catch broken links but miss the more insidious problem: links that still work but no longer point to the right content.

## Solution

Use the url-context-validator agent to check every link for both technical reachability and contextual correctness.

## Prerequisites

- The utility-agents plugin installed
- One or more documentation files containing URLs (markdown, HTML, or any text format with links)

## Steps

### 1. Identify the Documents to Validate

Decide the scope of your check. You can validate a single file, a directory of documentation, or an entire docs site. For a focused check on a specific file:

```
Check the links in docs/api-guide.md.
```

For a broader sweep across a documentation directory:

```
Validate all the references in the docs/ directory.
```

The url-context-validator agent activates and begins extracting URLs from the specified content using Read.

### 2. Review the Validation Report

The agent fetches each URL with WebFetch and compares the destination against the surrounding context. Results are grouped by severity:

- **Broken** -- the URL returns an error. The page is gone, the domain expired, or the server is unreachable. Action: replace with a working URL or remove the reference entirely.
- **Misaligned** -- the URL works but the destination content does not match what the link text or surrounding paragraph describes. A link labeled "authentication guide" that lands on a billing page is misaligned. Action: update the link text to match the destination, or find a URL that matches the original intent.
- **Stale** -- the URL works and is contextually reasonable, but the content is outdated. The agent uses WebSearch to check whether a more current version exists. Action: update the URL to point to the current version.
- **Fine** -- the URL works, matches its context, and points to current content. No action required.

Each entry includes the source file, line number, anchor text, destination content summary, and recommended action.

### 3. Apply Recommended Actions

Work through the findings by severity:

**For broken links**, either find a replacement URL or remove the reference. If the linked content has moved, search for the new location. The agent may suggest a replacement URL if it found one via WebSearch.

**For misaligned links**, decide whether the anchor text or the URL is wrong. If the destination content is valuable but the label is misleading, update the label. If the label is correct but the destination is wrong, find the right URL.

**For stale links**, update the URL to the version the agent recommends. Verify that the new version covers the same topic -- a link to "React 16 lifecycle methods" should not be replaced with a generic "React documentation" link.

### 4. Re-run Validation After Fixes

After applying changes, run the validator again on the same scope:

```
Re-check the links in docs/api-guide.md.
```

Confirm that previously broken, misaligned, or stale links now report as fine. New issues may appear if a replacement URL itself has problems -- iterate until the report is clean.

## Verification

After completing these steps you should have:

- [ ] A validation report with all links categorized by severity
- [ ] All broken links replaced or removed
- [ ] All misaligned links corrected (text or URL updated)
- [ ] All stale links updated to current versions
- [ ] A clean re-validation confirming no remaining issues

## Troubleshooting

**Links behind authentication** -- URLs that require login (internal wikis, private repositories, paywalled content) will appear as broken because WebFetch cannot authenticate. Note these as expected failures and verify them manually.

**Rate limiting during bulk checks** -- if you are validating hundreds of URLs, destination servers may rate-limit the requests. The agent handles this gracefully by reporting rate-limited URLs separately. Re-run validation for those URLs after a delay, or validate in smaller batches.

**Redirected URLs** -- a URL that redirects is not necessarily broken. The agent reports the final destination so you can decide whether to update the link to the canonical URL or keep the redirect. Prefer the canonical URL to avoid depending on redirect chains.
