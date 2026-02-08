---
title: "Debug a Slow Makefile"
description: "Diagnose and fix slow build times using Make's built-in debugging flags and systematic measurement"
weight: 1
---

# Debug a Slow Makefile

This guide walks through a systematic approach to diagnosing and fixing slow Makefile builds. The method uses Make's built-in diagnostic flags to identify the root cause before applying any optimization.

## Establish a Baseline

Before changing anything, measure the current state. Run a full clean build and an incremental build to understand where time is spent.

```bash
# Full build time
time make clean && time make all

# Incremental build time (single file change)
touch src/main.c && time make all
```

Record both numbers. A full build measures overall compilation speed. The incremental build reveals whether Make is doing unnecessary work. An incremental build that takes nearly as long as a full build indicates a dependency or rebuild problem.

## Inspect Planned Actions with Dry-Run

Use `make -n` to see what Make intends to execute without actually running anything.

```bash
make -n all
```

Review the output for red flags:

- Commands that compile files you did not change
- Redundant operations (compiling the same file twice)
- Missing parallelism (everything in sequence when it could be parallel)

After a successful build, run `make -n all` again. The output should say "nothing to be done." If Make still wants to rebuild targets, you have a timestamp or dependency problem.

```bash
make all
make -n all
# Expected: make: Nothing to be done for 'all'.
```

## Trace Why Targets Rebuild

If Make is rebuilding targets unnecessarily, use `make -d` to see its decision-making process. The output is verbose, so filter it.

```bash
# Find which files Make considers newer than their targets
make -d all 2>&1 | grep 'File.*is newer'

# Focus on a specific target
make -d myapp.o 2>&1 | grep -A5 'myapp.o'
```

For a more readable view on GNU Make 4.0 and later, use `--trace`:

```bash
make --trace all
```

This shows each rule that fires and what prerequisite triggered it. Look for targets that rebuild because of a header or generated file with a future timestamp.

## Check Variable Values

Unexpected variable expansion can cause performance issues. For example, using recursively-expanded variables (`=`) with `$(wildcard)` or `$(shell)` causes re-evaluation on every reference. Use `make -p` to inspect actual values.

```bash
# See the value of CFLAGS
make -p | grep "^CFLAGS"

# See all variables
make -p | grep "^[A-Z_]* ="
```

Look for variables defined with `=` (recursive) that contain expensive function calls. These re-evaluate every time the variable is referenced.

## Apply Targeted Fixes

After diagnosis, apply the fix that addresses the root cause you identified.

### Fix: Enable Parallel Builds

If the build is slow because compilation runs serially, enable parallelism.

```bash
# Run with 4 parallel jobs
make -j4 all

# Use all available cores
make -j$(nproc) all
```

To make this the default, add it to the Makefile:

```makefile
MAKEFLAGS += -j$(shell nproc)
```

Parallel builds require correct dependency declarations. If targets have undeclared dependencies, parallel builds may fail intermittently.

### Fix: Use Immediate Expansion

If `make -p` revealed recursive variables with expensive expansions, switch to immediate expansion.

```makefile
# Slow: re-evaluates $(wildcard) on every reference
SRCS = $(wildcard src/*.c)

# Fast: evaluates $(wildcard) once at parse time
SRCS := $(wildcard src/*.c)
```

Use `:=` (immediate) instead of `=` (recursive) for any variable whose value does not need to change during the build.

### Fix: Add Compiler Cache

If full rebuilds are frequent (such as switching branches), add ccache.

```makefile
CC := ccache gcc
CXX := ccache g++
```

ccache caches compilation results and returns cached output when the same file is compiled with the same flags. It has no effect on incremental builds but dramatically speeds up full rebuilds.

### Fix: Skip Dependency Parsing for Clean

If `make clean` is slow, it may be parsing auto-generated `.d` dependency files unnecessarily. Guard the include.

```makefile
ifneq ($(MAKECMDGOALS),clean)
  -include $(DEPS)
endif
```

### Fix: Resolve Timestamp Issues

If `make -d` showed files with unexpected timestamps (common with network filesystems, VMs, or generated files), the fix depends on the cause:

- **Clock skew**: synchronize clocks across build machines
- **Generated files**: use order-only prerequisites (`|`) for generated headers
- **Copy operations**: use `cp -p` to preserve timestamps

## Verify the Fix

After applying a fix, re-measure to confirm the improvement.

```bash
# Full build
time make clean && time make all

# Incremental build
touch src/main.c && time make all
```

Compare against your baseline numbers. If the improvement is not significant, return to the diagnosis step. The root cause may be different from what you assumed.

## Quick Diagnostic Reference

| Symptom | First Command | What to Look For |
|---------|---------------|------------------|
| Build is slow overall | `time make clean && time make all` | Where time is spent |
| Everything rebuilds | `make all && make -n all` | Should say "nothing to be done" |
| Single file change rebuilds too much | `make -d all 2>&1 \| grep 'is newer'` | Unexpected "newer" files |
| Variable has wrong value | `make -p \| grep "^VAR"` | Check `=` vs `:=` |
| Unknown execution order | `make --trace all` | Rule firing sequence |

## See Also

- [Skill Reference]({{< ref "reference/skills" >}}) -- makefile-debugging-optimization skill specification
- [Organize a Multi-Directory Build]({{< ref "howto/organize-multi-directory-build" >}}) -- fix parallelism issues in multi-directory builds
- [Skill Progression]({{< ref "explanation/skill-progression" >}}) -- how debugging fits the learning path
