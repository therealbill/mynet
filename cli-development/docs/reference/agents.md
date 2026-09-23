---
title: "Agents"
description: "Technical specifications for all cli-development agents"
weight: 1
---

# Agents

Technical specifications for the three agents in the cli-development plugin. For the reasoning behind how scope is split between them, see [Architecture](../../explanation/architecture/).

## cli-developer

**Model:** sonnet
**Color:** blue
**Tools:** Read, Write, Edit, Bash, Glob, Grep

### Scope

Go command-line tools: Cobra command trees, flags and subcommands, migrating hand-rolled `os.Args` or flag-package parsing onto Cobra, auditing an existing CLI's usability, configuration layering, output formats, exit codes, error messages, and shell completions. Node and Python CLIs on request. Hands interactive screens, and any Go project containing a Bubble Tea screen, to go-tui-developer; hands color, symbol, and visual hierarchy decisions to cli-ui-designer.

### Trigger Patterns

- "Build a CLI for managing database migrations"
- "Our CLI's --help output is a mess"
- "Add support for a config file"
- Proactive: a Go project parsing `os.Args` by hand when a new flag is requested

### Technology Defaults

- **Go:** Cobra, or urfave/cli when already in use
- **Node.js:** Commander or yargs
- **Python:** Click

### Conventions

- Configuration precedence: flags > environment variables > config file > compiled defaults; config under `$XDG_CONFIG_HOME/<tool>/`, falling back to `~/.config/<tool>/` when the variable is unset, empty, or not an absolute path; never `os.UserConfigDir()`, which returns a non-XDG path on macOS
- Data on stdout, diagnostics on stderr; structured output behind `--output json|yaml`
- Exit codes: 0 success, 1 error, 2 usage error, 130 interrupted; 141 is shell-reported for a SIGPIPE kill, never set by the tool itself
- Color only on a TTY with `NO_COLOR` unset or empty; `CLICOLOR_FORCE` or a nonzero `FORCE_COLOR` overrides and wins over `NO_COLOR`, `FORCE_COLOR=0` disables color
- SIGINT cancels through `context.Context`, then re-raises the signal (`signal.Reset` plus `syscall.Kill`) so the exit status is genuinely signaled, falling back to `os.Exit(130)` on Windows; checking stdout writes for SIGPIPE/EPIPE is dead code, since a broken pipe already kills the Go process
- `-` for stdin, `--` for end of options, `--dry-run` on every mutating command
- Helpful errors with the failing value and a "did you mean" suggestion
- Completions for bash, zsh, fish, and PowerShell, with dynamic completions for runtime values
- Non-interactive equivalent for every prompt

### Process

1. Design the command tree around composing workflows, and decide where each flag lives
2. Fix output formats, exit codes, and error conventions before writing handlers
3. Implement thin commands that delegate to library code
4. Add completions and verify `--help` per command
5. Test flag edge cases, piped I/O, and `NO_COLOR` on target platforms

### Output

Compiling code plus a summary of commands, flags, exit codes, config precedence, and open decisions. For reviews, a findings list ordered by user impact, each naming the exact flag, message, or help text to change.

### Do Not

- Put business logic in command handlers
- Require interactive input without a non-interactive equivalent
- Print unstructured text when `--output json` is set
- Emit color or progress animations into a pipe
- Ship without shell completions

---

## go-tui-developer

**Model:** opus
**Color:** cyan
**Tools:** Read, Write, Edit, Grep, Glob, Bash

### Scope

Interactive, full-screen, and inline TUIs in Go on the Charm stack, and any Go project that mixes TUI screens with plain commands, including its Cobra wiring. Plain Go CLIs with no interactive screens belong to cli-developer; choosing the palette or symbol vocabulary itself, even for a Bubble Tea app, belongs to cli-ui-designer — implementing a chosen palette stays here.

See also: [Build an Interactive TUI](../../howto/build-interactive-tui/) for the task-oriented walkthrough.

### Trigger Patterns

- "Build a TUI for browsing API responses"
- "Add an interactive selection prompt"
- "Wire the plain export command and the interactive status screen under one Cobra root"
- "Add theme support with dark mode detection"
- Proactive: a blocking call inside `Update` when a new key binding is requested

### Technology Stack

- **Bubble Tea** — MVU application framework (Init, Update, View)
- **Bubbles** — Components (list, table, viewport, textinput, spinner, progress, paginator)
- **Lip Gloss** — Styling and layout; no raw ANSI escape codes
- **Huh** — Structured form input
- **Cobra** — CLI command structure and flag parsing

### Charm Version Handling

The agent reads `go.mod` and follows the matching major. New projects default to v2 (`charm.land/.../v2`). v1 (`github.com/charmbracelet/...`) remains only when a required component has no v2 release — for example `76creates/stickers`; `bubble-table` and `promptkit` are v2-ready at their unsuffixed import paths.

| Concern | v1 | v2 |
| --- | --- | --- |
| View | `View() string` | `View() tea.View` |
| Alt screen | `tea.WithAltScreen()` | `v.AltScreen = true` |
| Key presses | `tea.KeyMsg` | `tea.KeyPressMsg` |
| Light/dark | `termenv.HasDarkBackground()` | `tea.RequestBackgroundColor` then `tea.BackgroundColorMsg.IsDark()` |
| Adaptive colors | `lipgloss.AdaptiveColor` | `lipgloss.LightDark` |
| Color profile | `termenv.ColorProfile()` | `tea.ColorProfileMsg` |
| Terminal background | `termenv.DefaultOutput().SetBackgroundColor` writes OSC 11 with no reset helper, so the agent writes the OSC 111 reset itself | `v.BackgroundColor`, with the renderer emitting the reset on shutdown |

