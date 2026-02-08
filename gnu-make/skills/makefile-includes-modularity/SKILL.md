---
name: Makefile Includes and Modularity
description: >-
  This skill should be used when the user asks about splitting a Makefile into
  multiple files, using the include or -include directive, organizing a large
  Makefile into modules, creating environment-specific configurations
  (dev/prod/staging), sharing Makefile configuration across projects, or when a
  Makefile exceeds 150 lines. Also applies when the user mentions config.mk,
  rules.mk, modular Makefile structure, or asks how to reduce Makefile size and
  improve maintainability.
version: 1.0.0
---

# Makefile Includes and Modularity

## Overview

**Suggest modular structure for large projects** (>150 lines). Use `include` directive to split concerns into focused files, not maintain everything in one monolithic Makefile. Even if user prefers "one file", explain trade-offs and show modular alternative.

## Core Anti-Pattern to Avoid

### ❌ Monolithic Makefile (Common but Unmaintainable)

```makefile
# DON'T DO THIS - 500+ line single file
# Lines 1-100: Variables and configuration
VERSION = 1.0.0
CC = gcc
CFLAGS = -Wall -Wextra
PREFIX = /usr/local
... (80 more variable definitions)

# Lines 101-250: Core library rules
lib_sources = $(wildcard src/core/*.c)
lib_objects = $(lib_sources:.c=.o)
... (140 more lines of lib rules)

# Lines 251-400: API server rules
api_sources = $(wildcard src/api/*.c)
... (140 more lines of API rules)

# Lines 401-500: CLI tool rules, install, clean, etc.
... (100 more lines)
```

**Problems:**
- **Hard to navigate**: Finding config means scrolling through 500 lines
- **Merge conflicts**: Team members editing same large file
- **Cannot reuse**: Can't share common parts with other projects
- **Mixed concerns**: Variables, rules, targets all tangled together
- **Intimidating**: 500 lines discourages contributions

**When user shows this pattern:** Don't accept it. Suggest modular structure even if they prefer "one file for simplicity".

## Correct Pattern: Modular Organization

### ✅ Split Into Focused Modules

```makefile
# Makefile (main orchestrator - 50 lines)
include config.mk
include rules.mk
include targets.mk

.PHONY: all clean install

all: $(TARGETS)

clean:
	$(RM) $(OBJECTS) $(TARGETS)
```

```makefile
# config.mk (100 lines) - All configuration
VERSION := 1.0.0
CC := gcc
CFLAGS := -Wall -Wextra -std=c11
PREFIX := /usr/local

# Component-specific configs
CORE_CFLAGS := $(CFLAGS) -fPIC
API_CFLAGS := $(CFLAGS) -pthread
```

```makefile
# rules.mk (150 lines) - Pattern rules
%.o: %.c
	$(CC) $(CFLAGS) -c -o $@ $<

$(LIBDIR)/%.a: $(OBJECTS)
	$(AR) rcs $@ $^
```

```makefile
# targets.mk (200 lines) - Specific targets
TARGETS := libcore.a api-server cli-tool

libcore.a: $(CORE_OBJECTS)
api-server: $(API_OBJECTS) libcore.a
cli-tool: $(CLI_OBJECTS) libcore.a
```

**Benefits:**
- **Easy navigation**: Config in config.mk, rules in rules.mk
- **Less conflicts**: Team edits different files
- **Reusable**: Can `include` config.mk in other projects
- **Focused files**: Each ~100-150 lines, single responsibility
- **Approachable**: Smaller files encourage understanding and contributions

## The `include` Directive

### Basic Syntax

```makefile
include filename
include file1.mk file2.mk file3.mk
include *.mk
include $(CONFIG_FILE)
```

**Behavior:**
- Make suspends reading current file
- Reads included files in order
- Resumes reading original file
- Variables and rules from included files are available

### `-include` for Optional Files

```makefile
# Include if exists, silent if not
-include local-config.mk
-include dev.mk
-include $(wildcard *.d)  # Auto-dependencies
```

**Use when:** File might not exist (local overrides, auto-generated deps).

## When to Suggest Modular Structure

### Size Thresholds

| Lines | Recommendation |
|-------|----------------|
| < 50 | Monolithic okay |
| 50-150 | Consider modules if multiple concerns |
| 150-300 | **Suggest modules** |
| 300+ | **Strongly recommend modules** |

### Multiple Concerns Checklist

Even if < 150 lines, suggest modules when project has:
- [ ] Multiple components (lib + server + cli)
- [ ] Environment configs (dev, staging, prod)
- [ ] Different flag sets per component
- [ ] Shared config with other projects
- [ ] Team collaboration (reduce conflicts)

**If 2+ checkboxes:** Suggest modular structure.

## Standard Module Categories

### config.mk - Configuration Variables

```makefile
# Project metadata
PROJECT := myapp
VERSION := 1.0.0

# Tools
CC := gcc
AR := ar
INSTALL := install

# Flags
CFLAGS := -Wall -Wextra -std=c11
LDFLAGS := -lpthread
PREFIX := /usr/local

# Paths
SRCDIR := src
BUILDDIR := build
BINDIR := $(PREFIX)/bin
```

**Contains:** Variables only, no rules or targets.

### rules.mk - Pattern Rules

```makefile
# Compilation rules
%.o: %.c
	$(CC) $(CFLAGS) -c -o $@ $<

# Library creation
%.a: $(OBJECTS)
	$(AR) rcs $@ $^

# Linking
%: %.o
	$(CC) -o $@ $^ $(LDFLAGS)
```

**Contains:** Pattern rules only, generic transformations.

