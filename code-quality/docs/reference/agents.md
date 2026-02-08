---
title: "Agents"
description: "Technical specifications for all code-quality agents"
weight: 1
---

# Agents

Agent specifications for the code-quality plugin. Five agents covering code review, architectural review, test management, end-to-end test automation, and accessibility auditing.

## Overview

| Agent | Model | Color | Primary Focus |
|-------|-------|-------|---------------|
| code-reviewer | sonnet | yellow | Code correctness, security, maintainability |
| architect-reviewer | opus | yellow | Dependency direction, boundary integrity, pattern consistency |
| test-writer-fixer | opus | cyan | Test creation, failure diagnosis, test repair |
| playwright-expert | sonnet | cyan | Playwright E2E test automation |
| web-accessibility-checker | sonnet | green | WCAG 2.2 compliance auditing |

---

## code-reviewer

Reviews recent code changes for correctness, security, and maintainability. Analyzes git diffs with full file context.

### Specification

| Field | Value |
|-------|-------|
| Name | code-reviewer |
| Model | sonnet |
| Color | yellow |
| Tools | Read, Grep, Glob, Bash |

### Trigger Conditions

code-reviewer activates when the user mentions:

- Reviewing recent code changes
- Checking code quality before committing or pushing
- Asking whether code "looks right"

Example phrases: "Review the changes I just made", "Does this look right?", "Check my code before I push"

### Capabilities

- Runs `git diff` (staged and unstaged) to identify changes
- Reads full file context of modified files, not just diff hunks
- Detects logic errors, off-by-one errors, nil/null dereference, race conditions
- Identifies injection vulnerabilities, exposed secrets, missing input validation, unsafe deserialization
- Flags maintainability issues: confusing names, duplicated logic, overloaded functions
- Checks OpenTelemetry instrumentation: broken traces from `context.Background()`/`context.TODO()` misuse, unbounded spans in loops, missing `span.RecordError`/`span.SetStatus`

### Review Process

1. Run `git diff` (staged and unstaged) to see what changed
2. Read the full context of modified files
3. Review for priorities: correctness, security, maintainability, observability
4. Report findings grouped by severity

### Output Format

Findings grouped by severity:

- **Critical** -- must fix before merging (bugs, security vulnerabilities)
- **Warnings** -- should fix (maintenance problems)
- **Suggestions** -- consider (non-urgent improvements)

Each finding includes: file path, line number, description, fix. Categories with no findings are omitted. Clean code produces "No issues found."

### Constraints

- Does not nitpick style or formatting
- Does not manufacture feedback when code is clean
- Reports only on changed code and its immediate context

---

## architect-reviewer

Reviews code changes for architectural consistency, pattern adherence, and structural integrity.

### Specification

| Field | Value |
|-------|-------|
| Name | architect-reviewer |
| Model | opus |
| Color | yellow |
| Tools | Read, Grep, Glob, Bash |

### Trigger Conditions

architect-reviewer activates when the user mentions:

- Architectural review of new services or modules
- Checking refactored code for structural consistency
- Changes to inter-layer communication (API to domain, domain to data access)

Example phrases: "Check the architecture", "Review the refactoring I did on the data access layer", "I changed how the API layer talks to the domain layer"

### Capabilities

- Evaluates dependency direction (handlers -> services -> repositories, not reversed)
- Detects boundary violations: logic placed in the wrong layer
- Checks pattern consistency against existing codebase conventions
- Assesses coupling between components that should be independent
- Identifies proportional abstraction problems: over-engineering (premature generalization, single-implementation interfaces) and under-engineering (copy-pasted logic, missing boundaries)

### Review Process

1. Map the change within the overall architecture -- identify affected layers and boundaries
2. Check dependency direction and coupling
3. Compare against existing patterns in the codebase
4. Assess long-term impact on maintainability and extensibility

### Output Format

- **Impact:** High / Medium / Low
- **Violations:** Specific issues with file references and explanations
- **Recommendations:** Concrete refactoring suggestions

### Constraints

- Only flags real structural problems, not stylistic preferences
- Keeps output concise
- Compares against the existing codebase rather than imposing external architectural opinions

---

## test-writer-fixer

Writes new tests, runs existing tests, diagnoses failures, and fixes broken tests while preserving test intent.

### Specification

| Field | Value |
|-------|-------|
| Name | test-writer-fixer |
| Model | opus |
| Color | cyan |
| Tools | Read, Write, Edit, Bash, Grep, Glob |

### Trigger Conditions

test-writer-fixer activates when the user mentions:

- Code changes that need corresponding tests (proactive trigger)
- Modules or functions lacking test coverage
- Failing or broken tests after a refactor

Example phrases: "I've updated the authentication logic", "Our payment module has no tests", "Tests are broken after my refactor"

### Capabilities

- Writes tests that assert on behavior, not implementation details
- Uses descriptive test names that read as specifications
- Follows project's existing test patterns, framework, and conventions
- Covers happy path, error conditions, and edge cases in that priority order
- Runs affected tests first, then the full suite
- Classifies failures into four categories: real bug, stale expectation, brittle test, environment issue
- Fixes broken tests while preserving original test intent
- Runs fixed tests multiple times to detect flakiness

### Failure Classification

| Category | Meaning | Action |
|----------|---------|--------|
| Real bug | Test caught a genuine code problem | Report the bug; do not modify the test |
| Stale expectation | Code behavior legitimately changed | Update assertion to match new behavior |
| Brittle test | Test depends on implementation details, timing, or ordering | Refactor for resilience while preserving intent |
| Environment issue | Missing dependencies, config, or state | Fix setup, not the test |

