---
title: "Migrate to TypeScript"
description: "Steps to incrementally migrate a JavaScript project to TypeScript using the typescript-pro agent"
weight: 3
---

# Migrate to TypeScript

This guide walks through the complete process of migrating a JavaScript project to TypeScript using the typescript-pro agent. The strategy is incremental: you start by allowing JavaScript files alongside TypeScript, convert files one at a time, and progressively tighten the type system.

## Prerequisites

- A working JavaScript project with a `package.json`
- Node.js 16 or later
- The programming-languages plugin installed

## Step 1: Initialize TypeScript Configuration

Ask the typescript-pro agent to set up the project for incremental migration:

```
Set up TypeScript for incremental migration of this JavaScript project
```

The agent will create a `tsconfig.json` with these key settings:

```json
{
  "compilerOptions": {
    "strict": true,
    "allowJs": true,
    "checkJs": false,
    "outDir": "./dist",
    "rootDir": "./src",
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "resolveJsonModule": true,
    "declaration": true
  },
  "include": ["src/**/*"]
}
```

Key decisions the agent makes:

- **`allowJs: true`** enables JavaScript and TypeScript files to coexist in the same project
- **`checkJs: false`** prevents the TypeScript compiler from flagging errors in JavaScript files you have not converted yet
- **`strict: true` from the start** ensures every new TypeScript file is written to the highest standard rather than accumulating `any` types that need cleanup later

## Step 2: Install TypeScript and Type Definitions

The agent will install the necessary dependencies:

```
npm install --save-dev typescript @types/node
```

For any third-party libraries your project uses, the agent checks for available `@types/` packages and installs them. If a library does not have type definitions, the agent creates a minimal declaration file in a `types/` directory.

## Step 3: Convert Entry Points First

Ask the agent to begin converting files:

```
Start migrating the JavaScript files to TypeScript, beginning with the main entry points
```

The agent follows this conversion order:

1. **Entry points** (`index.js`, `app.js`, `server.js`) -- these establish the type contracts that flow through the rest of the codebase
2. **Shared utilities and types** -- modules imported by many other files, so typing them provides the most leverage
3. **Business logic modules** -- the bulk of the codebase, converted file by file
4. **Tests** -- converted last, since they benefit from all the types already defined

For each file, the agent:

- Renames `.js` to `.ts` (or `.jsx` to `.tsx`)
- Adds type annotations to function parameters and return types
- Replaces `require()` with `import` statements
- Defines interfaces for object shapes used across the file

## Step 4: Replace `any` Types Incrementally

After the initial conversion, some types will be `any` because the agent could not infer them from context alone. Ask the agent to tighten these:

```
Replace the remaining any types in src/services/ with proper types
```

The agent's strategy for eliminating `any`:

- **API response types** -- defined as interfaces matching the actual response shape, with `unknown` at the boundary and a type guard or validation function to narrow
- **Event handler parameters** -- typed using the library's event type definitions (e.g., `React.ChangeEvent<HTMLInputElement>`)
- **Dynamic objects** -- replaced with `Record<string, T>` or a discriminated union, depending on whether the keys are known
- **Function parameters from JavaScript callers** -- typed with `unknown` and narrowed at the top of the function body

## Step 5: Enable Additional Strict Checks

Once the majority of files are converted, ask the agent to enable additional checks:

```
Enable stricter TypeScript checks now that most files are converted
```

The agent enables these incrementally:

- **`noUncheckedIndexedAccess`** -- array and object index access returns `T | undefined`, catching out-of-bounds bugs
- **`exactOptionalPropertyTypes`** -- distinguishes between a missing property and a property explicitly set to `undefined`
- **`noImplicitOverride`** -- requires the `override` keyword on subclass methods, catching typos in method names

Each additional check may surface new type errors. The agent fixes them rather than suppressing them.

## Step 6: Remove JavaScript Allowance

When all files have been converted, finalize the migration:

```
Finalize the TypeScript migration -- disable allowJs and verify the full build
```

The agent will:

- Set `allowJs: false` in `tsconfig.json`
- Verify no `.js` source files remain in the `src/` directory
- Run `tsc --noEmit` to confirm the entire project compiles cleanly
- Update `package.json` scripts to use `tsc` for building
- Add a `"type": "module"` field if the project uses ESM

## Step 7: Set Up CI Type Checking

The agent recommends adding a type-checking step to your CI pipeline:

```json
{
  "scripts": {
    "typecheck": "tsc --noEmit",
    "build": "tsc",
    "pretest": "npm run typecheck"
  }
}
```

This ensures type errors are caught before tests run and before code is merged.

## Handling Common Migration Challenges

### Dynamic `require()` Calls

If your JavaScript code uses dynamic `require()` with variable paths, the agent converts these to dynamic `import()` with explicit type annotations for each possible module.

### Untyped Third-Party Libraries

For libraries without `@types/` packages, the agent creates a declaration file at `types/<library-name>.d.ts` with the minimum type surface needed by your code. This file can be expanded as you use more of the library's API.

### Test Framework Types

The agent installs type definitions for your test framework (`@types/jest`, `@types/mocha`, etc.) and converts test files to use typed assertions and matchers.

## See Also

- [The TypeScript section covers replacing `any` types in existing TypeScript projects](../../howto/modernize-legacy-code/)
- [Full specification of the typescript-pro agent](../../reference/agents/)
- [What typescript-pro knows about the type system that a general agent does not](../../explanation/language-specialization/)
