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
| [code-quality](../code-quality/docs/) | Code review, testing, accessibility, architectural quality | 5 agents |
| [web-development](../web-development/docs/) | Frontend frameworks, UI/UX, fullstack web | 6 agents |
| [backend-development](../backend-development/docs/) | Backend architecture, Go services, SQL | 3 agents |
| [mcp-development](../mcp-development/docs/) | MCP server integration and registry | 2 agents |
| [programming-languages](../programming-languages/docs/) | Language-specific experts (C++, C#, JS, TS, Rust, Go, shell, zsh) | 8 agents |
| [mobile-development](../mobile-development/docs/) | iOS, Swift, cross-platform mobile | 3 agents |
| [game-development](../game-development/docs/) | Unity, Unreal Engine, game design | 5 agents |
| [ai-development](../ai-development/docs/) | AI/ML engineering, agent modernization | 1 agent, 1 skill |
| [desktop-development](../desktop-development/docs/) | Electron+Go hybrid desktop apps | 1 agent |
| [cli-development](../cli-development/docs/) | CLI tools, TUI frameworks, terminal UI | 3 agents |
| [devops-and-infra](../devops-and-infra/docs/) | DevOps automation, CI/CD, monitoring | 4 agents |
| [developer-tools](../developer-tools/docs/) | Git workflows, documentation, prototyping | 3 agents |
| [diataxis-docs](../diataxis-docs/docs/) | Diataxis documentation framework | 7 agents |
| [utility-agents](../utility-agents/docs/) | URL analysis, workflow orchestration | 3 agents |
| [research](../research/docs/) | Academic, technical, multi-source research | 5 agents |
| [utility-skills](../utility-skills/docs/) | Markdown formatting, cross-domain utilities | 1 skill |
| [gnu-make](../gnu-make/docs/) | GNU Make best practices | 5 skills |
| [timelord](../timelord/docs/) | Temporal.io expertise | 3 agents, 16 skills, 6 commands |
| [hugo-repo](../hugo-repo/docs/) | Hugo site management for repositories | 2 agents, 7 skills, 4 commands |