### Review Process

1. Identify affected tests via file relationships and imports
2. Run affected tests
3. Classify any failures
4. Fix tests according to classification
5. Re-run to confirm stability

### Constraints

- Does not delete or skip failing tests without understanding the failure
- Does not add `sleep`/`waitForTimeout` to fix timing issues
- Does not write tests that pass but validate nothing meaningful
- Does not over-mock to the point of testing only mock wiring

---

## playwright-expert

Expert in Playwright test automation for web applications. Writes, debugs, and maintains end-to-end test suites.

### Specification

| Field | Value |
|-------|-------|
| Name | playwright-expert |
| Model | sonnet |
| Color | cyan |
| Tools | Read, Write, Edit, Bash, Grep, Glob |

### Trigger Conditions

playwright-expert activates when the user mentions:

- Writing Playwright or E2E tests for web flows
- Debugging flaky or timing-out Playwright tests
- Setting up Playwright in a project or CI pipeline

Example phrases: "Write Playwright tests for checkout flow", "Playwright tests keep timing out", "Set up Playwright with our CI"

### Capabilities

- Writes E2E tests using Playwright's auto-waiting assertion model
- Uses selector priority: `getByRole` > `getByLabel` > `getByText` > `getByTestId` (never CSS selectors tied to implementation)
- Implements Page Object Model for pages touched by multiple tests
- Configures network mocking with `page.route()` for deterministic tests
- Sets up CI configuration: `--workers=auto`, `--retries=2` (CI only), HTML reporter, browser installation via `npx playwright install --with-deps`
- Debugs flaky tests using tracing (`trace: 'on-first-retry'`), DOM snapshots, and network logs

### Selector Priority

| Priority | Selector | Rationale |
|----------|----------|-----------|
| 1 | `getByRole` | Most resilient to DOM changes |
| 2 | `getByLabel` | Stable for form elements |
| 3 | `getByText` | Stable for content-driven elements |
| 4 | `getByTestId` | Fallback when semantic selectors insufficient |
| Never | CSS selectors (class, tag nesting) | Tied to implementation details |

### Review Process

1. Understand the user flows to test
2. Write tests using page objects and auto-waiting assertions
3. Run locally to verify
4. Configure for CI
5. If fixing flaky tests: reproduce first, then diagnose with tracing

### CI Configuration Defaults

| Setting | Value |
|---------|-------|
| Workers | `--workers=auto` |
| Retries | `--retries=2` (CI only, not local) |
| Reporter | HTML reporter for artifact collection |
| Browser install | `npx playwright install --with-deps` |
| Tracing | `trace: 'on-first-retry'` |

### Constraints

- Does not use `page.waitForTimeout()` -- always uses auto-waiting assertions
- Does not use CSS selectors tied to class names or tag nesting
- Does not share state between tests -- each test gets its own browser context
- Does not enable retries for local development runs

---

## web-accessibility-checker

Audits web content against WCAG 2.2 and provides specific, actionable fixes with code examples.

### Specification

| Field | Value |
|-------|-------|
| Name | web-accessibility-checker |
| Model | sonnet |
| Color | green |
| Tools | Read, Grep, Glob, Bash |

### Trigger Conditions

web-accessibility-checker activates when the user mentions:

- Accessibility audits for web pages or components
- Form accessibility validation
- WCAG compliance requirements

Example phrases: "Audit this page for accessibility", "Make sure this form is accessible", "We need to be WCAG AA compliant"

### Capabilities

- Checks semantic HTML: appropriate elements (`nav`, `main`, `article`, `button`), heading order, ARIA used only when native semantics are insufficient
- Verifies keyboard navigation: all interactive elements focusable, logical tab order, visible focus indicators, modal focus trapping
- Validates screen reader compatibility: meaningful `alt` text, form label associations (`<label for>`, `aria-labelledby`), error message associations (`aria-describedby`), live regions (`aria-live`)
- Checks color and visual: contrast ratios (4.5:1 normal text, 3:1 large text for AA), information not conveyed by color alone, focus indicator contrast
- Validates forms: visible labels on every input, programmatic required field indication, error messages identifying field and describing error, related input grouping with `fieldset`/`legend`

### Audit Process

1. Check semantic HTML structure and ARIA usage
2. Verify keyboard navigation and focus management
3. Validate screen reader compatibility
4. Check color contrast and visual presentation
5. Validate form labeling and error handling

### Output Format

For each issue:

- WCAG success criterion (e.g., "1.1.1 Non-text Content")
- Level (A, AA, or AAA)
- File and line number
- Description and fix with code example

Level A violations listed first, then AA. Level AAA flagged only when specifically requested.

### Contrast Requirements (WCAG AA)

| Text Type | Minimum Ratio |
|-----------|---------------|
| Normal text | 4.5:1 |
| Large text (18pt+ or 14pt+ bold) | 3:1 |
| UI components and graphical objects | 3:1 |

### Constraints

- Provides specific, actionable fixes -- not generic advice
- Flags `role="button"` on `<div>` as needing a `<button>` element instead
- Prioritizes Level A and AA; does not flag AAA unless requested
- Performs static source code analysis; does not execute JavaScript at runtime

## See Also

- [Architecture]({{< ref "explanation/architecture" >}}) -- why five specialized agents instead of one
- [Choosing the Right Agent]({{< ref "explanation/choosing-the-right-agent" >}}) -- decision guide
- [Getting Started]({{< ref "tutorials/getting-started" >}}) -- tutorial walkthrough
