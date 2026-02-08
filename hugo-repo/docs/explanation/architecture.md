---
title: "Architecture"
description: "How the 17 components of the hugo-repo plugin fit together across four component types"
weight: 1
---

# Architecture

The hugo-repo plugin contains 17 components organized into four types: 7 skills, 2 agents, 4 commands, and 4 templates. Each type serves a distinct role in the plugin's architecture, and the types are designed to complement each other rather than overlap.

## The four component types

### Skills: domain knowledge

Skills are passive knowledge components. They contain structured information about Hugo topics that Claude uses when responding to user questions. Skills do not take actions -- they inform Claude's understanding so that responses are accurate and complete.

The 7 skills cover the full Hugo lifecycle:

| Skill | Domain |
|-------|--------|
| hugo-fundamentals | Site structure, configuration, content organization, development workflow |
| hugo-content-authoring | Shortcodes, render hooks, taxonomies, page bundles |
| hugo-data-templates | Data directory, data-driven patterns, remote data |
| hugo-themes | Theme installation, customization, override patterns, asset pipeline |
| hugo-module-mounts | Multi-directory content aggregation, union file system |
| hugo-github-actions | GitHub Pages deployment, CI/CD workflows, path filtering |
| hugo-s3-deployment | AWS S3 hosting, CloudFront, OIDC authentication |

Skills are triggered by pattern matching against the user's input. When a user asks about Hugo shortcodes, the hugo-content-authoring skill activates and provides Claude with detailed knowledge about shortcode syntax, patterns, and best practices.

### Agents: autonomous work

Agents are active components that perform multi-step tasks using assigned tools. They analyze, plan, and execute, reading files, running commands, and making changes. Each agent has a specific model, a set of tools, and a defined process for how it approaches its domain.

| Agent | Role | Tools |
|-------|------|-------|
| hugo-site-architect | Site design and scaffolding | Read, Write, Edit, Bash, Glob, Grep, WebSearch, WebFetch |
| hugo-build-doctor | Build diagnostics and troubleshooting | Read, Bash, Glob, Grep |

The hugo-site-architect has a broader toolset because it creates files and modifies configuration. The hugo-build-doctor has a narrower, read-oriented toolset because its job is to investigate and diagnose, not to make changes without user approval.

### Commands: quick actions

Commands are user-invoked actions that follow a defined procedure. They are deterministic in their process -- each command has a fixed sequence of steps that it follows when invoked. Commands are the primary interface for common operations.

| Command | Action |
|---------|--------|
| `/hugo-init` | Scaffold a complete Hugo site from repository analysis |
| `/hugo-serve` | Start the Hugo development server |
| `/hugo-deploy` | Set up deployment to GitHub Pages or S3 |
| `/hugo-add-section` | Mount a new directory as a content section |

Commands call agents and use templates to accomplish their work. For example, `/hugo-init` leverages the hugo-site-architect agent's analysis capabilities and uses `hugo.toml.tmpl` to generate configuration.

### Templates: scaffolding

Templates are static files with placeholders that commands use to generate project files. They encode the plugin's recommended configuration patterns and best practices.

| Template | Generates |
|----------|-----------|
| `hugo.toml.tmpl` | Hugo site configuration with module mounts |
| `github-pages.yml.tmpl` | GitHub Actions workflow for GitHub Pages deployment |
| `s3-deploy.yml.tmpl` | GitHub Actions workflow for S3 deployment |
| `_index.md.tmpl` | Section index pages for content sections |

Templates are not used directly by users. Commands fill in their placeholders with project-specific values during scaffolding.

## Component flow

The typical workflow through the plugin's components follows this pattern:

1. A **command** initiates the action (user runs `/hugo-init`)
2. An **agent** performs analysis and planning (hugo-site-architect scans the repository)
3. **Templates** provide the scaffolding structure (`hugo.toml.tmpl` defines the configuration shape)
4. **Skills** inform the agent's decisions (hugo-module-mounts knowledge ensures correct mount configuration)

Not every interaction uses all four types. A user asking "how do Hugo shortcodes work?" only activates the hugo-content-authoring skill. A user reporting a build failure activates the hugo-build-doctor agent directly. The four types are layers available as needed, not a required pipeline.

## Component inventory

The complete set of 17 components:

| # | Type | Name |
|---|------|------|
| 1 | Skill | hugo-fundamentals |
| 2 | Skill | hugo-content-authoring |
| 3 | Skill | hugo-data-templates |
| 4 | Skill | hugo-themes |
| 5 | Skill | hugo-module-mounts |
| 6 | Skill | hugo-github-actions |
| 7 | Skill | hugo-s3-deployment |
| 8 | Agent | hugo-site-architect |
| 9 | Agent | hugo-build-doctor |
| 10 | Command | /hugo-init |
| 11 | Command | /hugo-serve |
| 12 | Command | /hugo-deploy |
| 13 | Command | /hugo-add-section |
| 14 | Template | hugo.toml.tmpl |
| 15 | Template | github-pages.yml.tmpl |
| 16 | Template | s3-deploy.yml.tmpl |
| 17 | Template | _index.md.tmpl |

## Related

- {{< ref "explanation/design-decisions" >}} -- Why the plugin is structured this way
- {{< ref "explanation/component-interaction" >}} -- Detailed interaction patterns between components
- {{< ref "reference/skills" >}} -- Full specifications for all 7 skills
- {{< ref "reference/agents" >}} -- Full specifications for both agents
- {{< ref "reference/commands" >}} -- Full specifications for all 4 commands
- {{< ref "reference/templates" >}} -- Full specifications for all 4 templates
