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
