---
name: Makefile Debugging and Optimization
description: >-
  This skill should be used when the user asks to debug a Makefile, diagnose
  slow builds, fix unexpected rebuild behavior, understand why make rebuilds
  everything, figure out why a target does not rebuild, check Makefile variable
  values, use make -n/-d/-p/--trace flags, optimize build performance, or
  troubleshoot any Makefile issue. Also applies when the user reports "build is
  slow", "everything rebuilds", or "target never rebuilds".
version: 1.0.0
---

# Makefile Debugging and Optimization

## Overview

**Teach systematic debugging first** using Make's built-in diagnostic flags. When users report issues (slow builds, incorrect behavior, mysterious rebuilds), start with debugging methodology, not immediate solutions. Make's debugging flags reveal the root cause - use them proactively.

## Core Anti-Pattern to Avoid

### ❌ Jump to Solutions Without Diagnosis

```
User: "My build is slow"
Agent: "Here are 8 optimizations to try:
  1. Use ccache
  2. Enable parallel builds
  3. Optimize compiler flags
  ..."
```

**Problems:**
- **Guessing**: No diagnosis of actual problem
- **Shotgun approach**: Try many things, hope one works
- **Misses root cause**: May not address actual issue
- **Not educational**: User doesn't learn debugging approach
- **Wastes time**: Implementing irrelevant fixes

**Even if suggestions are good, this misses teaching opportunity.**

## Correct Pattern: Systematic Debugging

### ✅ Debug First, Then Fix

```
User: "My build is slow"
Agent: "Let's debug systematically to find the root cause.

First, let's see what Make is actually doing:

1. Dry-run to see planned actions:
   make -n all

2. Check what's being rebuilt and why:
   make -d all 2>&1 | grep 'File.*is newer'

3. Time different scenarios:
   make clean && time make all    # Full build
   touch file.c && time make all  # Incremental

This will show us exactly what's slow. Then we can fix the specific issue.

[After seeing output:]
The debug output shows Make is rebuilding everything because...
Here's the specific fix..."
```

**Benefits:**
- **Accurate diagnosis**: See actual problem
- **Educational**: User learns debugging approach
- **Targeted fix**: Address root cause
- **Reproducible**: User can debug future issues
- **Professional**: This is how experts debug

## Make's Essential Debugging Flags

### Always Suggest These First

#### `make -n` (Dry-Run)

**Purpose:** Show what would be executed without running commands

```bash
make -n all

# See what would run for specific target
make -n install

# Check if Make thinks rebuild is needed
make all && make -n all  # Should show "nothing to be done"
```

**When to use:**
- "What will this Makefile do?"
- "Why is Make rebuilding?"
- Verify Makefile logic before running

#### `make -d` (Debug Mode)

**Purpose:** Show Make's decision-making process

```bash
make -d all 2>&1 | less

# Filter for specific insights
make -d all 2>&1 | grep "Considering\|File.*is newer"

# See why target rebuilt
make -d mytarget 2>&1 | grep mytarget
```

**When to use:**
- "Why did Make rebuild this?"
- "How is Make resolving dependencies?"
- "Why is Make not rebuilding when it should?"

#### `make -p` (Print Database)

**Purpose:** Show all rules, variables, and values

```bash
make -p | less

# See specific variable value
make -p | grep "^CC ="

# See all pattern rules
make -p | grep "^%"

# See automatic variables
make -p | grep "automatic"
```

**When to use:**
- "What's the value of this variable?"
- "What implicit rules exist?"
- "Why is this variable not what I expect?"

#### `make --trace` (GNU Make 4.0+)

**Purpose:** Show rule executions in real-time

```bash
make --trace all

# See what triggered rebuild
touch file.c
make --trace all
```

**When to use:**
- "What's the execution order?"
- "What triggered this rebuild?"
- More readable than `-d` for basic tracing

### Quick Reference Card

| Flag | Purpose | Common Use |
|------|---------|------------|
| `-n` | Dry-run (what would run) | "Will this work?" |
| `-d` | Debug (why decisions made) | "Why rebuilt?" |
| `-p` | Database (all rules/vars) | "What's the value?" |
| `--trace` | Trace execution | "What ran and why?" |

