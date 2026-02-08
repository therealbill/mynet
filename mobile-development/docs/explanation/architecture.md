---
title: "Architecture"
description: "Platform specialization and when to use each mobile agent"
weight: 1
---

# Architecture

The mobile-development plugin provides three agents that cover different layers of mobile development: feature implementation, platform strategy, and architectural review. This page explains how they relate, where they overlap, and when to use each one.

## Platform Specialization

Each agent owns a distinct concern in mobile development:

- **ios-developer** handles native iOS feature implementation. It writes SwiftUI views, configures state management with `@Observable`, integrates data persistence with SwiftData or Core Data, and builds with `xcodebuild`. Its scope is a single platform: iOS.

- **mobile-developer** handles cross-platform strategy. It decides whether you should build native or use React Native or Flutter, designs code sharing strategies, and implements cross-platform features targeting both iOS and Android. Its scope spans platforms but stays at the implementation level.

- **swift-expert** handles Swift architecture and language design. It reviews concurrency models, makes framework selection decisions (SwiftUI vs UIKit), designs type hierarchies, and evaluates server-side Swift trade-offs. Its scope is the Swift language and ecosystem, regardless of whether the code runs on iOS, macOS, or a server.

## The Scope Overlap

ios-developer and swift-expert both work with Swift code, but at different levels of concern.

ios-developer writes Swift code to implement features. When it creates a view model with `@Observable` and loads data with `async/await`, it follows established patterns without debating whether those patterns are the right choice. It builds the feature.

swift-expert evaluates whether those patterns are the right choice. When reviewing the same view model, it considers whether structured concurrency (`async let`, `TaskGroup`) would be better than unstructured `Task { }`, whether the view model needs actor isolation, and whether the error type design is robust enough. It reviews the architecture.

This is intentional. Feature implementation benefits from a focused agent that follows conventions and produces working code quickly. Architectural review benefits from a separate agent that questions assumptions and evaluates trade-offs with deeper reasoning.

## When to Use Each Agent

### Building a feature: ios-developer

Use ios-developer when you have a specific feature to implement on iOS. Examples:

- "Build a settings screen with profile editing"
- "Add pull-to-refresh to the product list"
- "Implement background data sync with BGTaskScheduler"

ios-developer reads your project, implements the feature with SwiftUI, includes error handling and accessibility, and builds to verify.

### Choosing a platform: mobile-developer

Use mobile-developer when you need to decide how to build a mobile app or when implementing cross-platform features. Examples:

- "Should we use React Native or go native for this project?"
- "We need iOS and Android with 80% feature parity"
- "Set up offline-first data sync for our React Native app"

mobile-developer evaluates your requirements, recommends a platform, and implements cross-platform code with clean platform splits.

### Making architecture decisions: swift-expert

Use swift-expert when you need to make or review a Swift architectural decision. Examples:

- "Review the concurrency model in this module"
- "Should I use SwiftUI or UIKit for custom collection layouts?"
- "We're getting Sendable warnings after enabling strict concurrency"

swift-expert reads existing code, makes a clear recommendation with trade-offs, implements the solution, and summarizes what was decided.

## Model Selection

The agents use different models based on the reasoning depth their work requires:

**ios-developer and mobile-developer use sonnet.** Feature implementation and platform strategy are focused tasks with clear patterns. Sonnet handles these efficiently -- it follows conventions, produces working code, and completes tasks without excessive deliberation.

**swift-expert uses opus.** Architectural decisions require deeper reasoning. Evaluating concurrency models, comparing framework trade-offs, and designing type hierarchies involve weighing multiple factors against each other. Opus provides the reasoning depth these decisions need.

This reflects a practical principle: use the most capable model where judgment matters most, and use the faster model where established patterns apply.

## What This Plugin Does Not Cover

**Android-specific native development.** There is no dedicated Kotlin or Android agent. If mobile-developer recommends native Android, you get the strategy but not a Kotlin implementation agent. Native Android features need to be implemented with general-purpose agents or a future Android-specific plugin.

**Web-based mobile apps.** Progressive web apps, responsive web, and mobile web development are handled by the web-development plugin. mobile-developer covers native and cross-platform native apps only.

**macOS applications.** Swift code that targets macOS is covered by the desktop-development plugin. swift-expert can review Swift architecture for any platform, but ios-developer is iOS-specific and does not handle macOS AppKit or Catalyst patterns.

## Cross-Plugin Relationships

The mobile-development agents work alongside agents from other plugins:

- **programming-languages plugin:** Use swift-pro for Swift code quality and simplification tasks that are not architectural in nature -- refactoring verbose code, improving naming, or applying Swift idioms.

- **desktop-development plugin:** Use electron-go-pro or the macOS agents for desktop applications. swift-expert can still review Swift architecture for macOS code, but ios-developer does not handle macOS-specific patterns.

- **code-quality plugin:** Use code-reviewer for general code review that goes beyond Swift-specific architecture -- test coverage, documentation, API design consistency.

## See Also

- [Agent Reference](../../reference/agents/) -- full specifications for all three agents
- [Getting Started](../../tutorials/getting-started/) -- tutorial showing ios-developer and swift-expert working together
- [Build an iOS Feature](../../howto/build-ios-feature/) -- practical guide for ios-developer
- [Choose Mobile Platform](../../howto/choose-mobile-platform/) -- practical guide for mobile-developer
