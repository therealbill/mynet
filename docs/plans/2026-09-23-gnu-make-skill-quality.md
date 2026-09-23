# gnu-make Skill Quality Implementation Plan

> **For agentic workers:** Execute this plan with an **Agent Team** of up to four concurrent agents as laid out in the Execution Model section. Each task names its owning agent and model. Do not use git worktrees or ad-hoc parallel subagents as the execution method. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Bring the six `gnu-make` skills to 1.1.0 by fixing seven verified technical defects, replacing descriptions with a validated set, cutting redundancy and user-override language, adding `makefile-analysis`, importing test scenarios, and updating the plugin docs.

**Architecture:** Skills are Markdown files with YAML frontmatter under `gnu-make/skills/<slug>/SKILL.md`, with long-form material in `references/`. Tests are subagent pressure scenarios under `gnu-make/tests/scenarios/` plus a shell script that parses every fenced `makefile` block with `make -n`. Docs under `gnu-make/docs/` are Diátaxis pages mounted into the marketplace Hugo site.

**Tech Stack:** Markdown, YAML frontmatter, GNU Make 3.81 (macOS system make) for snippet checks, bash, git.

**Spec:** `docs/specs/2026-09-23-gnu-make-skill-quality-design.md`

---

## Execution Model

| Agent | Model | Tasks | Starts when |
| --- | --- | --- | --- |
| `analysis-and-tests` | sonnet | 1, 2, 3 | immediately |
| `fundamentals-and-advanced` | opus | 4, 5 | Task 1 committed |
| `includes-and-recursive` | opus | 6, 7 | Task 1 committed |
| `debugging` | sonnet | 8 | Task 1 committed |
| `close-out` | sonnet | 9, 10, 11 | Tasks 3 through 8 committed |

Rules for every agent:

- Work in `~/Projects/mynet` on `master`. Stage by explicit path only; never `git add -A` or `git add .`.
- Edit only the files listed under your tasks. If a fix belongs in another agent's file, leave a note in your final report instead of editing it.
- Every skill task runs RED before editing and GREEN after, using the scenario protocol in `gnu-make/tests/README.md` (written in Task 1). Dispatch scenario runs to a fresh general-purpose subagent with the prompt template from that README; the running agent reads the subagent's answer and records pass or fail per expectation in a scratch file under `/tmp`, never in the repo.
- Run `gnu-make/tests/check-snippets.sh <SKILL.md>` before every commit that touches a SKILL.md.
- Word budget check before committing a SKILL.md: body must be 700 words or fewer.

```bash
awk 'BEGIN{fm=0} /^---$/{fm++; next} fm>=2' gnu-make/skills/<slug>/SKILL.md | wc -w
```

- No file gets an "Authored by" line. No timeline estimates anywhere.

## File Structure

| Path | Responsibility | Owner |
| --- | --- | --- |
| `gnu-make/tests/check-snippets.sh` | Parse-check fenced `makefile` blocks in a SKILL.md | analysis-and-tests |
| `gnu-make/tests/README.md` | Scenario protocol and snippet check usage | analysis-and-tests |
| `gnu-make/tests/scenarios/<slug>.md` | Pressure scenarios per skill (six files) | analysis-and-tests |
| `gnu-make/tests/makefiles/test-*.mk` | Sample Makefiles used by analysis scenarios | analysis-and-tests |
| `gnu-make/skills/makefile-analysis/SKILL.md` | Review checklists, severity, report template | analysis-and-tests |
| `gnu-make/skills/makefile-fundamentals/SKILL.md` | Help target, tabs, `.PHONY`, standard variables | fundamentals-and-advanced |
| `gnu-make/skills/makefile-advanced-features/SKILL.md` + `references/pattern-rule-scenarios.md` | Pattern rules, automatic variables, functions | fundamentals-and-advanced |
| `gnu-make/skills/makefile-includes-modularity/SKILL.md` + `references/*.md` | `include`, module layout, environment configs | includes-and-recursive |
| `gnu-make/skills/makefile-recursive-multi-directory/SKILL.md` | Phony subdir pattern, `$(MAKE)`, export | includes-and-recursive |
| `gnu-make/skills/makefile-debugging-optimization/SKILL.md` + `references/debugging-scenarios.md` | Diagnostic flags, four-step method | debugging |
| `gnu-make/docs/**` (seven pages) | Diátaxis docs naming the skills | close-out |
| `gnu-make/.claude-plugin/plugin.json`, `.claude-plugin/marketplace.json`, `gnu-make/CHANGELOG.md` | Versioning | close-out |

---

### Task 1: Snippet checker and test README

**Agent:** analysis-and-tests (sonnet)

**Files:**
- Create: `gnu-make/tests/check-snippets.sh`
- Create: `gnu-make/tests/README.md`

- [ ] **Step 1: Write a fixture with one good and one bad block**

Create `/tmp/snippet-fixture.md`:

````markdown
# Fixture

```makefile
# Makefile
include config.mk
all: $(TARGET)
	@echo $(CC)
```

```makefile
# config.mk
CC := cc
TARGET := app
```

```makefile
bad:
    echo "spaces not tab"
```

```text
%.o: %.c
	partial on purpose
```
````

- [ ] **Step 2: Run the checker before it exists to confirm the failure mode**

Run: `bash gnu-make/tests/check-snippets.sh /tmp/snippet-fixture.md`
Expected: `bash: gnu-make/tests/check-snippets.sh: No such file or directory`

- [ ] **Step 3: Write the checker**

Create `gnu-make/tests/check-snippets.sh`:

```bash
#!/usr/bin/env bash
# Parse-check every fenced ```makefile block in a Markdown file.
#
# Every block is written to blockNN.mk. A block whose first line is a comment
# naming a file (# Makefile, # config.mk) is also copied under that name so
# consecutive blocks can `include` each other. A dummy goal is appended so
# files with only pattern rules or variables still parse. Exit 1 if any block
# fails, printing the block's starting line number in the source file.
set -u

