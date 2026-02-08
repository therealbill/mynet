---
title: "Agents"
description: "Technical specifications for all developer-tools agents"
weight: 1
---

# Agents

Technical specifications for the three agents in the developer-tools plugin.

## documentation-engineer

- **Model:** opus
- **Color:** cyan
- **Tools:** Read, Write, Edit, Bash, Grep, Glob

### Trigger Patterns

- "REST API has 40 endpoints and zero documentation"
- "We need a docs site for our SDK"
- "Nobody reads our docs"

### Capabilities

- Documentation-as-code workflows where docs live in the repository and go through PR review
- Diataxis framework application for structuring content across tutorials, how-to guides, reference, and explanation
- Automated documentation generation from source code analysis
- Tooling selection and configuration: MkDocs, Docusaurus, Sphinx, or plain Markdown
- CI validation setup including link checkers, code sample validators, and schema drift detection
- Source code audit for routes, schemas, authentication patterns, and error handling

### Process

1. Audit current documentation state and codebase structure
2. Understand the target audience and their information needs
3. Structure content according to the Diataxis framework
4. Implement documentation with automation for generation and validation
5. Verify accuracy against actual code behavior

### Constraints

- Does not duplicate information already present in well-written code comments
- Does not choose complex infrastructure (MkDocs, Docusaurus) when a `docs/` directory with Markdown suffices
- Does not create empty templates or placeholder content
- Prioritizes usefulness over coverage -- fewer accurate docs over many shallow ones

---

## git-workflow-manager

- **Model:** opus
- **Color:** blue
- **Tools:** Read, Write, Edit, Bash, Grep, Glob

### Trigger Patterns

- "Set up our Git workflow"
- "Main branch history is a mess"
- "Automate our release process"

### Capabilities

- Branching strategy selection: trunk-based development, GitHub Flow, Git Flow
- Commit convention enforcement using Conventional Commits (`feat:`, `fix:`, `docs:`, `chore:`, `refactor:`, `test:`)
- Merge policy configuration: squash-merge for features, real merge commits for releases
- Automation setup: commit hooks, PR templates, labels, auto-merge rules
- Release management: semantic versioning, automated changelog generation, tag-driven deployments
- Conflict prevention through branch naming conventions and short-lived branches
- Agent-parallel workflow support using Git worktrees for concurrent Claude Code sessions

### Process

1. Assess current Git workflow, branch structure, and CI configuration
2. Identify the specific problem or gap to address
3. Implement the minimal change that solves the stated problem
4. Configure automation to enforce the new workflow
5. Verify end-to-end by testing the branch-commit-PR-merge cycle

### Constraints

- Does not overhaul workflows that are already working
- Does not force-push shared branches
- Does not configure Git Flow when trunk-based development suffices
- Implements minimal changes rather than complete workflow replacements

---

## rapid-prototyper

- **Model:** opus
- **Color:** green
- **Tools:** Write, Edit, Bash, Read, Glob

### Trigger Patterns

- "Create a phone anxiety app -- test the concept this week"
- "Meeting with investors -- need functional demo"
- "AI avatars are blowing up -- build something in a few days"

### Capabilities

- 6-day development cycle execution:

  - Days 1-2: Core features -- the minimum needed to test the central hypothesis
  - Days 3-4: Secondary features -- supporting functionality that makes the demo coherent
  - Day 5: Testing -- verify the critical path works end to end
  - Day 6: Polish and deploy -- visual cleanup, error handling, production deployment

- Leverage maximization with hosted services: Supabase for database and auth, Clerk for identity, Stripe for payments
- Critical path prioritization -- the single user journey that proves or disproves the hypothesis ships first
- Deliberate shortcuts marked with `TODO` comments for future resolution
- Demo-readiness as a non-negotiable requirement -- the prototype must be presentable

### Process

1. Clarify the core question the prototype is meant to answer
2. Choose the simplest technology stack that supports the core question
3. Scaffold the project structure within 30 minutes
4. Build the critical user journey end to end
5. Add demo polish -- loading states, error messages, visual consistency
6. Deploy to a publicly accessible URL and collect feedback

### Constraints

- Does not build infrastructure intended to outlive the prototype
- Does not add features beyond the core hypothesis being tested
- Does not optimize performance before validating the concept
- Does not treat the prototype as a production codebase

---

## Comparison

| Agent | Model | Scope | Primary Output |
|-------|-------|-------|----------------|
| documentation-engineer | opus | Documentation systems | Reference docs, doc site configuration, CI validation |
| git-workflow-manager | opus | Version control workflows | Branching strategy, commit hooks, PR templates, release automation |
| rapid-prototyper | opus | Application prototyping | Deployed MVP with critical path implemented in 6 days |
