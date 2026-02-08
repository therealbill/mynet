---
title: "Choose Mobile Platform"
description: "Use mobile-developer to decide between native and cross-platform approaches"
weight: 2
---

# Choose Mobile Platform

Use the mobile-developer agent to make an informed decision between native iOS/Android development and cross-platform frameworks like React Native or Flutter.

## Goal

Get a clear platform recommendation based on your project requirements, team skills, and performance needs -- not a list of pros and cons, but a definitive strategy with trade-offs stated.

## Step 1: Describe Your Requirements

Tell Claude Code about your project with enough detail for mobile-developer to make a recommendation:

- **Platforms:** Which platforms you need (iOS only, iOS + Android, iOS + Android + web)
- **Feature parity:** How similar the apps need to be across platforms (percentage estimate)
- **Team skills:** What your team knows (React/JS, Dart/Flutter, Swift, Kotlin)
- **Performance needs:** Whether the app is data-display, media-heavy, or real-time

```
We need iOS and Android apps with 80% feature parity.
Team has strong React experience. The app is mostly data display
with some camera integration for document scanning.
```

mobile-developer activates when you describe multi-platform requirements, native vs cross-platform decisions, or offline sync needs.

## Step 2: Review the Platform Recommendation

mobile-developer evaluates your requirements against its decision framework:

- **Cross-platform when:** Feature parity is high, team has web/Dart experience, app is data-display or form-driven
- **Native when:** Performance-critical (games, video, AR), relies heavily on platform-specific APIs, or the team has strong native developers
- **Flutter when:** Greenfield project with designers who want pixel-perfect control
- **React Native when:** Team has existing React/JS expertise or needs heavy web code sharing

The agent recommends one approach and explains why, with specific trade-offs for your situation.

```
Why React Native over Flutter for our case? What do we lose?
```

## Step 3: Review the Code Sharing Strategy

If the recommendation is cross-platform, mobile-developer outlines a code sharing strategy:

- **Shared layers:** Business logic, state management, API clients, data models
- **Isolated layers:** Platform-specific UI adaptations, native module bridges, navigation patterns
- **File organization:** Platform-specific file extensions (`.ios.ts` / `.android.ts`) for divergent behavior instead of runtime `Platform.OS` checks

For your document scanning requirement, expect the agent to flag camera integration as a native module boundary -- shared API surface, platform-specific implementation underneath.

## Step 4: Understand Platform-Specific Considerations

mobile-developer addresses the practical implications of the choice:

- **React Native:** Use Expo unless native module requirements force bare workflow. Use Turbo Modules (New Architecture) for new native modules. Target Hermes engine for performance.
- **Flutter:** Keep widget tree shallow. Use `const` constructors. Platform channels for native integration.
- **Both:** Respect each platform's navigation paradigm -- bottom tabs on iOS, drawer navigation common on Android.

## Verification

The decision is sound when:

- The chosen platform aligns with your team's existing skills (or the ramp-up cost is justified)
- Performance-critical features have a clear path to native performance (native modules, platform channels)
- The code sharing strategy separates shared logic from platform-specific code cleanly
- Offline requirements, if any, have a defined local database and sync strategy

## Troubleshooting

### The answer is "both native and cross-platform"

Some projects benefit from a hybrid approach -- cross-platform for most features, native for performance-critical screens. This is valid when:

- One feature (like video editing or AR) needs native performance while the rest is data-display
- You want cross-platform for feature velocity but native for a flagship feature

Ask mobile-developer to define the boundary:

```
Can we use React Native for most screens but native Swift/Kotlin for the document scanner?
```

### Requirements change after the decision

If your app grows beyond the original requirements (adding AR, real-time video, or platform-specific integrations), revisit the decision:

```
We originally chose React Native but now need ARKit integration -- should we reconsider?
```

mobile-developer re-evaluates based on the new requirements and may recommend adding native modules rather than switching platforms entirely.

### Need native iOS implementation after choosing native

If mobile-developer recommends native development, use ios-developer for the iOS implementation:

```
Build the product list screen as a native SwiftUI view
```

## See Also

- [Getting Started]({{< ref "tutorials/getting-started" >}}) -- end-to-end tutorial with ios-developer and swift-expert
- [Agent Reference]({{< ref "reference/agents" >}}) -- mobile-developer specification
- [Architecture]({{< ref "explanation/architecture" >}}) -- how the three agents cover different concerns
