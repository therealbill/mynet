---
title: "Understanding Diataxis in practice"
description: "How the four Diataxis documentation types work together, why type purity matters, and common pitfalls when applying the framework."
weight: 2
doc_type: explanation
prerequisites:
  - "Basic familiarity with documentation writing"
est_time: "15 minutes"
roles: ["developer", "technical writer", "architect"]
stability: stable
---

# Understanding Diataxis in Practice

The Diataxis framework organizes documentation into four types based on what the reader needs at that moment. This page explains how the four types work together as a system, why mixing types degrades documentation quality, and common pitfalls teams encounter when applying the framework.

## The Problem

Most documentation fails not because it is poorly written, but because it tries to do too many things at once. A page titled "Authentication Guide" might start with a conceptual overview, transition into step-by-step configuration, include a full API reference table, and end with a tutorial-style walkthrough. Each section may be individually well-written, but the reader looking for one specific thing must wade through everything else to find it.

The result is documentation that is simultaneously too much and not enough: too much for the user who needs a quick recipe, not enough for the user who needs a complete specification.

## Overview: The Four Types

Diataxis identifies four fundamentally different documentation needs, each requiring a different writing approach:

```
                  Studying              Working
              (acquiring knowledge)  (applying knowledge)
             +---------------------+---------------------+
  Practical  |     Tutorials       |   How-to Guides     |
  (steps)    | learning-oriented   |   task-oriented      |
             +---------------------+---------------------+
 Theoretical |   Explanation       |    Reference         |
  (knowledge)| understanding-      |   information-       |
             |   oriented          |     oriented         |
             +---------------------+---------------------+
```

Each type serves readers in a specific mode:

- **Tutorials** serve readers who are studying and need practical steps. They are learning-oriented. The reader does not yet know the tool and needs to build confidence by achieving a working result.

- **How-to guides** serve readers who are working and need practical steps. They are task-oriented. The reader knows the tool and needs a recipe to accomplish a specific goal.

- **Reference** serves readers who are working and need theoretical knowledge. It is information-oriented. The reader needs to look up a specific fact: a parameter type, an error code, a configuration option.

- **Explanation** serves readers who are studying and need theoretical knowledge. It is understanding-oriented. The reader wants to understand why the system works the way it does, what trade-offs were made, or how components relate.

## How the Types Work Together

The four types form a connected system. Each type addresses a gap the others deliberately leave:

**Tutorials create users.** A new user's first interaction with documentation is typically a tutorial. The tutorial's job is to get them to a working result as fast as possible, building confidence. It deliberately omits explanations of why things work and alternatives to the path shown. Those gaps are filled by the other three types.

**How-to guides make users productive.** Once a user understands the basics (from the tutorial), they face specific tasks: configure authentication, deploy to production, integrate with another service. Each how-to guide solves one task. It does not teach basics (the tutorial did that) or explain design decisions (explanation does that). It provides a recipe.

**Reference makes users self-sufficient.** As users gain experience, they need to look up specific details without reading narrative content. What parameters does this function accept? What error codes can this endpoint return? What is the default value of this configuration option? Reference provides facts without opinions.

**Explanation deepens understanding.** Experienced users eventually ask "why": Why does the system use event sourcing? Why is this API synchronous while that one is asynchronous? What trade-offs drove the architecture? Explanation builds the mental model that makes the facts in reference and the steps in how-to guides make sense.

### Cross-referencing Between Types

The types connect through deliberate cross-links:

- A tutorial's "Next steps" section links to how-to guides for common follow-up tasks, reference for the APIs used, and explanations for concepts mentioned
- A how-to guide's prerequisite section links to tutorials for setup knowledge, its steps reference API details in reference pages, and it links to explanations for conceptual background
- A reference page's "See Also" links to how-to guides that demonstrate usage patterns
- An explanation page links to how-to guides for practical application of the concepts discussed and reference pages for the technical details mentioned

These cross-links are not decorative. They are the mechanism by which a reader navigates from one mode (studying, working) and one orientation (practical, theoretical) to another.

## Why Type Purity Matters

Type purity means each documentation page serves exactly one of the four purposes. This is the most counterintuitive aspect of Diataxis and the one most frequently violated.

### The Case Against Mixing

Consider a how-to guide for configuring OAuth that includes three paragraphs explaining why OAuth works the way it does. The author's intention is helpful: give the reader context. The effect is harmful:

