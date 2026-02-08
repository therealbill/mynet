---
title: "Plugin JSON"
description: "Specification of the plugin.json manifest format"
weight: 3
---

# Plugin JSON

**File location:** `.claude-plugin/plugin.json` within each plugin directory
**Format:** JSON
**Stability:** Stable

## Purpose

The plugin manifest declares a plugin's identity, version, and metadata. Each plugin directory contains exactly one `plugin.json` file at `.claude-plugin/plugin.json`.

## Structure

```json
{
  "name": "string",
  "version": "string",
  "description": "string",
  "keywords": ["string"],
  "author": {
    "name": "string",
    "url": "string"
  },
  "repository": "string",
  "license": "string"
}
```

## Fields

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `name` | string | Yes | -- | Plugin identifier. Matches the plugin's directory name. Uses kebab-case. |
| `version` | string | Yes | -- | Semver version string (e.g., `1.0.0`). |
| `description` | string | Yes | -- | Short description of the plugin's purpose and contents. |
| `keywords` | string[] | No | `[]` | Discovery keywords for categorization and search. |
| `author` | object | No | -- | Author information. |
| `author.name` | string | No | -- | Author name. |
| `author.url` | string | No | -- | Author URL (e.g., GitHub profile). |
| `repository` | string | No | -- | Source repository URL. |
| `license` | string | No | -- | License identifier (e.g., `MIT`) or reference (e.g., `SEE LICENSE IN ../LICENSE`). |

## Constraints

- The `name` field value must match the containing directory name exactly.
- The `version` field follows [Semantic Versioning 2.0.0](https://semver.org/).
- The `name` field uses kebab-case (lowercase, hyphen-separated).

## Examples

### Minimal

```json
{
  "name": "code-quality",
  "version": "1.0.0",
  "description": "Code review, testing, accessibility, and architectural quality agents",
  "keywords": ["code-review", "testing", "accessibility", "architecture", "quality"]
}
```

### Extended

```json
{
  "name": "timelord",
  "version": "1.1.0",
  "description": "Temporal.io expertise for deploying, managing, and developing workflow applications with self-hosted Temporal clusters",
  "author": {
    "name": "Bill",
    "url": "https://github.com/therealbill"
  },
  "repository": "https://github.com/therealbill/mynet",
  "license": "SEE LICENSE IN ../LICENSE",
  "keywords": ["temporal", "temporal.io", "workflow", "orchestration"]
}
```

## Directory Layout

The `plugin.json` file sits inside a `.claude-plugin` subdirectory within the plugin's top-level directory.

```
<plugin-name>/
  .claude-plugin/
    plugin.json        <-- this file
  agents/
  skills/
  commands/
  templates/
```

## See Also

- [Plugin Catalog](../plugin-catalog/) -- complete listing of all plugins
- [Marketplace Manifest](../marketplace-manifest/) -- root marketplace manifest specification
- [Component Conventions](../component-conventions/) -- agent, skill, command, and template formats
