---
title: "Agent Reference"
description: "Complete technical specifications for all five active language agents in the programming-languages plugin"
weight: 1
---

# Agent Reference

This page provides the complete specification for each active agent in the programming-languages plugin.

## cpp-pro

| Property | Value |
|----------|-------|
| **Model** | sonnet |
| **Color** | yellow |
| **Tools** | Read, Write, Edit, Bash |
| **C++ Standard** | C++17 baseline, C++20/23 features where supported |
| **Guidelines** | C++ Core Guidelines |

### Capabilities

- RAII resource management with `unique_ptr` (default) and `shared_ptr` (shared ownership)
- `constexpr` and compile-time computation
- STL algorithms and ranges over raw loops
- Value semantics with move semantics for expensive types
- Structured bindings, `std::optional`, concepts
- Concurrency analysis and synchronization patterns
- Build system configuration (CMake)
- Sanitizer-guided verification (`-fsanitize=address,undefined`)

### Trigger Conditions

- "Modernize this C++ code to use smart pointers and RAII"
- "Add C++20 concepts to constrain these templates"
- "Help me fix these data races in my C++ server"

### Constraints

- Does not introduce `dynamic_cast` or RTTI unless the design requires it
- Does not add template metaprogramming complexity for one-time-use code
- Does not use C-style casts
- Does not modify third-party or generated code

---

## go-simplifier

| Property | Value |
|----------|-------|
| **Model** | opus |
| **Color** | green |
| **Tools** | Read, Write, Edit, Grep, Glob, Bash |
| **Standards** | Effective Go, Go Code Review Comments |
| **Default Scope** | Recently modified files |

### Capabilities

- Idiomatic Go patterns: early returns, switch over if/else chains, descriptive naming
- Deprecated API replacement (e.g., `ioutil` to `io`/`os`)
- Error handling convention enforcement
- Dead code and redundant check elimination
- `go vet` verification after changes

### Trigger Conditions

- "Clean up the Go code I just wrote"
- "Make this more idiomatic Go"
- "Simplify the Go handlers I just added"

### Constraints

- Never changes behavior -- only changes how code is expressed
- Does not modify generated files (`// Code generated` header), vendored code, CGo blocks, or `//go:build` directives
- Does not introduce `unsafe`, `reflect`, or complex generics to shorten code
- Preserves performance-sensitive code marked with comments

### Exclusion List

The following file types are never modified:

- Files with `// Code generated` header
- Vendored dependencies
- CGo blocks
- Build directive files

---

## javascript-pro

| Property | Value |
|----------|-------|
| **Model** | sonnet |
| **Color** | green |
| **Tools** | Read, Write, Edit, Bash |
| **Environments** | Node.js, browser |
| **Module Default** | ES modules (`import`/`export`) |

### Capabilities

- `async`/`await` conversion from callbacks and Promise chains
- Dual-format packaging (ESM and CJS) with proper exports map
- Memory leak identification: event listeners, closures, unbounded caches
- Event loop analysis and main-thread blocking detection
- JSDoc comments on exported functions
- Destructuring, optional chaining, nullish coalescing
- Variable scoping: `const` default, `let` when needed, never `var`

### Trigger Conditions

- "Refactor this to use async/await instead of callbacks"
- "This Node.js service keeps running out of memory"
- "Make this module work as both ESM and CJS"

### Constraints

- Does not use `eval`, `with`, or `new Function()`
- Does not suppress errors with empty catch blocks
- Does not add polyfills without confirming browser target matrix
- Does not mix module systems without explicit dual-export configuration

---

## typescript-pro

| Property | Value |
|----------|-------|
| **Model** | sonnet |
| **Color** | cyan |
| **Tools** | Read, Write, Edit, Bash |
| **Strict Mode** | `strict: true` always, no exceptions |
| **Type Preference** | `interface` for extensible shapes, `type` for unions and computed types |

### Capabilities

- Strict configuration enforcement
- `any` elimination: replaces with `unknown` at boundaries, concrete types internally
- Discriminated unions for type narrowing
- Advanced generics with conditional types and mapped types
- JS-to-TS migration with `allowJs` and incremental conversion
- Type inference optimization: explicit types on exports, inferred on internals
- `tsc --noEmit` verification after changes

### Trigger Conditions

- "Tighten the types and enable strict mode"
- "Create a typed API wrapper with full inference"
- "Migrate this JavaScript project to TypeScript"

### Constraints

- Does not use `as` type assertions to silence errors
- Does not create deeply nested conditional types without extracting named helpers
- Does not add `@ts-ignore` or `@ts-expect-error` without a tracking comment
- Does not over-type simple internal functions

---

## zsh-expert

| Property | Value |
|----------|-------|
| **Model** | sonnet |
| **Color** | cyan |
| **Tools** | Read, Write, Edit, Bash, Grep, Glob |
| **Config Standard** | XDG Base Directory Specification |
| **Plugin Preference** | zinit or manual loading over oh-my-zsh |

### Capabilities

- `.zshrc` / `.zprofile` / `.zshenv` configuration and proper loading order
- Startup performance profiling and optimization with `zprof`
- Custom completion functions using `compdef` and `_arguments`
- Plugin framework setup and management (zinit, manual)
- Zsh-native features: parameter expansion flags, glob qualifiers, `zmv`
- Prompt configuration and theme setup

### Configuration File Responsibilities

| File | Purpose |
|------|---------|
| `.zshenv` | Environment exports, set once per session |
| `.zprofile` | Login shell setup |
| `.zshrc` | Interactive shell: aliases, functions, plugins, completions, prompt |

### Trigger Conditions

- "My terminal takes forever to load -- optimize my .zshrc"
- "Build custom zsh completions for my CLI tool"
- "Set up a new zsh config with a good prompt and proper PATH handling"

### Constraints

- Does not use bash-isms (`declare` -- uses `typeset` instead)
- Does not install frameworks without discussing performance tradeoffs
- Does not put interactive-only settings in `.zshenv`
