---
title: "Choosing the Right Agent"
description: "When to use each web development agent"
weight: 2
---

# Choosing the Right Agent

How to select the right web-development agent for a task, handle overlapping responsibilities, and combine agents for complex workflows.

## Decision Matrix

| Task | Agent |
|------|-------|
| Building a UI component (any framework) | frontend-developer |
| React-specific component architecture, hooks, or state management | react-specialist |
| Vue 3 components, composables, Pinia, or Nuxt 3 | vue-expert |
| Next.js routing, rendering strategies, or Server Actions | nextjs-developer |
| End-to-end feature spanning UI + API + database | fullstack-developer |
| Design system, UX flow, visual specifications | ui-designer |

A more detailed breakdown by situation:

- **"I need a responsive layout"** -- frontend-developer. It handles responsive design across any framework and works with the project's existing styling system.
- **"This React component re-renders too much"** -- react-specialist. Re-render debugging requires understanding React's reconciliation model, memo boundaries, and state placement.
- **"Migrate our Vuex store to Pinia"** -- vue-expert. This is a Vue-ecosystem-specific task involving store decomposition, `storeToRefs`, and Pinia's API patterns.
- **"Set up ISR for our product pages"** -- nextjs-developer. Rendering strategy is a Next.js-specific decision involving `revalidate` configuration, cache behavior, and `generateStaticParams`.
- **"Build user registration with email verification"** -- fullstack-developer. Registration spans UI, API, database, and email -- it cannot be built in one layer alone.
- **"Users are dropping off at checkout"** -- ui-designer. Conversion analysis requires user journey mapping, friction identification, and UX redesign.

## The Overlap Question

**frontend-developer vs. react-specialist** -- The frontend-developer handles general UI implementation: responsive layouts, component composition, state management patterns that apply across frameworks. The react-specialist handles React-specific concerns: hook composition, Server/Client Component boundaries, React Query integration, React DevTools profiling. If the task is "build a card component," either agent works. If the task is "why does this component re-render when the parent's state changes," the react-specialist's React-specific knowledge produces a better answer.

The same logic applies to **frontend-developer vs. vue-expert**. General component work can go to the frontend-developer. Vue-specific concerns -- Composition API migration, `ref` vs. `reactive` decisions, Nuxt 3 configuration -- belong to the vue-expert.

**nextjs-developer vs. react-specialist** -- Both work with React, but at different levels. The react-specialist focuses on component-level concerns (hooks, state, rendering). The nextjs-developer focuses on application-level concerns (routing, rendering strategies, Server Actions, caching). If the task is about a component's internal architecture, use the react-specialist. If the task is about where data fetching happens, which routes are static, or how to structure layouts, use the nextjs-developer.

**nextjs-developer vs. fullstack-developer** -- The nextjs-developer works within the Next.js layer. It designs routes, configures rendering, and implements Server Actions. The fullstack-developer works across all layers. If the task requires database schema design, API route implementation, and frontend wiring as a single unit, the fullstack-developer handles it with Next.js as the frontend framework. If the task is purely about Next.js architecture (route groups, parallel routes, streaming), the nextjs-developer is more focused.

## Combining Agents

For complex features, agents can be used in sequence:

1. **ui-designer** first -- establish the design: user journey, component specifications, design tokens, interaction states
2. **react-specialist** or **vue-expert** second -- implement the components following the design specifications, with correct architecture and state management
3. **fullstack-developer** third -- wire the frontend to API routes and database operations, ensuring end-to-end data flow

This pipeline moves from design to implementation to integration. Each agent builds on the artifacts the previous one produced.

A simpler two-agent combination is common:

- **ui-designer + react-specialist** -- design and implement a component with both visual quality and architectural soundness
- **nextjs-developer + fullstack-developer** -- set up the Next.js route structure, then build the full-stack features within it

## Cross-Plugin Boundaries

Some tasks appear web-related but belong to other plugins:

- **Writing tests for components** -- use the code-quality plugin. Web-development agents build features; code-quality agents verify them.
- **Designing a REST API that the frontend will consume** -- use the backend-development plugin for API design, then the fullstack-developer or frontend-developer to build the frontend that calls it.
- **Setting up CI/CD for a Next.js app** -- use the cli-development or DevOps-related plugin. The nextjs-developer builds the application; deployment pipeline configuration is a separate concern.
