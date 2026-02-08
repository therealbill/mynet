---
title: "Set Up a Unity Project"
description: "Configure a Unity project with proper architecture using unity-game-developer for rendering pipeline selection, ScriptableObject data, and project structure"
weight: 2
---

# Set Up a Unity Project

Set up a Unity project with the correct rendering pipeline, folder structure, data architecture, input configuration, and assembly definitions using the unity-game-developer agent.

## Prerequisites

- Claude Code with the game-development plugin installed
- Unity 2022.3 LTS or later installed (Unity Hub recommended)
- A new Unity project created from the Unity Hub (any template -- the agent reconfigures it)

## Steps

### 1. Choose a Rendering Pipeline

Ask unity-game-developer for a pipeline recommendation:

```
I'm starting a Unity project for [describe your game and target platforms].
Which rendering pipeline should I use?
```

The agent evaluates based on your target:

- **URP (Universal Render Pipeline):** Default recommendation for mobile, most indie projects, and any project targeting a wide hardware range. URP is the pipeline the agent selects unless you have a specific reason for HDRP.
- **HDRP (High Definition Render Pipeline):** Only for PC and current-gen console projects that require advanced lighting, volumetric fog, or high-fidelity materials. HDRP does not support mobile.
- **Built-in pipeline:** Only for legacy projects already using it. Do not start new projects on built-in.

If your project targets mobile, the agent selects URP and configures it with baked lighting defaults, since real-time lights are expensive on mobile GPUs.

### 2. Configure Folder Structure

Ask the agent to set up the project structure:

```
Set up the folder structure for this project
```

unity-game-developer generates a folder hierarchy that separates concerns:

```
Assets/
  _Project/
    Art/
      Materials/
      Models/
      Textures/
      Animations/
    Audio/
      Music/
      SFX/
    Data/              # ScriptableObject assets
    Prefabs/
      Characters/
      Environment/
      UI/
    Scenes/
    Scripts/
      Runtime/
        Core/
        Gameplay/
        UI/
      Editor/
      Tests/
    Settings/          # Input actions, render pipeline assets
  Plugins/             # Third-party packages that aren't UPM
```

The `_Project` prefix sorts the project folder above Unity's default folders in the Project window. Third-party assets go in `Plugins/` to keep them separate from project code.

### 3. Set Up Assembly Definitions

Ask the agent to configure assembly definitions:

```
Add assembly definitions for the project
```

unity-game-developer creates `.asmdef` files that split compilation into independent assemblies:

- **`Game.Runtime`** -- core gameplay code, no editor dependencies
- **`Game.Editor`** -- custom inspectors and editor tools, references `Game.Runtime`
- **`Game.Tests`** -- test assemblies referencing `Game.Runtime` and Unity's test framework

Assembly definitions reduce recompilation time. Changing a script in `Game.Editor` does not trigger recompilation of `Game.Runtime`. On projects with more than a few dozen scripts, this saves significant iteration time.

### 4. Configure ScriptableObject Data Architecture

Ask the agent to set up the data layer:

```
Set up a ScriptableObject-based data architecture for game configuration
```

unity-game-developer creates base patterns for data-driven design:

- **Configuration ScriptableObjects** for game settings (player stats, enemy definitions, item data) with `[CreateAssetMenu]` attributes for easy asset creation
- **Event ScriptableObjects** for decoupled communication between systems (a `GameEvent` ScriptableObject that MonoBehaviours can raise and listen to without direct references)
- **Runtime Sets** for tracking active entities (a `RuntimeSet<T>` ScriptableObject that entities register with on enable and deregister on disable)

This architecture follows unity-game-developer's core principle: ScriptableObjects hold data, MonoBehaviours define behavior. Game designers modify ScriptableObject assets in the Inspector without touching code.

### 5. Configure the Input System

Ask the agent to set up input:

```
Configure the new Input System for my project
```

unity-game-developer installs and configures the new Input System package:

- Creates a `PlayerInputActions` asset with default action maps (Player, UI, Menu)
- Configures common actions (Move, Jump, Attack, Interact, Pause) with keyboard, gamepad, and optional touch bindings
- Generates a C# wrapper class for type-safe input access
- Sets the project's Active Input Handling to "Input System Package (New)" in Player Settings

The agent avoids the legacy `Input` class entirely. The new Input System supports rebinding, multiple input devices, and local multiplayer without custom code.

### 6. Verification

Run through this checklist to confirm the project is set up correctly:

- [ ] The correct rendering pipeline is active (check **Edit > Project Settings > Graphics**)
- [ ] Folder structure exists under `Assets/_Project/` with all expected subdirectories
- [ ] Assembly definitions compile without errors (check the Console for warnings)
- [ ] A sample ScriptableObject can be created via the Assets menu (**Assets > Create** shows your custom types)
- [ ] The Input System is active (check **Edit > Project Settings > Player > Active Input Handling**)
- [ ] Entering Play mode triggers no errors or warnings related to project configuration

## Troubleshooting

**"Input System" package not found:**

Open the Package Manager (**Window > Package Manager**), switch to Unity Registry, and search for "Input System." If it does not appear, your Unity version may be too old. Unity 2019.4+ supports the new Input System.

**Assembly definition errors about missing references:**

Ensure `Game.Editor` references `Game.Runtime` in its `.asmdef` inspector. Editor assemblies must explicitly reference runtime assemblies -- Unity does not infer this.

**ScriptableObject assets lose data on Play mode exit:**

This is expected Unity behavior. ScriptableObject changes made during Play mode persist only if you mark them dirty or use a custom editor tool. For runtime-only state, use MonoBehaviours or plain C# classes instead.

## See Also

- [Getting Started](../../tutorials/getting-started/) -- tutorial building a mechanic inside a Unity project
- [Agent Reference](../../reference/agents/) -- unity-game-developer specification
- [Architecture](../../explanation/architecture/) -- why engine-specific agents exist alongside the generalist game-developer
