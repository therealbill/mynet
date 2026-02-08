---
title: "Restructure your docs with the Diataxis pipeline"
description: "A hands-on tutorial that takes you from a messy documentation directory to a clean, four-type Diataxis structure using all seven diataxis-docs agents."
weight: 1
doc_type: tutorial
prerequisites:
  - "An existing documentation directory with at least 5 markdown files"
  - "The diataxis-docs plugin installed and available in your Claude Code session"
  - "Familiarity with markdown"
est_time: "45 minutes"
roles: ["developer", "technical writer"]
stability: stable
---

# Restructure Your Docs with the Diataxis Pipeline

This tutorial walks you through transforming an existing documentation directory into a clean Diataxis structure. You will use all seven agents in the diataxis-docs plugin: scan and classify your existing pages, plan the transformation, write content with four type-specific writers, and validate the result.

## What You Will Build

By the end of this tutorial, you will have:

- An inventory report classifying every existing doc page by Diataxis type
- A restructured `docs/` directory with `tutorials/`, `howto/`, `reference/`, and `explanation/` sections
- Content written or reorganized by the four type-specific writer agents
- A quality report confirming Diataxis compliance

## Prerequisites

- An existing documentation directory with at least 5 markdown files
- The diataxis-docs plugin installed and available in your Claude Code session
- Familiarity with markdown

## Step 1: Identify Your Documentation Directory

Locate the documentation directory you want to restructure. For this tutorial, assume your docs live at `my-project/docs/`.

Verify the directory exists and contains markdown files:

```
ls my-project/docs/*.md
```

You should see a list of your existing documentation files.

## Step 2: Run the Inventory Scan

Use the **doc-inventory** agent to scan and classify every page in your documentation directory.

```
Scan and classify the docs at my-project/docs/ by Diataxis type
```

The doc-inventory agent will:

1. Read every `.md` file in the directory
2. Classify each page as Tutorial, How-to, Reference, Explanation, or Mixed
3. Rate the quality match (strong, moderate, weak, mixed)
4. Identify gaps in your documentation coverage
5. Generate an inventory report at `docs/_reports/inventory.json`

### What Just Happened?

The inventory agent applied Diataxis classification signals to each page. Pages that walk through building something are classified as tutorials. Pages with numbered steps solving a specific problem are how-to guides. Pages listing API specs without advice are reference. Pages explaining why things work a certain way are explanation. Pages that mix multiple types are flagged for splitting.

### Checkpoint 1: Review the Inventory

Open the generated `docs/_reports/inventory.json` and verify:

- Every documentation file appears in the `pages` array
- Each page has a `type` classification
- The `gaps` array lists missing documentation types
- The `reorganization` array suggests specific moves

If any files are missing from the inventory, re-run the scan. The agent must scan everything without skipping files.

## Step 3: Plan the Transformation with the Orchestrator

Use the **diataxis-orchestrator** agent to create a transformation plan based on the inventory.

```
Plan a Diataxis transformation for my-project/docs/ based on the inventory at docs/_reports/inventory.json
```

The orchestrator will review the inventory and produce a plan that includes:

- Which pages move to which section unchanged
- Which mixed-type pages need splitting
- Which gaps need new content
- The order of operations for the transformation

### What Just Happened?

The orchestrator acts as a coordinator. It does not write documentation itself. Instead, it reads the inventory, identifies the work needed, and determines which specialist agents to invoke and in what sequence. Reference generation typically comes first because other doc types link to it.

### Checkpoint 2: Confirm the Plan

Review the orchestrator's proposed plan. Verify that:

- Mixed pages are identified for splitting with clear instructions on what content goes where
- Critical gaps (like a missing onboarding tutorial) are prioritized
- The plan does not propose deleting content, only reorganizing it

Confirm the plan before proceeding. The orchestrator will ask for your confirmation before making structural changes.

## Step 4: Create the Directory Structure

Set up the Diataxis directory structure that content will move into:

```
mkdir -p my-project/docs/tutorials
mkdir -p my-project/docs/howto
mkdir -p my-project/docs/reference
mkdir -p my-project/docs/explanation
```

Verify the directories exist:

```
ls -d my-project/docs/tutorials my-project/docs/howto my-project/docs/reference my-project/docs/explanation
```

You should see all four directories listed.

## Step 5: Generate Reference Documentation

Use the **doc-reference-gen** agent to create or reorganize your reference documentation.

```
Generate reference documentation for my-project from the source code at my-project/src/
```

The doc-reference-gen agent will:

- Identify public API surfaces, CLI commands, or configuration options
- Generate specification pages with consistent structure
- Enforce the "no advice" rule: no "you should" or "we recommend" statements
- Place output in `docs/reference/`

### What Just Happened?

Reference comes first in the pipeline because tutorials, how-tos, and explanations all link to reference pages. By generating reference first, the other writers have targets to cross-reference.

### Checkpoint 3: Verify Reference Content

Open one of the generated reference pages and check:

- It has frontmatter with title, description, stability, and version
- It uses consistent structure (signature, parameters, returns, errors)
- It contains no prescriptive language ("you should", "we recommend")
- It includes "See Also" links

