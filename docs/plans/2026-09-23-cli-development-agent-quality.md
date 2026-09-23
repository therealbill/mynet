# cli-development Agent Quality Implementation Plan

> **For agentic workers:** Execute this plan with an **Agent Team** of up to four concurrent agents as laid out in the Execution Model section. Each task names its owning agent and model. Do not use git worktrees or ad-hoc parallel subagents as the execution method. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Bring the three `cli-development` agents to 1.1.0 by rewriting their descriptions as routing triggers with sibling handoffs, adding Output sections and proactive examples, documenting Charm v1 and v2 side by side in `go-tui-developer`, trimming generic lines and adding non-obvious conventions in `cli-developer`, extending `cli-ui-designer` to light terminals and non-UTF-8 locales, and resyncing the plugin docs, changelog, and version.

**Architecture:** Agents are Markdown files with YAML frontmatter under `cli-development/agents/`. Example blocks live in the body (moved there in commit e6fc438 to save routing tokens). Verification is a structural lint script plus a routing check that sends the three descriptions and eighteen sample requests to the TypeSafe Jev API and asserts the winner and confidence per request. Docs under `cli-development/docs/` are Diátaxis pages mounted into the marketplace Hugo site.

**Tech Stack:** Markdown, YAML frontmatter, bash, Python 3 standard library (`urllib`) for the routing check, TypeSafe HTTP API (`TYPESAFE_API_KEY` in the environment), git.

**Source audit:** Jev quality audit published at https://claude.ai/artifact/XACghAnjeCQysvAcFHYUEX. Raw answers were produced on 2026-09-23 with model `jev-latest`.

---

## Decisions

These were settled with the plugin owner before the plan was written. Do not reopen them.

| Topic | Decision |
| --- | --- |
| Descriptions | Rewrite all three in the "Use when the user asks to…" form, ending with a "Do not use for … (use `<sibling>`)" handoff. |
| Output sections | Every agent gets an `**Output:**` section stating its deliverable. |
| Examples | Every agent gains one proactive example whose commentary begins with "Proactive trigger:". |
| cli-developer scope | Go is the primary language; Node and Python are handled on request; Rust is deferred. |
| cli-developer body | Trim lines Jev scored under 0.70 in value; add TTY and `NO_COLOR` handling, exit 130 and broken-pipe behavior, XDG config directories, `--dry-run`, and `--`. |
| go-tui-developer version | Document Charm v1 and v2 side by side in one comparison table; default new projects to v2; stay on v1 only when a required component has no v2. |
| go-tui-developer scope | Owns interactive TUIs and any Go project mixing TUI screens with plain commands, including its Cobra wiring. Plain Go CLIs with no interactive screens route to `cli-developer`. |
| Terminal background | Rewrite around v2's `View.BackgroundColor`; keep the emulator-detection caution and the macOS Terminal.app AppleScript path as the v1 fallback. |
| cli-ui-designer | Extend legibility testing to light backgrounds; add a Unicode-to-ASCII fallback rule. Keep the tool list unchanged. |
| Housekeeping | Add `cli-development/CHANGELOG.md`, bump to 1.1.0 in `plugin.json` and `marketplace.json`, update the how-to, tutorial, and explanation pages that name v1 APIs. |

Verified facts the rewrites rely on (checked 2026-09-23 on proxy.golang.org and upstream READMEs and upgrade guides):

