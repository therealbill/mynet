# Pattern Rule Scenarios

Reference material for the Advanced Makefile Features skill.

### Scenario 1: Compile C to Object Files

```makefile
%.o: %.c
	$(CC) $(CFLAGS) $(CPPFLAGS) -c -o $@ $<
```

Handles unlimited .c files automatically.

### Scenario 2: Compile C Directly to Executables

```makefile
# From baseline test: 8 repetitive rules → 1 pattern rule
APPS := app1 app2 app3 app4 app5 app6 app7 app8

all: $(APPS)

%: %.c
	$(CC) $(CFLAGS) -o $@ $<

clean:
	rm -f $(APPS)

.PHONY: all clean
```

### Scenario 3: Multiple Extensions

```makefile
# C files
%.o: %.c
	$(CC) $(CFLAGS) -c -o $@ $<

# C++ files
%.o: %.cpp
	$(CXX) $(CXXFLAGS) -c -o $@ $<

# Assembly files
%.o: %.s
	$(AS) $(ASFLAGS) -o $@ $<
```

### Scenario 4: Auto-Discovery with Wildcard

```makefile
# Automatically find all .c files
SRCS := $(wildcard *.c)
OBJS := $(SRCS:.c=.o)
DEPS := $(SRCS:.c=.d)

TARGET := myapp

all: $(TARGET)

$(TARGET): $(OBJS)
	$(CC) -o $@ $^

%.o: %.c
	$(CC) $(CFLAGS) -c -o $@ $<

# Auto-generate dependencies
%.d: %.c
	$(CC) -MM $(CFLAGS) $< > $@

-include $(DEPS)

clean:
	rm -f $(OBJS) $(DEPS) $(TARGET)

.PHONY: all clean
```

**Benefit:** Add new .c file → no Makefile changes needed.

### Scenario 5: Static Pattern Rules (When Subset Needs Pattern)

```makefile
OBJS := main.o utils.o parser.o

# Apply pattern rule only to these objects
$(OBJS): %.o: %.c
	$(CC) $(CFLAGS) -c -o $@ $<
```

**Use when:** Pattern applies to specific list, not all files.
