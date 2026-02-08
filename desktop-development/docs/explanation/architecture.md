---
title: "Architecture"
description: "Why Electron+Go for desktop applications"
weight: 1
---

# Architecture

The desktop-development plugin centers on a specific architectural bet: Electron for the UI, Go for the backend, running as separate OS processes on macOS. This page explains why that combination exists, where the boundaries are, and when to choose something else.

## Why Hybrid Electron+Go

Electron provides cross-platform UI built on web technologies -- HTML, CSS, JavaScript, React. Developers can reuse web skills and the massive npm ecosystem to build rich interfaces quickly. However, Electron's Node.js runtime is not well suited for CPU-intensive work, filesystem-heavy operations, or anything requiring low-level system access.

Go fills that gap. It compiles to a single static binary, starts in milliseconds, has excellent concurrency primitives, and integrates cleanly with C libraries and system APIs. Go handles the work that would be slow or awkward in Node.js: database access, file processing, background tasks, system integration.

The combination gives web-tech UI flexibility with systems-language performance. Neither layer compromises for the other.

## The IPC Boundary

Go runs as a separate operating system process. Electron communicates with it over a local network protocol -- HTTP for request/response, WebSocket for real-time push, or gRPC if the project already uses protobuf.

This separation is deliberate. Two processes with a network boundary between them can be:

- **Developed independently.** The Go API has a contract. The Electron UI consumes that contract. Either side can change its internals without breaking the other.
- **Tested independently.** The Go API is testable with standard HTTP test tools. The Electron UI is testable with a mock API server. Integration tests cover the boundary.
- **Debugged independently.** A bug in the Go layer does not require understanding the Electron renderer. A UI glitch does not require reading Go code.

The Electron main process owns the Go process lifecycle: it spawns the Go binary on startup, monitors its health on a 5-second interval, and terminates it on quit. The Go binary binds to `127.0.0.1` with a random port, preventing external network access.

## Why macOS Focus

electron-go-pro specializes in macOS for a practical reason: macOS has the most complex signing and distribution requirements of any desktop platform.

A hybrid Electron+Go app on macOS requires:

- Code signing the Go binary with a Developer ID certificate
- Code signing the Electron `.app` bundle (which includes the Go binary)
- Notarizing the signed bundle with Apple's notary service
- Stapling the notarization ticket to the `.app`
- Building a universal binary that runs on both Intel and Apple Silicon
- Respecting macOS sandboxing and entitlements for App Store distribution

Skipping any of these steps results in Gatekeeper blocking the app or Apple rejecting the submission. electron-go-pro encodes this knowledge so developers do not have to rediscover it through trial and error.

The Electron+Go architecture itself is not macOS-specific. The same IPC pattern, project structure, and build pipeline work on Windows and Linux. The macOS focus is about distribution complexity, not architectural limitation.

## When NOT to Use This Plugin

Not every application needs a graphical desktop shell.

- **Terminal-only applications.** If the app runs in a terminal and does not need windows, menus, or graphical UI, use the [cli-development](/cli-development/) plugin instead. It provides agents for CLI argument parsing, TUI frameworks like Bubble Tea, and terminal UI design. Adding Electron to a terminal app introduces unnecessary complexity.
- **Web applications.** If the app runs in a browser and does not need desktop integration (file system access, system tray, native menus), use the [web-development](/web-development/) plugin. Wrapping a web app in Electron just for the sake of a desktop icon is rarely worth the packaging and update overhead.
- **Lightweight utilities.** If the app is a simple menubar utility or system tray icon with no complex UI, a native Swift/AppKit approach may be simpler than the Electron+Go stack. The [mobile-development](/mobile-development/) plugin includes Swift expertise that applies to macOS native development.

## Cross-Plugin Collaboration

electron-go-pro does not work in isolation. It acts as an architect that coordinates specialists from other plugins:

- **go-architect** (backend-development plugin) handles all Go backend work -- API design, database operations, business logic. electron-go-pro defines the IPC contract; go-architect implements the Go side.
- **frontend-developer** and **react-specialist** (web-development plugin) handle the Electron renderer UI -- React components, state management, styling. electron-go-pro defines the security constraints; the frontend agents work within them.
- **typescript-pro** (programming-languages plugin) ensures type-safe interfaces between Go API responses and the TypeScript frontend.
- **playwright-expert** (code-quality plugin) handles end-to-end testing of the full hybrid stack.

This delegation pattern reflects how the Mynet marketplace works: plugins collaborate through Claude's dispatch layer rather than importing each other's code. electron-go-pro knows what work to delegate and to whom, but the specialist agents do not need to know they are operating inside a desktop application context.

## See Also

- [Agent Reference](../../reference/agents/) -- electron-go-pro specification
- [Set Up an Electron+Go Project](../../howto/set-up-electron-go-project/) -- practical setup steps
- [Getting Started](../../tutorials/getting-started/) -- build your first hybrid app
