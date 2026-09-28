# desktop-development Wails v3 Conversion Design

## Overview

**Plugin:** `desktop-development`
**Version:** 1.0.0 → 2.0.0
**Scope:** Replace the Electron+Go agent with a Wails v3 agent, add a `wails-v3` reference skill, add structural and routing tests, and rewrite the five Diátaxis pages. Electron is dropped entirely.

Inputs: the current plugin, the working Wails v3 mock-up at `~/Projects/Praetor/tactical-bridge-wails` (beta.25, React + Vite + TypeScript + Tailwind, Taskfile builds, `server` build tag), its rebuild analysis (`docs/analysis/2026-09-24-wails-rebuild-analysis.md`), and the official docs at https://v3.wails.io.

### Decisions made during design

| Decision | Choice |
| --- | --- |
| Plugin shape | One agent plus one reference skill |
| Electron | Dropped entirely; no migration trigger, no legacy agent |
| Platforms | macOS first with full defaults; Windows and Linux at Taskfile level only; no mobile |
| Frontend default | React + Vite + TypeScript + Tailwind via the `react` template; others on request |
| Beta drift | Pin one tag in the skill and in projects; read v3.wails.io before writing any API call; live docs win over the skill |
| Verification | agent-modernizer audit and routing with no Should-fix, live `wails3 init` scaffold check, TypeSafe analysis loop on the skill, Diátaxis validator, plugin validator |
| Plugin name and directory | Unchanged (`desktop-development`) |
| Branch | `wails-v3-conversion` off master, commits per phase, explicit-path staging |

## Non-goals

- No Electron support of any kind after this change.
- No iOS or Android guidance even though the mock-up scaffolds both.
- No signed or notarized build in this session; signing claims are written from the docs and marked as unverified where a signed bundle would be needed to prove them.
- No changes to other plugins. The tests read sibling agents' descriptions read-only.

## Structure after the change

```
desktop-development/
├── .claude-plugin/plugin.json          # 2.0.0; keywords wails, go, desktop, macos, native, webview
├── CHANGELOG.md                        # new; 2.0.0 entry names the verified Wails tag
├── agents/
│   └── wails-go-developer.md           # replaces electron-go-pro (deleted)
├── skills/
│   └── wails-v3/
│       ├── SKILL.md
│       └── references/
│           ├── project-layout.md
│           ├── native-features.md
│           ├── build-and-sign.md
│           └── beta-gotchas.md
├── tests/
│   ├── lint-agents.sh
│   └── route-check.py
└── docs/
    ├── _index.md
    ├── tutorials/getting-started.md
    ├── howto/set-up-wails-go-project.md   # replaces set-up-electron-go-project.md
    ├── reference/agents.md
    ├── reference/skills.md                # new
    └── explanation/architecture.md
```

The marketplace manifest entry for `desktop-development` gets the new version, description, and keywords.

## The agent: `wails-go-developer`

A builder that owns the whole Wails app: Go services, the webview frontend, native integration, and packaging. It delegates deep backend design to `go-architect`, React work to `react-specialist`, and end-to-end tests to `playwright-expert`, all of which exist in the marketplace.

### Frontmatter

| Field | Value |
| --- | --- |
| `name` | `wails-go-developer` |
| `model` | `opus` |
| `color` | `cyan` |
| `tools` | Read, Write, Edit, Bash, Grep, Glob, WebFetch |

### Description

Under 120 words, selection criteria with handoffs:

> Use when the user asks to build, restructure, or package a desktop app with Wails v3 and Go: services and generated bindings, windows and menus, system tray, notifications, global shortcuts, the dev loop, signing and DMG packaging. Also use when an existing Go program needs a desktop window. Do not use for terminal-only tools (use cli-developer), interactive terminal UIs (use go-tui-developer), or web apps with no desktop shell (use frontend-developer).

### Body (target under 3,000 characters excluding examples)

**Role.** One sentence: a senior desktop developer who builds macOS-first Wails v3 applications in Go with a React frontend, and ships them signed.

**Architecture defaults.**

