---
name: Advanced Makefile Features
description: >-
  This skill should be used when the user asks about Makefile pattern rules,
  automatic variables ($@, $<, $^, $?, $*), wildcard functions, patsubst,
  reducing Makefile duplication, compiling multiple C files with one rule, static
  pattern rules, conditional compilation with ifeq, or target-specific
  variables. Also applies when the user shows repetitive explicit rules that
  could be replaced with pattern rules, or asks about DRY Makefile practices.
  For basic Makefile creation and the ## help pattern, see
  makefile-fundamentals.
version: 1.0.0
---

# Advanced Makefile Features

## Overview

**Default to pattern rules and automatic variables**, not explicit repetitive rules. Pattern rules eliminate duplication, reduce errors, and scale to any project size. Even if user doesn't request them, suggest pattern rules proactively.

## Core Anti-Pattern to Avoid

### ❌ Repetitive Explicit Rules (Common but Unmaintainable)

```makefile
# DON'T DO THIS - even if user prefers it
main.o: main.c
	gcc -c main.c

utils.o: utils.c
	gcc -c utils.c

parser.o: parser.c
	gcc -c parser.c
```

**Problems:**
- **Violates DRY**: Same recipe repeated for every file
- **Doesn't scale**: 100 files = 100 rules to write and maintain
- **Error-prone**: Changing compile flags requires editing N places
- **More code**: Verbosity increases linearly with files
- **Harder to read**: Must scan N rules to understand pattern

**When user shows this pattern:** Don't just implement it. Explain limitations and suggest pattern rule alternative.

## Correct Pattern: Pattern Rules with Automatic Variables

### ✅ The DRY Way

```makefile
# Use this pattern by default
%.o: %.c
	$(CC) $(CFLAGS) -c -o $@ $<
```

**Benefits:**
- **DRY principle**: One rule handles infinite files
- **Scales perfectly**: Works for 1 file or 1000 files
- **Maintainable**: Change compile command in ONE place
- **Less code**: 3 lines instead of 3N lines
- **Clearer intent**: Shows the pattern, not repetitive noise

### Pattern Rule Syntax

```makefile
%.o: %.c
	$(CC) $(CFLAGS) -c -o $@ $<
```

**Components:**
- `%`: Matches any nonempty substring (the "stem")
- `%.o`: Target pattern (any file ending in .o)
- `%.c`: Prerequisite pattern (corresponding .c file)
- `$@`: Automatic variable = target name
- `$<`: Automatic variable = first prerequisite name

### Complete Example with Pattern Rules

```makefile
# Compiler and flags
CC := gcc
CFLAGS := -Wall -Wextra -std=c99

# Source files
SRCS := main.c utils.c parser.c
OBJS := $(SRCS:.c=.o)

# Target executable
TARGET := myapp

# Default target
all: $(TARGET)

# Link
$(TARGET): $(OBJS)
	$(CC) -o $@ $^

# Pattern rule for all .c → .o compilations
%.o: %.c
	$(CC) $(CFLAGS) -c -o $@ $<

clean:
	rm -f $(OBJS) $(TARGET)

.PHONY: all clean
```

**Result:** Works for 3 files or 300 files with zero changes to pattern rule.

## Automatic Variables - Essential Knowledge

### Core Automatic Variables

Always use these in recipes for DRY code:

```makefile
$@  # Target name
$<  # First prerequisite
$^  # All prerequisites (space-separated)
$?  # Prerequisites newer than target
$*  # Stem (matched by % in pattern rule)
```

### Usage Examples

```makefile
# Compilation: target = .o, first prereq = .c
%.o: %.c headers.h
	$(CC) $(CFLAGS) -c -o $@ $<
	# $@ = foo.o, $< = foo.c

# Linking: target = exe, all prereqs = .o files
myapp: $(OBJS)
	$(CC) -o $@ $^
	# $@ = myapp, $^ = main.o utils.o parser.o

# Archive: only changed files
lib.a: $(OBJS)
	ar r $@ $?
	# $? = only objects newer than lib.a
```

### Why Automatic Variables Matter

**Without automatic variables:**
```makefile
%.o: %.c
	gcc -c -o %.o %.c  # WRONG - % doesn't work here!
```

**With automatic variables:**
```makefile
%.o: %.c
	$(CC) $(CFLAGS) -c -o $@ $<  # CORRECT
```

Automatic variables provide the actual matched names at execution time.

## Functions for Flexibility

### wildcard - Automatic File Discovery

```makefile
# Find all .c files in current directory
SRCS := $(wildcard *.c)

# Find all .c files recursively
SRCS := $(wildcard *.c src/*.c lib/*.c)
```

**Use when:** You want Makefile to adapt to new/removed files automatically.

### patsubst - Pattern Substitution

```makefile
# Convert .c list to .o list
OBJS := $(patsubst %.c,%.o,$(SRCS))

# Or shorter syntax:
OBJS := $(SRCS:.c=.o)
```

### foreach - Iteration

```makefile
# Generate multiple targets
APPS := app1 app2 app3

all: $(APPS)

$(foreach app,$(APPS),$(eval $(app): $(app).c))
```

## When User Prefers Explicit Rules

When user prefers explicit rules, acknowledge their perspective, then explain the scalability and maintenance limitations (N files = N rules to edit for any flag change). Show both approaches side-by-side with the pattern rule version as the recommended default.

## "Keep It Simple" Requests

