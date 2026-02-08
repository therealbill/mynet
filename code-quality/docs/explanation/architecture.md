---
title: "Architecture"
description: "Why five specialized reviewers instead of one general-purpose agent"
weight: 1
---

# Architecture

The code-quality plugin uses five specialized agents rather than one general-purpose code quality agent. This page explains why the plugin is structured this way, how the agents relate to each other, and the reasoning behind model selection.

## The Quality Spectrum

Code quality is not a single concern. It spans a spectrum of distinct disciplines:

- **Code correctness** -- Does the code do what it should? Are there logic errors, security vulnerabilities, or missing error handling?
- **Architectural integrity** -- Does the code respect established boundaries, dependency direction, and patterns?
- **Test coverage** -- Are the important behaviors tested? Do tests catch regressions without being brittle?
- **End-to-end validation** -- Do user flows work correctly in a real browser environment?
- **Accessibility compliance** -- Can all users interact with the application, including those using assistive technology?

Each of these disciplines requires different tools, different depth of analysis, and different output formats. A single agent attempting to cover the entire spectrum would either be shallow across all areas or overwhelmed by the breadth of its instructions.

## Why Specialization

The five agents exist because each review type has fundamentally different requirements:

**Different tools.** code-reviewer needs `git diff` and file reading to analyze changes in context. test-writer-fixer needs Write and Edit tools to create and modify test files. playwright-expert needs Write and Bash to generate test files and run Playwright. web-accessibility-checker needs Grep and Glob to scan markup patterns across many files. Giving every tool to a single agent increases the chance of the wrong tool being used for a given task.

**Different expertise depth.** An architectural review requires understanding dependency graphs, coupling metrics, and design pattern tradeoffs. An accessibility audit requires knowledge of WCAG 2.2 success criteria, ARIA specifications, and screen reader behavior. These are deep, specialized knowledge domains. Separating them allows each agent's system prompt to go deep rather than wide.

**Different output formats.** code-reviewer produces severity-grouped findings (Critical/Warnings/Suggestions). architect-reviewer rates impact as High/Medium/Low with specific violations and recommendations. web-accessibility-checker reports by WCAG success criterion with conformance level. test-writer-fixer produces actual code -- test files and fixes. A unified output format would lose the precision that makes each agent's output actionable.

**Different trigger patterns.** "Review my changes" and "audit this page for accessibility" are fundamentally different requests. Separate agents with distinct trigger patterns mean Claude Code can route each request to the right specialist without ambiguity.

## Model Selection Rationale

The five agents use two different models, chosen based on the cognitive demands of each task.

**Opus (deeper reasoning):** architect-reviewer and test-writer-fixer use the opus model.

- architect-reviewer must hold an entire codebase's architectural patterns in context, trace dependency chains across multiple layers, and reason about long-term structural impact. This requires the kind of multi-step reasoning that benefits from a more capable model.
- test-writer-fixer must understand existing code behavior, write tests that capture that behavior precisely, classify failures by root cause, and modify tests without weakening their intent. Distinguishing a real bug from a stale expectation requires reasoning about what the code should do versus what it does.

**Sonnet (pattern matching):** code-reviewer, playwright-expert, and web-accessibility-checker use the sonnet model.

- code-reviewer applies well-defined rules to code diffs -- checking for nil dereference, injection patterns, race conditions, and broken traces. These are pattern-matching tasks where speed matters more than deep reasoning chains.
- playwright-expert applies Playwright best practices -- selector priority, auto-waiting patterns, page object structure, CI configuration. The expertise is in knowing the right patterns, not in multi-step reasoning.
- web-accessibility-checker matches HTML patterns against WCAG criteria -- missing alt text, insufficient contrast ratios, absent form labels. Each check is a well-defined rule applied to markup. Speed and thoroughness matter more than reasoning depth.

This split balances cost and quality. Opus is used only where the task genuinely requires deeper reasoning. Sonnet handles the high-throughput pattern matching where its speed is an advantage.

## The Review Workflow

The five agents complement each other rather than competing. A typical quality workflow applies them in sequence:

1. **code-reviewer** runs first, catching correctness and security issues in the code diff. This is the fastest feedback loop -- immediate issues that would cause bugs.
2. **architect-reviewer** runs on structural changes, checking that new code respects existing patterns and boundaries. This catches problems that are correct line-by-line but wrong in aggregate.
3. **test-writer-fixer** writes or updates tests for the changed code, ensuring regressions will be caught in the future. It builds on the assumption that the code is already correct (after code-reviewer) and well-structured (after architect-reviewer).
4. **playwright-expert** validates user-facing flows end-to-end in a real browser, catching integration issues that unit tests miss. This is the most expensive check and runs less frequently.
5. **web-accessibility-checker** audits the UI for WCAG compliance, ensuring the application is usable by all users. This runs whenever UI components change.

Not every change needs all five agents. A backend refactor might need code-reviewer, architect-reviewer, and test-writer-fixer. A UI component change might need code-reviewer, playwright-expert, and web-accessibility-checker. The agents are independently triggered based on what changed and what needs checking.

## Cross-Plugin Collaboration

The code-quality agents focus on review and testing. They do not write production code. For language-specific code quality guidance -- idiomatic patterns, language-specific anti-patterns, framework conventions -- the [programming-languages](/programming-languages/) plugin provides complementary expertise. code-reviewer catches generic bugs; a language specialist catches language-specific issues.

Similarly, playwright-expert writes E2E tests but does not build the web application. The [web-development](/web-development/) plugin provides agents for frontend development. The quality agents review and test what the development agents build.

## See Also

- [Agent Reference]({{< ref "reference/agents" >}}) -- specifications for all five agents
- [Choosing the Right Agent]({{< ref "explanation/choosing-the-right-agent" >}}) -- decision guide for selecting the right agent
- [Getting Started]({{< ref "tutorials/getting-started" >}}) -- tutorial walkthrough of the review-test-fix cycle