**Always mention these when user has Makefile problems.**

## Systematic Debugging Approach

### Step 1: Understand What's Happening

```bash
# What commands would run?
make -n target

# Is Make seeing what you expect?
make -p | grep VARIABLE_NAME

# What's the current state?
make --question target; echo $?  # 0=up-to-date, 1=outdated, 2=error
```

### Step 2: Identify the Problem

```bash
# For "it's slow" issues:
time make clean && time make all  # Baseline
touch one_file.c && time make all # Incremental
# Should be fast! If not, debug why

# For "it rebuilds everything" issues:
make all
make -n all  # Should say "nothing to be done"
# If not, something's wrong with dependencies

# For "it doesn't rebuild" issues:
touch dependency.h
make -d target 2>&1 | grep dependency.h
# Should show it considers the file
```

### Step 3: Find Root Cause

```bash
# Why is this file being rebuilt?
make -d target.o 2>&1 | grep -A5 "target.o"

# What triggered this?
make --trace all

# What does Make think this variable is?
make -p | grep "^VARNAME\s*="
```

### Step 4: Implement Targeted Fix

**Now** suggest specific solutions based on diagnosis, not guesses.

## Common Issues and How to Debug Them

### Issue: "Build is Slow"

**DON'T say:** "Use ccache, enable -j, optimize flags..."

**DO say:**

"Let's measure where time is spent:

```bash
# Baseline timing
time make clean && time make all

# Incremental build (should be fast)
touch src/main.c && time make all

# Is it compilation or linking?
make clean
time make VERBOSE=1 all  # Watch which steps are slow
```

What do these timings show? That tells us where to optimize."

**After diagnosis**, suggest targeted fixes.

### Issue: "Everything Rebuilds"

**DON'T say:** "Check your dependencies, use .d files..."

**DO say:**

"Let's see why Make thinks everything is out of date:

```bash
# Build once
make all

# Should be no-op
make -n all

# If it wants to rebuild, debug why:
make -d all 2>&1 | grep 'File.*is newer than'
```

This shows exactly which file timestamps are causing rebuilds."

### Issue: "Target Doesn't Rebuild When It Should"

**DON'T say:** "You need to add dependencies..."

**DO say:**

"Let's see what dependencies Make knows about:

```bash
# Print rule for target
make -p | grep -A10 "^target:"

# Check if Make sees file changes
touch dependency.h
make -d target 2>&1 | grep dependency.h

# Verify .d files if using auto-dependencies
cat build/target.d
```

This shows what Make knows. We can see what's missing."

## When to Suggest Optimizations

### After Understanding the Problem

**Timing showed slow compilation:**
→ Suggest `ccache`, `-j` flag, optimize flags

**Timing showed slow linking:**
→ Suggest split debug info, reduce -O level for dev

**Debug showed all files rebuild:**
→ Fix timestamp issue, check .d file generation

**Debug showed no parallelism:**
→ Suggest `.NOTPARALLEL` removal, add `-j` flag

**Debug showed expensive wildcard:**
→ Cache results, use immediate expansion `:=`

## Proactive Debugging Education

### When User Reports Issue

**Template response structure:**

1. **Acknowledge issue**
2. **Suggest debugging approach** (flags to use)
3. **Wait for or request output**
4. **Analyze output**
5. **Provide targeted fix**

**Example:**

"That rebuild behavior sounds incorrect. Let's debug it systematically:

```bash
# First, let's see what Make is actually doing:
make all
make -n all  # Should say 'nothing to be done'

# If it wants to rebuild, let's see why:
make -d all 2>&1 | grep 'File.*is newer'
```

Can you run these and share the output? That will show us exactly what's causing unnecessary rebuilds.

[After seeing output:]
Ah, I see the issue. Make thinks `config.h` is newer than the targets because... Here's how to fix it..."

### When Creating Makefiles

**Proactively include debugging support:**

