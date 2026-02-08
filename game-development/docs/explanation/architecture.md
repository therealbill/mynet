---
title: "Architecture"
description: "Design vs implementation vs engine-specific: the three-layer agent architecture of the game-development plugin"
weight: 1
---

# Architecture

The game-development plugin organizes its five agents into a three-layer architecture: design, implementation, and engine-specific. This page explains why those layers exist, how they interact, and where the 3d-artist fits in.

## Three Layers

The plugin separates game creation into three distinct concerns, each handled by a different type of agent.

**Design layer: game-designer.** This agent thinks about player experience. It defines mechanics, balances economies, designs progression curves, and reasons about reward psychology. Its output is specifications and parameter ranges, not code. game-designer does not know or care which engine implements its designs.

**Implementation layer: game-developer.** This agent thinks about code architecture. It selects engines, designs component systems, sets performance budgets, and recommends patterns like ECS, data-driven design, and simulation/presentation separation. game-developer works across engines. It produces architectural guidance and engine-agnostic code structures, not engine-specific API calls.

**Engine-specific layer: unity-game-developer and unreal-engine-developer.** These agents think in the idioms of a specific engine. unity-game-developer writes C# using MonoBehaviour lifecycle methods, ScriptableObjects, URP configuration, and Unity's physics API. unreal-engine-developer writes C++ with UCLASS macros, Gameplay Framework conventions, Nanite/Lumen configuration, and Unreal's replication system. Each agent produces code that compiles and runs in its target engine.

## Why Separate Design from Implementation

Game design and game programming are different disciplines that reason about different problems.

A game designer asks: "Does the jump feel responsive? Is the difficulty curve frustrating at level 12? Are players spending currency faster than they earn it?" These questions are about player experience, pacing, and psychology. The answers are parameter values, formulas, and interaction rules -- not code.

A game developer asks: "Should entities use composition or inheritance? Where does the simulation boundary sit? What is the frame budget for 200 enemies on screen?" These questions are about software architecture and runtime performance. The answers are patterns, data structures, and optimization strategies.

Collapsing these into a single agent would force it to context-switch between player psychology and memory allocation. By keeping them separate, each agent can go deep into its domain. game-designer produces a mechanic specification with tunable parameters; game-developer translates that specification into an architecture that supports runtime tuning without recompilation. Neither agent needs to do the other's job.

This mirrors how professional game studios work. Game designers create design documents. Programmers implement those documents in code. The handoff is explicit, and the artifacts are distinct.

## Why Engine-Specific Agents

Unity and Unreal Engine are not interchangeable skins over the same engine. They represent fundamentally different paradigms.

Unity uses C# with a component-based architecture centered on MonoBehaviours and ScriptableObjects. Its rendering pipeline choice (URP vs HDRP) determines the visual capability ceiling. Performance optimization means avoiding GC allocations, caching GetComponent calls, and using NonAlloc physics queries. The ecosystem revolves around the Asset Store and Unity Package Manager.

Unreal uses C++ for core systems with Blueprint for designer-facing logic. Its architecture is built on the Gameplay Framework (GameMode, GameState, PlayerState, PlayerController), which prescribes where different types of logic live. UE5 features like Nanite and Lumen have specific hardware requirements and fallback paths. Performance optimization means minimizing Tick usage, using Unreal Insights, and moving Blueprint hot paths to C++.

Generic advice like "use object pooling" or "separate data from behavior" is helpful at the architecture level -- and game-developer provides it. But actually implementing object pooling in Unity (ObjectPool<T>, custom pool with Instantiate/Destroy wrappers) looks nothing like implementing it in Unreal (TArray-based pool with SpawnActor/DestroyActor patterns). Engine-specific agents close the gap between architectural intent and working code.

## The 3d-artist as a Cross-Cutting Concern

The 3d-artist agent does not belong to any single layer. It operates across all of them.

During design, 3d-artist informs what is visually feasible within polygon and texture budgets for the target platform. A game designer might want 50 unique enemy types on screen simultaneously -- 3d-artist calculates whether the art budget supports that.

During implementation, 3d-artist defines the asset pipeline: export settings from DCC tools (Blender, Maya), import settings in the engine, naming conventions that batch tools depend on, and LOD strategies that the rendering system uses.

During engine-specific work, 3d-artist configures engine-specific import settings, material setups, and texture channel packing conventions. PBR texture setup for Unity's URP differs from setup for Unreal's material system.

This cross-cutting role reflects how technical art works in studios. A technical artist collaborates with designers (to scope visual ambition), programmers (to define the asset pipeline), and engine specialists (to configure import and rendering settings).

## How the Layers Collaborate

The typical workflow moves from left to right across the layers:

1. **game-designer** produces a mechanic specification with parameters, formulas, and interaction rules
2. **game-developer** translates that specification into an engine-agnostic architecture (state machines, data structures, component decomposition)
3. **unity-game-developer** or **unreal-engine-developer** implements the architecture in engine-specific code
4. **3d-artist** works alongside any layer to handle visual asset creation and pipeline configuration

Not every task requires all layers. A pure balancing question goes only to game-designer. A Unity GC optimization goes directly to unity-game-developer. An engine selection decision goes to game-developer. The layers exist to be used as needed, not as a mandatory pipeline.

## See Also

- [Choosing the Right Agent]({{< ref "explanation/choosing-the-right-agent" >}}) -- practical decision guide for which agent to invoke
- [Agent Reference]({{< ref "reference/agents" >}}) -- technical specifications for all five agents
- [Getting Started]({{< ref "tutorials/getting-started" >}}) -- tutorial demonstrating the design-to-implementation handoff
