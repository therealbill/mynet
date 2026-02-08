---
title: "Language Specialization"
description: "What each language agent knows that a general-purpose code agent does not, covering idioms, patterns, and toolchain conventions"
weight: 2
---

# Language Specialization

Each agent in the programming-languages plugin carries deep, language-specific knowledge that a general-purpose coding agent lacks. This document catalogs the specialized knowledge domains for each agent and explains why this depth matters for code quality.

## cpp-pro: Modern C++ Expertise

### RAII and Ownership Semantics

The cpp-pro agent understands C++'s resource ownership model at a level that goes beyond "use smart pointers." It knows:

- **When `unique_ptr` is correct vs. `shared_ptr` vs. raw pointers.** A raw pointer is appropriate for non-owning references where the lifetime is guaranteed by a parent object. `unique_ptr` is the default for owned resources. `shared_ptr` is reserved for genuinely shared ownership, because it adds reference-counting overhead and makes ownership graphs harder to reason about.

- **The Rule of Five.** If a class defines any of the five special member functions (destructor, copy constructor, copy assignment, move constructor, move assignment), it should define or explicitly default all five. The agent detects violations and either adds the missing functions or refactors the class so the compiler-generated defaults are correct.

- **Exception safety guarantees.** The agent writes code that provides the strong exception guarantee where possible (operations either succeed or leave state unchanged) and documents where only the basic guarantee is practical.

### Template Metaprogramming and Concepts

The agent knows how to constrain templates using C++20 concepts to produce clear error messages instead of pages of template instantiation failures. It understands:

- Writing `requires` clauses that express the minimum constraint needed
- Using `concept` definitions to name and reuse constraints
- When SFINAE is still necessary (pre-C++20 codebases) versus when concepts replace it
- The difference between syntactic constraints (the expression compiles) and semantic constraints (the expression does what we expect)

### Concurrency Patterns

The agent understands C++ concurrency beyond "add a mutex":

- Lock hierarchies to prevent deadlocks
- `std::atomic` with appropriate memory orderings (not just `memory_order_seq_cst` everywhere)
- When `std::jthread` and `std::stop_token` (C++20) simplify thread lifecycle management
- Thread sanitizer output interpretation to fix data races systematically

## go-simplifier: Idiomatic Go

### The Go Philosophy

The go-simplifier agent internalizes Go's design philosophy: simplicity, explicitness, and readability over cleverness. This is not just a style preference -- it affects concrete decisions:

- **Early returns over deep nesting.** A function with four levels of `if` nesting restructured to check error conditions first and return early. This is not just aesthetically cleaner; it aligns with how Go developers scan code.

- **Switch over if/else chains.** When dispatching on a value (HTTP methods, error types, states), a switch statement is both more readable and signals to the reader that the cases are exhaustive.

- **Named return values for documentation, not for naked returns.** The agent uses named returns to document what a function returns but avoids naked `return` statements because they obscure what values are being returned.

### Error Handling Conventions

Go's explicit error handling is its most distinctive feature. The agent knows:

- The `if err != nil { return ..., fmt.Errorf("context: %w", err) }` pattern and when wrapping is appropriate versus when it adds noise
- How to use `errors.Is` and `errors.As` for error inspection instead of string matching
- When sentinel errors are appropriate versus custom error types
- That error messages should not be capitalized or end with punctuation (per Go convention)

### Standard Library Patterns

The agent knows which standard library functions to prefer:

- `io.ReadAll` over the deprecated `ioutil.ReadAll`
- `json.NewDecoder` for HTTP bodies instead of `json.Unmarshal` on a pre-read buffer
- `context.Context` threading through call chains for cancellation and timeouts
- `http.HandlerFunc` adapter pattern for registering handlers

## javascript-pro: Runtime and Module Expertise

### Event Loop Model

The javascript-pro agent understands the JavaScript event loop at a level required for performance work:

- **Macrotasks vs. microtasks.** `setTimeout` schedules macrotasks; `Promise.then` schedules microtasks. Microtasks run before the next macrotask, which means a loop creating Promises can starve I/O.

- **Blocking the event loop.** CPU-intensive synchronous work in a `for` loop blocks all I/O, timers, and incoming connections. The agent identifies these patterns and recommends `worker_threads` or chunked processing with `setImmediate`.

- **Garbage collector pressure.** Creating many short-lived objects (common in functional-style code with lots of `map`/`filter`/`reduce` chains) can trigger frequent GC pauses. The agent identifies when object allocation patterns matter for latency-sensitive code.

### Module Systems

The agent deeply understands both ES modules and CommonJS:

- **Dual-format packaging.** Configuring `package.json` `"exports"` with `"import"` and `"require"` conditions, separate entry points, and correct `"type"` field.

- **Tree-shaking compatibility.** Writing code so that bundlers can eliminate unused exports. This requires understanding what makes a module "side-effect-free" and how the `"sideEffects"` field in `package.json` works.

- **Circular dependency resolution.** Diagnosing and fixing circular imports, which behave differently in ESM (live bindings, partial initialization) versus CJS (cached partial objects).

### Async Patterns

