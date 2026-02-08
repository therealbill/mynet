---
title: "Marketplace Manifest"
description: "Specification of the marketplace.json format"
weight: 2
---

# Marketplace Manifest

**File location:** `.claude-plugin/marketplace.json` at the repository root
**Format:** JSON
**Stability:** Stable

## Purpose

The marketplace manifest declares the marketplace identity and lists all plugins available within it. Each repository contains exactly one marketplace manifest.

## Structure

```json
{
  "name": "string",
  "owner": {
    "name": "string"
  },
  "plugins": [
    {
      "name": "string",
      "source": "string",
      "description": "string",
      "version": "string",
      "keywords": ["string"],
      "author": {
        "name": "string"
      }
    }
  ]
}
```

## Fields

### Root Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Marketplace identifier. Unique name for this marketplace instance. |
| `owner` | object | Yes | Marketplace owner information. |
| `owner.name` | string | Yes | Name of the marketplace owner. |
| `plugins` | array | Yes | Array of plugin entry objects. |

### Plugin Entry Fields

Each element in the `plugins` array has the following fields.

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `name` | string | Yes | -- | Plugin directory name. Matches the top-level directory containing the plugin. |
| `source` | string | Yes | -- | Relative path to the plugin directory from the repository root (e.g., `./code-quality`). |
| `description` | string | Yes | -- | Short description of the plugin's purpose and contents. |
| `version` | string | No | -- | Semver version string (e.g., `1.0.0`, `1.1.0`). |
| `keywords` | string[] | No | `[]` | Array of keywords for plugin discovery and categorization. |
| `author` | object | No | -- | Plugin author information. |
| `author.name` | string | No | -- | Name of the plugin author. |

## Constraints

- The `name` field in each plugin entry must correspond to an existing top-level directory in the repository.
- The `source` field uses a relative path prefixed with `./`.
- The `version` field, when present, follows [Semantic Versioning 2.0.0](https://semver.org/).
- The `plugins` array contains one entry per plugin; duplicate `name` values are not permitted.

## Example

```json
{
  "name": "mynet",
  "owner": {
    "name": "Bill"
  },
  "plugins": [
    {
      "name": "code-quality",
      "source": "./code-quality",
      "description": "Code review, testing, accessibility, and architectural quality agents",
      "version": "1.0.0",
      "keywords": [
        "code-review",
        "testing",
        "accessibility",
        "architecture",
        "quality"
      ]
    },
    {
      "name": "timelord",
      "source": "./timelord",
      "description": "Temporal.io expertise for deploying, managing, and developing workflow applications",
      "version": "1.1.0",
      "author": {
        "name": "Bill"
      },
      "keywords": [
        "temporal",
        "temporal.io",
        "workflow",
        "orchestration"
      ]
    }
  ]
}
```

## See Also

- [Plugin Catalog](../plugin-catalog/) -- complete listing of all plugins
- [Plugin JSON](../plugin-json/) -- per-plugin manifest specification
- [Component Conventions](../component-conventions/) -- agent, skill, command, and template formats
