---
title: "Build an iOS Feature"
description: "Use ios-developer to implement SwiftUI views with proper state management"
weight: 1
---

# Build an iOS Feature

Use the ios-developer agent to implement a complete SwiftUI feature with proper state management, data flow, error handling, and accessibility.

## Goal

Build a fully functional iOS feature that follows Apple's conventions, uses the Observation framework for state management, and handles loading, error, and empty states correctly.

## Step 1: Describe the Feature Requirements

Tell Claude Code what you want to build, including:

- **Data source:** Where the data comes from (REST API, local database, or both)
- **Interaction model:** What the user does with the feature (browse, create, edit, delete)
- **Navigation:** How this feature fits into the app's navigation hierarchy (tab, push, modal)

```
Build a settings screen that lets users edit their profile.
It reads from and writes to our REST API.
It should be pushed from the profile tab.
```

ios-developer activates when you describe SwiftUI views, iOS features, or data fetching patterns.

## Step 2: Review the Project Analysis

ios-developer reads your project before writing code. It examines:

- `Package.swift` or `.xcodeproj` for dependencies and deployment target
- Existing view files for naming conventions and patterns
- Data models for `Codable` and persistence patterns
- Navigation structure for how views are presented

If the agent asks clarifying questions -- answer them. Questions about deployment target, minimum iOS version, or existing patterns help the agent match your project conventions.

## Step 3: Review the SwiftUI View Hierarchy

After implementation, review what the agent generated. Check for:

- **View composition:** Each view file should be under 200 lines. Larger views should be split into extracted subviews.
- **State management:** `@State` for view-local state, `@Observable` view models for shared state, `@Environment` for app-wide dependencies. No `@EnvironmentObject` in new code.
- **Data flow:** The view model owns the data. Views read from the view model and call methods on it. Views do not contain business logic.
- **Error handling:** Domain-specific error enums, not bare `String` messages or generic `Error` catches.

```
Show me the view hierarchy and explain the state management choices
```

## Step 4: Review Accessibility Support

ios-developer includes accessibility by default. Verify:

- Interactive elements have `accessibilityLabel` values that describe their purpose
- `accessibilityHint` is present where the action is not obvious from the label
- Proper trait assignments (`.isButton`, `.isHeader`, `.isSelected`) on custom components
- Dynamic type support -- text scales with the user's preferred size

```
Does this feature pass a VoiceOver audit?
```

## Verification

The feature is complete when:

- `xcodebuild` compiles without errors or warnings
- The view handles loading state (spinner or skeleton), error state (message with retry action), and empty state (informative message)
- Navigation integrates with the existing app structure
- Accessibility labels are present on all interactive elements
- State management uses `@Observable` (not `ObservableObject`) for iOS 17+ targets

## Troubleshooting

### State management causes unexpected re-renders

If views re-render more than expected, the observation boundaries may be wrong. Common causes:

- An `@Observable` model that is too broad -- split it into focused models per feature
- Reading computed properties in the view body that depend on unrelated state changes
- Using `@State` for state that should be in a view model

Ask swift-expert to review the observation boundaries:

```
My settings view re-renders when unrelated state changes -- review the observation model
```

### Architecture decisions beyond ios-developer's scope

ios-developer makes implementation decisions, but architectural questions benefit from swift-expert:

- Whether to use actors for shared mutable state
- Structured vs unstructured concurrency patterns
- SwiftUI vs UIKit for a specific component
- Protocol design and type system decisions

```
Should this view model use an actor instead of @MainActor?
```

## See Also

- [Getting Started](../../tutorials/getting-started/) -- end-to-end tutorial building and reviewing an iOS feature
- [Agent Reference](../../reference/agents/) -- ios-developer specification
- [Architecture](../../explanation/architecture/) -- when to use ios-developer vs swift-expert
