---
title: "How to Find the Right Plugin for Your Task"
description: "Search the Mynet marketplace catalog by keyword, match plugins to your task type, and inspect plugin components before installing."
weight: 2
---

# How to Find the Right Plugin for Your Task

Identify which Mynet plugin provides the agents, skills, or commands you need for a specific task.

## Prerequisites

- The Mynet marketplace added to Claude Code (run `/plugin marketplace add therealbill/mynet` if not already registered)
- Familiarity with the plugin component types (agents, skills, commands) -- see the [Marketplace Explanation](../../explanation/) for background

## Steps

### 1. Browse the Plugin Catalog

Open the plugin manager and navigate to the Discover tab:

```
/plugin
```

The **Discover** tab lists all registered plugins with their name, description, and keywords. Scroll through the list or use the search/filter functionality to narrow results.

### 2. Search by Keyword

Use the Discover tab's search to filter plugins by keyword. Common search terms and what they match:

| Search term | Matches |
|---|---|
| `review` | code-quality |
| `deploy` | devops-and-infra |
| `react` | web-development |
| `go` | backend-development, programming-languages, desktop-development |
| `research` | research |
| `documentation` | developer-tools, diataxis-docs, utility-skills |

### 3. Match Your Task to a Plugin

Use this mapping to go from task type to the recommended plugin:

| Task | Plugin | Key Agents |
|---|---|---|
| Code review and architecture | `code-quality` | code-reviewer, architect-review |
| Writing tests or fixing failures | `code-quality` | test-writer-fixer |
| Frontend / React / Next.js development | `web-development` | (framework-specific agents) |
| Backend APIs and Go services | `backend-development` | (Go, SQL, API agents) |
| Deployment and CI/CD pipelines | `devops-and-infra` | (GitHub Actions, monitoring agents) |
| Language-specific coding (Rust, C++, TypeScript) | `programming-languages` | rust-pro, cpp-pro, typescript-pro |
| Research and literature review | `research` | comprehensive-researcher, report-generator |
| Documentation creation | `diataxis-docs` | doc-howto-writer, doc-tutorial-writer |
| Git workflow and prototyping | `developer-tools` | git-workflow-manager, rapid-prototyper |
| Hugo static sites | `hugo-repo` | hugo-site-architect, hugo-build-doctor |
| Makefile authoring | `gnu-make` | (Make best-practices skills) |

### 4. Inspect Plugin Components

Before installing, check what a plugin contains. You can inspect components from the `/plugin` UI by selecting a plugin in the Discover tab to see its agents, skills, and commands.

You can also browse the marketplace repository directly for detailed component definitions:

```bash
# List agents in a plugin
ls code-quality/agents/

# Read a specific agent definition
cat code-quality/agents/code-reviewer.md
```

The agent's YAML frontmatter shows its name, description, required tools, and model. The body contains the system prompt that defines the agent's behavior.

### 5. Check Plugin Version and Metadata

Plugin metadata (version, keywords, description) is visible in the Discover tab when you select a plugin. For detailed inspection, read the plugin manifest:

```bash
cat code-quality/.claude-plugin/plugin.json
```

```json
{
  "name": "code-quality",
  "version": "1.0.0",
  "description": "Code review, testing, accessibility, and architectural quality agents",
  "keywords": ["code-review", "testing", "accessibility", "architecture", "quality"]
}
```

## Verify It Works

After identifying the right plugin:

- [ ] The plugin's description matches your task
- [ ] At least one agent or skill addresses your specific need
- [ ] You have reviewed the agent's frontmatter to confirm it fits your use case

## Troubleshooting

**No plugins match your keyword:**

- Try broader terms: "web" instead of "nextjs", "code" instead of "linting"
- Browse the full catalog manually -- descriptions often cover capabilities not listed in keywords

**Multiple plugins seem to fit:**

- Prefer the domain-specific plugin over general-purpose ones (e.g., `programming-languages` for Rust work over `code-quality`)
- Install both and use agents from each as needed -- see [How to Use Plugins Together](../use-plugins-together/)

**Plugin exists but has no agents for your need:**

- Check the `skills/` directory -- the capability may be a skill rather than an agent
- Check the `commands/` directory for slash commands

## Next Steps

- [How to Install a Plugin](../../howto/install-a-plugin/) -- install the plugin you found
- [How to Use Plugins Together](../../howto/use-plugins-together/) -- combine agents across plugins
- [Plugin Catalog](../../reference/plugin-catalog/) -- complete listing of all plugins
- [Marketplace Manifest](../../reference/marketplace-manifest/) -- marketplace.json schema
