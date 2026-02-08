---
title: "Skill Progression Path"
description: "How the five gnu-make skills build on each other in a deliberate learning path from fundamentals to debugging mastery"
weight: 2
---

# Skill Progression Path

The five skills in the gnu-make plugin follow a deliberate progression. Each skill builds on concepts from the previous one, and together they form a complete learning path from first Makefile to expert debugging. This progression is not accidental; it mirrors the way Make expertise develops in practice.

## Level 1: Fundamentals

**Skill:** makefile-fundamentals

The foundation covers what every Makefile must have, regardless of project size or complexity.

- **TAB characters** for recipe indentation. This is a hard requirement of Make's parser. Spaces produce the notorious "missing separator" error and are the most common stumbling point for newcomers.
- **`.PHONY` declarations** for non-file targets. Without this, a file named `clean` in the project directory silently prevents `make clean` from running. This is a correctness issue, not a stylistic preference.
- **`.DELETE_ON_ERROR`** for build robustness. When a recipe fails partway through writing a target file, this directive removes the partial output so the next build does not incorrectly consider it up to date.
- **The `##` help pattern** for discoverability. Four lines of boilerplate give the Makefile a `make help` command that documents every target. New team members can orient themselves immediately.
- **Standard variables** (`CC`, `CFLAGS`, `LDFLAGS`) for consistency. These names are conventional across the Make ecosystem and allow compiler or flag changes in a single place.
- **Basic pattern rules** (`%.o: %.c`) and automatic variables (`$@`, `$<`, `$^`) as a first introduction. The fundamentals skill shows that one pattern rule replaces N explicit rules, but does not go into advanced usage.

The fundamentals skill takes a firm position: these practices are non-negotiable. Even when a user asks for a "simple" Makefile, TAB characters, `.PHONY`, and the help target should be present. Simplicity means fewer features, not fewer correctness guarantees.

## Level 2: Advanced Features

**Skill:** makefile-advanced-features

Once the basics are solid, the advanced features skill deepens the user's command of DRY patterns. The central lesson is counterintuitive to many newcomers: pattern rules produce less code, not more complexity.

- **Pattern rules in depth**: the `%` stem, how Make matches targets to prerequisites, and when to use static pattern rules for a specific subset of files.
- **All automatic variables**: `$@` (target), `$<` (first prerequisite), `$^` (all prerequisites), `$?` (changed prerequisites), `$*` (stem). Understanding these eliminates hardcoded filenames from recipes.
- **`$(wildcard)`** for automatic file discovery. The Makefile adapts to new source files without manual edits.
- **`$(patsubst)` and `$(foreach)`** for transformations and iteration within variable assignments.
- **Conditional compilation** with `ifeq` and `ifdef` for debug/release mode switching.
- **Target-specific variables** for per-file flag overrides without breaking the pattern rule.

The core anti-pattern this skill corrects is repetitive explicit rules. When a user writes separate compilation rules for each `.c` file, the advanced features skill demonstrates that a single pattern rule handles unlimited files with less code and zero copy-paste inconsistencies. The mantra: "pattern rules are simpler -- less code, fewer places to edit."

## Level 3: Multi-Directory Builds

**Skill:** makefile-recursive-multi-directory

Projects grow beyond a single directory. This skill teaches how to coordinate builds across subdirectories without sacrificing parallelism.

- **The phony target pattern**: each subdirectory becomes a Make target, enabling `make -j` to build them in parallel.
- **Why shell loops fail**: a for-loop in a recipe is a single shell command, invisible to Make's job scheduler. The `-j` flag has no effect, and `-k` (continue on error) breaks.
- **`$(MAKE)` for recursive invocations**: ensures that flags like `-n`, `-t`, and `-q` propagate to sub-makes. Hardcoding `make` silently breaks dry-run and touch modes.
- **Dependency declarations**: `app: lib` tells Make to build `lib` before `app`, while still building `app` and `tests` in parallel when both depend on `lib`.
- **Variable export**: passing configuration from the root Makefile to subdirectory Makefiles via `export` or command-line overrides.

