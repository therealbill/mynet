---
title: "Getting Started with Web Development"
description: "Build a React component with react-specialist and style it with ui-designer guidance"
weight: 1
---

# Getting Started with Web Development

Walk through building a user profile card component from scratch using two complementary agents: react-specialist for component architecture and ui-designer for visual design guidance.

## What You'll Learn

- Use react-specialist to design a component's props interface, state management, and TypeScript types
- Use ui-designer to establish design tokens, spacing, and interaction states
- Understand how the two agents complement each other in a typical frontend workflow

## Prerequisites

- A React project with TypeScript configured (Create React App, Vite, or Next.js all work)
- The web-development plugin installed in your Claude Code environment

## Step 1: Describe the Component You Need

Start with a plain description of what you want to build. Be specific about the data it needs to display and the interactions it should support:

```
I need a user profile card that shows an avatar, name, role, and a short bio.
It should have an "Edit" button that reveals inline editing fields.
The card appears in a sidebar, so it needs to be compact.
```

This description gives the agents enough context to make architectural and design decisions. Mention the layout context (sidebar, grid, modal) because it affects sizing and responsive behavior.

## Step 2: Trigger react-specialist for Component Architecture

The react-specialist agent activates when the task involves React-specific component design. Phrases like "break it up," "component structure," or "props interface" are natural triggers. For this tutorial, describing a React component is enough:

```
Design the React component architecture for this profile card.
```

The react-specialist analyzes your requirements and produces:

- A **props interface** with TypeScript types -- `UserProfileCardProps` with fields for `avatar`, `name`, `role`, `bio`, and callback handlers
- A **state strategy** -- local `useState` for the edit mode toggle, since no sibling component needs this state
- A **component decomposition** -- the card split into `UserAvatar`, `UserInfo`, and `InlineEditForm` sub-components, each with a focused responsibility
- **Data state handling** -- loading skeleton, error fallback, and empty/missing-avatar states alongside the happy path

### Checkpoint

Review the component structure. You should see:

- A clear props interface with narrow, well-typed props
- Local state for UI concerns (edit mode), not lifted unnecessarily
- Sub-components that are small enough to understand at a glance
- All four data states accounted for: loading, error, empty, success

## Step 3: Review the Component in Detail

Examine what react-specialist produced. The component follows several React patterns worth noting:

- **Composition over configuration** -- instead of one large component with boolean props like `isEditing`, `showAvatar`, `showBio`, the card composes smaller pieces. The `InlineEditForm` is a separate component that replaces `UserInfo` when editing, rather than a conditional branch buried inside a monolith.
- **TypeScript discriminated unions** -- if the card has variants (e.g., compact vs. expanded), the props use a discriminated union rather than optional fields that create impossible states.
- **Server Component compatibility** -- the outer card is a Server Component by default. Only the edit button and form carry the `"use client"` directive, keeping the JavaScript bundle small.

If anything looks wrong -- too many props, state in the wrong place, missing error handling -- ask the react-specialist to iterate:

```
The bio field can be very long. Add truncation with a "show more" toggle.
```

## Step 4: Ask ui-designer for Visual Design Guidance

Now shift from architecture to appearance. The ui-designer agent activates on design-related requests:

```
Design the visual style for this profile card -- spacing, colors, and interaction states.
```

The ui-designer provides:

- **Design tokens** -- a spacing scale (4px grid), color palette (neutral scale plus one accent color), border radius, and shadow values expressed as Tailwind classes or CSS custom properties
- **Layout specifications** -- the avatar size, text hierarchy (name as `text-lg font-semibold`, role as `text-sm text-muted-foreground`), and internal padding
- **Interaction states** -- hover effect on the edit button, focus ring for keyboard navigation, disabled state while the edit form submits
- **Responsive behavior** -- how the card adapts if the sidebar collapses or the viewport shrinks
- **Edge cases** -- what happens with a missing avatar (initials fallback), a very long name (truncation with title attribute), or an empty bio (placeholder text)

### Checkpoint

Verify the design covers:

- Default, hover, focus, active, disabled, loading, and error states for all interactive elements
- Accessible color contrast ratios (WCAG 2.1 AA minimum)
- A spacing system based on consistent units, not arbitrary pixel values

## Step 5: See the Complete Styled Component

Combine the architecture from react-specialist with the design from ui-designer. The result is a component that is:

- **Architecturally sound** -- small sub-components, typed props, correct state placement
- **Visually consistent** -- design tokens ensure the card fits into any existing design system
- **Accessible** -- semantic HTML, keyboard navigation, focus management on the edit toggle
- **Complete** -- all data states and interaction states specified and implemented

The final component includes the `UserProfileCard` container, `UserAvatar` with initials fallback, `UserInfo` with truncation, and `InlineEditForm` with validation and submit states.

## What You Learned

- **react-specialist** handles the structural concerns: how to decompose a component, where state belongs, what the props interface looks like, and how to handle data states. It thinks in React patterns.
- **ui-designer** handles the visual concerns: spacing, color, typography, interaction states, and accessibility. It thinks in design systems.

The two agents address different dimensions of the same component. Using them together produces a result that is neither architecturally messy with good visuals nor well-structured with no design consideration.

## Next Steps

- [Build a React Component]({{< relref "../howto/build-react-component" >}}) -- detailed how-to for component creation workflows
- [Set Up a Next.js Project]({{< relref "../howto/set-up-nextjs-project" >}}) -- if your project uses Next.js
- [Agent Reference]({{< relref "../reference/agents" >}}) -- full specifications for all 6 web-development agents
- [Choosing the Right Agent]({{< relref "../explanation/choosing-the-right-agent" >}}) -- when to use which agent
