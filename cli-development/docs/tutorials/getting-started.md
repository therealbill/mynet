---
title: "Getting Started with CLI Development"
description: "Build a CLI command with cli-developer and add TUI components with go-tui-developer"
weight: 1
---

# Getting Started with CLI Development

This tutorial walks you through building a CLI tool and adding interactive terminal UI components. You will use two agents — **cli-developer** for the command structure and **go-tui-developer** for interactive elements — and see how they work together in a typical development workflow.

By the end, you will have a CLI tool with a well-designed command tree, proper flag conventions, and an interactive selection list powered by Bubble Tea.

## Prerequisites

- A Go project initialized with `go mod init`, or willingness to start one
- Basic familiarity with Go (the agents handle framework specifics)
- A terminal emulator that supports at least ANSI 256 colors (most modern terminals do)

## Step 1: Describe Your CLI Tool Concept

Start by telling **cli-developer** what you want to build. Be specific about what the tool manages and who uses it.

Example prompt:

> "Build a CLI for managing database migrations. It should support creating new migrations, running them, rolling back, and showing migration status."

The agent will ask clarifying questions if needed — what database engines to support, whether the tool is standalone or embedded in a larger project, and what output formats matter. Answer these to shape the design.

## Step 2: Review the Command Tree Design

cli-developer produces a command tree — the hierarchy of commands, subcommands, and flags that define how users interact with your tool. For the migration example, this might look like:

```
migrate
├── create <name>          --dir ./migrations
├── up                     --steps N  --dry-run
├── down                   --steps N  --dry-run
├── status                 --output json|table
└── version
```

Review this structure for:

- **Progressive disclosure** — Common operations are top-level commands. Rarely-used options are flags, not separate commands.
- **GNU-style flags** — Long flags with double dashes (`--dry-run`), short aliases where intuitive (`-n` for `--steps`).
- **Predictable patterns** — If `up` takes `--steps`, `down` takes `--steps` too. Consistency across subcommands.

The agent designs the command tree using Cobra as the default Go CLI framework. It also defines the configuration precedence: flags override environment variables, which override config file values, which override built-in defaults.

## Step 3: See the Thin Wiring Pattern

cli-developer follows a strict separation: command handlers parse input and call library code. Business logic never lives in the handler.

The generated code looks like this:

```go
// cmd/up.go — thin wiring, no business logic
var upCmd = &cobra.Command{
    Use:   "up",
    Short: "Apply pending migrations",
    RunE: func(cmd *cobra.Command, args []string) error {
        steps, _ := cmd.Flags().GetInt("steps")
        dryRun, _ := cmd.Flags().GetBool("dry-run")

        // Delegate to library code
        return app.MigrateUp(cfg.DB, steps, dryRun)
    },
}
```

```go
// internal/app/migrate.go — business logic, testable independently
func MigrateUp(db *sql.DB, steps int, dryRun bool) error {
    // All migration logic here
}
```

This pattern keeps handlers testable, reusable, and easy to refactor. The agent also generates shell completions for bash, zsh, fish, and PowerShell — these are not optional.

## Step 4: Add a TUI Component with go-tui-developer

Now switch to **go-tui-developer** to add an interactive element. For example, you want `migrate status` to show an interactive list of migrations that users can browse and inspect.

Example prompt:

> "Add interactive selection prompt to the migrate status command. Users should see a list of migrations, scroll through them, and press enter to see details."

go-tui-developer works with the Charm stack — Bubble Tea for the application framework, Bubbles for pre-built components, and Lip Gloss for styling.

## Step 5: Review the Bubble Tea Model

go-tui-developer produces code following the Model-View-Update (MVU) pattern:

- **Init** — Returns the initial command (e.g., fetch migration list from the database)
- **Update** — Handles messages (key presses, window resize, data loaded) and returns the updated model
- **View** — Renders the current model state as a string for the terminal

```go
type statusModel struct {
    list     list.Model    // Bubbles list component
    selected *Migration
    width    int
    height   int
}

func (m statusModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
        m.list.SetSize(msg.Width, msg.Height-2)
    case tea.KeyMsg:
        if msg.String() == "enter" {
            m.selected = m.list.SelectedItem().(*Migration)
        }
    }
    // ...
}
```

The agent integrates this with Cobra so that `migrate status --interactive` launches the TUI, while `migrate status` still outputs a plain table or JSON. This preserves the non-interactive path for scripting and pipelines.

The package layout follows a clean structure:

```
cmd/                    # Cobra command definitions
internal/tui/           # Bubble Tea models and components
internal/tui/styles/    # Lip Gloss style definitions
internal/app/           # Business logic (shared by CLI and TUI)
main.go
```

## What You Learned

- **cli-developer** designs the command tree, flag conventions, configuration precedence, and output formats. It produces thin command handlers that delegate to library code.
- **go-tui-developer** adds interactive terminal UI components using the Bubble Tea MVU pattern, with proper window resize handling and Cobra integration.
- The two agents share a common package layout. Business logic in `internal/app/` is used by both CLI handlers and TUI models.

## Next Steps

- [Build an Interactive TUI]({{< ref "/cli-development/howto/build-interactive-tui" >}}) — Detailed guide for Bubble Tea TUI development
- [Design CLI Visual Style]({{< ref "/cli-development/howto/design-cli-visual-style" >}}) — Add visual polish with cli-ui-designer
- [Agent Reference]({{< ref "/cli-development/reference/agents" >}}) — Full specifications for all three agents

## Cross-Plugin References

- **programming-languages** plugin — For Go code quality, idioms, and best practices beyond CLI/TUI specifics
- **backend-development/go-architect** agent — When your CLI tool is part of a larger Go service architecture
