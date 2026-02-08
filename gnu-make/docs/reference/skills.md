---
title: "Skill Reference"
description: "Technical specifications for all five gnu-make plugin skills including triggers, topics, patterns, and anti-patterns"
weight: 1
---

# Skill Reference

The gnu-make plugin contains five skills. Each skill provides knowledge injection for a specific area of GNU Make expertise. Skills activate automatically based on trigger patterns in the user's request.

---

## makefile-fundamentals

**Version:** 1.0.0

**Trigger description:** Activates when the user asks to create, write, generate, fix, or review a Makefile; resolve a "missing separator" error; add build targets; set up a help target; or implement self-documenting Makefiles with the `##` comment pattern.

**Topics covered:**

- Self-documenting `##` help pattern with `grep` and `awk`
- TAB character requirement for recipe indentation
- `.PHONY` declarations for non-file targets
- `.DELETE_ON_ERROR` for build robustness
- Automatic variables: `$@`, `$<`, `$^`, `$*`
- Standard variables: `CC`, `CFLAGS`, `LDFLAGS`, `TARGET`, `SOURCES`, `OBJECTS`
- Basic pattern rules (`%.o: %.c`)
- Substitution syntax for variable transformation (`$(SOURCES:.c=.o)`)

**Key patterns taught:**

| Pattern | Purpose |
|---------|---------|
| `help: ## Show help` with `grep`/`awk` recipe | Self-documenting help target |
| `.PHONY: all clean test` | Declare non-file targets |
| `.DELETE_ON_ERROR:` | Remove targets on recipe failure |
| `$(SOURCES:.c=.o)` | Derive object list from source list |

**Anti-patterns corrected:**

- Spaces instead of TAB characters in recipes
- Missing `.PHONY` declarations for non-file targets
- Hardcoded filenames in recipes instead of automatic variables
- Missing help target
- Missing `.DELETE_ON_ERROR` declaration

---

## makefile-advanced-features

**Version:** 1.0.0

**Trigger description:** Activates when the user asks about pattern rules, automatic variables in depth, wildcard functions, `patsubst`, reducing Makefile duplication, compiling multiple C files with one rule, static pattern rules, conditional compilation with `ifeq`, or target-specific variables. Also activates when the user shows repetitive explicit rules.

**Topics covered:**

- Pattern rule syntax and the `%` stem
- All automatic variables: `$@`, `$<`, `$^`, `$?`, `$*`
- `$(wildcard)` for automatic file discovery
- `$(patsubst)` for pattern substitution
- `$(foreach)` for iteration
- Static pattern rules for subset matching
- Conditional compilation with `ifeq`/`ifdef`
- Target-specific variable overrides
- Debug vs release flag switching with `DEBUG ?= 0`

**Key patterns taught:**

| Pattern | Purpose |
|---------|---------|
| `%.o: %.c` with `$(CC) $(CFLAGS) -c -o $@ $<` | Generic compilation rule |
| `SRCS := $(wildcard *.c)` | Auto-discover source files |
| `OBJS := $(SRCS:.c=.o)` | Derive object files from sources |
| `$(OBJS): %.o: %.c` | Static pattern rule for a subset |
| `ifeq ($(DEBUG),1)` | Conditional flag switching |
| `parser.o: CFLAGS += -Wno-unused-function` | Target-specific variable |

**Anti-patterns corrected:**

- Repetitive explicit rules for similar targets (N files = N rules)
- Hardcoded filenames in recipes instead of automatic variables
- Manually listing source files when `$(wildcard)` is appropriate
- Repeating substitution logic instead of using `$(patsubst)`
- Using `make` instead of `$(MAKE)` in recursive contexts

---

## makefile-recursive-multi-directory

**Version:** 1.0.0

**Trigger description:** Activates when the user asks about multi-directory Makefile builds, recursive make, building subdirectories, using `make -C`, structuring a project with multiple Makefiles, fixing `make -j` parallelism issues, replacing shell loops with phony targets, exporting variables to sub-makes, or using the `$(MAKE)` variable.

**Topics covered:**

- Phony target pattern for subdirectory builds
- Shell loop anti-pattern and its limitations
- `$(MAKE) -C $@` for recursive invocations
- Dependency declarations between subdirectories
- Variable export with `export` directive
- Command-line variable passing to sub-makes
- `.EXPORT_ALL_VARIABLES` special target
- `unexport` for selective exclusion
- Parallel build enabling with `make -j`
- Error handling differences between loops and phony targets

**Key patterns taught:**