Beyond basic `async`/`await`, the agent knows:

- `Promise.allSettled` for operations where partial failure is acceptable
- `AbortController` and `AbortSignal` for cancellation
- Async iterators (`for await...of`) for streaming data
- The `pipeline` utility from `stream/promises` for composing Node.js streams without listener leaks

## typescript-pro: Type System Mastery

### Advanced Type System Features

The typescript-pro agent operates at a level of type system sophistication that is inaccessible to a general-purpose agent:

- **Discriminated unions.** Defining a `type` field on each variant and using `switch` on it for exhaustive narrowing. The agent knows how to structure data types so TypeScript's control flow analysis eliminates the need for runtime type checks.

- **Conditional types.** `T extends U ? X : Y` patterns for type-level branching. The agent uses these to build type-safe API wrappers where the return type depends on the input type.

- **Mapped types.** `{ [K in keyof T]: Transform<T[K]> }` for transforming object types systematically. Used for creating readonly versions, optional versions, or pick/omit utilities.

- **Template literal types.** `\`${Prefix}/${Path}\`` for type-safe string manipulation, commonly used in routing and API endpoint definitions.

- **Infer keyword.** `T extends Promise<infer U> ? U : T` for extracting types from generic positions. The agent uses this for utility types that unwrap Promises, extract function parameters, or decompose complex types.

### Strict Mode Configuration

The agent knows every flag under `strict: true` and what each one catches:

- `strictNullChecks` -- prevents `null`/`undefined` from being assignable to every type
- `strictFunctionTypes` -- enables contravariant function parameter checking
- `strictBindCallApply` -- types `bind`, `call`, and `apply` correctly
- `strictPropertyInitialization` -- ensures class properties are assigned in the constructor
- `noImplicitAny` -- requires explicit types where inference fails
- `noImplicitThis` -- requires explicit `this` parameter types
- `useUnknownInCatchVariables` -- types catch clause variables as `unknown` instead of `any`
- `alwaysStrict` -- emits `"use strict"` in every file

### Migration Strategy Knowledge

The agent understands the practical realities of migrating large JavaScript codebases:

- Which files to convert first for maximum type-safety leverage
- How to write `.d.ts` declaration files for untyped libraries
- How `allowJs` and `checkJs` interact and when to enable each
- Strategies for handling dynamically-typed patterns (event emitters, configuration objects) that resist static typing

## zsh-expert: Shell Ecosystem Knowledge

### Completion System

The zsh completion system is one of the most powerful and least documented features of any shell. The agent knows:

- **`compdef` function registration.** Associating a completion function with a command name so the function is called when the user presses Tab.

- **`_arguments` specification syntax.** The DSL for declaring subcommands, flags with arguments, mutually exclusive options, and file completion for specific arguments. This syntax (`'(-v --verbose)'{-v,--verbose}'[Enable verbose output]'`) is dense and error-prone without deep familiarity.

- **`_describe` and `_values` helpers.** For completing from a list of options with descriptions, or for completing colon-separated value specifications.

- **Caching with `_store_cache` and `_retrieve_cache`.** For completions that are expensive to compute (such as listing remote resources), storing the results and refreshing them on a timer.

### Startup Optimization

The agent understands zsh's loading sequence and how to optimize each phase:

- **Login vs. interactive vs. script shells.** Which configuration files load in which order (`.zshenv` always, `.zprofile` for login, `.zshrc` for interactive) and what belongs in each.

- **Lazy-loading with `zinit wait`.** Deferring plugin loading until after the prompt is drawn, so the user sees their shell immediately while plugins load in the background.

- **`compinit` caching.** The `compinit` call rebuilds the completion system from all installed completion functions. By checking the dump file's modification time and only rebuilding daily, startup can be reduced by hundreds of milliseconds.

- **`zprof` profiling.** Wrapping the configuration in `zmodload zsh/zprof` at the top and `zprof` at the bottom to get a time-sorted breakdown of every function call during startup.

### XDG Base Directory Convention

The agent follows XDG conventions for configuration file placement:

- `$XDG_CONFIG_HOME/zsh/` (defaulting to `~/.config/zsh/`) for configuration
- `$XDG_DATA_HOME/zsh/` for history and other data
- `$XDG_CACHE_HOME/zsh/` for completion caches and compiled files

This keeps the home directory clean and makes configuration portable across machines.

## Why Depth Matters

A general agent can apply formatting rules and catch obvious issues. But the improvements that matter most -- ownership correctness in C++, idiomatic error handling in Go, event loop awareness in JavaScript, type-level programming in TypeScript, completion system authoring in zsh -- require depth that can only come from specialization.

Each agent in this plugin represents a concentration of language-specific knowledge that would be diluted in a general-purpose tool. The result is that every suggestion is grounded in the conventions, patterns, and best practices of the specific language being reviewed.

## See Also

- {{< ref "explanation/architecture" >}} -- The architectural rationale for one agent per language
- {{< ref "reference/agents" >}} -- Complete technical specifications for each agent
- {{< ref "tutorials/getting-started" >}} -- Hands-on experience with agent specialization
