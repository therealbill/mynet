# Beta Gotchas

**Verified against:** `v3.0.0-beta.25`, 2026-09-28. Each entry names the behavior, the consequence, and the workaround, with its source. Re-check every entry on a tag bump; delete the ones that stop being true.

## Version and drift

- **Go 1.25 or newer.** `go.mod` from the scaffold declares `go 1.25.0`. An older toolchain fails at `go mod tidy`. Source: scaffold.
- **Breaking changes ship inside the beta series** (macOS coordinate normalization, template removals, WebView2 loader changes across recent betas). Pin `go.mod` and the CLI to one tag; run `wails3 update cli` only as a deliberate bump. Source: owner's rebuild analysis R4; `release_notes.md` in the module.
- **The docs site blocks scripted reads.** `curl` and fetch tools get HTTP 403 from v3.wails.io pages; only `llms.txt` loads. Read `docs/mpress/content` from the wails repo at the tag. Source: checked 2026-09-28.
- **Docs and tag disagree on signing configuration.** Docs say set `SIGN_IDENTITY`, `KEYCHAIN_PROFILE`, `ENTITLEMENTS` in `build/darwin/Taskfile.yml`; the scaffold's `sign` tasks call `wails3 tool sign` and take identity from `wails3 setup signing` or `CLI_ARGS`. Follow the scaffold. Source: `build/darwin/Taskfile.yml` versus `guides/build/signing.md`.
- **Docs place `config.yml` at the root; the scaffold puts it at `build/config.yml`.** Source: `guides/dev/project-structure.md` versus scaffold.
- **`TitleBar` appears as both a constant and a struct in the docs.** Both are the same thing at this tag: `MacWindow.TitleBar` has type `MacTitleBar` (a struct), and `application.MacTitleBarHiddenInset` is a package-level `var` holding a predeclared `MacTitleBar` value, not a constant. Use the predeclared value as the scaffold does; hand-build a `MacTitleBar` literal only to vary its fields. Source: `pkg/application/webview_window_options.go:631` and `:844`, scaffold `main.go`, `features/windows/options.md`.
- **There is no `react-ts` template.** `react` is the TypeScript template; `react-js` is JavaScript. Source: `wails3 init -l`.

## Webview

- **Custom `wails` scheme.** WKWebView serves the app through `setURLSchemeHandler:forURLScheme:@"wails"`, so the page origin is not `http(s)`. OAuth redirect URIs must be `http(s)`, so run OAuth in the system browser (`app.Browser.OpenURL`) with a loopback listener in Go and never inside the webview. Source: `pkg/application/webview_window_darwin.go:181`.
- **Test `localStorage` and cookies per app.** Persistence under the custom scheme was not verified here; keep app state in Go and treat web storage as a cache. Source: owner's rebuild analysis R4 A9.
- **Safari feature set.** WKWebView, not Chromium; no Chromium-only APIs; test `contenteditable` early if you build an editor. Source: owner's analysis.
- **Translucency needs a private-API build.** `MacBackdropTranslucent`, `MacBackdropTransparent`, and `MacBackdropLiquidGlass` keep the webview opaque unless built with `-tags private_mac_apis`. CSS `backdrop-filter` over your own gradient needs no tag. Source: `features/windows/options.md`, `guides/build/macos.md`, constants at `pkg/application/webview_window_options.go:510-516`.

## Native features

- **Notifications need a packaged, signed bundle on macOS.** `task dev`, `task run`, and `bin/<APP_NAME>` never authorize; test from `task package` output signed with a real identity. Source: `features/notifications/overview.md` platform table.
- **URL schemes and single-instance URL relay only work from a packaged `.app`.** Handle both `ApplicationLaunchedWithUrl` and `OnSecondInstanceLaunch`. Source: `guides/distribution/custom-protocols.md`, `examples/single-instance-url-scheme`.
- **No window-state persistence.** Save and restore `Bounds()` yourself. Source: owner's analysis R4; no such option in `WebviewWindowOptions`.
- **Hide shortcut plus terminate-after-last-window quits the app.** Leave `ApplicationShouldTerminateAfterLastWindowClosed` false for tray-resident or hotkey-driven apps. Source: `features/keyboard/global-shortcuts.md`.
- **Global shortcut registered twice errors and keeps the first callback.** Unregister before re-registering. Source: same page.
- **Not available at all:** Touch Bar, macOS Services menu, CoreSpotlight indexing, dock bounce, custom dock menu, `.pkg` output. Source: module `pkg/services` and `pkg/application` surface; owner's analysis.
- **Server mode drops tray, dialogs, windows, shortcuts, notifications.** Every native adapter needs a `server` build-tag stub. Source: `guides/server-build.md`.

## Dev loop and build

- **`task dev` reliability is an open GA blocker.** Treat `task build` as the acceptance step and `task dev` as convenience. Source: owner's analysis R4 A10.
- **`dev_mode.ignore.file` excludes `*_test.go` and the watcher does not rebuild on JSON changes** unless `"*.json"` is added to `watched_extension`. Source: `build/config.yml`.
- **`common:update:build-assets` overwrites hand edits** to `Info.plist`, NSIS, and `.desktop` files. Edit `build/config.yml` instead. Source: `build/config.yml` header comment.
- **`generate:bindings` runs `-clean=true`**, deleting `frontend/bindings` first; parallel universal builds are serialized with `run: once` for that reason. Source: `build/Taskfile.yml`.
- **`mac.LoadResource` and `mac.ResourceFS` return `ErrNotInAppBundle`** outside a bundle, so bundle resources are invisible under `task dev`. Source: `guides/build/macos.md`.
- **Cross-compiled macOS binaries are unsigned and ad-hoc signing is applied by `package`.** Distribution still needs Developer ID signing on a Mac. Source: `build/darwin/Taskfile.yml`.
