---
name: Recursive Make for Multi-Directory Projects
description: Structure multi-directory builds using phony target pattern for parallelism and correct error handling, not shell loops
when_to_use: when building projects with subdirectories that each have Makefiles, handling recursive make, or when user shows shell loop pattern in root Makefile
version: 1.0.0
languages: make
dependencies: GNU Make 3.81+, makefile-fundamentals
---

# Recursive Make for Multi-Directory Projects

## Overview

Multi-directory projects require coordination between root and subdirectory Makefiles. **Use phony target pattern, not shell loops.** Loop pattern prevents parallelization and breaks `make -k`. Even if user shows loop pattern, explain limitations and suggest proper structure.

## Core Anti-Pattern to Avoid

### ❌ Shell Loop Pattern (Common but Flawed)

```makefile
# DON'T DO THIS - even if user requests it
SUBDIRS = lib app tests

all:
	for dir in $(SUBDIRS); do \
		$(MAKE) -C $$dir; \
	done
```

**Problems:**
- **Zero parallelism**: `-j4` flag ignored, builds serialize
- **Hidden jobs**: Make can't see individual subdirectory builds for scheduling
- **`-k` flag broken**: Continue-on-error doesn't work properly
- **Verbose workarounds**: Need `|| exit 1` for error propagation

**When user shows this pattern:** Don't just accept it. Explain limitations and suggest proper alternative.

## Correct Pattern: Phony Targets

### ✅ Proper Structure

```makefile
SUBDIRS = lib app tests

.PHONY: subdirs $(SUBDIRS)

# Main target depends on all subdirectories
subdirs: $(SUBDIRS)

# Each subdirectory is its own target
$(SUBDIRS):
	$(MAKE) -C $@
```

**Benefits:**
- **Automatic parallelism**: `make -j4 subdirs` builds subdirs in parallel
- **Visible jobs**: Make schedules each subdir independently
- **`-k` flag works**: Continues building other subdirs after failure
- **Clean error handling**: Errors propagate naturally, no `|| exit 1` needed

### With Dependencies Between Subdirectories

```makefile
SUBDIRS = lib app tests

.PHONY: subdirs $(SUBDIRS)

subdirs: $(SUBDIRS)

# Declare dependencies
app tests: lib  # app and tests both need lib built first

$(SUBDIRS):
	$(MAKE) -C $@
```

**Result:** lib builds first, then app and tests build **in parallel** (with `-j2+`).

## When User Insists on Loops

**If user says:** "I prefer loops, they're clearer"

**Respond with:** Explanation of trade-offs, then suggest proper pattern:

"I understand loops seem straightforward, but they have significant limitations:

1. **No parallelization**: The entire loop is one shell command. `make -j4` can't parallelize across subdirectories because Make doesn't see them as separate jobs.

2. **Error handling complexity**: Need `|| exit 1` for each command, and `-k` flag won't work as expected.

3. **Scalability**: As your project grows, serial builds become painfully slow.

The phony target pattern solves all these issues and is the standard approach in professional projects. Here's the equivalent structure that enables parallel builds..."

Then show the correct pattern above.

## Quick Reference

| Approach | Parallelizes with -j? | -k flag works? | Complexity |
|----------|----------------------|----------------|------------|
| Shell loop | ❌ No | ❌ No | High (needs || exit) |
| Phony targets | ✅ Yes | ✅ Yes | Low (clean) |

**Always suggest phony target pattern**, even if user shows loop preference.

## Variable Export

### Passing Variables to Subdirectories

**Option 1: Export directive** (recommended for few variables)

```makefile
VERSION = 1.2.3
CFLAGS = -O2

export VERSION CFLAGS

subdirs: $(SUBDIRS)
	# Variables automatically available in sub-makes
```

**Option 2: Command-line passing**

```makefile
$(SUBDIRS):
	$(MAKE) -C $@ VERSION=$(VERSION) CFLAGS="$(CFLAGS)"
```

**Option 3: Export all** (use sparingly)

```makefile
export  # Exports all variables

# Or use special target
.EXPORT_ALL_VARIABLES:
```

