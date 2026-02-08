---
title: "Marketplace Architecture"
description: "How the Mynet marketplace discovers, loads, and organizes plugins"
weight: 1
---

# Marketplace Architecture

The Mynet marketplace solves a specific problem: how does Claude Code find and load the right specialist for a given task, when there are 64 agents, 30 skills, and 10 commands spread across 19 plugins? The answer is a two-layer discovery system -- a manifest that lists plugins, and a convention-based directory structure that exposes components within each plugin.

## The discovery problem

Claude Code needs to know what capabilities are available before a user ever asks a question. Without a discovery mechanism, every plugin would need to register itself through some initialization step, creating ordering dependencies and startup complexity. The marketplace sidesteps this entirely with a static, declarative approach.

## How discovery works

Discovery happens in two stages: plugin enumeration, then component scanning.

```
marketplace.json          Plugin directories         Component directories
+-----------------+       +-------------------+      +------------------+
| plugins:        |       | backend-dev/      |      | agents/          |
|   - backend-dev-+------>|   .claude-plugin/  |      |   sql-pro.md     |
|   - gnu-make    |       |   plugin.json     +----->|   go-architect.md|
|   - timelord    |       |   agents/         |      | skills/          |
|   ...           |       |   skills/         |      | commands/        |
+-----------------+       +-------------------+      +------------------+
```

**Stage 1: Manifest enumeration.** The root `.claude-plugin/marketplace.json` file lists every plugin by name and source path. Claude Code reads this single file to learn that 19 plugins exist and where each one lives. The manifest also carries descriptions and keywords, which provide a first-pass filter -- but the real dispatch happens at the component level.

**Stage 2: Convention-based scanning.** Within each plugin directory, Claude Code looks for well-known subdirectories: `agents/`, `skills/`, `commands/`, and `templates/`. Any `.md` file in `agents/` is an agent definition. Any `SKILL.md` file inside a subdirectory of `skills/` is a skill definition. Any `.md` file in `commands/` is a slash command. This convention-over-configuration approach means no registration code is needed -- the file system *is* the registry.

## How dispatch works

Once all components are loaded, Claude Code uses their `description` fields to match user intent.

For agents, the description field explains what the agent does and when it should be used. When a user asks "optimize this SQL query," Claude reads all agent descriptions and selects `sql-pro` because its description mentions query optimization, execution plans, and CTEs. The description field is not documentation -- it is the dispatch mechanism. A poorly written description means the agent never gets selected.

For skills, the description field lists trigger conditions. The `makefile-fundamentals` skill fires when a user asks to "create a Makefile" or "fix a missing separator error" because those phrases appear in its description. Skills inject knowledge into whichever agent is currently active, rather than taking over the conversation.

Commands bypass the matching system entirely. They register as `/command-name` and activate only through explicit user invocation.

## Why a flat structure

All 19 plugins sit at the same directory level. There is no nesting, no plugin categories, and no dependency chains between plugins. This design reflects a deliberate choice.

A hierarchical structure -- say, grouping `backend-development`, `programming-languages`, and `cli-development` under a "backend" parent -- would impose an organizational opinion that may not match how users think. A user writing a Go CLI tool needs components from at least three plugins. Hierarchy would either force them to navigate a tree or require complex cross-category references.

The flat structure also means any plugin can be installed independently. Removing the `game-development` plugin has zero effect on `web-development`. Plugins collaborate not through import chains but through Claude's dispatch layer -- a `code-reviewer` agent from `code-quality` can review code that a `typescript-pro` agent from `programming-languages` just wrote, without either plugin knowing the other exists.

## Independent versioning

Each plugin maintains its own version in its `plugin.json` file. The marketplace manifest optionally mirrors the version, but the plugin's own manifest is authoritative. This separation exists because plugins evolve at different rates. The `timelord` plugin is at version 1.1.0 while most others are at 1.0.0. Tying all plugins to a single marketplace version would either slow down active plugins or force meaningless version bumps on stable ones.

The cost of independent versioning is that there is no single "marketplace version" that guarantees a known-good combination of plugins. In practice, this has not been a problem because plugins do not depend on each other's internals.

## Related

- [Plugin Catalog]({{< ref "reference/plugin-catalog" >}}) -- complete listing of all plugins and their components
- [Marketplace Manifest]({{< ref "reference/marketplace-manifest" >}}) -- field-by-field reference for marketplace.json
