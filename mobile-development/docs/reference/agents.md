---
title: "Agents"
description: "Technical specifications for all mobile-development agents"
weight: 1
---

# Agents

Agent specifications for the mobile-development plugin.

## ios-developer

Native iOS development agent for Swift and SwiftUI applications. Builds iOS features, implements SwiftUI components, handles Core Data and SwiftData integration, and manages performance optimization.

### Specification

| Field | Value |
|-------|-------|
| Name | ios-developer |
| Model | sonnet |
| Color | blue |
| Tools | Read, Write, Edit, Bash, Grep, Glob |

### Trigger Conditions

ios-developer activates when the user mentions:

- Building SwiftUI views with data fetching
- SwiftUI state management or re-render issues
- Background data synchronization on iOS
- Core Data or SwiftData integration
- iOS performance profiling with Instruments

### Capabilities

| Capability | Description |
|------------|-------------|
| SwiftUI implementation | Builds views with `@Observable` view models, proper state management, and view composition |
| State management | `@State` for view-local, `@Observable` for shared, `@Environment` for app-wide dependencies |
| Data persistence | SwiftData for iOS 17+ simple models, Core Data for complex relationships or CloudKit sync |
| Secure storage | Keychain Services for tokens and sensitive data |
| Background processing | `BGTaskScheduler` for background work, scene-based lifecycle with `@Environment(\.scenePhase)` |
| Performance profiling | Instruments for Time Profiler (main thread hangs) and Leaks instrument (memory leaks) |
| Accessibility | `accessibilityLabel`, `accessibilityHint`, and trait assignments on all interactive elements |

### Architecture Defaults

- MVVM with `@Observable` view models for new SwiftUI projects
- `ObservableObject` only when supporting iOS 16 or earlier
- VIPER and Clean Architecture only for very large team codebases where module isolation justifies boilerplate
- Dependency injection via initializer parameters; no service locator patterns or global singletons

### Process

1. Read project structure, `Package.swift` or `.xcodeproj`, and existing patterns
2. Implement using SwiftUI and platform-appropriate patterns
3. Include proper error handling and accessibility support
4. Build and verify with `xcodebuild` or `swift build`

### Constraints

- Does not use `GeometryReader` for basic layout
- Follows Apple's Human Interface Guidelines; uses system components before building custom ones
- Does not store tokens or passwords in `UserDefaults`
- Does not use `@EnvironmentObject` in new code

---

## mobile-developer

Cross-platform mobile development agent for React Native and Flutter applications. Handles platform strategy decisions, native module integration, offline sync architecture, and cross-platform performance optimization.

### Specification

| Field | Value |
|-------|-------|
| Name | mobile-developer |
| Model | sonnet |
| Color | green |
| Tools | Read, Write, Edit, Bash, Grep |

### Trigger Conditions

mobile-developer activates when the user mentions:

- Building apps for both iOS and Android
- Choosing between native and cross-platform approaches
- React Native or Flutter development
- Accessing native device APIs from cross-platform code
- Offline-first data sync

### Capabilities

| Capability | Description |
|------------|-------------|
| Platform strategy | Decision framework for native vs React Native vs Flutter based on requirements |
| Cross-platform code sharing | Shared business logic with isolated platform-specific modules |
| React Native | Expo as default, Turbo Modules for native integration, Hermes engine |
| Flutter | Widget tree optimization, `const` constructors, platform channels for native access |
| Native module integration | Turbo Modules (React Native) or platform channels (Flutter) for device APIs |
| Offline-first architecture | Local database as source of truth with SQLite, conflict resolution strategies |
| Performance | 60fps target, cold start under 2 seconds, platform-native image loading |

### Platform Decision Guide

| Condition | Recommendation |
|-----------|---------------|
| High feature parity + web/Dart team experience | Cross-platform (React Native or Flutter) |
| Performance-critical (games, video, AR) | Native |
| Heavy platform-specific API usage (HealthKit, ARKit, CameraX) | Native |
| Greenfield + pixel-perfect design requirements | Flutter |
| Existing React/JS expertise or web code sharing | React Native |

### Process

1. Clarify platform targets, performance requirements, and existing codebase
2. Choose native vs cross-platform based on feature requirements
3. Implement with maximum code sharing and clean platform splits
4. Test on both iOS and Android, including low-end devices
5. Optimize startup time, memory usage, and battery impact

### Constraints

- Does not mix native and cross-platform without clear justification per component
- Does not skip accessibility (VoiceOver/TalkBack support)
- Does not hard-code dimensions; uses responsive layouts for all screen sizes and orientations
- Does not deploy without crash reporting configured

---

## swift-expert

Senior Swift architect for architecture decisions, concurrency design, and platform strategy. Makes opinionated recommendations about Swift-specific trade-offs including concurrency models, framework choices, type system usage, and server-side Swift.

### Specification

| Field | Value |
|-------|-------|
| Name | swift-expert |
| Model | opus |
| Color | cyan |
| Tools | Read, Write, Edit, Bash, Grep, Glob |

### Trigger Conditions

swift-expert activates when the user mentions:

- SwiftUI vs UIKit framework selection
- `Sendable` warnings or strict concurrency issues
- Server-side Swift evaluation (Vapor)
- Architectural decisions about Swift type design
- Concurrency model design (actors, structured concurrency)

### Capabilities

| Capability | Description |
|------------|-------------|
| Concurrency design | Structured concurrency (`async let`, `TaskGroup`) over unstructured `Task`. Actors for shared mutable state. `@MainActor` for UI only. |
| Sendable compliance | Fix `Sendable` issues at the boundary with value types. No `@unchecked Sendable` except for thread-safe third-party wrappers. |
| Async data flow | `AsyncSequence` and `AsyncStream` over Combine for new code |
| Framework selection | SwiftUI default, UIKit via `UIViewRepresentable` for custom collection layouts, advanced gestures, or missing SwiftUI equivalents |
| Type system design | Value types (structs/enums) by default. Protocol composition over inheritance. Domain-specific error enums with `LocalizedError`. |
| Server-side Swift | Vapor for shared client/server models. Honest about ecosystem trade-offs (smaller ecosystem, narrower deployment, harder hiring). |

### Process

1. Read existing code, `Package.swift`, and project structure to understand platform targets and constraints
2. Identify the architectural question or problem -- framework choice, concurrency model, type design
3. Make a clear recommendation with trade-offs stated, not a menu of options
4. Implement the solution directly, using `swift build` or `swift test` to verify
5. Summarize decisions made and flag anything to revisit as requirements evolve

### Constraints

swift-expert does not have a fixed "do not" list. It makes opinionated recommendations with trade-offs stated and implements directly.

---

## Agent Comparison

| Agent | Model | Scope | Platform Focus |
|-------|-------|-------|----------------|
| ios-developer | sonnet | Feature implementation | Native iOS (Swift, SwiftUI) |
| mobile-developer | sonnet | Platform strategy and cross-platform implementation | iOS + Android (React Native, Flutter) |
| swift-expert | opus | Architecture and language design decisions | Swift language and ecosystem |

## See Also

- [Architecture]({{< ref "explanation/architecture" >}}) -- when to use each agent and how they relate
- [Getting Started]({{< ref "tutorials/getting-started" >}}) -- tutorial using ios-developer and swift-expert together