- The reader who needs the steps must read through explanation to find them
- The reader who wants the explanation gets an incomplete version embedded in a recipe
- The page is longer than necessary for either purpose
- Maintenance becomes harder because changes to the OAuth architecture require updating a how-to guide

The same content, split into a how-to guide (steps only, linking to the explanation) and an explanation page (design rationale only, linking to the how-to), serves both readers better while being easier to maintain.

### Type Mixing Patterns

The most common type violations, from most to least frequent:

**Explanation in how-to guides.** The author adds "why" paragraphs between steps. The fix: extract explanation content into its own page and add a link.

**Advice in reference.** The reference author writes "You should always validate inputs" instead of "Input validation is the caller's responsibility." The fix: rewrite as factual statements and move prescriptive content to a how-to guide.

**Reference in tutorials.** The tutorial includes a complete parameter table for an API instead of showing only the parameters needed for the current step. The fix: show only what the learner needs and link to the reference page for the full specification.

**Steps in explanation.** The explanation includes "First, configure X, then run Y" to illustrate a concept. The fix: describe the concept without imperative instructions and link to the how-to guide for the actual steps.

**Teaching in how-to guides.** The how-to guide starts with "Before we begin, let me explain..." and spends several paragraphs teaching basics. The fix: add a prerequisite link to the tutorial and assume the knowledge.

## Common Pitfalls

### Pitfall: Treating Diataxis as a Filing System

Diataxis is not just a directory structure. Creating `tutorials/`, `howto/`, `reference/`, and `explanation/` folders and moving files into them does not make documentation Diataxis-compliant. The content of each page must match its type. A page in the `reference/` directory that contains advice is a misclassified how-to, not a reference page.

### Pitfall: One Page Per Feature

Teams sometimes organize documentation by feature: "Authentication" gets one page covering the tutorial, how-to, reference, and explanation for authentication. This produces the mixed-type pages Diataxis is designed to prevent. Instead, authentication may appear in several places: a getting-started tutorial that includes authentication setup, a how-to guide for configuring OAuth, a reference page for the auth API, and an explanation of the security model.

### Pitfall: Tutorials That Are How-to Guides

A "tutorial" titled "How to set up your first project" is likely a how-to guide, not a tutorial. Tutorials guide learners through building something, emphasizing learning. How-to guides help experienced users accomplish a goal, emphasizing efficiency. The distinction is the reader's starting point: if they need to learn, it is a tutorial. If they need to accomplish, it is a how-to.

### Pitfall: Explanation Without Practical Links

Explanation pages that discuss concepts without linking to practical content (how-to guides, reference) become academic. The reader finishes understanding why something works and has no path to applying that understanding. Every explanation page should connect back to the practical types.

### Pitfall: Reference That Teaches

Reference pages that include "Getting Started" sections or "Quick Start" examples are mixing reference with tutorial content. A reference page for a function documents its parameters, return values, and error codes. A minimal syntax example is appropriate. A workflow example is not.

## Trade-offs of the Diataxis Model

| Benefit | Cost |
|---------|------|
| Readers find what they need faster | More pages to maintain |
| Each page has a clear, testable purpose | Content is sometimes duplicated across types |
| Type constraints improve writing quality | Authors must learn to recognize and separate types |
| Cross-links create navigation paths | Cross-links require maintenance |
| Structure scales with project complexity | Initial restructuring takes significant effort |

## Common Misconceptions

**"Every feature needs all four types."**
Some features may only need reference documentation. Others may need a how-to guide and an explanation but no tutorial. The four types are a framework, not a checklist. Create documentation based on what readers need, not on filling every quadrant.

**"Explanation is optional."**
Explanation is the most frequently omitted type because it does not directly help users accomplish tasks. However, explanation is what makes experienced users effective. Understanding why a system works the way it does prevents misuse that how-to steps alone cannot.

**"Short docs do not need Diataxis."**
Even a small project benefits from type separation. A README that mixes "quick start" (tutorial), "usage" (how-to), and "API" (reference) sections serves each purpose poorly. Separating them, even into sections of the same page, improves clarity.

## Related

- {{< ref "reference/agents" >}} -- Specifications for the four type-specific writer agents
- {{< ref "explanation/architecture" >}} -- Why diataxis-docs uses one agent per type
- {{< ref "explanation/orchestration-model" >}} -- How the orchestrator sequences type-specific writing
- {{< ref "howto/restructure-docs-to-diataxis" >}} -- Practical guide to applying Diataxis to existing docs
- {{< ref "tutorials/getting-started" >}} -- Hands-on tutorial demonstrating the full Diataxis pipeline
