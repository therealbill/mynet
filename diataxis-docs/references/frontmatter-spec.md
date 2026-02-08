# Frontmatter Specification

Use this reference when validating or creating documentation page frontmatter.

## Required Fields

Every documentation page must include these frontmatter fields:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `title` | string | Yes | Page title. Should reflect the doc type (see naming below). |
| `summary` | string | Yes | One-sentence description of what this page covers. |
| `doc_type` | string | Yes | One of: `tutorial`, `how-to`, `reference`, `explanation` |
| `stability` | string | Yes | One of: `stable`, `beta`, `experimental`, `deprecated` |

## Recommended Fields

| Field | Type | Applicable Types | Description |
|-------|------|-----------------|-------------|
| `prerequisites` | list | tutorial, how-to | Links to required prior knowledge or setup |
| `est_time` | string | tutorial, how-to | Estimated completion time (e.g., "15 minutes") |
| `roles` | list | all | Target audience roles (e.g., `["developer", "operator"]`) |
| `version` | string | reference | API/tool version this documents |

## Title Naming Conventions

Each doc type follows a specific naming pattern:

| Type | Pattern | Examples |
|------|---------|----------|
| Tutorial | "Build/Create your first [thing]" | "Build your first workflow", "Create your first API" |
| How-to | "How to [verb phrase]" | "How to configure OAuth", "How to deploy to production" |
| Reference | "[Thing] reference" or "[Thing] API" | "CLI reference", "REST API", "Configuration reference" |
| Explanation | "Understanding [concept]" | "Understanding the event system", "Architecture overview" |

## Stability Values

| Value | Meaning | Visual indicator |
|-------|---------|-----------------|
| `stable` | API/feature is stable, breaking changes follow semver | None (default) |
| `beta` | API may change, feedback welcome | Beta badge |
| `experimental` | API is unstable, may be removed | Warning banner |
| `deprecated` | Scheduled for removal, migration path provided | Deprecation notice with replacement link |

## Example Frontmatter by Type

### Tutorial

```yaml
---
title: "Build your first data pipeline"
summary: "A hands-on tutorial that takes you from installation to a working data pipeline in 90 minutes."
doc_type: tutorial
prerequisites:
  - "Python 3.10+ installed"
  - "Basic command-line familiarity"
est_time: "90 minutes"
roles: ["developer"]
stability: stable
---
```

### How-to

```yaml
---
title: "How to configure retry policies"
summary: "Set up automatic retry with exponential backoff for failed workflow activities."
doc_type: how-to
prerequisites:
  - "Completed the [Getting Started tutorial](../tutorials/getting-started.md)"
  - "A running workflow worker"
est_time: "10 minutes"
roles: ["developer"]
stability: stable
---
```

### Reference

```yaml
---
title: "Workflow API Reference"
summary: "Complete specification of the Workflow REST API including all endpoints, parameters, and error codes."
doc_type: reference
stability: stable
version: "2.4.0"
---
```

### Explanation

```yaml
---
title: "Understanding workflow execution semantics"
summary: "How the workflow engine achieves exactly-once execution through event sourcing and deterministic replay."
doc_type: explanation
prerequisites:
  - "Familiarity with basic workflow concepts (see [Getting Started](../tutorials/getting-started.md))"
est_time: "20 minutes"
roles: ["developer", "architect"]
stability: stable
---
```

## Validation Rules

When validating frontmatter:

1. **All required fields present** — title, summary, doc_type, stability
2. **doc_type is valid** — one of the four Diataxis types
3. **stability is valid** — one of: stable, beta, experimental, deprecated
4. **Title matches type convention** — tutorials say "Build/Create", how-tos say "How to", etc.
5. **Tutorials and how-tos have prerequisites** — even if the list is empty, the field should exist
6. **Tutorials and how-tos have est_time** — readers need to know the time commitment
7. **Reference pages have version** — readers need to know which version is documented
8. **Deprecated pages have migration link** — in the body, not just the frontmatter
