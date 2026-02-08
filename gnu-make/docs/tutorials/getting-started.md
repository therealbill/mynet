---
title: "Getting Started with GNU Make"
description: "Create a professional Makefile from scratch with self-documenting help, pattern rules, and essential best practices"
weight: 1
---

# Getting Started with GNU Make

In this tutorial you will create a Makefile from scratch for a small C project. By the end you will have a working build system with a self-documenting help target, pattern rules for compilation, and the essential best practices that every professional Makefile needs.

## What You Will Build

A Makefile that compiles a multi-file C project with the following capabilities:

- A `help` target that prints available commands
- Pattern rules that scale to any number of source files
- Proper `.PHONY` declarations for correctness
- Standard variables for compiler and flags

## Prerequisites

- A Unix-like environment (Linux, macOS, or WSL)
- GNU Make installed (`make --version` to verify)
- A C compiler such as gcc or clang
- A text editor that can insert literal TAB characters

## Step 1: Set Up the Project

Create a small project with two C source files.

```
myproject/
  main.c
  utils.c
  utils.h
```

Create `main.c`:

```c
#include <stdio.h>
#include "utils.h"

int main(void) {
    printf("Result: %d\n", add(2, 3));
    return 0;
}
```

Create `utils.h`:

```c
#ifndef UTILS_H
#define UTILS_H
int add(int a, int b);
#endif
```

Create `utils.c`:

```c
#include "utils.h"

int add(int a, int b) {
    return a + b;
}
```

## Step 2: Create the Makefile Foundation

Create a file named `Makefile` (capital M, no extension) in the project root. Open it in your editor and add the following content. Every indented line in a recipe **must** use a literal TAB character, not spaces. Using spaces will produce a "missing separator" error.

```makefile
# Variables
CC = gcc
CFLAGS = -Wall -O2
TARGET = myapp
SOURCES = main.c utils.c
OBJECTS = $(SOURCES:.c=.o)
```

This block defines the standard variables. `CC` and `CFLAGS` are conventional names that Make and other tools recognize. The `OBJECTS` line uses substitution syntax to convert `main.c utils.c` into `main.o utils.o`.

## Step 3: Add the Self-Documenting Help Target

The `##` help pattern is a best practice that makes your Makefile self-documenting. Add the following below the variables:

```makefile
# Non-negotiable declarations
.PHONY: all clean help
.DELETE_ON_ERROR:
```

`.PHONY` tells Make that `all`, `clean`, and `help` are not actual files. Without this, if a file named `clean` ever exists in your directory, `make clean` would silently do nothing. `.DELETE_ON_ERROR` removes partially-written target files when a recipe fails, preventing corrupted build artifacts.

Now add the help target:

```makefile
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
```

This target scans the Makefile for lines matching the pattern `target: ## description` and formats them into a help listing. The `@` prefix suppresses printing the command itself.

## Step 4: Add Build Targets with Help Comments

Add the core build targets, each annotated with a `##` comment:

```makefile
all: $(TARGET) ## Build the project

$(TARGET): $(OBJECTS)
	$(CC) $(CFLAGS) $^ -o $@

clean: ## Remove build artifacts
	rm -f $(OBJECTS) $(TARGET)
```

Notice the automatic variables in the linking rule:

- `$^` expands to all prerequisites (`main.o utils.o`)
- `$@` expands to the target name (`myapp`)

These variables keep the recipe generic and maintainable. If you add more source files to `SOURCES`, the linking rule needs no changes.

## Step 5: Add a Pattern Rule

Instead of writing separate compilation rules for each `.c` file, add a single pattern rule:

```makefile
%.o: %.c
	$(CC) $(CFLAGS) -c $< -o $@
```

The `%` acts as a wildcard. When Make needs `main.o`, it matches `%.o` with stem `main`, looks for `main.c` as the prerequisite, and runs the recipe. The automatic variable `$<` expands to the first prerequisite (the `.c` file).

This single rule handles any number of source files. Without it, you would need a separate rule for every `.c` file in the project.

## Step 6: Verify the Complete Makefile

Your complete Makefile should look like this:

```makefile
# Variables
CC = gcc
CFLAGS = -Wall -O2
TARGET = myapp
SOURCES = main.c utils.c
OBJECTS = $(SOURCES:.c=.o)

# Declarations
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

## Checkpoint: Test Your Makefile

Run `make help` to see the self-documenting output:

```
$ make help
all                  Build the project
clean                Remove build artifacts
help                 Show this help message
```

If you see this output, the `##` pattern is working correctly.

Now build and run the project:

```
$ make all
gcc -Wall -O2 -c main.c -o main.o
gcc -Wall -O2 -c utils.c -o utils.o
gcc -Wall -O2 main.o utils.o -o myapp

$ ./myapp
Result: 5
```

Verify that incremental builds work. Run `make all` again without changing anything:

```
$ make all
make: 'myapp' is up to date.
```

Make detects that nothing changed and skips the build. Now touch one file and rebuild:

```
$ touch utils.c
$ make all
gcc -Wall -O2 -c utils.c -o utils.o
gcc -Wall -O2 main.o utils.o -o myapp
```

Only the changed file recompiles. This is the fundamental value of Make: it rebuilds only what is necessary.

Finally, clean up:

```
$ make clean
rm -f main.o utils.o myapp
```

## What You Learned

This tutorial covered the core concepts from the **makefile-fundamentals** and **makefile-advanced-features** skills:

- **TAB characters** are required for recipe indentation (spaces cause "missing separator" errors)
- **`.PHONY`** prevents conflicts with files that share target names
- **`.DELETE_ON_ERROR`** removes corrupted targets when recipes fail
- **The `##` help pattern** makes Makefiles self-documenting with `make help`
- **Standard variables** (`CC`, `CFLAGS`) provide a single place to change settings
- **Pattern rules** (`%.o: %.c`) handle any number of files with one rule
- **Automatic variables** (`$@`, `$<`, `$^`) keep recipes generic and maintainable

## Next Steps

- [Debug a Slow Makefile]({{< relref "../howto/debug-slow-makefile" >}}) -- learn to diagnose build issues with Make's built-in flags
- [Organize a Multi-Directory Build]({{< relref "../howto/organize-multi-directory-build" >}}) -- scale your build across subdirectories
- [Split a Large Makefile]({{< relref "../howto/split-large-makefile" >}}) -- modularize as your project grows
- [Skill Reference]({{< relref "../reference/skills" >}}) -- technical specifications for all five skills
