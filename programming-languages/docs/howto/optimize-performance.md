---
title: "Optimize Performance"
description: "Steps to optimize code performance in C++, Go, JavaScript, TypeScript, and Zsh using language-specific agents"
weight: 2
---

# Optimize Performance

Each language agent has a specific performance focus aligned with what matters most in that language's runtime. This guide covers the performance optimization workflow for each supported language.

## C++ Performance with cpp-pro

### Focus Areas

- RAII and deterministic destruction to eliminate resource leaks
- Move semantics to avoid unnecessary copies
- Concurrency correctness with proper synchronization
- Compile-time computation with `constexpr`

### Steps

1. Identify the performance-sensitive code. Ask cpp-pro to analyze:

   ```
   Profile the hot path in src/engine.cpp for unnecessary copies and allocation overhead
   ```

2. The agent applies these optimizations:

   - Replaces copy-by-value with move semantics for large objects
   - Converts runtime computations to `constexpr` where inputs are known at compile time
   - Replaces raw loops with STL algorithms that can be auto-vectorized
   - Identifies lock contention and suggests lock-free alternatives or finer-grained locking

3. The agent builds with sanitizers to verify that optimizations do not introduce data races or undefined behavior.

### What the Agent Will Not Do

The agent does not apply micro-optimizations that sacrifice readability (such as manual loop unrolling or bit manipulation tricks) unless you specifically request them and explain the performance requirement.

## Go Performance with go-simplifier

### Focus Areas

- Idiomatic patterns that the Go compiler optimizes well
- Eliminating unnecessary allocations
- Proper use of goroutines and channels
- Interface compliance without reflection

### Steps

1. Ask go-simplifier to review for performance-relevant patterns:

   ```
   Simplify the Go code in pkg/server/ and flag any patterns that cause unnecessary allocations
   ```

2. The agent applies idiomatic patterns that happen to be faster:

   - Pre-sized slices with `make([]T, 0, n)` instead of growing via `append`
   - `strings.Builder` instead of repeated string concatenation
   - Early returns that reduce the amount of code executing in the common case
   - `sync.Pool` for frequently allocated and discarded objects when appropriate

3. The agent runs `go vet` to verify correctness after changes.

### What the Agent Will Not Do

The agent prioritizes clarity over raw speed. It will not introduce `unsafe` operations, reflection-based shortcuts, or complex generics solely for performance.

## JavaScript Performance with javascript-pro

### Focus Areas

- Event loop optimization to avoid blocking the main thread
- Garbage collection pressure reduction
- Memory leak identification from event listeners and closures
- Efficient async patterns

### Steps

1. Ask javascript-pro to analyze performance:

   ```
   Identify performance bottlenecks in lib/data-processor.js, especially event loop blocking and memory leaks
   ```

2. The agent applies these optimizations:

   - Moves CPU-intensive work to `worker_threads` or breaks it into chunks with `setImmediate` to avoid blocking
   - Identifies unbounded caches, event listener leaks, and closure retention issues
   - Replaces `Promise.all()` with `Promise.allSettled()` where partial failure handling improves resilience
   - Uses streaming APIs (`ReadableStream`, `pipeline`) instead of buffering entire payloads in memory

3. The agent runs the project's test suite and linter to verify behavior is preserved.

### What the Agent Will Not Do

The agent will not add micro-optimizations like replacing `forEach` with `for` loops without evidence of measurable impact. It focuses on algorithmic and architectural improvements.

## TypeScript Performance with typescript-pro

### Focus Areas

- Type narrowing to help the compiler eliminate runtime checks
- Discriminated unions for efficient branching
- Avoiding runtime type checking overhead from excessive type guards
- Build-time type computation to reduce emitted JavaScript

### Steps

1. Ask typescript-pro to review for type-level performance:

   ```
   Optimize the type narrowing in src/router.ts to reduce runtime overhead from type guards
   ```

2. The agent applies these optimizations:

   - Replaces chains of `typeof` / `instanceof` checks with discriminated unions that narrow in a single comparison
   - Converts runtime type validation to compile-time checks where possible
   - Simplifies complex conditional types that slow down the TypeScript compiler (large union distributions, deep recursive types)
   - Uses `const` assertions and literal types to push computation to compile time

3. The agent runs `tsc --noEmit` to verify type correctness and may report compiler performance improvement for large projects.

### What the Agent Will Not Do

The agent will not sacrifice type safety for performance. It will not introduce `any` types or `as` assertions to skip runtime checks.

## Zsh Startup Performance with zsh-expert

### Focus Areas

- Shell startup time optimization
- Lazy-loading plugins and completions
- Efficient script execution patterns

### Steps

1. Ask zsh-expert to profile your shell startup:

   ```
   Optimize my slow .zshrc startup -- it takes over 3 seconds to load
   ```

2. The agent applies these optimizations:

   - Profiles startup with `zprof` to identify the slowest components
   - Defers plugin loading with zinit's `wait` ice or manual `autoload` patterns
   - Moves completion initialization (`compinit`) to a cached or lazy-loaded setup
   - Replaces heavy framework features (oh-my-zsh) with lightweight alternatives that load only what you use
   - Moves environment exports from `.zshrc` to `.zshenv` so they are set once per session, not per shell

3. The agent tests the new configuration by sourcing it in a subshell to verify no errors.

### What the Agent Will Not Do

The agent will not remove plugins without discussing the tradeoff. It explains the startup cost of each plugin and lets you decide what to keep.

## See Also

- [Modernizing legacy patterns often improves performance as a side effect](../../howto/modernize-legacy-code/)
- [Why language-specific knowledge matters for performance](../../explanation/language-specialization/)
- [Complete agent specifications including tool access](../../reference/agents/)
