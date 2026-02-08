# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Purpose

Mynet is a Claude Code plugin marketplace. Each top-level directory is an independent plugin with its own versioning, licensing, and domain focus. The marketplace manifest at `.claude-plugin/marketplace.json` ties them together.

## Architecture

```
mynet/
├── .claude-plugin/marketplace.json    # Marketplace manifest — lists all plugins
├── <plugin-name>/                     # Each plugin is a self-contained directory
│   ├── .claude-plugin/plugin.json     # Plugin metadata (name, version, description, keywords)
│   ├── agents/                        # Agent definitions (*.md with YAML frontmatter)
│   ├── skills/                        # Skill definitions (*/SKILL.md with YAML frontmatter)
│   ├── commands/                      # Slash command definitions (*.md)
│   ├── templates/                     # Code generation templates (optional)
│   ├── tools/                         # CLI tools or utilities (optional)
│   └── docs/                          # Plugin documentation (optional)
```

## Plugin Component Conventions

- **`.claude-plugin/plugin.json`** — Plugin manifest with name, version, description, keywords
- **`agents/*.md`** — Agent definitions with YAML frontmatter (name, description, tools, model) and system prompt body
- **`skills/*/SKILL.md`** — Skill definitions with YAML frontmatter (name, description, trigger patterns) and instruction body
- **`commands/*.md`** — Slash command definitions with YAML frontmatter and implementation instructions
- **`templates/`** — Code generation templates, format depends on the plugin's domain

## Adding a New Plugin

1. Create a new top-level directory named after the plugin
2. Add `.claude-plugin/plugin.json` with at minimum: name, version, description
3. Add the plugin entry to the root `.claude-plugin/marketplace.json`
4. Each plugin is independently versioned using semver with its own CHANGELOG.md

## Marketplace Manifest

The root `.claude-plugin/marketplace.json` contains:

- `name` — Marketplace identifier (`mynet`)
- `owner` — Marketplace owner
- `plugins[]` — Array of plugin entries with name, source path, description, and optional version/keywords

## Documentation

When a plugin includes documentation, it should follow the Diataxis framework (tutorials, how-to guides, explanation, reference). Use Diataxis skills/agents when creating or modifying plugin docs.

## Licensing

The repository uses a custom aggregation-prohibiting license (see root `LICENSE`). It prohibits bundling in multi-plugin collections from multiple sources but permits marketplace listing with reference to the original repo and individual redistribution with attribution. Individual plugins may have additional licensing terms in their own LICENSE files.
