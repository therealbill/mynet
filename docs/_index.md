---
title: "Mynet Plugin Marketplace"
description: "Documentation for the Mynet Claude Code plugin marketplace"
weight: 1
type: docs
---

# Mynet Plugin Marketplace

A curated collection of Claude Code plugins for software engineering, research, documentation, and development operations.

Mynet organizes specialized Claude agents, skills, commands, and templates into domain-focused plugins that can be installed independently or combined for comprehensive coverage.

## Quick Start

### Installation

```bash
# Add the Mynet marketplace (one-time setup)
/plugin marketplace add therealbill/mynet

# Install a plugin from the marketplace
/plugin install code-quality@mynet --scope project

# Or install from a local directory
claude plugin install ./code-quality
```

### Using a Plugin

Once installed, plugin components activate automatically:

- **Agents** are dispatched when your question matches their expertise
- **Skills** trigger when you ask about topics they cover
- **Commands** are invoked with `/command-name`
- **Templates** provide starter files for common patterns

## Documentation

### Tutorials

Step-by-step guides for getting started:

- [Your First Plugin](tutorials/getting-started/) — Install a plugin and use its components

### How-To Guides

Task-oriented guides:

- [Install a Plugin](howto/install-a-plugin/)
- [Find the Right Plugin](howto/find-right-plugin/)
- [Use Plugins Together](howto/use-plugins-together/)

### Reference

Technical specifications:

- [Plugin Catalog](reference/plugin-catalog/) — All available plugins
- [Marketplace Manifest](reference/marketplace-manifest/) — marketplace.json format
- [Plugin JSON](reference/plugin-json/) — plugin.json format
- [Component Conventions](reference/component-conventions/) — Agent, skill, command, template formats

### Explanation

Conceptual documentation:

- [Marketplace Architecture](explanation/marketplace-architecture/) — How the marketplace works
- [Plugin Design Philosophy](explanation/plugin-design-philosophy/) — Why this structure
- [Licensing](explanation/licensing/) — License model

## Plugins

| Plugin | Description | Components |
|--------|-------------|------------|
| [code-quality](plugins/code-quality/) | Code review, testing, accessibility, architectural quality | 5 agents |
| [web-development](plugins/web-development/) | Frontend frameworks, UI/UX, fullstack web | 6 agents |
| [backend-development](plugins/backend-development/) | Backend architecture, Go services, SQL | 3 agents |
| [mcp-development](plugins/mcp-development/) | MCP server integration and registry | 2 agents |
| [programming-languages](plugins/programming-languages/) | Language-specific experts (C++, C#, JS, TS, Rust, Go, shell, zsh) | 8 agents |
| [mobile-development](plugins/mobile-development/) | iOS, Swift, cross-platform mobile | 3 agents |
| [game-development](plugins/game-development/) | Unity, Unreal Engine, game design | 5 agents |
| [ai-development](plugins/ai-development/) | AI/ML engineering, agent modernization | 1 agent, 1 skill |
| [desktop-development](plugins/desktop-development/) | Electron+Go hybrid desktop apps | 1 agent |
| [cli-development](plugins/cli-development/) | CLI tools, TUI frameworks, terminal UI | 3 agents |
| [devops-and-infra](plugins/devops-and-infra/) | DevOps automation, CI/CD, monitoring | 4 agents |
| [developer-tools](plugins/developer-tools/) | Git workflows, documentation, prototyping | 3 agents |
| [diataxis-docs](plugins/diataxis-docs/) | Diataxis documentation framework | 7 agents |
| [utility-agents](plugins/utility-agents/) | URL analysis, workflow orchestration | 3 agents |
| [research](plugins/research/) | Academic, technical, multi-source research | 5 agents |
| [utility-skills](plugins/utility-skills/) | Markdown formatting, cross-domain utilities | 1 skill |
| [gnu-make](plugins/gnu-make/) | GNU Make best practices | 5 skills |
| [timelord](plugins/timelord/) | Temporal.io expertise | 3 agents, 16 skills, 6 commands |
| [hugo-repo](plugins/hugo-repo/) | Hugo site management for repositories | 2 agents, 7 skills, 4 commands |
