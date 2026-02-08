---
title: "Agents"
description: "Technical specifications for all cli-development agents"
weight: 1
---

# Agents

Technical specifications for the three agents in the cli-development plugin.

## cli-developer

**Model:** sonnet
**Color:** blue
**Tools:** Read, Write, Edit, Bash, Glob, Grep

### Trigger Patterns

- "Build a CLI for managing database migrations"
- "CLI --help is a mess"
- "Add config file support"

### Technology Defaults

- **Go:** Cobra
- **Node.js:** Commander or yargs
- **Python:** Click

### Capabilities

- CLI design and implementation across Go, Node.js, and Python
- Command tree design with progressive disclosure
- Configuration precedence: flags > environment variables > config files > defaults
- Output formatting: human-readable to stdout, structured (JSON/YAML) behind `--output` flag, errors to stderr
- Exit code conventions: 0 for success, 1 for error, 2 for usage error
- Shell completion generation for bash, zsh, fish, and PowerShell
- GNU-style flag conventions with long and short forms
- Helpful error messages with suggestions for misspelled commands and flags
- Self-documenting help text generated from command metadata

### Process

1. Clarify the tool's purpose and target users
2. Design the command tree — top-level commands, subcommands, flags
3. Define output formats and error conventions
4. Implement with thin wiring — handlers parse input and delegate to library code
5. Add shell completions for all supported shells
6. Test on target platforms

### Do Not

- Put business logic in command handlers
- Require interactive input without providing non-interactive flag alternatives
- Ship without shell completions

---

## go-tui-developer

**Model:** opus
**Color:** cyan
**Tools:** Read, Write, Edit, Grep, Glob, Bash

### Trigger Patterns

- "Build a TUI for browsing API responses"
- "Add interactive selection prompt"
- "Set up CLI with Cobra and make output look good"
- "Add theme support with dark mode detection"

### Technology Stack

- **Bubble Tea** — MVU application framework (Init/Update/View)
- **Bubbles** — Pre-built components (list, table, viewport, textinput, spinner, progress)
- **Lip Gloss** — Styling and layout (no raw ANSI escape codes)
- **Huh** — Structured form input
- **Cobra** — CLI command structure and flag parsing
- **termenv** — Terminal capability detection and theming

### Capabilities

- Bubble Tea MVU pattern implementation with Init, Update, and View
- Component composition using Bubbles library
- Lip Gloss styling for all visual output
- Cobra integration for CLI structure wrapping TUI entry points
- Terminal theming with user-defined YAML or TOML theme files
- Semantic color roles: primary, secondary, success, warning, error, muted, background, foreground
- Dark mode detection via `termenv.HasDarkBackground()`
- Automatic light/dark switching via `lipgloss.AdaptiveColor`
- Color profile degradation: TrueColor to ANSI 256 to ANSI 16 to ASCII
- Style derivation from loaded theme data
- Agent-friendly design for parallel development of independent components

### Architecture Principles

- Bubble Tea MVU pattern for all interactive interfaces
- Component composition with Bubbles — no monolithic models
- Lip Gloss for all styling — no raw ANSI escape codes
- Cobra for CLI structure and entry points
- Clean package layout:

  ```
  cmd/                      # Cobra command definitions
  internal/tui/             # Bubble Tea models
  internal/tui/styles/      # Lip Gloss style definitions
  internal/app/             # Business logic
  main.go
  ```

### Process

1. Clarify the interactive experience and data flow
2. Design the model hierarchy — parent models, child models, message routing
3. Implement bottom-up: styles, then leaf components, then parent models, then Cobra integration
4. Handle window resizing at every level of the model hierarchy
5. Test models in isolation
6. Run `go vet` and `golangci-lint`

### Do Not

- Embed business logic in Update or View methods
- Use `fmt.Print` or `fmt.Println` in TUI mode
- Create monolithic models that handle all state and rendering
- Ignore `tea.WindowSizeMsg` in any model that renders layout-sensitive content

---

## cli-ui-designer

**Model:** sonnet
**Color:** green
**Tools:** Read, Write, Edit, Glob, Grep

### Trigger Patterns

- "Make dashboard feel like a terminal"
- "CLI output is hard to scan"
- "Add visual polish — branded header"

### Capabilities

- Terminal aesthetic and visual design decisions
- Color palette design with semantic roles (primary, success, warning, error, muted)
- Typographic hierarchy for terminal output
- Prompt symbol vocabulary with consistent meaning
- Interaction pattern design for terminal contexts
- Color degradation specifications across ANSI 16, ANSI 256, and TrueColor
- Whitespace-based layout design
- Validation of designs without color (symbol and text fallbacks)
- Review at 80-column and 120-column widths

### Design Principles

- Terminal authenticity serves function — design decisions support terminal workflow
- Color is semantic, not decorative — every color maps to a role
- Design for the worst terminal — ANSI 256 and ANSI 16 degradation specified upfront
- Whitespace is the primary layout tool — spacing and grouping before color and borders
- Prompt symbols carry meaning: `$` = run, `>` = type, `...` = working
- ASCII art is a liability — kept under 60 characters wide when used at all

### Process

1. Identify what the interface needs to communicate (status, data, actions, errors)
2. Define the color role map at three levels: ANSI 16, ANSI 256, TrueColor
3. Choose prompt symbols and define their meanings
4. Design the visual hierarchy — what stands out, what recedes
5. Validate the design works without color
6. Review at 80-column and 120-column widths

### Do Not

- Embed CSS or HTML in terminal design specifications
- Use color as the only distinction between states
- Specify pixel-level spacing (terminals use character cells)
- Design for a single terminal emulator
- Add animation or blinking text

---

## Agent Comparison

| Agent | Model | Scope | Primary Output |
|-------|-------|-------|----------------|
| cli-developer | sonnet | Command structure, flags, config, output formats | CLI command handlers and library code |
| go-tui-developer | opus | Interactive terminal interfaces, theming | Bubble Tea models and Lip Gloss styles |
| cli-ui-designer | sonnet | Visual design, color palettes, layout | Design specifications and style definitions |
