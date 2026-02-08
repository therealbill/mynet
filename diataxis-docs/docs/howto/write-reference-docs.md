---
title: "How to write reference documentation"
description: "Generate complete, advice-free reference documentation from source code using doc-reference-gen with consistent structure and full API coverage."
weight: 3
doc_type: how-to
prerequisites:
  - "The diataxis-docs plugin installed"
  - "Access to the source code, CLI tool, or API being documented"
  - "Knowledge of available introspection tools (typedoc, pydoc, go doc, etc.)"
est_time: "15 minutes"
roles: ["developer", "technical writer"]
stability: stable
---

# How to Write Reference Documentation

Generate complete, advice-free reference documentation from source code using the **doc-reference-gen** agent.

**Goal:** Produce reference pages that document every public API surface with consistent structure, no opinions, and no prescriptive language.

## Prerequisites

- The diataxis-docs plugin installed and available
- Access to the source code, CLI tool, or API being documented
- Knowledge of what introspection tools are available (typedoc, pydoc, go doc, cargo doc, OpenAPI specs)

## Steps

### 1. Identify the public API surface

Determine what needs reference documentation:

- **REST API**: endpoints, request/response schemas, error codes
- **CLI tool**: commands, subcommands, flags, exit codes
- **SDK/library**: public functions, types, interfaces, constants
- **Configuration**: all config options, types, defaults, constraints

### 2. Identify available introspection tools

Check what tools can extract API information automatically:

- TypeScript/JavaScript: `typedoc`, JSDoc comments
- Python: `pydoc`, type hints, docstrings
- Go: `go doc`, godoc comments
- Rust: `cargo doc`, rustdoc comments
- REST API: OpenAPI/Swagger spec files, GraphQL schema
- CLI: `--help` output, man pages

### 3. Invoke the doc-reference-gen agent

Pass the source path and introspection tool information:

```
Generate reference documentation for <project> at <source-path>.
Introspection tools available: <list tools>.
Place output at <docs-path>/reference/
```

The agent reads the source code, extracts public API surfaces, and generates structured reference pages.

### 4. Verify the "no advice" rule

This is the critical quality check for reference documentation. Scan every generated page for prescriptive language:

- **Remove** any instance of "you should," "we recommend," "it is best to," or "consider using"
- **Replace** prescriptive statements with factual ones. For example, change "You should validate inputs before calling this function" to "Input validation is the caller's responsibility"
- **Move** genuine advice to a how-to guide and link to it from the reference page's "See Also" section

### 5. Verify completeness

Check that every public symbol, command, or endpoint has a reference entry:

- Every public function or method is documented
- Every parameter has type, required/optional, default value, and constraints
- Every error code has its condition and HTTP status (for APIs)
- Every CLI flag has its long form, short form, type, and default
- Rate limits, quotas, and size constraints are documented

### 6. Verify consistent structure

All reference entries of the same kind must follow the same structure. For API endpoints:

```
# Endpoint Name
## Signature
## Parameters (table: name, type, required, default, description)
## Returns
## Errors (table: code, condition, status)
## Limits
## Examples (syntax only, not workflows)
## See Also
```

For CLI commands:

```
# Command Name
## Synopsis
## Arguments
## Options (table: flag, type, default, description)
## Exit Codes
## Examples
## See Also
```

### 7. Add cross-links

Reference pages link to:

- How-to guides that demonstrate usage of the documented API or command
- Related reference pages for associated types or commands
- Explanation pages that discuss the design rationale

Use Hugo ref shortcodes: `{{</* ref "howto/some-guide" */>}}`

### 8. Set version and stability

Every reference page must include in its frontmatter:

- `stability`: one of `stable`, `beta`, `experimental`, `deprecated`
- `version`: the version of the API or tool being documented

Deprecated entries must include a migration path in the page body.

## Verify It Works

Correctly written reference documentation meets these criteria:

- Every public API surface has a documented entry
- Zero instances of prescriptive language ("you should," "we recommend")
- All entries of the same kind follow identical structure
- Parameters include type, required/optional, default, and constraints
- Error codes are documented with conditions
- Examples show syntax, not workflows or use cases

## Troubleshooting

**The agent generated advice in reference pages.**
Re-run the agent on the specific page with explicit instructions to enforce the "no advice" rule. Alternatively, manually extract advice into a how-to guide and replace it with a factual statement plus a "See Also" link.

**Some public symbols are missing from the generated docs.**
Check whether the symbols are exported or public. The agent documents the public API surface. If symbols are intentionally public but were missed, point the agent at the specific source files containing them.

**The generated structure is inconsistent across pages.**
Invoke the agent again with a specific template to follow. Reference the structure from the first correctly generated page as the template for the rest.

**Examples show workflows instead of syntax.**
Replace workflow examples with minimal syntax demonstrations. A reference example for a function shows how to call it with representative arguments, not how to build a feature using it.

## See Also

- [doc-reference-gen agent specification](../../reference/agents/)
- [Why reference is information-oriented with no advice](../../explanation/diataxis-in-practice/)
- [Validate reference pages pass quality checks](../../howto/validate-doc-quality/)
- [Extract reference content from mixed-type pages](../../howto/restructure-docs-to-diataxis/)
