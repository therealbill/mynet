# Debugging Scenarios and Quick Reference

Reference material for the Makefile Debugging and Optimization skill.

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
