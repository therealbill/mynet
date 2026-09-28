# Changelog

All notable changes to desktop-development will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2026-09-28

Verified against Wails `v3.0.0-beta.25` with a live `wails3 init -t react` scaffold and `task build` on macOS.

### Added

- `wails-go-developer` agent: macOS-first Wails v3 applications in Go with a React frontend, one process, generated bindings as the only contract, thin adapters under `app/`, and a `server` build-tag stub for every native feature
- `wails-v3` skill with four references: project layout and Taskfile contract, native features, build and signing, and beta gotchas, each with a "Verified against" block
- `tests/lint-agents.sh` structural check and `tests/route-check.py` Jev routing simulation with Electron and Swift negatives
- `docs/reference/skills.md` and `docs/howto/set-up-wails-go-project.md`

### Changed

- Tutorial, reference, explanation, and index pages rewritten for Wails v3
- Plugin description and keywords now name Wails, Go, and the webview

### Removed

- `electron-go-pro` agent and the Electron+Go how-to; Electron is no longer supported by this plugin

## [1.0.0] - 2026-02-08

### Added

- `electron-go-pro` agent
- Diátaxis documentation set
