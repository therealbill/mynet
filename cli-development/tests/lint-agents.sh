#!/usr/bin/env bash
# Structural checks for cli-development agent files.
# Usage: cli-development/tests/lint-agents.sh [agent.md ...]
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
  examples="$(grep -c '^<example>' "$f")"

  grep -q 'Use when' <<<"$desc" \
    && ok "$name: description has 'Use when'" \
    || bad "$name: description lacks 'Use when'"

  grep -q 'Do not use' <<<"$desc" \
    && ok "$name: description has a 'Do not use' handoff" \
    || bad "$name: description lacks a 'Do not use' handoff"

  grep -qE '\(use [a-z-]+\)' <<<"$desc" \
    && ok "$name: handoff names a sibling agent" \
    || bad "$name: handoff does not name a sibling agent"

  grep -q '^\*\*Output:\*\*' <<<"$body" \
    && ok "$name: body has an Output section" \
    || bad "$name: body lacks an Output section"

  grep -q 'Proactive trigger:' "$f" \
    && ok "$name: has a proactive example" \
    || bad "$name: has no proactive example"

  [ "$examples" -ge 3 ] && [ "$examples" -le 5 ] \
    && ok "$name: has $examples examples" \
    || bad "$name: has $examples examples (want 3 to 5)"

  [ "$(printf '%s' "$body" | wc -c)" -lt 10000 ] \
    && ok "$name: body under 10000 chars" \
    || bad "$name: body is 10000 chars or more"

  grep -qE '\b(I am|I will|I would)\b' <<<"$body" \
    && bad "$name: body uses first person" \
    || ok "$name: body is second person"
done
exit $fail
