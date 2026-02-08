---
title: "Modernize Legacy Code"
description: "Steps to modernize legacy code patterns using language-specific agents for C++, JavaScript, and TypeScript"
weight: 1
---

# Modernize Legacy Code

This guide covers three common legacy modernization tasks, each handled by a different language agent. Pick the section matching your language and follow the steps.

## Modernize C++ Raw Pointers to Smart Pointers

Use the cpp-pro agent to replace manual memory management with RAII and smart pointers.

### When to Use This

Your codebase has raw `new`/`delete` calls, manual resource management, or classes that violate the Rule of Five.

### Steps

1. Identify files with raw ownership. Look for patterns like `new`, `delete`, raw pointer returns from factory functions, and destructors that free resources manually.

2. Ask cpp-pro to modernize the ownership:

   ```
   Modernize the memory management in src/buffer.cpp to use smart pointers and RAII
   ```

3. The agent will apply these transformations:

   - Factory functions return `std::unique_ptr<T>` instead of `T*`
   - Owned members become `std::unique_ptr` (sole owner) or `std::shared_ptr` (shared ownership)
   - Non-owning references use raw pointers or `std::reference_wrapper` -- the agent does not blindly wrap every pointer
   - Classes with manual destructors are refactored so the destructor can be defaulted

4. Review the agent's summary for ABI or API changes. Converting a raw pointer return to `unique_ptr` changes the function signature and all callers must be updated.

5. The agent runs the build with `-fsanitize=address,undefined` to verify no memory errors remain. If your project uses CMake, it will add the flags to a debug configuration rather than modifying the release build.

### Common Pitfalls

- **Shared ownership disguised as unique ownership.** If multiple parts of the code hold onto a raw pointer, the agent may initially choose `unique_ptr` and then detect the sharing during compilation. It will escalate to `shared_ptr` with a comment explaining the shared ownership.
- **C-style APIs that require raw pointers.** The agent uses `.get()` at interop boundaries and documents why the raw pointer is safe at that call site.

## Convert JavaScript Callbacks to async/await

Use the javascript-pro agent to replace callback chains and raw Promise `.then()` patterns with `async`/`await`.

### When to Use This

Your Node.js code uses callback-style APIs (like the old `fs.readFile(path, callback)` pattern) or deeply nested `.then().catch()` chains.

### Steps

1. Identify callback-heavy files. Common signs: function parameters named `cb`, `callback`, or `done`; nested `.then()` chains more than two levels deep; error-first callback patterns `(err, result)`.

2. Ask javascript-pro to convert to async/await:

   ```
   Refactor the callback chains in lib/data-loader.js to use async/await
   ```

3. The agent will apply these transformations:

   - Callback-based function calls replaced with their Promise-based equivalents (e.g., `fs.promises.readFile` instead of `fs.readFile`)
   - `.then().catch()` chains flattened into sequential `await` calls with `try/catch` blocks
   - Error handling consolidated -- each `try/catch` block covers a logical unit of work rather than wrapping every individual `await`
   - Parallel operations that were sequentialized by nested callbacks restored to true parallelism with `Promise.all()`

4. Verify the runtime supports async/await natively. The agent checks your `package.json` for the `engines.node` field and your bundler configuration for target settings.

5. Run the existing test suite to verify behavior is preserved. The agent executes `npm test` or the project's configured test command.

### Common Pitfalls

- **Losing parallelism.** Converting `Promise.all([a(), b()])` into sequential `await a(); await b();` is a regression. The agent preserves parallel execution where the original code ran operations concurrently.
- **Error semantics changes.** Callback-style error handling sometimes swallows errors silently. The agent flags these cases rather than replicating the silent swallowing.

## Replace TypeScript `any` Types with Strict Typing

Use the typescript-pro agent to eliminate `any` types and enable strict mode in your TypeScript project.

### When to Use This

Your `tsconfig.json` has `strict: false` or omits it, and your codebase has `any` types scattered through function parameters, return types, or variable declarations.

### Steps

1. Get a count of `any` usage to understand the scope. Ask the agent to assess the current state:

   ```
   How many any types are in this TypeScript project and where are they concentrated?
   ```

2. Enable strict mode incrementally. Ask the agent to start:

   ```
   Tighten the types in this project and enable strict mode
   ```

3. The agent works in this order:

   - Enables `strict: true` in `tsconfig.json`
   - Replaces `any` at API boundaries with `unknown` and adds type guards
   - Replaces internal `any` types with concrete types inferred from usage
   - Converts type assertions (`as SomeType`) to type narrowing with discriminated unions where possible
   - Adds explicit return types to exported functions

4. Review compile errors. Enabling strict mode often surfaces real bugs -- the agent will fix the type errors it introduces rather than suppressing them with `@ts-ignore`.

5. The agent runs `tsc --noEmit` to verify the type system is fully satisfied.

### Common Pitfalls

- **Third-party types.** Some libraries have `any` in their type definitions. The agent wraps these with typed adapter functions rather than trying to modify `node_modules`.
- **Generic constraints.** Replacing `any` with `unknown` in generic type parameters may require adding constraints. The agent adds the minimum constraint needed.

## See Also

- {{< ref "howto/migrate-to-typescript" >}} -- Full project migration from JavaScript to TypeScript
- {{< ref "howto/optimize-performance" >}} -- Performance optimization after modernization
- {{< ref "reference/agents" >}} -- Complete specifications for all language agents
