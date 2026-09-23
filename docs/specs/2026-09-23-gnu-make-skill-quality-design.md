# gnu-make Skill Quality Design

## Overview

**Plugin:** `gnu-make`
**Version:** 1.0.0 → 1.1.0
**Scope:** Fix verified technical defects in the five existing skills, replace their descriptions, cut redundancy and user-override language, add a sixth skill (`makefile-analysis`) that owns review checklists, import test scenarios, and update the plugin docs.

This design follows a review of the same skills in their two lineages (this plugin and the `make` plugin in ck-intelligence-agency). The review used TypeSafe's Jev model for rubric judgments over each SKILL.md, a trigger-routing experiment over the descriptions, and manual reproduction of every technical claim against GNU Make 3.81. Full report: https://claude.ai/artifact/LaXy5KH8PcVZMzisTQvUJq

### Findings the design responds to

- Seven technical defects reproduce as written. Three produce wrong behavior: the help target prints filenames once a Makefile uses `include`; the debugging grep `'File.*is newer'` never matches `make -d` output; the `-include dev.mk` / `-include prod.mk` "alternative" includes both files.
- Every body is two to three times the writing-skills target, and Jev rated redundancy at the top level in all five with confidence 0.84 to 1.00.
- Four skills instruct the agent to argue past a user's stated preference (Jev probability 0.81 to 0.95).
- `name:` fields use display names with spaces, which the agent-skills specification disallows.
- The fork's descriptions route the "missing separator" and "repetitive rules" requests to the wrong skill; a validated replacement set routes 14 of 14 sample requests correctly.
- No test scenarios exist in this plugin; the ck plugin has one per skill.
- There is no analysis/review skill, so each specialized skill carries its own review checklist.

### Decisions made during design

| Decision | Choice |
| --- | --- |
| Automated Jev harness in the repo | Not included. Testing uses subagent scenarios plus a snippet parse check. |
| Review ownership | New `makefile-analysis` skill owns all review checklists. |
| Test scenarios | Imported from the ck plugin into `gnu-make/tests/`. |
| Firmness | Firm on correctness (tabs, `.PHONY`, `.DELETE_ON_ERROR`, `$(MAKE)`); one trade-off explanation then defer on style (pattern rules, modular includes, phony vs loop). |
| Approach | Edit the fork in place; port analysis from the ck plugin. |
| Spec and plan location | `docs/specs/` and `docs/plans/` at the repo root. |

## Non-goals

- The ck-intelligence-agency `plugins/make` copy is not edited. Once this ships it should be retired from that marketplace; that is a separate task in that repo.
- No new Make topics beyond the specific additions listed per skill.
- No change to the plugin's Diátaxis docs structure, only content updates where the skills changed.

## Structure after the change

```
gnu-make/
  .claude-plugin/plugin.json          # version 1.1.0; description names review/analysis
  CHANGELOG.md                        # new
  skills/
    makefile-fundamentals/SKILL.md
    makefile-advanced-features/SKILL.md
    makefile-advanced-features/references/pattern-rule-scenarios.md
    makefile-includes-modularity/SKILL.md
    makefile-includes-modularity/references/module-patterns.md
    makefile-includes-modularity/references/shared-configuration.md
    makefile-recursive-multi-directory/SKILL.md
    makefile-debugging-optimization/SKILL.md
    makefile-debugging-optimization/references/debugging-scenarios.md
    makefile-analysis/SKILL.md        # new
  tests/
    README.md                         # how to run a RED/GREEN pass and the snippet check
    check-snippets.sh                 # extracts fenced makefile blocks, runs make -n -f on each
    scenarios/makefile-fundamentals.md
    scenarios/makefile-advanced-features.md
    scenarios/makefile-includes-modularity.md
    scenarios/makefile-recursive-multi-directory.md
    scenarios/makefile-debugging-optimization.md
    scenarios/makefile-analysis.md
    makefiles/test-good.mk
    makefiles/test-large.mk
    makefiles/test-multidir.mk
    makefiles/test-simple.mk
  docs/                               # seven pages updated for six skills
```

## Rules applied to every skill

### Frontmatter

- `name:` is the hyphenated directory slug (`makefile-fundamentals`, not `Makefile Fundamentals`).
- `description:` is replaced with the text below. These were validated by a routing experiment: 14 of 14 sample requests reached the intended skill, with the three weakest prompts moving to full confidence.
- `version:` is `1.1.0`.