The reset returns the terminal to its configured default background, not a saved snapshot of the previous color.

### Architecture Principles

- One `tea.Model` per screen or panel, composed by embedding child models
- Bubbles components wrapped in domain models, never reimplemented
- Lip Gloss for all styling, `tea.WindowSizeMsg` handled in every layout model
- Cobra commands stay thin; every TUI command has a non-interactive form
- Package layout:

  ```
  cmd/                      # Cobra command definitions
  internal/tui/             # Bubble Tea models, one file per screen
  internal/tui/styles/      # Lip Gloss style definitions
  internal/app/             # Business logic
  main.go
  ```

### Theming

- User-defined YAML or TOML themes with semantic color roles, validated at load
- Separate light and dark palettes per theme, selected once the background is known
- Styles derived from the theme into a `Styles` struct passed to models; switching rebuilds the struct and propagates it with a custom `tea.Msg`
- Terminal background set via OSC 11; emulator detection first; AppleScript on macOS Terminal.app only as a v1 fallback

### Process

1. Read `go.mod` to determine the Charm major, or choose v2 for a new project
2. Design the model hierarchy: which screens, which components, how messages flow
3. Implement bottom-up: styles, then leaf components, then parent models, then Cobra wiring
4. Test models by constructing them directly and calling `Update` with synthetic messages; for golden-file tests of full-screen output use `github.com/charmbracelet/x/exp/teatest` on v1 or `github.com/charmbracelet/x/exp/teatest/v2` on v2
5. Run `go vet` and `golangci-lint` before delivering

### Output

Compiling, lint-clean Go code in the package layout above, plus a short summary listing each model and its file, the message types that flow between them, the keybindings, and the Charm major targeted and why. When asked for a design rather than code, deliver the model hierarchy and message flow as a list before any code.

### Do Not

- Embed business logic in `Update` or `View`
- Use `fmt.Print` in TUI mode
- Create monolithic models
- Ignore `tea.WindowSizeMsg`
- Mix v1 and v2 imports in one module

---

## cli-ui-designer

**Model:** sonnet
**Color:** green
**Tools:** Read, Write, Edit, Glob, Grep

### Scope

How CLI and TUI output should look: palettes and semantic color roles, colorblind-safe status colors, visual hierarchy for dense output, prompt and status symbol vocabulary, ASCII fallbacks for symbols or box-drawing that render as boxes or misalign, spacing conventions, ASCII branding, and terminal-styled web dashboards. Produces design decisions and specs, not implementation.

See also: [Design CLI Visual Style](../../howto/design-cli-visual-style/) for the task-oriented walkthrough.

### Trigger Patterns

- "Make this dashboard feel like a terminal"
- "The output of our CLI is hard to scan"
- "Add a branded header and status indicators"
- Proactive: output with no visual hierarchy when a new output line is requested

### Design Principles

- Terminal authenticity serves function
- Color is semantic, not decorative; never the sole indicator
- Design for the worst terminal and both backgrounds: a light-background and a dark-background value for every role at TrueColor and ANSI 256, and a palette index at ANSI 16; also specify the no-color rendering the tool uses when `NO_COLOR` is set or stdout is not a terminal
- Whitespace is the primary layout tool
- Prompt symbols carry meaning: `$` run, `>` type, `...` working
- Unicode is not guaranteed: every symbol gets an ASCII fallback pair, both padded to the same field width; the switching rule resolves the effective locale as `LC_ALL`, else `LC_CTYPE`, else `LANG`, and uses Unicode only when its codeset names UTF-8, matched case-insensitively with the hyphen optional; box-drawing glyphs such as `─` are East Asian Ambiguous width and take two cells in CJK locales, breaking alignment; the pointer fallback is `->`
- ASCII art is a liability: branding only, under 60 columns so it survives an 80-column terminal with indentation, with a plain-text fallback

### Process

1. Identify what the interface must communicate and define the color role map as principle 3 specifies
2. Choose prompt and status symbols, and the ASCII fallback for each
3. Design the visual hierarchy with spacing and indentation before reaching for borders
4. Validate that every visual distinction survives without color and without Unicode
5. Review at 80-column and 120-column widths; the design must work at both

### Output

A Markdown design spec: the color role map as principle 3 specifies, including the no-color rendering; the symbol vocabulary with its ASCII fallback and the switching rule; the spacing and hierarchy rules; and an 80-column and a 120-column mock of the primary screen as plain text. Implementation code is out of scope.

### Do Not

- Embed CSS, HTML, or Go in design guidance
- Use color as the only way to distinguish states
- Specify pixel-level spacing for terminal output; work in character cells and line heights, and keep a terminal-styled web dashboard on a character grid even when its CSS uses rem
- Design for a specific terminal emulator; target the intersection of capabilities
- Add animation or blinking effects unless the user explicitly asks

---

## Agent Comparison

| Agent | Model | Scope | Primary Output |
|-------|-------|-------|----------------|
| cli-developer | sonnet | Go CLIs with no interactive screens; Node and Python on request | Command code and a command summary |
| go-tui-developer | opus | Interactive TUIs and hybrid TUI-plus-CLI Go projects, Charm v1 and v2 | Bubble Tea models, Lip Gloss styles, Cobra wiring |
| cli-ui-designer | sonnet | Visual design decisions for terminal output | Markdown design spec with mocks |

### Routing Rules

- A Go CLI with no interactive screen goes to cli-developer even if it uses Cobra.
- Any project with at least one Bubble Tea screen goes to go-tui-developer, including its plain commands.
- Palette and symbol decisions go to cli-ui-designer even when the request names Bubble Tea; its spec is then implemented by one of the other two.
