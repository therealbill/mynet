# desktop-development Wails v3 Conversion Implementation Plan

> **For agentic workers:** Execute this plan with an **Agent Team** of up to four concurrent agents as laid out in the Execution Model section. Each task names its owning agent and model. Do not use git worktrees or ad-hoc parallel subagents as the execution method. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the Electron+Go agent in `desktop-development` with a Wails v3 agent (`wails-go-developer`), add a `wails-v3` reference skill whose every claim is verified against the `v3.0.0-beta.25` scaffold and module source, add structural and Jev routing tests, rewrite the five Diátaxis pages, and ship it as 2.0.0.

**Architecture:** Agents are Markdown files with YAML frontmatter under `desktop-development/agents/`; example blocks live in the body (repo convention since commit e6fc438). The skill is `SKILL.md` plus four reference files under `references/`. Verification is a structural lint, a Jev routing check over sibling descriptions, the ai-development `audit-agents.py` audit, a Jev quality loop over the skill files, the Diátaxis cross-link validator, and the plugin validator. Docs under `desktop-development/docs/` are Diátaxis pages mounted into the marketplace Hugo site.

**Tech Stack:** Markdown, YAML frontmatter, bash, Python 3 standard library (`urllib`), TypeSafe HTTP API (`TYPESAFE_API_KEY` in the environment, model `jev-latest`), Wails `v3.0.0-beta.25` CLI and module, Task, Go 1.25+, git.

**Spec:** `docs/specs/2026-09-28-desktop-development-wails-v3-design.md`.

---

## Decisions

Settled with the plugin owner. Do not reopen them.

| Topic | Decision |
| --- | --- |
| Plugin shape | One agent (`wails-go-developer`) plus one reference skill (`wails-v3`) |
| Electron | Dropped entirely; `electron-go-pro.md` is deleted; no migration trigger |
| Platforms | macOS first with full defaults; Windows and Linux at Taskfile-target level only; no mobile |
| Frontend default | React + Vite + TypeScript + Tailwind via the built-in `react` template (there is no `react-ts` template) |
| Beta drift | Pin `v3.0.0-beta.25`; precedence for facts is module source and scaffold at the tag, then the docs at the tag, then the skill |
| Docs source | The live site returns 403 to scripted fetches. Read the docs at the tag from the local sparse clone (see Verified Facts) and the module source in the Go module cache |
| Verification bar | agent-modernizer audit with no Must-fix and no Should-fix, lint and routing checks green, Jev quality loop on the skill with findings confirmed by reading, Diátaxis validator, plugin validator |
| Branch | `wails-v3-conversion` (already created off master); one commit per task; explicit-path staging |

## Verified Facts

Checked on 2026-09-28 on this machine (`wails3 version` = `v3.0.0-beta.25`, Go 1.27.1, Task installed) by running `wails3 init -n wailsprobe -t react -d wailsprobe -mod example.com/wailsprobe` and `task build`. Every reference file cites these; do not re-derive them from memory.

**Local sources (absolute paths):**

- Docs at the tag: `/private/tmp/claude-501/-Users-bill-Projects-mynet/8934198f-2ad7-4e6f-823c-41527251b7fd/scratchpad/wails-docs/docs/mpress/content/` (English pages are the top-level directories `features/`, `guides/`, `concepts/`, `reference/`, `tutorials/`; language subdirectories such as `de/`, `fr/` are translations, ignore them). If the directory is missing, recreate it: `git clone --filter=blob:none --sparse --depth 1 --branch v3.0.0-beta.25 https://github.com/wailsapp/wails wails-docs && cd wails-docs && git sparse-checkout set docs/mpress/content`.
- Module source: `~/go/pkg/mod/github.com/wailsapp/wails/v3@v3.0.0-beta.25/` (`pkg/application`, `pkg/services/{notifications,dock,sqlite,kvstore,log,fileserver}`, `pkg/updater`, `examples/`).
- Scaffold: `/private/tmp/claude-501/-Users-bill-Projects-mynet/8934198f-2ad7-4e6f-823c-41527251b7fd/scratchpad/wailsprobe/wailsprobe/` (built; `bin/wailsprobe` is 9.3 MB). Recreate with the `wails3 init` command above if missing; note `-d wailsprobe -n wailsprobe` nests the project one level deeper than `-d`.
- Owner's mock-up: `~/Projects/Praetor/tactical-bridge-wails` (read-only; `app/` adapters, `main_desktop.go` and `main_server.go` build-tag split, `Taskfile.yml` with `bindings` and `web` targets).

**Scaffold layout (`react` template):**

```
.gitignore  go.mod  go.sum  main.go  greetservice.go  README.md  Taskfile.yml
build/
  config.yml  Taskfile.yml  appicon.png  appicon.icon/
  darwin/   Assets.car  Info.plist  Info.dev.plist  Taskfile.yml  icons.icns  dmg-background.png  dmg-file-icon.{icns,png}
  windows/  Taskfile.yml  icon.ico  info.json  wails.exe.manifest  nsis/  msix/
  linux/    Taskfile.yml  desktop  appimage/  nfpm/
  ios/  android/  docker/
frontend/
  index.html  package.json  vite.config.ts  tsconfig.json  .npmrc  public/  src/{App.tsx,main.tsx,vite-env.d.ts}
  bindings/   (generated; gitignored? no: only dist and node_modules are ignored)
  dist/       (generated, gitignored)
bin/          (generated, gitignored)
```

- `go.mod`: `go 1.25.0`, `require github.com/wailsapp/wails/v3 v3.0.0-beta.25`.
- `main.go`: `//go:embed all:frontend/dist`; `application.RegisterEvent[string]("time")` in `init()`; `application.New(application.Options{Name, Description, Services: []application.Service{application.NewService(&GreetService{})}, Assets: application.AssetOptions{Handler: application.AssetFileServerFS(assets)}, Mac: application.MacOptions{ApplicationShouldTerminateAfterLastWindowClosed: true}})`; `app.Window.NewWithOptions(application.WebviewWindowOptions{Title, Width: 1000, Height: 618, Mac: application.MacWindow{InvisibleTitleBarHeight: 50, Backdrop: application.MacBackdropTranslucent, TitleBar: application.MacTitleBarHiddenInset}, BackgroundColour: application.NewRGB(6, 7, 15), URL: "/"})`; `app.Event.Emit("time", now)` from a goroutine; `app.Run()`.
- `frontend/package.json`: `@wailsio/runtime: latest`, React 18, Vite 8, `@vitejs/plugin-react` 6, TypeScript 5. Scripts `dev`, `build:dev` (`tsc && vite build --minify false --mode development`), `build` (`tsc && vite build --mode production`). No Tailwind in the template; it is added per project.
- Generated bindings land in `frontend/bindings/<module path>/<service>.ts` plus `index.ts`, and typed events in `frontend/bindings/github.com/wailsapp/wails/v3/internal/eventcreate.ts` and `eventdata.d.ts`.
- `.gitignore`: `.task`, `bin`, `frontend/dist`, `frontend/node_modules`, Linux appimage build dir, WebView2 bootstrapper, iOS and Android overlays.

**Taskfile targets (root `Taskfile.yml`):** `build`, `package`, `run`, `dev` (`wails3 dev -config ./build/config.yml -port {{.VITE_PORT}}`, Vite port default 9245), `build:server`, `run:server`, `build:docker`, `run:docker`, `setup:docker`, `install:msix:tools`. Root vars: `APP_NAME`, `BIN_DIR: bin`, `PACKAGE_MANAGER` (npm default), `VITE_PORT`, `GOOS` (dispatches to `{{.GOOS}}:build`).

**Common tasks (`build/Taskfile.yml`, namespace `common:`):** `install:frontend:deps`, `build:frontend` (runs `generate:bindings` first), `generate:bindings` (`wails3 generate bindings -f '{{.BUILD_FLAGS}}' -clean=true -ts -i`), `generate:icons`, `dev:frontend`, `update:build-assets` (`wails3 update build-assets -name ... -config config.yml -dir .`), `build:server` (`-tags server,production -trimpath -ldflags="-w -s"`, output `bin/<APP_NAME>-server`), `run:server`, `build:docker`, `run:docker`, `setup:docker`. There is no `bindings` shortcut at the root; the owner's mock-up adds one (`wails3 generate bindings -ts -i -clean -d frontend/bindings`).

**Darwin tasks (`build/darwin/Taskfile.yml`, namespace `darwin:`):** `build` (native or Docker), `build:native` (internal; `-tags production -trimpath -buildvcs=false -ldflags="-w -s"`, `CGO_ENABLED=1`, `MACOSX_DEPLOYMENT_TARGET=12.0`), `build:universal` (amd64 + arm64 then `lipo`), `package` (`create:app:bundle`, ad-hoc `codesign --force --deep --sign -`), `package:universal`, `package:dmg`, `create:dmg` (`wails3 tool package --format dmg ...` with vars `DMG_BACKGROUND`, `DMG_VOLUME_ICON`, `DMG_FILE_ICON`, `DMG_WINDOW_WIDTH`, `DMG_WINDOW_HEIGHT`, `DMG_FILES`), `run` (builds `bin/<APP_NAME>.dev.app` with `Info.dev.plist`, ad-hoc signs, launches), `sign` (`wails3 tool sign --input bin/<APP_NAME>.app {{.CLI_ARGS}}`), `sign:notarize` (`... --notarize`). The scaffold ships **no** `entitlements.plist`; `wails3 setup entitlements` generates `build/darwin/entitlements.plist` and `entitlements.dev.plist`.

**Docs-versus-tag discrepancies found (record in `beta-gotchas.md`):**

