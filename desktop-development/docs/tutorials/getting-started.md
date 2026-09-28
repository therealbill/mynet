---
title: "Getting Started with Desktop Development"
description: "Build your first Wails v3 desktop application with the wails-go-developer agent"
weight: 1
---

# Getting Started with Desktop Development

Build a note-taking desktop application for macOS using a Wails v3 frontend and Go backend, guided by the wails-go-developer agent.

## What You'll Build

By the end of this tutorial, you will have:

- Triggered the wails-go-developer agent with a project description
- Reviewed an architecture proposal for a Wails v3 app with generated Go-to-UI bindings
- Scaffolded a project with a Go service and typed events
- Built and launched a working desktop application on macOS

## Prerequisites

- Claude Code CLI installed and authenticated (run `claude --version` to verify)
- The desktop-development plugin installed in your project's `.claude/settings.json`
- The `wails3` CLI installed and pinned to `v3.0.0-beta.25` (`wails3 version` to verify)
- Go 1.25 or newer installed (`go version` to verify)
- Node.js 18+ installed (`node --version` to verify)
- Xcode Command Line Tools installed (`xcode-select --install` if needed)

## Step 1: Start a Conversation

Open Claude Code in your project directory and describe what you want to build:

```
I want to build a macOS desktop app with Wails v3 and Go.
It's a note-taking app with local storage.
```

Claude Code matches your request to the wails-go-developer agent based on the mention of Wails, Go, and desktop application development. The agent activates and begins coordinating the project.

## Step 2: Answer the Agent's Questions

wails-go-developer asks clarifying questions about your application requirements before proposing an architecture. Expect questions like:

- What data does the app store? (Notes with title, body, timestamps)
- Does it need real-time features? (No -- simple CRUD is sufficient, with updates pushed to the UI through a typed event)
- What about search or filtering? (Full-text search on note content)
- Any specific macOS integrations? (Dark mode support, native menu bar)

Answer these questions. The agent uses your responses to make architectural decisions about the Go service, the SQLite schema, and the build pipeline.

## Step 3: Review the Architecture Proposal

wails-go-developer proposes a project architecture. For a note-taking app, expect something like:

- **Go backend:** A `NoteService` under `app/` with methods that take `context.Context` first and return `error`, so the frontend gets a cancellable promise and a rejected promise on failure. The `Note` type is defined once under `internal/domain`. Notes persist in SQLite through the Wails sqlite service, with full-text search via SQLite FTS5.
- **React frontend:** `frontend/src/App.tsx` calls the Wails-generated TypeScript bindings directly -- no hand-written client, no HTTP API, no IPC bridge.
- **Events:** A typed `note:saved` event, registered with `application.RegisterEvent[Note]("note:saved")` in `init()`, fires whenever a note is saved so the UI updates without polling.

The agent explains why it chose generated bindings over a hand-rolled API -- one Go process and one typed contract eliminate the boundary plumbing a separate HTTP layer would need.

Review the proposal. Ask questions or request changes before proceeding.

### Checkpoint

At this point you should have:

- Triggered wails-go-developer with your project description
- Answered the agent's clarifying questions
- Received and reviewed an architecture proposal

If wails-go-developer did not activate, verify the desktop-development plugin is listed in your `.claude/settings.json` and restart Claude Code.

## Step 4: Scaffold the Project

Tell the agent to proceed:

```
That architecture looks good. Scaffold the project.
```

wails-go-developer confirms the `wails3` CLI is pinned to `v3.0.0-beta.25` and Go is 1.25 or newer, then scaffolds the project:

```
wails3 init -n notes -t react -d notes -mod example.com/notes
```

It deletes the scaffold's example `greetservice.go`, and delegates domain modeling to go-architect and the React UI to react-specialist or frontend-developer. You see files created across both layers:

- `notes/main.go` -- `application.New`, registers `NoteService`, and calls `application.RegisterEvent[Note]("note:saved")` in `init()`
- `notes/app/note_service.go` -- `NoteService` with CRUD methods, each taking `context.Context` first and returning `error`
- `notes/internal/domain/note.go` -- the `Note` type and validation rules
- `notes/internal/store/` -- SQLite persistence behind an interface
- `notes/frontend/src/App.tsx` -- React UI for the note-taking interface
- `notes/frontend/bindings/` -- generated TypeScript bindings, committed to the repo
- `notes/go.mod` -- pinned to `github.com/wailsapp/wails/v3 v3.0.0-beta.25`
- `notes/Taskfile.yml` -- `dev`, `build`, and `package` targets
- `notes/build/config.yml` -- product info and the dev-mode file watcher

## Step 5: Verify the Generated Files

Check that the key files exist and contain the expected structure:

```
Show me the contents of app/note_service.go and frontend/src/App.tsx
```

Verify:

- `note_service.go`'s methods take `context.Context` as their first parameter and return `error`
- `main.go`'s `init()` calls `application.RegisterEvent[Note]("note:saved")`
- `App.tsx` imports and calls the generated binding directly, such as `import { SaveNote } from "../bindings/notes/app/noteservice"` -- no fetch calls, no preload script, no IPC bridge

### Checkpoint

At this point you should have:

- A complete project structure with a Go backend and a React frontend in one process
- Generated TypeScript bindings as the only contract between Go and the UI
- A `note:saved` typed event registered in `init()`
- A `Taskfile.yml` with `dev` and `build` targets

If any files are missing, ask wails-go-developer to regenerate the specific component.

## Step 6: Build and Run

Ask the agent to build and launch the app:

```
Build and run the app
```

The agent runs the dev loop first, then the production build:

1. Runs `task dev`, which rebuilds the Go binary on change and gives the frontend Vite hot reload
2. A window appears with the note-taking interface. Create a note, edit it, and watch the UI update through the `note:saved` event.
3. Stops the dev loop and runs `task build`, which produces the production binary at `bin/notes`

Verify a note persists by closing and reopening the app built at `bin/notes`.

### Checkpoint

At this point you should have:

- A running desktop application launched through `task dev`
- The ability to create and retrieve notes through the UI, with updates delivered by the `note:saved` event
- A production binary at `bin/notes` produced by `task build`

## What You Learned

In this tutorial, you:

- **Triggered wails-go-developer** by describing a desktop application with Wails v3 and Go -- the agent activated based on those keywords
- **Reviewed an architecture proposal** that used one Go process and generated bindings as the sole contract between the `NoteService` and the React UI
- **Observed agent delegation** as wails-go-developer coordinated with go-architect for domain modeling and react-specialist for the React UI
- **Verified the generated bindings and the typed `note:saved` event** instead of a hand-rolled IPC layer
- **Built and ran** the application through `task dev`, then produced a production binary with `task build`

## Next Steps

- [Set Up a Wails+Go Project](../../howto/set-up-wails-go-project/) -- configure native features, signing, and DMG packaging for distribution
- [Architecture](../../explanation/architecture/) -- understand why Wails v3 is a strong combination for desktop apps
- [Agent Reference](../../reference/agents/) -- full specification of the wails-go-developer agent
