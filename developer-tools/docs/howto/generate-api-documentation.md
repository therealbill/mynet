---
title: "Generate API Documentation"
description: "Use documentation-engineer to generate docs from source code"
weight: 2
---

# Generate API Documentation

Generate structured API reference documentation from an existing codebase using documentation-engineer.

## Goal

By the end of this guide, your project will have:

- API reference documentation generated from your source code
- Documentation organized by resource and endpoint
- CI checks configured to keep docs in sync with code changes

## Step 1: Point documentation-engineer at Your Codebase

Activate the agent by describing what needs documentation:

```
We need a docs site for our SDK
```

```
REST API has 40 endpoints and zero documentation
```

```
Generate docs for the API in src/routes/
```

Be specific about which part of the codebase needs documentation. If your project has both a REST API and a CLI tool, tell the agent which one to document first. documentation-engineer works better when focused on a single documentation target per session.

The agent begins with a codebase audit rather than generating docs immediately. It reads your project structure, identifies the technology stack, and locates the relevant source files before writing anything.

## Step 2: Review the Audit

documentation-engineer produces an audit covering:

- **Routes and endpoints** -- every HTTP route with its method, path, and handler location
- **Request and response schemas** -- parameter types, validation rules, and response shapes inferred from code
- **Authentication patterns** -- middleware, token validation, role checks
- **Error handling** -- error response formats, status codes, and common error conditions
- **Dependencies** -- external services, databases, and third-party APIs that endpoints interact with

Review the audit for accuracy. If the agent missed endpoints (common with dynamic route registration or plugin-based architectures), point them out. The agent re-scans with that guidance.

## Step 3: Review Generated Reference Docs

After the audit, documentation-engineer generates reference documentation structured by resource:

- Each resource gets its own section (Users, Orders, Products, etc.)
- Each endpoint within a resource includes method, path, description, parameters, request body, response format, authentication requirement, and example
- Error responses are documented alongside success responses
- Authentication requirements are noted per endpoint, not just globally

The generated docs follow the Diataxis reference quadrant -- they describe what things are and how they behave, without tutorial-style walkthroughs or explanatory digressions. If you need tutorials or conceptual documentation alongside the reference, use the [diataxis-docs](/diataxis-docs/) plugin for the broader documentation structure.

The agent selects a documentation tooling approach based on your project:

- **Markdown in `docs/`** -- for projects where simplicity is paramount
- **MkDocs** -- for Python projects or teams wanting Material for MkDocs theming
- **Docusaurus** -- for JavaScript/TypeScript projects wanting React-powered docs
- **Sphinx** -- for Python projects needing autodoc from docstrings

If the agent picks a tool that does not fit your needs, specify your preference.

## Step 4: Set Up CI Checks

documentation-engineer configures CI validation to prevent docs from drifting out of sync:

- **Link checker** -- verifies all internal and external links resolve
- **Code sample validator** -- ensures code examples in docs actually compile or run
- **Schema drift detection** -- compares documented endpoints against actual route definitions and flags discrepancies
- **Build verification** -- ensures the documentation site builds without errors on every PR

These checks run as part of your existing CI pipeline. The agent adds a workflow step rather than creating a separate pipeline.

## Verification

Confirm the documentation is accurate:

1. Pick 3-5 endpoints at random and compare the documented behavior against the actual API
2. Verify request parameters match what the code validates
3. Verify response shapes match what the code returns
4. Run the CI checks and confirm they pass
5. Make a small API change (add a query parameter) without updating docs -- verify CI catches the drift

If all checks pass and the documentation matches observed behavior, the docs are production-ready.

## Troubleshooting

**Missing endpoints.** If your framework uses dynamic route registration, decorator-based routing, or plugin architectures, the agent may not detect all routes through static analysis. List the missing endpoints explicitly and the agent adds them. For large APIs, provide an OpenAPI spec or route dump if available.

**Docs are too sparse.** documentation-engineer generates from what it can infer from code. If your endpoints lack descriptive handler names, clear parameter types, or structured error responses, the generated docs reflect that sparsity. Improve the source code clarity first, then regenerate.

**When to use diataxis-docs instead.** If you need a complete documentation site with tutorials, how-to guides, explanations, and reference -- not just API reference -- use the [diataxis-docs](/diataxis-docs/) plugin. documentation-engineer handles the reference quadrant. diataxis-docs handles the full Diataxis framework across all four quadrants.

**Duplicate information.** documentation-engineer avoids duplicating information already present in well-written code comments. If your code has comprehensive JSDoc or docstring coverage, the agent extracts and structures that content rather than paraphrasing it. It will not create a separate description that says the same thing as the inline comment.