1. `guides/build/signing.md` and `guides/build/macos.md` say to set `SIGN_IDENTITY`, `KEYCHAIN_PROFILE`, `ENTITLEMENTS` vars in `build/darwin/Taskfile.yml`. The beta.25 scaffold's `sign` tasks instead call `wails3 tool sign`, which reads identity and keychain profile from `wails3 setup signing` (stored in `~/.config/wails/defaults.yaml`); overrides go through `CLI_ARGS`: `task darwin:sign -- --identity "Developer ID Application: ..." --entitlements build/darwin/entitlements.plist`.
2. `guides/dev/project-structure.md` shows `config.yml` at the project root; the scaffold puts it at `build/config.yml`.
3. `features/windows/options.md` shows `TitleBar` as a `MacTitleBar` struct; the scaffold uses the constant `application.MacTitleBarHiddenInset` and builds. Both spellings appear in the source at this tag; prefer what the scaffold uses.
4. The docs site's `.md` and page URLs return HTTP 403 to `curl` and WebFetch; the `llms.txt` index loads. Use the local clone at the tag.

**API facts from the module source at the tag:**

- `application.Options` fields include `Name`, `Description`, `Icon`, `Mac`, `Windows`, `Linux`, `Services`, `Assets`, `Logger`, `LogLevel`, `KeyBindings`, `OnShutdown`, `ShouldQuit`, `FileAssociations`, `SingleInstance *SingleInstanceOptions`, `Server ServerOptions`. There is no `Protocols` field; URL schemes are declared under `protocols:` in `build/config.yml`.
- `application.MacOptions`: `ActivationPolicy`, `ApplicationShouldTerminateAfterLastWindowClosed`.
- `application.WebviewWindowOptions` fields include `Name`, `Title`, `Width`, `Height`, `MinWidth`, `MinHeight`, `URL`, `Frameless`, `Hidden`, `AlwaysOnTop`, `HideOnFocusLost`, `HideOnEscape`, `BackgroundColour`, `StartState`, `KeyBindings`, `DevToolsEnabled`, `Mac MacWindow`.
- `application.MacWindow`: `Backdrop`, `TitleBar`, `InvisibleTitleBarHeight`, `Appearance`, `WindowLevel`, `CollectionBehavior`, `WindowClass`, `PanelPreferences`, `LiquidGlass`, `TabbingMode`.
- Managers on `*application.App`: `Window`, `Event`, `Menu`, `SystemTray`, `Dialog`, `KeyBinding`, `GlobalShortcut`, `Clipboard`, `Browser`, `Screen`, `Env`, `ContextMenu`, plus `Updater` and `Logger`.
- Services: `notifications.New()` with `RequestNotificationAuthorization`, `CheckNotificationAuthorization`, `SendNotification`, `SendNotificationWithActions`, `RegisterNotificationCategory`, `OnNotificationResponse`; `dock.New()` with `SetBadge`, `RemoveBadge`, `HideAppIcon`, `ShowAppIcon`; `sqlite.New()`/`NewWithConfig` with `Execute`, `Query`, `Prepare`; `kvstore`; `log`; `fileserver`.
- Updater: `app.Updater.Init(updater.Config{CurrentVersion, Providers, PublicKey, CheckInterval})`, providers `github.New(github.Config{Repository, ChecksumAsset})`, `keygen`, `appcast`, `endpoint`; `CheckAndInstall(ctx)`, `Check`, `DownloadAndInstall`, `Restart`; events `updater.EventUpdateAvailable`, `EventDownloadProgress`, `EventUpdateReady`, `EventError`.
- WKWebView serves the app under the custom `wails` URL scheme (`[config setURLSchemeHandler:delegate forURLScheme:@"wails"]` in `pkg/application/webview_window_darwin.go:181`). OAuth providers require an `http(s)` redirect URI, so the redirect cannot land in the webview; run OAuth in the system browser with a loopback listener.
- Global shortcuts use Carbon hot keys on macOS and need no Accessibility permission; `Register` returns an error on a duplicate within the app; with `ApplicationShouldTerminateAfterLastWindowClosed: true` a hide-shortcut on the only window quits the app.
- Notifications on macOS require user authorization and a packaged, signed bundle; `CheckNotificationAuthorization` returns `true` unconditionally on Windows and Linux.
- Server mode: `-tags server`; `ServerOptions{Host, Port, TLS, ...}`; `/health` endpoint; no system tray or native dialogs; each browser tab is a window named `browser-N`.
- Single instance: `SingleInstanceOptions{UniqueID, OnSecondInstanceLaunch func(SecondInstanceData), AdditionalData, EncryptionKey}`; `events.Common.ApplicationLaunchedWithUrl` with `e.Context().URL()` for URL-scheme launches.
- Typed events: `application.RegisterEvent[T]("name")` at init; Vite plugin `import wails from '@wailsio/runtime/plugins/vite'`; frontend imports typed creators from the generated bindings; `-tags strictevents` warns on unregistered events.
- Private macOS APIs (webview transparency for translucent backdrops, Liquid Glass grouping, programmatic inspector) need `-tags private_mac_apis`; without it the options are accepted and the private parts are no-ops.

## Execution Model

| Agent | Model | Tasks | Starts when |
| --- | --- | --- | --- |
| `harness-and-close-out` | sonnet | 1, then 5, 6, 7, 8 | Task 1 immediately; Task 5 after Tasks 3 and 4 are committed; Tasks 6 through 8 after Tasks 2 through 5 are committed |
| `agent-rewrite` | opus | 2 | Task 1 committed |
| `skill-layout-build` | sonnet | 3 | Task 1 committed |
| `skill-native-gotchas` | opus | 4 | Task 1 committed |

Rules for every agent:

- Work in `~/Projects/mynet` on branch `wails-v3-conversion`. Stage by explicit path only; never `git add -A` or `git add .`.
- Edit only the files listed under your tasks. If a fix belongs in another agent's file, leave a note in your final report instead of editing it.
- `agent-rewrite` loads the `plugin-dev:agent-development` skill before editing the agent file. `skill-layout-build` and `skill-native-gotchas` load `plugin-dev:skill-development`. `harness-and-close-out` loads `typesafe:typesafe-ai` before Task 5.
- `skill-native-gotchas` dispatches the `go-architect` subagent to review every Go snippet in `native-features.md` against the module source before committing (CLAUDE.md: Go code is reviewed by go-architect).
- Every fact in a reference file comes from the Verified Facts section, the local docs clone, the module source, or the scaffold. Cite the source page or file in the `Verified against` block at the top of each reference. A claim you cannot trace is removed, not softened.
- Run `desktop-development/tests/lint-agents.sh` before any commit that touches the agent file. It must print no `FAIL` line.
- No file gets an "Authored by" line. No timeline or effort estimates anywhere.
- Commit messages are one line, imperative, no trailer.

## File Structure

| Path | Responsibility | Owner |
| --- | --- | --- |
| `desktop-development/tests/lint-agents.sh` | Structural checks on the agent file | harness-and-close-out |
| `desktop-development/tests/route-check.py` | Jev routing simulation over sibling descriptions with expected winners | harness-and-close-out |
| `desktop-development/tests/README.md` | How to run the checks | harness-and-close-out |
| `desktop-development/agents/wails-go-developer.md` | The agent (new); `agents/electron-go-pro.md` deleted | agent-rewrite |
| `desktop-development/skills/wails-v3/SKILL.md` | Pin policy, doc-reading rule, reference table | skill-layout-build |
| `desktop-development/skills/wails-v3/references/project-layout.md` | Scaffold tree, Taskfile targets, config, build-tag stub pattern | skill-layout-build |
| `desktop-development/skills/wails-v3/references/build-and-sign.md` | Packaging, signing, notarization, DMG, updater, other platforms | skill-layout-build |
| `desktop-development/skills/wails-v3/references/native-features.md` | Windows, menus, tray, notifications, shortcuts, dock, events, streams, dialogs, single instance, URL schemes | skill-native-gotchas |
| `desktop-development/skills/wails-v3/references/beta-gotchas.md` | Verified-against list of beta edges and docs-versus-tag discrepancies | skill-native-gotchas |
| `desktop-development/docs/**` | Five Diátaxis pages plus `reference/skills.md` (new) and `howto/set-up-wails-go-project.md` (replaces the Electron how-to) | harness-and-close-out via Diátaxis agents |
| `desktop-development/CHANGELOG.md`, `desktop-development/.claude-plugin/plugin.json`, `.claude-plugin/marketplace.json` | Versioning | harness-and-close-out |

---

### Task 1: Tests, RED baseline, and Electron removal probes

**Agent:** harness-and-close-out (sonnet)

**Files:**
- Create: `desktop-development/tests/lint-agents.sh`
- Create: `desktop-development/tests/route-check.py`
- Create: `desktop-development/tests/README.md`

- [ ] **Step 1: Write the lint script**

