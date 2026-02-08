---
title: "Agents"
description: "Technical specifications for all game-development plugin agents"
weight: 1
---

# Agents

Agent specifications for the game-development plugin. This plugin contains five agents covering game design, general game development, engine-specific implementation, and 3D art production.

## 3d-artist

Technical 3D art specialist for game-ready asset creation and art pipeline efficiency.

### Specification

| Field | Value |
|-------|-------|
| Name | 3d-artist |
| Model | sonnet |
| Color | cyan |
| Tools | Read, Write, Edit, Bash |

### Trigger Conditions

3d-artist activates when the user mentions:

- Optimizing 3D assets for a target platform (polygon budgets, LOD)
- Texture workflows, channel packing, or PBR material setup
- Asset pipeline configuration (Blender to Unity, DCC to engine export)
- UV layout, texel density, or texture atlasing
- Naming conventions for 3D assets

### Capabilities

| Capability | Description |
|------------|-------------|
| Polygon budgets | Sets per-category triangle budgets based on target platform. Mobile hero: 5-15k tris. PC/console hero: 30-100k tris. |
| LOD strategy | Plans 3-4 LOD levels (full detail, 50%, 25%, billboard/impostor) from the start of asset creation. |
| Texture optimization | Configures texture atlases for shared-material props. Packs channels (roughness/metallic/AO into one texture RGB). Power-of-two resolutions. |
| UV layout | Maximizes texel density consistency. Prioritizes player-facing surfaces. Minimizes UV island count. |
| Export pipeline | Defines FBX/glTF export settings including scale, axis orientation, and smoothing group handling. |
| Naming conventions | Enforces consistent naming (SM_Name_LOD0, T_Name_BaseColor) for pipeline and batch tool compatibility. |

### Example Interactions

```
User: "My character models are 50k tris each and the game runs poorly on phones"
Agent: Establishes polygon budgets per asset category and creates an LOD strategy for mobile hardware.

User: "How should I set up my PBR textures for Unreal?"
Agent: Recommends texture channel packing, resolution targets, and material setup for Unreal's rendering pipeline.

User: "We need a consistent pipeline from Blender to Unity for our art team"
Agent: Defines export settings, naming conventions, and import configurations for a Blender-to-Unity workflow.
```

---

## game-designer

Game design specialist for mechanics, systems balancing, progression, and player experience.

### Specification

| Field | Value |
|-------|-------|
| Name | game-designer |
| Model | sonnet |
| Color | magenta |
| Tools | Read, Write, Edit |

### Trigger Conditions

game-designer activates when the user mentions:

- Game mechanics design or core loop definition
- Economy balancing, currency sources and sinks
- Progression systems, XP curves, or leveling
- Difficulty systems (dynamic or discrete)
- Reward psychology or engagement mechanics
- Prototyping or playtesting methodology

### Capabilities

| Capability | Description |
|------------|-------------|
| Core loop design | Defines the action-reward-progression loop before secondary systems. |
| Economy balancing | Maps currency sources and sinks, calculates earn rates and spending opportunities, uses spreadsheet models. |
| Progression curves | Designs exponential or polynomial XP curves with content gates at milestones. Front-loads early levels for onboarding. |
| Difficulty systems | Evaluates dynamic difficulty adjustment vs discrete settings based on genre and audience. |
| Reward psychology | Applies variable ratio reinforcement (randomized rewards) bounded to avoid frustration, mixed with fixed schedules. |
| Prototyping guidance | Recommends paper prototypes and spreadsheet simulations before committing to code implementation. |

### Example Interactions

```
User: "How should I structure XP and leveling so it doesn't feel grindy?"
Agent: Designs an XP curve that front-loads early levels and uses diminishing returns with content-gated milestones.

User: "Players accumulate gold way too fast and there's nothing worth buying"
Agent: Analyzes the economy's sources and sinks, identifies the imbalance, and designs corrective sinks.

User: "Should I use discrete difficulty levels or dynamic difficulty adjustment?"
Agent: Evaluates both approaches against the game's genre, audience, and core experience goals.
```