Pattern rules are simpler — less code, fewer places to edit, impossible to have copy-paste inconsistencies. When user equates "simple" with "explicit rules for each file," reframe: more lines of identical code is more complexity, not less. Show a concrete line-count comparison.

## Common Pattern Rule Scenarios

For detailed scenarios covering C compilation, direct-to-executable builds, multiple extensions, auto-discovery with wildcard, and static pattern rules, see `references/pattern-rule-scenarios.md`.

## Conditional Compilation

### Mode-Based Flags

```makefile
DEBUG ?= 0

ifeq ($(DEBUG),1)
    CFLAGS += -g -O0 -DDEBUG
else
    CFLAGS += -O2 -DNDEBUG
endif

%.o: %.c
	$(CC) $(CFLAGS) -c -o $@ $<
```

**Usage:**
```bash
make              # Release mode
make DEBUG=1      # Debug mode
```

### Target-Specific Variables

```makefile
%.o: %.c
	$(CC) $(CFLAGS) -c -o $@ $<

# Special flags for specific files
parser.o: CFLAGS += -Wno-unused-function
```

## Common Mistakes

| Mistake | Fix | Why |
|---------|-----|-----|
| Explicit rules for similar targets | Use pattern rule | DRY, scales better |
| Hardcoded filenames in recipes | Use automatic variables ($@, $<) | Reusable, error-free |
| Not using wildcard for file lists | Use $(wildcard *.c) | Adapts to new files |
| Repeating substitution patterns | Use patsubst or $(VAR:.c=.o) | Cleaner, maintainable |
| Writing make instead of $(MAKE) | Always use $(MAKE) | Flags pass through |

## Proactive Guidance - Always Apply

**When creating Makefiles:**
- ✅ Default to pattern rules
- ✅ Use automatic variables ($@, $<, $^)
- ✅ Use wildcard for automatic discovery when appropriate
- ✅ Show explicit rules as comparison only

**When reviewing Makefiles:**
- 🔍 Look for repetitive explicit rules
- 💡 Suggest pattern rule refactoring
- 📊 Show before/after line counts
- 🎯 Explain scalability benefits

**When user shows explicit rules:**
- 🤔 Acknowledge their approach
- ⚠️ Explain maintainability concerns
- ✅ Show pattern rule alternative
- 📈 Provide concrete scaling example

**When facing time pressure:**
- ⚡ Pattern rules are FASTER to write (3 lines vs 30)
- 🎯 Don't let "quick" become "verbose"
- ✅ Pattern rules are the shortcut, not the explicit rules

## Red Flags - Review for Anti-Patterns

When reviewing any Makefile:

- [ ] Multiple rules with same recipe pattern?
- [ ] Hardcoded filenames in recipes instead of $@ or $<?
- [ ] File lists hardcoded instead of using wildcard?
- [ ] Copy-paste between similar rules?
- [ ] User preference for explicit rules not challenged?

**If any red flags present:** Explain pattern rule benefits, even if user didn't ask for it.

## Real-World Impact

### Example: Medium-Size C Project (50 files)

**Explicit rules approach:**
```makefile
# 150 lines (50 files × 3 lines each)
main.o: main.c
	gcc -Wall -c main.c

utils.o: utils.c
	gcc -Wall -c utils.c

# ... 48 more rules
```

**Pattern rules approach:**
```makefile
# 3 lines (regardless of file count)
%.o: %.c
	$(CC) $(CFLAGS) -c -o $@ $<
```

**Comparison:**
- Lines of code: 150 → 3 (50x reduction)
- Places to edit for flag changes: 50 → 1 (50x easier maintenance)
- New file overhead: +3 lines → +0 lines (automatic)
- Potential copy-paste errors: 50 opportunities → 0 opportunities

### Concrete Scenario: Adding Debug Flags

**With explicit rules:**
- Must edit 50 compilation lines
- Risk of missing some or making typos
- 5-10 minutes of careful editing

**With pattern rules:**
- Change $(CFLAGS) definition once
- Impossible to make inconsistent
- 10 seconds

## The Bottom Line

**Pattern rules and automatic variables are not "advanced features" - they're essential best practices.**

Don't default to explicit rules because they seem "simpler" or the user is "in a hurry". Pattern rules:
- Have LESS code (simpler!)
- Are FASTER to write (quicker!)
- Scale perfectly (maintainable!)
- Eliminate errors (reliable!)

**Always suggest pattern rules first.** Show explicit rules only as a "what not to do" comparison.

## Quick Reference Card

### Most Common Patterns

```makefile
# C compilation
%.o: %.c
	$(CC) $(CFLAGS) -c -o $@ $<

# C to executable
%: %.c
	$(CC) $(CFLAGS) -o $@ $<

# Auto-discovery
SRCS := $(wildcard *.c)
OBJS := $(SRCS:.c=.o)

# Linking
$(TARGET): $(OBJS)
	$(CC) -o $@ $^
```

### Essential Automatic Variables

```makefile
$@  # Target file name
$<  # First prerequisite
$^  # All prerequisites
$?  # Changed prerequisites
$*  # Pattern stem
```

### When to Use Pattern Rules

✅ Always - for any repetitive pattern
✅ Compilation rules (.c → .o, .cpp → .o)
✅ Simple build rules (.c → executable)
✅ Any N-to-N transformation
✅ As default recommendation

❌ Never avoid them due to "complexity" - they're simpler!
