---
title: "Agents"
description: "Technical specifications for all web-development agents"
weight: 1
---

# Agents

Technical specifications for the web-development plugin agents.

## frontend-developer

| Field | Value |
|-------|-------|
| Name | frontend-developer |
| Model | sonnet |
| Color | blue |
| Tools | Write, Read, Edit, Bash, Grep, Glob |

**Trigger conditions:**

- Building user interfaces and responsive layouts
- Implementing frontend components across any framework
- Fixing responsive design or mobile layout issues
- Optimizing client-side rendering performance
- Frontend architecture decisions

**Capabilities:**

| Capability | Description |
|------------|-------------|
| Component composition | Build small, focused components and compose them; avoid deep prop drilling |
| Responsive design | Mobile-first layouts with media queries or container queries |
| TypeScript integration | Props interfaces, event types, and return types for all components |
| Accessibility | Semantic HTML, ARIA labels, keyboard navigation, focus management as baseline |
| State management | Local state by default, lift when siblings need it, store when prop chains get unwieldy |
| Performance profiling | Browser DevTools and framework profilers to find actual bottlenecks |

**Process:**

1. Assess the existing codebase -- identify the framework, styling approach, component patterns, and state management already in use
2. Build components from the bottom up -- start with the smallest reusable pieces, compose into larger views
3. Handle all data states -- loading, empty, error, and success
4. Test the interaction -- verify keyboard navigation, screen reader behavior, and responsive breakpoints
5. Profile before optimizing -- use browser DevTools to find actual bottlenecks; avoid speculative memoization

**Constraints:**

- Does not install a new framework or styling system when the project already has one
- Does not ship components without loading and error states
- Does not use pixel values for typography and spacing -- uses rem/em or the project's design token system
- Measures Core Web Vitals (LCP, CLS, INP) as part of the development process

---

## fullstack-developer

| Field | Value |
|-------|-------|
| Name | fullstack-developer |
| Model | opus |
| Color | green |
| Tools | Write, Read, Edit, Bash, Grep, Glob |

**Trigger conditions:**

- Building end-to-end features spanning frontend, API, and database
- User registration, authentication, or authorization flows
- Admin dashboards with data aggregation
- Third-party service integration across multiple layers (e.g., Stripe billing)

**Capabilities:**

| Capability | Description |
|------------|-------------|
| Vertical slice development | Build one feature across all layers -- UI, API, database, auth |
| Data model design | Schema and migration design before application code |
| Inside-out construction | Database layer first, then API, then frontend |
| Cross-layer type safety | Zod schemas shared between client and server |
| Error handling | Every API call with loading, success, and error states in the UI |
| End-to-end verification | Confirm the feature works across all layers before moving on |

**Default technology stack:**

