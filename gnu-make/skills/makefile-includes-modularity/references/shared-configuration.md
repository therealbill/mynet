# Shared Configuration Across Projects

Reference material for the Makefile Includes and Modularity skill.

### Common Configuration File

**Project Structure:**
```
common/
  └── common.mk       # Shared across all projects
project1/
  ├── Makefile
  └── src/
project2/
  ├── Makefile
  └── src/
```

```makefile
# common/common.mk - Shared configuration
CC := gcc
CFLAGS := -Wall -Wextra -std=c11
PREFIX := /usr/local

%.o: %.c
	$(CC) $(CFLAGS) -c -o $@ $<
```

```makefile
# project1/Makefile
include ../common/common.mk

TARGET := app1
SOURCES := main.c utils.c
OBJECTS := $(SOURCES:.c=.o)

# Project-specific additions
CFLAGS += -DAPP_NAME=\"app1\"

$(TARGET): $(OBJECTS)
	$(CC) -o $@ $^ $(LDFLAGS)
```

**Benefits:**
- Change compiler flags once, affects all projects
- Consistent build process across projects
- DRY principle at project suite level
