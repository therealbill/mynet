---
title: "Prototype a Game Mechanic"
description: "Rapidly prototype a game mechanic using game-designer for specification and game-developer for engine-agnostic implementation"
weight: 1
---

# Prototype a Game Mechanic

Rapidly define, validate, and implement a minimum viable game mechanic using game-designer for the design specification and game-developer for engine-agnostic implementation guidance.

## Prerequisites

- Claude Code with the game-development plugin installed
- A game concept with at least one mechanic you want to test
- A spreadsheet application (Google Sheets, Excel, or similar) for simulation
- A game engine installed (Unity, Unreal, or Godot) for the coding phase

## Steps

### 1. Define the Mechanic

Describe the mechanic to Claude Code with enough context for game-designer to activate:

```
I need to design a grappling hook mechanic for a 2D action game.
The player should be able to swing from attachment points and launch
themselves. It needs to feel physics-driven but controllable.
```

game-designer asks clarifying questions about the mechanic's role in the core loop, the target audience, and how it interacts with other systems. Answer these to narrow the design space.

The agent produces a mechanic specification containing:

- **Core parameters:** rope length, swing speed, launch velocity multiplier, attachment range, input window
- **Parameter ranges:** minimum and maximum values for each parameter, with a recommended starting point
- **Interaction rules:** what happens on attach, during swing, on release, on collision
- **Edge cases:** what if the player attaches while falling, while already attached, or at extreme angles

### 2. Validate with a Paper Prototype or Spreadsheet

Before writing any code, game-designer recommends validating the numbers. Ask:

```
Help me set up a spreadsheet simulation to test the grappling hook physics
```

The agent defines a spreadsheet model that calculates:

- Swing arc trajectories at different rope lengths
- Launch velocities at different release angles
- Time-to-target for gap-crossing scenarios at varying distances
- Whether the parameter ranges produce fun or frustrating outcomes

Plug in the recommended starting values and test several scenarios. Adjust parameters until the spreadsheet shows the behavior you want. This catches problems like launch velocities that overshoot every platform or rope lengths that make small gaps trivial.

This step takes 10-15 minutes and saves hours of in-engine iteration.

### 3. Get Engine-Agnostic Architecture

Once the parameters feel right on paper, ask for implementation guidance:

```
How should I architect this grappling hook mechanic? I want it to work
with any physics engine. Give me the architecture, not engine-specific code.
```

Claude Code routes this to the game-developer agent. It provides:

- **State machine:** Idle, Aiming, Attached, Swinging, Launching -- with defined transitions
- **Data-driven configuration:** all parameters externalized into a data file or configuration object
- **Physics approach:** constraint-based (distance joint) vs force-based (apply tangential force each frame), with tradeoffs for each
- **Separation of simulation and presentation:** the swing physics runs independently of the visual rope rendering

game-developer recommends composition over inheritance: the grappling hook is a component attached to the player entity, not a subclass of the player.

### 4. Implement the Minimum Viable Version

Ask game-developer for the simplest possible implementation that tests the core feel:

```
What is the absolute minimum I need to implement to test whether
the grappling hook feels good?
```

The agent strips the mechanic down to:

- A point-to-point distance constraint (no rope rendering, no aiming UI)
- Fixed attachment points placed manually in the scene (no dynamic detection)
- Release launches the player at their current velocity (no velocity multiplier tuning yet)
- Keyboard input only (no gamepad, no touch)

This minimum viable mechanic can be built in 30-60 minutes in any engine. The goal is to answer one question: does swinging and launching feel satisfying at a fundamental physics level?

### 5. Test and Iterate on Parameters

Play the minimum viable mechanic and take notes:

- Does the swing arc feel too wide or too tight? Adjust rope length.
- Does the launch feel too weak or too powerful? Adjust velocity handling.
- Is attachment too easy or too hard to trigger? Adjust attachment range.
- Does the player feel in control during the swing? Adjust damping.

Feed your observations back to game-designer:

```
The swing feels too floaty and the launch is too powerful. Players overshoot
every platform. The rope length feels right though.
```

The agent adjusts the parameter specification based on your playtest feedback, explains why the observed behavior occurred, and recommends specific numerical changes.

### 6. Decide: Commit or Kill

After 2-3 iteration rounds, make a decision:

- **Commit:** The core feel is right. Move to full production implementation with your engine-specific agent (unity-game-developer or unreal-engine-developer).
- **Pivot:** The feel is close but the mechanic needs a different approach. Return to step 1 with revised constraints.
- **Kill:** The mechanic does not produce the intended experience. Stop investing time and move to the next mechanic idea.

game-designer's prototyping philosophy is explicit: test mechanics with minimum viable implementations before committing to full production. A killed prototype that took 2 hours is cheaper than a fully produced mechanic that gets cut after 2 weeks.

## Troubleshooting

**game-designer produces too much detail too early:**

Ask the agent to focus on the minimum viable version. Remind it that you are prototyping, not producing.

**Spreadsheet simulation does not match in-engine behavior:**

The spreadsheet uses simplified physics. Expect directional accuracy (the parameter trends are correct) but not exact numerical parity. Use the spreadsheet to eliminate obviously wrong parameter ranges, then fine-tune in-engine.

**The mechanic feels different at different frame rates:**

Ask game-developer to verify that your implementation uses fixed-timestep physics (FixedUpdate in Unity, Tick with fixed delta in Unreal) rather than frame-dependent calculations.

## See Also

- [Getting Started](../../tutorials/getting-started/) -- full tutorial walking through design and implementation
- [Architecture](../../explanation/architecture/) -- why design and implementation are separate agents
- [Choosing the Right Agent](../../explanation/choosing-the-right-agent/) -- when to use game-designer vs game-developer
