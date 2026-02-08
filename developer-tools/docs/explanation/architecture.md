---
title: "Architecture"
description: "The developer workflow spectrum — version control, documentation, and prototyping"
weight: 1
---

# Architecture

The developer-tools plugin addresses three distinct concerns in the software development lifecycle: version control, documentation, and prototyping. This document explains why these concerns are separated into three agents, how they relate to each other, and the design philosophies behind each one.

## The Developer Workflow Spectrum

Software development involves a spectrum of activities that range from maintaining existing systems to creating entirely new ones. developer-tools organizes this spectrum into three zones:

**Version control** (git-workflow-manager) handles the coordination layer. How does code move from a developer's machine to production? What conventions prevent merge conflicts, preserve history, and automate releases? This is the connective tissue of any team workflow -- it does not produce user-visible output, but when it breaks, everything slows down.

**Documentation** (documentation-engineer) handles the knowledge layer. What does the code do? How do you use it? What decisions were made and why? Documentation bridges the gap between code that exists and people who need to understand it. Without documentation, knowledge lives only in the heads of the people who wrote the code.

**Prototyping** (rapid-prototyper) handles the creation layer. Can this idea work? Is there demand for this? What does the critical user journey feel like? Prototyping is the fastest path from hypothesis to evidence, trading durability for speed.

These three zones are ordered by their relationship to time. Version control manages the past and present of a codebase. Documentation captures the present for future readers. Prototyping explores the future as quickly as possible.

## Why Three Agents, Not One

A single "developer tools" agent would need to context-switch between maintaining Git workflow configurations, generating structured documentation, and scaffolding throwaway applications. These activities require different mental models:

- **git-workflow-manager** thinks in terms of processes, policies, and automation. It reads your repository's branch structure, CI configuration, and team dynamics. Its output is configuration files, hooks, and workflow definitions. It operates conservatively -- implementing the minimal change that solves the stated problem.

- **documentation-engineer** thinks in terms of audiences, information architecture, and accuracy. It reads your source code to understand what the software does, then structures that understanding for human consumption. Its output is prose, organized according to the Diataxis framework. It operates precisely -- producing no documentation rather than inaccurate documentation.

- **rapid-prototyper** thinks in terms of hypotheses, leverage, and speed. It reads a product idea and selects the fastest path to a deployed demo. Its output is running code with deliberate shortcuts. It operates aggressively -- choosing hosted services over self-managed infrastructure, skipping tests for non-critical paths, and marking technical debt with `TODO` rather than resolving it.

An agent optimized for conservative minimal changes cannot simultaneously be optimized for aggressive speed. An agent focused on information architecture operates differently from one focused on process automation. Separating these concerns allows each agent to be opinionated about its domain without compromising the others.

## The 6-Day Cycle Philosophy

rapid-prototyper follows a fixed 6-day cycle that treats disposability as a feature, not a compromise.

Days 1 and 2 focus exclusively on core features -- the minimum needed to test the central hypothesis. If the hypothesis is "people will pay for AI-generated meal plans," the core feature is generating a meal plan and showing a payment screen. Nothing else matters until this works end to end.

Days 3 and 4 add secondary features that make the demo coherent -- user accounts, settings, a landing page. These features exist to support the core journey, not to be valuable on their own.

Day 5 is testing, but only of the critical path. The prototype does not need comprehensive test coverage. It needs the one journey that an investor, user, or stakeholder will walk through to evaluate the idea.

Day 6 is polish and deployment. Loading states, error messages, visual consistency, and a publicly accessible URL. A prototype that cannot be demonstrated does not exist.

The cycle is deliberately short because prototypes that take longer than a week tend to accumulate emotional investment. Developers start treating the prototype as a real codebase, adding features, refactoring for cleanliness, optimizing performance -- all activities that are correct for production software and counterproductive for hypothesis testing. The 6-day constraint prevents this drift by making the timeline feel temporary from the start.

Shortcuts are not hidden. Every expedient decision is marked with a `TODO` comment explaining what the production version would do differently. This makes the prototype honest about its limitations while keeping velocity high.

## Documentation as Code

documentation-engineer treats documentation as a code artifact. This is more than a workflow preference -- it is a design principle that shapes how the agent generates, validates, and maintains docs.

Documentation lives in the repository alongside the code it describes. It goes through pull request review. It is validated by CI -- link checkers verify references, code sample validators ensure examples compile, and schema drift detection catches documentation that no longer matches the API it describes.

This approach solves the most common documentation failure mode: drift. Documentation maintained in a wiki, a separate repository, or a hosted platform has no mechanical connection to the code. When the code changes, nothing forces the documentation to update. Over time, the docs become wrong, developers stop trusting them, and new developers stop reading them.

By storing docs in the same repository and validating them in the same CI pipeline, documentation-engineer makes drift visible and fixable through normal development workflows. A PR that changes an API endpoint and does not update the corresponding documentation fails CI, just like a PR that breaks a test.

The agent selects documentation tooling (Markdown, MkDocs, Docusaurus, Sphinx) based on the project's technology stack and complexity needs, but the docs-as-code principle applies regardless of the tooling choice.

## Cross-Plugin Relationships

developer-tools connects to other plugins in the marketplace:

- **[diataxis-docs](/diataxis-docs/)** provides agents and skills for Diataxis-structured documentation. documentation-engineer handles the reference quadrant (API docs from source code); diataxis-docs handles the full four-quadrant structure when a project needs tutorials, how-to guides, and explanations alongside reference material.

- **[code-quality](/code-quality/)** provides code review and test automation. git-workflow-manager sets up the branching and PR workflow; code-quality's code-reviewer agent operates within that workflow to review changes before they merge.

These plugins are complementary, not overlapping. developer-tools handles the workflow infrastructure. Other plugins handle the activities that occur within that infrastructure.