```bash
#!/usr/bin/env bash
# Structural checks for desktop-development agent files.
# Usage: desktop-development/tests/lint-agents.sh [agent.md ...]
# Exits 1 if any check fails.
set -u
here="$(cd "$(dirname "$0")" && pwd)"
files=("$@")
[ ${#files[@]} -eq 0 ] && files=("$here"/../agents/*.md)
fail=0

ok()   { printf 'ok   %s\n' "$1"; }
bad()  { printf 'FAIL %s\n' "$1"; fail=1; }

for f in "${files[@]}"; do
  name="$(basename "$f" .md)"
  desc="$(awk '/^description: >/{flag=1; next} /^[a-z]+:/{flag=0} flag' "$f" | tr '\n' ' ')"
  body="$(awk 'BEGIN{fm=0} /^---$/{fm++; next} fm>=2' "$f")"
  body_no_examples="$(printf '%s\n' "$body" | awk '/<example>/{s=1} !s{print} /<\/example>/{s=0}')"
  examples="$(grep -c '^<example>' "$f")"
  words="$(printf '%s' "$desc" | wc -w | tr -d ' ')"

  grep -q 'Use when' <<<"$desc" \
    && ok "$name: description has 'Use when'" \
    || bad "$name: description lacks 'Use when'"

  grep -q 'Do not use' <<<"$desc" \
    && ok "$name: description has a 'Do not use' handoff" \
    || bad "$name: description lacks a 'Do not use' handoff"

  grep -qE '\(use [a-z-]+\)' <<<"$desc" \
    && ok "$name: handoff names a sibling agent" \
    || bad "$name: handoff does not name a sibling agent"

  [ "$words" -le 120 ] \
    && ok "$name: description is $words words" \
    || bad "$name: description is $words words (limit 120)"

  grep -q '<example>' <<<"$desc" \
    && bad "$name: example blocks inside description" \
    || ok "$name: no example blocks inside description"

  grep -q '^\*\*Output:\*\*' <<<"$body" \
    && ok "$name: body has an Output section" \
    || bad "$name: body lacks an Output section"

  grep -q 'Proactive trigger:' "$f" \
    && ok "$name: has a proactive example" \
    || bad "$name: has no proactive example"

  [ "$examples" -ge 2 ] && [ "$examples" -le 5 ] \
    && ok "$name: has $examples examples" \
    || bad "$name: has $examples examples (want 2 to 5)"

  chars="$(printf '%s' "$body_no_examples" | wc -c | tr -d ' ')"
  [ "$chars" -lt 3000 ] \
    && ok "$name: body is $chars chars excluding examples" \
    || bad "$name: body is $chars chars excluding examples (limit 3000)"

  grep -qE '\b(I am|I will|I would|I have)\b' <<<"$body" \
    && bad "$name: body uses first person" \
    || ok "$name: body is second person"

  grep -qiE 'electron' "$f" \
    && bad "$name: mentions Electron" \
    || ok "$name: no Electron mention"
done
exit $fail
```

- [ ] **Step 2: Write the routing check**

```python
#!/usr/bin/env python3
"""Routing simulation for the desktop-development agent.

Sends the descriptions of wails-go-developer and its marketplace siblings plus
sample requests to the TypeSafe Jev API as one Choice question per request and
checks the winner and confidence.

Usage: TYPESAFE_API_KEY=... desktop-development/tests/route-check.py
Exit 0 when every must-pass probe routes as expected with confidence >= THRESHOLD,
1 otherwise, 2 when the key is missing.
"""
import json, os, re, sys, time, urllib.error, urllib.request

API = "https://api.typesafe.ai/v1/systemone"
HERE = os.path.dirname(os.path.abspath(__file__))
ROOT = os.path.abspath(os.path.join(HERE, "..", ".."))
THRESHOLD = 0.70

# name -> path relative to the marketplace root
AGENTS = {
    "wails-go-developer": "desktop-development/agents/wails-go-developer.md",
    "cli-developer": "cli-development/agents/cli-developer.md",
    "go-tui-developer": "cli-development/agents/go-tui-developer.md",
    "frontend-developer": "web-development/agents/frontend-developer.md",
    "react-specialist": "web-development/agents/react-specialist.md",
    "go-architect": "backend-development/agents/go-architect.md",
}

# (request, expected winner; "none" means no listed agent; None is informational)
PROBES = [
    ("I want a macOS desktop app in Go for tracking incidents, with a real window and menus", "wails-go-developer"),
    ("Add a menu bar icon and native notifications to our Wails app", "wails-go-developer"),
    ("Package this Wails app for distribution with signing and a DMG", "wails-go-developer"),
    ("Our Go tool has a CLI; can you give it a small window that shows status?", "wails-go-developer"),
    ("Set up wails3 dev and generate the TypeScript bindings for our services", "wails-go-developer"),
    ("Register a global shortcut so the app window pops up from anywhere", "wails-go-developer"),
    ("Build a TUI for browsing API responses with Bubble Tea", "go-tui-developer"),
    ("Build a CLI for managing database migrations in Go", "cli-developer"),
    ("Add shell completions for zsh and fish to our Go CLI", "cli-developer"),
    ("Create a Next.js marketing site with SSG", "frontend-developer"),
    ("Our React table re-renders on every keystroke", "react-specialist"),
    ("Design the service boundaries and API for our Go backend", "go-architect"),
    ("Build a macOS app with Electron and a Go backend", "none"),
    ("How should the Electron renderer talk to our Go process over IPC?", "none"),
    ("Package our Electron app with electron-builder and notarize it", "none"),
    ("Build an iOS app with Swift", "none"),
    ("Add dark mode to our Wails app's React frontend", None),
]


def parse_description(path):
    """Return the description value, handling scalar, folded `>` and `|` forms."""
    lines = open(path, encoding="utf-8").read().split("\n")
    if lines[0].strip() != "---":
        sys.exit(f"no frontmatter in {path}")
    i = 1
    while i < len(lines) and lines[i].strip() != "---":
        m = re.match(r"^description:\s*(.*)$", lines[i])
        if m:
            val = m.group(1).strip()
            if val in (">", "|", ">-", "|-"):
                block = []
                i += 1
                while i < len(lines) and lines[i].strip() != "---" and not re.match(r"^[A-Za-z_][\w-]*:", lines[i]):
                    block.append(lines[i].strip())
                    i += 1
                return " ".join(b for b in block if b)
            return val.strip('"\'')
        i += 1
    sys.exit(f"no description in {path}")


def call(state, questions):
    key = os.environ.get("TYPESAFE_API_KEY")
    if not key:
        print("SKIP: TYPESAFE_API_KEY is not set")
        sys.exit(2)
    body = json.dumps({"state": state, "model": "jev-latest", "questions": questions}).encode()
    for attempt in range(8):
        req = urllib.request.Request(API, data=body, headers={
            "Authorization": f"Bearer {key}", "Content-Type": "application/json"})
        try:
            with urllib.request.urlopen(req, timeout=120) as r:
                return json.loads(r.read())
        except urllib.error.HTTPError as e:
            if e.code in (429, 500, 502, 503, 504, 520, 521, 522, 523, 524, 529):
                time.sleep(2 ** attempt)
                continue
            sys.exit(f"HTTP {e.code}: {e.read().decode()[:500]}")
    sys.exit("retries exhausted")


def main():
    agents = {}
    for name, rel in AGENTS.items():
        path = os.path.join(ROOT, rel)
        if not os.path.exists(path):
            print(f"note: {rel} missing; probes expecting {name} become informational")
            continue
        agents[name] = parse_description(path)
    state = {"agents": agents, "requests": {f"R{i}": p[0] for i, p in enumerate(PROBES)}}
    questions = {
        f"route_R{i}": {
            "type": "choice",
            "instructions": {
                "request_id": f"R{i}",
                "question": "Which agent in `agents` should handle `requests.<request_id>`, judging only from the agent descriptions?",
            },
            "criteria": {**{n: None for n in agents},
                         "none": "No listed agent fits; a different plugin or the general assistant should handle it"},
        }
        for i in range(len(PROBES))
    }
    answers = call(state, questions)["answers"]
    failures = 0
    print(f"{'result':6} {'winner':20} {'conf':5} {'expected':20} request")
    for i, (request, expected) in enumerate(PROBES):
        a = answers[f"route_R{i}"]
        winner, conf = a["choice"], a["confidence"]
        if expected is None or (expected != "none" and expected not in agents):
            status = "info"
        elif winner == expected and conf >= THRESHOLD:
            status = "ok"
        else:
            status = "FAIL"
            failures += 1
        print(f"{status:6} {winner:20} {conf:.2f}  {expected or '-':20} {request[:70]}")
    print(f"\n{failures} failing must-pass probes (threshold {THRESHOLD})")
    sys.exit(1 if failures else 0)


if __name__ == "__main__":
    main()
```

- [ ] **Step 3: Write the tests README**

```markdown
# desktop-development tests

Two checks, both run from the repository root.

## lint-agents.sh

Structural checks on every agent file: a "Use when" description under 120 words with a
"Do not use for … (use `<sibling>`)" handoff and no example blocks, an `**Output:**` section,
two to five examples with one marked `Proactive trigger:`, a body under 3,000 characters
excluding examples, second person, and no Electron mention.

    desktop-development/tests/lint-agents.sh

## route-check.py

Sends the `wails-go-developer` description and five sibling descriptions from other plugins
to the TypeSafe Jev API with seventeen sample requests. Each request is one Choice question;
the winner must match the expected agent at confidence 0.70 or higher. Electron and Swift
requests must route to `none`. Needs `TYPESAFE_API_KEY`.

    desktop-development/tests/route-check.py

The ai-development audit is the third check and lives in that plugin:

    python3 ai-development/skills/agent-modernizer/scripts/audit-agents.py desktop-development/agents/
    python3 ai-development/skills/agent-modernizer/scripts/audit-agents.py --route desktop-development/agents/
```

- [ ] **Step 4: Run the RED baseline against the current agent**

Run:

```bash
chmod +x desktop-development/tests/lint-agents.sh desktop-development/tests/route-check.py
desktop-development/tests/lint-agents.sh desktop-development/agents/electron-go-pro.md
python3 ai-development/skills/agent-modernizer/scripts/audit-agents.py desktop-development/agents/electron-go-pro.md
```

Expected: the lint prints `FAIL electron-go-pro: description lacks 'Use when'`, `FAIL electron-go-pro: description lacks a 'Do not use' handoff`, `FAIL electron-go-pro: mentions Electron`, and possibly others. The audit prints a `description_triggering` finding (the description is a capability summary). Record the lint FAIL lines and the audit findings table in your final report; they are the baseline the rewrite must clear.