## Step 6: Write the Onboarding Tutorial

Use the **doc-tutorial-writer** agent to create or restructure your getting-started content.

```
Write an onboarding tutorial for my-project that takes beginners from installation to a working example
```

The doc-tutorial-writer agent will:

- Design a 90-120 minute learning journey
- Write step-by-step instructions with complete code and expected outputs
- Add checkpoints every 3-5 steps
- Include a troubleshooting section
- Link to the reference pages generated in the previous step

### Checkpoint 4: Verify the Tutorial

Open the generated tutorial and check:

- The "What you'll build" section describes a concrete outcome
- Every step includes the complete command or code to run
- Expected output is shown after each command
- Checkpoints appear every 3-5 steps
- No choices or alternatives are presented (one path only)

## Step 7: Write How-to Guides

Use the **doc-howto-writer** agent to create how-to guides for your most common user tasks.

```
Write how-to guides for the top tasks users need to do with my-project
```

The doc-howto-writer agent will:

- Create one guide per task, each under 1800 words
- Use numbered steps with a single goal per guide
- Include verification and troubleshooting sections
- Link to reference for API details and explanations for conceptual background

### What Just Happened?

How-to guides differ from tutorials. Tutorials teach beginners by walking them through a full project. How-to guides help experienced users accomplish specific tasks quickly. Each guide assumes the reader already understands the basics and needs a recipe, not a lesson.

## Step 8: Write Explanation Pages

Use the **doc-explanation-writer** agent to create conceptual and architectural documentation.

```
Write explanation docs covering the architecture and design decisions of my-project
```

The doc-explanation-writer agent will:

- Explain why the system is designed the way it is
- Discuss trade-offs and alternatives that were considered
- Build mental models with diagrams and analogies
- Link to how-to guides for practical application and reference for details

### Checkpoint 5: Review All Content Types

At this point, all four Diataxis types should have content. Verify the directory structure:

```
ls my-project/docs/tutorials/
ls my-project/docs/howto/
ls my-project/docs/reference/
ls my-project/docs/explanation/
```

Each directory should contain at least one content file.

## Step 9: Validate with the Cross-link Validator

Use the **doc-crosslink-validator** agent to check the quality of your restructured documentation.

```
Validate the documentation at my-project/docs/ for Diataxis compliance
```

The doc-crosslink-validator agent will:

- Check that every page has complete frontmatter
- Verify each page fits one Diataxis type cleanly (no type mixing)
- Confirm cross-links exist between types (tutorials link to reference, how-tos link to explanations)
- Find broken links and orphaned pages
- Generate a quality report at `docs/_reports/quality.json`

### Checkpoint 6: Review the Quality Report

Open `docs/_reports/quality.json` and check:

- **Errors** (must fix): missing frontmatter, type violations, broken links
- **Warnings** (should fix): weak cross-linking, orphaned pages
- **Info** (nice to have): additional cross-link opportunities

Fix any errors reported before considering the transformation complete.

## What You Built

You have transformed a documentation directory into a clean Diataxis structure:

- **Inventory report** classifying all original content
- **Reference pages** with complete, advice-free API specifications
- **Onboarding tutorial** guiding beginners through a working project
- **How-to guides** solving specific user tasks with numbered steps
- **Explanation pages** covering architecture and design decisions
- **Quality report** confirming Diataxis compliance

## Next Steps

- {{< ref "howto/restructure-docs-to-diataxis" >}} -- Repeat this process on other documentation directories
- {{< ref "howto/validate-doc-quality" >}} -- Run validation checks regularly as docs evolve
- {{< ref "reference/agents" >}} -- Review the full specification of all seven agents
- {{< ref "explanation/architecture" >}} -- Understand why the plugin uses seven specialized agents
- {{< ref "explanation/diataxis-in-practice" >}} -- Learn the principles behind Diataxis type separation

## Troubleshooting

**The inventory scan missed some files.**
Verify the docs path is correct and the files have `.md` extensions. The doc-inventory agent scans for `.md`, `.mdx`, and `.rst` files. Files in hidden directories (starting with `.`) may be skipped.

**A page was classified as "mixed" but seems correct.**
Mixed classification means the page contains content belonging to multiple Diataxis types. This is common and expected. Review the agent's recommendation for how to split the content. A page that explains a concept and then gives step-by-step instructions should be split into an explanation page and a how-to guide.

**The reference generator included advice.**
Re-run the doc-reference-gen agent on the specific page. The "no advice" rule is strict: any "you should", "we recommend", or prescriptive language must be removed or moved to a how-to guide.

**The quality report shows many cross-link warnings.**
This is normal after initial restructuring. Cross-links are added iteratively. Focus on fixing errors first, then address warnings by adding links between related pages across different Diataxis types.

**The orchestrator did not invoke all agents.**
The orchestrator only invokes agents needed for the specific transformation. If your docs already have strong reference coverage, it may skip doc-reference-gen and focus on gaps. Review the plan to confirm the right agents are being used.
