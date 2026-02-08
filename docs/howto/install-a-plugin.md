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
- The Mynet marketplace repository cloned locally

## Steps

### 1. Browse Available Plugins

Open the marketplace manifest to see all available plugins:

```bash
cat .claude-plugin/marketplace.json
```

Each entry lists a plugin's name, description, and keywords. Identify the plugin you want to install. For example, `code-quality` provides code review and testing agents.

### 2. Install from the Marketplace

Add the plugin to your project's Claude Code configuration. In your project's `.claude/settings.json`, add the plugin source path:

```json
{
  "plugins": [
    {
      "name": "code-quality",
      "source": "/path/to/claude-plugins/code-quality"
    }
  ]
}
```

Replace `/path/to/claude-plugins/` with the absolute path to your local Mynet repository clone.

### 3. Install from a Local Directory

If you have a plugin directory outside the marketplace, point directly to it:

```json
{
  "plugins": [
    {
      "name": "my-custom-plugin",
      "source": "/absolute/path/to/my-custom-plugin"
    }
  ]
}
```

The directory must contain a `.claude-plugin/plugin.json` manifest at minimum.

### 4. Restart Your Claude Code Session

Close and reopen Claude Code, or start a new session in your project directory. Plugins load at session startup.

### 5. List Installed Plugins

Confirm the plugin loaded by checking which agents, skills, and commands are available. Invoke an agent from the installed plugin to verify:

```
@code-reviewer Review the current file for issues.
```

If the agent responds, the plugin is installed correctly.

### 6. Remove a Plugin

Delete the plugin entry from your project's `.claude/settings.json`:

```json
{
  "plugins": []
}
```

Restart your Claude Code session. The plugin's agents, skills, and commands will no longer be available.

## Verify It Works

After installation, run a quick check:

- [ ] The plugin's agents respond when invoked by name (e.g., `@code-reviewer`)
- [ ] The plugin's skills trigger on matching patterns
- [ ] The plugin's slash commands appear in the command palette (if the plugin defines commands)

## Troubleshooting

**Agent not found after installation:**

- Confirm the `source` path in your settings is an absolute path, not relative
- Verify the plugin directory contains `.claude-plugin/plugin.json`
- Restart your Claude Code session -- plugins load only at startup

**Plugin loads but skills do not trigger:**

- Check that the skill's trigger patterns in `SKILL.md` match your input
- Verify the skill directory structure follows `skills/<skill-name>/SKILL.md`

**Permission errors when loading a plugin:**

- Ensure the plugin directory and all files are readable by your user
- On macOS, check that the directory is not quarantined by Gatekeeper

## Next Steps

- [How to Find the Right Plugin for Your Task]({{< ref "howto/find-right-plugin" >}}) -- search and match plugins to your workflow
- [How to Use Plugins Together]({{< ref "howto/use-plugins-together" >}}) -- combine agents from multiple plugins
- [Plugin JSON Reference]({{< ref "reference/plugin-json" >}}) -- full plugin manifest schema
- [Component Conventions]({{< ref "reference/component-conventions" >}}) -- agent, skill, command, and template formats