| Skill | Description |
| --- | --- |
| makefile-fundamentals | Use when creating a Makefile, fixing a broken one, adding build targets, or hitting a 'missing separator' error; also when a Makefile lacks a help target, .PHONY declarations, or standard variables. |
| makefile-advanced-features | Use when a Makefile has repetitive explicit rules for similar files (main.o: main.c, utils.o: utils.c, ...), when compiling many sources of the same type, or when asked what $@, $<, $^, $* or % mean. |
| makefile-includes-modularity | Use when a single Makefile is large (roughly 150+ lines), mixes several components, or carries dev/staging/prod configuration in ifeq blocks; also when build configuration is shared across projects. |
| makefile-recursive-multi-directory | Use when a project has subdirectories with their own Makefiles, when a root Makefile loops over SUBDIRS with a shell for-loop, or when make -j gives no speedup on a multi-directory build. |
| makefile-debugging-optimization | Use when a build rebuilds everything, fails to rebuild after a change, is slow, or behaves unexpectedly and the cause is not yet known. |
| makefile-analysis | Use when asked to analyze, review, audit, assess, or check an existing Makefile for quality or best practices. |

### Body

- Target 500 to 700 words. Content beyond that moves to the skill's `references/` directory.
- Keep: Overview, When to Use, one complete example that runs as written, Quick Reference, Common Mistakes.
- Remove: Core Anti-Pattern narrative, Proactive Guidance, Red Flags, Real-World Impact, The Bottom Line, and every "When user insists/prefers" script.
- Remove emoji markers.
- Firmness: tabs, `.PHONY`, `.DELETE_ON_ERROR`, and `$(MAKE)` are stated once as correctness requirements. Pattern rules, modular includes, and phony-vs-loop get one trade-off paragraph that ends with following the user's stated choice. The phrases "even if user requests it", "don't accept", and "user preference not challenged" do not appear.
- Every fenced `makefile` block passes `make -n -f` (see Testing). Blocks that are intentionally partial are fenced as `text`, not `makefile`.

## Per-skill changes

### makefile-fundamentals

- Help target recipe becomes:
  ```
  @grep -hE '^[a-zA-Z0-9_./-]+:.*## ' $(MAKEFILE_LIST) | \
      awk 'BEGIN {FS = ":.*## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
  ```
  in both places it appears. One sentence explains that `-h` suppresses filename prefixes once `include` adds files to `MAKEFILE_LIST`, and that the character class admits targets such as `docs/site`.
- `LDFLAGS = -lm` becomes `LDLIBS = -lm`; link recipes use `$(CC) $(LDFLAGS) $^ $(LDLIBS) -o $@`.
- `.DELETE_ON_ERROR` is added to the numbered Essential Best Practices list.
- The Automatic Variables and Pattern Rules sections shrink to a short pointer at `makefile-advanced-features`.
- Proactive Suggestions and Red Flags are removed. "Handling Keep It Simple Requests" stays as the correctness statement.

### makefile-advanced-features

- The `$(MAKE)` row leaves Common Mistakes.
- The complete example gains a build directory with an order-only prerequisite (`$(OBJS): | $(BUILDDIR)`) and the `-MMD -MP` dependency idiom (`CFLAGS += -MMD -MP`, `-include $(OBJS:.o=.d)`). The older `%.d: %.c` rule is removed from `references/pattern-rule-scenarios.md`.
- "When User Prefers Explicit Rules", "Keep It Simple Requests", Proactive Guidance, Red Flags, and The Bottom Line are replaced by one trade-off paragraph. Quick Reference Card stays.

### makefile-includes-modularity

- "Alternative: Conditional with `-include`" is deleted.
- The Correct Pattern example becomes one buildable three-file set (Makefile, config.mk, rules.mk) whose targets all have recipes. The `%: %.o` and `%.a: $(OBJECTS)` rules are removed everywhere.
- Add `$(dir $(lastword $(MAKEFILE_LIST)))` for locating sibling includes and a note that `include` search paths come from `-I`.
- Benefits Summary, The Bottom Line, Red Flags, and Proactive Guidance are removed. The size-threshold table stays.

### makefile-recursive-multi-directory

- `clean` uses the phony pattern: `CLEANDIRS = $(SUBDIRS:%=clean-%)`, `.PHONY: clean $(CLEANDIRS)`, `clean: $(CLEANDIRS)`, `$(CLEANDIRS): ; $(MAKE) -C $(@:clean-%=%) clean`. The bare loop is removed.
- One paragraph names non-recursive make as the alternative for tightly coupled trees, with a pointer.
- Real-World Performance Impact becomes a one-line worked example labeled as such; "8x" does not appear as a fact.
- "When User Insists on Loops", Red Flags, and The Bottom Line are removed.

### makefile-debugging-optimization

