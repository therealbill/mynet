---
name: Hugo Module Mounts
description: >-
  This skill should be used when the user asks about Hugo module mounts,
  mounting multiple directories into Hugo's content tree, aggregating docs from
  scattered directories, using hugo mod init, configuring module.mounts in
  hugo.toml, building a unified documentation site from a monorepo, mapping
  subdirectory docs/ folders into a single Hugo site, or understanding Hugo's
  union file system. Also applies when the user has a repository with multiple
  docs/ directories that need to be presented as a single site. This is the
  core skill for multi-directory content aggregation.
version: 1.0.0
---

## Overview

Hugo's module mounts map arbitrary directories into Hugo's content tree using the `[[module.mounts]]` configuration. This creates a union file system — multiple source directories appear as a single content tree without symlinks or build-time copying. This is Hugo's official mechanism for multi-directory content aggregation.

Use module mounts when a repository has documentation scattered across multiple subdirectories (per-plugin, per-service, per-tool) and needs a unified Hugo site.

## Prerequisites

Hugo modules require Go and module initialization:

```bash
# Initialize Hugo modules (required once per project)
hugo mod init github.com/username/repo
```

This creates `go.mod` and `go.sum` files. Commit both to the repository.

## Basic Mount Configuration

In `hugo.toml`, mounts map a `source` directory to a `target` location in Hugo's virtual filesystem:

```toml
[[module.mounts]]
source = 'content'
target = 'content'

[[module.mounts]]
source = 'static'
target = 'static'

[[module.mounts]]
source = 'layouts'
target = 'layouts'

[[module.mounts]]
source = 'data'
target = 'data'

[[module.mounts]]
source = 'assets'
target = 'assets'

[[module.mounts]]
source = 'i18n'
target = 'i18n'

[[module.mounts]]
source = 'archetypes'
target = 'archetypes'
```

When you define any mount, Hugo's default mounts are replaced. Always include the standard mounts above when adding custom mounts, or the default directories will stop working.

## Multi-Directory Content Aggregation

The core pattern: mount each subdirectory's `docs/` folder into the content tree.

### Example: Plugin Marketplace

Repository structure:

```
repo/
├── hugo.toml
├── content/           # Root site content (home page, about, etc.)
├── plugin-a/
│   ├── docs/          # Plugin A documentation
│   │   ├── _index.md
│   │   └── usage.md
│   └── src/
├── plugin-b/
│   ├── docs/
│   │   ├── _index.md
│   │   └── api.md
│   └── src/
└── plugin-c/
    ├── docs/
    └── src/
```

Mount configuration:

```toml
# Standard mounts (required when defining any custom mounts)
[[module.mounts]]
source = 'content'
target = 'content'

[[module.mounts]]
source = 'static'
target = 'static'

[[module.mounts]]
source = 'layouts'
target = 'layouts'

[[module.mounts]]
source = 'data'
target = 'data'

[[module.mounts]]
source = 'assets'
target = 'assets'

# Plugin documentation mounts
[[module.mounts]]
source = 'plugin-a/docs'
target = 'content/plugins/plugin-a'

[[module.mounts]]
source = 'plugin-b/docs'
target = 'content/plugins/plugin-b'

[[module.mounts]]
source = 'plugin-c/docs'
target = 'content/plugins/plugin-c'
```

Result: Hugo sees a unified content tree:

```
content/
├── _index.md                    # From content/
├── plugins/
│   ├── plugin-a/
│   │   ├── _index.md           # From plugin-a/docs/
│   │   └── usage.md
│   ├── plugin-b/
│   │   ├── _index.md           # From plugin-b/docs/
│   │   └── api.md
│   └── plugin-c/
│       └── ...
```

### Example: Monorepo Services

```toml
# Service documentation mounts
[[module.mounts]]
source = 'services/api-gateway/docs'
target = 'content/services/api-gateway'

[[module.mounts]]
source = 'services/auth-service/docs'
target = 'content/services/auth-service'

[[module.mounts]]
source = 'services/worker/docs'
target = 'content/services/worker'
```

### Example: CLI Tool Collection

```toml
# Tool documentation mounts
[[module.mounts]]
source = 'tools/migrate/docs'
target = 'content/tools/migrate'

[[module.mounts]]
source = 'tools/validate/docs'
target = 'content/tools/validate'

[[module.mounts]]
source = 'tools/generate/docs'
target = 'content/tools/generate'
```

## Section Index Pages

Each mounted directory needs an `_index.md` to appear as a proper section. The parent directory also needs one.

Create `content/plugins/_index.md` (or `content/services/_index.md`):

```markdown
---
title: "Plugins"
description: "Documentation for all plugins"
weight: 10
---

Browse documentation for individual plugins below.
```

Each mounted directory needs its own `_index.md` (inside the source `docs/` folder):

```markdown
---
title: "Plugin A"
description: "Plugin A documentation"
weight: 1
---

Plugin A overview content.
```

## Mount Precedence

When files conflict (same path from multiple mounts), the first mount wins:

```toml
# This mount takes precedence for any conflicting paths
[[module.mounts]]
source = 'content'
target = 'content'

# This mount's files are used only when they don't conflict with above
[[module.mounts]]
source = 'shared-content'
target = 'content'
```

Order mounts from highest to lowest priority.

## Mounting Non-Content Directories

Mounts work for any Hugo component directory, not just content:

```toml
# Mount shared layouts
[[module.mounts]]
source = 'layouts'
target = 'layouts'

[[module.mounts]]
source = 'shared/layouts'
target = 'layouts'

# Mount data files from a subdirectory
[[module.mounts]]
source = 'data'
target = 'data'

[[module.mounts]]
source = 'plugin-a/data'
target = 'data/plugins/plugin-a'
```

Valid mount targets: `archetypes`, `assets`, `content`, `data`, `i18n`, `layouts`, `static`.

## Adding a New Mount

When adding documentation for a new subdirectory:

1. Create the `docs/` directory with `_index.md`
2. Add the mount entry to `hugo.toml`
3. Ensure the parent section has `_index.md`
4. Run `hugo server` to verify

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Defining custom mounts without standard mounts | Always include mounts for `content`, `static`, `layouts`, `data`, `assets` when adding any custom mount |
| Missing `_index.md` in mounted directories | Every mounted directory needs `_index.md` to be a proper Hugo section |
| Missing parent section `_index.md` | If mounting to `content/plugins/foo`, ensure `content/plugins/_index.md` exists |
| Forgetting `hugo mod init` | Module mounts require module initialization |
| Absolute paths in source | Use paths relative to the repository root |
| Not committing `go.mod` and `go.sum` | These are required for reproducible builds |

## Debugging Mounts

```bash
# Verify Hugo sees all mounted content
hugo list all

# Check which files Hugo resolves
hugo config mounts

# Verbose build output to see mount resolution
hugo --verbose
```
