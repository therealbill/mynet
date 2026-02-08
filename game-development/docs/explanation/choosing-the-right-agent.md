---
title: "Choosing the Right Agent"
description: "Decision guide for selecting the appropriate game-development agent based on your task"
weight: 2
---

# Choosing the Right Agent

The game-development plugin has five agents. Each handles a different aspect of game creation. This page explains what each agent is for and how they combine in common workflows.

## By Domain

### game-designer

Use game-designer when the question is about player experience, not code.

- Designing a new mechanic (what it does, how it feels, what parameters control it)
- Balancing an economy (currency earn rates, spending opportunities, inflation prevention)
- Designing a progression system (XP curves, level pacing, content gates)
- Choosing between difficulty approaches (dynamic adjustment vs discrete settings)
- Evaluating whether a mechanic is working (playtesting methodology, what to measure)
- Prototyping strategy (paper prototypes, spreadsheet simulations before coding)

game-designer produces specifications, parameter ranges, formulas, and interaction rules. It does not produce code. Its output is the input for game-developer or an engine-specific agent.

### game-developer

Use game-developer when the question is about architecture or applies across multiple engines.

- Selecting an engine for a new project (Unity vs Unreal vs Godot vs web)
- Designing the overall system architecture (ECS, component patterns, data-driven design)
- Setting performance budgets (frame time targets, memory ceilings, entity count limits)
- Choosing between architectural approaches (object pooling strategies, spatial partitioning algorithms, state machine patterns)
- Networking architecture (server authority model, delta compression, interest management)
- Any implementation question that is not specific to a single engine

game-developer produces engine-agnostic architecture guidance. It recommends patterns and structures. When the work needs engine-specific code, it hands off to unity-game-developer or unreal-engine-developer.

### unity-game-developer

Use unity-game-developer when you are working in Unity and need Unity-specific answers.

- Setting up a Unity project (folder structure, assembly definitions, rendering pipeline)
- Writing C# that uses Unity APIs (MonoBehaviour, ScriptableObject, Input System, Physics2D)
- Solving Unity-specific performance problems (GC spikes from LINQ or GetComponent, draw call batching, physics layer matrix)
- Choosing between Unity systems (URP vs HDRP, UI Toolkit vs uGUI, Addressables vs Resources)
- Configuring Unity packages (new Input System, Cinemachine, TextMeshPro)

unity-game-developer writes C# code that compiles in Unity. It knows Unity's serialization rules, lifecycle methods, and the practical differences between Unity's overlapping systems.

### unreal-engine-developer

Use unreal-engine-developer when you are working in Unreal and need UE-specific answers.

- Setting up an Unreal project (Gameplay Framework class hierarchy, module structure)
- Deciding the C++/Blueprint boundary for a system (core logic in C++, designer tuning in Blueprint)
- Configuring UE5 features (Nanite for static meshes, Lumen for dynamic GI, World Partition for open worlds)
- Implementing multiplayer replication (UPROPERTY(Replicated), OnRep callbacks, Server/Client RPCs)
- Using GAS (Gameplay Ability System) for complex ability interactions
- Solving Unreal-specific performance problems (excessive Tick usage, synchronous asset loads, Blueprint hot paths)

unreal-engine-developer writes C++ following Epic's coding standards and creates Blueprint-friendly interfaces. It knows the Gameplay Framework roles and when to use (or avoid) UE5 features based on target hardware.

### 3d-artist

Use 3d-artist when the question involves visual assets or the art pipeline.

- Setting polygon budgets for a target platform (mobile vs PC vs console, by asset category)
- Planning LOD levels for real-time assets
- Configuring texture workflows (channel packing, atlas strategies, resolution targets)
- Setting up the DCC-to-engine pipeline (Blender/Maya export settings, engine import settings)
- Defining naming conventions for assets (meshes, textures, materials)
- Optimizing existing assets that exceed performance budgets

3d-artist works with any engine. It handles the visual production side that the programming agents do not cover.

## Common Workflows

### New game from scratch

Start with game-designer to define the core loop and key mechanics. Move to game-developer for engine selection and architecture. Then use the appropriate engine-specific agent for implementation. Bring in 3d-artist when you need visual assets.

```
game-designer  -->  game-developer  -->  unity-game-developer
(core loop)        (engine choice)       (implementation)
                                              |
                                         3d-artist
                                    (assets in parallel)
```

### Mechanic prototyping

Use game-designer for the mechanic specification with tunable parameters. Use game-developer for engine-agnostic architecture. Test with a minimum viable implementation before committing to full production.

```
game-designer  -->  game-developer
(parameters)       (architecture + MVP)
```

### Performance optimization

Go directly to the engine-specific agent. Unity GC problems need unity-game-developer. Unreal Tick overhead needs unreal-engine-developer. Only involve game-developer if the problem is architectural (wrong data structure choice, missing spatial partitioning) rather than engine-specific.

```
unity-game-developer    (Unity perf problems)
unreal-engine-developer (Unreal perf problems)
game-developer          (architectural perf problems)
```

### Art pipeline setup

Start with 3d-artist for the pipeline definition. Involve the engine-specific agent for import configuration on the engine side.

```
3d-artist  -->  unity-game-developer or unreal-engine-developer
(export)        (import settings)
```

### Balancing an existing game

Use game-designer exclusively. Economy problems, progression pacing, and difficulty tuning are design problems, not code problems. game-designer analyzes the system and recommends parameter changes that can be applied through existing configuration without code modifications.

```
game-designer
(analysis + parameter adjustment)
```

## When Agents Overlap

Some questions sit on the boundary between agents. Here are the tiebreakers:

- **"How do I implement X?"** -- If X is engine-specific, use the engine-specific agent. If X is engine-agnostic, use game-developer.
- **"Should I use X or Y?"** -- If both options are within one engine (URP vs HDRP), use the engine-specific agent. If the options span engines or are architectural (ECS vs inheritance), use game-developer.
- **"Is this mechanic working?"** -- If the question is about player experience and feel, use game-designer. If the question is about whether the code is correct, use the engine-specific agent.
- **"My game is too slow"** -- If you know the engine, use the engine-specific agent to profile and fix. If you do not know whether the problem is architectural, start with game-developer for diagnosis.

## See Also

- [Architecture]({{< ref "explanation/architecture" >}}) -- the three-layer design behind these agent roles
- [Agent Reference]({{< ref "reference/agents" >}}) -- technical specifications for all five agents
- [Prototype a Game Mechanic]({{< ref "howto/prototype-game-mechanic" >}}) -- the designer-to-developer workflow in practice