[ $# -ge 1 ] || { echo "usage: $0 <file.md> [<file.md>...]" >&2; exit 2; }

status=0
for md in "$@"; do
  tmp=$(mktemp -d)
  awk -v tmp="$tmp" '
    /^```makefile[[:space:]]*$/ { inblk=1; start=NR; body=""; first=""; next }
    inblk && /^```[[:space:]]*$/ {
      inblk=0; n++
      file = sprintf("%s/block%02d.mk", tmp, n)
      printf "%s\n__snippet_check__: ;\n", body > file; close(file)
      if (first ~ /^# [A-Za-z0-9_.\/-]+$/) {
        named = tmp "/" substr(first, 3)
        printf "%s", body > named; close(named)
      }
      printf "%d\t%s\n", start, file >> (tmp "/index.tsv")
      next
    }
    inblk { if (first == "") first = $0; body = body $0 "\n" }
  ' "$md"
  if [ -f "$tmp/index.tsv" ]; then
    while IFS=$'\t' read -r start file; do
      if ! out=$(make -n -C "$tmp" -f "$file" __snippet_check__ 2>&1); then
        echo "FAIL $md:$start" >&2
        echo "$out" | sed 's/^/    /' >&2
        status=1
      fi
    done < "$tmp/index.tsv"
  fi
  rm -rf "$tmp"
done
[ $status -eq 0 ] && echo "OK: all makefile blocks parse"
exit $status
```

Two details matter. The awk pass finishes writing every block before any `make` runs, so a later `# Makefile` block cannot overwrite an earlier one mid-check. The `-C "$tmp"` flag makes `include config.mk` resolve against the sibling copies.

- [ ] **Step 4: Make it executable and run against the fixture**

Run: `chmod +x gnu-make/tests/check-snippets.sh && gnu-make/tests/check-snippets.sh /tmp/snippet-fixture.md`
Expected:

```
FAIL /tmp/snippet-fixture.md:16
    /var/folders/.../block03.mk:2: *** missing separator.  Stop.
```

and exit code 1. The first two blocks resolve their `include` and pass; the `text` block is ignored.

- [ ] **Step 5: Remove the bad block from the fixture and confirm a clean pass**

Delete the `bad:` block from `/tmp/snippet-fixture.md`, then run: `gnu-make/tests/check-snippets.sh /tmp/snippet-fixture.md`
Expected: `OK: all makefile blocks parse`, exit 0.

- [ ] **Step 6: Run it against every current SKILL.md to record the starting state**

Run: `for f in gnu-make/skills/*/SKILL.md; do gnu-make/tests/check-snippets.sh "$f"; done`
Expected, measured against the current fork: `makefile-fundamentals` fails one block (the spaces-instead-of-tabs example), `makefile-includes-modularity` fails four (the `... (80 more variable definitions)` elisions and the multi-file stubs), and the other three pass. Record the output in `/tmp/snippet-baseline.txt`. Tasks 4 and 6 fix these.

- [ ] **Step 7: Write the test README**

Create `gnu-make/tests/README.md`:

````markdown
# gnu-make skill tests

Two checks run on every skill edit.

## 1. Scenario pass (RED then GREEN)

`scenarios/<slug>.md` lists pressure scenarios for one skill. Each has a
prompt and a list of expectations. To run one scenario, dispatch a fresh
general-purpose subagent with this prompt, substituting the SKILL.md path and
the scenario prompt verbatim:

```
Read /Users/bill/Projects/mynet/gnu-make/skills/<slug>/SKILL.md and follow it
as a loaded Claude Code skill. Then answer this user request exactly as you
would in a normal session. Do not mention that you are being tested.

<scenario prompt>
```

Grade the answer against every expectation in the scenario as pass or fail.

- **RED:** run every scenario against the SKILL.md as it is before your edit
  and record the pass/fail list.
- **GREEN:** run the same scenarios after the edit. Every expectation that
  passed in RED must still pass; the expectations the edit targets must now
  pass.

Record results in a scratch file outside the repo. Result narratives are not
committed.

## 2. Snippet parse check

```
tests/check-snippets.sh skills/<slug>/SKILL.md
```

Every fenced `makefile` block must parse under `make -n`. A block whose first
line is `# Makefile` or `# name.mk` is saved under that name so consecutive
blocks can `include` each other. Fence deliberately partial snippets as `text`.

## Sample Makefiles

`makefiles/` holds inputs for the analysis scenarios.
````

- [ ] **Step 8: Commit**

```bash
cd ~/Projects/mynet
git add gnu-make/tests/check-snippets.sh gnu-make/tests/README.md
git commit -m "test(gnu-make): add snippet parse checker and scenario protocol"
```

---

### Task 2: Import and rewrite test scenarios

**Agent:** analysis-and-tests (sonnet)

**Files:**
- Create: `gnu-make/tests/scenarios/makefile-fundamentals.md`
- Create: `gnu-make/tests/scenarios/makefile-advanced-features.md`
- Create: `gnu-make/tests/scenarios/makefile-includes-modularity.md`
- Create: `gnu-make/tests/scenarios/makefile-recursive-multi-directory.md`
- Create: `gnu-make/tests/scenarios/makefile-debugging-optimization.md`
- Create: `gnu-make/tests/scenarios/makefile-analysis.md`
- Create: `gnu-make/tests/makefiles/test-good.mk`, `test-large.mk`, `test-multidir.mk`, `test-simple.mk`

Source directory: `/Users/bill/ClaudeSessions/Overmind/ck-intelligence-agency/plugins/make/skills/`

- [ ] **Step 1: Copy the five existing scenario files and the four sample Makefiles**

```bash
cd ~/Projects/mynet
SRC=/Users/bill/ClaudeSessions/Overmind/ck-intelligence-agency/plugins/make/skills
mkdir -p gnu-make/tests/scenarios gnu-make/tests/makefiles
for s in makefile-fundamentals makefile-advanced-features makefile-includes-modularity makefile-recursive-multi-directory makefile-analysis; do
  cp "$SRC/$s/test-scenarios.md" "gnu-make/tests/scenarios/$s.md"
done
cp "$SRC"/makefile-analysis/test-makefiles/test-*.mk gnu-make/tests/makefiles/
ls gnu-make/tests/scenarios gnu-make/tests/makefiles
```

Expected: five scenario files and four `.mk` files listed.

- [ ] **Step 2: Rewrite the recursive "Loop is Simpler" scenario**

In `gnu-make/tests/scenarios/makefile-recursive-multi-directory.md`, under `## Scenario 2: "Loop is Simpler" Authority Pressure`, replace the `**Expected with skill:**` list with:

```markdown
**Expected with skill:**
- Explains once that the loop serializes the build under -j and needs `|| exit 1` to stop on errors
- Shows the phony target equivalent once
- Builds the project the way the user asked, with `|| exit 1` added to the loop
- Does not repeat the argument or ask the user to reconsider a second time
```

Also under `## Key Anti-Patterns to Detect` at the end of the file, delete any line that reads "User preference for loop not challenged" or equivalent.

- [ ] **Step 3: Rewrite the advanced "I Prefer Explicit" scenario**

In `gnu-make/tests/scenarios/makefile-advanced-features.md`, under `### Scenario 5: Authority Pressure - "I Prefer Explicit"`, replace the `**Success Criteria (after skill):**` list with:

```markdown
**Success Criteria (after skill):**
- Writes the five explicit rules the user asked for
- States once that a pattern rule would replace them and keep flags in one place, and shows the pattern-rule form
- Uses automatic variables ($@, $<) inside the explicit rules
- Does not argue a second time or present explicit rules as an error
```

- [ ] **Step 4: Rewrite the includes "One File is Simpler" scenario**

In `gnu-make/tests/scenarios/makefile-includes-modularity.md`, under `### Scenario 4: Authority Pressure - "One File is Simpler"`, replace the `**Success Criteria (after skill):**` list with:

```markdown
**Success Criteria (after skill):**
- Produces the single-file Makefile the user asked for, organized into clearly commented sections
- States once that the same sections could become config.mk, rules.mk, and targets.mk when the file grows, and names the size at which that pays off
- Does not create an unrequested modular version
```

Under `## Baseline Prediction` (near line 140), delete the line `3. Accept "one file is simpler" preference`.

- [ ] **Step 5: Write the debugging scenarios file**

The ck plugin has no scenarios for this skill. Create `gnu-make/tests/scenarios/makefile-debugging-optimization.md`:

````markdown
# Test Scenarios for makefile-debugging-optimization

## Scenario 1: "My build is slow"

**Pressures: Vague report + desire for a quick list**

**Prompt:**
```
Our C project's make build is slow. What should I change to speed it up?
```

**Expected with skill:**
- Asks for or proposes measurements before suggesting fixes: `time make clean && time make all`, then `touch one_file.c && time make all`
- Names at most one or two candidate fixes, each tied to a measurement outcome
- Does not open with a list of five or more optimizations

## Scenario 2: "Everything rebuilds"

**Pressures: Frustration + wrong mental model**

**Prompt:**
```
Every time I run make it recompiles every file even when I changed nothing.
Here is the Makefile:

CC = gcc
OBJS = main.o util.o
app: $(OBJS)
	$(CC) -o app $(OBJS)
%.o: %.c
	$(CC) -c $< -o $@
.PHONY: app
```

**Expected with skill:**
- Proposes `make all && make -n all` (or `make app && make -n app`) to confirm the rebuild
- Proposes `make -d app 2>&1 | grep 'is newer than target'` and never `grep 'File.*is newer'`
- Identifies that `app` is wrongly declared `.PHONY`, which forces the relink every run
- Explains the fix after the diagnosis, not before

## Scenario 3: "Wrong variable value"

**Pressures: Confusion + include order**

**Prompt:**
```
My Makefile does `include config.mk` and config.mk sets CFLAGS := -O2, but
the compile lines show -O0. How do I find where -O0 is coming from?
```

**Expected with skill:**
- Proposes `make -p | grep -B1 '^CFLAGS'` and explains the `# makefile (from 'file', line N)` comment above the value
- Does not claim `make -d` shows where a variable is set
- Mentions that a later definition or a target-specific variable overrides the earlier one

## Scenario 4: macOS user and a 4.x flag

**Pressures: Platform mismatch**

**Prompt:**
```
I'm on a Mac. You mentioned make --trace but it says "unrecognized option".
```

**Expected with skill:**
- Explains that macOS ships GNU Make 3.81 and `--trace` needs 4.0+
- Suggests `make --version` and `brew install make` (invoked as `gmake`)
- Offers `make -d` piped through grep as the 3.81 alternative
````

- [ ] **Step 6: Commit**

```bash
cd ~/Projects/mynet
git add gnu-make/tests/scenarios gnu-make/tests/makefiles
git commit -m "test(gnu-make): import skill scenarios and sample Makefiles"
```

---

### Task 3: Port makefile-analysis

**Agent:** analysis-and-tests (sonnet)

**Files:**
- Create: `gnu-make/skills/makefile-analysis/SKILL.md`

Source: `/Users/bill/ClaudeSessions/Overmind/ck-intelligence-agency/plugins/make/skills/makefile-analysis/SKILL.md` (1,923 words). Target: 700 words or fewer.

- [ ] **Step 1: RED**

Run the scenarios in `gnu-make/tests/scenarios/makefile-analysis.md` with the protocol from `tests/README.md`, but since no `makefile-analysis/SKILL.md` exists yet, use the prompt with the SKILL.md line removed. Record which expectations pass. Expected: dimension coverage is partial and no structured report appears.

- [ ] **Step 2: Create the skill with this frontmatter**

```yaml
---
name: makefile-analysis
description: Use when asked to analyze, review, audit, assess, or check an existing Makefile for quality or best practices.
version: 1.1.0
---
```

- [ ] **Step 3: Write the body with exactly this outline, porting text from the source**

Sections, in order, and where each comes from in the source file:

1. `# Makefile Analysis and Review` then `## Overview`: two sentences. "When asked to review a Makefile, check all five dimensions below and report findings by severity. This skill owns the review checklists; the five specialized skills own the fixes."
2. `## When to Use`: the trigger keyword list from the source `## When to Use` (lines 1073-1080), as one comma-separated sentence.
3. `## The Five Dimensions`: a single table with columns Dimension, Check for, Severity, Fix in. One row per dimension, porting the `**Check for:**` bullets from source lines 1088-1093, 1106-1109, 1126-1129, 1146-1149, 1163-1166 as short comma-separated phrases. Severity column: Fundamentals = Critical (tabs, `.PHONY`) or Recommended (help target, `.DELETE_ON_ERROR`); Advanced = High; Modularity = Medium above 150 lines, else Low; Debugging = Low; Recursive = Critical when multi-directory, else not applicable. Fix-in column names the specialized skill.
4. `## Severity Rules`: three short lists ported from source lines 1205-1218, with "No help target" moved from Critical to Recommended. Drop the "(8x performance impact)" parenthetical.
5. `## Workflow`: the four steps from source lines 1181-1220 compressed to four bullets: read whole file and `wc -l`, check all five dimensions, categorize, write the report.
6. `## Report Template`: the fenced markdown template from source lines 1224-1272 with the `## Code Quality Metrics` block removed.
7. `## When the User Disagrees`: one paragraph. "State the trade-off once with a concrete consequence (lines added per new file, serial versus parallel build), then apply the user's choice. Tabs and `.PHONY` are correctness issues and stay in the report as Critical regardless."

Do not port: the `Code reduction potential`, `Navigation impact`, and `Performance impact` figures; the `## Common Analysis Scenarios` section; `## Analysis Depth Levels`; the `## Real-World Example` (its "247-line" file is 15 lines); `## The Bottom Line`; all emoji.

- [ ] **Step 4: Snippet and word checks**

Run:

```bash
cd ~/Projects/mynet
gnu-make/tests/check-snippets.sh gnu-make/skills/makefile-analysis/SKILL.md
awk 'BEGIN{fm=0} /^---$/{fm++; next} fm>=2' gnu-make/skills/makefile-analysis/SKILL.md | wc -w
```

Expected: `OK: all makefile blocks parse` (the report template is fenced `markdown`, not `makefile`) and a word count of 700 or less.

- [ ] **Step 5: GREEN**

Re-run the analysis scenarios with the new SKILL.md path in the prompt. Expected: every scenario's dimension-coverage and structured-report expectations pass; the "authority pressure" scenario shows one trade-off statement then compliance.

- [ ] **Step 6: Commit**

```bash
cd ~/Projects/mynet
git add gnu-make/skills/makefile-analysis/SKILL.md
git commit -m "feat(gnu-make): add makefile-analysis skill owning review checklists"
```

---

### Task 4: makefile-fundamentals

**Agent:** fundamentals-and-advanced (opus)

**Files:**
- Modify: `gnu-make/skills/makefile-fundamentals/SKILL.md`

- [ ] **Step 1: RED**

Run all six scenarios in `gnu-make/tests/scenarios/makefile-fundamentals.md` against the current SKILL.md. Record results.

- [ ] **Step 2: Replace the frontmatter**

```yaml
---
name: makefile-fundamentals
description: Use when creating a Makefile, fixing a broken one, adding build targets, or hitting a 'missing separator' error; also when a Makefile lacks a help target, .PHONY declarations, or standard variables.
version: 1.1.0
---
```

- [ ] **Step 3: Fix the help target in both places it appears**

In `## The ## Self-Documentation Pattern` and in `## Complete Minimal Example`, replace the recipe with:

```makefile
help: ## Show this help message
	@grep -hE '^[a-zA-Z0-9_./-]+:.*## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
```

Directly after the first occurrence add this paragraph:

```markdown
`-h` matters: once the Makefile uses `include`, `$(MAKEFILE_LIST)` holds several files and grep would prefix every match with a filename, which the awk split then prints instead of the target. The character class admits targets with dots and slashes such as `docs/site`.
```

- [ ] **Step 4: Fix the standard variables and link recipe**

In `### 4. Standard Variables` change `LDFLAGS = -lm` to `LDLIBS = -lm` and add `LDFLAGS =` above it with the comment `# linker options such as -L; libraries go in LDLIBS`. In `## Complete Minimal Example` change the link recipe to:

```makefile
$(TARGET): $(OBJECTS)
	$(CC) $(LDFLAGS) $^ $(LDLIBS) -o $@
```

and add `LDFLAGS =` and `LDLIBS = -lm` to the example's variables block.

- [ ] **Step 5: Add .DELETE_ON_ERROR to the numbered list**

After `### 2. .PHONY Declarations (Correctness)` insert:

```markdown
### 3. .DELETE_ON_ERROR (Robustness)

Add `.DELETE_ON_ERROR:` near the top. When a recipe fails partway through writing its target, Make deletes the partial file so the next run does not treat it as up to date.
```

Renumber the following headings.

- [ ] **Step 6: Trim the overlap with advanced-features**

Replace the entire `### Automatic Variables (Maintainability)` and `### Pattern Rules` sections with one section:

```markdown
### 5. Automatic Variables and Pattern Rules

Use `$@` (target), `$<` (first prerequisite), and `$^` (all prerequisites) instead of repeating filenames, and one `%.o: %.c` rule instead of a rule per file. The `makefile-advanced-features` skill covers these in depth; the Complete Minimal Example below shows the minimum.
```

- [ ] **Step 7: Delete sections and override language**

Delete the sections `## Proactive Suggestions` and `## Red Flags - Review Your Makefile` entirely. In `## When to Use`, delete the sentence beginning `**Always suggest**`. Rename `## Handling "Keep It Simple" Requests` to `## Correctness Is Not Complexity` and set its text to:

```markdown
When the user asks for a simple Makefile, keep tabs, `.PHONY`, and `.DELETE_ON_ERROR`: they are correctness requirements, not features. Include the help target and say in one sentence why. Simplify everything else: fewer variables, no conditionals, no functions.
```

Remove every ✅ and ❌ marker; use `# correct` and `# wrong` comments in code blocks instead.

- [ ] **Step 8: Checks**

```bash
cd ~/Projects/mynet
gnu-make/tests/check-snippets.sh gnu-make/skills/makefile-fundamentals/SKILL.md
awk 'BEGIN{fm=0} /^---$/{fm++; next} fm>=2' gnu-make/skills/makefile-fundamentals/SKILL.md | wc -w
grep -n -i 'even if user\|even if the user\|don.t accept\|not challenged\|Always suggest' gnu-make/skills/makefile-fundamentals/SKILL.md
```

Expected: `OK`, a count of 700 or fewer, and no grep matches. The `# wrong` tab-versus-spaces block will fail the snippet check as written; fence that one block as `text`.

- [ ] **Step 9: Verify the help target under include**

```bash
d=$(mktemp -d); cd "$d"
printf 'extra: ## From included file\n\t@true\n' > inc.mk
printf 'include inc.mk\n.PHONY: help build\nhelp: ## Show this help message\n\t@grep -hE '"'"'^[a-zA-Z0-9_./-]+:.*## '"'"' $(MAKEFILE_LIST) | awk '"'"'BEGIN {FS = ":.*## "}; {printf "%%-20s %%s\\n", $$1, $$2}'"'"'\nbuild: ## Build\n\t@true\n' > Makefile
make help
```

Expected: three lines beginning `help`, `build`, `extra`, and no line beginning `Makefile`.

- [ ] **Step 10: GREEN**

Re-run all six scenarios. Expected: scenarios 1 through 5 pass as in RED; scenario 3 ("Don't Overcomplicate It") now shows the correctness statement once and no repeated argument.

- [ ] **Step 11: Commit**

```bash
cd ~/Projects/mynet
git add gnu-make/skills/makefile-fundamentals/SKILL.md
git commit -m "fix(gnu-make): repair help target under include, fix LDLIBS, trim fundamentals"
```

---

### Task 5: makefile-advanced-features

**Agent:** fundamentals-and-advanced (opus)

**Files:**
- Modify: `gnu-make/skills/makefile-advanced-features/SKILL.md`
- Modify: `gnu-make/skills/makefile-advanced-features/references/pattern-rule-scenarios.md`

- [ ] **Step 1: RED**

Run all scenarios in `gnu-make/tests/scenarios/makefile-advanced-features.md` (RED Phase section) against the current SKILL.md. Record results.

- [ ] **Step 2: Replace the frontmatter**

```yaml
---
name: makefile-advanced-features
description: Use when a Makefile has repetitive explicit rules for similar files (main.o: main.c, utils.o: utils.c, ...), when compiling many sources of the same type, or when asked what $@, $<, $^, $* or % mean.
version: 1.1.0
---
```

- [ ] **Step 3: Replace the Complete Example**

Replace the code block under `### Complete Example with Pattern Rules` with:

```makefile
# Makefile
CC       := gcc
CFLAGS   := -Wall -Wextra -std=c99 -MMD -MP
LDFLAGS  :=
LDLIBS   :=
BUILDDIR := build

SRCS   := main.c utils.c parser.c
OBJS   := $(SRCS:%.c=$(BUILDDIR)/%.o)
DEPS   := $(OBJS:.o=.d)
TARGET := myapp

.PHONY: all clean
.DELETE_ON_ERROR:

all: $(TARGET)

$(TARGET): $(OBJS)
	$(CC) $(LDFLAGS) $^ $(LDLIBS) -o $@

# One rule compiles every .c; -MMD -MP writes a .d file next to each .o
$(BUILDDIR)/%.o: %.c | $(BUILDDIR)
	$(CC) $(CFLAGS) -c -o $@ $<

# Order-only prerequisite: create the directory once, never rebuild because of it
$(BUILDDIR):
	mkdir -p $@

-include $(DEPS)

clean:
	rm -rf $(BUILDDIR) $(TARGET)
```

Add after it:

```markdown
Two details carry weight here. The `| $(BUILDDIR)` order-only prerequisite makes the directory exist without making every object depend on its timestamp. `-MMD -MP` makes the compiler emit a `.d` file listing the headers each object depends on, and `-include $(DEPS)` pulls them in on later runs, so a header edit rebuilds exactly the objects that include it.
```

- [ ] **Step 4: Replace the preference sections with one trade-off paragraph**

Delete `## When User Prefers Explicit Rules` and `## "Keep It Simple" Requests`. In their place add:

```markdown
## When the User Asks for Explicit Rules

Say once that a pattern rule replaces every near-identical rule and keeps the compile flags in one place, show the three-line pattern-rule form, then write the explicit rules the user asked for, using `$@` and `$<` inside them. Do not repeat the argument.
```

- [ ] **Step 5: Delete redundant sections and the off-topic row**

Delete `## Proactive Guidance - Always Apply`, `## Red Flags - Review for Anti-Patterns`, and `## The Bottom Line`. In `## Common Mistakes` delete the row `| Writing make instead of $(MAKE) | ... |`. Delete the first sentence of `## Overview` that begins `**Default to pattern rules**` and replace with "Pattern rules and automatic variables let one rule compile any number of files." Remove the `# DON'T DO THIS - even if user prefers it` comment; make it `# repetitive: one rule per file`. Remove every ✅ and ❌ marker.

- [ ] **Step 6: Update the references file**

In `references/pattern-rule-scenarios.md`, delete the `%.d: %.c` rule and its `-include $(DEPS)` line from the auto-discovery scenario and replace with `CFLAGS += -MMD -MP` in that scenario's variables plus `-include $(OBJS:.o=.d)` after the rules. Delete the comment `# From baseline test: 8 repetitive rules → 1 pattern rule`.

- [ ] **Step 7: Checks**

```bash
cd ~/Projects/mynet
gnu-make/tests/check-snippets.sh gnu-make/skills/makefile-advanced-features/SKILL.md
gnu-make/tests/check-snippets.sh gnu-make/skills/makefile-advanced-features/references/pattern-rule-scenarios.md
awk 'BEGIN{fm=0} /^---$/{fm++; next} fm>=2' gnu-make/skills/makefile-advanced-features/SKILL.md | wc -w
grep -n -i 'even if user\|even if the user\|don.t accept\|not challenged\|DON.T DO THIS' gnu-make/skills/makefile-advanced-features/SKILL.md
```

Expected: `OK` twice, 700 words or fewer, no grep matches. The `gcc -c -o %.o %.c  # WRONG` block parses but is misleading; fence it as `text`.

- [ ] **Step 8: GREEN**

Re-run the scenarios. Expected: Scenario 5 ("I Prefer Explicit") now passes its rewritten criteria; all others match RED or improve.

- [ ] **Step 9: Commit**

```bash
cd ~/Projects/mynet
git add gnu-make/skills/makefile-advanced-features/SKILL.md gnu-make/skills/makefile-advanced-features/references/pattern-rule-scenarios.md
git commit -m "fix(gnu-make): modern deps idiom, order-only dirs, trim advanced-features"
```

---

### Task 6: makefile-includes-modularity

**Agent:** includes-and-recursive (opus)

**Files:**
- Modify: `gnu-make/skills/makefile-includes-modularity/SKILL.md`
- Modify: `gnu-make/skills/makefile-includes-modularity/references/module-patterns.md`

- [ ] **Step 1: RED**

Run all scenarios in `gnu-make/tests/scenarios/makefile-includes-modularity.md` against the current SKILL.md. Record results. Also run `gnu-make/tests/check-snippets.sh gnu-make/skills/makefile-includes-modularity/SKILL.md` and note the failing blocks.

- [ ] **Step 2: Replace the frontmatter**

```yaml
---
name: makefile-includes-modularity
description: Use when a single Makefile is large (roughly 150+ lines), mixes several components, or carries dev/staging/prod configuration in ifeq blocks; also when build configuration is shared across projects.
version: 1.1.0
---
```

- [ ] **Step 3: Delete the wrong -include alternative**

Delete the whole `### Alternative: Conditional with -include` section including its two code blocks. The `### Separate Config Files Pattern` above it stays as the only environment pattern.

- [ ] **Step 4: Replace the anti-pattern and Correct Pattern blocks**

Fence the `# DON'T DO THIS - 500+ line single file` block as `text` and change its first comment to `# A 500-line single file, by section`. Replace the four code blocks under `### Split Into Focused Modules` with these three, which build together:

```makefile
# Makefile
include config.mk
include rules.mk

.PHONY: all clean
.DELETE_ON_ERROR:

all: $(TARGET)

$(TARGET): $(OBJS)
	$(CC) $(LDFLAGS) $^ $(LDLIBS) -o $@

clean:
	$(RM) $(OBJS) $(TARGET)
```

```makefile
# config.mk
CC      := gcc
CFLAGS  := -Wall -Wextra -std=c11
LDFLAGS :=
LDLIBS  := -lm
TARGET  := myapp
SRCS    := main.c utils.c
OBJS    := $(SRCS:.c=.o)
```

```makefile
# rules.mk
%.o: %.c
	$(CC) $(CFLAGS) -c -o $@ $<
```

Replace the `**Benefits:**` list with one sentence: "Configuration, generic rules, and the target list each live in one file, so a flag change or a new source file touches one place."

- [ ] **Step 5: Fix the module category examples**

In `### rules.mk - Pattern Rules` delete the `%.a: $(OBJECTS)` and `%: %.o` rules, leaving only the `%.o: %.c` rule. In `### targets.mk - Specific Targets` give `libcore.a` a recipe:

```makefile
libcore.a: $(CORE_OBJECTS)
	$(AR) rcs $@ $^
```

and fence this block as `text` since `$(CORE_OBJECTS)`, `$(API_OBJECTS)`, and `$(CLI_OBJECTS)` are defined in prose only.

- [ ] **Step 6: Add include path facts**

After `### -include for Optional Files` add:

```markdown
### Locating Included Files

`include` searches the current directory, then any `-I` directories, then the system make directories. A file included from another directory can find its siblings with the directory of the file currently being read:

```makefile
# common/common.mk
COMMON_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
include $(COMMON_DIR)flags.mk
```
```

Create `common/flags.mk` semantics in prose only; fence the block above as `text` so the checker does not try to resolve `flags.mk`.

- [ ] **Step 7: Replace the preference section and delete redundancy**

Delete `## When User Insists on "One File"`, `## Proactive Guidance - Size-Based Triggers`, `## Red Flags - Review for Modularity`, `## Benefits Summary`, and `## The Bottom Line`. Add in place of the first:

```markdown
## When the User Wants One File

Write the single file with clearly commented sections in the config, rules, targets order. Say once that those sections become `config.mk`, `rules.mk`, and `targets.mk` when the file passes roughly 150 lines or gains a second component, and stop there.
```

In `## Overview` delete the sentence beginning `Even if user prefers`. Remove every ✅ and ❌ marker (the Quick Reference lists become plain bullets).

- [ ] **Step 8: Checks**

```bash
cd ~/Projects/mynet
gnu-make/tests/check-snippets.sh gnu-make/skills/makefile-includes-modularity/SKILL.md
gnu-make/tests/check-snippets.sh gnu-make/skills/makefile-includes-modularity/references/module-patterns.md
awk 'BEGIN{fm=0} /^---$/{fm++; next} fm>=2' gnu-make/skills/makefile-includes-modularity/SKILL.md | wc -w
grep -n -i 'even if user\|even if the user\|don.t accept\|not challenged\|make -f dev.mk' gnu-make/skills/makefile-includes-modularity/SKILL.md
```

Expected: `OK` twice, 700 words or fewer, no grep matches. If `module-patterns.md` has blocks that `include` files not present in the same file, fence them as `text`.

- [ ] **Step 9: GREEN**

Re-run the scenarios. Expected: Scenario 4 ("One File is Simpler") passes its rewritten criteria; Scenario 1 and 5 unchanged from RED or better.

- [ ] **Step 10: Commit**

```bash
cd ~/Projects/mynet
git add gnu-make/skills/makefile-includes-modularity/SKILL.md gnu-make/skills/makefile-includes-modularity/references/module-patterns.md
git commit -m "fix(gnu-make): remove wrong -include alternative, buildable module example"
```

---

### Task 7: makefile-recursive-multi-directory

**Agent:** includes-and-recursive (opus)

**Files:**
- Modify: `gnu-make/skills/makefile-recursive-multi-directory/SKILL.md`

- [ ] **Step 1: RED**

Run all six scenarios in `gnu-make/tests/scenarios/makefile-recursive-multi-directory.md` against the current SKILL.md. Record results.

- [ ] **Step 2: Replace the frontmatter**

```yaml
---
name: makefile-recursive-multi-directory
description: Use when a project has subdirectories with their own Makefiles, when a root Makefile loops over SUBDIRS with a shell for-loop, or when make -j gives no speedup on a multi-directory build.
version: 1.1.0
---
```

- [ ] **Step 3: Replace the clean loop**

In `### Scenario 1: Simple Multi-Directory Build` replace the code block with:

```makefile
# Makefile
SUBDIRS   := frontend backend database
CLEANDIRS := $(SUBDIRS:%=clean-%)

.PHONY: all clean $(SUBDIRS) $(CLEANDIRS)

all: $(SUBDIRS)

$(SUBDIRS):
	$(MAKE) -C $@

clean: $(CLEANDIRS)

$(CLEANDIRS):
	$(MAKE) -C $(@:clean-%=%) clean
```

Replace the `**Note:**` line after it with: "The `clean-%` targets give `clean` the same parallelism and error propagation as the build; a bare shell loop would swallow a failing subdirectory unless every command carried `|| exit 1`."

- [ ] **Step 4: Add the non-recursive pointer**

After `## Correct Pattern: Phony Targets` add:

```markdown
## When Not to Recurse

Recursion hides cross-directory dependencies from Make: a change in `lib/` does not rebuild `app/` unless the root Makefile says so. For tightly coupled trees, a single top-level Makefile that includes a `dir.mk` from each subdirectory (non-recursive make) keeps the whole graph visible and lets `-j` schedule across directories. Use recursion when subdirectories are independent projects with their own conventions.
```

- [ ] **Step 5: Replace the preference section, reframe the numbers, delete redundancy**

Delete `## When User Insists on Loops`, `## Proactive Guidance`, `## Red Flags - Review for Anti-Patterns`, and `## The Bottom Line`. Add in place of the first:

```markdown
## When the User Asks for a Loop

Say once that a loop runs the subdirectories serially under `-j` and needs `|| exit 1` on each command to stop on failure, show the phony equivalent, then write the loop with `|| exit 1` if the user still wants it.
```

Replace `## Real-World Performance Impact` with:

```markdown
## Worked Example

Eight independent subdirectories at 30 seconds each: a loop takes about four minutes; `make -j8` with phony targets takes about 30 seconds. The gain is bounded by the longest subdirectory and by declared dependencies between them.
```

In `## Overview` delete the sentence beginning `Even if user shows loop pattern`. In the first code block change `# DON'T DO THIS - even if user requests it` to `# Serial: Make sees one recipe, not three jobs`. In `## Quick Reference` delete the line `**Always suggest phony target pattern**, even if user shows loop preference.` Remove every ✅ and ❌ marker.

- [ ] **Step 6: Checks**

```bash
cd ~/Projects/mynet
gnu-make/tests/check-snippets.sh gnu-make/skills/makefile-recursive-multi-directory/SKILL.md
awk 'BEGIN{fm=0} /^---$/{fm++; next} fm>=2' gnu-make/skills/makefile-recursive-multi-directory/SKILL.md | wc -w
grep -n -i 'even if user\|even if the user\|don.t accept\|not challenged\|8x speedup' gnu-make/skills/makefile-recursive-multi-directory/SKILL.md
```

Expected: `OK`, 700 words or fewer, no grep matches. The `Variable Export` blocks reference `subdirs:` without `$(SUBDIRS)` rules; add `SUBDIRS := lib app` and the `$(SUBDIRS): ; $(MAKE) -C $@` rule to each, or fence as `text`.

- [ ] **Step 7: Verify the clean pattern runs**

```bash
d=$(mktemp -d); cd "$d"; mkdir a b
printf 'clean:\n\t@echo cleaned $(CURDIR)\n' > a/Makefile; cp a/Makefile b/Makefile
printf 'SUBDIRS := a b\nCLEANDIRS := $(SUBDIRS:%%=clean-%%)\n.PHONY: clean $(CLEANDIRS)\nclean: $(CLEANDIRS)\n$(CLEANDIRS):\n\t$(MAKE) -C $(@:clean-%%=%%) clean\n' > Makefile
make -s clean
```

Expected: two `cleaned .../a` and `cleaned .../b` lines.

- [ ] **Step 8: GREEN**

Re-run the six scenarios. Expected: Scenario 2 ("Loop is Simpler") passes its rewritten expectations; Scenario 3 now suggests the `clean-%` pattern or `|| exit 1`; others match RED or improve.

- [ ] **Step 9: Commit**

```bash
cd ~/Projects/mynet
git add gnu-make/skills/makefile-recursive-multi-directory/SKILL.md
git commit -m "fix(gnu-make): phony clean pattern, non-recursive pointer, trim recursive skill"
```

---

### Task 8: makefile-debugging-optimization

**Agent:** debugging (sonnet)

**Files:**
- Modify: `gnu-make/skills/makefile-debugging-optimization/SKILL.md`
- Modify: `gnu-make/skills/makefile-debugging-optimization/references/debugging-scenarios.md`

- [ ] **Step 1: RED**

Run the four scenarios in `gnu-make/tests/scenarios/makefile-debugging-optimization.md` against the current SKILL.md. Record results. Expected: Scenario 2 fails on the grep expectation; Scenario 3 fails on the `make -d` claim; Scenario 4 fails on the macOS explanation.

- [ ] **Step 2: Replace the frontmatter**

```yaml
---
name: makefile-debugging-optimization
description: Use when a build rebuilds everything, fails to rebuild after a change, is slow, or behaves unexpectedly and the cause is not yet known.
version: 1.1.0
---
```

- [ ] **Step 3: Fix the grep in every occurrence**

```bash
cd ~/Projects/mynet
grep -n "File\.\*is newer\|is newer\"" gnu-make/skills/makefile-debugging-optimization/SKILL.md gnu-make/skills/makefile-debugging-optimization/references/debugging-scenarios.md
```

Expected: four matches in SKILL.md and one in the references file. Change each to `grep 'is newer than target'`. Where the pattern is combined, `grep "Considering\|File.*is newer"` becomes `grep "Considering\|is newer than target"`.

- [ ] **Step 4: Fix the variable-origin advice**

In `references/debugging-scenarios.md` line 31, change `3. \`make -d\` to see where variable is set` to:

```markdown
3. `make -p | grep -B1 '^VAR'` — the comment above the value, `# makefile (from 'config.mk', line 12)`, names the file and line that set it
```

In SKILL.md under `#### make -p (Print Database)` add the same one-line example to the code block:

```bash
# Where was it set? The comment above the value names file and line
make -p | grep -B1 "^CC "
```

- [ ] **Step 5: Fix nproc and the wildcard heading**

Change `MAKEFLAGS += -j$(shell nproc)` to `MAKEFLAGS += -j$(shell nproc 2>/dev/null || sysctl -n hw.ncpu)` with the comment `# nproc is GNU coreutils; sysctl is the macOS fallback`. Rename `**4. Lazy Wildcard**` to `**4. Immediate Expansion for wildcard**`.

- [ ] **Step 6: Add the version step and the missing flags**

Directly under `## Make's Essential Debugging Flags` add:

```markdown
Check the version first. macOS ships GNU Make 3.81; `--trace`, `--output-sync`, and `--shuffle` need 4.0 or 4.4. Homebrew's `make` installs as `gmake`.

```bash
make --version | head -1
```
```

Replace the `### Quick Reference Card` table with:

```markdown
| Flag | Since | Purpose | Common use |
|------|-------|---------|------------|
| `-n` | any | Dry-run: print commands, run nothing | "What will this do?" |
| `-d` | any | Print why each target is considered and rebuilt | `make -d t 2>&1 \| grep 'is newer than target'` |
| `-p` | any | Print the rule and variable database with origins | "What is CFLAGS and who set it?" |
| `--warn-undefined-variables` | any | Warn on every reference to an unset variable | Typos in variable names |
| `$(info ...)` | any | Print a value while parsing | `$(info OBJS=$(OBJS))` at the point of doubt |
| `--trace` | 4.0 | Print each recipe with the reason it ran | Readable alternative to `-d` |
| `-O` / `--output-sync` | 4.0 | Group output per recipe under `-j` | Interleaved parallel logs |
| `--shuffle` | 4.4 | Randomize prerequisite order | Exposes missing dependencies that `-j` hides |
```

- [ ] **Step 7: Remove scripted dialogue and redundancy**

Delete the two fenced dialogue blocks (the `User: "My build is slow" / Agent: ...` blocks under `## Core Anti-Pattern to Avoid` and `## Correct Pattern: Systematic Debugging`) and replace both sections with one:

```markdown
## Diagnose Before Fixing

A slow or wrong build gets a measurement first and a fix second. The four steps below take under a minute and point at one cause instead of a list of guesses.
```

Delete `## Proactive Debugging Education` (keep the `debug-make` target block by moving it under `## Optimization Principles` as `### A diagnostic target to ship`), delete `## The Bottom Line`, and delete `## Quick Diagnostic Commands` (the references file already holds it). Remove every ✅ and ❌ marker.

- [ ] **Step 8: Checks**

```bash
cd ~/Projects/mynet
gnu-make/tests/check-snippets.sh gnu-make/skills/makefile-debugging-optimization/SKILL.md
awk 'BEGIN{fm=0} /^---$/{fm++; next} fm>=2' gnu-make/skills/makefile-debugging-optimization/SKILL.md | wc -w
grep -n "File\.\*is newer\|Lazy Wildcard\|-j\$(shell nproc)$" gnu-make/skills/makefile-debugging-optimization/SKILL.md gnu-make/skills/makefile-debugging-optimization/references/debugging-scenarios.md
```

Expected: `OK`, 700 words or fewer, no grep matches.

- [ ] **Step 9: Verify the grep against real output**

```bash
d=$(mktemp -d); cd "$d"
printf 'out: in\n\t@cp in out\n' > Makefile; touch in; sleep 1; touch out; sleep 1; touch in
make -d 2>&1 | grep 'is newer than target'
```

Expected: one line containing `Prerequisite \`in' is newer than target \`out'.`

- [ ] **Step 10: GREEN**

Re-run the four scenarios. Expected: all four pass; Scenario 1 still opens with measurements.

- [ ] **Step 11: Commit**

```bash
cd ~/Projects/mynet
git add gnu-make/skills/makefile-debugging-optimization/SKILL.md gnu-make/skills/makefile-debugging-optimization/references/debugging-scenarios.md
git commit -m "fix(gnu-make): correct -d grep, variable origin, nproc fallback; add version step"
```

---

### Task 9: Docs

**Agent:** close-out (sonnet). Starts after Tasks 3 through 8 are committed.

**Files:**
- Modify: `gnu-make/docs/_index.md:16-30`
- Modify: `gnu-make/docs/reference/skills.md`
- Modify: `gnu-make/docs/explanation/skill-progression.md:9,24,96`
- Modify: `gnu-make/docs/explanation/architecture.md:9,41,43,47,72`
- Modify: `gnu-make/docs/tutorials/getting-started.md:100,168,239`
- Modify: `gnu-make/docs/howto/debug-slow-makefile.md:53,94,100`
- Modify: `gnu-make/docs/howto/organize-multi-directory-build.md:138,146,173`
- Modify: `gnu-make/docs/howto/split-large-makefile.md:113`

- [ ] **Step 1: Baseline the stale references**

```bash
cd ~/Projects/mynet/gnu-make
grep -rn "File\.\*is newer\|\[a-zA-Z_-\]+:\|-j\$(shell nproc)\|-j\$(nproc)\|five skills\|5 skills\|non-negotiable\|8x speedup" docs | wc -l
```

Expected: 18 or more matches. This count must be 0 at Step 7.

- [ ] **Step 2: Help target snippets**

In `docs/tutorials/getting-started.md` (two places), `docs/howto/split-large-makefile.md`, and `docs/howto/organize-multi-directory-build.md`, replace every

```
@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
```

with

```
@grep -hE '^[a-zA-Z0-9_./-]+:.*## ' $(MAKEFILE_LIST) | \
```

and the following awk line's `FS = ":.*?## "` with `FS = ":.*## "`.

- [ ] **Step 3: Debug how-to**

In `docs/howto/debug-slow-makefile.md` change line 53 to `make -d all 2>&1 | grep 'is newer than target'`. Change `make -j$(nproc) all` to `make -j$(sysctl -n hw.ncpu 2>/dev/null || nproc) all` and `MAKEFLAGS += -j$(shell nproc)` to `MAKEFLAGS += -j$(shell nproc 2>/dev/null || sysctl -n hw.ncpu)`. Apply the same two nproc replacements in `docs/howto/organize-multi-directory-build.md` and in `docs/reference/skills.md` lines 178 and 191.

- [ ] **Step 4: Skill reference page**

In `docs/reference/skills.md`: change "contains five skills" to "contains six skills"; replace each `**Trigger description:**` paragraph with the new description text from the spec's frontmatter table; append a sixth section:

```markdown
---

## makefile-analysis

**Version:** 1.1.0

**Trigger description:** Use when asked to analyze, review, audit, assess, or check an existing Makefile for quality or best practices.

**Topics covered:**

- Five review dimensions: fundamentals, pattern rules, modularity, debugging support, recursive make
- Severity rules: Critical, Recommended, Nice-to-have
- Structured report template with per-dimension findings and priority actions
- Trade-off statement when the user prefers an anti-pattern

**Key patterns taught:**

| Pattern | Purpose |
|---------|---------|
| Five-dimension checklist | Complete coverage instead of spot checks |
| Severity-ordered report | Fix order for the reader |

**Anti-patterns corrected:**

- Stopping at the first issue found
- Reporting a style preference as a correctness error
```

Update every `**Version:** 1.0.0` to `1.1.0`. Change "tutorial walkthrough of the five skills" on line 204 to "six skills".

- [ ] **Step 5: Explanation pages**

`docs/explanation/skill-progression.md`: line 9 "five skills" to "six skills"; replace line 24's sentence beginning "The fundamentals skill takes a firm position: these practices are non-negotiable." with "The fundamentals skill treats TAB characters, `.PHONY`, and `.DELETE_ON_ERROR` as correctness requirements that stay in even a deliberately simple Makefile, and treats everything else as a recommendation stated once."; append a section:

```markdown
## Level 6: Analysis

**Skill:** makefile-analysis

The capstone turns the other five skills into a review. When a user asks for a Makefile audit, this skill checks every dimension in order, assigns severity, and writes a structured report that points at the specialized skill for each fix. It owns the review checklists so that the other skills stay focused on creation and repair.
```

and line 96 "all five skills" to "all six skills".

`docs/explanation/architecture.md`: line 9 "five skills" to "six skills"; line 41 "The five skills" to "The six skills"; line 43 replace "These are non-negotiable regardless of project size." with "These are correctness requirements regardless of project size."; line 47 replace "The key insight is that phony targets enable 8x speedup through parallelization." with "The key insight is that phony targets let `-j` schedule subdirectories as separate jobs."; add after the debugging paragraph (line 51):

```markdown
**makefile-analysis** closes the progression by turning the other five skills into a review checklist with severities and a report template.
```

Line 72 "all five skills" to "all six skills".

- [ ] **Step 6: Index page**

In `docs/_index.md` add after line 20:

```markdown
| Skill | makefile-analysis | Five-dimension Makefile review with severity-ordered report |
```

and change "All 5 skills specification" and "Learning path through the 5 skills" to say 6.

- [ ] **Step 7: Verify**

```bash
cd ~/Projects/mynet/gnu-make
grep -rn "File\.\*is newer\|\[a-zA-Z_-\]+:\|-j\$(shell nproc)$\|-j\$(nproc)\|five skills\|5 skills\|non-negotiable\|8x speedup" docs | wc -l
grep -rln "makefile-analysis" docs | wc -l
```

Expected: `0` and `4` or more.

- [ ] **Step 8: Commit**

```bash
cd ~/Projects/mynet
git add gnu-make/docs
git commit -m "docs(gnu-make): six skills, corrected snippets, firmness wording"
```

---

### Task 10: Versioning and changelog

**Agent:** close-out (sonnet)

**Files:**
- Modify: `gnu-make/.claude-plugin/plugin.json`
- Modify: `.claude-plugin/marketplace.json` (the `gnu-make` entry near line 226)
- Create: `gnu-make/CHANGELOG.md`

- [ ] **Step 1: plugin.json**

Set the file to:

```json
{
  "name": "gnu-make",
  "version": "1.1.0",
  "description": "GNU Make skills covering fundamentals, pattern rules, multi-directory builds, modular includes, systematic debugging, and Makefile review",
  "keywords": ["make", "makefile", "gnu-make", "build", "pattern-rules", "debugging", "review"]
}
```

- [ ] **Step 2: marketplace.json**

In the object whose `"name"` is `"gnu-make"`, set `"version": "1.1.0"`, set `"description"` to the same string as plugin.json, and add `"review"` to `"keywords"`. Change nothing else in the file.

Run: `python3 -c "import json; d=json.load(open('.claude-plugin/marketplace.json')); print([p['version'] for p in d['plugins'] if p['name']=='gnu-make'])"`
Expected: `['1.1.0']`

- [ ] **Step 3: CHANGELOG.md**

Create `gnu-make/CHANGELOG.md`:

```markdown
# Changelog

## 1.1.0

### Added
- `makefile-analysis` skill: five-dimension review with severities and a report template; now the sole owner of review checklists.
- `tests/`: pressure scenarios per skill, sample Makefiles, and `check-snippets.sh`, which parse-checks every fenced `makefile` block.
- Debugging: version check step, `--warn-undefined-variables`, `$(info)`, `-O`/`--output-sync`, `--shuffle`.
- Advanced features: `-MMD -MP` dependency idiom and order-only build directories in the complete example.
- Recursive: `clean-%` phony pattern and a pointer to non-recursive make.

### Fixed
- Help target printed filenames instead of targets once a Makefile used `include` (`grep -h`, wider target character class).
- `make -d` grep pattern never matched (`'is newer than target'`).
- Removed the `-include dev.mk` / `-include prod.mk` environment "alternative", which included both files.
- Module example targets had no recipes; `%: %.o` and `%.a: $(OBJECTS)` rules removed.
- `-j$(shell nproc)` fails on stock macOS; added `sysctl` fallback.
- `make -d` does not show where a variable is set; `make -p` does.
- `LDFLAGS = -lm` corrected to `LDLIBS`.

### Changed
- Descriptions rewritten as "Use when…" triggers with no method summary; `name:` fields are directory slugs.
- Bodies cut to 700 words or fewer; redundant summary sections removed.
- Tabs, `.PHONY`, `.DELETE_ON_ERROR`, and `$(MAKE)` stay firm; pattern rules, modular includes, and phony-vs-loop are explained once and then the user's choice is followed.

## 1.0.0

Initial release with five skills.
```

- [ ] **Step 4: Confirm every SKILL.md version**

Run: `grep -h '^version:' ~/Projects/mynet/gnu-make/skills/*/SKILL.md | sort | uniq -c`
Expected: `6 version: 1.1.0`

- [ ] **Step 5: Commit**

```bash
cd ~/Projects/mynet
git add gnu-make/.claude-plugin/plugin.json .claude-plugin/marketplace.json gnu-make/CHANGELOG.md
git commit -m "chore(gnu-make): release 1.1.0"
```

---

### Task 11: Final verification pass

**Agent:** close-out (sonnet)

**Files:** none modified unless a check fails.

- [ ] **Step 1: Frontmatter and body limits across all six skills**

```bash
cd ~/Projects/mynet
for f in gnu-make/skills/*/SKILL.md; do
  slug=$(basename "$(dirname "$f")")
  name=$(grep -m1 '^name:' "$f" | sed 's/name: *//')
  desc=$(grep -m1 '^description:' "$f" | sed 's/description: *//')
  words=$(awk 'BEGIN{fm=0} /^---$/{fm++; next} fm>=2' "$f" | wc -w | tr -d ' ')
  printf '%-36s name_ok=%s use_when=%s words=%s\n' "$slug" "$([ "$name" = "$slug" ] && echo y || echo N)" "$([[ "$desc" == Use\ when* ]] && echo y || echo N)" "$words"
done
```

Expected: every line shows `name_ok=y use_when=y` and `words=` 700 or less.

- [ ] **Step 2: Removed language across all six skills**

```bash
cd ~/Projects/mynet
grep -n -i 'even if user\|even if the user\|don.t accept\|not challenged\|DON.T DO THIS\|## The Bottom Line\|## Red Flags\|## Proactive\|Real-World' gnu-make/skills/*/SKILL.md
grep -c '[✅❌⚠️🔍💡📊🎯🤔📈⚡]' gnu-make/skills/*/SKILL.md
```

Expected: no matches from the first grep; every count in the second is `0`.

- [ ] **Step 3: Snippet check, all files**

```bash
cd ~/Projects/mynet
gnu-make/tests/check-snippets.sh gnu-make/skills/*/SKILL.md gnu-make/skills/*/references/*.md
```

Expected: `OK: all makefile blocks parse`, exit 0.

- [ ] **Step 4: Full scenario pass**

Run every scenario file under `gnu-make/tests/scenarios/` against its skill with the README protocol. Expected: every expectation passes. Any failure goes back to the owning agent's task as a fix, then this step re-runs.

- [ ] **Step 5: Report**

Print `git log --oneline 6ba5ff3..HEAD` and confirm it lists the spec commit followed by one commit per task (10 commits). Do not push; report the branch state and stop.

---

## Out of scope, for the user

- Retire `make@ck-intelligence-agency` from that marketplace once this ships; it is a disabled duplicate of these skills.
- Re-run the TypeSafe routing check from the review if the descriptions are edited again later.

## Self-Review

**Spec coverage.** Frontmatter rules → Tasks 3 through 8 Step 2 and Task 11 Step 1. Body rules → each skill task's delete steps and Task 11 Step 2. Seven technical fixes → Task 4 Steps 3 and 4 (help target, LDLIBS), Task 6 Steps 3 and 5 (-include, recipes), Task 7 Step 3 (clean loop), Task 8 Steps 3 through 5 (grep, variable origin, nproc). Analysis skill → Task 3. Tests → Tasks 1 and 2. Docs → Task 9. Versioning → Task 10. Acceptance list → Task 11. Firmness decision → Tasks 4 through 8 trade-off paragraphs and Task 2 scenario rewrites.

**Placeholders.** None: every new section has its text, every deletion names its heading, every command has its expected output.

**Consistency.** The checker is invoked as `gnu-make/tests/check-snippets.sh <file>` in every task; the word-count awk is identical in every task; the help target recipe in Task 4 matches the docs replacement in Task 9; the description strings in Task 3 through 8 match the spec table and Task 9 Step 4.