- One Go process. Wails-generated TypeScript bindings are the only contract between Go and the UI. No HTTP API, no hand-written client, no duplicated types.
- Every entity is defined once in Go under `internal/domain`, which imports neither Wails nor any store.
- `app/` holds thin Wails adapters, one file per native concern (window, menu, tray, notifications, shortcuts), so beta API churn stays local.
- Nothing under `internal/` imports Wails, so a CLI or TUI entry point under `cmd/` reuses it unchanged.
- Go pushes to the UI with typed events; no polling timers. Streams carry token-by-token output.
- Local state, when needed, lives in SQLite through the Wails sqlite service; the UI never touches storage.
- OAuth runs in the system browser with a loopback listener. The webview has no cookie support and a custom origin, so it is never the OAuth surface.
- Every native feature has a `server` build tag stub so `wails3 build -tags server` produces a browser-served build for tests and demos.

**Process.**

1. Pin the Wails tag from the `wails-v3` skill in `go.mod` and confirm the `wails3` CLI matches and Go is 1.25 or newer.
2. Read the skill references for the features in play, then the matching v3.wails.io page for any API about to be written.
3. Scaffold with `wails3 init -t react`; never hand-roll the tree.
4. Put services in `app/`, domain types in `internal/domain`, one adapter file per native feature with its `server` stub.
5. Run `task dev` once and `task build` before reporting; run `wails3 generate bindings` after any service signature change.

**Output.** What was built, the pinned tag, the Taskfile targets that ran, and any native feature left unverified because it needs a signed bundle.

**Do Not.**

- Hand-roll an HTTP or WebSocket API between Go and the UI.
- Put business logic in `app/` adapters.
- Use `latest` or a floating version for the Wails module or CLI.
- Promise Touch Bar, the macOS Services menu, or Spotlight indexing; Wails has none of them.
- Claim notifications or the custom URL scheme work from an unsigned or unpackaged build.

### Examples (in the body, four)

1. New app: "I want a macOS desktop app in Go for tracking incidents" → scaffold, services, bindings, events.
2. Native features: "Add a menu bar icon and native notifications to our Wails app" → tray adapter, notifications adapter, signed-bundle caveat.
3. Packaging: "Package this for distribution with signing and a DMG" → universal binary, codesign, notarytool, DMG task.
4. Proactive: user with a Go CLI asks for "a small window that shows the status" → agent proposes Wails with the CLI kept under `cmd/`, marked `Proactive trigger:` in the commentary.

## The skill: `wails-v3`

A reference skill. Description triggers on Wails v3 questions from any agent or the user: bindings, services, events, streams, windows, menus, tray, notifications, shortcuts, Taskfile builds, signing, the updater, and "does Wails v3 support X".

### SKILL.md

Short. Three parts:

1. **Pin policy.** The verified tag (`v3.0.0-beta.25`, the mock-up's tag, unless the scaffold check in this session shows the installed CLI is newer, in which case the newer tag is verified and recorded). Projects pin `go.mod` and the CLI to one tag and upgrade deliberately. A bump re-runs the scaffold check and updates every "Verified against" line.
2. **Doc-reading rule.** The skill is the map; v3.wails.io is the truth. Read the matching page before writing any API call. If they disagree, the page wins and the skill gets a fix.
3. **Reference table.** Which file to open for which need.

### References

Each file opens with `Verified against: <tag>, <date>` and lists the doc pages it was written from.

**`project-layout.md`**: the tree `wails3 init -t react` produces; what `main.go`, `build/config.yml`, `build/Taskfile.yml`, and the per-platform Taskfiles own; the `task` targets (`dev`, `build`, `package`, `run`, plus a `bindings` target to add); the `//go:embed all:frontend/dist` asset pattern; the `server` build tag stub pattern with a two-file example; where `cmd/` entry points go; what `dev_mode` in `build/config.yml` watches and why `*_test.go` is ignored.

**`native-features.md`**: one section per feature with the Go entry point, the option or event names that matter, and the constraint. Windows (`WebviewWindowOptions`, hidden-inset title bar, frameless, multi-window with per-window hash routes); menus (roles, why the app menu is hand-built to carry Preferences); system tray (template icon, attached window, dynamic menu rebuild with debounce); notifications (actions, authorization, signed-bundle requirement); global and in-app shortcuts; dock badge; events (application, window, custom typed with `RegisterEvent`); streams; dialogs; single instance; custom URL scheme via `protocols` in `build/config.yml`.

**`build-and-sign.md`**: universal binary, `codesign` with hardened runtime and entitlements, `notarytool` submit and `stapler`, DMG, the updater and its manifest, and the Windows (NSIS, MSIX) and Linux (AppImage, deb) paths at Taskfile-target level only.

**`beta-gotchas.md`**: the verified-against list. Each entry: behavior, consequence, workaround. Contents: no cookies and a custom `wails://` origin (OAuth in the system browser); notifications need a signed bundle to authorize; no window-state persistence (use `Bounds` plus a store); Go 1.25 minimum; breaking changes ship inside the beta series (pin); `localStorage` under the custom scheme must be tested per app; no Touch Bar, Services menu, or CoreSpotlight; dev-mode reliability is an open GA blocker so `task build` is the acceptance step, not `task dev`.

### Left out on purpose

React, Vite, Tailwind, and Go idioms the model already knows; mobile; Electron.

## Docs

All five pages rewritten by the Diátaxis agents from the finished agent and skill, plus a new `reference/skills.md`.

- **Tutorial**: the note-taking app again, now with `wails3 init`, one service, generated bindings, one typed event, and a `task dev` run. Checkpoints kept.
- **How-to** `set-up-wails-go-project.md`: scaffold, add a service, wire an event, add tray and notifications with `server` stubs, sign and package. Troubleshooting covers the gotchas that bite first: Go version, notifications unauthorized, OAuth in the webview, stale bindings.
- **Reference**: `agents.md` for `wails-go-developer`; `skills.md` for `wails-v3` (trigger phrases, reference files, verified tag).
- **Explanation**: why one process and generated bindings instead of two processes and IPC; what the mock-up analysis measured (283 MB Electron bundle against a 10 MB Wails app; about 80 percent of the old code was boundary plumbing or dead); why macOS first; the beta trade-off and the pin policy.
- **Index**: components table and links.

## Tests

`desktop-development/tests/`, same shape as cli-development:

- **`lint-agents.sh`**: description has "Use when" and a "Do not use for … (use `<sibling>`)" handoff; body has an `**Output:**` line; 2 to 5 examples with one marked `Proactive trigger:`; body under 3,000 characters excluding examples; no first person.
- **`route-check.py`**: one Jev Choice per probe over the descriptions of `wails-go-developer`, `cli-developer`, `go-tui-developer`, `frontend-developer`, `react-specialist`, and `go-architect`, plus `none`. Probe sets: positives (new Wails app, add a tray, sign a DMG, "give my Go tool a window"); siblings that must route elsewhere (a Bubble Tea TUI, a plain Cobra CLI, a Next.js site, a Go API server); Electron negatives that must route to `none`. Threshold 0.70, exit 1 on any must-pass failure, exit 2 without a key.

## Verification

In order, all before the branch is offered for merge:

1. **Scaffold check.** `wails3 init -t react` in the scratchpad at the installed CLI's tag, `task build` once. Every path, target, config key, and file name in the skill is checked against that output. A claim that does not hold is corrected or removed. The tag that was built becomes the verified tag.
2. **Agent audit.** `ai-development/skills/agent-modernizer/scripts/audit-agents.py` on the new agent: no Must-fix, no Should-fix. `--route desktop-development/agents` and `tests/route-check.py` both pass.
3. **TypeSafe loop on the skill.** A Jev question set over `SKILL.md` and each reference, run before and after edits with findings confirmed by reading: does each section state a decision or constraint rather than a topic list; is every claim marked verified traceable to a named doc page or the scaffold output; does the skill description read as selection criteria; does any reference teach general React, Go, or shell knowledge. The question set lives in the plan, not in the plugin.
4. **Docs.** Diátaxis cross-link validator over `desktop-development/docs`.
5. **Plugin.** `plugin-validator` on `desktop-development`.

## Risks

- **The installed CLI may be newer than beta.25.** The scaffold check resolves this: whatever tag builds is the verified tag, and the mock-up's patterns are re-checked against it.
- **Docs and beta may disagree.** The doc-reading rule and the "live docs win" clause put the fix on the skill, not the project.
- **Signing claims cannot be proven here.** They are written from the signing guide and marked unverified in `build-and-sign.md` where a signed bundle would be needed.
