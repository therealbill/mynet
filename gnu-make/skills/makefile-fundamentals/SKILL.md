---
name: Makefile Fundamentals
description: >-
  This skill should be used when the user asks to create a Makefile, write a
  Makefile, generate a Makefile template, fix a broken Makefile, resolve a
  "missing separator" error, add build targets, set up a help target, or
  implement self-documenting Makefiles with the ## comment pattern. Also applies
  when reviewing Makefiles for best practices like .PHONY declarations, tab
  characters, and .DELETE_ON_ERROR. For pattern rules and automatic variables in
  depth, see the makefile-advanced-features skill.
version: 1.0.0
---

# Makefile Fundamentals

## Overview

Professional Makefiles are **self-documenting** using the `##` comment pattern to generate `make help` output. This pattern plus core best practices (`.PHONY`, automatic variables, proper tabs) creates maintainable build systems.

## When to Use

**Use this skill when:**
- Creating any new Makefile
- Fixing Makefile errors ("missing separator", broken targets)
- Adding targets to existing Makefiles
- User requests "simple" Makefile (simplicity ≠ skipping correctness)
- Reviewing or improving build systems

**Always suggest** the `##` help pattern, even if user didn't request it. It's a best practice that adds minimal complexity for significant value.

## The ## Self-Documentation Pattern

**Every Makefile should have a self-documenting help target:**

```makefile
.PHONY: help
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
```

Then document each target with `##` comments:

```makefile
.PHONY: build test clean

build: ## Build the project
	go build -o server main.go

test: ## Run all tests
	go test ./...

clean: ## Remove build artifacts
	rm -f server
```

Running `make help` shows:
```
build                Build the project
test                 Run all tests
clean                Remove build artifacts
```

### Why This Matters

- **Discoverability**: New team members run `make help` to see available targets
- **Documentation**: Comments stay next to code, never get out of sync
- **Standard**: Widely adopted pattern in professional projects
- **Low cost**: 4 lines of boilerplate, one `##` comment per target

## Essential Best Practices

### 1. Tab Characters (Critical)

**Recipes MUST use TAB characters for indentation**, not spaces.

```makefile
# ✅ CORRECT - TAB character
all:
	gcc main.c -o myapp

# ❌ WRONG - spaces cause "missing separator" error
all:
    gcc main.c -o myapp
```

**When fixing "missing separator" errors**: Replace spaces with TAB character.

### 2. .PHONY Declarations (Correctness)

**Declare non-file targets as `.PHONY`** to avoid conflicts:

```makefile
.PHONY: all clean test install help

all: myapp

clean:
	rm -f *.o myapp

test:
	./run_tests.sh
```

**Why**: Without `.PHONY`, if a file named `clean` or `test` exists, `make clean` won't run.

**This is correctness, not complexity.** Always include `.PHONY` even if user wants "simple" Makefile.

### 3. Automatic Variables (Maintainability)

Use automatic variables instead of hardcoding:

```makefile
# ✅ GOOD - uses automatic variables
%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

myapp: main.o utils.o
	$(CC) $(CFLAGS) $^ -o $@

# ❌ BAD - hardcoded, error-prone
main.o: main.c
	gcc -Wall main.c -o main.o

myapp: main.o utils.o
	gcc -Wall main.o utils.o -o myapp
```

**Key automatic variables:**
- `$@` = target name
- `$<` = first prerequisite
- `$^` = all prerequisites
- `$*` = stem (in pattern rules)

### 4. Standard Variables

Use conventional variable names:

```makefile
CC = gcc
CFLAGS = -Wall -O2
LDFLAGS = -lm
TARGET = myapp
SOURCES = main.c utils.c
OBJECTS = $(SOURCES:.c=.o)
```

### 5. Pattern Rules

Use pattern rules for scalability:

```makefile
# ✅ Handles any number of .c files
%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

# ❌ Doesn't scale
main.o: main.c
	$(CC) $(CFLAGS) -c main.c -o main.o

utils.o: utils.c
	$(CC) $(CFLAGS) -c utils.c -o utils.o
```

## Quick Reference

| Feature | Why | When to Skip |
|---------|-----|--------------|
| `##` help pattern | Discoverability | Never - always suggest |
| `.PHONY` | Correctness | Never - required for correctness |
| Automatic variables | Maintainability | Never - best practice |
| Variables (CC, CFLAGS) | Flexibility | Never - standard convention |
| Pattern rules | Scalability | Small projects with 1-2 files |
| `.DELETE_ON_ERROR` | Robustness | Never - prevents corruption |

## Complete Minimal Example

```makefile
# Variables
CC = gcc
CFLAGS = -Wall -O2
TARGET = myapp
SOURCES = main.c utils.c
OBJECTS = $(SOURCES:.c=.o)

# Targets
.PHONY: all clean help
.DELETE_ON_ERROR:

all: $(TARGET) ## Build the project

$(TARGET): $(OBJECTS)
	$(CC) $(CFLAGS) $^ -o $@

%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@

clean: ## Remove build artifacts
	rm -f $(OBJECTS) $(TARGET)

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
```

## Handling "Keep It Simple" Requests

When the user asks for simplicity, include .PHONY (correctness), tab characters (requirement), and the help target (discoverability) as non-negotiable fundamentals. Simplify by reducing variables and skipping advanced features like conditionals and functions. Explain the distinction: "These are essentials for correctness, not complexity."

## Common Mistakes

| Mistake | Fix | Why It Matters |
|---------|-----|----------------|
| Spaces instead of tabs | Replace with TAB | Causes "missing separator" error |
| No .PHONY | Add `.PHONY: clean test` etc | Prevents conflicts with files |
| Hardcoded filenames | Use `$@`, `$<`, `$^` | Maintainability and DRY |
| No help target | Add `##` pattern | Team discoverability |
| Missing .DELETE_ON_ERROR | Add to top | Prevents corrupted targets on failure |

## Proactive Suggestions

**When creating Makefiles:** Always include help target in initial version.

**When fixing Makefiles:** Suggest adding help target if missing: "While fixing this, I recommend adding a help target using the `##` pattern for better documentation."

**When adding targets:** Include `##` comments for new targets and suggest adding help infrastructure if not present.

## Red Flags - Review Your Makefile

- [ ] No help target with `##` pattern?
- [ ] Targets not declared in `.PHONY`?
- [ ] Hardcoded filenames instead of automatic variables?
- [ ] Spaces used instead of tabs in recipes?
- [ ] No `.DELETE_ON_ERROR` declaration?

**If any red flags present, suggest improvements.** These are standard practices, not "overcomplication."
