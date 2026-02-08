---
title: "Agents"
description: "Desktop development agent specifications"
weight: 1
---

# Agents

Agent specifications for the desktop-development plugin.

## electron-go-pro

Expert Electron+Go hybrid desktop application architect for macOS. Coordinates Go backend services with Electron UI, handles IPC design, native macOS integration, code signing, and packaging.

### Specification

| Field | Value |
|-------|-------|
| Name | electron-go-pro |
| Model | opus |
| Color | cyan |
| Tools | Read, Write, Edit, Bash, Grep, Glob, Task, WebSearch, WebFetch |

### Trigger Conditions

electron-go-pro activates when the user mentions:

- Building desktop applications with Electron and Go
- IPC design between Go and Electron processes
- Packaging or distributing hybrid Electron+Go applications
- Code signing and notarization for macOS
- Coordinating a Go backend with an Electron frontend

### Capabilities

| Capability | Description |
|------------|-------------|
| Project architecture | Proposes and scaffolds hybrid Electron+Go project structures |
| IPC layer design | Configures HTTP, WebSocket, or gRPC communication between Go and Electron |
| Native macOS integration | Menu bar, dock badges, notification center, dark/light mode, universal binary |
| Code signing configuration | Signs both Go binary and Electron `.app` bundle, configures notarization |
| Electron packaging | Embeds Go binary in Electron app resources, builds with electron-builder or electron-forge |
| Delegation to specialist agents | Routes backend work to go-architect, frontend work to frontend-developer or react-specialist |

### Delegation Pattern

electron-go-pro delegates domain-specific work to specialist agents from other plugins:

| Delegate | Plugin | Responsibility |
|----------|--------|---------------|
| go-architect | backend-development | Go backend code: API design, database operations, business logic, system integration |
| frontend-developer | web-development | Electron renderer UI: React components, state management, real-time updates |
| react-specialist | web-development | React-specific patterns within the Electron renderer |
| typescript-pro | programming-languages | Type-safe interfaces between Go API responses and the TypeScript frontend |
| playwright-expert | code-quality | End-to-end testing of the full hybrid stack |

### Security Defaults

electron-go-pro enforces these Electron security settings on every project:

| Setting | Value |
|---------|-------|
| Context isolation | Enabled |
| Node integration | Disabled in all renderers |
| Preload scripts | Required as the only bridge to Node/IPC |
| Content Security Policy | No `unsafe-inline`, no `unsafe-eval` |
| Web security | Enabled |
| Remote module | Disabled |

### IPC Defaults

| Protocol | Use Case |
|----------|----------|
| HTTP | Request/response communication (default) |
| WebSocket | Real-time push from Go to Electron |
| gRPC | Only when the project already uses protobuf |

IPC binds to `127.0.0.1` with a random available port. The Go binary writes the port to stdout. The Electron main process reads the port and passes it to the renderer via the preload bridge.

### Build Pipeline

1. Go binary: cross-compile for `darwin/arm64` and `darwin/amd64`, combine with `lipo` into a universal binary
2. Electron shell: package with electron-builder or electron-forge
3. Embed Go binary in the Electron app's resources directory
4. Code sign the Go binary and the Electron `.app` bundle separately
5. Notarize the final `.app` with `notarytool`
6. Auto-update delivers the combined package; Electron handles updates for both binaries

### Example Interactions

```
User: "I want to build a macOS desktop app with an Electron frontend and Go backend"
Agent: Activates, asks about app requirements, proposes architecture with IPC strategy

User: "How should the Go backend communicate with the Electron renderer?"
Agent: Designs IPC architecture, recommends HTTP or WebSocket based on requirements

User: "Package this app for macOS with code signing and notarization"
Agent: Configures dual code signing, notarization, and the hybrid build pipeline
```

## See Also

- [Architecture]({{< ref "explanation/architecture" >}}) -- design decisions behind the Electron+Go approach
- [Set Up an Electron+Go Project]({{< ref "howto/set-up-electron-go-project" >}}) -- step-by-step project setup