### targets.mk - Specific Targets

```makefile
# Concrete targets
TARGETS := libcore.a api-server cli-tool

libcore.a: $(CORE_OBJECTS)

api-server: $(API_OBJECTS) libcore.a
	$(CC) -o $@ $^ $(LDFLAGS)

cli-tool: $(CLI_OBJECTS) libcore.a
	$(CC) -o $@ $^ $(LDFLAGS)
```

**Contains:** Specific targets and their dependencies.

### install.mk - Installation Rules

```makefile
.PHONY: install uninstall

install: all
	$(INSTALL) -d $(BINDIR)
	$(INSTALL) -m 755 $(TARGETS) $(BINDIR)

uninstall:
	$(RM) $(addprefix $(BINDIR)/,$(TARGETS))
```

**Contains:** Installation, uninstallation, distribution.

## Environment-Specific Configurations

### Separate Config Files Pattern

```makefile
# Makefile - select configuration
ENV ?= dev
include config/$(ENV).mk
include rules.mk

all: $(TARGET)
```

```makefile
# config/dev.mk
CFLAGS := -g -O0 -DDEBUG
TARGET := myapp-dev
```

```makefile
# config/prod.mk
CFLAGS := -O3 -DNDEBUG
LDFLAGS := -s
TARGET := myapp
```

**Usage:**
```bash
make              # Uses dev (default)
make ENV=prod     # Uses production config
make ENV=staging  # Uses staging config
```

**Benefits:**
- Adding new environment = creating new file
- Clear separation of concerns
- Easy to compare configs (separate files vs if/else blocks)

### Alternative: Conditional with `-include`

```makefile
# Makefile - defaults
CFLAGS := -Wall -Wextra

# Override with environment-specific
-include dev.mk
-include prod.mk

all: $(TARGET)
```

```bash
# Build with specific config
make -f dev.mk    # or
touch dev.mk && make  # if dev.mk exists, it's included
```

## Shared Configuration Across Projects

For patterns on sharing configuration across multiple projects using `include ../common/common.mk`, see `references/shared-configuration.md`.

## When User Insists on "One File"

When the user prefers a single file, acknowledge the convenience for small projects, then explain the trade-offs at scale: navigating 500 lines vs opening the right focused file, merge conflicts when multiple people edit one file, and inability to share configuration across projects. Create both versions with the modular structure as the default, and let the user compare.

## Proactive Guidance - Size-Based Triggers

### Creating New Makefile

**If project description suggests >150 lines:**
- Suggest modular structure from the start
- Show example file breakdown
- Explain benefits before creating monolithic

**If project is clearly simple (< 50 lines):**
- Monolithic is fine
- Still mention: "For larger projects, we'd split into modules"

### Reviewing Existing Makefile

**If file is >150 lines:**
- Identify logical sections
- Suggest specific split: "Lines 1-50 → config.mk, 51-100 → rules.mk, etc."
- Show refactored example

**If file is 50-150 lines but has multiple concerns:**
- Mention modular option
- Explain when it becomes worthwhile

## Common Module Patterns

Four patterns (Three-File Foundation, Component-Based, Environment + Component, Shared Common) are documented in `references/module-patterns.md`. Choose based on project complexity.

## Red Flags - Review for Modularity

When reviewing any Makefile:

- [ ] File is >150 lines?
- [ ] Multiple components (lib, server, cli)?
- [ ] Multiple environment configs mixed with if/else?
- [ ] Repeated patterns that could be shared?
- [ ] Team collaboration (git history shows conflicts)?
- [ ] User says "I know it's messy" or "hard to find things"?

**If 2+ red flags present:** Strongly suggest modular structure.

## Benefits Summary

### Monolithic (< 50 lines)
✅ Simple, everything visible
✅ No include overhead
✅ Good for small, focused projects

### Modular (> 150 lines)
✅ Easier navigation (focused files)
✅ Better for teams (less conflicts)
✅ Reusable across projects
✅ Clearer organization
✅ Easier to maintain
✅ Simpler to understand (smaller chunks)

## Quick Reference

### When to Use `include`

**Use for:**
- ✅ Large projects (>150 lines)
- ✅ Multiple components
- ✅ Environment-specific configs
- ✅ Shared configuration across projects
- ✅ Team collaboration (reduce conflicts)
- ✅ Separating concerns (config vs rules vs targets)

**Don't need for:**
- ❌ Simple projects (< 50 lines)
- ❌ Single component, single config
- ❌ Personal projects with no reuse

### Standard File Names

```makefile
config.mk        # Variables, configuration
rules.mk         # Pattern rules
targets.mk       # Specific targets
install.mk       # Installation rules
test.mk          # Testing targets
dev.mk           # Development config
prod.mk          # Production config
common.mk        # Shared across projects
```

### Module Size Guidelines

- config.mk: 50-150 lines (variables)
- rules.mk: 50-200 lines (patterns)
- targets.mk: 100-300 lines (targets)
- Each module: Single responsibility

**If a module exceeds 300 lines:** Consider splitting further.

## The Bottom Line

**Monolithic Makefiles don't scale.** A 50-line file is fine. A 500-line file is unmaintainable.

Don't accept "one file is simpler" for large projects. Modular structure:
- REDUCES complexity (smaller focused files)
- IMPROVES navigation (clear file purposes)
- ENABLES reuse (share across projects)
- PREVENTS conflicts (team works on different files)

**Suggest modular structure proactively at 150+ lines.** Show concrete comparison. Explain trade-offs even if user has preference.

"One file" feels simpler until you're searching through 500 lines. Then modular becomes obviously better.
