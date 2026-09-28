---
title: "Skills"
description: "Wails v3 skill specification"
weight: 2
---

# Skills

Skill specifications for the desktop-development plugin.

## wails-v3

Reference material for Wails `v3.0.0-beta.25` covering the scaffold layout, the Taskfile contract, native feature entry points, and the beta edges.

### Specification

| Field | Value |
|-------|-------|
| Name | wails-v3 |
| Verified tag | `v3.0.0-beta.25` (2026-09-28) |
| Reference files | 4 (`project-layout.md`, `native-features.md`, `build-and-sign.md`, `beta-gotchas.md`) |

### Trigger Conditions

The wails-v3 skill applies when working with Wails v3 in Go:

- Running `wails3 init`
- Writing services and generated bindings
- Typed events and streams
- Windows and menus
- System tray, notifications, global shortcuts, dock badges
- Taskfile builds (`task dev`, `task build`, `task package`)
- Signing and notarization, DMG packaging
- The updater
- Server mode
- Questions of the form "does Wails v3 support X", "why won't my notification show", "how do I sign and notarize my Wails app", or "wails3 dev isn't picking up my changes"
- Checking or bumping a Wails project's pinned beta tag

It does not apply to general React/Vite/Tailwind work with no Wails API involved, or to iOS/Android builds.

### Pin Policy

| Item | Value |
|------|-------|
| Verified tag | `v3.0.0-beta.25` |
| Verified date | 2026-09-28 |
| Pin targets | `go.mod` and the `wails3` CLI, pinned to the same tag |
| Upgrade discipline | Deliberate bump only, never `latest` |
| Bump procedure | Re-run the scaffold check (`wails3 init -t react`, `task build`) and update every `Verified against` block in `references/` |
| Go version | 1.25 or newer, required by the module |

### Source Precedence

When sources disagree, precedence is:

1. Scaffold and module source at the tag (`$(go env GOMODCACHE)/github.com/wailsapp/wails/v3@v3.0.0-beta.25/`: `pkg/application`, `pkg/services`, `pkg/updater`, `examples/`).
2. Docs at the tag (`https://github.com/wailsapp/wails/tree/v3.0.0-beta.25/docs/mpress/content`, sparse-cloned; the live site at v3.wails.io returns HTTP 403 to scripted fetches).
3. This skill.

If a docs page and the source disagree, the source wins and the discrepancy is noted in `references/beta-gotchas.md`.

### Reference Files

| File | Covers |
|------|--------|
| `references/project-layout.md` | What `wails3 init -t react` creates, Taskfile targets, `build/config.yml`, the `server` build-tag stub pattern, where a CLI entry point goes |
| `references/native-features.md` | Windows, menus, tray, notifications, shortcuts, dock, events, streams, dialogs, single instance, URL schemes |
| `references/build-and-sign.md` | Universal binary, entitlements, signing, notarization, DMG, updater, Windows and Linux packaging |
| `references/beta-gotchas.md` | What does not work yet, what needs a signed bundle, docs-versus-tag discrepancies |

### Beta Gotchas

Each entry in `references/beta-gotchas.md` is verified against `v3.0.0-beta.25` (2026-09-28) and names the behavior, the consequence, and the workaround.

| Gotcha | Section |
|--------|---------|
| Go 1.25 or newer | Version and drift |
| Breaking changes ship inside the beta series | Version and drift |
| The docs site blocks scripted reads | Version and drift |
| Docs and tag disagree on signing configuration | Version and drift |
| Docs place `config.yml` at the root; the scaffold puts it at `build/config.yml` | Version and drift |
| `TitleBar` appears as both a constant and a struct in the docs | Version and drift |
| There is no `react-ts` template | Version and drift |
| Custom `wails` scheme | Webview |
| Test `localStorage` and cookies per app | Webview |
| Safari feature set | Webview |
| Translucency needs a private-API build | Webview |
| Notifications need a packaged, signed bundle on macOS | Native features |
| URL schemes and single-instance URL relay only work from a packaged `.app` | Native features |
| No window-state persistence | Native features |
| Hide shortcut plus terminate-after-last-window quits the app | Native features |
| Global shortcut registered twice errors and keeps the first callback | Native features |
| Not available at all: Touch Bar, macOS Services menu, CoreSpotlight indexing, dock bounce, custom dock menu, `.pkg` output | Native features |
| Server mode drops tray, dialogs, windows, shortcuts, notifications | Native features |
| `task dev` reliability is an open GA blocker | Dev loop and build |
| `dev_mode.ignore.file` excludes `*_test.go` and the watcher does not rebuild on JSON changes | Dev loop and build |
| `common:update:build-assets` overwrites hand edits | Dev loop and build |
| `generate:bindings` runs `-clean=true` | Dev loop and build |
| `mac.LoadResource` and `mac.ResourceFS` return `ErrNotInAppBundle` | Dev loop and build |
| Cross-compiled macOS binaries are unsigned and ad-hoc signing is applied by `package` | Dev loop and build |

### Not Covered

React, Vite, Tailwind, and Go idioms; iOS and Android.

## See Also

- [Architecture](../../explanation/architecture/) -- design decisions behind the Wails v3 approach
- [Set Up a Wails + Go Project](../../howto/set-up-wails-go-project/) -- step-by-step project setup
- [Agents](../agents/) -- the wails-go-developer agent that consumes this skill