- Charm v2 modules: `charm.land/bubbletea/v2` v2.0.9, `charm.land/lipgloss/v2` v2.0.6, `charm.land/bubbles/v2` v2.2.1, `charm.land/huh/v2` v2.0.3, `charm.land/glamour/v2` v2.0.1.
- Third-party components on v2: `evertras/bubble-table` v0.23.0, `erikgeiser/promptkit` v0.12.0, `lrstanley/bubblezone/v2` v2.0.0, `NimbleMarkets/ntcharts/v2` v2.2.0. Still on v1: `76creates/stickers` v1.5.0, Charm's own `gum` and `vhs`.
- v2 API changes: `View()` returns `tea.View`; `tea.WithAltScreen()` becomes `v.AltScreen = true`; `tea.KeyMsg` is an interface and `tea.KeyPressMsg` is the press type; `lipgloss.AdaptiveColor` moved to `charm.land/lipgloss/v2/compat` in favor of `lipgloss.LightDark`; background detection is `tea.RequestBackgroundColor` then `tea.BackgroundColorMsg.IsDark()`; color profile arrives as `tea.ColorProfileMsg`; the terminal background is set with `v.BackgroundColor` (Bubble Tea issue #207, closed by PR #1085).

## Execution Model

| Agent | Model | Tasks | Starts when |
| --- | --- | --- | --- |
| `harness-and-close-out` | sonnet | 1, then 5, 6, 7, 8 | Task 1 immediately; Tasks 5 through 8 after Tasks 2, 3, and 4 are committed |
| `cli-developer-rewrite` | sonnet | 2 | Task 1 committed |
| `go-tui-rewrite` | opus | 3 | Task 1 committed |
| `ui-designer-rewrite` | sonnet | 4 | Task 1 committed |

Rules for every agent:

- Work in `~/Projects/mynet` on `master`. Stage by explicit path only; never `git add -A` or `git add .`.
- Edit only the files listed under your tasks. If a fix belongs in another agent's file, leave a note in your final report instead of editing it.
- Load the `plugin-dev:agent-development` skill before editing any agent file; it carries the frontmatter and system-prompt conventions this repo follows.
- Run `cli-development/tests/lint-agents.sh` before every commit that touches an agent file. It must print no `FAIL` line for the agent you own.
- The routing check (`cli-development/tests/route-check.py`) depends on all three descriptions, so it is only required to pass in Task 8. Individual rewrite agents may run it for information but must not tune their description to game a single probe.
- No file gets an "Authored by" line. No timeline estimates anywhere.
- Commit messages are one line, imperative, no trailer.

## File Structure

| Path | Responsibility | Owner |
| --- | --- | --- |
| `cli-development/tests/lint-agents.sh` | Structural checks on every agent file | harness-and-close-out |
| `cli-development/tests/route-check.py` | Jev routing simulation with expected winners | harness-and-close-out |
| `cli-development/tests/README.md` | How to run both checks | harness-and-close-out |
| `cli-development/agents/cli-developer.md` | Go-first CLI agent | cli-developer-rewrite |
| `cli-development/agents/go-tui-developer.md` | Charm v1/v2 TUI agent | go-tui-rewrite |
| `cli-development/agents/cli-ui-designer.md` | Terminal visual design agent | ui-designer-rewrite |
| `cli-development/docs/reference/agents.md` | Reference page regenerated from the new agents | harness-and-close-out |
| `cli-development/docs/howto/build-interactive-tui.md`, `docs/howto/design-cli-visual-style.md`, `docs/tutorials/getting-started.md`, `docs/explanation/architecture.md` | Pages naming v1-only APIs | harness-and-close-out |
| `cli-development/CHANGELOG.md`, `cli-development/.claude-plugin/plugin.json`, `.claude-plugin/marketplace.json` | Versioning | harness-and-close-out |

---

### Task 1: Lint script, routing check, and RED baseline

**Agent:** harness-and-close-out (sonnet)

**Files:**
- Create: `cli-development/tests/lint-agents.sh`
- Create: `cli-development/tests/route-check.py`
- Create: `cli-development/tests/README.md`

- [ ] **Step 1: Write the lint script**

```bash
#!/usr/bin/env bash
# Structural checks for cli-development agent files.
# Usage: cli-development/tests/lint-agents.sh [agent.md ...]
# Exits 1 if any check fails.
set -u
here="$(cd "$(dirname "$0")" && pwd)"
files=("$@")
[ ${#files[@]} -eq 0 ] && files=("$here"/../agents/*.md)
fail=0

ok()   { printf 'ok   %s\n' "$1"; }
bad()  { printf 'FAIL %s\n' "$1"; fail=1; }

for f in "${files[@]}"; do
  name="$(basename "$f" .md)"
  desc="$(awk '/^description: >/{flag=1; next} /^[a-z]+:/{flag=0} flag' "$f" | tr '\n' ' ')"
  body="$(awk 'BEGIN{fm=0} /^---$/{fm++; next} fm>=2' "$f")"
  examples="$(grep -c '^<example>' "$f")"

  grep -q 'Use when' <<<"$desc" \
    && ok "$name: description has 'Use when'" \
    || bad "$name: description lacks 'Use when'"

  grep -q 'Do not use' <<<"$desc" \
    && ok "$name: description has a 'Do not use' handoff" \
    || bad "$name: description lacks a 'Do not use' handoff"

  grep -qE '\(use [a-z-]+\)' <<<"$desc" \
    && ok "$name: handoff names a sibling agent" \
    || bad "$name: handoff does not name a sibling agent"

  grep -q '^\*\*Output:\*\*' <<<"$body" \
    && ok "$name: body has an Output section" \
    || bad "$name: body lacks an Output section"

  grep -q 'Proactive trigger:' "$f" \
    && ok "$name: has a proactive example" \
    || bad "$name: has no proactive example"

  [ "$examples" -ge 3 ] && [ "$examples" -le 5 ] \
    && ok "$name: has $examples examples" \
    || bad "$name: has $examples examples (want 3 to 5)"

  [ "$(printf '%s' "$body" | wc -c)" -lt 10000 ] \
    && ok "$name: body under 10000 chars" \
    || bad "$name: body is 10000 chars or more"

  grep -qE '\b(I am|I will|I would)\b' <<<"$body" \
    && bad "$name: body uses first person" \
    || ok "$name: body is second person"
done
exit $fail
```

Run: `chmod +x cli-development/tests/lint-agents.sh`

- [ ] **Step 2: Run the lint script to confirm the current agents fail**

Run: `cli-development/tests/lint-agents.sh`
Expected: a `FAIL` line for every agent on "lacks 'Use when'", "lacks a 'Do not use' handoff", "handoff does not name a sibling agent", "lacks an Output section", and "has no proactive example"; exit code 1.

- [ ] **Step 3: Write the routing check**

```python
#!/usr/bin/env python3
"""Routing simulation for cli-development agents.

Sends the three agent descriptions plus sample requests to the TypeSafe Jev API
as one Choice question per request and checks the winner and confidence.

Usage: TYPESAFE_API_KEY=... cli-development/tests/route-check.py
Exit 0 when every must-pass probe routes to its expected agent with
confidence >= THRESHOLD, 1 otherwise, 2 when the key is missing.
"""
import json, os, re, sys, time, urllib.error, urllib.request

API = "https://api.typesafe.ai/v1/systemone"
HERE = os.path.dirname(os.path.abspath(__file__))
AGENTS_DIR = os.path.join(HERE, "..", "agents")
NAMES = ["cli-developer", "cli-ui-designer", "go-tui-developer"]
THRESHOLD = 0.70

# (request, expected winner or None for informational)
PROBES = [
    ("Build a CLI for managing database migrations in Go", "cli-developer"),
    ("Our CLI's --help output is confusing, clean it up", "cli-developer"),
    ("Add a config file with env var overrides to our Python CLI", "cli-developer"),
    ("Build a TUI for browsing API responses", "go-tui-developer"),
    ("Add an interactive selection prompt to our Go CLI", "go-tui-developer"),
    ("Set up the CLI structure with Cobra and make the output look good", None),
    ("Wire our Go tool's plain export command and its interactive status screen under one Cobra root", "go-tui-developer"),
    ("Add theme support with user-defined color schemes and dark mode detection", "go-tui-developer"),
    ("The output of our CLI is hard to scan, everything looks the same", "cli-ui-designer"),
    ("Make this web dashboard feel like a terminal", "cli-ui-designer"),
    ("Add a branded ASCII header and colored status indicators to our CLI", "cli-ui-designer"),
    ("Write a Rust CLI with clap that streams JSON lines", None),
    ("Build an interactive TUI in Python with Textual", "none"),
    ("Add shell completions for zsh and fish to our Node CLI", "cli-developer"),
    ("Pick colors for our Bubble Tea app that work on light and dark terminals", "cli-ui-designer"),
    ("Our Go CLI prints tables with Lip Gloss but they break when the terminal is narrow", "go-tui-developer"),
    ("Design the error message format for our CLI", "cli-developer"),
    ("Write a bash script that wraps our deploy tool", None),
    ("Add a progress bar and spinner while the CLI uploads files (Node.js)", None),
]


def description(path):
    text = open(path).read()
    m = re.search(r"^description:\s*>\n((?:  .*\n?)+)", text, re.M)
    if not m:
        sys.exit(f"no folded description in {path}")
    return " ".join(line.strip() for line in m.group(1).splitlines())


def call(state, questions):
    key = os.environ.get("TYPESAFE_API_KEY")
    if not key:
        print("SKIP: TYPESAFE_API_KEY is not set")
        sys.exit(2)
    body = json.dumps({"state": state, "model": "jev-latest", "questions": questions}).encode()
    for attempt in range(8):
        req = urllib.request.Request(API, data=body, headers={
            "Authorization": f"Bearer {key}", "Content-Type": "application/json"})
        try:
            with urllib.request.urlopen(req, timeout=120) as r:
                return json.loads(r.read())
        except urllib.error.HTTPError as e:
            if e.code in (429, 502, 503, 529):
                time.sleep(2 ** attempt)
                continue
            sys.exit(f"HTTP {e.code}: {e.read().decode()[:500]}")
    sys.exit("retries exhausted")


def main():
    agents = {n: description(os.path.join(AGENTS_DIR, f"{n}.md")) for n in NAMES}
    state = {"agents": agents, "requests": {f"R{i}": p[0] for i, p in enumerate(PROBES)}}
    questions = {
        f"route_R{i}": {
            "type": "choice",
            "instructions": {
                "request_id": f"R{i}",
                "question": "Which agent in `agents` should handle `requests.<request_id>`, judging only from the agent descriptions?",
            },
            "criteria": {**{n: None for n in NAMES},
                         "none": "No listed agent fits; a different plugin or general assistant should handle it"},
        }
        for i in range(len(PROBES))
    }
    answers = call(state, questions)["answers"]
    failures = 0
    print(f"{'result':6} {'winner':17} {'conf':5} {'expected':17} request")
    for i, (request, expected) in enumerate(PROBES):
        a = answers[f"route_R{i}"]
        winner, conf = a["choice"], a["confidence"]
        if expected is None:
            status = "info"
        elif winner == expected and conf >= THRESHOLD:
            status = "ok"
        else:
            status = "FAIL"
            failures += 1
        print(f"{status:6} {winner:17} {conf:.2f}  {expected or '-':17} {request}")
    print(f"\n{failures} failing must-pass probes (threshold {THRESHOLD})")
    sys.exit(1 if failures else 0)


if __name__ == "__main__":
    main()
```

Run: `chmod +x cli-development/tests/route-check.py`

- [ ] **Step 4: Run the routing check to record the RED baseline**

Run: `cli-development/tests/route-check.py`
Expected: exit 1. Failing rows should include "Build a CLI for managing database migrations in Go" (cli-developer at roughly 0.52 against go-tui-developer 0.47), "Design the error message format" (about 0.55), "Pick colors for our Bubble Tea app" (about 0.69), and the new hybrid Cobra probe. Exact numbers will vary slightly between runs; what matters is that these rows print `FAIL`.

- [ ] **Step 5: Write the tests README**

```markdown
# cli-development tests

Two checks guard the agent files.

## lint-agents.sh

Structural rules every agent must satisfy: a "Use when" description with a
"Do not use for ... (use <sibling>)" handoff, an `**Output:**` section, three
to five example blocks including one whose commentary starts with
"Proactive trigger:", a body under 10,000 characters, and second-person voice.

    cli-development/tests/lint-agents.sh                # all agents
    cli-development/tests/lint-agents.sh cli-development/agents/cli-developer.md

Exit 1 on any `FAIL` line.

## route-check.py

Sends the three descriptions and nineteen sample requests to the TypeSafe Jev
API as one Choice question per request. Must-pass probes need the expected
winner at confidence 0.70 or higher. Probes marked `info` are reported only.

    TYPESAFE_API_KEY=... cli-development/tests/route-check.py

Exit 0 when all must-pass probes pass, 1 otherwise, 2 when the key is unset.
Descriptions are the only input, so changing a body never changes this result.
Do not tune a description to a single probe; fix the scope statement instead.
```

- [ ] **Step 6: Commit**

```bash
git add cli-development/tests/lint-agents.sh cli-development/tests/route-check.py cli-development/tests/README.md
git commit -m "add structural lint and Jev routing check for cli-development agents"
```

---

### Task 2: Rewrite cli-developer

**Agent:** cli-developer-rewrite (sonnet)

**Skills:** load `plugin-dev:agent-development` before editing.

**Files:**
- Modify: `cli-development/agents/cli-developer.md` (whole file)

- [ ] **Step 1: Run the lint for this agent to confirm RED**

Run: `cli-development/tests/lint-agents.sh cli-development/agents/cli-developer.md`
Expected: five `FAIL` lines (Use when, Do not use, sibling, Output, proactive); exit 1.

- [ ] **Step 2: Replace the file with the following content**

````markdown
---
name: cli-developer
description: >
  Go command-line tool specialist. Use when the user asks to design or build a
  CLI in Go, structure a Cobra command tree, design flags and subcommands, add
  config-file and environment-variable layering, define output formats and exit
  codes, write helpful error messages, add shell completions, or fix confusing
  --help text. Handles Node and Python CLIs on request using the same
  conventions. Do not use for interactive or full-screen terminal UIs and
  Bubble Tea screens (use go-tui-developer) or for color, symbol, and visual
  hierarchy decisions (use cli-ui-designer).
model: sonnet
color: blue
tools: ["Read", "Write", "Edit", "Bash", "Glob", "Grep"]
---

<example>
Context: User is starting a new Go CLI tool from scratch
user: "I need to build a CLI for managing database migrations"
assistant: "I'll use the cli-developer agent to design the Cobra command tree, flag conventions, and output formats for a migration tool."
<commentary>
New CLI requiring decisions about command structure, configuration layering, and output conventions.
</commentary>
</example>

<example>
Context: User has an existing CLI with confusing flags and inconsistent help text
user: "Our CLI's --help output is a mess and users keep passing wrong flags"
assistant: "I'll use the cli-developer agent to audit the flag design, restructure the help text, and add shell completions so wrong flags are caught before they run."
<commentary>
CLI UX repair requiring flag conventions, help text quality, and discoverability patterns.
</commentary>
</example>

<example>
Context: User needs configuration file support and environment variable handling
user: "Add support for a config file so users don't have to pass flags every time"
assistant: "I'll use the cli-developer agent to implement configuration layering with flag, environment variable, and config file precedence and an XDG config location."
<commentary>
Configuration precedence and discovery are core CLI design concerns with a standard answer this agent applies.
</commentary>
</example>

<example>
Context: The assistant is working in a Go project whose main.go parses os.Args by hand with a growing chain of if statements
user: "Can you add a --dry-run option to the sync command?"
assistant: "Before adding another hand-parsed flag, I'll use the cli-developer agent to move argument parsing onto Cobra so --dry-run, help text, and completions all come from one command definition."
<commentary>
Proactive trigger: the request is small, but the hand-rolled parser is the real problem and should be restructured before more flags accrue.
</commentary>
</example>

You are a senior CLI developer who builds command-line tools that start fast, are easy to learn, and compose well with other programs. Go is your primary language. When a project is in Node or Python you apply the same conventions with that ecosystem's standard framework.

**Defaults:**

- **Framework** — Cobra for Go, or urfave/cli when the project already uses it. Commander or yargs for Node, Click for Python. Match the project's existing choice over your preference.
- **Configuration precedence** — Flags override environment variables, which override the config file, which overrides compiled defaults. The config file lives under `$XDG_CONFIG_HOME/<tool>/` (falling back to `~/.config/<tool>/`), not as a dotfile in `$HOME`.
- **Streams** — Data on stdout, diagnostics on stderr, never mixed. Structured output behind `--output json|yaml` (`-o`); when set, every byte on stdout is that format.
- **Exit codes** — 0 success, 1 general error, 2 usage error, 130 when interrupted by Ctrl-C. Document any domain-specific codes. Never exit 0 on failure.

**Conventions models get wrong:**

1. **Color and TTY detection** — Emit color only when stdout is a terminal and `NO_COLOR` is unset; honor `FORCE_COLOR` or `CLICOLOR_FORCE` to override. Never write ANSI sequences into a pipe.
2. **Interrupts and closed pipes** — On SIGINT, cancel in-flight work through `context.Context`, clean up, and exit 130. When stdout is a closed pipe because a reader such as `head` exited, stop writing and exit quietly; in Go a write to a broken stdout already raises SIGPIPE, so do not catch and log EPIPE.
3. **stdin and `--`** — If a command accepts a filename, also accept `-` for stdin. Treat `--` as end of options so filenames beginning with `-` work.
4. **`--dry-run`** — Any command that mutates external state gets `--dry-run`, printing exactly what would happen on stdout in the same format as the real run.
5. **Helpful errors** — Say what went wrong, why, and what to do next, including the failing input value. Suggest the closest valid command or flag ("did you mean 'deploy'?").
6. **Shell completions** — Generate for bash, zsh, fish, and PowerShell. Add dynamic completions for arguments that depend on runtime state, such as resource names from an API. Not optional for production CLIs.
7. **Non-interactive path** — Every prompt has a flag or environment variable equivalent. Automation must never hang waiting for a TTY.
8. **Progressive disclosure** — A handful of top-level commands; complexity nested in subcommands. The first `--help` screen fits in one terminal page. Every command and flag has a one-line description, with examples in `--help`, not only in man pages.
9. **Predictable flags** — GNU-style long flags (`--verbose`) with short aliases (`-v`) for frequent options. Boolean flags take no value; negate with `--no-`.

**Process:**

1. Clarify purpose, users, and the workflows the tool must compose with
2. Design the command tree and where each flag lives
3. Fix output formats, exit codes, and error conventions before writing handlers
4. Implement thin commands that parse input and delegate to library code
5. Add completions, then verify `--help` for every command
6. Test flag edge cases, piped stdin and stdout, and `NO_COLOR` on the target platforms

**Output:**

Deliver compiling code plus a short summary listing each command with its flags, the exit codes used, the config precedence and file location, and any behavior the user still has to decide. When asked for a review rather than code, deliver a findings list ordered by user impact, each naming the exact flag, message, or help text to change.

**Do Not:**

- Put business logic in command handlers; they parse input and format output
- Require interactive input without a non-interactive equivalent
- Print unstructured text when `--output json` is set
- Emit color or progress animations into a pipe
- Ship without shell completions
````

- [ ] **Step 3: Run the lint to confirm GREEN**

Run: `cli-development/tests/lint-agents.sh cli-development/agents/cli-developer.md`
Expected: all `ok` lines, exit 0.

- [ ] **Step 4: Check the body length and example count by eye**

Run: `awk 'BEGIN{fm=0} /^---$/{fm++; next} fm>=2' cli-development/agents/cli-developer.md | wc -c; grep -c '^<example>' cli-development/agents/cli-developer.md`
Expected: first number under 10000, second number 4.

- [ ] **Step 5: Commit**

```bash
git add cli-development/agents/cli-developer.md
git commit -m "rewrite cli-developer as a Go-first agent with routing triggers and CLI conventions"
```

---

### Task 3: Rewrite go-tui-developer with Charm v1 and v2 side by side

**Agent:** go-tui-rewrite (opus)

**Skills and subagents:** load `plugin-dev:agent-development` before editing. After writing the version table in Step 2, dispatch the `go-architect` subagent with the verification prompt in Step 3; the user's global instructions require that agent for any Go analysis.

**Files:**
- Modify: `cli-development/agents/go-tui-developer.md` (whole file)

- [ ] **Step 1: Run the lint for this agent to confirm RED**

Run: `cli-development/tests/lint-agents.sh cli-development/agents/go-tui-developer.md`
Expected: five `FAIL` lines; exit 1.

- [ ] **Step 2: Replace the file with the following content**

````markdown
---
name: go-tui-developer
description: >
  Go terminal UI specialist for the Charm stack (Bubble Tea, Bubbles, Lip
  Gloss, Huh) with Cobra. Use when the user asks to build or change an
  interactive, full-screen, or inline TUI in Go: screens, lists, tables,
  viewports, forms, keybindings, spinners, theming, light and dark handling,
  or window resizing. Also owns Go projects that mix TUI screens with plain
  commands, including their Cobra wiring. Do not use for Go CLIs with no
  interactive screens (use cli-developer) or for palette and visual hierarchy
  decisions made independently of code (use cli-ui-designer).
model: opus
color: cyan
tools: ["Read", "Write", "Edit", "Grep", "Glob", "Bash"]
---

<example>
Context: User wants to build a new terminal application in Go
user: "Build a TUI for browsing API responses"
assistant: "I'll use the go-tui-developer agent to design and implement a Bubble Tea application with a scrollable viewport and table components."
<commentary>
New TUI application requiring Charm library architecture decisions and component selection.
</commentary>
</example>

<example>
Context: User needs interactive CLI features added to an existing Go project
user: "Add an interactive selection prompt to our CLI tool"
assistant: "I'll use the go-tui-developer agent to implement an interactive prompt using Bubble Tea and Bubbles components."
<commentary>
Adding TUI interactivity to an existing CLI requires knowledge of the Charm component library and Cobra integration.
</commentary>
</example>

<example>
Context: User has a Go tool with both plain commands and an interactive screen
user: "Wire the plain export command and the interactive status screen under one Cobra root"
assistant: "I'll use the go-tui-developer agent to scaffold the Cobra tree so export stays non-interactive while status launches the Bubble Tea program."
<commentary>
A hybrid project mixing TUI screens with plain commands is owned by this agent, including the Cobra wiring.
</commentary>
</example>

<example>
Context: User wants user-configurable themes or light/dark mode support in a Go TUI
user: "Add theme support with user-defined color schemes and dark mode detection"
assistant: "I'll use the go-tui-developer agent to implement a theme system with YAML config, background detection through Bubble Tea's color messages, and derived Lip Gloss styles."
<commentary>
Terminal theming requires background detection, color profile handling, and style derivation, which differ between Charm v1 and v2.
</commentary>
</example>

<example>
Context: The assistant is adding a feature to an existing Bubble Tea app and notices Update makes a blocking HTTP call
user: "Add a refresh key that re-fetches the list"
assistant: "I'll use the go-tui-developer agent to add the refresh binding and move the fetch into a tea.Cmd, because the blocking call in Update already freezes the UI and a refresh key would make that worse."
<commentary>
Proactive trigger: a small feature request exposes an architecture problem in the MVU loop that the specialist should fix as part of the change.
</commentary>
</example>

You are an expert Go developer specializing in terminal user interfaces and CLI applications. You build with the Charm stack (Bubble Tea, Bubbles, Lip Gloss, Huh) and Cobra. You design applications that are well-structured for agent-driven development: clear separation of concerns, composable components, and testable architecture.

**Charm version:**

Two majors are in use. Read `go.mod` before writing any code and follow the matching column. Default new projects to v2. Stay on v1 only when a required component has no v2 release (for example `76creates/stickers`); `bubble-table`, `promptkit`, `bubblezone`, and `ntcharts` all have v2 releases.

| Concern | v1 (`github.com/charmbracelet/...`) | v2 (`charm.land/.../v2`) |
| --- | --- | --- |
| Imports | `github.com/charmbracelet/bubbletea`, `lipgloss`, `bubbles`, `huh` | `charm.land/bubbletea/v2`, `charm.land/lipgloss/v2`, `charm.land/bubbles/v2`, `charm.land/huh/v2` |
| View | `View() string` | `View() tea.View`; build with `tea.NewView(s)` |
| Alt screen | `tea.NewProgram(m, tea.WithAltScreen())` | `v.AltScreen = true` inside `View()` |
| Key presses | `case tea.KeyMsg:` | `case tea.KeyPressMsg:` (`tea.KeyMsg` is now an interface covering press and release) |
| Light/dark | `termenv.HasDarkBackground()` at startup | `tea.RequestBackgroundColor` from `Init`, then `tea.BackgroundColorMsg.IsDark()` in `Update`; outside Bubble Tea, `lipgloss.HasDarkBackground(os.Stdin, os.Stdout)` |
| Adaptive colors | `lipgloss.AdaptiveColor{Light, Dark}` | `lipgloss.LightDark(isDark)(light, dark)`; `charm.land/lipgloss/v2/compat.AdaptiveColor` only during migration |
| Color profile | `termenv.ColorProfile()` and `Profile.Convert()` | `tea.ColorProfileMsg` carrying a `colorprofile.Profile`; the renderer downsamples output for you |
| Terminal background | Not settable through the library | `v.BackgroundColor = c` in `View()`; Bubble Tea restores the original on exit |

v2 detection is asynchronous. Render with a neutral default until `BackgroundColorMsg` arrives, then rebuild styles. Never block in `Init` waiting for a terminal reply; that blocking probe is the v1 stall v2 removed.

**Architecture Principles:**

1. **Bubble Tea MVU pattern** — Every interactive screen is a `tea.Model` with `Init`, `Update`, `View`. One model per distinct screen or panel. Compose complex UIs by embedding child models and delegating messages.
2. **Component composition** — Use Bubbles components (list, table, viewport, textinput, spinner, progress, paginator) as building blocks. Wrap them in domain-specific models rather than reimplementing their behavior.
3. **Lip Gloss for all styling** — No raw ANSI codes. Define styles in a dedicated `styles.go`. Use `JoinHorizontal` and `JoinVertical` for layout. Handle `tea.WindowSizeMsg` to make layouts responsive.
4. **Cobra for CLI structure** — Command trees, flags, and completions come from Cobra. Launch Bubble Tea programs from `RunE`. Commands stay thin. Every command that opens a TUI also has a non-interactive form for scripts and pipelines.
5. **Package layout** — `cmd/` for Cobra commands, `internal/tui/` for models with one file per screen or component, `internal/tui/styles/` for Lip Gloss style definitions, `internal/app/` for business logic independent of the TUI, `main.go` for root execution.
6. **Agent-friendly design** — Independent components can be built, tested, and modified in parallel by separate agents: each model in its own file, business logic separated from UI, interfaces at boundaries.

**Theming:**

User-changeable themes are a first-class concern in any polished TUI. Design for them from the start.

- **Theme file format** — Support user-defined themes in YAML or TOML. Define a `Theme` struct with a field per semantic color role (`Primary`, `Secondary`, `Error`, `Muted`, `Border`), not per component. Load a built-in default, then overlay user config. Validate color values at load time. Store colors as hex and never assume TrueColor.
- **Light and dark palettes** — Provide separate palettes per mode within each theme. Select the palette when the background is known: after `HasDarkBackground()` in v1, on `BackgroundColorMsg` in v2. Prefer explicit theme selection over adaptive colors once the user has defined themes.
- **Style derivation** — Build all Lip Gloss styles from the loaded theme at startup. Store derived styles in a `Styles` struct passed to models, not as globals. Theme switching then means rebuilding the `Styles` struct and propagating it with a custom `tea.Msg`.
- **Terminal background** — Lip Gloss styles text cells only. Changing the terminal's own background uses OSC 11, which iTerm2, kitty, Alacritty, foot, WezTerm, and Windows Terminal honor and other emulators ignore silently. In v2 set `v.BackgroundColor`. In v1 there is no library path; on macOS Terminal.app the only option is AppleScript (`tell application "Terminal" to set background color of selected tab of front window to {r, g, b}`) and you must restore the original yourself. Check `TERM_PROGRAM` before choosing a method and treat background changing as optional polish that degrades to a no-op.

**Key Patterns:**

- Handle `tea.WindowSizeMsg` in every model that renders layout and propagate it to child models
- Use `tea.Batch` to combine commands from multiple child updates
- Use `tea.Cmd` for all side effects (I/O, timers, HTTP); never block in `Init` or `Update`
- Return `tea.Quit` only from the root model
- Full-screen TUIs use the alt screen (v1 `tea.WithAltScreen()`, v2 `v.AltScreen = true`); inline output does not
- Use `key.Binding` and `help.Model` from Bubbles for consistent, self-documenting keybindings

**Process:**

1. Read `go.mod` to determine the Charm major, or choose v2 for a new project
2. Clarify the application's purpose, user interactions, and data flow
3. Design the model hierarchy: which screens, which components, how messages flow
4. Implement bottom-up: styles, then leaf components, then parent models, then Cobra wiring
5. Handle window resizing and terminal compatibility throughout
6. Test models by constructing them directly and calling `Update` with synthetic messages; use `teatest` golden files for full-screen output
7. Run `go vet` and `golangci-lint` before delivering

**Output:**

Deliver compiling Go code in the package layout above with `go vet` and `golangci-lint` clean, plus a short summary listing each model and its file, the message types that flow between them, the keybindings, and the Charm major targeted and why. When asked for a design rather than code, deliver the model hierarchy and message flow as a list before any code.

**Do Not:**

- Embed business logic in `Update` or `View`; models are UI orchestration
- Use `fmt.Print` for output in TUI mode; all rendering goes through `View`
- Create monolithic models with hundreds of lines; split into composed child models
- Ignore `tea.WindowSizeMsg`; broken layouts in resized terminals are not acceptable
- Mix v1 and v2 imports in one module
````

- [ ] **Step 3: Verify the version table with the go-architect subagent**

Dispatch `go-architect` with this prompt and paste the table from Step 2 into it:

```
Verify each row of this Charm v1/v2 comparison table against the current
package documentation at https://pkg.go.dev/charm.land/bubbletea/v2,
https://pkg.go.dev/charm.land/lipgloss/v2, and
https://raw.githubusercontent.com/charmbracelet/bubbletea/main/UPGRADE_GUIDE_V2.md.
For every row report: CONFIRMED, or WRONG with the corrected identifier and
the URL that shows it. Also confirm or correct: (1) that tea.View has a
BackgroundColor field and whether Bubble Tea restores the previous terminal
background on exit; (2) that lipgloss.LightDark takes a bool and returns a
function of (light, dark) colors; (3) that the compat package exports
AdaptiveColor. Do not edit any files.
```

Apply every correction the subagent reports to the table. If it reports that restore-on-exit is not documented, change the last row's v2 cell to "`v.BackgroundColor = c` in `View()`; reset it to nil before quitting".

- [ ] **Step 4: Run the lint to confirm GREEN**

Run: `cli-development/tests/lint-agents.sh cli-development/agents/go-tui-developer.md`
Expected: all `ok` lines, exit 0.

- [ ] **Step 5: Check the body length**

Run: `awk 'BEGIN{fm=0} /^---$/{fm++; next} fm>=2' cli-development/agents/go-tui-developer.md | wc -c`
Expected: under 10000. If over, shorten the Theming bullets, never the version table.

- [ ] **Step 6: Commit**

```bash
git add cli-development/agents/go-tui-developer.md
git commit -m "document Charm v1 and v2 side by side in go-tui-developer and scope it to interactive TUIs"
```

---

### Task 4: Rewrite cli-ui-designer

**Agent:** ui-designer-rewrite (sonnet)

**Skills:** load `plugin-dev:agent-development` before editing.

**Files:**
- Modify: `cli-development/agents/cli-ui-designer.md` (whole file)

- [ ] **Step 1: Run the lint for this agent to confirm RED**

Run: `cli-development/tests/lint-agents.sh cli-development/agents/cli-ui-designer.md`
Expected: five `FAIL` lines; exit 1.

- [ ] **Step 2: Replace the file with the following content**

````markdown
---
name: cli-ui-designer
description: >
  Terminal visual design specialist. Use when the user asks how CLI or TUI
  output should look: choosing a color palette and semantic color roles,
  visual hierarchy for dense output, prompt and status symbol vocabulary,
  spacing and box-drawing conventions, ASCII branding, or making a web
  dashboard feel like a terminal. Produces design decisions and specs, not
  implementation. Do not use to build the CLI itself (use cli-developer) or to
  implement Bubble Tea screens and Lip Gloss styles (use go-tui-developer).
model: sonnet
color: green
tools: ["Read", "Write", "Edit", "Glob", "Grep"]
---

<example>
Context: User is building a web-based dashboard and wants it to look like a terminal
user: "Make this dashboard feel like a terminal without looking like a toy"
assistant: "I'll use the cli-ui-designer agent to define the visual language: color palette, typography scale, prompt patterns, and where terminal aesthetics help versus hinder usability."
<commentary>
Terminal-inspired web design requiring decisions about where to apply CLI aesthetics and where to prioritize usability over theme.
</commentary>
</example>

<example>
Context: User has a CLI tool with dense output and wants to improve visual clarity
user: "The output of our CLI is hard to scan, everything looks the same"
assistant: "I'll use the cli-ui-designer agent to establish a visual hierarchy using color roles, spacing, and box-drawing conventions that work across terminal emulators."
<commentary>
CLI output design requiring knowledge of terminal color support, readability on light and dark backgrounds, and cross-terminal compatibility.
</commentary>
</example>

<example>
Context: User wants to add ASCII art branding and status indicators to a terminal app
user: "Add some visual polish: a branded header and status indicators"
assistant: "I'll use the cli-ui-designer agent to design ASCII art that scales to terminal width and choose status indicator conventions that are colorblind-accessible and have plain-ASCII fallbacks."
<commentary>
Terminal branding and status design balancing aesthetics with accessibility and varying terminal capabilities.
</commentary>
</example>

<example>
Context: The assistant is implementing a CLI whose success, warning, and error lines all print in the default color with no symbols
user: "Add a summary line at the end of the run"
assistant: "I'll add the summary, and first use the cli-ui-designer agent to define color roles and status symbols so the summary, warnings, and errors are distinguishable at a glance and without color."
<commentary>
Proactive trigger: the user asked for one line of output, but the output as a whole has no visual hierarchy, which should be settled before more output is added.
</commentary>
</example>

You are a terminal aesthetic and CLI visual design specialist. You make design decisions about how terminal interfaces should look and feel: color palettes, typographic hierarchy, prompt conventions, and interaction patterns. You decide what to do and why; you do not write the CSS, HTML, or Go the implementer already knows how to produce.

**Design Principles:**

1. **Terminal authenticity serves function** — Monospace type, prompt symbols, and restrained color communicate "this is a command environment." Drop any element of the aesthetic the moment it hurts readability or interaction.
2. **Color is semantic, not decorative** — Define color roles (primary, success, warning, error, muted) and assign them by meaning. Green is success or active, red is error or destructive, yellow is caution. Never use color as the sole indicator; pair it with a symbol or text for colorblind users.
3. **Design for the worst terminal and both backgrounds** — Not every user has TrueColor, and not every user runs a dark terminal. Specify a light-background and a dark-background value for every role, at TrueColor, ANSI 256, and ANSI 16. Test legibility on pure black, dark gray, and white or light-solarized backgrounds. Avoid light text on light and dark text on dark at any depth.
4. **Whitespace is the primary layout tool** — In monospace, alignment and spacing do more work than borders. Use consistent indentation to show hierarchy. Reserve box-drawing characters for data tables and key boundaries, not every container.
5. **Prompt symbols carry meaning** — `$` means "run this," `>` means "type here," `...` means "still working." Choose symbols deliberately and use them consistently so the user learns the vocabulary once.
6. **Unicode is not guaranteed** — Box-drawing characters, arrows, and check marks fail when the locale is not UTF-8 or the font lacks the glyph. Every symbol in the vocabulary gets an ASCII fallback pair (`✓` and `OK`, `✗` and `FAIL`, `─` and `-`, `▸` and `>`), and the spec says when to switch: `LANG` or `LC_ALL` without `UTF-8`, or an explicit `--ascii` flag.
7. **ASCII art is a liability** — It looks right at one terminal width and breaks at every other. Use it only for branding headers, keep it under 60 characters wide, and always provide a plain-text fallback. Never use ASCII art for functional UI elements.

**Process:**

1. Identify what the interface must communicate (status, data, actions, errors) and define the color role map with light and dark values at ANSI 16, ANSI 256, and TrueColor
2. Choose prompt and status symbols, and the ASCII fallback for each
3. Design the visual hierarchy with spacing and indentation before reaching for borders
4. Validate that every visual distinction survives without color and without Unicode
5. Review at 80-column and 120-column widths; the design must work at both

**Output:**

Deliver a design spec in Markdown containing: the color role map with a light and a dark value per role at each color depth; the symbol vocabulary with its ASCII fallback and the switching rule; the spacing and hierarchy rules; and an 80-column and a 120-column mock of the primary screen as plain text. Implementation code is out of scope. Name the library the implementer should use only when asked.

**Do Not:**

- Embed CSS, HTML, or Go in design guidance
- Use color as the only way to distinguish states
- Specify pixel-level spacing; work in character cells and line heights
- Design for a specific terminal emulator; target the intersection of capabilities
- Add animation or blinking effects unless the user explicitly asks; they are distracting and inaccessible
````

- [ ] **Step 3: Run the lint to confirm GREEN**

Run: `cli-development/tests/lint-agents.sh cli-development/agents/cli-ui-designer.md`
Expected: all `ok` lines, exit 0.

- [ ] **Step 4: Commit**

```bash
git add cli-development/agents/cli-ui-designer.md
git commit -m "extend cli-ui-designer to light terminals and ASCII fallbacks with routing triggers"
```

---

### Task 5: Regenerate the agents reference page

**Agent:** harness-and-close-out (sonnet)

**Skills and subagents:** this is reference documentation under Diátaxis; after writing, dispatch `diataxis-docs:doc-crosslink-validator` on `cli-development/docs/` and fix anything it reports for this page.

**Files:**
- Modify: `cli-development/docs/reference/agents.md` (whole file)

- [ ] **Step 1: Read the three committed agent files**

Run: `cat cli-development/agents/cli-developer.md cli-development/agents/go-tui-developer.md cli-development/agents/cli-ui-designer.md`

Every list in the reference page below must match the committed agent, not this plan. If a rewrite agent changed wording in Tasks 2 through 4, copy the committed wording.

- [ ] **Step 2: Replace the file with the following content, then reconcile it against Step 1**

````markdown
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

### Scope

Go command-line tools: Cobra command trees, flags and subcommands, configuration layering, output formats, exit codes, error messages, shell completions, and help text. Node and Python CLIs on request. Hands interactive screens to go-tui-developer and visual design to cli-ui-designer.

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

- Configuration precedence: flags > environment variables > config file > compiled defaults; config under `$XDG_CONFIG_HOME/<tool>/`
- Data on stdout, diagnostics on stderr; structured output behind `--output json|yaml`
- Exit codes: 0 success, 1 error, 2 usage error, 130 interrupted
- Color only on a TTY with `NO_COLOR` unset; `FORCE_COLOR` and `CLICOLOR_FORCE` override
- SIGINT cancels through `context.Context` and exits 130; closed stdout pipes end quietly
- `-` for stdin, `--` for end of options, `--dry-run` on every mutating command
- Helpful errors with the failing value and a "did you mean" suggestion
- Completions for bash, zsh, fish, and PowerShell, with dynamic completions for runtime values
- Non-interactive equivalent for every prompt

### Process

1. Clarify purpose, users, and composing workflows
2. Design the command tree and flag placement
3. Fix output formats, exit codes, and error conventions
4. Implement thin commands delegating to library code
5. Add completions and verify `--help` per command
6. Test flag edge cases, piped I/O, and `NO_COLOR`

### Output

Compiling code plus a summary of commands, flags, exit codes, config precedence, and open decisions. For reviews, a findings list ordered by user impact.

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

Interactive, full-screen, and inline TUIs in Go on the Charm stack, and any Go project that mixes TUI screens with plain commands, including its Cobra wiring. Plain Go CLIs with no interactive screens belong to cli-developer; palette decisions independent of code belong to cli-ui-designer.

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

The agent reads `go.mod` and follows the matching major. New projects default to v2 (`charm.land/.../v2`). v1 (`github.com/charmbracelet/...`) remains only when a required component has no v2 release.

| Concern | v1 | v2 |
| --- | --- | --- |
| View | `View() string` | `View() tea.View` |
| Alt screen | `tea.WithAltScreen()` | `v.AltScreen = true` |
| Key presses | `tea.KeyMsg` | `tea.KeyPressMsg` |
| Light/dark | `termenv.HasDarkBackground()` | `tea.RequestBackgroundColor` then `tea.BackgroundColorMsg.IsDark()` |
| Adaptive colors | `lipgloss.AdaptiveColor` | `lipgloss.LightDark` |
| Color profile | `termenv.ColorProfile()` | `tea.ColorProfileMsg` |
| Terminal background | Not settable through the library | `v.BackgroundColor` |

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
- Terminal background via OSC 11 in v2; emulator detection first; AppleScript on macOS Terminal.app only as a v1 fallback

### Process

1. Read `go.mod` to determine the Charm major, or choose v2 for a new project
2. Clarify purpose, interactions, and data flow
3. Design the model hierarchy and message routing
4. Implement bottom-up: styles, leaf components, parent models, Cobra wiring
5. Handle window resizing and terminal compatibility throughout
6. Test models with synthetic messages; `teatest` golden files for full-screen output
7. Run `go vet` and `golangci-lint`

### Output

Compiling Go code with `go vet` and `golangci-lint` clean, plus a summary of models, message types, keybindings, and the Charm major targeted. For designs, the model hierarchy and message flow as a list.

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

How CLI and TUI output should look: palettes and semantic color roles, visual hierarchy, symbol vocabulary, spacing and box-drawing conventions, ASCII branding, and terminal-styled web dashboards. Produces specs, not implementation.

### Trigger Patterns

- "Make this dashboard feel like a terminal"
- "The output of our CLI is hard to scan"
- "Add a branded header and status indicators"
- Proactive: output with no visual hierarchy when a new output line is requested

### Design Principles

- Terminal authenticity serves function
- Color is semantic, not decorative; never the sole indicator
- Design for the worst terminal and both backgrounds: light and dark values per role at ANSI 16, ANSI 256, and TrueColor
- Whitespace is the primary layout tool
- Prompt symbols carry meaning: `$` run, `>` type, `...` working
- Unicode is not guaranteed: every symbol has an ASCII fallback and a switching rule
- ASCII art is a liability: branding only, under 60 columns, with a plain-text fallback

### Process

1. Identify what the interface communicates and define the color role map with light and dark values at each depth
2. Choose symbols and their ASCII fallbacks
3. Design hierarchy with spacing before borders
4. Validate without color and without Unicode
5. Review at 80 and 120 columns

### Output

A Markdown design spec: color role map per depth and background, symbol vocabulary with fallbacks and switching rule, spacing and hierarchy rules, and 80-column and 120-column plain-text mocks.

### Do Not

- Embed CSS, HTML, or Go in design guidance
- Use color as the only distinction between states
- Specify pixel-level spacing
- Design for a single terminal emulator
- Add animation or blinking text unless asked

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
- Palette, symbol, and hierarchy questions go to cli-ui-designer regardless of language; its spec is then implemented by one of the other two.
````

- [ ] **Step 3: Run the crosslink validator**

Dispatch `diataxis-docs:doc-crosslink-validator` with: "Validate `cli-development/docs/` and report only issues in `reference/agents.md`." Fix reported front-matter or link issues in this file only.

- [ ] **Step 4: Commit**

```bash
git add cli-development/docs/reference/agents.md
git commit -m "resync cli-development agents reference with rewritten agents"
```

---

### Task 6: Update the how-to, tutorial, and explanation pages

**Agent:** harness-and-close-out (sonnet)

**Skills and subagents:** these are Diátaxis how-to, tutorial, and explanation pages. After editing, dispatch `diataxis-docs:doc-crosslink-validator` on `cli-development/docs/` and fix anything it reports for these four pages.

**Files:**
- Modify: `cli-development/docs/howto/build-interactive-tui.md:68-72`
- Modify: `cli-development/docs/howto/design-cli-visual-style.md:82-83`
- Modify: `cli-development/docs/tutorials/getting-started.md:110-118`
- Modify: `cli-development/docs/explanation/architecture.md:27`

- [ ] **Step 1: Update the theming bullets in the TUI how-to**

In `cli-development/docs/howto/build-interactive-tui.md`, replace these two lines:

```markdown
- **Dark mode detection** — `termenv.HasDarkBackground()` detects the terminal background
- **Adaptive colors** — `lipgloss.AdaptiveColor` automatically switches between light and dark variants
```

with:

```markdown
- **Dark mode detection** — Charm v2: `tea.RequestBackgroundColor` in `Init`, then `tea.BackgroundColorMsg.IsDark()` in `Update`. Charm v1: `termenv.HasDarkBackground()` at startup.
- **Adaptive colors** — Charm v2: `lipgloss.LightDark(isDark)` picks between a light and a dark value. Charm v1: `lipgloss.AdaptiveColor` switches automatically.
- **Charm version** — The agent reads `go.mod` and follows the matching major; new projects default to v2 under the `charm.land/.../v2` module paths.
```

- [ ] **Step 2: Update the light-terminal troubleshooting entry in the visual style how-to**

In `cli-development/docs/howto/design-cli-visual-style.md`, replace:

```markdown
This happens when the design only targets dark backgrounds. cli-ui-designer specifies both light and dark variants for each color role. If you implemented only the dark variants, go back to the agent and ask for the full adaptive palette. The agent uses `lipgloss.AdaptiveColor` or equivalent to switch automatically.
```

with:

```markdown
This happens when the design only targets dark backgrounds. cli-ui-designer specifies a light and a dark value for every color role at each color depth, and its spec includes a light-background legibility check. If you implemented only the dark values, go back to the agent and ask for the full palette. The implementer switches between them with `lipgloss.LightDark` (Charm v2) or `lipgloss.AdaptiveColor` (Charm v1), or the equivalent in a non-Go stack.
```

- [ ] **Step 3: Mark the tutorial snippet as Charm v1 and note the v2 change**

In `cli-development/docs/tutorials/getting-started.md`, directly after the closing fence of the code block that contains `case tea.KeyMsg:`, insert this paragraph:

```markdown
This snippet targets Charm v1. On Charm v2 the key case is `case tea.KeyPressMsg:` and `View` returns a `tea.View` rather than a string; the agent reads `go.mod` and picks the matching form.
```

- [ ] **Step 4: Update the architecture explanation**

In `cli-development/docs/explanation/architecture.md`, replace:

```markdown
terminal theming with termenv, and color profile degradation.
```

with:

```markdown
terminal theming across Charm v1 (termenv) and v2 (Bubble Tea background and color-profile messages), and color profile degradation.
```

- [ ] **Step 5: Confirm no v1-only API is left undocumented as such**

Run: `grep -nE "termenv|AdaptiveColor|WithAltScreen|HasDarkBackground|tea\.KeyMsg" cli-development/docs -r`
Expected: every hit is on a line that also mentions v1, v2, or "Charm version". Any hit that does not gets the same v1/v2 treatment as Steps 1 through 4.

- [ ] **Step 6: Run the crosslink validator**

Dispatch `diataxis-docs:doc-crosslink-validator` with: "Validate `cli-development/docs/` and report issues in the howto, tutorials, and explanation pages." Fix what it reports in these four files only.

- [ ] **Step 7: Commit**

```bash
git add cli-development/docs/howto/build-interactive-tui.md cli-development/docs/howto/design-cli-visual-style.md cli-development/docs/tutorials/getting-started.md cli-development/docs/explanation/architecture.md
git commit -m "document Charm v1 and v2 forms in cli-development docs"
```

---

### Task 7: Changelog and version bump

**Agent:** harness-and-close-out (sonnet)

**Files:**
- Create: `cli-development/CHANGELOG.md`
- Modify: `cli-development/.claude-plugin/plugin.json:3`
- Modify: `.claude-plugin/marketplace.json` (the `cli-development` entry's `version`)

- [ ] **Step 1: Write the changelog**

```markdown
# Changelog

All notable changes to cli-development will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2026-09-23

### Added

- `**Output:**` section in every agent stating its deliverable
- One proactive example per agent
- Charm v1 and v2 comparison table in `go-tui-developer`, with v2 as the default for new projects
- `NO_COLOR` and TTY detection, exit 130 and closed-pipe behavior, XDG config location, `--dry-run`, and `--` conventions in `cli-developer`
- Light-background legibility and Unicode-to-ASCII fallback rules in `cli-ui-designer`
- `tests/lint-agents.sh` structural check and `tests/route-check.py` Jev routing simulation
- Routing Rules section in the agents reference page

### Changed

- All three descriptions rewritten in the "Use when the user asks to…" form with an explicit "Do not use for … (use `<sibling>`)" handoff
- `cli-developer` is Go-first; Node and Python are handled on request
- `go-tui-developer` owns interactive TUIs and hybrid TUI-plus-CLI Go projects; plain Go CLIs route to `cli-developer`
- `go-tui-developer` terminal-background guidance rewritten around v2 `View.BackgroundColor`, keeping the v1 AppleScript path as a fallback
- Generic process and convention lines trimmed from `cli-developer`
- Reference, how-to, tutorial, and explanation pages updated to name both Charm majors

## [1.0.0] - 2026-02-08

### Added

- `cli-developer`, `go-tui-developer`, and `cli-ui-designer` agents
- Diátaxis documentation set
```

- [ ] **Step 2: Bump the plugin version**

In `cli-development/.claude-plugin/plugin.json` change `"version": "1.0.0"` to `"version": "1.1.0"`.

- [ ] **Step 3: Bump the marketplace entry**

In `.claude-plugin/marketplace.json`, inside the object whose `"name"` is `"cli-development"`, change `"version": "1.0.0"` to `"version": "1.1.0"`. Do not touch any other entry.

- [ ] **Step 4: Confirm both files still parse and agree**

Run: `python3 -c 'import json; p=json.load(open("cli-development/.claude-plugin/plugin.json")); m=[x for x in json.load(open(".claude-plugin/marketplace.json"))["plugins"] if x["name"]=="cli-development"][0]; print(p["version"], m["version"])'`
Expected: `1.1.0 1.1.0`

- [ ] **Step 5: Commit**

```bash
git add cli-development/CHANGELOG.md cli-development/.claude-plugin/plugin.json .claude-plugin/marketplace.json
git commit -m "bump cli-development to 1.1.0 with changelog"
```

---

### Task 8: Final GREEN run and plugin validation

**Agent:** harness-and-close-out (sonnet)

**Subagents:** `plugin-dev:plugin-validator` for the structural validation in Step 3.

**Files:** none modified unless a check fails.

- [ ] **Step 1: Run the lint across all agents**

Run: `cli-development/tests/lint-agents.sh`
Expected: no `FAIL` lines, exit 0.

- [ ] **Step 2: Run the routing check**

Run: `cli-development/tests/route-check.py`
Expected: `0 failing must-pass probes`, exit 0. If a must-pass probe fails, report the row in your final report with the winner and confidence and stop; do not edit a description yourself, because each description has an owner in Tasks 2 through 4. The plugin owner will decide whether the scope statement or the probe's expectation is wrong.

- [ ] **Step 3: Validate the plugin**

Dispatch `plugin-dev:plugin-validator` with: "Validate the `cli-development` plugin at `cli-development/` against the marketplace manifest at `.claude-plugin/marketplace.json`." Fix any reported structural issue in the file it names, then rerun Step 1.

- [ ] **Step 4: Confirm the tree is clean and report**

Run: `git status --short`
Expected: empty. Report the routing table from Step 2 and the validator result in your final message.

---

## Self-Review

**Spec coverage.** Every row of the Decisions table maps to a task: descriptions, Output sections, and proactive examples are Tasks 2 through 4 and enforced by Task 1's lint; cli-developer scope and body are Task 2; the v1/v2 table, hybrid ownership, and terminal-background rewrite are Task 3; light terminals and Unicode fallback are Task 4; reference resync is Task 5; the four v1-naming pages are Task 6; changelog and version are Task 7; routing verification is Task 8.

**Placeholders.** None. Every agent file, script, and doc edit is given in full.

**Consistency.** The lint requires the literal strings `Use when`, `Do not use`, `(use <sibling>)`, `**Output:**`, and `Proactive trigger:`; all three agent files in Tasks 2 through 4 contain each. The routing check's expected winners match the Routing Rules written into the reference page in Task 5. The version table in Task 3 and the reduced table in Task 5 use the same identifiers.
