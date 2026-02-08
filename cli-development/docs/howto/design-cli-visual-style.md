---
title: "Design CLI Visual Style"
description: "Use cli-ui-designer to create terminal aesthetics and visual hierarchies"
weight: 2
---

# Design CLI Visual Style

This guide walks through using **cli-ui-designer** to make CLI output scannable, visually coherent, and functional across terminal environments.

## Goal

Define a visual language for your CLI tool — color roles, prompt symbols, typographic hierarchy, and layout conventions — that works across terminal capabilities and remains functional without color.

## Step 1: Identify Communication Needs

Before triggering the agent, inventory what your interface needs to communicate:

- **Status** — Success, failure, warnings, in-progress
- **Data** — Tables, lists, key-value pairs, structured output
- **Actions** — What the user should run next, available commands
- **Errors** — What went wrong, why, and what to do about it

Write these down. The agent uses this inventory to design a visual hierarchy where the most important information stands out and secondary details recede.

## Step 2: Trigger cli-ui-designer

The agent activates on prompts about terminal aesthetics and visual design. Good trigger phrases:

- "Make dashboard feel like a terminal"
- "CLI output is hard to scan"
- "Add visual polish — branded header"

Describe the specific output you want to improve. Include an example of current output if possible — the agent works best when it can see what needs to change.

## Step 3: Review the Color Role Map

cli-ui-designer defines colors by semantic role, not by name. The standard roles are:

| Role | Purpose | Example usage |
|------|---------|---------------|
| **primary** | Brand identity, key actions | Tool name, primary headings |
| **success** | Completed, passed, healthy | "Migration applied", check marks |
| **warning** | Attention needed, degraded | "3 migrations pending", caution |
| **error** | Failed, broken, blocked | Error messages, failure indicators |
| **muted** | Secondary, de-emphasized | Timestamps, metadata, help text |

The agent specifies each role at three levels: TrueColor (24-bit hex), ANSI 256 (index), and ANSI 16 (named). This ensures the design degrades gracefully across terminal capabilities.

## Step 4: Review Prompt Symbol Vocabulary

Symbols communicate state without color. The agent defines a consistent vocabulary:

| Symbol | Meaning |
|--------|---------|
| `$` | Run this command |
| `>` | Type this input |
| `...` | Working, in progress |
| `*` | Important, attention |
| `-` | List item, detail |

These symbols remain meaningful even when color is stripped. The agent ensures every piece of information that uses color also has a non-color indicator — a symbol, a label, or a textual prefix.

## Step 5: Validate Without Color

The most important check: does the output make sense with color completely removed?

Pipe your CLI output through `cat` or set `NO_COLOR=1` (if your tool respects the [NO_COLOR](https://no-color.org/) convention). Every status, every distinction, every grouping should still be clear from symbols, labels, whitespace, and text alone.

The agent designs with this constraint from the start. Color enhances meaning but never creates it.

## Verification

- **ANSI 16** — Set `TERM=xterm` and verify the color palette is readable on both dark and light backgrounds
- **ANSI 256** — Test with a 256-color terminal profile. The intermediate palette should be distinct from the 16-color fallback.
- **TrueColor** — Test with `COLORTERM=truecolor`. The full palette should display as designed.
- **80-column width** — All output fits without wrapping that breaks readability. Tables truncate gracefully or switch to vertical layout.
- **120-column width** — Extra width is used for alignment and breathing room, not wasted.

## Troubleshooting

**Colors unreadable on light terminals**
This happens when the design only targets dark backgrounds. cli-ui-designer specifies both light and dark variants for each color role. If you implemented only the dark variants, go back to the agent and ask for the full adaptive palette. The agent uses `lipgloss.AdaptiveColor` or equivalent to switch automatically.

**ASCII art breaks at narrow widths**
The agent limits ASCII art to under 60 characters wide for this reason. If you have wider art, ask the agent for a compact alternative. In general, ASCII art is treated as a liability — it often looks wrong in different font sizes, line heights, or terminal emulators.

**Output hard to scan despite colors**
Color alone does not create hierarchy. The agent uses whitespace as the primary layout tool — blank lines between sections, indentation for sub-items, alignment for related values. If output is hard to scan, the issue is usually spacing and grouping, not color. Ask the agent to redesign the whitespace structure.
