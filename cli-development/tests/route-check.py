#!/usr/bin/env python3
"""Routing simulation for cli-development agents.

Sends the three agent descriptions plus sample requests to the TypeSafe Jev API
as one Choice question per request and checks the winner and confidence.

Usage: TYPESAFE_API_KEY=... cli-development/tests/route-check.py
Exit 0 when every must-pass probe routes to its expected agent with
confidence >= THRESHOLD, 1 otherwise, 2 when the key is missing.
"""
import json, os, re, sys, time, urllib.error, urllib.request

API = "https://api.typesafe.ai/v1/systemone"
HERE = os.path.dirname(os.path.abspath(__file__))
AGENTS_DIR = os.path.join(HERE, "..", "agents")
NAMES = ["cli-developer", "cli-ui-designer", "go-tui-developer"]
THRESHOLD = 0.70

# (request, expected winner or None for informational)
PROBES = [
    ("Build a CLI for managing database migrations in Go", "cli-developer"),
    ("Our CLI's --help output is confusing, clean it up", "cli-developer"),
    ("Add a config file with env var overrides to our Python CLI", "cli-developer"),
    ("Build a TUI for browsing API responses", "go-tui-developer"),
    ("Add an interactive selection prompt to our Go CLI", "go-tui-developer"),
    ("Set up the CLI structure with Cobra and make the output look good", None),
    ("Wire our Go tool's plain export command and its interactive status screen under one Cobra root", "go-tui-developer"),
    ("Add theme support with user-defined color schemes and dark mode detection", "go-tui-developer"),
    ("The output of our CLI is hard to scan, everything looks the same", "cli-ui-designer"),
    ("Make this web dashboard feel like a terminal", "cli-ui-designer"),
    ("Add a branded ASCII header and colored status indicators to our CLI", "cli-ui-designer"),
    ("Write a Rust CLI with clap that streams JSON lines", None),
    ("Build an interactive TUI in Python with Textual", "none"),
    ("Add shell completions for zsh and fish to our Node CLI", "cli-developer"),
    ("Pick colors for our Bubble Tea app that work on light and dark terminals", "cli-ui-designer"),
    ("Our Go CLI prints tables with Lip Gloss but they break when the terminal is narrow", "go-tui-developer"),
    ("Design the error message format for our CLI", "cli-developer"),
    ("Write a bash script that wraps our deploy tool", None),
    ("Add a progress bar and spinner while the CLI uploads files (Node.js)", None),
]


def description(path):
    text = open(path).read()
    m = re.search(r"^description:\s*>\n((?:  .*\n?)+)", text, re.M)
    if not m:
        sys.exit(f"no folded description in {path}")
    return " ".join(line.strip() for line in m.group(1).splitlines())


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
            if e.code in (429, 502, 503, 529):
                time.sleep(2 ** attempt)
                continue
            sys.exit(f"HTTP {e.code}: {e.read().decode()[:500]}")
    sys.exit("retries exhausted")


def main():
    agents = {n: description(os.path.join(AGENTS_DIR, f"{n}.md")) for n in NAMES}
    state = {"agents": agents, "requests": {f"R{i}": p[0] for i, p in enumerate(PROBES)}}
    questions = {
        f"route_R{i}": {
            "type": "choice",
            "instructions": {
                "request_id": f"R{i}",
                "question": "Which agent in `agents` should handle `requests.<request_id>`, judging only from the agent descriptions?",
            },
            "criteria": {**{n: None for n in NAMES},
                         "none": "No listed agent fits; a different plugin or general assistant should handle it"},
        }
        for i in range(len(PROBES))
    }
    answers = call(state, questions)["answers"]
    failures = 0
    print(f"{'result':6} {'winner':17} {'conf':5} {'expected':17} request")
    for i, (request, expected) in enumerate(PROBES):
        a = answers[f"route_R{i}"]
        winner, conf = a["choice"], a["confidence"]
        if expected is None:
            status = "info"
        elif winner == expected and conf >= THRESHOLD:
            status = "ok"
        else:
            status = "FAIL"
            failures += 1
        print(f"{status:6} {winner:17} {conf:.2f}  {expected or '-':17} {request}")
    print(f"\n{failures} failing must-pass probes (threshold {THRESHOLD})")
    sys.exit(1 if failures else 0)


if __name__ == "__main__":
    main()
