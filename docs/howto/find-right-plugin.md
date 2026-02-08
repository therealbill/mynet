---
title: "How to Find the Right Plugin for Your Task"
description: "Search the Mynet marketplace catalog by keyword, match plugins to your task type, and inspect plugin components before installing."
weight: 2
---

# How to Find the Right Plugin for Your Task

Identify which Mynet plugin provides the agents, skills, or commands you need for a specific task.

## Prerequisites

- The Mynet marketplace repository cloned locally
- Familiarity with the plugin component types (agents, skills, commands) -- see the [Marketplace Explanation](../../explanation/) for background

## Steps

### 1. Browse the Plugin Catalog

Open the marketplace manifest to see all registered plugins:

```bash
cat .claude-plugin/marketplace.json | jq '.plugins[] | {name, description, keywords}'
```

Expected output (abbreviated):

```json
{
  "name": "code-quality",
  "description": "Code review, testing, accessibility, and architectural quality agents",
  "keywords": ["code-review", "testing", "accessibility", "architecture", "quality"]
}
{
  "name": "devops-and-infra",
  "description": "DevOps automation, CI/CD, monitoring, and infrastructure agents",
  "keywords": ["devops", "ci-cd", "github-actions", "prometheus", "monitoring"]
}
```

### 2. Search by Keyword

Filter plugins using a keyword that matches your task:

```bash
cat .claude-plugin/marketplace.json | jq '.plugins[] | select(.keywords[] | test("review")) | .name'
```

This returns every plugin whose keywords match "review". Common search terms:

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

Before installing, check what a plugin actually contains:

```bash
# List agents
ls code-quality/agents/

# List skills
ls code-quality/skills/

# Read a specific agent definition
cat code-quality/agents/code-reviewer.md
```

The agent's YAML frontmatter shows its name, description, required tools, and model. The body contains the system prompt that defines the agent's behavior.

### 5. Check Plugin Version and Metadata

Read the plugin manifest for version and compatibility information:

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

**jq command not found:**

- Install jq: `brew install jq` (macOS) or `apt install jq` (Linux)
- Alternatively, open the JSON file directly in your editor

## Next Steps

- [How to Install a Plugin]({{< ref "howto/install-a-plugin" >}}) -- install the plugin you found
- [How to Use Plugins Together]({{< ref "howto/use-plugins-together" >}}) -- combine agents across plugins
- [Plugin Catalog]({{< ref "reference/plugin-catalog" >}}) -- complete listing of all plugins
- [Marketplace Manifest]({{< ref "reference/marketplace-manifest" >}}) -- marketplace.json schema