- Every `grep 'File.*is newer'` becomes `grep 'is newer than target'` in SKILL.md and `references/debugging-scenarios.md`.
- "`make -d` to see where variable is set" becomes `make -p`, citing the `# makefile (from 'Makefile', line N)` comment.
- `-j$(shell nproc)` becomes `-j$(shell nproc 2>/dev/null || sysctl -n hw.ncpu)`.
- The "Lazy Wildcard" heading becomes "Immediate expansion for wildcard".
- A `make --version` step precedes every 4.x-only flag, with a note that macOS ships GNU Make 3.81. The flags table adds `$(info)`, `--warn-undefined-variables`, `-O`/`--output-sync` (4.0), and `--shuffle` (4.4).
- The scripted dialogue blocks are removed. The four-step approach and the flags table stay.

### makefile-analysis (new)

- Ported from the ck-intelligence-agency plugin and cut to the word target.
- Contains: the five-dimension checklist (absorbing the Red Flags removed from the other skills), severity rules, and the report template.
- "No help target" is Recommended, not Critical. The "247-line Makefile" example and the "9 focused files" figure are removed.
- Cross-references the five specialized skills by name for fixes.

## Testing

Every skill edit follows the writing-skills Iron Law: baseline, edit, re-run.

1. **Scenario import.** The six `test-scenarios.md` files from `ck-intelligence-agency/plugins/make/skills/*/` move to `gnu-make/tests/scenarios/<skill>.md`; the four sample Makefiles move to `gnu-make/tests/makefiles/`. Result narratives are not imported.
2. **Scenario rewrite.** Before any skill edit, the scenarios that encode the override stance are rewritten: "user insists on loops" and "user prefers explicit rules" expect one trade-off explanation followed by compliance with the user's choice.
3. **RED.** For each skill, a subagent runs its scenarios against the current SKILL.md; the runner records which expectations pass.
4. **GREEN.** After the edit, the same scenarios run again; every expectation that passed still passes, and the targeted expectations now pass.
5. **Snippet check.** `tests/check-snippets.sh <SKILL.md>` extracts every fenced `makefile` block into one temp directory per SKILL.md. A block whose first line is a comment naming a file (`# Makefile`, `# config.mk`) is saved under that name so multi-file examples in consecutive blocks resolve their `include` lines; other blocks get sequential names. The script appends a dummy target `__snippet_check__: ;` to each file and runs `make -n -f <file> __snippet_check__`, which parses the whole file without needing a default goal. A non-zero exit fails the task and the script prints the block's line number in the SKILL.md. Snippets that are deliberately partial are fenced as `text`, not `makefile`.
6. **`tests/README.md`** documents both loops in under 200 words.

## Docs

Seven pages under `gnu-make/docs/` name the skills and update after skill content settles:

- `reference/skills.md`: sixth entry; trigger descriptions replaced with the new text.
- `explanation/skill-progression.md`: analysis added as the capstone; "non-negotiable" language softened to match the firmness decision.
- `explanation/architecture.md` and `_index.md`: six skills.
- `howto/debug-slow-makefile.md`, `howto/organize-multi-directory-build.md`, `howto/split-large-makefile.md`: the `make -d` grep and help-target snippets corrected so the docs match the skills.

## Versioning

- `gnu-make/.claude-plugin/plugin.json` and the `gnu-make` entry in `.claude-plugin/marketplace.json` go to 1.1.0; the description mentions analysis/review.
- New `gnu-make/CHANGELOG.md` with one 1.1.0 entry: the seven fixes, the new skill, the description rewrite, the firmness change, and the tests directory.
- Every SKILL.md `version:` is 1.1.0.

## Agent team

Four concurrent agents, then one sequential close-out. Each agent owns its files and does not edit another agent's files.

| Agent | Work | Model |
| --- | --- | --- |
| analysis-and-tests | Writes `tests/check-snippets.sh` first, then imports and rewrites scenarios, then ports and trims `makefile-analysis` | sonnet |
| fundamentals-and-advanced | `makefile-fundamentals`, `makefile-advanced-features`, and the overlap trim between them | opus |
| includes-and-recursive | `makefile-includes-modularity`, `makefile-recursive-multi-directory` | opus |
| debugging | `makefile-debugging-optimization` and its references file | sonnet |
| close-out (after all four) | Docs pages, versioning, CHANGELOG, final scenario pass across all six skills | sonnet |

The three skill-editing agents wait for `check-snippets.sh` before their GREEN step. The plan does not include timeline estimates.

## Acceptance

- All six SKILL.md files have hyphenated `name:` values and the descriptions above.
- No SKILL.md body exceeds 700 words; none contains the removed section headings or the override phrases.
- Every fenced `makefile` block in every SKILL.md passes `make -n -f`.
- Each skill's scenarios pass GREEN with no regression from RED.
- The help target snippet, run under a Makefile with one `include`, prints target names, not filenames.
- `make -d` on a stale target, piped through the skill's grep, prints the prerequisite line.
- Docs name six skills and contain no `File.*is newer` grep.
- `plugin.json`, `marketplace.json`, `CHANGELOG.md`, and all `version:` fields read 1.1.0.
