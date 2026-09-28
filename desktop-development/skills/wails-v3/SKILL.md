---
name: wails-v3
description: >
  This skill should be used when working with Wails v3 in Go: running "wails3 init", writing
  services and generated bindings, typed events, streams, windows and menus, system tray,
  notifications, global shortcuts, dock badges, Taskfile builds ("task dev", "task build",
  "task package"), signing and notarization, DMG packaging, the updater, server mode, or
  questions like "does Wails v3 support X", "why won't my notification show", "how do I sign
  and notarize my Wails app", or "wails3 dev isn't picking up my changes". Also applies when a
  Wails project's pinned beta tag must be checked or bumped. Do not use for general
  React/Vite/Tailwind work with no Wails API involved, or iOS/Android builds.
---

# Wails v3

Reference material for Wails `v3.0.0-beta.25` that the model cannot infer: the scaffold layout,
the Taskfile contract, native feature entry points, and the beta edges. The skill is the map;
the sources below are the truth.

## Pin policy

- Verified tag: `v3.0.0-beta.25` (2026-09-28). Projects pin `go.mod` and the `wails3` CLI to the
  same tag and upgrade deliberately, never with `latest`.
- Wails ships near-daily betas and has shipped breaking changes inside the beta series (see
  `references/beta-gotchas.md`). A bump re-runs the scaffold check (`wails3 init -t react`,
  `task build`) and updates every `Verified against` block in `references/`.
- Go 1.25 or newer is required by the module (`go.mod`; see `references/project-layout.md`).

## Where the truth is

Precedence when sources disagree: the scaffold and module source at the tag, then the docs at the
tag, then this skill. Before writing any API call, read the matching page or source file:

- Module source: `$(go env GOMODCACHE)/github.com/wailsapp/wails/v3@v3.0.0-beta.25/` (`pkg/application`, `pkg/services`, `pkg/updater`, `examples/`).
- Docs at the tag: `https://github.com/wailsapp/wails/tree/v3.0.0-beta.25/docs/mpress/content` (sparse-clone it; the live site at v3.wails.io blocks scripted fetches with HTTP 403, and its pages lag or lead the tag — see `references/beta-gotchas.md`).
- If a page and the source disagree, follow the source and note the discrepancy in
  `references/beta-gotchas.md`.

## References

| Need | Open |
| --- | --- |
| What `wails3 init -t react` creates, Taskfile targets, `build/config.yml`, the `server` build-tag stub pattern, where a CLI entry point goes | `references/project-layout.md` |
| Windows, menus, tray, notifications, shortcuts, dock, events, streams, dialogs, single instance, URL schemes | `references/native-features.md` |
| Universal binary, entitlements, signing, notarization, DMG, updater, Windows and Linux packaging | `references/build-and-sign.md` |
| What does not work yet, what needs a signed bundle, docs-versus-tag discrepancies | `references/beta-gotchas.md` |

## Not covered

React, Vite, Tailwind, and Go idioms the model already knows; iOS and Android.
