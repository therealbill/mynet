---
title: "Set Up a Wails v3 Go Project"
description: "Scaffold, add services and events, wire native features, and package a Wails v3 app"
weight: 1
---

# Set Up a Wails v3 Go Project

**Goal**: Scaffold a Wails v3 application with a React frontend, add a Go service with generated
bindings, wire a typed event, add tray and notification adapters, and produce a signed,
notarized DMG.

## Prerequisites

- Go 1.25 or newer
- Node.js and npm
- The `wails3` CLI installed and pinned to the same tag as `go.mod` (`v3.0.0-beta.25` or later,
  matching your project)
- A macOS machine for the packaging and signing steps; a Developer ID Application certificate
  and an app-specific password for notarization
- Familiarity with Go, React, and the Wails v3 concepts covered in the `wails-v3` skill
  (services, bindings, typed events)

## Steps

### 1. Scaffold the project

```bash
wails3 init -n incidents -t react -d . -mod github.com/example/incidents
cd incidents
```

`-t react` is the TypeScript template; there is no separate `react-ts` template. `-d` with `-n`
nests the project at `<dir>/<name>`. The scaffold includes `main.go`, `go.mod` pinned to
`go 1.25.0`, the root `Taskfile.yml`, `build/config.yml`, and a `frontend/` tree with
`@wailsio/runtime`, React 18, and Vite. Delete the example `greetservice.go` once your own
service exists.

Never hand-build this tree yourself; the generated `Taskfile.yml` and `build/config.yml` are
what every later step in this guide depends on.

### 2. Add Tailwind and the Wails Vite plugin

The template does not include Tailwind or the Wails Vite plugin, and typed events need the
plugin to work on the frontend.

```bash
cd frontend
npm i -D tailwindcss postcss autoprefixer
npx tailwindcss init -p
cd ..
```

Edit `frontend/vite.config.ts` and add the Wails plugin alongside the React plugin:

```ts
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { wails } from "@wailsio/runtime/plugins/vite";

export default defineConfig({
  plugins: [react(), wails()],
});
```

### 3. Add a service and regenerate bindings

Services are plain Go structs under `app/`. Create `app/incidentservice.go`:

```go
package app

import "context"

type IncidentService struct{}

func (s *IncidentService) List(ctx context.Context) ([]Incident, error) {
	// domain logic belongs in internal/domain; keep this adapter thin
	return nil, nil
}
```

Register it in `main.go`:

```go
app := application.New(application.Options{
	Services: []application.Service{application.NewService(&app.IncidentService{})},
	// ...
})
```

Regenerate the TypeScript bindings after every service signature change:

```bash
wails3 generate bindings -ts -i -clean -d frontend/bindings
```

`-clean` deletes `frontend/bindings` before writing, so the command always leaves a fresh,
complete set of files. Add it as a root Taskfile shortcut so the team runs one command:

```yaml
  bindings:
    summary: Regenerate TypeScript bindings
    cmds:
      - wails3 generate bindings -ts -i -clean -d frontend/bindings
```

Expected result: `frontend/bindings/<module>/app/incidentservice.ts` exists and exports typed
functions for each exported method. Commit `frontend/bindings`; `.gitignore` does not exclude it.

### 4. Register a typed event

Register the event type in `init()` so both the Go and the generated frontend event creator
exist before anything emits or listens:

```go
func init() {
	application.RegisterEvent[IncidentCreated]("incident:created")
}
```

Emit it from Go once the service does the work:

```go
app.Event.Emit("incident:created", IncidentCreated{ID: id})
```

Consume it on the frontend with the generated typed creator (requires the `wails()` plugin from
step 2):

```ts
import { Events } from "@wailsio/runtime";
import { IncidentCreated } from "../bindings/github.com/wailsapp/wails/v3/internal/eventcreate";

Events.On(IncidentCreated, (e) => setIncidents((prev) => [...prev, e.data]));
```

Build with `-tags strictevents` during development to catch emits or listeners on names that
were never registered.

### 5. Add tray and notifications adapters

Each native feature needs a `//go:build !server` implementation and a `//go:build server` stub,
so `task build:server` still compiles. Create `app/tray.go`:

```go
//go:build !server

package app

func StartTray(app *application.App, win *application.WebviewWindow) *application.SystemTray {
	tray := app.SystemTray.New()
	tray.SetTemplateIcon(trayIconPNG)
	tray.AttachWindow(win)
	return tray
}
```

And its stub, `app/tray_server.go`:

```go
//go:build server

package app

func StartTray(app *application.App, win *application.WebviewWindow) *application.SystemTray {
	return nil
}
```