---

## game-developer

Generalist game development architect covering engine selection, game architecture, performance optimization, and cross-platform strategy.

### Specification

| Field | Value |
|-------|-------|
| Name | game-developer |
| Model | sonnet |
| Color | blue |
| Tools | Read, Write, Edit, Bash |

### Trigger Conditions

game-developer activates when the user mentions:

- Choosing a game engine for a new project
- Game system architecture (ECS, component patterns, data-driven design)
- Cross-engine or engine-agnostic design questions
- Performance optimization strategy (frame budgets, memory budgets, object pooling)
- Networking architecture for multiplayer games

### Capabilities

| Capability | Description |
|------------|-------------|
| Engine selection | Evaluates Unity, Unreal, Godot, and web frameworks against project requirements (team skills, target platform, scope). |
| Architecture patterns | Recommends composition over inheritance (ECS), data-driven design, and simulation/presentation separation. |
| Performance strategy | Sets frame and memory budgets. Recommends object pooling, spatial partitioning, and GC minimization. |
| Cross-platform guidance | Advises on platform-specific constraints and abstraction strategies. |
| Process discipline | Enforces profile-before-optimizing methodology. Recommends engine profilers for actual hotspot identification. |

### Engine Selection Defaults

| Engine | Recommended For |
|--------|----------------|
| Unity | Mobile, 2D, indie scope, rapid prototyping. C# ecosystem. |
| Unreal | AAA visuals, large teams, C++/Blueprint workflow. |
| Godot | Lightweight 2D/3D, open-source requirements, GDScript simplicity. |
| Web (Phaser/Three.js/PlayCanvas) | Browser-first distribution. |

### Example Interactions

```
User: "I'm building a 2D roguelike with online co-op -- what engine should I use?"
Agent: Evaluates engine options against 2D focus, procedural generation, and networked multiplayer requirements.

User: "My game drops to 20 FPS when there are more than 200 enemies on screen"
Agent: Diagnoses the bottleneck and recommends spatial partitioning, object pooling, or LOD strategies.

User: "How should I structure my inventory and crafting system so it's extensible?"
Agent: Designs a data-driven architecture using composition patterns appropriate for the engine.
```

---

## unity-game-developer

Unity-specific game developer for C# scripting, URP/HDRP rendering, mobile optimization, and Unity project architecture.

### Specification

| Field | Value |
|-------|-------|
| Name | unity-game-developer |
| Model | sonnet |
| Color | green |
| Tools | Read, Write, Edit, Bash |

### Trigger Conditions

unity-game-developer activates when the user mentions:

- Unity project setup or architecture decisions
- URP vs HDRP rendering pipeline selection
- Unity-specific C# scripting patterns (MonoBehaviour, ScriptableObject)
- Unity performance problems (GC spikes, draw calls, batching)
- Unity packages and systems (Input System, Addressables, UI Toolkit)

### Capabilities

| Capability | Description |
|------------|-------------|
| Render pipeline selection | Defaults to URP for mobile and indie. HDRP for PC/console with advanced lighting needs. Built-in for legacy only. |
| Data architecture | ScriptableObjects for shared game data. MonoBehaviours for runtime behavior. |
| Input configuration | New Input System package with rebinding support. Legacy Input class avoided. |
| UI systems | UI Toolkit for editor tools and HUD. uGUI for world-space and animated game UI. |
| Async patterns | UniTask or Awaitable (Unity 2023+) for complex async flows. Coroutines for simple delays. |
| Asset loading | Addressables for any project with significant asset volume. |

### Performance Rules

| Rule | Rationale |
|------|-----------|
| Cache GetComponent results | GetComponent allocates and is slow in Update loops. |
| Avoid LINQ in hot paths | LINQ methods allocate enumerator objects every call. |
| Use NonAlloc physics queries | RaycastNonAlloc and OverlapSphereNonAlloc avoid per-frame heap allocations. |
| Configure physics layer matrix | Skip unnecessary collision checks between unrelated layers. |
| Bake lighting on mobile | Real-time lights are expensive on mobile GPUs. |
| Object pool runtime instantiation | Use ObjectPool<T> or custom pools for projectiles, particles, enemies. |