```makefile
# Debug and optimization targets
.PHONY: debug-make help-debug

debug-make:
	@echo "==> Dry-run (what would be executed):"
	@$(MAKE) -n all
	@echo ""
	@echo "==> Variable values:"
	@$(MAKE) -p | grep "^CC\|^CFLAGS\|^LDFLAGS"

help-debug:
	@echo "Debugging targets:"
	@echo "  make -n target     - Show what would run"
	@echo "  make -d target     - Show decision-making"
	@echo "  make -p | grep VAR - Show variable value"
	@echo "  make --trace       - Trace execution"
```

**Mention in help target:**

```makefile
help: ## Show this help
	@echo "Debugging:"
	@echo "  make -n all        - Dry-run to see planned commands"
	@echo "  make -d all        - Debug why Make is rebuilding"
	@echo "  make debug-make    - Show Makefile diagnostics"
```

## Common Debugging Scenarios

### Scenario 1: Mysterious Rebuilds

**Steps:**

1. `make all && make -n all` (should be no-op)
2. If rebuilding: `make -d all 2>&1 | grep "is newer"`
3. Check timestamps: `ls -lt file1 file2 target`
4. Find root cause (clock skew, generated files, etc.)

### Scenario 2: Target Never Rebuilds

**Steps:**

1. `make -p | grep -A5 "^target:"` (check rule exists)
2. `touch dependency && make -d target` (see if detected)
3. Check if target is phony accidentally
4. Verify prerequisite paths are correct

### Scenario 3: Variable Has Wrong Value

**Steps:**

1. `make -p | grep "^VAR ="`
2. Check if recursively vs simply expanded (= vs :=)
3. `make -d` to see where variable is set
4. Check include order (last definition wins)

### Scenario 4: Slow Incremental Builds

**Steps:**

1. `time make clean && time make all` (baseline)
2. `touch file.c && time make all` (should be <5s)
3. If slow: `make --trace all` (see what runs)
4. Profile compilation vs linking vs wildcards

## Optimization Principles

### Measure Before Optimizing

```bash
# Don't guess - measure
time make clean && time make all

# Measure incremental
touch src/main.c
time make all

# Profile with timestamps
make all
stat target dependencies
```

### Common Optimizations (After Diagnosis)

**1. Parallel Builds** (if compilation is slow)
```makefile
MAKEFLAGS += -j$(shell nproc)
```

**2. Compiler Cache** (if recompiling same files)
```makefile
CC := ccache gcc
```

**3. Optimize Dev Flags** (if dev build is slow)
```makefile
# Development
CFLAGS := -g -Og  # Not -O0

# Conditional sanitizers
ifdef SANITIZE
  CFLAGS += -fsanitize=address
endif
```

**4. Lazy Wildcard** (if startup is slow)
```makefile
# Faster: immediate expansion
SRCS := $(wildcard src/*.c)

# Not: recursive expansion
SRCS = $(wildcard src/*.c)  # Re-evaluates every time
```

**5. Conditional Dependency Inclusion**
```makefile
# Don't parse .d files for clean
ifneq ($(MAKECMDGOALS),clean)
  -include $(DEPS)
endif
```

## The Bottom Line

**Don't guess - debug.** Make has powerful diagnostic flags that reveal exactly what's happening. Use them first:

1. **Problem reported** → Suggest debugging flags
2. **Run flags, see output** → Understand root cause
3. **Root cause identified** → Provide targeted fix
4. **User learns** → Can debug future issues themselves

Teaching debugging methodology is more valuable than providing quick fixes. Users who learn `make -n`, `make -d`, and `make -p` can solve their own problems.

**Always frame debugging as first step, not afterthought.**

## Quick Diagnostic Commands

**Keep these handy:**

```bash
# Is Make doing what I expect?
make -n target

# Why did Make rebuild this?
make -d target 2>&1 | grep target

# What's this variable's value?
make -p | grep "^VAR ="

# Time everything
time make clean && time make all

# Check incremental build
touch file.c && time make all  # Should be <5s

# Trace execution
make --trace all
```

Suggest these proactively when user has issues. Make debugging flags ARE the debugging methodology.