Repeat the same pair for notifications: `app/notifications.go` under `!server` registers the
`notifications` service and calls `RequestNotificationAuthorization`; `app/notifications_server.go`
under `server` is a no-op. Wire both from a single `startDesktop` function guarded the same way,
as shown in the `wails-v3` skill's project-layout reference.

### 6. Build

```bash
task build
```

This runs the per-OS build target (`darwin:build:native` on macOS) with `-tags production
-trimpath -ldflags="-w -s"` and produces `bin/<APP_NAME>`. Run this before reporting any change
as done; treat `task dev` as a convenience, not a correctness check.

Expected result: a binary in `bin/` with no build errors, and `task build:server` also succeeding
so the server-mode fallback still compiles.

### 7. Sign, notarize, and package

The scaffold ships no entitlements file, so generate one first:

```bash
wails3 setup entitlements
```

Record your signing identity and notarization profile once:

```bash
security find-identity -v -p codesigning
xcrun notarytool store-credentials "my-profile" \
  --apple-id you@example.com --team-id TEAMID --password app-specific-password
wails3 setup signing
```

Then sign, notarize, and package a DMG:

```bash
task darwin:sign:notarize
task darwin:package:dmg
```

`darwin:sign:notarize` packages the `.app`, signs it with the hardened runtime, submits it, and
staples the ticket. `darwin:package:dmg` builds the styled DMG from the signed `.app`.

## Verify it works

```bash
spctl --assess --verbose=2 "bin/<APP_NAME>.app"
```

Expected output:

```
bin/<APP_NAME>.app: accepted
source=Notarized Developer ID
```

Open the app from the DMG (not `bin/<APP_NAME>` directly) and confirm the tray icon appears and a
test notification is delivered.

Success! The app is scaffolded, has a working service with generated bindings, pushes a typed
event to the UI, and ships as a signed, notarized DMG.

## Troubleshooting

### Problem: `go mod tidy` or the build fails with a toolchain error
**Symptom**: `go: go.mod requires go >= 1.25.0 (running go 1.2x.x)`
**Cause**: The scaffold's `go.mod` declares `go 1.25.0`; an older toolchain cannot build it.
**Solution**: Install Go 1.25 or newer and confirm with `go version` before retrying.

### Problem: Notifications never appear
**Symptom**: `RequestNotificationAuthorization` returns `true` but no notification is shown, or
the permission prompt never appears at all.
**Cause**: On macOS, notifications authorize only from a packaged, signed bundle. `task dev`,
`task run`, and the bare `bin/<APP_NAME>` binary never show them.
**Solution**: Test notifications from a `task package` or `task darwin:sign:notarize` build, not
from the dev loop.

### Problem: An OAuth redirect never reaches the app
**Symptom**: The login flow hangs or errors inside the app window after redirecting.
**Cause**: The webview is served under a custom `wails` scheme, not `http(s)`, so an OAuth
provider's redirect URI cannot target it.
**Solution**: Open the OAuth URL in the system browser (`app.Browser.OpenURL`) and complete the
flow against a loopback HTTP listener in Go, then hand the result back to the app.

### Problem: The frontend calls a binding that no longer matches the Go signature
**Symptom**: A TypeScript type error or a runtime argument mismatch after changing a service
method, especially after a Wails tag bump.
**Cause**: Bindings are generated output; they go stale the moment a service signature or the
Wails module version changes and are not regenerated automatically.
**Solution**: Rerun `wails3 generate bindings -ts -i -clean -d frontend/bindings` (or `task
bindings` if you added the shortcut). `-clean` removes the old `frontend/bindings` tree first, so
the regenerated output is always complete; commit the result.

### Problem: `task common:update:build-assets` overwrote a manual `Info.plist` edit
**Symptom**: A hand-added `Info.plist` key (for example a custom `CFBundleURLTypes` entry) is
gone after a build.
**Cause**: `common:update:build-assets` regenerates `Info.plist`, the NSIS script, and the
`.desktop` file from `build/config.yml` every time it runs, discarding direct edits to the
generated files.
**Solution**: Make the change in `build/config.yml` (for example under `protocols:`) and let
`common:update:build-assets` regenerate the plist from it, instead of editing `Info.plist`
directly.

## Next steps

Now that the project builds and packages, you might want to:
- Add global shortcuts, dock badges, or additional windows — see the `native-features.md`
  reference in the `wails-v3` skill
- Set up the updater for future releases
- Package for Windows and Linux using the Taskfile targets described in the
  `build-and-sign.md` reference

## See also

- The `wails-v3` skill and its `references/project-layout.md`, `references/native-features.md`,
  `references/build-and-sign.md`, and `references/beta-gotchas.md`
- The `wails-go-developer` agent reference for architecture defaults and process