- [ ] **Step 5: Run the routing check with the old agent absent from the roster**

The roster names `wails-go-developer`, which does not exist yet, so its probes print `info`. Run:

```bash
python3 desktop-development/tests/route-check.py
```

Expected: `note: desktop-development/agents/wails-go-developer.md missing; ...`, sibling probes `ok`, Electron probes route to `none` or `frontend-developer`. Record the table; the Electron rows are the ones that must read `none` after Task 2.

- [ ] **Step 6: Commit**

```bash
git add desktop-development/tests/lint-agents.sh desktop-development/tests/route-check.py desktop-development/tests/README.md
git commit -m "add desktop-development lint and Jev routing checks with RED baseline"
```

---

### Task 2: Write wails-go-developer and delete electron-go-pro

**Agent:** agent-rewrite (opus). Load `plugin-dev:agent-development` first.

**Files:**
- Create: `desktop-development/agents/wails-go-developer.md`
- Delete: `desktop-development/agents/electron-go-pro.md`

- [ ] **Step 1: Write the agent file exactly as below**

```markdown
---
name: wails-go-developer
description: >
  Use when the user asks to build, restructure, or package a desktop app with Wails v3 and Go:
  services and generated bindings, windows and menus, system tray, notifications, global
  shortcuts, the dev loop, signing and DMG packaging. Also use when an existing Go program
  needs a desktop window. Do not use for terminal-only tools (use cli-developer), interactive
  terminal UIs (use go-tui-developer), or web apps with no desktop shell (use frontend-developer).
model: opus
color: cyan
tools: ["Read", "Write", "Edit", "Bash", "Grep", "Glob", "WebFetch"]
---

<example>
Context: User is starting a desktop application in Go
user: "I want a macOS desktop app in Go for tracking incidents, with a real window and menus"
assistant: "I'll use the wails-go-developer agent to scaffold a Wails v3 project, define the incident service with generated bindings, and wire the window and app menu."
<commentary>
New desktop app in Go with a native window: Wails v3 project layout, services, bindings, and events are this agent's core.
</commentary>
</example>

<example>
Context: An existing Wails v3 app needs background presence
user: "Add a menu bar icon and native notifications to our Wails app"
assistant: "I'll use the wails-go-developer agent to add a tray adapter with a dynamic menu and a notifications adapter, each behind a server build-tag stub, and note that notifications authorize only from a signed bundle."
<commentary>
Native macOS integration through Wails services, with the beta constraints the model would not guess.
</commentary>
</example>

<example>
Context: The app works in development and must ship
user: "Package this Wails app for distribution with signing and a DMG"
assistant: "I'll use the wails-go-developer agent to run the universal build, set up entitlements and the Developer ID identity, sign and notarize the bundle, and produce the DMG through the darwin Taskfile."
<commentary>
Packaging and signing are Taskfile-driven in Wails v3 and differ from what the docs describe at this beta.
</commentary>
</example>

<example>
Context: The assistant is working in a Go repository that already has a Cobra CLI under cmd/
user: "Can you give the tool a small window that shows the sync status?"
assistant: "A status window means a desktop shell. I'll use the wails-go-developer agent to add a Wails main that reuses the existing internal packages and keeps the CLI under cmd/ untouched."
<commentary>
Proactive trigger: the request is small, but it changes the program from a CLI to a desktop app, and the shared-internals layout must be set before the first window exists.
</commentary>
</example>

You are a senior desktop developer who builds macOS-first Wails v3 applications in Go with a React frontend and ships them signed. You default to macOS conventions and keep Windows and Linux builds working through the generated Taskfiles.

**Architecture defaults:**

- One Go process. Wails-generated TypeScript bindings are the only contract between Go and the UI: no HTTP API, no hand-written client, no duplicated types.
- Every entity is defined once in Go under `internal/domain`, which imports neither Wails nor any store.
- `app/` holds thin Wails adapters, one file per native concern (window, menu, tray, notifications, shortcuts), so beta API churn stays local.
- Nothing under `internal/` imports Wails; a CLI or TUI under `cmd/` reuses it unchanged.
- Go pushes to the UI with typed events registered at init; no polling timers. Streams carry token-by-token output.
- Local state, when needed, lives in SQLite through the Wails sqlite service; the UI never touches storage.
- OAuth runs in the system browser with a loopback listener. The webview is served under a custom `wails` scheme and cannot receive an `http(s)` redirect.
- Every native feature has a `server` build-tag stub so `task build:server` produces a browser-served build for tests and demos.

**Process:**

1. Pin the Wails tag from the `wails-v3` skill in `go.mod`, confirm `wails3 version` matches, and confirm Go 1.25 or newer.
2. Read the skill references for the features in play, then the docs page at the tag for any API about to be written. The scaffold and module source win over the docs when they disagree.
3. Scaffold with `wails3 init -t react`; never hand-roll the tree. Add Tailwind and the Vite Wails plugin after.
4. Put services in `app/`, domain types in `internal/domain`, one adapter file per native feature with its `server` stub.
5. Run `wails3 generate bindings` after any service signature change, `task dev` once, and `task build` before reporting.

**Output:** What was built, the pinned tag, the Taskfile targets that ran, and any native feature left unverified because it needs a signed bundle.

**Do Not:**

- Hand-roll an HTTP or WebSocket API between Go and the UI.
- Put business logic in `app/` adapters.
- Use `latest` or a floating version for the Wails module or CLI.
- Promise Touch Bar, the macOS Services menu, or Spotlight indexing; Wails has none of them.
- Claim notifications or the custom URL scheme work from an unsigned or unpackaged build.
```

- [ ] **Step 2: Delete the Electron agent**

```bash
git rm -q desktop-development/agents/electron-go-pro.md
```

- [ ] **Step 3: Run the lint**

Run: `desktop-development/tests/lint-agents.sh`
Expected: every line starts with `ok`. If the description exceeds 120 words or the body exceeds 3,000 characters, trim wording (not facts) until it passes.

- [ ] **Step 4: Run the audit**

Run: `python3 ai-development/skills/agent-modernizer/scripts/audit-agents.py desktop-development/agents/`
Expected: no Must fix and no Should fix rows. `Consider` rows for `Possible:` judgments are acceptable only after you read the cited lines and record in your report why the text does not support them. A `description_triggering` score of 3.0 is expected because the description carries a sibling handoff.

- [ ] **Step 5: Run the routing check**

Run: `python3 desktop-development/tests/route-check.py`
Expected: `0 failing must-pass probes`. If a Wails probe loses to a sibling, strengthen the description with the concrete phrasing from that probe rather than removing the probe.

- [ ] **Step 6: Commit**

```bash
git add desktop-development/agents/wails-go-developer.md
git commit -m "replace electron-go-pro with wails-go-developer"
```

---

### Task 3: SKILL.md, project-layout.md, and build-and-sign.md

**Agent:** skill-layout-build (sonnet). Load `plugin-dev:skill-development` first.

**Files:**
- Create: `desktop-development/skills/wails-v3/SKILL.md`
- Create: `desktop-development/skills/wails-v3/references/project-layout.md`
- Create: `desktop-development/skills/wails-v3/references/build-and-sign.md`

- [ ] **Step 1: Write SKILL.md**

```markdown
---
name: wails-v3
description: >
  This skill should be used when working with Wails v3 in Go: "wails3 init", services and
  generated bindings, typed events, streams, windows and menus, system tray, notifications,
  global shortcuts, dock badges, Taskfile builds, signing and notarization, DMG, the updater,
  server mode, or when asked "does Wails v3 support X". Also applies when a Wails project's
  pinned beta tag must be checked or bumped.
---

# Wails v3

Reference material for Wails `v3.0.0-beta.25` that the model cannot infer: the scaffold layout,
the Taskfile contract, native feature entry points, and the beta edges. The skill is the map;
the sources below are the truth.

## Pin policy

- Verified tag: `v3.0.0-beta.25` (2026-09-28). Projects pin `go.mod` and the `wails3` CLI to the
  same tag and upgrade deliberately, never with `latest`.
- Wails ships near-daily betas and has shipped breaking changes inside the beta series. A bump
  re-runs the scaffold check (`wails3 init -t react`, `task build`) and updates every
  `Verified against` block in `references/`.
- Go 1.25 or newer is required by the module.

## Where the truth is

Precedence when sources disagree: the scaffold and module source at the tag, then the docs at the
tag, then this skill. Before writing any API call, read the matching page or source file:

- Module source: `$(go env GOMODCACHE)/github.com/wailsapp/wails/v3@v3.0.0-beta.25/` (`pkg/application`, `pkg/services`, `pkg/updater`, `examples/`).
- Docs at the tag: `https://github.com/wailsapp/wails/tree/v3.0.0-beta.25/docs/mpress/content` (sparse-clone it; the live site at v3.wails.io blocks scripted fetches with HTTP 403, and its pages lag or lead the tag).
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

React, Vite, Tailwind, and Go idioms the model already knows; iOS and Android; Electron.
```

- [ ] **Step 2: Write references/project-layout.md**

```markdown
# Project Layout and Build Contract

**Verified against:** `v3.0.0-beta.25`, 2026-09-28, by running `wails3 init -n wailsprobe -t react -d wailsprobe -mod example.com/wailsprobe` and `task build` (9.3 MB binary, about 8 seconds including `npm install`). Sources: the generated files, `guides/dev/project-structure.md`, `concepts/build-system.md`, `guides/server-build.md`, `reference/cli.md`.

## Scaffold

