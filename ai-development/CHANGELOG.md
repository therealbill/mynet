# Changelog

All notable changes to ai-development will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2026-09-28

### Added

- `skills/agent-modernizer/scripts/audit-agents.py`: deterministic checks in code (required fields, name format, model alias, color, example structure and placement, body and description lengths, bullet counts) with severity and action composed by rules
- `skills/agent-modernizer/references/question-catalog.json`: fourteen typed judgments (Nouls and Scores) with explicit criteria, sent to TypeSafe Jev when `TYPESAFE_API_KEY` is set and answerable manually otherwise
- `--route DIR` example self-routing check: each agent's `user:` lines are routed over the sibling descriptions as one Choice per request; proactive examples are informational
- `--compare OLD NEW` rewrite regression check for dropped domain knowledge, changed role, and unsupported new claims
- Batch summary table with per-severity counts and a computed action per agent
- `Possible:` findings for judgments in the 0.4-0.7 band that must be confirmed by reading

### Changed

- Example blocks are expected in the body, after the frontmatter and before the role statement; examples inside `description` are now a Should-fix finding instead of their absence there being a Must-fix
- Description limit of 120 words added, since every description is loaded for routing on every turn
- Severity names unified to Must fix, Should fix, Consider across the skill, reference, and docs
- The single "redundancy" anti-pattern split into duplicate coverage and mirrored "Do Not" lists; the combined question sat near 0.5 for most agents
- Phantom-reference check uses the whole marketplace's agent roster, so cross-plugin references to real agents no longer count as phantom
- Example count range widened to 2-5 to match plugins that carry a proactive example
- `SKILL.md` process rewritten around run, confirm by reading, report, rewrite, validate; added a Do Not section
- `references/audit-criteria.md` gains the deterministic measurement table, the judgment table with flag rules, and a manual answering guide

## [1.0.0] - 2026-02-08

### Added

- `ai-engineer` agent
- `agent-modernizer` skill with audit criteria reference
- Diátaxis documentation set
