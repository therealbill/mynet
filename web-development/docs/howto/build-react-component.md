---
title: "Build a React Component"
description: "Use react-specialist to design and implement React components"
weight: 1
---

# Build a React Component

Create a well-architected React component using the react-specialist agent, from initial description through implementation and iteration.

## Goal

Produce a React component with a clean props interface, correct state management, TypeScript types, and coverage for all data states (loading, error, empty, success).

## Prerequisites

- A React project with TypeScript
- The web-development plugin installed

## Steps

### 1. Describe the Component's Purpose and Data Requirements

Provide the react-specialist with a clear description that includes:

- **What the component displays** -- the data it renders and the layout context it appears in
- **What interactions it supports** -- clicks, form input, drag-and-drop, keyboard shortcuts
- **What data it consumes** -- props from a parent, API responses, local state

Be specific about the data shape. "A table of users" is less useful than "a sortable table of users with columns for name, email, role, and last active date, with row selection for bulk actions."

### 2. Trigger react-specialist

The agent activates on React-specific tasks. Natural trigger phrases include:

- "Design the component architecture for..."
- "Break this component into smaller pieces"
- "This component re-renders too often"
- "Should I use useState or useReducer here?"
- "Help me choose between Zustand and React Query"

Describe what you need in plain language:

```
Design a sortable data table component for displaying user records with bulk selection.
```

### 3. Review the Component Structure

The react-specialist produces a component with:

- **Props interface** -- TypeScript types with narrow, well-named props. Discriminated unions for component variants rather than optional boolean flags.
- **State strategy** -- local `useState` for UI state (sort direction, selected rows), `useReducer` for complex state transitions, React Query for server data.
- **Hook composition** -- custom hooks like `useSortableData` or `useRowSelection` that encapsulate logic and can be tested independently.
- **Sub-component breakdown** -- `TableHeader`, `TableRow`, `SelectionToolbar` as separate components with their own props interfaces.

Verify that the structure follows these patterns:

- State lives where it is consumed, not lifted unnecessarily
- Server data uses React Query or Server Components, not manual `useEffect` + `useState`
- The component handles loading, error, empty, and success states

### 4. Iterate if Needed

If the component needs refinement, ask the react-specialist directly:

- **Decomposition**: "This component is still too large -- split the filtering logic into a separate hook."
- **Performance**: "The table re-renders on every keystroke in the search box -- fix the re-render boundary."
- **State management**: "Multiple components need access to the selection state -- should I lift it or use a store?"

The agent adjusts the architecture based on your feedback, explaining the tradeoff behind each decision.

## Verification

After completing these steps you should have:

- [ ] A TypeScript props interface with no `any` types
- [ ] State management appropriate to the complexity (useState, useReducer, or external store)
- [ ] All four data states rendered: loading skeleton, error message, empty state, and populated view
- [ ] Custom hooks extracted for reusable logic
- [ ] Sub-components small enough to understand at a glance

## Troubleshooting

**The component re-renders too often.**
Ask the react-specialist to profile the component. Common causes: creating new object/array references on every render (move them outside the component or wrap in `useMemo`), passing unstable callbacks (wrap in `useCallback`), or placing state too high in the tree.

**Props are being drilled through multiple layers.**
This signals a composition problem, not a state management problem. Ask the react-specialist to restructure using the compound component pattern or React context scoped to the component subtree.

**The task requires API routes or database changes.**
The react-specialist focuses on the React layer. If the feature spans the full stack -- UI, API, and database -- use the fullstack-developer agent instead. It builds all layers together so the contracts stay aligned.
