# Native Features

**Verified against:** `v3.0.0-beta.25`, 2026-09-28. Sources: module source under `pkg/application` and `pkg/services`, `examples/{systray-menu,notifications,global-shortcuts,single-instance-url-scheme,events,updater}`, docs pages `features/windows/options.md`, `features/menus/systray.md`, `features/notifications/overview.md`, `features/keyboard/global-shortcuts.md`, `features/platform/dock.md`, `guides/events-reference.md`, `guides/single-instance.md`, `guides/distribution/custom-protocols.md`, `guides/server-build.md`, and the owner's mock-up `app/` package. Each section ends with the constraint that is not obvious from the API.

Adapter rule: one file per feature under `app/`, each with a `//go:build !server` twin that is a no-op under `//go:build server`. Keep the Wails types at the edge; hand domain code plain functions and channels.

## Windows

```go
win := app.Window.NewWithOptions(application.WebviewWindowOptions{
	Name: "main", Title: "App", Width: 1280, Height: 820, MinWidth: 960, MinHeight: 640,
	BackgroundColour: application.NewRGB(10, 10, 10),
	URL:              "/",                      // hash route per window: "/#/present"
	Mac: application.MacWindow{
		TitleBar:                application.MacTitleBarHiddenInset, // traffic lights inset, no title text
		InvisibleTitleBarHeight: 52,                                  // drag region; only with hidden or transparent title bar
	},
})
win.Show(); win.Focus(); win.Hide(); win.Restore()
win.EmitEvent("nav", "/incidents/new")           // window-scoped event
win.RegisterHook(events.Common.WindowClosing, func(e *application.WindowEvent) { win.Hide(); e.Cancel() })
```

Other options: `Frameless`, `AlwaysOnTop`, `Hidden`, `HideOnFocusLost`, `HideOnEscape` (tray popups), `StartState`, `KeyBindings map[string]func(application.Window)`, `DevToolsEnabled`. Mac: `Backdrop` (`MacBackdropNormal`, `MacBackdropTranslucent`, `MacBackdropTransparent`, `MacBackdropLiquidGlass`; the last three need `private_mac_apis` for webview transparency), `WindowLevel`, `CollectionBehavior` (bitmask; `MacWindowCollectionBehaviorCanJoinAllSpaces | MacWindowCollectionBehaviorFullScreenAuxiliary` for Spotlight-style panels), `WindowClass: MacWindowClassPanel` with `PanelPreferences: MacPanelPreferences{NonActivating: true}` for a panel that does not activate the app.

Constraint: there is no window-state persistence. Save `win.Bounds()` yourself (kvstore service or your store) and restore before `Show`.

## Menus

```go
m := application.NewMenu()
appMenu := m.AddSubmenu(appName)                    // hand-built so it can carry Preferences…
appMenu.AddRole(application.About)
appMenu.AddSeparator()
appMenu.Add("Preferences…").SetAccelerator("CmdOrCtrl+,").OnClick(func(*application.Context) { onPrefs() })
appMenu.AddSeparator()
appMenu.AddRole(application.ServicesMenu)
appMenu.AddSeparator()
appMenu.AddRole(application.Hide); appMenu.AddRole(application.HideOthers); appMenu.AddRole(application.UnHide)
appMenu.AddSeparator()
appMenu.AddRole(application.Quit)
inc := m.AddSubmenu("Incident")
inc.Add("New Incident").SetAccelerator("CmdOrCtrl+I").OnClick(func(*application.Context) { onNew() })
m.AddRole(application.EditMenu); m.AddRole(application.WindowMenu); m.AddRole(application.HelpMenu)
app.Menu.Set(m)
```

Item types: `Add`, `AddSeparator`, `AddCheckbox(label, checked)`, `AddRadio(label, checked)`, `AddSubmenu`, `AddRole`. `ctx.ClickedMenuItem()` gives the item for `SetLabel`, `SetEnabled`, `Checked`. Constraint: `application.AppMenu` role cannot be extended; build the app menu by hand from the roles above when it needs a Preferences item. Role menus need a running app, so keep menu construction out of unit tests.

## System tray

```go
tray := app.SystemTray.New()
tray.SetTemplateIcon(iconTemplatePNG)              // macOS: black and clear only, 18 to 22 px, file named *Template.png
tray.SetTooltip("App"); tray.SetLabel("3")          // label shows text beside the icon on macOS
tray.SetMenu(menu)                                  // rebuild and SetMenu again to update; debounce bursts
tray.AttachWindow(win); tray.WindowOffset(8); tray.WindowDebounce(200 * time.Millisecond)
tray.OnClick(func() { … }); tray.OnRightClick(func() { … })
```

Set `Mac.ActivationPolicy: application.ActivationPolicyAccessory` for a menu-bar-only app with no Dock icon. Constraint: not available in server mode. Rebuilding the menu on every event thrashes; coalesce with a 200 ms debouncer as the mock-up does.

## Notifications

```go
ns := notifications.New()                          // github.com/wailsapp/wails/v3/pkg/services/notifications
app := application.New(application.Options{Services: []application.Service{application.NewService(ns)}})
ok, err := ns.CheckNotificationAuthorization()
if !ok { ok, err = ns.RequestNotificationAuthorization() }
ns.RegisterNotificationCategory(notifications.NotificationCategory{
	ID: "incident", Actions: []notifications.NotificationAction{{ID: "VIEW", Title: "View"}, {ID: "ACK", Title: "Acknowledge"}},
})
ns.SendNotificationWithActions(notifications.NotificationOptions{
	ID: "inc-42", Title: "P1 opened", Subtitle: "checkout", Body: "…", CategoryID: "incident",
	Data: map[string]any{"route": "/incidents/42"},
})
ns.OnNotificationResponse(func(r notifications.NotificationResult) {
	if r.Error != nil { return }
	switch r.Response.ActionIdentifier { case "VIEW": open(r.Response.UserInfo["route"].(string)) }
})
```

