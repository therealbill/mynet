---
title: "How to validate documentation quality"
description: "Check documentation quality using doc-crosslink-validator, interpret the quality report, and fix issues by category."
weight: 4
doc_type: how-to
prerequisites:
  - "The diataxis-docs plugin installed"
  - "An existing documentation directory structured with Diataxis sections"
est_time: "10 minutes"
roles: ["developer", "technical writer"]
stability: stable
---

# How to Validate Documentation Quality

Check documentation quality using the **doc-crosslink-validator** agent, interpret the quality report, and fix issues by category.

**Goal:** Run a quality validation on your documentation, understand the report output, and resolve errors, warnings, and recommendations.

## Prerequisites

- The diataxis-docs plugin installed and available
- A documentation directory with content organized in Diataxis sections (tutorials, howto, reference, explanation)

## Steps

### 1. Invoke the cross-link validator

Run the **doc-crosslink-validator** agent against your documentation directory:

```
Validate the documentation at <your-docs-path> for Diataxis compliance
```

The agent scans all documentation files and generates two outputs:

- A machine-readable report at `docs/_reports/quality.json`
- A human-readable markdown summary

### 2. Review the summary counts

The report begins with a summary showing:

- Total pages scanned
- Pages by Diataxis type
- Error count (must fix)
- Warning count (should fix)
- Info count (nice to have)

Focus on the error count first. Errors indicate violations that break Diataxis compliance.

### 3. Fix frontmatter errors

Frontmatter errors are the most common issue. Every page requires:

- `title` -- page title following the naming convention for its type
- `description` -- one-sentence summary
- `doc_type` -- one of: `tutorial`, `how-to`, `reference`, `explanation`
- `stability` -- one of: `stable`, `beta`, `experimental`, `deprecated`

For tutorials and how-to guides, also include:

- `prerequisites` -- list of required prior knowledge or setup
- `est_time` -- estimated completion time

For reference pages, also include:

- `version` -- the version of the API or tool documented

### 4. Fix type purity violations

Type purity violations mean a page mixes content from multiple Diataxis types. Common violations:

- **Advice in reference**: statements like "you should" or "we recommend" in a reference page. Replace with factual language or move the advice to a how-to guide.
- **Explanation in how-to**: paragraphs explaining why something works a certain way inside a task-focused guide. Extract to an explanation page and link to it.
- **Steps in explanation**: numbered instructions inside a conceptual page. Move the steps to a how-to guide and link to it.
- **API specs in tutorial**: detailed parameter tables inside a learning-oriented page. Link to the reference page instead.

### 5. Fix broken links

Broken link errors include:

- Internal links pointing to pages that do not exist
- Relative paths that resolve incorrectly after restructuring
- Hugo ref shortcodes with incorrect paths

Fix each broken link by updating the path to the correct location. Use Hugo ref shortcodes for internal links to catch broken references at build time: `{{</* ref "howto/some-guide" */>}}`

### 6. Address orphaned pages

Orphaned pages are documentation files that no other page links to. Users cannot discover them through navigation.

For each orphaned page:

- Add a link from a related page in another Diataxis section
- Add it to the section's `_index.md` if it has a page listing
- Consider whether the page is still needed -- outdated pages should be removed or marked deprecated

### 7. Improve cross-linking

The validator checks that pages link to related content across Diataxis types. Minimum cross-link targets:

- **Tutorials**: 3 or more links to reference, how-to guides, and explanations
- **How-to guides**: 5 or more links to prerequisites, reference, related how-tos, and explanations
- **Reference pages**: 2 or more links to how-to guides showing usage
- **Explanation pages**: 4 or more links to how-to guides, reference, and related concepts

Add links in "Next steps," "See Also," or inline where the reader would benefit from additional context.

### 8. Re-run validation

After fixing issues, run the validator again to confirm all errors are resolved:

```
Validate the documentation at <your-docs-path> for Diataxis compliance
```

Repeat until the error count reaches zero.

## Verify It Works

A passing validation shows:

- Zero errors in the quality report
- Every page has complete, valid frontmatter
- No type purity violations detected
- All internal links resolve correctly
- No orphaned pages exist
- Cross-link counts meet the minimum thresholds

## Troubleshooting

**The validator reports type violations in pages that seem correct.**
Review the specific violation. A how-to guide that says "the reason this works is..." contains explanation content. Even a single paragraph of the wrong type triggers a violation. Extract the content to the correct type and link to it.

**Cross-link counts seem too strict.**
The thresholds are minimums for well-connected documentation. If a page genuinely has nothing to link to in another type, that may indicate a gap in documentation coverage rather than a linking problem.

**The validator cannot read certain files.**
Check file permissions and encoding. The agent reads `.md`, `.mdx`, and `.rst` files. Files with non-UTF-8 encoding or binary content will fail to scan.

**Validation passes but navigation still feels disconnected.**
The validator checks for the presence of cross-links but not their placement or quality. Manually review whether links appear at the point where a reader would need them, not just in a "See Also" section at the bottom.

## See Also

- [doc-crosslink-validator agent specification](../../reference/agents/)
- [Understanding type purity and why it matters](../../explanation/diataxis-in-practice/)
- [Full restructuring workflow including validation](../../howto/restructure-docs-to-diataxis/)
- [Tutorial covering the end-to-end pipeline including validation](../../tutorials/getting-started/)
