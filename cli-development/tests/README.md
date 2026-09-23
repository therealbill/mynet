# cli-development tests

Two checks guard the agent files.

## lint-agents.sh

Structural rules every agent must satisfy: a "Use when" description with a
"Do not use for ... (use <sibling>)" handoff, an `**Output:**` section, three
to five example blocks including one whose commentary starts with
"Proactive trigger:", a body under 10,000 characters, and second-person voice.

    cli-development/tests/lint-agents.sh                # all agents
    cli-development/tests/lint-agents.sh cli-development/agents/cli-developer.md

Exit 1 on any `FAIL` line.

## route-check.py

Sends the three descriptions and nineteen sample requests to the TypeSafe Jev
API as one Choice question per request. Must-pass probes need the expected
winner at confidence 0.70 or higher. Probes marked `info` are reported only.

    TYPESAFE_API_KEY=... cli-development/tests/route-check.py

Exit 0 when all must-pass probes pass, 1 otherwise, 2 when the key is unset.
Descriptions are the only input, so changing a body never changes this result.
Do not tune a description to a single probe; fix the scope statement instead.
