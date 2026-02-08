---
title: "Build an Interactive TUI"
description: "Use go-tui-developer to create Bubble Tea terminal interfaces"
weight: 1
---

# Build an Interactive TUI

This guide walks through adding an interactive terminal user interface to a Go application using **go-tui-developer**.

## Goal

Add an interactive terminal UI — browsing data, selecting items, displaying live status — to a Go application using the Charm stack (Bubble Tea, Bubbles, Lip Gloss).

## Step 1: Describe the Interaction

Tell go-tui-developer what the user needs to do interactively. Focus on the experience, not the implementation.

Good trigger phrases:

- "Build a TUI for browsing API responses"
- "Add interactive selection prompt"
- "Set up CLI with Cobra and make output look good"
- "Add theme support with dark mode detection"

Be specific about data flow — where does the data come from, what can the user do with it, and what happens when they make a selection.

## Step 2: Trigger go-tui-developer

The agent activates on prompts involving interactive terminal interfaces, Bubble Tea, Bubbles components, Lip Gloss styling, or TUI theming. Describe your use case and the agent will design the component architecture.

## Step 3: Review the Model Hierarchy

go-tui-developer designs a hierarchy of Bubble Tea models. Each model follows the MVU (Model-View-Update) pattern with three methods:

- **Init** — Returns the initial command to run (data fetching, timer setup)
- **Update** — Receives messages, updates state, returns commands
- **View** — Renders the current state as a string

The agent composes models using Bubbles components:

- **list.Model** — Scrollable, filterable lists
- **table.Model** — Tabular data with column sorting
- **viewport.Model** — Scrollable text content
- **textinput.Model** — Single-line text input
- **spinner.Model** — Loading indicators
- **progress.Model** — Progress bars

Review how messages flow between parent and child models. Parent models receive messages first and delegate to the active child. This is where most complexity lives.

## Step 4: Review the Package Layout

The generated code follows a standard structure:

```
cmd/                      # Cobra command definitions
internal/tui/             # Bubble Tea models and update logic
internal/tui/styles/      # Lip Gloss style definitions, theme loading
internal/app/             # Business logic (database, API, file I/O)
main.go
```

Business logic stays in `internal/app/`. TUI models in `internal/tui/` call into `internal/app/` for data and operations but never contain domain logic themselves. This separation keeps models testable and allows the same logic to serve both interactive and non-interactive paths.

## Step 5: Add Theming

If your TUI needs theming support, go-tui-developer can add:

- **Theme files** — User-defined themes in YAML or TOML with semantic color roles (primary, secondary, success, warning, error, muted, background, foreground)
- **Dark mode detection** — `termenv.HasDarkBackground()` detects the terminal background
- **Adaptive colors** — `lipgloss.AdaptiveColor` automatically switches between light and dark variants
- **Color profile degradation** — Styles degrade gracefully from TrueColor to ANSI 256 to ANSI 16 to plain ASCII

Theme loading happens at startup. Styles are derived from the loaded theme, not hardcoded.

## Verification

After the agent generates the TUI code, verify:

- **Window resizing** — Resize the terminal while the TUI is running. Every component should adjust. The agent handles `tea.WindowSizeMsg` and propagates dimensions to child models.
- **80-column width** — Run the TUI at 80 columns wide. All content should be readable without horizontal truncation breaking meaning.
- **120-column width** — Run at 120 columns. The layout should use the extra space without becoming sparse or unbalanced.
- **Color degradation** — Test with `COLORTERM=` unset and `TERM=xterm` to simulate limited color support. The interface should remain usable.

## Troubleshooting

**Layout breaks on terminal resize**
The most common issue is missing `tea.WindowSizeMsg` handling. Every model that renders layout-sensitive content must respond to this message and propagate new dimensions to its child components. go-tui-developer handles this by default, but if you modify the code manually, ensure every model in the hierarchy processes resize messages.

**When to use cli-developer instead**
If the interaction is non-interactive — the user provides all input as flags and arguments, the tool produces output, and exits — use **cli-developer** instead. TUIs are for stateful, interactive sessions where the user makes decisions during execution. A tool that prints a report does not need Bubble Tea.

**Huh forms for structured input**
For collecting structured input (multi-field forms, confirmations, selections), go-tui-developer uses the Huh library rather than building custom input flows. Mention "form" or "survey" in your prompt to trigger this.
