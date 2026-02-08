---
title: "Set Up a Hugo Site with Module Mounts"
description: "How to set up a Hugo site that aggregates documentation from multiple directories in a repository using module mounts"
weight: 1
---

# Set Up a Hugo Site with Module Mounts

This guide covers setting up a Hugo site for a repository where documentation lives in multiple subdirectories. The hugo-site-architect agent and hugo-module-mounts skill handle this scenario.

## When to use this

Use this approach when your repository has documentation spread across multiple directories, such as:

- A monorepo with per-service `docs/` folders
- A plugin marketplace with per-plugin documentation
- A tool collection where each tool has its own docs

## Steps

### 1. Ask the hugo-site-architect agent to analyze your repository

Tell Claude what you need:

```
Set up a Hugo site for this repo -- we have docs in each plugin directory
```

The hugo-site-architect agent scans the repository structure using its Read, Glob, and Grep tools to find all directories containing markdown documentation. It identifies the repository type and proposes a content hierarchy.

### 2. Review the proposed structure

The agent presents its findings before making changes:

- All discovered documentation directories
- The proposed mount mapping (which directories mount where)
- The recommended site hierarchy
- Theme recommendation based on the project's needs

Confirm the structure or request adjustments before proceeding.

### 3. Let the agent scaffold the site

Once you approve the structure, the agent:

- Runs `hugo new site . --force` and `hugo mod init`
- Generates `hugo.toml` with all module mounts
- Creates section index pages (`_index.md`) for every section

The generated `hugo.toml` includes both standard mounts and custom content mounts. The standard mounts are required whenever custom mounts are defined:

```toml
# Standard mounts (always required with custom mounts)
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
source = 'archetypes'
target = 'archetypes'

# Custom content mounts
[[module.mounts]]
source = 'plugin-a/docs'
target = 'content/plugins/plugin-a'

[[module.mounts]]
source = 'plugin-b/docs'
target = 'content/plugins/plugin-b'
```

### 4. Verify the content tree

Run Hugo's list command to confirm all mounted content is visible:

```bash
hugo list all
```

Every mounted directory should appear under its target path. If any content is missing, check that:

- The source directory exists and contains markdown files
- The source directory has an `_index.md` file
- Parent sections (such as `content/plugins/`) have their own `_index.md`
- The mount entry uses a path relative to the repository root

### 5. Install and configure a theme

The agent recommends a theme based on your project's needs. Common documentation themes include Book, Docsy, and Geekdoc. Install using Hugo modules:

```bash
hugo mod get github.com/theNewDynamic/gohugo-theme-ananke
```

Configure in `hugo.toml`:

```toml
[module]
  [[module.imports]]
    path = 'github.com/theNewDynamic/gohugo-theme-ananke'
```

### 6. Start the development server

Preview the site locally:

```bash
hugo server --buildDrafts
```

Navigate to `http://localhost:1313/` and verify all sections appear in the navigation and content renders correctly.

## Key considerations

- **Standard mounts are required**: When you define any custom mount, Hugo's default mounts are replaced. You must explicitly include mounts for content, static, layouts, data, assets, and archetypes.
- **Every section needs `_index.md`**: Both the mounted directories and their parent sections need section index pages.
- **Mount precedence**: When files conflict across mounts, the first mount listed takes precedence.
- **Commit `go.mod` and `go.sum`**: These files are required for reproducible builds in CI.

## Related

- {{< ref "howto/add-content-section" >}} -- Add more directories after initial setup
- {{< ref "reference/agents" >}} -- hugo-site-architect agent specification
- {{< ref "reference/skills" >}} -- hugo-module-mounts skill specification
- {{< ref "explanation/architecture" >}} -- How skills and agents work together