`wails3 init -n <name> -t react -d <dir> -mod <module> -productidentifier <com.example.app> -productcompany "<Company>"`.
`-d` plus `-n` nests the project at `<dir>/<name>`. Built-in templates: `react`, `react-js`, `svelte`, `svelte-js`, `vue`, `vue-js`, `vanilla` (default), `vanilla-js`, `ios`. There is no `react-ts`; `react` is the TypeScript template.

```
main.go            application.New, window, app.Run
greetservice.go    example service; delete it
go.mod             go 1.25.0; require github.com/wailsapp/wails/v3 v3.0.0-beta.25
Taskfile.yml       root tasks; includes build/Taskfile.yml and the per-platform files
build/
  config.yml       info (product name, identifier, version), dev_mode watcher, protocols, fileAssociations
  Taskfile.yml     common tasks: frontend deps, bindings, icons, server build, docker
  darwin/          Info.plist, Info.dev.plist, icons.icns, Assets.car, dmg assets, Taskfile.yml
  windows/         icon.ico, info.json, wails.exe.manifest, nsis/, msix/, Taskfile.yml
  linux/           desktop, appimage/, nfpm/, Taskfile.yml
  ios/ android/ docker/
frontend/
  package.json     @wailsio/runtime, react 18, vite 8, typescript 5; scripts dev, build:dev, build
  vite.config.ts   add the Wails plugin here (see Typed events below)
  src/             App.tsx, main.tsx
  bindings/        generated by wails3 generate bindings; committed
  dist/            built frontend, embedded by main.go; gitignored
bin/               build output; gitignored
```

`.gitignore` covers `.task`, `bin`, `frontend/dist`, `frontend/node_modules`, and the mobile overlays. It does not ignore `frontend/bindings`; keep bindings committed so a checkout builds without running the generator first.

## main.go shape

```go
//go:embed all:frontend/dist
var assets embed.FS

func init() {
	application.RegisterEvent[string]("time") // typed events register at init
}

func main() {
	app := application.New(application.Options{
		Name:        "wailsprobe",
		Description: "…",
		Services:    []application.Service{application.NewService(&GreetService{})},
		Assets:      application.AssetOptions{Handler: application.AssetFileServerFS(assets)},
		Mac:         application.MacOptions{ApplicationShouldTerminateAfterLastWindowClosed: true},
	})
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Window 1", Width: 1000, Height: 618,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(6, 7, 15),
		URL:              "/",
	})
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
```

Set `ApplicationShouldTerminateAfterLastWindowClosed: false` for tray-resident apps; with `true`, hiding the only window from a global shortcut quits the app.

## Recommended package layout

```
main.go              application.New, registers services, calls startDesktop
main_desktop.go      //go:build !server  — window, menu, tray, shortcuts, notifications
main_server.go       //go:build server   — startDesktop is a no-op
app/                 thin Wails adapters, one file per native concern; services live here
internal/domain/     entities and rules; imports neither Wails nor a store
internal/store/      persistence behind interfaces
cmd/<tool>/          optional CLI or TUI entry points that import internal/ only
```

The build-tag split, from the owner's mock-up:

```go
// main_desktop.go
//go:build !server

func startDesktop(app *application.App, deps Deps) func() {
	win := app.Window.NewWithOptions(mainWindowOptions())
	app.Menu.Set(buildMenu(app, win))
	tray := startTray(app, win)
	stopShortcuts := startShortcuts(app, win)
	stopNotify := startNotifications(app, deps)
	return func() { stopNotify(); stopShortcuts(); tray.Stop() }
}
```

```go
// main_server.go
//go:build server

func startDesktop(*application.App, Deps) func() { return func() {} }
```

`task build:server` compiles with `-tags server,production`; the result serves the same frontend and bindings over HTTP (`ServerOptions{Host, Port}`, default `localhost:8080`, `/health` endpoint) with each browser tab as a window named `browser-N`. Server mode has no tray, no native dialogs, and no native windows.

## Taskfile targets

Run with `task <name>` or `wails3 task <name>`.

| Target | Does |
| --- | --- |
| `dev` | `wails3 dev -config ./build/config.yml -port 9245`; rebuilds Go on change, Vite HMR for the frontend |
| `build` | `{{GOOS}}:build`; on macOS `darwin:build:native` with `-tags production -trimpath -ldflags="-w -s"`, `CGO_ENABLED=1`, deployment target 12.0 |
| `package` | `{{GOOS}}:package`; on macOS creates `bin/<APP_NAME>.app` and ad-hoc signs it |
| `run` | on macOS builds `bin/<APP_NAME>.dev.app` with `Info.dev.plist` and launches it |
| `build:server`, `run:server` | `-tags server` binary at `bin/<APP_NAME>-server` |
| `build:docker`, `run:docker`, `setup:docker` | container images for server mode and cross-compilation |
| `common:generate:bindings` | `wails3 generate bindings -f '<build flags>' -clean=true -ts -i`; runs before every frontend build |
| `common:generate:icons` | `.icns` and `.ico` from `build/appicon.png`, `Assets.car` from `build/appicon.icon` |
| `common:update:build-assets` | regenerates `Info.plist`, NSIS, `.desktop` from `build/config.yml`; overwrites hand edits |
| `darwin:build:universal`, `darwin:package:universal`, `darwin:package:dmg`, `darwin:sign`, `darwin:sign:notarize` | see `build-and-sign.md` |

Add a root shortcut for the generator the way the mock-up does:

```yaml
  bindings:
    summary: Regenerate TypeScript bindings
    cmds:
      - wails3 generate bindings -ts -i -clean -d frontend/bindings
```

Root vars: `APP_NAME`, `BIN_DIR` (`bin`), `PACKAGE_MANAGER` (`npm`; `bun`, `pnpm`, `yarn` supported), `VITE_PORT` (`9245`), `GOOS`.

## build/config.yml

```yaml
info:
  companyName: "…"
  productName: "…"
  productIdentifier: "com.example.app"   # CFBundleIdentifier; also the notifications and single-instance id
  description: "…"
  copyright: "…"
  version: "0.1.0"
dev_mode:
  root_path: .
  log_level: warn
  debounce: 1000
  ignore:
    dir: [.git, node_modules, frontend, bin]
    file: [.DS_Store, .gitignore, .gitkeep, "*_test.go"]
    watched_extension: ["*.go", "*.js", "*.ts"]   # add "*.json" if Go embeds JSON fixtures
    git_ignore: true
  executes:
    - cmd: wails3 build DEV=true
      type: blocking
    - cmd: wails3 task common:dev:frontend
      type: background
    - cmd: wails3 task run
      type: primary
      exit_policy: shutdown
protocols:            # custom URL schemes; consumed by packagers, not by application.Options
  - scheme: myapp
    description: "My App links"
fileAssociations: []
```

After editing `info` or `protocols`, run `task common:update:build-assets`.

## Bindings and typed events

- Services are plain structs registered with `application.NewService(&T{})`. Exported methods become bindings; a `context.Context` first parameter gives the frontend a cancellable promise; returned `error` becomes a rejected promise.
- Optional lifecycle: `ServiceStartup(ctx, options) error` (registration order; an error aborts startup), `ServiceShutdown() error` (reverse order), `ServiceName() string`.
- Generated files: `frontend/bindings/<module>/<service>.ts` and `index.ts`; typed events in `frontend/bindings/github.com/wailsapp/wails/v3/internal/eventcreate.ts` and `eventdata.d.ts`.
- Typed events: `application.RegisterEvent[T]("name")` in `init()`; emit with `app.Event.Emit("name", value)` or `window.EmitEvent`; `-tags strictevents` warns on unregistered names. Frontend: add `wails()` from `@wailsio/runtime/plugins/vite` to `vite.config.ts` plugins, then `Events.On(TypedCreator, handler)` from `@wailsio/runtime`.
- Frontend call: `import { Greet } from "../bindings/<module>/greetservice"` then `await Greet(name)`.

## Adding React tooling the template omits

`npm i -D tailwindcss postcss autoprefixer` and the usual `tailwind.config.ts` and `postcss.config.js`; `npm i react-router-dom zustand` for per-window hash routes and state. Use hash routing: every window loads `/` from the embedded assets and picks its screen from `#/route`.
```

- [ ] **Step 3: Write references/build-and-sign.md**

```markdown
# Build, Sign, and Distribute

**Verified against:** `v3.0.0-beta.25`, 2026-09-28. Sources: the scaffold's `build/darwin/Taskfile.yml`, `guides/build/macos.md`, `guides/build/signing.md`, `guides/updater.md`, `reference/cli.md`. Signing and notarization were not executed in this session (no Developer ID present); commands are transcribed from the Taskfile and marked where a signed bundle is needed to prove behavior.

## macOS

### Build and package

```bash
task build                              # bin/<APP_NAME>, host arch, production flags
task darwin:build:universal             # amd64 + arm64, lipo into bin/<APP_NAME>
task package                            # bin/<APP_NAME>.app, ad-hoc signed
task darwin:package:universal           # universal .app
task darwin:package:dmg                 # styled .dmg from the .app (macOS only)
```

The `.app` contains `Contents/MacOS/<APP_NAME>`, `Contents/Resources/icons.icns` (and `Assets.car` when present), and `Contents/Info.plist` copied from `build/darwin/Info.plist`. `Info.plist` and `Info.dev.plist` are regenerated by `task common:update:build-assets` from `build/config.yml`; edit `config.yml`, not the plist, unless you also stop running that task.

Extra read-only files ship under `Contents/Resources/`: add a copy step to `darwin:create:app:bundle` and read them with `mac.LoadResource("path")` or `mac.ResourceFS()` from `github.com/wailsapp/wails/v3/pkg/mac`. Both return `mac.ErrNotInAppBundle` when the binary runs outside a bundle, so `task dev` and `bin/<APP_NAME>` cannot see them. Never write into the bundle after signing.

