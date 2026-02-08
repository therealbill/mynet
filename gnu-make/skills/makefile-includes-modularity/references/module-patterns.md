# Module Patterns

Reference material for the Makefile Includes and Modularity skill.

### Pattern 1: Three-File Foundation

```makefile
# Makefile (main)
include config.mk
include rules.mk

all: $(TARGET)
	@echo "Build complete"
```

**Use when:** 150-300 line monolith, clear config vs rules split.

### Pattern 2: Component-Based

```makefile
# Makefile
include config.mk
include core.mk
include api.mk
include cli.mk

all: core api cli
```

**Use when:** Multiple components, each with own rules.

### Pattern 3: Environment + Component

```makefile
# Makefile
ENV ?= dev
include config/$(ENV).mk
include core.mk
include api.mk

all: $(TARGETS)
```

**Use when:** Multiple components AND multiple environments.

### Pattern 4: Shared Common

```makefile
# Makefile
include ../common/common.mk  # Shared across projects
include project-specific.mk

all: $(TARGET)
```

**Use when:** Multiple projects sharing configuration.
