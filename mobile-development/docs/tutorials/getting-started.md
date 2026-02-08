---
title: "Getting Started with Mobile Development"
description: "Build an iOS feature with ios-developer and refine Swift code with swift-expert"
weight: 1
---

# Getting Started with Mobile Development

Build a product list feature for an iOS app using the ios-developer agent, then bring in swift-expert to review the concurrency model and architectural decisions.

## What You'll Build

By the end of this tutorial, you will have:

- Triggered the ios-developer agent with a feature description
- Reviewed a SwiftUI view with an `@Observable` view model and async data loading
- Asked swift-expert to audit the concurrency and architecture
- Understood which agent to use for building versus reviewing Swift code

This takes about 15 minutes.

## Prerequisites

- Claude Code CLI installed and authenticated (run `claude --version` to verify)
- The mobile-development plugin installed in your project's `.claude/settings.json`
- Xcode 15+ installed with iOS 17+ SDK (`xcodebuild -version` to verify)
- Basic familiarity with Swift syntax and SwiftUI concepts

## Step 1: Describe the Feature

Open Claude Code in your iOS project directory and describe the feature you want to build:

```
Build me a SwiftUI view that fetches and displays a list of products from our REST API
```

Claude Code matches your request to the ios-developer agent based on the mention of SwiftUI and data fetching. The agent activates and begins planning the implementation.

## Step 2: Review the Agent's Approach

ios-developer reads your project structure first -- it looks at `Package.swift` or `.xcodeproj`, existing view patterns, and your data model conventions. Then it implements the feature using SwiftUI and platform-appropriate patterns.

For a product list, expect the agent to generate:

- **View model:** An `@Observable` class with `products` array, `isLoading` flag, and `error` state
- **Data loading:** An `async` method wrapping `URLSession` with `Codable` model decoding
- **SwiftUI view:** A `List` with loading, error, and empty states handled explicitly
- **Error types:** A domain-specific error enum conforming to `LocalizedError`

The agent also adds accessibility support -- `accessibilityLabel` on product cells and proper trait assignments for interactive elements.

### Checkpoint

At this point you should have:

- A `ProductListView.swift` with a SwiftUI view
- A `ProductViewModel.swift` with an `@Observable` view model
- A `Product.swift` model conforming to `Codable`
- Error handling with specific error types, not generic `Error` catches

If ios-developer did not activate, verify the mobile-development plugin is listed in your `.claude/settings.json` and restart Claude Code.

## Step 3: Build and Verify

The agent builds the project using `xcodebuild` or `swift build` to verify compilation. Check that:

- The project compiles without warnings
- The view handles all three states: loading spinner, error message with retry, and populated list
- Accessibility labels are present on interactive elements

```
Build the project and show me any warnings
```

## Step 4: Ask swift-expert to Review

Now bring in the architectural reviewer. Ask a question that triggers swift-expert:

```
Review the concurrency model in ProductViewModel -- should I use structured concurrency differently?
```

swift-expert activates because the question is about concurrency architecture, not feature implementation. This agent uses the opus model, which provides deeper architectural reasoning than ios-developer's sonnet model.

## Step 5: Review the Architectural Recommendations

swift-expert reads the generated code and makes specific recommendations. For a product list view model, expect feedback like:

- Whether `Task { }` in `.task` is the right pattern or if `async let` would be better for parallel loads
- Whether the view model needs `@MainActor` isolation or if `@Observable` handles it
- Whether the error type design is sufficient or should use typed throws
- Whether `AsyncSequence` would be better than a one-shot fetch for data that updates

The agent makes a clear recommendation with trade-offs stated -- not a menu of options. It implements the changes directly if you agree.

### Checkpoint

At this point you should have:

- A working product list feature built by ios-developer
- Architectural feedback from swift-expert on the concurrency model
- An understanding of why certain patterns were chosen over alternatives

## What You Learned

In this tutorial, you:

- **Triggered ios-developer** by describing a SwiftUI feature -- the agent built a complete view with view model, data loading, error handling, and accessibility
- **Observed the agent's process** -- it read your project structure first, then implemented using your existing conventions
- **Triggered swift-expert** by asking an architecture question -- this agent reviewed the code at a higher level than implementation
- **Saw the scope difference** between the two agents: ios-developer builds features, swift-expert reviews architectural decisions

## Next Steps

- [Build an iOS Feature]({{< ref "howto/build-ios-feature" >}}) -- detailed guide for implementing SwiftUI features with proper state management
- [Choose Mobile Platform]({{< ref "howto/choose-mobile-platform" >}}) -- use mobile-developer to decide between native and cross-platform
- [Agent Reference]({{< ref "reference/agents" >}}) -- full specifications for all three mobile agents
- [Architecture]({{< ref "explanation/architecture" >}}) -- understand when to use each agent

### Cross-Plugin References

- **programming-languages plugin** -- use swift-pro for Swift code quality and simplification tasks
- **code-quality plugin** -- use code-reviewer for general code review beyond Swift-specific architecture
