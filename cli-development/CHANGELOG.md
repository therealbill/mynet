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
- `go-tui-developer` terminal-background guidance rewritten around v2 `View.BackgroundColor` and v1 termenv, with AppleScript kept only as the macOS Terminal.app route
- Generic process and convention lines trimmed from `cli-developer`
- Reference, how-to, tutorial, and explanation pages updated to name both Charm majors
- Code review corrected color-override precedence, SIGINT and SIGPIPE behavior, XDG fallback rules, the locale rule for Unicode fallbacks, and the Charm v1 and v2 terminal-background facts

## [1.0.0] - 2026-02-08

### Added

- `cli-developer`, `go-tui-developer`, and `cli-ui-designer` agents
- Diátaxis documentation set