DMG vars in `build/darwin/Taskfile.yml`: `DMG_BACKGROUND`, `DMG_VOLUME_ICON`, `DMG_FILE_ICON`, `DMG_WINDOW_WIDTH` (540), `DMG_WINDOW_HEIGHT` (380), `DMG_FILES` (`"Name=path,..."`).

### Entitlements

The scaffold ships no entitlements file. Generate both with the wizard:

```bash
wails3 setup entitlements    # writes build/darwin/entitlements.plist and entitlements.dev.plist
```

Development needs JIT, unsigned executable memory, and debugging; production needs only network. The App Store preset enables the sandbox.

### Signing and notarization

At this tag the darwin Taskfile does **not** read `SIGN_IDENTITY` or `KEYCHAIN_PROFILE` vars, whatever the docs say. `darwin:sign` and `darwin:sign:notarize` call `wails3 tool sign --input bin/<APP_NAME>.app [--notarize] {{.CLI_ARGS}}`, and `wails3 tool sign` takes identity and keychain profile from `wails3 setup signing` (stored in `~/.config/wails/defaults.yaml`) or from flags.

```bash
security find-identity -v -p codesigning          # confirm a "Developer ID Application" identity
xcrun notarytool store-credentials "my-profile" \
  --apple-id you@example.com --team-id TEAMID --password app-specific-password
wails3 setup signing                              # records identity and profile for wails3 tool sign
task darwin:sign                                  # package, then sign with hardened runtime
task darwin:sign:notarize                         # package, sign, submit, staple
task darwin:sign -- --identity "Developer ID Application: Co (TEAMID)" --entitlements build/darwin/entitlements.plist
spctl --assess --verbose=2 "bin/<APP_NAME>.app"   # verify
```

Unverified in this session: that `wails3 tool sign` staples the ticket and that notarization succeeds with the wizard's production entitlements. Say so in reports until a signed build has been produced.

"App is damaged and can't be opened" means the bundle is unsigned or its signature broke; `xattr -cr <app>` clears quarantine for local testing only.

### Private macOS APIs

`-tags private_mac_apis` (or `wails3 build -tags private_mac_apis`, `EXTRA_TAGS=private_mac_apis wails3 dev`) enables webview transparency behind translucent and transparent backdrops, Liquid Glass grouping, and programmatic inspector opening. Without the tag the options are accepted and the private parts are no-ops; the webview stays opaque. A CSS `backdrop-filter` over the app's own gradient needs none of this.

## Updater

```go
gh, err := github.New(github.Config{Repository: "org/app", ChecksumAsset: "SHA256SUMS"})
if err := app.Updater.Init(updater.Config{
	CurrentVersion: version,
	Providers:      []updater.Provider{gh},
	PublicKey:      publicKey,          // ed25519 public key embedded with go:embed; sign releases with the private key
	CheckInterval:  6 * time.Hour,      // optional periodic check
}); err != nil { … }
go app.Updater.CheckAndInstall(ctx)    // opens the default update window, downloads, verifies, stages
```

Imports: `github.com/wailsapp/wails/v3/pkg/updater` and `.../pkg/updater/providers/github` (also `keygen`, `appcast`, `endpoint`). States: `idle`, `checking`, `up-to-date`, `available`, `downloading`, `verifying`, `installing`, `ready`, `error`. Events: `updater.EventUpdateAvailable`, `EventDownloadProgress`, `EventUpdateReady`, `EventError`. `Restart` waits for the helper to reach `application.New`; raise `HelperReadyTimeout` if startup is slow. The GitHub provider picks the asset by `GOOS`/`GOARCH` substrings; a universal `.app` zip needs a custom `AssetMatcher`. Releases must be signed and the bundle must be Developer ID signed for the swapped binary to launch.

## Windows and Linux (Taskfile level)

| Target | Output |
| --- | --- |
| `task build GOOS=windows`, `task package GOOS=windows` | `bin/<APP_NAME>.exe`, NSIS installer from `build/windows/nsis/` |
| `task windows:package:msix` (after `task install:msix:tools`) | MSIX from `build/windows/msix/` |
| `task build GOOS=linux`, `task package GOOS=linux` | binary, AppImage from `build/linux/appimage/`, deb and rpm via nfpm from `build/linux/nfpm/` |
| `task setup:docker` then `task build GOOS=darwin` on Linux or Windows | cross-compiled macOS binary, unsigned; sign on a Mac |

`build/config.yml` `protocols` and `fileAssociations` feed the NSIS macros, the MSIX manifest, `CFBundleURLTypes`, and the `.desktop` file when `common:update:build-assets` runs. Windows signing uses `wails3 tool sign --certificate` or `--thumbprint`; Linux packages sign with a PGP key via `--pgp-key`.
```

- [ ] **Step 4: Check every path and target against the scaffold**

Run from the scaffold directory:

```bash
cd /private/tmp/claude-501/-Users-bill-Projects-mynet/8934198f-2ad7-4e6f-823c-41527251b7fd/scratchpad/wailsprobe/wailsprobe
task --list-all | sort > /tmp/targets.txt
for t in build package run dev build:server run:server common:generate:bindings common:generate:icons common:update:build-assets darwin:build:universal darwin:package:universal darwin:package:dmg darwin:sign darwin:sign:notarize; do grep -q "^\* $t:" /tmp/targets.txt && echo "ok $t" || echo "MISSING $t"; done
ls build/darwin build/config.yml frontend/vite.config.ts frontend/bindings
```

Expected: every target prints `ok`; every path exists. Fix any reference that names something missing.

- [ ] **Step 5: Commit**

```bash
git add desktop-development/skills/wails-v3/SKILL.md desktop-development/skills/wails-v3/references/project-layout.md desktop-development/skills/wails-v3/references/build-and-sign.md
git commit -m "add wails-v3 skill with project layout and build-and-sign references"
```

---

### Task 4: native-features.md and beta-gotchas.md

**Agent:** skill-native-gotchas (opus). Load `plugin-dev:skill-development` first. Before committing, dispatch `go-architect` to review every Go snippet in `native-features.md` against `~/go/pkg/mod/github.com/wailsapp/wails/v3@v3.0.0-beta.25/pkg/application` and `pkg/services`, and apply its corrections.

**Files:**
- Create: `desktop-development/skills/wails-v3/references/native-features.md`
- Create: `desktop-development/skills/wails-v3/references/beta-gotchas.md`

- [ ] **Step 1: Write references/native-features.md**

```markdown
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

Other options: `Frameless`, `AlwaysOnTop`, `Hidden`, `HideOnFocusLost`, `HideOnEscape` (tray popups), `StartState`, `KeyBindings map[string]func(application.Window)`, `DevToolsEnabled`. Mac: `Backdrop` (`MacBackdropNormal`, `Translucent`, `Transparent`, `LiquidGlass`; the last three need `private_mac_apis` for webview transparency), `WindowLevel`, `CollectionBehavior` (bitmask; `CanJoinAllSpaces | FullScreenAuxiliary` for Spotlight-style panels), `WindowClass: MacWindowClassPanel` with `PanelPreferences{NonActivating: true}` for a panel that does not activate the app.

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

Streams (`app.HandleStream`, `guides/streams.md`) carry byte streams with backpressure over the asset server without a TCP port; use them for token-by-token model output. Dialogs: `app.Dialog.Info()`, `.Question()`, `.Error()`, file open and save with `.SetTitle().SetMessage().AddButton().Show()`; not available in server mode.

## Single instance and URL schemes

```go
SingleInstance: &application.SingleInstanceOptions{
	UniqueID: "com.example.app",
	OnSecondInstanceLaunch: func(d application.SecondInstanceData) { /* d.Args, d.WorkingDir; focus the window */ },
	EncryptionKey: key,                             // 32 bytes; without it treat d.Args as untrusted
},
app.Event.OnApplicationEvent(events.Common.ApplicationLaunchedWithUrl, func(e *application.ApplicationEvent) { url := e.Context().URL() })
```

Declare schemes under `protocols:` in `build/config.yml` and run `task common:update:build-assets`; `application.Options` has no `Protocols` field. Constraint: scheme launches only work from a packaged `.app` whose `Info.plist` carries `CFBundleURLTypes`; handle both the cold path (`ApplicationLaunchedWithUrl`) and the warm path (`OnSecondInstanceLaunch` with the URL in `Args`).

## Server mode

`-tags server` keeps services, bindings, events, and streams; drops windows, tray, dialogs, shortcuts, notifications. `Server: application.ServerOptions{Host, Port, TLS, WebSocketOriginPatterns}`; `/health` returns `{"status":"ok"}`. Use it for Playwright tests and demos in a browser.
```

- [ ] **Step 2: Write references/beta-gotchas.md**

```markdown
# Beta Gotchas

**Verified against:** `v3.0.0-beta.25`, 2026-09-28. Each entry names the behavior, the consequence, and the workaround, with its source. Re-check every entry on a tag bump; delete the ones that stop being true.

## Version and drift