### Example Interactions

```
User: "I'm starting a Unity mobile game -- how should I structure the project?"
Agent: Sets up folder structure, recommends URP, establishes ScriptableObject-based data architecture.

User: "I'm getting frame hitches every few seconds from garbage collection"
Agent: Identifies allocation hotspots and applies object pooling, cached references, and non-allocating physics queries.

User: "Should I use URP or HDRP for my project?"
Agent: Evaluates tradeoffs based on target platforms and visual requirements.
```

---

## unreal-engine-developer

Unreal Engine specialist for C++ gameplay programming, Blueprint integration, rendering pipeline configuration, and UE5-specific systems.

### Specification

| Field | Value |
|-------|-------|
| Name | unreal-engine-developer |
| Model | sonnet |
| Color | yellow |
| Tools | Read, Write, Edit, Bash |

### Trigger Conditions

unreal-engine-developer activates when the user mentions:

- Unreal Engine project architecture or Gameplay Framework usage
- C++ vs Blueprint decisions for game systems
- Nanite, Lumen, or other UE5-specific rendering features
- Unreal replication, RPCs, or multiplayer networking
- Gameplay Ability System (GAS) implementation
- Unreal-specific performance profiling or optimization

### Capabilities

| Capability | Description |
|------------|-------------|
| C++/Blueprint boundary | Core systems and performance-critical code in C++. Designer-facing tuning and prototyping in Blueprint. UFUNCTION(BlueprintCallable) for the bridge. |
| Gameplay Framework | GameMode for rules, GameState for shared state, PlayerState for per-player data, PlayerController for input. Roles are not collapsed. |
| Nanite | Static meshes with high geometric complexity. Not for skeletal meshes, translucent materials, or unsupported hardware. |
| Lumen | Dynamic GI on capable hardware (current-gen consoles, mid+ PC). Fallback to screen-space or baked for weaker targets. |
| World Partition | Open-world and large maps. Not used for small linear levels. |
| GAS | Gameplay Ability System for complex ability/effect interactions. Not recommended for simple action games. |
| Replication | UPROPERTY(Replicated) with OnRep callbacks for state. RPCs (Server/Client/Multicast) for events, not continuous state. |

### Performance Rules

| Rule | Rationale |
|------|-----------|
| Use Unreal Insights for profiling | More detailed than stat commands for identifying bottlenecks. |
| Minimize Tick usage | Avoid Tick on Actors that do not need per-frame updates. Use Timers or event-driven patterns. |
| Set tick intervals | Actors that need Tick but tolerate reduced frequency should have explicit tick intervals. |
| Async asset loading | Use StreamableManager. Synchronous loads cause hitches. |
| Move hot paths to C++ | Blueprint nativization is not a substitute for C++ performance in critical code. |
| Use UPROPERTY() for UObject pointers | Prevents garbage collection of referenced objects. |

### Example Interactions

```
User: "Should I implement my combat system in C++ or Blueprint?"
Agent: Recommends the C++/Blueprint split based on system complexity and team workflow.

User: "How do I replicate my inventory system to clients?"
Agent: Designs replication strategy using UPROPERTY(Replicated), RPCs, and authority checks.

User: "Lumen is tanking my frame rate on the target hardware"
Agent: Diagnoses Lumen performance and recommends fallback strategies for constrained hardware.
```

## See Also

- [Architecture]({{< ref "explanation/architecture" >}}) -- design decisions behind the three-layer agent structure
- [Choosing the Right Agent]({{< ref "explanation/choosing-the-right-agent" >}}) -- decision guide for agent selection
- [Getting Started]({{< ref "tutorials/getting-started" >}}) -- tutorial using game-designer and unity-game-developer together
