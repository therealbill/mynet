---
title: "Getting Started with Desktop Development"
description: "Build your first Electron+Go desktop application with the electron-go-pro agent"
weight: 1
---

# Getting Started with Desktop Development

Build a note-taking desktop application for macOS using an Electron frontend and Go backend, guided by the electron-go-pro agent.

## What You'll Build

By the end of this tutorial, you will have:

- Triggered the electron-go-pro agent with a project description
- Reviewed an architecture proposal for a hybrid Electron+Go app
- Scaffolded a project with IPC communication between Go and Electron
- Built and launched a working desktop application on macOS

This takes about 20 minutes.

## Prerequisites

- Claude Code CLI installed and authenticated (run `claude --version` to verify)
- The desktop-development plugin installed in your project's `.claude/settings.json`
- Go 1.21+ installed (`go version` to verify)
- Node.js 18+ installed (`node --version` to verify)
- Xcode Command Line Tools installed (`xcode-select --install` if needed)

## Step 1: Start a Conversation

Open Claude Code in your project directory and describe what you want to build:

```
I want to build a macOS desktop app with an Electron frontend and Go backend.
It's a note-taking app with local storage.
```

Claude Code matches your request to the electron-go-pro agent based on the mention of Electron, Go, and desktop application development. The agent activates and begins coordinating the project.

## Step 2: Answer the Agent's Questions

electron-go-pro asks clarifying questions about your application requirements before proposing an architecture. Expect questions like:

- What data does the app store? (Notes with title, body, timestamps)
- Does it need real-time features? (No -- simple CRUD is sufficient)
- What about search or filtering? (Full-text search on note content)
- Any specific macOS integrations? (Dark mode support, native menu bar)

Answer these questions. The agent uses your responses to make architectural decisions about the IPC layer, database schema, and build pipeline.

## Step 3: Review the Architecture Proposal

electron-go-pro proposes a project architecture. For a note-taking app, expect something like:

- **Go backend:** HTTP API on localhost with a random port, SQLite database for note storage, full-text search via SQLite FTS5
- **Electron frontend:** React renderer with context isolation enabled, preload script bridging IPC to the Go API
- **Process lifecycle:** Electron main process spawns the Go binary as a child process, health-checks it on a 5-second interval, and tears it down on quit

The agent explains why it chose HTTP over WebSocket for this app -- request/response fits CRUD operations, and there is no real-time push requirement.

Review the proposal. Ask questions or request changes before proceeding.

### Checkpoint

At this point you should have:

- Triggered electron-go-pro with your project description
- Answered the agent's clarifying questions
- Received and reviewed an architecture proposal

If electron-go-pro did not activate, verify the desktop-development plugin is listed in your `.claude/settings.json` and restart Claude Code.

## Step 4: Scaffold the Project

Tell the agent to proceed:

```
That architecture looks good. Scaffold the project.
```

electron-go-pro generates the project structure. It delegates backend Go work to go-architect and frontend React/Electron work to frontend-developer or react-specialist. You see files created across both layers:

- `backend/main.go` -- Go entry point with HTTP server and SQLite setup
- `backend/api/` -- Route handlers for notes CRUD and search
- `backend/db/` -- SQLite schema migrations and query functions
- `frontend/src/main/` -- Electron main process with Go child process management
- `frontend/src/preload/` -- Preload script exposing IPC bridge
- `frontend/src/renderer/` -- React components for the note-taking UI
- `package.json` -- Electron and frontend dependencies
- `go.mod` -- Go module definition
- `Makefile` -- Build targets for both Go and Electron

## Step 5: Verify the Generated Files

Check that the key files exist and contain the expected structure:

```
Show me the contents of backend/main.go and frontend/src/main/index.ts
```

Verify:

- `backend/main.go` starts an HTTP server on a random localhost port and writes the port number to stdout
- `frontend/src/main/index.ts` reads the port from the Go process's stdout and passes it to the renderer via the preload bridge
- The preload script exposes a typed API object, not raw Node.js access

### Checkpoint

At this point you should have:

- A complete project structure with Go backend and Electron frontend
- IPC communication configured over HTTP on localhost
- A Makefile with build targets for both layers

If any files are missing, ask electron-go-pro to regenerate the specific component.

## Step 6: Build and Run

Ask the agent to build and launch the app:

```
Build and run the app
```

The agent executes the build pipeline:

1. Compiles the Go binary with `go build`
2. Installs Node dependencies with `npm install`
3. Launches Electron in development mode, which spawns the Go binary

A window appears with the note-taking interface. Create a note, edit it, and verify it persists when you close and reopen the app.

### Checkpoint

At this point you should have:

- A running desktop application on macOS
- The ability to create and retrieve notes through the UI
- Confirmed that data persists in the local SQLite database

## What You Learned

In this tutorial, you:

- **Triggered electron-go-pro** by describing a desktop application with Electron and Go -- the agent activated based on those keywords
- **Reviewed an architecture proposal** that separated concerns between the Go backend (data, business logic) and Electron frontend (UI presentation)
- **Observed agent delegation** as electron-go-pro coordinated with go-architect for backend code and frontend-developer for React components
- **Verified IPC communication** between Go and Electron over HTTP on localhost
- **Built and ran** a hybrid desktop application from a single build pipeline

## Next Steps

- [Set Up an Electron+Go Project]({{< ref "howto/set-up-electron-go-project" >}}) -- configure IPC, code signing, and packaging for distribution
- [Architecture]({{< ref "explanation/architecture" >}}) -- understand why Electron+Go is a strong combination for desktop apps
- [Agent Reference]({{< ref "reference/agents" >}}) -- full specification of the electron-go-pro agent
