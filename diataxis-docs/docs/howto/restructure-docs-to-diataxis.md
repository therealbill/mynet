---
title: "How to restructure docs to Diataxis"
description: "Transform existing documentation into the four-type Diataxis model using doc-inventory, the orchestrator, type-specific writers, and the cross-link validator."
weight: 1
doc_type: how-to
prerequisites:
  - "An existing documentation directory"
  - "The diataxis-docs plugin installed"
  - "Completed the [Getting Started tutorial]({{< ref \"tutorials/getting-started\" >}})"
est_time: "30 minutes"
roles: ["developer", "technical writer"]
stability: stable
---

# How to Restructure Docs to Diataxis

Transform an existing documentation directory into the four-type Diataxis model: tutorials, how-to guides, reference, and explanation.

**Goal:** Take a flat or loosely organized documentation directory and reorganize it into clean Diataxis sections with proper type separation and cross-linking.

## Prerequisites

- An existing documentation directory with markdown files
- The diataxis-docs plugin installed and available
- Familiarity with the Diataxis model (see {{< ref "explanation/diataxis-in-practice" >}})

## Steps

### 1. Scan and classify existing docs

Invoke the **doc-inventory** agent to scan your documentation directory and classify every page.

```
Scan and classify all docs at <your-docs-path> by Diataxis type
```

The agent generates `docs/_reports/inventory.json` containing each page's classified type, quality rating, and any gaps it detects.

Review the inventory output. Pay attention to pages classified as "mixed" -- these need splitting.

### 2. Create the target directory structure

Set up the four Diataxis sections:

```
mkdir -p <your-docs-path>/tutorials
mkdir -p <your-docs-path>/howto
mkdir -p <your-docs-path>/reference
mkdir -p <your-docs-path>/explanation
```

### 3. Plan the transformation

Invoke the **diataxis-orchestrator** agent with the inventory results:

```
Plan a Diataxis restructuring for <your-docs-path> based on the inventory
```

The orchestrator produces a plan listing:

- Pages that move to a section unchanged
- Mixed pages that need splitting into multiple pages
- Gaps requiring new content
- The agent invocation order

Confirm the plan before proceeding.

### 4. Split mixed-type pages

For each page classified as "mixed," extract content into separate pages by type:

- **Tutorial content** (step-by-step learning journeys) moves to `tutorials/`
- **How-to content** (task-focused recipes) moves to `howto/`
- **Reference content** (specifications, API details) moves to `reference/`
- **Explanation content** (conceptual discussion, design rationale) moves to `explanation/`

Use the appropriate type-specific writer agent for each split. For example, invoke **doc-reference-gen** to rewrite extracted API details as a proper reference page.

### 5. Move single-type pages to their correct section

Pages that match one Diataxis type cleanly can move directly:

```
mv <your-docs-path>/api-reference.md <your-docs-path>/reference/api-reference.md
```

Update any internal links that reference the old path.

### 6. Fill gaps with new content

For each gap identified in the inventory:

- Missing onboarding tutorial: invoke **doc-tutorial-writer**
- Missing how-to for a common task: invoke **doc-howto-writer**
- Missing reference for a public API: invoke **doc-reference-gen**
- Missing architectural explanation: invoke **doc-explanation-writer**

### 7. Add cross-links between types

Each Diataxis type should link to related content in other types:

- Tutorials link to reference pages and how-to guides in their "Next steps" section
- How-to guides link to prerequisite tutorials, reference for API details, and explanations for background
- Reference pages link to how-to guides that demonstrate usage
- Explanation pages link to how-to guides for practical application and reference for specifications

Use Hugo ref shortcodes for internal links: `{{</* ref "reference/api" */>}}`

### 8. Validate the restructured docs

Invoke the **doc-crosslink-validator** agent:

```
Validate the documentation at <your-docs-path> for Diataxis compliance
```

Fix any errors (missing frontmatter, type violations, broken links) before considering the restructuring complete.

## Verify It Works

After completing the restructuring:

1. The `docs/_reports/quality.json` report shows zero errors
2. Every page has complete frontmatter with `doc_type` set to one of the four types
3. Each section directory contains at least one content page
4. No page mixes content from multiple Diataxis types
5. Cross-links connect related content across sections

## Troubleshooting

**Some pages do not fit any Diataxis type cleanly.**
Content like changelogs, contributing guides, or FAQs sit outside the four types. Keep them at the top level of your docs directory. The Diataxis model covers product documentation, not all project files.

**The inventory classifies a page differently than expected.**
Review the classification signals. A page you consider a "tutorial" may actually be a how-to if it solves a specific problem rather than teaching through building. Read {{< ref "explanation/diataxis-in-practice" >}} for the distinctions.

**Splitting a page creates two pages that feel incomplete.**
This is expected. After splitting, each page needs to be expanded to stand on its own. A how-to extracted from a mixed page may need a proper verification section. An explanation extracted may need trade-off discussion added.

**Too many cross-link warnings in the quality report.**
Address cross-link warnings iteratively. Start with the highest-traffic pages and add links between related content in different sections. Not every page needs links to all other types.

## See Also

- {{< ref "tutorials/getting-started" >}} -- Full tutorial walking through the entire pipeline
- {{< ref "howto/validate-doc-quality" >}} -- Detailed guide on interpreting and fixing quality reports
- {{< ref "reference/agents" >}} -- Specifications for all seven agents
- {{< ref "explanation/diataxis-in-practice" >}} -- The principles behind Diataxis type separation
- {{< ref "explanation/orchestration-model" >}} -- How the orchestrator coordinates the restructuring workflow