- **Go 1.25 or newer.** `go.mod` from the scaffold declares `go 1.25.0`. An older toolchain fails at `go mod tidy`. Source: scaffold.
- **Breaking changes ship inside the beta series** (macOS coordinate normalization, template removals, WebView2 loader changes across recent betas). Pin `go.mod` and the CLI to one tag; run `wails3 update cli` only as a deliberate bump. Source: owner's rebuild analysis R4; `release_notes.md` in the module.
- **The docs site blocks scripted reads.** `curl` and fetch tools get HTTP 403 from v3.wails.io pages; only `llms.txt` loads. Read `docs/mpress/content` from the wails repo at the tag. Source: checked 2026-09-28.
- **Docs and tag disagree on signing configuration.** Docs say set `SIGN_IDENTITY`, `KEYCHAIN_PROFILE`, `ENTITLEMENTS` in `build/darwin/Taskfile.yml`; the scaffold's `sign` tasks call `wails3 tool sign` and take identity from `wails3 setup signing` or `CLI_ARGS`. Follow the scaffold. Source: `build/darwin/Taskfile.yml` versus `guides/build/signing.md`.
- **Docs place `config.yml` at the root; the scaffold puts it at `build/config.yml`.** Source: `guides/dev/project-structure.md` versus scaffold.
- **`TitleBar` appears as both a constant and a struct in the docs.** The scaffold uses `application.MacTitleBarHiddenInset` and builds; prefer that spelling. Source: scaffold `main.go` versus `features/windows/options.md`.
- **There is no `react-ts` template.** `react` is the TypeScript template; `react-js` is JavaScript. Source: `wails3 init -l`.

## Webview

- **Custom `wails` scheme.** WKWebView serves the app through `setURLSchemeHandler:forURLScheme:@"wails"`, so the page origin is not `http(s)`. OAuth redirect URIs must be `http(s)`, so run OAuth in the system browser (`app.Browser.OpenURL`) with a loopback listener in Go and never inside the webview. Source: `pkg/application/webview_window_darwin.go:181`.
- **Test `localStorage` and cookies per app.** Persistence under the custom scheme was not verified here; keep app state in Go and treat web storage as a cache. Source: owner's rebuild analysis R4 A9.
- **Safari feature set.** WKWebView, not Chromium; no Chromium-only APIs; test `contenteditable` early if you build an editor. Source: owner's analysis.
- **Translucency needs a private-API build.** `MacBackdropTranslucent`, `Transparent`, and `LiquidGlass` keep the webview opaque unless built with `-tags private_mac_apis`. CSS `backdrop-filter` over your own gradient needs no tag. Source: `features/windows/options.md`, `guides/build/macos.md`.

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
```

- [ ] **Step 3: Check every Go identifier against the module source**

Run:

```bash
M=~/go/pkg/mod/github.com/wailsapp/wails/v3@v3.0.0-beta.25
for id in MacTitleBarHiddenInset InvisibleTitleBarHeight HideOnFocusLost HideOnEscape MacWindowClassPanel ActivationPolicyAccessory SetTemplateIcon AttachWindow WindowOffset WindowDebounce RegisterNotificationCategory SendNotificationWithActions OnNotificationResponse RegisterEvent OnApplicationEvent ApplicationLaunchedWithUrl SingleInstanceOptions OnSecondInstanceLaunch GlobalShortcut KeyBinding SetBadge HideAppIcon HandleStream; do
  grep -rq "$id" $M/pkg && echo "ok $id" || echo "MISSING $id"
done
grep -rn 'func (m \*Menu) Add\(Role\|Separator\|Checkbox\|Radio\|Submenu\)' $M/pkg/application/menu.go | head
```

Expected: every identifier prints `ok`. Anything `MISSING` is corrected from the source or removed. Then dispatch `go-architect` with: "Review every Go snippet in `desktop-development/skills/wails-v3/references/native-features.md` against the Wails module at `$M`. Report each identifier, method receiver, or field that does not exist at this tag, with the correct spelling. Do not edit files." Apply its corrections.

- [ ] **Step 4: Commit**

```bash
git add desktop-development/skills/wails-v3/references/native-features.md desktop-development/skills/wails-v3/references/beta-gotchas.md
git commit -m "add wails-v3 native-features and beta-gotchas references"
```

---

### Task 5: TypeSafe quality loop on the skill

**Agent:** harness-and-close-out (sonnet). Load `typesafe:typesafe-ai` first. Starts after Tasks 3 and 4 are committed.

**Files:**
- Modify (only if findings are confirmed): any file under `desktop-development/skills/wails-v3/`
- Scratch: `/private/tmp/claude-501/-Users-bill-Projects-mynet/8934198f-2ad7-4e6f-823c-41527251b7fd/scratchpad/skill-check.py` (not committed)

- [ ] **Step 1: Write the question script in the scratchpad**

```python
#!/usr/bin/env python3
"""Jev quality judgments over the wails-v3 skill files. Not part of the plugin."""
import json, os, sys, time, urllib.error, urllib.request

API = "https://api.typesafe.ai/v1/systemone"
ROOT = os.path.expanduser("~/Projects/mynet/desktop-development/skills/wails-v3")
FILES = ["SKILL.md", "references/project-layout.md", "references/native-features.md",
         "references/build-and-sign.md", "references/beta-gotchas.md"]

QUESTIONS = {
    "topic_lists": {"type": "noul",
        "instructions": "Does `file.body` contain lists whose items only name a feature, API, or technology without stating what to do, prefer, or avoid?",
        "criteria": {"true": {"what": "Bare inventories such as '- Windows\\n- Menus\\n- Tray' with no decision or constraint attached."},
                     "false": {"what": "Every list item carries a command, a default, a constraint, or a consequence."}}},
    "teaches_general": {"type": "noul",
        "instructions": "Does `file.body` spend lines explaining general React, Vite, Tailwind, Go, shell, or git knowledge that a senior developer already has, rather than Wails-specific facts?",
        "criteria": {"true": {"examples": ["how useState works", "what go mod tidy does", "how to install npm packages in general"]},
                     "false": {"what": "Mentions of general tools are one-line pointers in service of a Wails-specific step."}}},
    "unsourced_claims": {"type": "noul",
        "instructions": "Does `file.body` state a fact about Wails behavior, a command, a path, or a limitation that is not covered by the 'Verified against' block or by a named source (a docs page, a source file, the scaffold, or the owner's analysis)?",
        "criteria": {"true": {"what": "A behavioral claim with no traceable source, for example a limit, a default, or a 'does not work' statement without a Source or page name."},
                     "false": {"what": "Every behavioral claim is inside a section whose sources are named in the Verified against block, or carries its own Source note."}}},
    "hedged_or_vague": {"type": "noul",
        "instructions": "Does `file.body` contain vague or hedged guidance such as 'consider', 'may want to', 'as appropriate', or 'handle edge cases' where a concrete rule could be stated?",
        "criteria": {"true": {"examples": ["Consider adding error handling as appropriate."]},
                     "false": {"what": "Guidance names the action, the default, or the exact condition."}}},
    "decision_density": {"type": "score",
        "instructions": "Judging `file.body`, how much of its content is commands, defaults, constraints, and consequences rather than description?",
        "criteria": ["Mostly descriptive prose about what Wails is or offers.",
                     "Mixed; several sections describe features without a rule or constraint.",
                     "Mostly rules, commands, and constraints; a few descriptive passages.",
                     "Nearly every paragraph is a command, a default, a constraint, or a consequence."]},
}
SKILL_ONLY = {
    "description_triggering": {"type": "score",
        "instructions": "Judging only `file.description`, how well does it tell a router when to load this skill?",
        "criteria": ["Summarizes what the skill contains with no triggering phrases.",
                     "Names a domain but no concrete requests or phrasings.",
                     "Names concrete requests, phrasings, or situations.",
                     "Names concrete triggers and also says when not to load it."]},
}


def call(state, questions):
    key = os.environ.get("TYPESAFE_API_KEY") or sys.exit("TYPESAFE_API_KEY not set")
    body = json.dumps({"state": state, "model": "jev-latest", "questions": questions}).encode()
    for attempt in range(8):
        req = urllib.request.Request(API, data=body, headers={"Authorization": f"Bearer {key}", "Content-Type": "application/json"})
        try:
            with urllib.request.urlopen(req, timeout=120) as r:
                return json.loads(r.read())["answers"]
        except urllib.error.HTTPError as e:
            if e.code in (429, 500, 502, 503, 504, 520, 521, 522, 523, 524, 529):
                time.sleep(2 ** attempt); continue
            sys.exit(f"HTTP {e.code}: {e.read().decode()[:300]}")
    sys.exit("retries exhausted")


def split_frontmatter(text):
    if not text.startswith("---"):
        return "", text
    end = text.index("\n---", 3)
    return text[3:end], text[end + 4:]


for rel in FILES:
    text = open(os.path.join(ROOT, rel), encoding="utf-8").read()
    fm, body = split_frontmatter(text)
    state = {"file": {"path": rel, "description": fm, "body": body}}
    qs = dict(QUESTIONS)
    if rel == "SKILL.md":
        qs.update(SKILL_ONLY)
    answers = call(state, qs)
    print(f"\n== {rel}")
    for qid, a in answers.items():
        if a["type"] == "noul":
            p = a["noul"]; flag = "FLAG" if p >= 0.7 else ("check" if p >= 0.4 else "ok")
            print(f"{flag:5} {qid:22} p={p:.2f}")
        else:
            print(f"{'':5} {qid:22} score={a['score']:.2f}/3 conf={a['confidence']:.2f}")
