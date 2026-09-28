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
