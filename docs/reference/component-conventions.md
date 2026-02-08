---
title: "Component Conventions"
description: "Specifications for agent, skill, command, and template formats"
weight: 4
---

# Component Conventions

**Stability:** Stable

This document specifies the file format, location, and front matter schema for each plugin component type: agents, skills, commands, and templates.

## Agents

### File Location

```
<plugin>/agents/<agent-name>.md
```

Each agent is a single Markdown file within the plugin's `agents/` directory.

### Format

Markdown with YAML front matter. The body after the front matter contains the agent's system prompt.

### Front Matter Fields

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `name` | string | Yes | -- | Agent identifier. Uses kebab-case. |
| `description` | string | Yes | -- | Agent description, including trigger conditions and usage examples. |
| `model` | string | No | -- | Model selection: `"opus"`, `"sonnet"`, or `"haiku"`. |
| `color` | string | No | -- | Display color hint for UI rendering. |
| `tools` | string[] | No | `[]` | Tools accessible to the agent. Valid values include: `"Read"`, `"Grep"`, `"Glob"`, `"Bash"`, `"Write"`, `"Edit"`, `"WebSearch"`, `"WebFetch"`, `"Task"`. |

### Description Format

The `description` field supports XML-style `<example>` blocks to define trigger conditions.

```yaml
description: |
  Agent description text.

  <example>
  Context: [situation description]
  user: "[user message]"
  assistant: "[assistant response using the agent]"
  <commentary>
  [Why this triggers the agent]
  </commentary>
  </example>
```

### Agent File Example

```markdown
---
name: my-agent
description: |
  Agent that handles a specific domain.

  <example>
  Context: User is working on a project
  user: "Analyze this code"
  assistant: "(delegates to my-agent)"
  <commentary>
  Domain-specific request triggers this agent.
  </commentary>
  </example>
model: sonnet
tools:
  - Read
  - Grep
  - Glob
  - Bash
---

System prompt content goes here. This is the instruction set
provided to the agent when it is invoked.
```

---

## Skills

### File Location

```
<plugin>/skills/<skill-name>/SKILL.md
```

Each skill is a `SKILL.md` file inside a named subdirectory within the plugin's `skills/` directory. The subdirectory name matches the skill identifier.

### Format

Markdown with YAML front matter. The body after the front matter contains the skill's instruction content.

### Front Matter Fields

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `name` | string | Yes | -- | Human-readable skill name. |
| `description` | string | Yes | -- | Detailed description of the skill's purpose and trigger conditions. |
| `version` | string | No | -- | Semver version string. |

### Skill File Example

```markdown
---
name: My Skill
description: |
  Activates when the user requests a specific operation.
  Covers the following areas: area-one, area-two.
version: 1.0.0
---

Skill instruction content goes here. This defines the
behavior and knowledge applied when the skill is triggered.
```

---

## Commands

### File Location

```
<plugin>/commands/<command-name>.md
```

Each command is a single Markdown file within the plugin's `commands/` directory.

### Format

Markdown with YAML front matter. The body after the front matter contains the command's implementation instructions.

### Front Matter Fields

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| `name` | string | Yes | -- | Command name. Invoked as `/<name>` in the CLI. |
| `description` | string | Yes | -- | Brief description of the command's function. |
| `arguments` | object[] | No | `[]` | Array of argument definitions. |
| `arguments[].name` | string | Yes (if arguments present) | -- | Argument name. |
| `arguments[].description` | string | Yes (if arguments present) | -- | Argument description. |
| `arguments[].required` | boolean | No | `false` | Whether the argument is required for invocation. |

### Command File Example

```markdown
---
name: my-command
description: Performs a specific operation on the target
arguments:
  - name: target
    description: The target resource to operate on
    required: true
  - name: format
    description: Output format (json, yaml, text)
    required: false
---

Implementation instructions for the command go here.
This defines what the command does when invoked.
```

### Invocation

Commands are invoked via the CLI slash-command interface as `/<command-name>`. Arguments, when defined, are passed after the command name.

---

## Templates

### File Location

```
<plugin>/templates/<template-name>.tmpl
```

Other file extensions are permitted depending on the template's target output format.

### Format

Templates are file templates used for code generation. There is no standardized front matter schema for templates. The format and syntax of each template depend on the target output language or framework.

### Template Example

```
<plugin>/templates/
  workflow.go.tmpl
  activity.go.tmpl
  config.yaml.tmpl
  Makefile.tmpl
```

---

## Directory Structure Summary

A fully populated plugin directory contains the following structure.

```
<plugin-name>/
  .claude-plugin/
    plugin.json
  agents/
    <agent-name>.md
  skills/
    <skill-name>/
      SKILL.md
  commands/
    <command-name>.md
  templates/
    <template-name>.tmpl
  docs/
```

All subdirectories (`agents/`, `skills/`, `commands/`, `templates/`, `docs/`) are optional. A plugin contains only the component directories relevant to its contents.

## See Also

- [Plugin Catalog]({{< ref "plugin-catalog" >}}) -- complete listing of all plugins
- [Plugin JSON]({{< ref "plugin-json" >}}) -- per-plugin manifest specification
- [Marketplace Manifest]({{< ref "marketplace-manifest" >}}) -- root marketplace manifest specification
