# Validation and Inventory Report Templates

Use these templates when generating reports for doc-inventory and doc-crosslink-validator.

## Inventory JSON Template

Generate as `docs/_reports/inventory.json`:

```json
{
  "scanned_at": "ISO-8601 timestamp",
  "docs_path": "/path/to/docs",
  "total_pages": 42,
  "by_type": {
    "tutorial": 2,
    "how-to": 8,
    "reference": 25,
    "explanation": 4,
    "mixed": 3
  },
  "pages": [
    {
      "path": "docs/getting-started.md",
      "title": "Getting Started",
      "type": "tutorial",
      "quality": "strong",
      "issues": [],
      "word_count": 1200,
      "has_frontmatter": true,
      "links_to": ["docs/api/client.md"],
      "linked_from": ["docs/index.md"]
    }
  ],
  "gaps": [
    {
      "type": "tutorial",
      "severity": "critical",
      "description": "No onboarding tutorial exists",
      "suggestion": "Create docs/tutorials/00-onboarding.md"
    }
  ],
  "reorganization": [
    {
      "action": "move",
      "from": "docs/setup.md",
      "to": "docs/tutorials/setup.md",
      "reason": "Content is tutorial-style walkthrough"
    },
    {
      "action": "split",
      "file": "docs/advanced.md",
      "into": ["docs/how-to/advanced-config.md", "docs/explanation/architecture.md"],
      "reason": "Mixes task instructions with architectural concepts"
    }
  ],
  "recommendations": [
    "Create onboarding tutorial as highest priority",
    "Split mixed content in docs/advanced.md",
    "Add front-matter to all pages",
    "Generate complete API reference from source"
  ]
}
```

### Quality ratings

- **strong** — Page cleanly matches one Diataxis type, well-structured, complete
- **moderate** — Page mostly matches, minor issues (missing troubleshooting, weak cross-links)
- **weak** — Page nominally fits but has significant gaps (no checkpoints in tutorial, no verification in how-to)
- **mixed** — Page contains substantial content from 2+ types and should be split

### Gap severity levels

- **critical** — Missing onboarding tutorial, no reference for core API
- **high** — Top user tasks lack how-to guides, key concepts unexplained
- **medium** — Secondary features undocumented
- **low** — Edge cases or advanced topics lacking docs

---

## Quality Report JSON Template

Generate as `docs/_reports/quality.json`:

```json
{
  "validated_at": "ISO-8601 timestamp",
  "docs_path": "/path/to/docs",
  "total_pages": 42,
  "summary": {
    "errors": 5,
    "warnings": 12,
    "info": 8,
    "passed": 25
  },
  "by_type": {
    "tutorial": { "total": 3, "errors": 0, "warnings": 1 },
    "how-to": { "total": 15, "errors": 2, "warnings": 5 },
    "reference": { "total": 20, "errors": 2, "warnings": 4 },
    "explanation": { "total": 4, "errors": 1, "warnings": 2 }
  },
  "issues": [
    {
      "file": "docs/how-to/setup-auth.md",
      "severity": "error",
      "type": "missing_frontmatter",
      "message": "Missing required field: est_time",
      "line": 1
    },
    {
      "file": "docs/reference/api/users.md",
      "severity": "error",
      "type": "type_violation",
      "message": "Reference contains advice: 'You should always validate...'",
      "line": 45,
      "context": "You should always validate email addresses..."
    },
    {
      "file": "docs/tutorials/getting-started.md",
      "severity": "warning",
      "type": "broken_link",
      "message": "Link to non-existent page: ../how-to/advanced-config.md",
      "line": 120
    },
    {
      "file": "docs/how-to/deploy.md",
      "severity": "info",
      "type": "missing_crosslink",
      "message": "No link to related reference docs/reference/cli/deploy.md",
      "suggestion": "Add link to CLI reference for 'deploy' command"
    }
  ],
  "link_graph": {
    "orphans": ["docs/old/deprecated.md"],
    "dead_links": [
      { "from": "docs/how-to/example.md", "to": "docs/nonexistent.md" }
    ],
    "cross_type_links": {
      "tutorial_to_reference": 15,
      "howto_to_explanation": 8,
      "reference_to_howto": 12,
      "explanation_to_howto": 6
    }
  },
  "recommendations": [
    "Fix 5 errors before publishing",
    "Add front-matter to 3 pages",
    "Remove advice from 2 reference pages",
    "Fix 4 broken links",
    "Add cross-links: how-tos should link to related reference"
  ]
}
```

---

## Quality Report Markdown Template

Generate as `docs/_reports/quality-report.md`:

```markdown
# Documentation Quality Report

**Validated**: [timestamp]
**Total pages**: [N]

## Summary

- [N] pages passed all checks
- [N] warnings need attention
- [N] errors must be fixed

## By Type

| Type | Total | Passed | Warnings | Errors |
|------|-------|--------|----------|--------|
| Tutorial | N | N | N | N |
| How-to | N | N | N | N |
| Reference | N | N | N | N |
| Explanation | N | N | N | N |

## Critical Issues (Must Fix)

### 1. [Category]

**Files affected**: N

- `path/to/file.md` — [specific issue]
- `path/to/file.md:line` — [specific issue with context]

### 2. [Category]

...

## Warnings

### Broken links: N

- path/to/file.md:line -> path/to/target.md (doesn't exist)

### Missing cross-links: N

- `path/to/file.md` should link to `path/to/related.md`

### Orphaned pages: N

- `path/to/orphan.md` — not linked from anywhere

## Link Graph Analysis

Cross-type links:

- Tutorial -> Reference: N
- How-to -> Explanation: N
- Reference -> How-to: N
- Explanation -> How-to: N

## Recommendations

1. [Prioritized action]
2. [Prioritized action]
3. ...

## Top Quality Pages

- `path/to/excellent.md` — [why it's good]
```

---

## Issue Types Reference

| Type | Severity | Description |
|------|----------|-------------|
| `missing_frontmatter` | error | Required frontmatter field missing |
| `no_frontmatter` | error | No frontmatter block at all |
| `type_violation` | error | Content violates its Diataxis type (advice in reference, etc.) |
| `broken_link` | error | Internal link points to non-existent file |
| `invalid_stability` | error | Stability value not one of: stable, beta, experimental, deprecated |
| `missing_crosslink` | warning | Page should link to related content but doesn't |
| `orphaned_page` | warning | Page has no incoming links |
| `weak_structure` | warning | Missing expected sections (no troubleshooting in how-to, no checkpoint in tutorial) |
| `over_length` | warning | How-to exceeds 1800 words, page exceeds 3000 words |
| `crosslink_opportunity` | info | Page could benefit from additional cross-type links |
| `split_candidate` | info | Mixed-type page that should be split |
| `terminology_inconsistency` | info | Same concept referred to by different names |
