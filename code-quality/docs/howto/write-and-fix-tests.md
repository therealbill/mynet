---
title: "Write and Fix Tests"
description: "Use test-writer-fixer to create tests and repair broken test suites"
weight: 2
---

# Write and Fix Tests

Add test coverage to untested code or repair a broken test suite using the test-writer-fixer agent.

## Prerequisites

- Claude Code with the code-quality plugin installed
- A project with a test framework already configured (test-writer-fixer follows existing conventions)

## Steps

### 1. Identify Code That Needs Tests

Determine what needs testing. Common scenarios:

- **New code with no tests:** A module, function, or API endpoint that was written without corresponding tests.
- **Failing tests after a refactor:** Tests that broke because the code under test changed behavior or structure.
- **Low coverage area:** A critical path (authentication, payment, data validation) that lacks sufficient test coverage.

### 2. Trigger test-writer-fixer

Use context-specific phrases to activate the agent:

For writing new tests:

```
Our payment module has no tests
```

```
Write tests for the authentication middleware
```

For fixing broken tests:

```
Tests are broken after my refactor -- can you fix them?
```

For proactive test maintenance after code changes:

```
I've updated the authentication logic
```

test-writer-fixer identifies affected test files via import relationships, reads your project's existing test patterns, and either creates new tests or diagnoses and fixes failing ones.

### 3. Understand Failure Classification

When tests fail, test-writer-fixer classifies each failure before acting:

- **Real bug:** The test caught a genuine problem in the production code. The agent reports the bug and does not modify the test. You need to fix the code.
- **Stale expectation:** The code behavior legitimately changed, but the test still asserts the old behavior. The agent updates the assertion to match the new, correct behavior.
- **Brittle test:** The test depends on implementation details (internal state, execution order, exact timing). The agent refactors the test to assert on observable behavior while still catching the same class of bugs.
- **Environment issue:** Missing dependencies, database state, configuration files, or service connectivity. The agent fixes the test setup, not the test logic.

This classification prevents the common mistake of weakening tests just to make them pass.

### 4. Review the Fixes

After test-writer-fixer modifies tests, review the changes:

- **Check that test intent is preserved.** Read the test name and the surrounding test suite. A test named `"rejects expired tokens"` should still verify token expiration after the fix, even if the assertion details changed.
- **Verify behavioral assertions.** Tests should assert on return values, side effects, and error conditions -- not on internal method calls or implementation structure.
- **Confirm naming conventions.** New tests should follow your project's naming style. test-writer-fixer mirrors the patterns it finds in existing tests.

### 5. Verification

Run the full test suite to confirm everything passes:

```
Run all the tests
```

test-writer-fixer runs fixed tests multiple times to detect flakiness. If a test passes intermittently, the agent identifies the source of non-determinism and fixes it.

A clean run means the test suite is stable and the changes are safe to commit.

## Troubleshooting

**Tests that are fundamentally brittle:**

Some tests are brittle by design -- they assert on snapshot output, timestamp values, or random data. test-writer-fixer refactors these to use deterministic comparisons (freezing time, seeding randomness, using partial matching on snapshots). If the brittleness is inherent to what is being tested, the agent adds clear comments explaining why the test may need periodic updates.

**Mock-heavy tests that test nothing:**

Tests that mock every dependency end up verifying mock wiring rather than real behavior. test-writer-fixer identifies over-mocked tests and suggests which mocks to remove in favor of real implementations or lightweight fakes. The goal is to test the behavior boundary, not the internal call graph.

**Tests that require external services:**

test-writer-fixer does not set up external databases or APIs. If tests need infrastructure, the agent creates setup instructions and configures test fixtures or in-memory alternatives where possible. It distinguishes between unit tests (no external dependencies) and integration tests (require running services).

## See Also

- [Getting Started](../../tutorials/getting-started/) -- tutorial walkthrough of the review-test-fix cycle
- [Agent Reference](../../reference/agents/) -- test-writer-fixer specification
- [Choosing the Right Agent](../../explanation/choosing-the-right-agent/) -- when to use test-writer-fixer versus playwright-expert