| Pattern | Purpose |
|---------|---------|
| `.PHONY: $(SUBDIRS)` with `$(MAKE) -C $@` | Parallelizable subdirectory builds |
| `app: lib` (prerequisite declaration) | Inter-directory dependency ordering |
| `export VERSION CFLAGS` | Pass variables to sub-makes |
| `$(MAKE) -C $@ VERSION=$(VERSION)` | Override sub-make variables |

**Anti-patterns corrected:**

- Shell for-loop for building subdirectories (prevents `-j` parallelization)
- Using hardcoded `make` instead of `$(MAKE)` for recursive calls
- Missing `.PHONY` declarations for subdirectory targets
- Missing dependency declarations between dependent subdirectories
- Missing `|| exit 1` in loop-based builds (error swallowing)

---

## makefile-includes-modularity

**Version:** 1.0.0

**Trigger description:** Activates when the user asks about splitting a Makefile into multiple files, using the `include` or `-include` directive, organizing a large Makefile, creating environment-specific configurations, sharing Makefile configuration across projects, or when a Makefile exceeds 150 lines.

**Topics covered:**

- `include` directive syntax and behavior
- `-include` for optional files (no error on missing)
- Size thresholds for suggesting modularization (150+ lines)
- Standard module naming: `config.mk`, `rules.mk`, `targets.mk`, `install.mk`
- Environment-specific configuration with `config/dev.mk`, `config/prod.mk`
- `ENV ?= dev` with `include config/$(ENV).mk` for environment selection
- Shared configuration across projects via `include ../common/common.mk`
- Module size guidelines (50-300 lines per module)
- Four module patterns: three-file foundation, component-based, environment + component, shared common

**Key patterns taught:**

| Pattern | Purpose |
|---------|---------|
| `include config.mk rules.mk targets.mk` | Split monolith into focused modules |
| `-include local.mk` | Optional local overrides |
| `ENV ?= dev` / `include config/$(ENV).mk` | Environment-specific builds |
| `include ../common/common.mk` | Cross-project shared configuration |

**Anti-patterns corrected:**

- Monolithic Makefiles exceeding 150 lines
- Mixed concerns (variables, rules, targets) in a single large file
- Inline `ifeq`/`else` blocks for environment switching instead of separate config files
- Duplicated configuration across multiple project Makefiles

---

## makefile-debugging-optimization

**Version:** 1.0.0

**Trigger description:** Activates when the user asks to debug a Makefile, diagnose slow builds, fix unexpected rebuild behavior, understand why Make rebuilds everything, figure out why a target does not rebuild, check variable values, use `make -n`/`-d`/`-p`/`--trace` flags, or optimize build performance.

**Topics covered:**

- `make -n` (dry-run): show planned commands without execution
- `make -d` (debug): show Make's decision-making process
- `make -p` (print database): show all rules, variables, and their values
- `make --trace` (GNU Make 4.0+): show rule executions in real-time
- `make --question` (exit code 0=up-to-date, 1=outdated, 2=error)
- Systematic debugging methodology: understand, identify, find root cause, fix
- Baseline measurement with `time make clean && time make all`
- Incremental build measurement with `touch file && time make all`
- Parallel builds with `MAKEFLAGS += -j$(shell nproc)`
- Compiler caching with `CC := ccache gcc`
- Immediate expansion (`:=`) vs recursive expansion (`=`) performance
- Conditional dependency inclusion to skip `.d` parsing during `make clean`

**Key patterns taught:**

| Pattern | Purpose |
|---------|---------|
| `make -n target` | Preview what would execute |
| `make -d target 2>&1 \| grep 'is newer'` | Find what triggers rebuilds |
| `make -p \| grep "^VAR ="` | Inspect variable values |
| `make --trace all` | Trace rule execution order |
| `MAKEFLAGS += -j$(shell nproc)` | Default parallel builds |
| `CC := ccache gcc` | Compiler output caching |

**Anti-patterns corrected:**

- Jumping to optimization suggestions without diagnosis
- Guessing at the cause of slow builds (shotgun approach)
- Using recursive expansion (`=`) for expensive `$(wildcard)` or `$(shell)` calls
- Parsing `.d` dependency files during `make clean`
- Missing parallelism due to undeclared dependencies or shell loops

## See Also

- [Getting Started](../../tutorials/getting-started/) -- tutorial walkthrough of the five skills
- [Architecture](../../explanation/architecture/) -- why skills are organized this way
- [Skill Progression](../../explanation/skill-progression/) -- the learning path through the skills
