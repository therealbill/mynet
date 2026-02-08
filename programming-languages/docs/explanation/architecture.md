---
title: "Architecture"
description: "Why the programming-languages plugin uses one specialized agent per language instead of a single general-purpose coding agent"
weight: 1
---

# Architecture

The programming-languages plugin provides five separate agents, each dedicated to a single language or language family. This document explains why that architecture was chosen over the alternative of a single "code improver" agent that handles all languages.

## The Case Against a General Coding Agent

A general-purpose code review agent can identify surface-level issues in any language: naming conventions, dead code, overly long functions. But language-specific improvement requires knowledge that does not transfer between languages.

Consider these examples:

- **C++ ownership semantics.** Knowing when to use `unique_ptr` versus `shared_ptr` versus a raw non-owning pointer requires understanding C++'s ownership model, the Rule of Five, and RAII. These concepts do not exist in Go or JavaScript.

- **Go's error handling philosophy.** Go deliberately chose explicit error returns over exceptions. A general agent might "improve" Go error handling by suggesting a try/catch abstraction, which would be anti-idiomatic. The go-simplifier agent knows that explicit `if err != nil` checks are the intended pattern and focuses on making them cleaner, not eliminating them.

- **JavaScript's event loop.** Optimizing JavaScript performance requires understanding the event loop, microtask queue, and garbage collector behavior. None of this applies to C++ or Go performance work.

- **TypeScript's type system.** Discriminated unions, conditional types, mapped types, and template literal types form a type-level programming language within TypeScript. A general agent might use `any` to make code compile; the typescript-pro agent eliminates `any` types and uses the type system to catch bugs at compile time.

- **Zsh's completion system.** The `compdef`/`_arguments` API for zsh completions is entirely unique to Z Shell. No other language in this plugin uses anything similar.

A general agent would need to hold all of this context simultaneously, diluting its depth in each area. Specialized agents load only the knowledge relevant to their language, which produces higher-quality output.

## Why go-simplifier Uses the Opus Model

All other agents in this plugin use the sonnet model, which provides fast iteration and strong code generation. The go-simplifier agent uses opus, which has stronger reasoning capabilities.

Simplification is harder than generation. Writing new code requires following patterns; simplifying existing code requires judging which of several valid alternatives is genuinely clearer. This judgment demands stronger reasoning:

- Deciding whether a helper function reduces or increases cognitive load
- Evaluating whether early returns make a function clearer or fragment its logic
- Determining whether an abstraction is warranted or premature

The opus model's stronger reasoning capability makes these judgment calls more reliably. The tradeoff is slower response time, which is acceptable because go-simplifier operates on recently modified code (a small scope) rather than entire codebases.

## Tool Access Differences

Not all agents have the same tools. The tool assignments reflect what each agent needs:

| Agent | Additional Tools | Reason |
|-------|-----------------|--------|
| go-simplifier | Grep, Glob | Needs to find recently modified files across the codebase and search for usage patterns |
| zsh-expert | Grep, Glob | Needs to search for configuration files, find existing completions, and locate shell scripts |
| cpp-pro | (base set) | Works on files the user points it to; build system handles discovery |
| javascript-pro | (base set) | Works on files the user points it to; npm scripts handle discovery |
| typescript-pro | (base set) | Works on files the user points it to; tsconfig handles discovery |

The base tool set for all agents is Read, Write, Edit, and Bash. Agents that need to discover files on their own (rather than being pointed to specific files) also have Grep and Glob.

## Relationship to the Code-Quality Plugin

The programming-languages plugin handles language-specific code improvement. For language-agnostic concerns -- such as code structure, documentation quality, test coverage, and design patterns -- the code-quality plugin provides separate agents.

The two plugins are complementary:

- Use a programming-languages agent when the improvement requires language-specific knowledge (ownership models, type systems, runtime behavior)
- Use a code-quality agent when the improvement applies to any language (function length, naming clarity, test organization)

You can use both in sequence: first run a language agent to apply idiomatic patterns, then run a code-quality agent to assess overall structure.

## Design Principles

The architecture follows these principles:

- **One agent, one language.** Each agent is an expert in exactly one language. This prevents dilution and ensures every suggestion is grounded in that language's conventions.

- **Behavior preservation.** All agents change how code is expressed, not what it does. This makes the agents safe to use on production code.

- **Verification after changes.** Each agent runs the language's standard verification tool (`go vet`, `tsc --noEmit`, sanitizers, `eslint`) after making changes. The agent does not consider its work done until verification passes.

- **Transparent reasoning.** Every agent summarizes what it changed and why. The user can accept, reject, or refine any individual change.