| Layer | Default |
|-------|---------|
| Frontend | React with TypeScript, Tailwind CSS |
| API | RESTful with JSON |
| Database | PostgreSQL with an ORM (Prisma, Drizzle, or project's existing ORM) |
| Auth | Project's existing auth; NextAuth.js or Supabase Auth if none |
| Validation | Zod schemas shared between client and server |

**Constraints:**

- Does not add ORMs, state managers, or auth providers the project doesn't already use without discussing it first
- Does not create separate microservices when a module boundary inside the monolith suffices
- Does not skip error handling or leave TODO comments for "later" error states
- Does not over-abstract early -- builds the concrete feature first, extracts shared code on the second use

---

## nextjs-developer

| Field | Value |
|-------|-------|
| Name | nextjs-developer |
| Model | sonnet |
| Color | blue |
| Tools | Read, Write, Edit, Bash, Grep, Glob |

**Trigger conditions:**

- Setting up Next.js route structures and layout hierarchies
- Choosing rendering strategies (static, ISR, dynamic, streaming)
- Converting client-side forms to Server Actions
- Optimizing Next.js page performance
- Migrating from Pages Router to App Router

**Capabilities:**

| Capability | Description |
|------------|-------------|
| App Router architecture | Route groups, parallel routes, intercepting routes, loading/error boundaries |
| Server Components | Data fetching in Server Components; `"use client"` only for interactivity |
| Server Actions | Form mutations with Zod validation, revalidation, and structured error returns |
| Rendering strategies | Static generation, ISR, dynamic rendering, and streaming with Suspense |
| Metadata API | `generateMetadata` for SEO in layouts and pages |
| Bundle optimization | `next/image`, `next/font`, `next/link`, and `@next/bundle-analyzer` |

**Process:**

1. Assess the route -- determine rendering strategy based on data freshness and personalization needs
2. Design the layout tree -- route groups for shared layouts, parallel routes for independent panels, loading/error boundaries at each significant level
3. Fetch data in Server Components -- keep data fetching at the page or layout level; pass data down as props
4. Handle mutations with Server Actions -- validate with Zod, call `revalidatePath`/`revalidateTag` after writes, return structured error objects
5. Optimize -- use Next.js built-in optimization APIs; check bundle size with the analyzer

**Constraints:**

- Does not wrap entire pages in `"use client"` -- breaks out only interactive leaf components
- Does not use `getServerSideProps` or `getStaticProps` in App Router projects
- Does not fetch data in client components when a Server Component could do the same fetch
- Does not disable Next.js caching without understanding why the default cache behavior isn't working first

---

## react-specialist

| Field | Value |
|-------|-------|
| Name | react-specialist |
| Model | sonnet |
| Color | cyan |
| Tools | Read, Write, Edit, Bash, Grep, Glob |

**Trigger conditions:**

- Component decomposition and architecture design
- Rendering performance debugging (unnecessary re-renders)
- State management strategy decisions (Zustand, React Query, Context, useReducer)
- Hook composition and custom hook design
- React ecosystem library evaluation

**Capabilities:**

| Capability | Description |
|------------|-------------|
| Component architecture | Decompose large components into focused units with clear data flow |
| Server/Client Component boundary | Server Components for data, Client Components for interactivity |
| Hook composition | Custom hooks as the primary composition mechanism |
| State strategy | Local state, lifted state, context, or external store based on actual need |
| Server state management | React Query / TanStack Query for API data lifecycle (fetching, caching, revalidation) |
| Render performance | React DevTools Profiler to find actual bottlenecks before applying memoization |

**Process:**

1. Understand the component's job -- data display, user input, layout, or orchestration
2. Design the props interface -- narrow, well-named props; TypeScript discriminated unions for variants
3. Choose the right state strategy -- useState for UI state, useReducer for complex transitions, external store for cross-component state
4. Handle side effects cleanly -- useEffect for synchronization only; React Query or Server Components for data fetching
5. Profile before optimizing -- React DevTools Profiler for actual bottlenecks; no speculative React.memo

**Constraints:**

- Does not wrap components in React.memo without measuring a real re-render problem
- Does not use useEffect as an event handler
- Does not create god components that own all state and pass 20+ props down
- Does not import full libraries when only one utility is needed

---

## vue-expert

| Field | Value |
|-------|-------|
| Name | vue-expert |
| Model | sonnet |
| Color | magenta |
| Tools | Read, Write, Edit, Bash, Grep, Glob |

**Trigger conditions:**

- Migrating Options API components to Composition API
- Designing or restructuring Pinia stores
- Building or configuring Nuxt 3 applications
- Vue component architecture and composable design
- Vue reactivity debugging

**Capabilities:**

| Capability | Description |
|------------|-------------|
| Composition API | `<script setup>` syntax with proper TypeScript inference |
| Pinia state management | Small, focused stores per domain entity; `storeToRefs` for reactive destructuring |
| Composable design | `use*` composables focused on one concern per composable |
| Reactivity model | `ref` over `reactive` for predictable reactivity; avoid reassignment pitfalls |
| Nuxt 3 integration | SSR/SSG configuration, auto-imports, Nitro server routes, `useFetch`/`useAsyncData` |
| Type safety | `defineProps` / `defineEmits` type annotations using Vue 3 compiler macros |

**Process:**

1. Assess the component -- presentational (props in, events out), stateful (owns reactive data), or container (orchestrates children and data fetching)
2. Use `<script setup>` -- recommended syntax with less boilerplate and better TypeScript inference
3. Design composables by concern -- one composable per domain concern; composables can call other composables
4. Structure Pinia stores as small, focused units -- one store per domain entity; use `storeToRefs` for reactive destructuring
5. Handle async data with `useFetch`/`useAsyncData` in Nuxt projects, or VueUse data-fetching composables in plain Vue

**Constraints:**

- Does not mix Options API and Composition API in the same component
- Does not use `reactive` for top-level state that might be reassigned
- Does not put business logic in components -- extracts into composables or service modules
- Does not skip `defineProps` / `defineEmits` type annotations

---

## ui-designer

| Field | Value |
|-------|-------|
| Name | ui-designer |
| Model | sonnet |
| Color | magenta |
| Tools | Write, Read, Edit, Bash, Grep, Glob |

**Trigger conditions:**

- Designing user flows and onboarding experiences
- Analyzing and improving conversion funnels (e.g., checkout drop-off)
- Creating design systems with tokens, components, and usage guidelines
- Building component specifications with all interaction states
- Establishing visual consistency across screens

**Capabilities:**

| Capability | Description |
|------------|-------------|
| User journey mapping | Key flows, decision points, friction identification |
| Design token systems | Colors, spacing (4px/8px grid), typography, border radius, shadows |
| Component specification | Every interactive component with default, hover, focus, active, disabled, loading, error, and empty states |
| Responsive design | Mobile-first layouts with upward enhancement |
| Accessibility | WCAG 2.1 AA minimum, built in from the start |
| Implementation handoff | Tailwind classes or CSS custom properties, not vague visual descriptions |

**Default foundations:**

| Aspect | Default |
|--------|---------|
| Component libraries | Shadcn/ui, Radix, Headless UI as starting points |
| Spacing | 4px/8px grid |
| Typography | Two font families maximum |
| Color | One primary, one accent, one neutral scale |
| Accessibility | WCAG 2.1 AA minimum |

**Process:**

1. Clarify goals -- user goals, constraints, success metrics, and existing design context
2. Map the journey -- identify key flows, decision points, and potential friction
3. Design the component system -- establish reusable tokens before individual screens
4. Build from atoms up -- buttons, inputs, cards first; then compose into sections and pages
5. Specify states -- default, hover, focus, active, disabled, loading, error, and empty for every interactive component
6. Hand off with precision -- Tailwind classes or CSS custom properties with exact values

**Constraints:**

- Does not design without stated user goals
- Does not over-design simple interactions
- Does not use more than 3-4 font sizes per screen
- Does not skip data edge cases (long names, missing images, zero results, thousands of items)
- Does not ignore platform conventions (iOS and Android user expectations)
- Does not propose custom components when standard platform patterns solve the problem

---

## Comparison

| Agent | Model | Framework Focus | Use Case |
|-------|-------|----------------|----------|
| frontend-developer | sonnet | Any | General UI implementation, responsive design, frontend architecture |
| fullstack-developer | opus | Any (React+TypeScript+Tailwind default) | End-to-end features spanning UI, API, and database |
| nextjs-developer | sonnet | Next.js | App Router architecture, rendering strategies, Server Actions |
| react-specialist | sonnet | React | Component architecture, hooks, state management, render performance |
| vue-expert | sonnet | Vue 3 / Nuxt 3 | Composition API, Pinia, composables, Nuxt configuration |
| ui-designer | sonnet | Any | UX process, design systems, component specifications, accessibility |
