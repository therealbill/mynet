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