Also `SendNotification` (no actions), `UpdateNotification` by ID, `RemoveNotification`, `Sound`, `Attachments` (absolute paths; write embedded assets to disk first), `ThreadID`, `InterruptionLevel`, `Schedule`. Constraint: on macOS notifications authorize only from a **packaged and signed** bundle; `task dev` and the bare binary never show them, and `CheckNotificationAuthorization` reports `true` unconditionally on Windows and Linux. Deliver off the event goroutine and dedupe repeats.

## Shortcuts

```go
// in-app: fires only when a window of this app has focus
app.KeyBinding.Add("CmdOrCtrl+K", func(w application.Window) { … })
// global: fires system-wide, Carbon hot keys on macOS, no Accessibility permission
if err := app.GlobalShortcut.Register("CmdOrCtrl+Shift+I", func() { win.Show(); win.Focus() }); err != nil { … }
app.GlobalShortcut.Unregister("CmdOrCtrl+Shift+I"); app.GlobalShortcut.IsRegistered(...); app.GlobalShortcut.GetAll()
```

Accelerator format: `CmdOrCtrl+Shift+G`, `Ctrl+Alt+K`, `Cmd+Option+Space`, `Super+D`, function keys. Constraints: registering the same chord twice in one app returns an error and keeps the first callback (unregister first to change it); on macOS a chord owned by another app still registers, on Windows and X11 it fails; hot keys bind to physical key positions; with `ApplicationShouldTerminateAfterLastWindowClosed: true` a hide shortcut on the only window quits the app.

## Dock

```go
d := dock.New()                                     // github.com/wailsapp/wails/v3/pkg/services/dock
application.NewService(d)
d.SetBadge("3"); d.SetBadge(""); d.RemoveBadge(); d.HideAppIcon(); d.ShowAppIcon()
```

Constraint: badge only; there is no dock bounce or custom dock menu API at this tag.

## Events

```go
func init() { application.RegisterEvent[IncidentCreated]("incident:created") }   // typed; panics on a type conflict
app.Event.Emit("incident:created", ev)              // to every window
win.EmitEvent("nav", "/x")                          // to one window
app.Event.On("ui:action", func(e *application.CustomEvent) { … e.Data … e.Sender })
app.Event.OnApplicationEvent(events.Common.ApplicationStarted, func(*application.ApplicationEvent) { … })
app.Event.OnApplicationEvent(events.Common.ThemeChanged, func(e *application.ApplicationEvent) { e.Context().IsDarkMode() })
```

Frontend: `Events.On(IncidentCreated, e => e.data)` with the typed creator from the generated bindings once `wails()` from `@wailsio/runtime/plugins/vite` is in `vite.config.ts`; untyped `Events.On("name", …)` also works. Bridge pattern: subscribe once to your domain event bus and re-emit every event under its kind name. `-tags strictevents` warns on unregistered names in development.

## Streams and dialogs

Streams (`app.HandleStream(name string, handler application.StreamHandler)`, `guides/streams.md`) carry byte streams with backpressure over the asset server without a TCP port; use them for token-by-token model output.

Message dialogs are `app.Dialog.Info()`, `.Question()`, `.Warning()`, `.Error()`, each returning `*application.MessageDialog`. `SetTitle` and `SetMessage` chain, but `AddButton` returns `*application.Button` (for `OnClick`, `SetAsDefault`, `SetAsCancel`), so `Show()` cannot chain off it:

```go
d := app.Dialog.Info().SetTitle("Saved").SetMessage("Incident 42 written.")
d.AddButton("OK").SetAsDefault()
d.Show()
```

File dialogs are a separate builder: `app.Dialog.OpenFile()` and `.SaveFile()` have no `AddButton` and no `Show`; they terminate with `PromptForSingleSelection() (string, error)` or, for open, `PromptForMultipleSelection() ([]string, error)`. `SaveFile()` has `SetMessage` but no `SetTitle`. Constraint: dialogs are not available in server mode.

## Single instance and URL schemes

```go
SingleInstance: &application.SingleInstanceOptions{
	UniqueID: "com.example.app",
	OnSecondInstanceLaunch: func(d application.SecondInstanceData) { /* d.Args, d.WorkingDir; focus the window */ },
	EncryptionKey: key,                             // [32]byte, not a slice; without it treat d.Args as untrusted
},
app.Event.OnApplicationEvent(events.Common.ApplicationLaunchedWithUrl, func(e *application.ApplicationEvent) { url := e.Context().URL() })
```

Declare schemes under `protocols:` in `build/config.yml` and run `task common:update:build-assets`; `application.Options` has no `Protocols` field. Constraint: scheme launches only work from a packaged `.app` whose `Info.plist` carries `CFBundleURLTypes`; handle both the cold path (`ApplicationLaunchedWithUrl`) and the warm path (`OnSecondInstanceLaunch` with the URL in `Args`).

## Server mode

`-tags server` keeps services, bindings, events, and streams; drops windows, tray, dialogs, shortcuts, notifications. `Server: application.ServerOptions{Host, Port, TLS, WebSocketOriginPatterns}`; `/health` returns `{"status":"ok"}`. Use it for Playwright tests and demos in a browser.