### Unexport Specific Variables

```makefile
export VERSION
unexport TEMP_VAR  # Don't pass this down
```

## Using $(MAKE) Variable

**Always use `$(MAKE)` for recursive invocations**, not hardcoded `make`:

```makefile
# ✅ CORRECT
$(SUBDIRS):
	$(MAKE) -C $@

# ❌ WRONG
$(SUBDIRS):
	make -C $@
```

**Why:** `$(MAKE)` ensures `-t`, `-n`, and `-q` flags work correctly in recursive contexts.

## Common Scenarios

### Scenario 1: Simple Multi-Directory Build

```makefile
SUBDIRS = frontend backend database

.PHONY: all clean $(SUBDIRS)

all: $(SUBDIRS)

clean:
	for dir in $(SUBDIRS); do \
		$(MAKE) -C $$dir clean; \
	done

$(SUBDIRS):
	$(MAKE) -C $@
```

**Note:** Clean often uses loop (no parallelism needed), build uses phony pattern (parallelism desired).

### Scenario 2: Complex Dependencies

```makefile
SUBDIRS = core utils app tests docs

.PHONY: all $(SUBDIRS)

all: $(SUBDIRS)

# Dependency chain
utils: core
app: core utils
tests: core utils app
docs: app

$(SUBDIRS):
	$(MAKE) -C $@
```

**Result:** Maximum parallelism while respecting build order.

### Scenario 3: Go Multi-Binary Project

```makefile
BINARIES = cmd/server cmd/cli cmd/worker

.PHONY: build $(BINARIES)

build: $(BINARIES)

$(BINARIES):
	go build -o bin/$(notdir $@) ./$@

clean:
	rm -rf bin/
```

**Note:** Even without subdirectory Makefiles, phony pattern enables parallelism.

## Error Handling

### Phony Pattern (Automatic)

```makefile
$(SUBDIRS):
	$(MAKE) -C $@
# Errors stop build automatically
# No special handling needed
```

### Continue on Error

```bash
make -k subdirs  # Builds all subdirs even if some fail
```

Works correctly with phony pattern, **doesn't work properly with loops**.

## Common Mistakes

| Mistake | Fix | Why |
|---------|-----|-----|
| Using loop for builds | Use phony target pattern | Enables parallelization |
| Using `make` instead of `$(MAKE)` | Always use `$(MAKE)` | Flags pass through correctly |
| Not declaring subdirs as .PHONY | Add to .PHONY | Prevents conflicts |
| Forgetting dependencies | Add prerequisite lines | Controls build order |
| Using loop when phony works | Only loop for targets like clean | Phony enables -j flag |

## Proactive Guidance

**When creating multi-directory builds:**
- Always suggest phony target pattern first
- Explain parallelization benefits
- Show dependency declaration

**When user shows loop pattern:**
- Acknowledge their approach
- Explain limitations (parallelization, -k flag)
- Suggest phony target alternative
- Show concrete performance difference

**When debugging "make -j doesn't work":**
- First check: are they using loops?
- Explain why loops serialize
- Provide phony target solution

## Red Flags - Review for Anti-Patterns

- [ ] Shell loop for building subdirectories?
- [ ] Using `make` instead of `$(MAKE)`?
- [ ] Subdirectories not in `.PHONY`?
- [ ] Missing dependency declarations?
- [ ] User preference for loop not challenged?

**If any red flags present:** Explain phony target pattern as proper solution, even if user didn't ask for it.

## Real-World Performance Impact

**Example:** Project with 8 subdirectories, each taking 30 seconds to build.

**Loop pattern:**
```bash
make all  # Takes 4 minutes (8 × 30s, serial)
```

**Phony pattern:**
```bash
make -j8 all  # Takes ~30 seconds (parallel)
```

**8x speedup** by using correct pattern. Mention this when explaining to users.

## The Bottom Line

**Shell loops prevent parallelization.** Even if user requests loops, explain the limitation and suggest phony target pattern. The performance difference matters for any non-trivial project, and the pattern is cleaner anyway.

Don't accept loop patterns without explanation - that's leaving 8x performance on the table.