```

- [ ] **Step 2: Run it before any edits and record the output**

Run: `python3 /private/tmp/claude-501/-Users-bill-Projects-mynet/8934198f-2ad7-4e6f-823c-41527251b7fd/scratchpad/skill-check.py`
Expected: a block per file. Save the output to `.../scratchpad/skill-check-before.txt` and paste it into your final report.

- [ ] **Step 3: Confirm every FLAG and check line by reading the file**

For each `FLAG` (p at or above 0.7) and each `check` (0.4 to 0.7), open the file, find the lines that decide it, and record: confirmed (with the lines) or dropped (with why the text does not support it). For confirmed `unsourced_claims`, either add the source or delete the claim. For confirmed `topic_lists` or `hedged_or_vague`, rewrite the lines as a rule or command. For `decision_density` below 2.0 on any reference, rewrite descriptive passages as rules. For `description_triggering` below 2.5 on `SKILL.md`, add concrete phrasings to the description while keeping it under 1,024 characters.

- [ ] **Step 4: Rerun and compare**

Run the script again and save to `skill-check-after.txt`. Expected: no `FLAG` lines remain; every `check` line is either gone or has a recorded reading that dismisses it. Put both outputs side by side in your report.

- [ ] **Step 5: Commit any edits**

```bash
git add desktop-development/skills/wails-v3/SKILL.md desktop-development/skills/wails-v3/references/project-layout.md desktop-development/skills/wails-v3/references/native-features.md desktop-development/skills/wails-v3/references/build-and-sign.md desktop-development/skills/wails-v3/references/beta-gotchas.md
git commit -m "tighten wails-v3 skill after Jev quality review"
```

If nothing changed, skip the commit and say so.

---

### Task 6: Rewrite the Diátaxis pages

**Agent:** harness-and-close-out (sonnet), dispatching the Diátaxis agents. Starts after Tasks 2 through 5 are committed.

**Files:**
- Modify: `desktop-development/docs/_index.md`
- Modify: `desktop-development/docs/tutorials/getting-started.md`
- Delete: `desktop-development/docs/howto/set-up-electron-go-project.md`
- Create: `desktop-development/docs/howto/set-up-wails-go-project.md`
- Modify: `desktop-development/docs/reference/agents.md`
- Create: `desktop-development/docs/reference/skills.md`
- Modify: `desktop-development/docs/explanation/architecture.md`

- [ ] **Step 1: Dispatch four Diátaxis agents concurrently with these prompts**

Common preamble for every prompt: "Ground truth is `desktop-development/agents/wails-go-developer.md`, `desktop-development/skills/wails-v3/SKILL.md`, its four references, and `docs/specs/2026-09-28-desktop-development-wails-v3-design.md`. Read them first. Keep Hugo front matter (`title`, `description`, `weight`, and `bookCollapseSection` where present). No Electron anywhere except one sentence in the explanation page's history. No authorship lines, no timeline estimates, no git commands. Edit only the file named."

`diataxis-docs:doc-tutorial-writer` (sonnet) for `docs/tutorials/getting-started.md`: "Rewrite the tutorial so the reader builds the same note-taking app with Wails v3: install `wails3` at `v3.0.0-beta.25` and Go 1.25+, run `wails3 init -n notes -t react -d notes -mod example.com/notes`, replace `greetservice.go` with a `NoteService` under `app/` whose methods take `context.Context` first and return `error`, register one typed event with `application.RegisterEvent[Note]("note:saved")` in `init()`, run `task dev`, call the generated binding from `frontend/src/App.tsx`, and finish with `task build` producing `bin/notes`. Keep the same length and checkpoint structure as the current page; the agent name is `wails-go-developer`."

`diataxis-docs:doc-howto-writer` (sonnet) for `docs/howto/set-up-wails-go-project.md` (new; delete the Electron how-to yourself before dispatch with `git rm -q desktop-development/docs/howto/set-up-electron-go-project.md`): "Write a how-to titled 'Set Up a Wails v3 Go Project' with front matter `title`, `description` ('Scaffold, add services and events, wire native features, and package a Wails v3 app'), `weight: 1`. Steps: scaffold with the `react` template; add Tailwind and the Wails Vite plugin; add a service under `app/` and regenerate bindings with `wails3 generate bindings -ts -i -clean -d frontend/bindings`; register a typed event; add tray and notifications adapters each with a `//go:build server` stub; run `task build`; sign with `wails3 setup entitlements`, `wails3 setup signing`, `task darwin:sign:notarize`, and `task darwin:package:dmg`. Troubleshooting: Go version too old, notifications never appear (needs a packaged signed bundle), OAuth redirect cannot reach the webview (use the system browser), stale bindings after a signature change, `task common:update:build-assets` overwrote a plist edit. Under 1,800 words."

`diataxis-docs:doc-reference-gen` (sonnet) for `docs/reference/agents.md` and `docs/reference/skills.md` (new, front matter `title: "Skills"`, `description: "Wails v3 skill specification"`, `weight: 2`) and `docs/_index.md`: "Regenerate `agents.md` from `wails-go-developer.md` (fields table, trigger conditions from the description, architecture defaults, process, output, do-not list, delegation table naming go-architect, react-specialist, playwright-expert). Generate `skills.md` from `SKILL.md` and the references: trigger phrases, verified tag, precedence rule, reference file table with one line each, and a table of the beta gotchas by title. Update `_index.md`'s component table (agent `wails-go-developer`, skill `wails-v3`) and its links (the how-to is now `howto/set-up-wails-go-project/`, add `reference/skills/`). Pure specification, no advice."

`diataxis-docs:doc-explanation-writer` (opus) for `docs/explanation/architecture.md`: "Rewrite as 'Why Wails v3 and Go' with front matter `description: "One process, generated bindings, and a macOS-first build pipeline"`. Sections: why one Go process with generated bindings instead of two processes and an IPC layer (the owner's rebuild analysis measured a 283 MB Electron bundle against a 10 MB Wails app and found about 80 percent of the old code was boundary plumbing or dead); the adapter boundary (`app/` thin, `internal/` never imports Wails, `server` build tag) and what it buys for CLI and TUI reuse and browser-based tests; why macOS first (signing, notarization, notifications needing a signed bundle) while Windows and Linux stay at Taskfile level; the beta trade-off and the pin policy with source precedence (scaffold and module, then docs at the tag, then the skill); when not to use this plugin (terminal-only tools go to cli-development, web apps to web-development, native Swift utilities to mobile-development); cross-plugin collaboration. One sentence may note that the plugin previously targeted Electron and dropped it in 2.0.0."

- [ ] **Step 2: Validate**

Dispatch `diataxis-docs:doc-crosslink-validator` (sonnet) over `desktop-development/docs` with: "Check front matter, type separation, that every relative link resolves under the Hugo section layout, that no page says examples belong in a description, that the how-to is named `set-up-wails-go-project`, that the tag `v3.0.0-beta.25` and the template name `react` are consistent, and that no page mentions Electron outside the explanation page's history sentence. Report file and line for every problem." Fix every reported problem yourself.

- [ ] **Step 3: Commit**

```bash
git add desktop-development/docs/_index.md desktop-development/docs/tutorials/getting-started.md desktop-development/docs/howto/set-up-wails-go-project.md desktop-development/docs/reference/agents.md desktop-development/docs/reference/skills.md desktop-development/docs/explanation/architecture.md
git commit -m "rewrite desktop-development docs for Wails v3"
```

(`git rm` already staged the deletion of the Electron how-to.)

---

### Task 7: Changelog, manifests, and version

**Agent:** harness-and-close-out (sonnet)

**Files:**
- Create: `desktop-development/CHANGELOG.md`
- Modify: `desktop-development/.claude-plugin/plugin.json`
- Modify: `.claude-plugin/marketplace.json` (the `desktop-development` entry only)

- [ ] **Step 1: Write the changelog**

```markdown
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
```

- [ ] **Step 2: Update plugin.json**

```json
{
  "name": "desktop-development",
  "version": "2.0.0",
  "description": "Wails v3 desktop application development in Go with a React frontend, macOS-first signing and packaging",
  "keywords": [
    "desktop",
    "go",
    "macos",
    "native",
    "wails",
    "webview"
  ]
}
```

- [ ] **Step 3: Update the marketplace entry**

Edit only the object whose `"name"` is `"desktop-development"` in `.claude-plugin/marketplace.json`: set `"version": "2.0.0"`, the same `description` as above, and the same keyword list. Run `python3 -c "import json;json.load(open('.claude-plugin/marketplace.json'))"` to confirm it still parses.

- [ ] **Step 4: Commit**

```bash
git add desktop-development/CHANGELOG.md desktop-development/.claude-plugin/plugin.json .claude-plugin/marketplace.json
git commit -m "bump desktop-development to 2.0.0 with changelog"
```

---

### Task 8: Final verification

**Agent:** harness-and-close-out (sonnet)

**Files:** none modified unless a check fails.

- [ ] **Step 1: Run every check**

```bash
cd ~/Projects/mynet
desktop-development/tests/lint-agents.sh
python3 desktop-development/tests/route-check.py
python3 ai-development/skills/agent-modernizer/scripts/audit-agents.py desktop-development/agents/
python3 ai-development/skills/agent-modernizer/scripts/audit-agents.py --route desktop-development/agents/
grep -rniE 'electron' desktop-development --include='*.md' --include='*.json' --include='*.sh' --include='*.py' | grep -v CHANGELOG.md | grep -v 'explanation/architecture.md'
```

Expected: lint all `ok`; routing `0 failing must-pass probes`; audit no Must fix and no Should fix; `--route` `0 of N non-proactive example requests do not route`; the grep prints nothing.

- [ ] **Step 2: Validate the plugin**

Dispatch `plugin-dev:plugin-validator` (haiku) over `desktop-development` with the note that example blocks live in agent bodies by marketplace convention. Expected: PASS.

- [ ] **Step 3: Re-check the skill against the scaffold one last time**

```bash
cd /private/tmp/claude-501/-Users-bill-Projects-mynet/8934198f-2ad7-4e6f-823c-41527251b7fd/scratchpad/wailsprobe/wailsprobe && task build 2>&1 | tail -3 && ls -la bin/
```

Expected: build succeeds; binary present. If the scratchpad is gone, recreate with the `wails3 init` command from Verified Facts, then build.

- [ ] **Step 4: Report**

Final report lists: every check with its exact final output line, the Jev before and after tables from Task 5, the RED baseline from Task 1 beside the final audit, and any claim in the references left marked unverified (signing and notarization behavior). Do not merge or push; leave the branch for the owner.
