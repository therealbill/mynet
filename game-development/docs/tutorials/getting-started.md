---
title: "Getting Started with Game Development"
description: "Design a jump mechanic with game-designer and implement it in Unity with unity-game-developer"
weight: 1
---

# Getting Started with Game Development

Design a tunable jump mechanic using the game-designer agent, then implement it in Unity using the unity-game-developer agent. This tutorial demonstrates the design-to-implementation handoff that defines the game-development plugin's workflow.

## What You'll Build

By the end of this tutorial, you will have:

- Used game-designer to define a jump mechanic with tunable parameters and feel targets
- Reviewed a design document specifying jump height, hang time, coyote time, and input buffering
- Used unity-game-developer to implement the mechanic in Unity with proper C# architecture
- Tested the mechanic in Unity's Play mode and adjusted parameters through ScriptableObjects

This takes about 25 minutes.

## Prerequisites

- Claude Code CLI installed and authenticated (run `claude --version` to verify)
- The game-development plugin installed in your project's `.claude/settings.json`
- Unity 2022.3 LTS or later installed with a new 2D URP project open
- Basic familiarity with Unity's Editor (Scene view, Inspector, Play mode)

## Step 1: Design the Mechanic

Open Claude Code and describe the mechanic you want to design:

```
I'm building a 2D platformer. I need a jump mechanic that feels responsive
and satisfying. The game is aimed at casual players, not hardcore precision platforming.
```

Claude Code matches your request to the game-designer agent based on the mention of game mechanics and player experience goals. The agent focuses on design parameters, not code.

game-designer asks clarifying questions:

- What is the player character's movement speed? (Moderate -- not a speedrunner game)
- Should the player have double jump or just single? (Single jump for now, extensible later)
- What is the target platform? (Mobile and PC)
- How forgiving should the jump window be? (Very -- casual audience)

Answer these questions. The agent uses your responses to produce a mechanic design document.

## Step 2: Review the Design Document

game-designer produces a design specification. For a casual-friendly jump mechanic, expect parameters like:

- **Jump height:** 3-4 tiles (tunable via ScriptableObject)
- **Time to apex:** 0.35-0.45 seconds (controls how floaty the jump feels)
- **Fall multiplier:** 1.5-2.5x gravity on descent (makes landing feel snappy)
- **Coyote time:** 80-150ms (grace period after walking off a ledge where jump still registers)
- **Input buffer:** 100-150ms (pressing jump slightly before landing queues the input)
- **Variable jump height:** Releasing the jump button early reduces upward velocity, giving the player height control

The agent explains the reasoning behind each parameter. Coyote time and input buffering make the jump feel responsive even on mobile touchscreens where input timing is imprecise. The fall multiplier prevents the floaty feeling that pure parabolic arcs create.

Review the design. Ask the agent to adjust any parameters that do not fit your game's feel. For example:

```
Make it less floaty. I want a snappier jump with less hang time at the apex.
```

### Checkpoint

At this point you should have:

- Triggered game-designer with your mechanic description
- Received a design document with tunable parameters and their rationale
- Adjusted any parameters that did not match your target feel

If game-designer did not activate, verify the game-development plugin is listed in your `.claude/settings.json` and restart Claude Code.

## Step 3: Hand Off to Implementation

Tell Claude Code you are ready to implement:

```
Implement this jump mechanic in Unity. Use proper architecture with
ScriptableObjects for the tunable parameters.
```

Claude Code routes this to the unity-game-developer agent based on the Unity-specific request. The agent receives the design parameters from the previous conversation context and translates them into Unity architecture decisions.

unity-game-developer generates the following files:

- **`JumpSettings.cs`** -- A ScriptableObject holding all tunable parameters (jump height, time to apex, fall multiplier, coyote time, input buffer duration)
- **`PlayerJump.cs`** -- A MonoBehaviour implementing the jump state machine with physics-based movement
- **`GroundCheck.cs`** -- A MonoBehaviour using Physics2D.OverlapCircleNonAlloc to detect ground contact without per-frame allocations

The agent makes specific Unity decisions:

- Calculates gravity and initial velocity from jump height and time-to-apex rather than using magic numbers
- Uses the new Input System for jump input with a dedicated `PlayerInputActions` asset
- Implements coyote time and input buffering as timers in Update, with physics applied in FixedUpdate
- Exposes parameters through the ScriptableObject so designers can tune without recompiling

## Step 4: Review the Generated Code

Ask to see the core implementation:

```
Show me PlayerJump.cs and JumpSettings.cs
```

Verify these architectural choices:

- `JumpSettings` inherits from `ScriptableObject` and uses `[CreateAssetMenu]` for easy creation in the Unity Editor
- `PlayerJump` references `JumpSettings` via a `[SerializeField]` field, not a hardcoded dependency
- Ground detection uses `Physics2D.OverlapCircleNonAlloc` with a pre-allocated results array
- The jump calculation derives gravity from `jumpHeight` and `timeToApex` using the formula: `gravity = (2 * jumpHeight) / (timeToApex * timeToApex)`
- Input buffering stores the timestamp of the last jump press and checks against it each frame

### Checkpoint

At this point you should have:

- Three C# scripts generated by unity-game-developer
- A ScriptableObject for tunable parameters separate from the MonoBehaviour logic
- Physics calculations derived from designer-friendly parameters, not raw force values

If any scripts are missing, ask unity-game-developer to regenerate the specific component.

## Step 5: Test in Unity

Create the runtime assets and test:

```
How do I set this up in my Unity scene to test it?
```

unity-game-developer walks you through:

1. Create a `JumpSettings` asset via **Assets > Create > Game > Jump Settings**
2. Set the initial values from the design document (jump height: 3.5, time to apex: 0.4, fall multiplier: 2.0, coyote time: 0.1, input buffer: 0.12)
3. Add `PlayerJump` and `GroundCheck` components to your player GameObject
4. Assign the `JumpSettings` asset to the PlayerJump component's settings field
5. Set up a ground layer and assign it to your platform GameObjects
6. Enter Play mode and test the jump

Try adjusting parameters in the Inspector while the game is running. Changes to the ScriptableObject take effect immediately because the code reads from it each frame rather than caching values at startup.

### Checkpoint

At this point you should have:

- A player character that jumps in Play mode
- Tunable parameters adjustable in real time through the ScriptableObject
- Coyote time and input buffering working (test by walking off a ledge and pressing jump slightly late)

## What You Learned

In this tutorial, you:

- **Used game-designer for mechanic design** -- the agent produced a parameter specification with rationale rooted in player psychology, not a code implementation
- **Observed the design-to-implementation handoff** -- game-designer defined what the mechanic should feel like; unity-game-developer translated that into Unity-specific architecture
- **Saw ScriptableObject-based data architecture** -- tunable parameters live in data assets, not hardcoded in MonoBehaviours, enabling iteration without recompilation
- **Applied Unity performance conventions** -- NonAlloc physics queries, cached references, and Input System usage followed unity-game-developer's built-in guidance

## Next Steps

- [Prototype a Game Mechanic](../../howto/prototype-game-mechanic/) -- rapid iteration workflow for testing mechanics before full implementation
- [Set Up a Unity Project](../../howto/set-up-unity-project/) -- complete Unity project architecture beyond a single mechanic
- [Choosing the Right Agent](../../explanation/choosing-the-right-agent/) -- when to use game-designer vs game-developer vs engine-specific agents
- [Agent Reference](../../reference/agents/) -- full specifications for all five agents
