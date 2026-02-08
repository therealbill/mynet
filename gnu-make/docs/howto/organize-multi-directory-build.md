---
title: "Organize a Multi-Directory Build"
description: "Set up a multi-directory Makefile build using phony targets for correct parallelization and dependency management"
weight: 2
---

# Organize a Multi-Directory Build

This guide shows how to set up a root Makefile that coordinates builds across multiple subdirectories. The approach uses phony targets instead of shell loops, which enables parallel builds with `make -j` and correct error handling with `make -k`.

## Identify the Problem

A common but flawed pattern uses a shell for-loop to iterate over subdirectories:

```makefile
# Anti-pattern: prevents parallelization
SUBDIRS = lib app tests

all:
	for dir in $(SUBDIRS); do \
		$(MAKE) -C $$dir; \
	done
```

This loop runs each subdirectory build sequentially. The `make -j4` flag has no effect because Make sees a single shell command, not separate targets. You also lose `make -k` (continue on error) behavior and need manual `|| exit 1` for error propagation.

## Set Up Phony Targets

Replace the loop with individual phony targets. Each subdirectory becomes its own Make target.

```makefile
SUBDIRS = lib app tests

.PHONY: all $(SUBDIRS)

all: $(SUBDIRS)

$(SUBDIRS):
	$(MAKE) -C $@
```

Each subdirectory is now a separate target that Make can schedule independently. Running `make -j4 all` will build up to four subdirectories in parallel.

Always use `$(MAKE)` instead of hardcoded `make` for recursive invocations. The `$(MAKE)` variable ensures that flags such as `-n` (dry-run), `-t` (touch), and `-q` (question) propagate correctly to sub-makes.

## Declare Dependencies Between Subdirectories

If one subdirectory depends on another (for example, `app` requires `lib` to be built first), declare it as a prerequisite.

```makefile
SUBDIRS = lib app tests

.PHONY: all $(SUBDIRS)

all: $(SUBDIRS)

# app and tests both depend on lib
app: lib
tests: lib

$(SUBDIRS):
	$(MAKE) -C $@
```

With this declaration, `make -j4 all` will:

1. Build `lib` first (both `app` and `tests` depend on it)
2. Build `app` and `tests` in parallel (they are independent of each other)

Make resolves the dependency graph automatically. You do not need to manage ordering manually.

## Pass Variables to Subdirectories

Subdirectory Makefiles often need access to root-level variables. There are two approaches.

### Export Specific Variables

Use the `export` directive for variables that sub-makes should inherit:

```makefile
VERSION = 1.2.3
CFLAGS = -Wall -O2

export VERSION CFLAGS

$(SUBDIRS):
	$(MAKE) -C $@
```

### Pass on the Command Line

For variables that should override sub-make values, pass them directly:

```makefile
$(SUBDIRS):
	$(MAKE) -C $@ VERSION=$(VERSION)
```

Command-line variables override definitions inside the sub-make's Makefile, while exported variables only set defaults.

## Add Clean and Other Multi-Directory Targets

For targets like `clean` that need to run in every subdirectory, a loop is acceptable because parallelism is not needed for cleanup. However, you can also use the phony pattern.

```makefile
.PHONY: clean

clean:
	for dir in $(SUBDIRS); do \
		$(MAKE) -C $$dir clean; \
	done
```

Or, using the phony target pattern for consistency:

```makefile
CLEAN_SUBDIRS = $(addprefix clean-,$(SUBDIRS))

.PHONY: clean $(CLEAN_SUBDIRS)

clean: $(CLEAN_SUBDIRS)

$(CLEAN_SUBDIRS):
	$(MAKE) -C $(patsubst clean-%,%,$@) clean
```

The second approach allows `make -j clean` to run cleanup in parallel across subdirectories, which can be faster for large projects.

## Enable Parallel Builds

Once phony targets and dependencies are declared, parallel builds work automatically:

```bash
# Build with 4 parallel jobs
make -j4 all

# Build with all available cores
make -j$(nproc) all
```

For a project with 8 independent subdirectories each taking 30 seconds, sequential builds take 4 minutes while `make -j8` completes in approximately 30 seconds.

To set a default parallelism level in the Makefile:

```makefile
MAKEFLAGS += -j$(shell nproc)
```

## Complete Example

Here is a complete root Makefile for a project with three subdirectories and declared dependencies.

```makefile
SUBDIRS = lib app tests

.PHONY: all clean help $(SUBDIRS)

all: $(SUBDIRS) ## Build all subdirectories

# Dependency declarations
app: lib
tests: lib app

$(SUBDIRS):
	$(MAKE) -C $@

clean: ## Clean all subdirectories
	for dir in $(SUBDIRS); do \
		$(MAKE) -C $$dir clean; \
	done

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
```

Each subdirectory (`lib/`, `app/`, `tests/`) should contain its own Makefile with its own targets, pattern rules, and variables. The root Makefile only orchestrates the build order and parallelism.

## Verify Correct Behavior

Test that parallelism and dependencies work as expected:

```bash
# Dry-run to see the build plan
make -n all

# Build with verbose output to verify ordering
make -j4 all

# Verify that lib builds before app and tests
make -j4 all 2>&1 | grep "Entering directory"
```

The "Entering directory" messages should show `lib` before `app` and `tests`, with `app` and `tests` interleaved (indicating parallel execution).
