---
title: "Set Up an Electron+Go Project"
description: "Configure IPC, build pipeline, and packaging for an Electron+Go hybrid app"
weight: 1
---

# Set Up an Electron+Go Project

Set up a working Electron+Go desktop application with IPC communication, a unified build pipeline, and macOS code signing.

## Prerequisites

- Claude Code with the desktop-development plugin installed
- Go 1.21+ and Node.js 18+ installed
- An Apple Developer account (for code signing)
- Xcode Command Line Tools installed

## Steps

### 1. Scaffold the Project

Ask electron-go-pro to create the project structure:

```
Create an Electron+Go desktop app with HTTP IPC and SQLite storage
```

The agent generates a Go backend, Electron frontend, IPC bridge, and Makefile. It delegates Go scaffolding to go-architect and Electron/React scaffolding to frontend-developer.

Verify the output includes:

- `backend/` directory with `main.go` and an HTTP API package
- `frontend/` directory with Electron main process, preload script, and React renderer
- `go.mod` and `package.json` at their respective roots
- A `Makefile` with `build`, `dev`, and `package` targets

### 2. Configure IPC Communication

electron-go-pro defaults to HTTP for request/response and WebSocket for real-time push. For most applications, HTTP is sufficient.

The IPC flow works as follows:

1. Electron main process spawns the Go binary as a child process
2. Go binary starts an HTTP server on `127.0.0.1` with a random available port
3. Go binary writes the port number to stdout
4. Electron main process reads the port and passes it to the renderer via the preload bridge
5. The renderer calls the Go API through the preload-exposed fetch wrapper

To switch to WebSocket IPC for real-time features, ask:

```
Add WebSocket IPC alongside the HTTP API for real-time updates
```

To use gRPC instead (only if the project already uses protobuf):

```
Replace HTTP IPC with gRPC using these proto definitions
```

### 3. Set Up the Electron Frontend

The Electron frontend follows strict security defaults:

- Context isolation enabled on all BrowserWindows
- Node integration disabled in all renderers
- A preload script serves as the only bridge between the renderer and Node/IPC
- Content Security Policy blocks `unsafe-inline` and `unsafe-eval`

To customize the renderer framework or add macOS-specific features:

```
Add native dark mode detection and a macOS-standard menu bar to the Electron shell
```

electron-go-pro delegates UI work to frontend-developer or react-specialist while maintaining the security boundary.

### 4. Configure the Build Pipeline

The Makefile coordinates both layers. Ask the agent to set it up:

```
Configure the build pipeline for both Go and Electron
```

The build pipeline performs these steps in order:

1. **Go binary:** Cross-compile with `GOOS=darwin GOARCH=arm64` and `GOARCH=amd64`, then combine into a universal binary with `lipo`
2. **Electron shell:** Install dependencies with `npm ci`, build the renderer with the configured bundler, package with electron-builder or electron-forge
3. **Embed:** Copy the Go universal binary into the Electron app's resources directory
4. **Result:** A single `.app` bundle containing both the Electron shell and the Go binary

Run `make build` to execute the full pipeline, or `make dev` for a development build that skips universal binary creation.

### 5. Set Up Code Signing for macOS

Code signing is required for distribution. Both the Go binary and the Electron shell must be signed separately. Ask the agent to configure signing:

```
Set up code signing and notarization for macOS distribution
```

electron-go-pro configures:

- **Go binary signing:** `codesign --sign "Developer ID Application: ..."` applied to the Go binary before embedding
- **Electron app signing:** electron-builder or electron-forge handles signing the `.app` bundle, including all embedded binaries
- **Notarization:** `notarytool` submits the signed `.app` to Apple for notarization, and `stapler` attaches the notarization ticket

Provide your Apple Developer Team ID and signing identity when the agent asks. The agent stores these in environment variables, not in source code.

### 6. Verification

Run the app and confirm IPC communication works end to end:

```bash
make dev
```

Verify:

- [ ] The Electron window opens without security warnings
- [ ] The Go backend starts and logs its listening port
- [ ] API requests from the renderer reach the Go backend and return responses
- [ ] The Go process terminates cleanly when the Electron window closes
- [ ] On `make build`, the resulting `.app` bundle launches without Gatekeeper warnings (if code signed)

## Troubleshooting

**Go binary fails to start:**

Check that the Go binary path in the Electron main process matches the actual location in the app resources. In development, this is typically `backend/` relative to the project root. In production, it is inside the `.app` bundle's `Resources` directory.

**IPC connection refused:**

The Go server binds to `127.0.0.1` with a random port. Ensure the Electron main process reads the port from the Go process's stdout before the renderer attempts to connect. A race condition here means the renderer tries to connect before Go is ready.

**Code signing errors:**

Ensure your Apple Developer certificate is installed in Keychain Access and not expired. Run `security find-identity -v -p codesigning` to list valid signing identities. The Go binary must be signed before being embedded in the Electron app.

## See Also

- [Getting Started](../../tutorials/getting-started/) -- tutorial walkthrough of building a complete app
- [Architecture](../../explanation/architecture/) -- why Electron+Go and how the IPC boundary works
- For terminal-only Go applications without a GUI, use the [cli-development](/cli-development/) plugin instead