The key insight is quantitative. For a project with 8 subdirectories, each taking 30 seconds to build, the phony target pattern with `make -j8` completes in 30 seconds instead of 4 minutes. That is an 8x speedup from a structural change, not a compiler optimization.

## Level 4: Modular Organization

**Skill:** makefile-includes-modularity

As individual Makefiles grow past 150 lines, they become hard to navigate and cause merge conflicts. This skill teaches how to split them.

- **The `include` directive**: Make suspends reading the current file, reads the included file, then resumes. Variables and rules from included files are available in the main Makefile.
- **`-include` for optional files**: local developer overrides (`local.mk`) and auto-generated dependency files (`.d`) should use `-include` to avoid errors when the file does not exist.
- **Standard module names**: `config.mk` for variables, `rules.mk` for pattern rules, `targets.mk` for specific targets, `install.mk` for installation logic.
- **Environment-specific configurations**: `config/dev.mk` and `config/prod.mk` selected by `ENV ?= dev` with `include config/$(ENV).mk`. Adding an environment means adding a file, not editing conditional blocks.
- **Cross-project shared configuration**: `include ../common/common.mk` lets multiple projects share compiler settings, pattern rules, and conventions.

The threshold of 150 lines is a guideline, not a hard rule. The real signal is mixed concerns: when configuration, pattern rules, and specific targets are interleaved in a single file, splitting improves clarity regardless of line count.

## Level 5: Debugging and Optimization

**Skill:** makefile-debugging-optimization

The final skill teaches the diagnostic methodology that separates proficient Make users from experts.

- **`make -n` (dry-run)**: see what commands would execute without running them. Essential for verifying Makefile logic before committing changes.
- **`make -d` (debug mode)**: see Make's decision-making process. Why it considers a target out of date, which rule it selects, how it resolves dependencies.
- **`make -p` (print database)**: dump all rules, variables, and their values. The definitive answer to "what does Make think this variable is?"
- **`make --trace` (GNU Make 4.0+)**: a more readable alternative to `-d` that shows which rules fire and what triggered them.
- **Systematic methodology**: measure first (`time make`), understand the problem (`make -n`, `make -d`), identify the root cause, then apply a targeted fix.

The core anti-pattern this skill corrects is jumping to solutions without diagnosis. When a user reports "my build is slow," the instinct is to suggest parallelism, ccache, and flag optimization. The debugging skill insists on measurement and diagnosis first. The actual problem might be a timestamp issue, an undeclared dependency, or a recursively-expanded variable re-evaluating `$(wildcard)` on every reference. Without diagnosis, optimizations are guesses.

## The Progression Is Intentional

Each level in this progression assumes knowledge from the previous levels:

- **Advanced features** assumes you understand TAB characters, `.PHONY`, and basic pattern rules from fundamentals.
- **Multi-directory builds** assumes you can write correct individual Makefiles using pattern rules and automatic variables from advanced features.
- **Modular organization** assumes you have Makefiles worth splitting, which implies multi-directory or complex single-directory builds.
- **Debugging** assumes you have a build system complex enough to need diagnosis, and that you understand the constructs (pattern rules, variables, includes) that the diagnostic output reveals.

This progression means a user working through the skills in order builds a complete, layered understanding of GNU Make. But the skills also work independently. A user who only needs to debug a slow build can use the debugging skill without having studied modularity. The progression describes the ideal learning path, not a prerequisite chain.

## See Also

- [Skill Reference]({{< ref "reference/skills" >}}) -- technical specifications for all five skills
- [Architecture]({{< ref "explanation/architecture" >}}) -- how the plugin is structured
- [Debug a Slow Makefile]({{< ref "howto/debug-slow-makefile" >}}) -- practical application of the debugging skill
