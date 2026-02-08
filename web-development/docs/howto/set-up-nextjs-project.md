---
title: "Set Up a Next.js Project"
description: "Use nextjs-developer to configure App Router, layouts, and rendering strategies"
weight: 2
---

# Set Up a Next.js Project

Configure a Next.js project with correct App Router architecture, layout hierarchy, and rendering strategies using the nextjs-developer agent.

## Goal

Produce a Next.js project structure with properly organized routes, shared layouts, loading and error boundaries, and appropriate rendering strategies (static, ISR, dynamic, or streaming) for each route.

## Prerequisites

- A Next.js 14+ project (or intent to create one)
- The web-development plugin installed

## Steps

### 1. Describe Your Project

Provide the nextjs-developer with your project's requirements:

- **Pages and routes** -- list the main pages (home, dashboard, settings, product detail, etc.)
- **Data sources** -- where data comes from (database, CMS, external API) and how fresh it needs to be
- **Authentication** -- which routes are public and which require login
- **Multi-tenant or single-tenant** -- whether the app serves multiple organizations with separate data

Be explicit about data freshness. "Product pages" could mean statically generated at build time (catalog that changes weekly) or dynamically rendered per request (inventory that changes every minute). The rendering strategy depends on this.

### 2. Trigger nextjs-developer for Route Structure Design

The agent activates on Next.js-specific tasks. Natural trigger phrases include:

- "Set up the route structure for..."
- "Design the App Router layout for..."
- "Should this page use static generation or ISR?"
- "Convert this form to server actions"
- "Product pages are slow"

Describe your project:

```
Set up the App Router structure for a SaaS app with public marketing pages,
authenticated dashboard, team settings, and a billing page.
```

### 3. Review the App Router Layout Tree

The nextjs-developer produces a route structure with:

- **Route groups** -- `(marketing)` for public pages, `(app)` for authenticated pages, each with their own root layout
- **Shared layouts** -- a dashboard layout with sidebar navigation that wraps all authenticated routes
- **Dynamic segments** -- `[teamId]` for multi-tenant routes, `[slug]` for content pages
- **Parallel routes** -- independent panels that load and error independently (e.g., a dashboard with a stats panel and an activity feed)
- **Loading and error boundaries** -- `loading.tsx` at each significant route level for streaming, `error.tsx` for graceful error recovery

Verify the structure includes:

- `layout.tsx` files at appropriate nesting levels (not just one root layout)
- `loading.tsx` for routes that fetch data
- `error.tsx` for routes where errors should be caught without crashing the entire page
- `not-found.tsx` for dynamic routes that may resolve to missing content

### 4. Configure Rendering Strategies Per Route

The nextjs-developer recommends a rendering strategy for each route based on your data freshness requirements:

- **Static** (`generateStaticParams`) -- content that changes rarely (marketing pages, documentation)
- **ISR** (`revalidate: 3600`) -- content that updates periodically (product catalog, blog posts)
- **Dynamic** (`force-dynamic`) -- content that must be fresh on every request (user dashboard, real-time data)
- **Streaming** (Suspense boundaries with `loading.tsx`) -- pages with a fast shell and slower data sections

Each route gets the least dynamic strategy that meets its freshness needs. Static is the default; dynamic is the last resort.

## Verification

After completing these steps you should have:

- [ ] A route structure that reflects your application's page hierarchy
- [ ] Route groups separating public and authenticated sections
- [ ] Layouts that share navigation and chrome without re-rendering on route changes
- [ ] Loading boundaries at every route that fetches data
- [ ] Error boundaries at every route where failures should be contained
- [ ] Rendering strategies documented per route, with the reasoning for each choice
- [ ] Server Components by default -- `"use client"` only on interactive leaf components

## Troubleshooting

**Client vs. Server Component boundaries are confusing.**
The rule is simple: Server Components by default. Add `"use client"` only when the component needs `useState`, `useEffect`, event handlers, or browser APIs. If you are unsure, ask the nextjs-developer to analyze a specific component and recommend the correct boundary.

**A page uses `getServerSideProps` or `getStaticProps`.**
Those are Pages Router APIs. In the App Router, data fetching happens directly in Server Components using `async` functions. Ask the nextjs-developer to migrate the page to the App Router pattern.

**The feature requires API routes, database schema, and frontend.**
The nextjs-developer focuses on the Next.js layer -- routing, rendering, and Server Components. If the feature spans the full stack (UI + API + database), use the fullstack-developer agent instead. It builds all layers together and uses Next.js as the frontend where appropriate.
