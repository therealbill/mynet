---
title: "Architecture"
description: "Framework specialization and the web development agent model"
weight: 1
---

# Architecture

Why the web-development plugin uses six specialized agents instead of one general-purpose web agent, and how their responsibilities relate to each other.

## The Framework Specialization Model

Web development spans a wide range of concerns: visual design, component architecture, framework-specific patterns, and full-stack integration. A single "web developer" agent would need to hold opinions about React hooks, Vue composables, Next.js rendering strategies, Pinia stores, design tokens, and database schema design simultaneously. The system prompt would be enormous, the opinions would conflict (React's state model differs from Vue's reactivity model), and the agent would be mediocre at everything rather than strong at anything.

The web-development plugin splits this space into six agents, each with a focused domain, a tight set of opinions, and a clear "Do Not" list. A focused agent produces better results because its system prompt contains only the patterns, constraints, and defaults relevant to its domain. The react-specialist never needs to think about Pinia stores. The vue-expert never needs to consider React Server Components. This focus means each agent's opinions are internally consistent and its recommendations are specific rather than generic.

## The Scope Spectrum

The six agents sit along a spectrum from narrow visual focus to broad cross-layer integration:

**ui-designer** operates at the design level. It handles user journey mapping, design tokens, component specifications with all interaction states, and accessibility requirements. It produces design decisions -- spacing, color, typography, state definitions -- but does not implement framework-specific code.

**frontend-developer** operates at the implementation level. It takes design decisions and turns them into components, handling responsive layouts, state management, and client-side performance. It is framework-agnostic: it works with whatever framework the project already uses.

**react-specialist**, **vue-expert**, and **nextjs-developer** operate at the framework-specific level. They bring deep knowledge of one framework's patterns, idioms, and ecosystem. The react-specialist knows when to use `useReducer` vs. `useState`. The vue-expert knows when `ref` is safer than `reactive`. The nextjs-developer knows when ISR is better than dynamic rendering. This specificity produces recommendations that are correct for the framework rather than generically plausible.

**fullstack-developer** operates across all layers. It builds vertical slices -- one feature, from database schema through API routes to frontend components. It uses the opus model because cross-layer reasoning (keeping types aligned between client and server, designing transactions that span multiple database operations, wiring error states from API through UI) benefits from stronger analytical capability.

## Cross-Agent Collaboration

The agents do not coordinate with each other through a shared protocol. Instead, they collaborate implicitly through the codebase. A typical workflow might proceed:

1. **ui-designer** establishes the design system -- tokens, component specs, interaction states
2. **react-specialist** (or **vue-expert**) implements the components following the design specifications
3. **fullstack-developer** wires the components to API routes and database operations

Each agent reads what the previous one produced and builds on it. The design tokens the ui-designer creates become the Tailwind classes the react-specialist uses. The component interfaces the react-specialist defines become the contracts the fullstack-developer implements on the backend.

This implicit coordination works because the agents operate on the same codebase and follow consistent conventions. The ui-designer produces Tailwind classes or CSS custom properties, not abstract color names. The react-specialist produces TypeScript interfaces, not informal prop descriptions. These concrete artifacts serve as the handoff mechanism.

## Model Selection Rationale

Five of the six agents use the sonnet model. The fullstack-developer uses opus.

The distinction is about the nature of the reasoning required. The sonnet-powered agents work within a single domain: React component architecture, Vue composable design, Next.js rendering strategies, or visual design specifications. Their tasks require deep expertise in one area but do not require reasoning across multiple technical domains simultaneously.

The fullstack-developer's tasks are fundamentally cross-domain. Building a user registration feature requires simultaneously reasoning about database constraints (unique email, password hashing), API contract design (request/response shapes, error codes), authentication flows (session management, token lifecycle), and frontend state management (form validation, optimistic updates, error display). A mistake in any layer can create a bug that manifests in another layer. The opus model's stronger cross-domain reasoning reduces these integration errors.

## What Is Not Covered

The web-development plugin focuses on frontend implementation, framework-specific patterns, UI design, and full-stack feature development. Several adjacent concerns are handled by other plugins:

- **Backend system design** -- API architecture, microservice patterns, message queues, and infrastructure concerns belong to the backend-development plugin
- **End-to-end testing** -- Test strategy, test implementation, and quality assurance belong to the code-quality plugin
- **CLI tools and developer experience** -- Build tooling, developer scripts, and CLI applications belong to the cli-development plugin

The boundary is practical: if the task is "build a feature that includes a frontend," use web-development agents. If the task is "design the backend architecture that the frontend will consume," use backend-development agents. If the task is "write tests for the feature," use code-quality agents. The fullstack-developer bridges the web and backend domains for tasks that cannot be cleanly separated, but it builds from the frontend's perspective -- it starts with the user-facing feature and works backward to the data layer.
