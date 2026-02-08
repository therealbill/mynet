---
title: "How to Install a Plugin"
description: "Install a Mynet plugin into your Claude Code project from the marketplace or a local directory, verify it works, and remove it when no longer needed."
weight: 1
---

# How to Install a Plugin

Install a Mynet plugin so its agents, skills, and commands are available in your Claude Code session.

## Prerequisites

- Claude Code installed and configured
- A Claude Code project initialized (a working directory with code)
- The Mynet marketplace added to Claude Code (run `/plugin marketplace add therealbill/mynet` if not already registered)

## Steps

### 1. Browse Available Plugins

Open the plugin manager and browse available plugins:

```
/plugin
```

Navigate to the **Discover** tab to see all plugins from registered marketplaces. Each entry lists a plugin's name, description, and keywords. Identify the plugin you want to install. For example, `code-quality` provides code review and testing agents.

### 2. Install from the Marketplace

Install a plugin directly using the `/plugin` command:

```
/plugin install code-quality@mynet --scope project
```

The `--scope project` flag makes the plugin available to anyone who works in this project directory. Omit it to install for your user only.

The plugin loads immediately -- no restart is needed.

### 3. Install from a Local Directory

If you have a plugin directory outside the marketplace, install it by path:

```
claude plugin install /absolute/path/to/my-custom-plugin
```

The directory must contain a `.claude-plugin/plugin.json` manifest at minimum.

### 4. List Installed Plugins

Confirm the plugin loaded by checking the Installed tab:

```
/plugin
```

Navigate to the **Installed** tab to see all active plugins. You can also invoke an agent from the installed plugin to verify:

```
@code-reviewer Review the current file for issues.
```

If the agent responds, the plugin is installed correctly.

### 5. Remove a Plugin

Uninstall a plugin using the `/plugin` command:

```
/plugin uninstall code-quality@mynet
```

The plugin's agents, skills, and commands will no longer be available.

## Verify It Works

After installation, run a quick check:

- [ ] The plugin's agents respond when invoked by name (e.g., `@code-reviewer`)
- [ ] The plugin's skills trigger on matching patterns
- [ ] The plugin's slash commands appear in the command palette (if the plugin defines commands)

## Troubleshooting

**Agent not found after installation:**

- Verify the plugin appears in the `/plugin` Installed tab
- Verify the plugin directory contains `.claude-plugin/plugin.json`
- Try reinstalling with `/plugin install plugin-name@mynet`

**Plugin loads but skills do not trigger:**

- Check that the skill's trigger patterns in `SKILL.md` match your input
- Verify the skill directory structure follows `skills/<skill-name>/SKILL.md`

**Permission errors when loading a plugin:**

- Ensure the plugin directory and all files are readable by your user
- On macOS, check that the directory is not quarantined by Gatekeeper

## Next Steps

- [How to Find the Right Plugin for Your Task](../../howto/find-right-plugin/) -- search and match plugins to your workflow
- [How to Use Plugins Together](../../howto/use-plugins-together/) -- combine agents from multiple plugins
- [Plugin JSON Reference](../../reference/plugin-json/) -- full plugin manifest schema
- [Component Conventions](../../reference/component-conventions/) -- agent, skill, command, and template formats
