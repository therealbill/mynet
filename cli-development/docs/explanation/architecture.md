---
title: "Architecture"
description: "CLI vs TUI distinction and when to use each agent"
weight: 1
---

# Architecture

The cli-development plugin contains three agents that cover different aspects of terminal application development. This document explains the distinctions between them, why the responsibilities are split the way they are, and how they combine in practice.

## CLI vs TUI

The fundamental distinction in this plugin is between CLI and TUI:

- **CLI (Command-Line Interface)** — The user provides all input upfront as arguments and flags, the tool processes that input, produces output, and exits. The interaction is non-interactive: input goes in, output comes out, the process terminates. CLIs compose well in pipelines (`grep | sort | uniq`), are scriptable, and work over SSH without special terminal support.

- **TUI (Terminal User Interface)** — The user interacts with a stateful, visual application that takes over the terminal. The application responds to keypresses in real time, redraws the screen, and maintains state between interactions. TUIs are interactive sessions where the user makes decisions during execution — browsing a list, selecting items, filling forms, watching live data.

A single tool can have both. A database migration tool might expose `migrate status` as a CLI command that prints a table, and `migrate status --interactive` as a TUI that lets the user browse and inspect migrations. The CLI path serves scripting and automation; the TUI path serves human exploration.

## Why Three Agents

Each agent addresses a distinct concern:

**cli-developer** handles the structural shell of a command-line tool. This includes the command hierarchy (top-level commands, subcommands, flags), configuration precedence (flags over env vars over config files over defaults), output format conventions (human-readable by default, structured with `--output json`), error handling (stderr, exit codes), and shell completions. The agent produces the scaffolding that makes a CLI tool feel professional and predictable. It works across Go (Cobra), Node.js (Commander/yargs), and Python (Click).

**go-tui-developer** handles interactive terminal interfaces built with the Charm stack in Go. This is a specialized domain — the Bubble Tea framework uses a Model-View-Update pattern that requires careful model decomposition, message routing, and component composition. The agent understands Bubbles components (list, table, viewport, textinput, spinner, progress), Lip Gloss styling, Huh forms, terminal theming with termenv, and color profile degradation. It uses the opus model because complex TUI development involves significant state management and multi-component orchestration.

**cli-ui-designer** handles visual design decisions for terminal output. This agent does not write implementation code. It produces design specifications — color role maps, prompt symbol vocabularies, typographic hierarchies, whitespace structures — that guide implementation. The distinction is important: design decisions (which color means "error", how to indicate progress, what symbol means "run this") are separate from the code that renders them.

## The Design-Implementation Split

cli-ui-designer and go-tui-developer represent a deliberate separation of design and implementation.

cli-ui-designer answers questions like: What color palette communicates the right tone? What does the prompt vocabulary look like? How does the visual hierarchy guide the eye? These are design decisions that apply regardless of the implementation technology. A color role map designed by cli-ui-designer works whether the output is rendered with Lip Gloss, ANSI escape codes in Python, or chalk in Node.js.

go-tui-developer answers questions like: How should the Bubble Tea model hierarchy work? Which Bubbles components compose to create this interaction? How does the theme file get loaded and applied to Lip Gloss styles? These are implementation decisions specific to the Charm stack in Go.

In practice, a workflow might start with cli-ui-designer to define the visual language, then move to go-tui-developer to implement it. The design agent produces specifications; the implementation agent produces code.

## When Agents Combine

A typical terminal application might involve all three agents at different stages:

1. **cli-developer** designs the command tree and flag conventions. The result is a Cobra-based CLI with proper help text, shell completions, and configuration handling.

2. **cli-ui-designer** defines the visual language — color roles for status indicators, prompt symbols for different action types, whitespace rules for output grouping.

3. **go-tui-developer** implements interactive components — a Bubble Tea model for browsing data, Lip Gloss styles derived from the design spec, theme support for user customization.

The shared package layout keeps these concerns organized:

```
cmd/                      # Cobra commands (cli-developer)
internal/tui/             # Bubble Tea models (go-tui-developer)
internal/tui/styles/      # Lip Gloss styles from design spec (cli-ui-designer + go-tui-developer)
internal/app/             # Business logic (shared)
```

Not every project needs all three. A non-interactive CLI tool might use only cli-developer. A data dashboard might skip cli-developer entirely and go straight to go-tui-developer. A project that already has code but needs visual refinement might use only cli-ui-designer.

## Model Selection

go-tui-developer uses the **opus** model. TUI development involves complex state management — multiple concurrent models, message routing between components, window resize propagation through hierarchies, and component lifecycle management. The higher capability model handles this orchestration more reliably.

cli-developer and cli-ui-designer use the **sonnet** model. CLI structure and visual design are focused domains where the agent applies well-understood patterns (Cobra conventions, GNU flag style, ANSI color semantics) without the multi-layered state management that TUI development demands.

## Cross-Plugin Relationships

The cli-development plugin intersects with several other plugins:

- **programming-languages** — For Go code quality, idiomatic patterns, and language-specific best practices beyond CLI/TUI concerns. cli-developer and go-tui-developer produce Go code, but the programming-languages plugin covers broader Go expertise.

- **backend-development/go-architect** — When a CLI tool is part of a larger Go service. The go-architect agent handles service architecture (HTTP servers, database layers, dependency injection) while cli-developer handles the command-line entry points.

- **desktop-development** — When a terminal interface is not enough. If requirements include rich graphics, mouse interaction beyond basic click support, or multi-window layouts, the desktop-development plugin covers GUI frameworks. The transition point: if the terminal is the right environment, use cli-development; if you need a desktop application, use desktop-development.
