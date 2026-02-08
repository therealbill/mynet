# Plugin Documentation Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Create Diataxis-compliant documentation for all 18 plugins (skip timelord) plus the marketplace root, with Hugo-ready front matter and directory structure.

**Architecture:** Each plugin gets a `docs/` directory with four Diataxis subdirectories (tutorials, howto, reference, explanation), each with `_index.md` section pages. Depth scales by tier: Compact (1-3 components), Standard (3-6), Comprehensive (7+). The marketplace root `docs/` covers cross-plugin concepts. All files use Hugo-compatible front matter (title, description, weight). Use the Write tool for all file creation.

**Tech Stack:** Markdown with YAML front matter. No code dependencies.

**Design document:** Brainstorming session from 2026-02-08. Timelord docs at `timelord/docs/content/` serve as the style reference.

**Reference files for conventions:**

- Timelord overview: `timelord/docs/content/_index.md`
- Timelord reference style: `timelord/docs/content/reference/agent-reference.md`
- Timelord how-to style: `timelord/docs/content/how-to/scale-workers.md`
- Timelord explanation style: `timelord/docs/content/explanation/how-temporal-works.md`

**Execution note:** Use the Write tool (not echo/cat heredocs in Bash) for all file creation. Use Bash only for mkdir and git commits.

**Tier assignments:**

| Tier | Plugins |
|------|---------|
| Compact (1-3 components) | desktop-development, mcp-development, ai-development, utility-skills |
| Standard (3-6 components) | backend-development, cli-development, developer-tools, utility-agents, mobile-development, devops-and-infra, code-quality, research, gnu-make, game-development, web-development |
| Comprehensive (7+ components) | hugo-repo, programming-languages, diataxis-docs |
| Root | marketplace (root docs/) |

**Cross-reference map:**

- code-quality ↔ programming-languages
- backend-development ↔ devops-and-infra
- cli-development ↔ programming-languages
- developer-tools ↔ diataxis-docs
- hugo-repo ↔ diataxis-docs
- research ↔ ai-development
- gnu-make ↔ developer-tools
- mcp-development ↔ any plugin using MCP tools

---

## Group A: Marketplace Root Documentation

### Task 1: Marketplace root docs scaffolding and overview

**Files:**

- Create: `docs/_index.md`
- Create: `docs/tutorials/_index.md`
- Create: `docs/howto/_index.md`
- Create: `docs/reference/_index.md`
- Create: `docs/explanation/_index.md`

**Step 1:** Create directories.

```bash
mkdir -p docs/tutorials docs/howto docs/reference docs/explanation
```

**Step 2:** Create `docs/_index.md`:

```markdown
---
title: "Mynet Plugin Marketplace"
description: "Documentation for the Mynet Claude Code plugin marketplace"
type: docs
---

# Mynet Plugin Marketplace

A curated collection of Claude Code plugins for software engineering, research, documentation, and development operations.

Mynet organizes specialized Claude agents, skills, commands, and templates into domain-focused plugins that can be installed independently or combined for comprehensive coverage.

## Quick Start

### Installation

```bash
# Install a single plugin
claude plugin install code-quality

# Or install from local directory
claude plugin install ./code-quality
```

### Using a Plugin

Once installed, plugin components activate automatically:

- **Agents** are dispatched when your question matches their expertise
- **Skills** trigger when you ask about topics they cover
- **Commands** are invoked with `/command-name`
- **Templates** provide starter files for common patterns

## Documentation

### Tutorials

Step-by-step guides for getting started:

- [Your First Plugin](tutorials/getting-started/) — Install a plugin and use its components

### How-To Guides

Task-oriented guides:

- [Install a Plugin](howto/install-a-plugin/)
- [Find the Right Plugin](howto/find-right-plugin/)
- [Use Plugins Together](howto/use-plugins-together/)

### Reference

Technical specifications:

- [Plugin Catalog](reference/plugin-catalog/) — All available plugins
- [Marketplace Manifest](reference/marketplace-manifest/) — marketplace.json format
- [Plugin JSON](reference/plugin-json/) — plugin.json format
- [Component Conventions](reference/component-conventions/) — Agent, skill, command, template formats

### Explanation

Conceptual documentation:

- [Marketplace Architecture](explanation/marketplace-architecture/) — How the marketplace works
- [Plugin Design Philosophy](explanation/plugin-design-philosophy/) — Why this structure
- [Licensing](explanation/licensing/) — License model
```

**Step 3:** Create section index pages:

`docs/tutorials/_index.md`:
```markdown
---
title: "Tutorials"
description: "Step-by-step guides for learning the Mynet marketplace"
weight: 1
---

Tutorials walk you through using the marketplace from scratch. Follow them in order to build practical experience.
```

`docs/howto/_index.md`:
```markdown
---
title: "How-To Guides"
description: "Task-oriented guides for common marketplace operations"
weight: 2
---

How-to guides solve specific problems. Each guide has one goal, numbered steps, and verification.
```

`docs/reference/_index.md`:
```markdown
---
title: "Reference"
description: "Technical specifications for marketplace formats and conventions"
weight: 3
---

Reference documentation covers the exact formats and fields for marketplace components. No opinions — just specifications.
```

`docs/explanation/_index.md`:
```markdown
---
title: "Explanation"
description: "Conceptual documentation about marketplace design and architecture"
weight: 4
---

Explanation documents help you understand why the marketplace works the way it does, the design decisions behind it, and how to think about plugin architecture.
```

**Step 4:** Commit.

```bash
git add docs/_index.md docs/tutorials/_index.md docs/howto/_index.md docs/reference/_index.md docs/explanation/_index.md
git commit -m "add marketplace root docs scaffolding and overview"
```

---

### Task 2: Marketplace root tutorial

**Files:**

- Create: `docs/tutorials/getting-started.md`

**Content:** A ~1000 word tutorial walking through installing the code-quality plugin, triggering the code-reviewer agent, and understanding what happened. Use the code-quality plugin as the example because code review is universally relatable.

**Key sections:** Prerequisites (Claude Code installed) → Install the plugin → Trigger an agent → Understand the output → What you learned → Next steps.

**Step 1:** Create the file with the Write tool.

**Step 2:** Commit.

```bash
git add docs/tutorials/getting-started.md
git commit -m "add marketplace getting-started tutorial"
```

---

### Task 3: Marketplace root how-to guides

**Files:**

- Create: `docs/howto/install-a-plugin.md`
- Create: `docs/howto/find-right-plugin.md`
- Create: `docs/howto/use-plugins-together.md`

**Content guidelines:**

- `install-a-plugin.md`: Steps for installing from marketplace, from local directory, verifying installation, listing installed plugins, removing a plugin. ~400 words.
- `find-right-plugin.md`: How to browse the catalog, search by keyword, match plugins to your task. Reference the plugin catalog. ~400 words.
- `use-plugins-together.md`: How agents from different plugins collaborate (e.g., research agent feeds into report-generator, code-quality reviews what programming-languages writes). ~500 words.

**Step 1:** Create all three files with the Write tool.

**Step 2:** Commit.

```bash
git add docs/howto/
git commit -m "add marketplace how-to guides"
```

---

### Task 4: Marketplace root reference docs

**Files:**

- Create: `docs/reference/plugin-catalog.md`
- Create: `docs/reference/marketplace-manifest.md`
- Create: `docs/reference/plugin-json.md`
- Create: `docs/reference/component-conventions.md`

**Content guidelines:**

- `plugin-catalog.md`: Table of all 19 plugins with name, description, component counts (agents/skills/commands/templates), keywords, version. Extract data from `.claude-plugin/marketplace.json` and each plugin's component directories.
- `marketplace-manifest.md`: Specification of the root `.claude-plugin/marketplace.json` format. Fields, types, examples. Pure specification — no advice.
- `plugin-json.md`: Specification of `.claude-plugin/plugin.json` format. Fields (name, version, description, author, repository, license, keywords), types, required vs optional.
- `component-conventions.md`: Specification of agent frontmatter (name, description with examples, model, color, tools), skill frontmatter (name, description, version), command frontmatter (name, description, arguments), template conventions.

**Step 1:** Create all four files with the Write tool.

**Step 2:** Commit.

```bash
git add docs/reference/
git commit -m "add marketplace reference documentation"
```

---

### Task 5: Marketplace root explanation docs

**Files:**

- Create: `docs/explanation/marketplace-architecture.md`
- Create: `docs/explanation/plugin-design-philosophy.md`
- Create: `docs/explanation/licensing.md`

**Content guidelines:**

- `marketplace-architecture.md`: How the marketplace works — directory structure, component discovery, how Claude Code finds and loads plugins. Why a flat marketplace with independent plugins rather than a monolithic system. ~600 words.
- `plugin-design-philosophy.md`: Why plugins have agents, skills, commands, and templates as distinct types. When to create each type. The agent specialization model — why narrow agents outperform general ones. The skill trigger model — how descriptions drive activation. ~800 words.
- `licensing.md`: The aggregation-prohibiting license — what it allows, what it prevents, why. Individual plugin redistribution rules. ~400 words. Reference the root LICENSE file.

**Step 1:** Create all three files with the Write tool.

**Step 2:** Commit.

```bash
git add docs/explanation/
git commit -m "add marketplace explanation documentation"
```

---

## Group B: Tier 1 (Compact) Plugin Documentation

Each Tier 1 plugin gets ~8 files: _index.md, 4 section indexes, 1 tutorial, 1 how-to, 1 reference, 1 explanation.

### Task 6: desktop-development docs (1 agent: electron-go-pro)

**Files:**

- Create: `desktop-development/docs/_index.md`
- Create: `desktop-development/docs/tutorials/_index.md`
- Create: `desktop-development/docs/tutorials/getting-started.md`
- Create: `desktop-development/docs/howto/_index.md`
- Create: `desktop-development/docs/howto/set-up-electron-go-project.md`
- Create: `desktop-development/docs/reference/_index.md`
- Create: `desktop-development/docs/reference/agents.md`
- Create: `desktop-development/docs/explanation/_index.md`
- Create: `desktop-development/docs/explanation/architecture.md`

**Content focus:**

- Overview: Electron+Go hybrid desktop app development for macOS
- Tutorial: ~500 words. Trigger electron-go-pro agent with a project description, walk through what it produces
- How-to: Set up an Electron+Go project with IPC layer
- Reference: electron-go-pro agent — model (sonnet), tools, triggers, capabilities (IPC design, code signing, packaging)
- Explanation: Why hybrid Electron+Go rather than pure Electron or pure native. When to use this plugin vs cli-development for terminal apps.
- Cross-refs: cli-development (alternative for terminal apps), programming-languages/go-simplifier (Go code quality)

**Step 1:** Create directories and all files with the Write tool.

**Step 2:** Commit.

```bash
git add desktop-development/docs/
git commit -m "add desktop-development plugin documentation"
```

---

### Task 7: mcp-development docs (2 agents: mcp-integration-engineer, mcp-registry-navigator)

**Files:**

- Create: `mcp-development/docs/_index.md`
- Create: `mcp-development/docs/tutorials/_index.md`
- Create: `mcp-development/docs/tutorials/getting-started.md`
- Create: `mcp-development/docs/howto/_index.md`
- Create: `mcp-development/docs/howto/find-and-evaluate-mcp-servers.md`
- Create: `mcp-development/docs/reference/_index.md`
- Create: `mcp-development/docs/reference/agents.md`
- Create: `mcp-development/docs/explanation/_index.md`
- Create: `mcp-development/docs/explanation/architecture.md`

**Content focus:**

- Overview: MCP server discovery, integration, and multi-server orchestration
- Tutorial: ~500 words. Use mcp-registry-navigator to find an MCP server, then mcp-integration-engineer to configure it
- How-to: Find and evaluate MCP servers for trustworthiness, maintenance, security
- Reference: Both agents — models, tools, triggers, capabilities
- Explanation: What MCP is, why dedicated agents for it, when to use registry-navigator vs integration-engineer

**Step 1:** Create directories and all files with the Write tool.

**Step 2:** Commit.

```bash
git add mcp-development/docs/
git commit -m "add mcp-development plugin documentation"
```

---

### Task 8: ai-development docs (1 agent: ai-engineer, 1 skill: agent-modernizer)

**Files:**

- Create: `ai-development/docs/_index.md`
- Create: `ai-development/docs/tutorials/_index.md`
- Create: `ai-development/docs/tutorials/getting-started.md`
- Create: `ai-development/docs/howto/_index.md`
- Create: `ai-development/docs/howto/audit-agent-definitions.md`
- Create: `ai-development/docs/reference/_index.md`
- Create: `ai-development/docs/reference/agents.md`
- Create: `ai-development/docs/reference/skills.md`
- Create: `ai-development/docs/explanation/_index.md`
- Create: `ai-development/docs/explanation/architecture.md`

**Content focus:**

- Overview: AI/ML feature implementation and agent definition modernization
- Tutorial: ~500 words. Use agent-modernizer skill to audit an agent definition, walk through recommendations
- How-to: Audit and modernize agent definitions for best practices
- Reference: ai-engineer agent (capabilities, model, tools) + agent-modernizer skill (trigger phrases, version)
- Explanation: Two-pronged purpose — building AI features (ai-engineer) and maintaining agent quality (agent-modernizer). Cross-ref: research plugin for gathering data before AI implementation

**Step 1:** Create directories and all files with the Write tool.

**Step 2:** Commit.

```bash
git add ai-development/docs/
git commit -m "add ai-development plugin documentation"
```

---

### Task 9: utility-skills docs (1 skill: markdown-nested-codeblocks)

**Files:**

- Create: `utility-skills/docs/_index.md`
- Create: `utility-skills/docs/tutorials/_index.md`
- Create: `utility-skills/docs/tutorials/getting-started.md`
- Create: `utility-skills/docs/howto/_index.md`
- Create: `utility-skills/docs/howto/nest-code-blocks-correctly.md`
- Create: `utility-skills/docs/reference/_index.md`
- Create: `utility-skills/docs/reference/skills.md`
- Create: `utility-skills/docs/explanation/_index.md`
- Create: `utility-skills/docs/explanation/architecture.md`

**Content focus:**

- Overview: Cross-domain utility skills for markdown and documentation
- Tutorial: ~500 words. Write markdown containing nested code blocks, trigger the skill, see it fix the fencing
- How-to: Correctly nest code blocks using the k+1 backtick rule
- Reference: markdown-nested-codeblocks skill — trigger phrases, version, the k+1 rule
- Explanation: Why nested codeblocks break, the backtick counting rule, why this is a skill rather than an agent

**Step 1:** Create directories and all files with the Write tool.

**Step 2:** Commit.

```bash
git add utility-skills/docs/
git commit -m "add utility-skills plugin documentation"
```

---

## Group C: Tier 2 (Standard) Plugin Documentation

Each Tier 2 plugin gets ~11 files: _index.md, 4 section indexes, 1 tutorial, 2-3 how-tos, 1 reference (per component type), 1-2 explanation pages.

### Task 10: code-quality docs (5 agents)

**Files:**

- Create: `code-quality/docs/_index.md`
- Create: `code-quality/docs/tutorials/_index.md`
- Create: `code-quality/docs/tutorials/getting-started.md`
- Create: `code-quality/docs/howto/_index.md`
- Create: `code-quality/docs/howto/review-code-before-committing.md`
- Create: `code-quality/docs/howto/write-and-fix-tests.md`
- Create: `code-quality/docs/howto/audit-accessibility.md`
- Create: `code-quality/docs/reference/_index.md`
- Create: `code-quality/docs/reference/agents.md`
- Create: `code-quality/docs/explanation/_index.md`
- Create: `code-quality/docs/explanation/architecture.md`
- Create: `code-quality/docs/explanation/choosing-the-right-agent.md`

**Content focus:**

- Overview: Code review, testing, accessibility, and architectural quality
- Agents: code-reviewer, architect-reviewer, test-writer-fixer, playwright-expert, web-accessibility-checker
- Tutorial: ~800 words. Make a code change, trigger code-reviewer, then test-writer-fixer. Show the review → test → fix cycle.
- How-tos: (1) Review code before committing, (2) Write and fix tests with test-writer-fixer, (3) Audit accessibility with web-accessibility-checker
- Reference: All 5 agents with model, tools, triggers, capabilities
- Explanation: Why 5 specialized reviewers instead of 1. The review spectrum (code → architecture → tests → a11y). Cross-ref: programming-languages for language-specific expertise
- Decision guide: When to use code-reviewer vs architect-reviewer vs just asking Claude

**Step 1:** Create directories and all files with the Write tool.

**Step 2:** Commit.

```bash
git add code-quality/docs/
git commit -m "add code-quality plugin documentation"
```

---

### Task 11: web-development docs (6 agents)

**Files:**

- Create: `web-development/docs/_index.md`
- Create: `web-development/docs/tutorials/_index.md`
- Create: `web-development/docs/tutorials/getting-started.md`
- Create: `web-development/docs/howto/_index.md`
- Create: `web-development/docs/howto/build-react-component.md`
- Create: `web-development/docs/howto/set-up-nextjs-project.md`
- Create: `web-development/docs/reference/_index.md`
- Create: `web-development/docs/reference/agents.md`
- Create: `web-development/docs/explanation/_index.md`
- Create: `web-development/docs/explanation/architecture.md`
- Create: `web-development/docs/explanation/choosing-the-right-agent.md`

**Content focus:**

- Agents: frontend-developer, fullstack-developer, nextjs-developer, react-specialist, ui-designer, vue-expert
- Tutorial: ~800 words. Build a React component with react-specialist, style it with ui-designer guidance
- How-tos: (1) Build a React component, (2) Set up a Next.js project
- Reference: All 6 agents
- Explanation: Framework specialization model, when to use react-specialist vs frontend-developer vs fullstack-developer

**Step 1:** Create directories and all files with the Write tool.

**Step 2:** Commit.

```bash
git add web-development/docs/
git commit -m "add web-development plugin documentation"
```

---

### Task 12: backend-development docs (3 agents)

**Files:**

- Create: `backend-development/docs/_index.md`
- Create: `backend-development/docs/tutorials/_index.md`
- Create: `backend-development/docs/tutorials/getting-started.md`
- Create: `backend-development/docs/howto/_index.md`
- Create: `backend-development/docs/howto/design-api-architecture.md`
- Create: `backend-development/docs/howto/optimize-sql-queries.md`
- Create: `backend-development/docs/reference/_index.md`
- Create: `backend-development/docs/reference/agents.md`
- Create: `backend-development/docs/explanation/_index.md`
- Create: `backend-development/docs/explanation/architecture.md`

**Content focus:**

- Agents: backend-architect, go-architect, sql-pro
- Tutorial: ~800 words. Design an API with backend-architect, implement in Go with go-architect
- How-tos: (1) Design API architecture, (2) Optimize SQL queries with sql-pro
- Reference: All 3 agents
- Explanation: Architecture-first vs code-first. When to use backend-architect (design) vs go-architect (implementation). Cross-ref: devops-and-infra for deployment

**Step 1:** Create directories and all files with the Write tool.

**Step 2:** Commit.

```bash
git add backend-development/docs/
git commit -m "add backend-development plugin documentation"
```

---

### Task 13: cli-development docs (3 agents)

**Files:**

- Create: `cli-development/docs/_index.md`
- Create: `cli-development/docs/tutorials/_index.md`
- Create: `cli-development/docs/tutorials/getting-started.md`
- Create: `cli-development/docs/howto/_index.md`
- Create: `cli-development/docs/howto/build-interactive-tui.md`
- Create: `cli-development/docs/howto/design-cli-visual-style.md`
- Create: `cli-development/docs/reference/_index.md`
- Create: `cli-development/docs/reference/agents.md`
- Create: `cli-development/docs/explanation/_index.md`
- Create: `cli-development/docs/explanation/architecture.md`

**Content focus:**

- Agents: cli-developer, go-tui-developer, cli-ui-designer
- Tutorial: ~800 words. Build a CLI command with cli-developer, add TUI with go-tui-developer
- How-tos: (1) Build an interactive TUI with Bubble Tea, (2) Design terminal visual style
- Reference: All 3 agents
- Explanation: CLI vs TUI distinction, when to use each agent. Cross-ref: programming-languages for Go expertise

**Step 1:** Create directories and all files with the Write tool.

**Step 2:** Commit.

```bash
git add cli-development/docs/
git commit -m "add cli-development plugin documentation"
```

---

### Task 14: developer-tools docs (3 agents)

**Files:**

- Create: `developer-tools/docs/_index.md`
- Create: `developer-tools/docs/tutorials/_index.md`
- Create: `developer-tools/docs/tutorials/getting-started.md`
- Create: `developer-tools/docs/howto/_index.md`
- Create: `developer-tools/docs/howto/set-up-git-workflow.md`
- Create: `developer-tools/docs/howto/generate-api-documentation.md`
- Create: `developer-tools/docs/reference/_index.md`
- Create: `developer-tools/docs/reference/agents.md`
- Create: `developer-tools/docs/explanation/_index.md`
- Create: `developer-tools/docs/explanation/architecture.md`

**Content focus:**

- Agents: documentation-engineer, git-workflow-manager, rapid-prototyper
- Tutorial: ~800 words. Set up a git workflow with git-workflow-manager, then generate docs with documentation-engineer
- How-tos: (1) Set up branching strategy and automation, (2) Generate API documentation from source
- Reference: All 3 agents
- Explanation: The developer workflow spectrum — version control, documentation, prototyping. Cross-ref: diataxis-docs for structured documentation

**Step 1:** Create directories and all files with the Write tool.

**Step 2:** Commit.

```bash
git add developer-tools/docs/
git commit -m "add developer-tools plugin documentation"
```

---

### Task 15: utility-agents docs (3 agents)

**Files:**

- Create: `utility-agents/docs/_index.md`
- Create: `utility-agents/docs/tutorials/_index.md`
- Create: `utility-agents/docs/tutorials/getting-started.md`
- Create: `utility-agents/docs/howto/_index.md`
- Create: `utility-agents/docs/howto/validate-urls-in-documentation.md`
- Create: `utility-agents/docs/howto/extract-all-urls-from-codebase.md`
- Create: `utility-agents/docs/reference/_index.md`
- Create: `utility-agents/docs/reference/agents.md`
- Create: `utility-agents/docs/explanation/_index.md`
- Create: `utility-agents/docs/explanation/architecture.md`

**Content focus:**

- Agents: episode-orchestrator, url-context-validator, url-link-extractor
- Tutorial: ~800 words. Use url-link-extractor to inventory URLs in a project, then url-context-validator to check them
- How-tos: (1) Validate URLs in documentation for freshness and relevance, (2) Extract all URLs from a codebase for migration
- Reference: All 3 agents
- Explanation: Utility vs domain-specific agents. The URL analysis pipeline (extract → validate → report).

**Step 1:** Create directories and all files with the Write tool.

**Step 2:** Commit.

```bash
git add utility-agents/docs/
git commit -m "add utility-agents plugin documentation"
```

---

### Task 16: mobile-development docs (3 agents)

**Files:**

- Create: `mobile-development/docs/_index.md`
- Create: `mobile-development/docs/tutorials/_index.md`
- Create: `mobile-development/docs/tutorials/getting-started.md`
- Create: `mobile-development/docs/howto/_index.md`
- Create: `mobile-development/docs/howto/build-ios-feature.md`
- Create: `mobile-development/docs/howto/choose-mobile-platform.md`
- Create: `mobile-development/docs/reference/_index.md`
- Create: `mobile-development/docs/reference/agents.md`
- Create: `mobile-development/docs/explanation/_index.md`
- Create: `mobile-development/docs/explanation/architecture.md`

**Content focus:**

- Agents: ios-developer, mobile-developer, swift-expert
- Tutorial: ~800 words. Build a simple iOS feature with ios-developer, refine Swift code with swift-expert
- How-tos: (1) Build an iOS feature, (2) Choose between native, cross-platform, and web approaches
- Reference: All 3 agents
- Explanation: Platform specialization — when to use ios-developer vs mobile-developer vs swift-expert

**Step 1:** Create directories and all files with the Write tool.

**Step 2:** Commit.

```bash
git add mobile-development/docs/
git commit -m "add mobile-development plugin documentation"
```

---

### Task 17: devops-and-infra docs (4 agents)

**Files:**

- Create: `devops-and-infra/docs/_index.md`
- Create: `devops-and-infra/docs/tutorials/_index.md`
- Create: `devops-and-infra/docs/tutorials/getting-started.md`
- Create: `devops-and-infra/docs/howto/_index.md`
- Create: `devops-and-infra/docs/howto/set-up-github-actions-ci.md`
- Create: `devops-and-infra/docs/howto/configure-prometheus-monitoring.md`
- Create: `devops-and-infra/docs/reference/_index.md`
- Create: `devops-and-infra/docs/reference/agents.md`
- Create: `devops-and-infra/docs/explanation/_index.md`
- Create: `devops-and-infra/docs/explanation/architecture.md`

**Content focus:**

- Agents: devops-automator, github-actions-expert, performance-monitor, prometheus-expert
- Tutorial: ~800 words. Set up a GitHub Actions CI pipeline with github-actions-expert, add monitoring with prometheus-expert
- How-tos: (1) Set up GitHub Actions CI/CD, (2) Configure Prometheus monitoring
- Reference: All 4 agents
- Explanation: The DevOps pipeline — build → deploy → monitor. Cross-ref: backend-development for what gets deployed

**Step 1:** Create directories and all files with the Write tool.

**Step 2:** Commit.

```bash
git add devops-and-infra/docs/
git commit -m "add devops-and-infra plugin documentation"
```

---

### Task 18: research docs (5 agents)

**Files:**

- Create: `research/docs/_index.md`
- Create: `research/docs/tutorials/_index.md`
- Create: `research/docs/tutorials/getting-started.md`
- Create: `research/docs/howto/_index.md`
- Create: `research/docs/howto/conduct-literature-review.md`
- Create: `research/docs/howto/produce-research-report.md`
- Create: `research/docs/reference/_index.md`
- Create: `research/docs/reference/agents.md`
- Create: `research/docs/explanation/_index.md`
- Create: `research/docs/explanation/architecture.md`
- Create: `research/docs/explanation/research-pipeline.md`

**Content focus:**

- Agents: academic-researcher, comprehensive-researcher, technical-researcher, research-synthesizer, report-generator
- Tutorial: ~800 words. Run a comprehensive research investigation: dispatch comprehensive-researcher, synthesize with research-synthesizer, generate report with report-generator
- How-tos: (1) Conduct an academic literature review, (2) Produce a polished research report from findings
- Reference: All 5 agents
- Explanation: (1) Architecture — the research pipeline (gather → synthesize → report). (2) Research pipeline — why 5 agents rather than 1, the specialization rationale. Cross-ref: ai-development for applying research findings

**Step 1:** Create directories and all files with the Write tool.

**Step 2:** Commit.

```bash
git add research/docs/
git commit -m "add research plugin documentation"
```

---

### Task 19: gnu-make docs (5 skills)

**Files:**

- Create: `gnu-make/docs/_index.md`
- Create: `gnu-make/docs/tutorials/_index.md`
- Create: `gnu-make/docs/tutorials/getting-started.md`
- Create: `gnu-make/docs/howto/_index.md`
- Create: `gnu-make/docs/howto/debug-slow-makefile.md`
- Create: `gnu-make/docs/howto/organize-multi-directory-build.md`
- Create: `gnu-make/docs/howto/split-large-makefile.md`
- Create: `gnu-make/docs/reference/_index.md`
- Create: `gnu-make/docs/reference/skills.md`
- Create: `gnu-make/docs/explanation/_index.md`
- Create: `gnu-make/docs/explanation/architecture.md`
- Create: `gnu-make/docs/explanation/skill-progression.md`

**Content focus:**

- Skills: makefile-fundamentals, makefile-advanced-features, makefile-recursive-multi-directory, makefile-includes-modularity, makefile-debugging-optimization
- Tutorial: ~800 words. Create a Makefile from scratch using fundamentals skill, add a pattern rule with advanced-features
- How-tos: (1) Debug a slow Makefile, (2) Organize a multi-directory build, (3) Split a monolithic Makefile into modules
- Reference: All 5 skills with trigger phrases, versions, what they cover
- Explanation: (1) Architecture — skills-only plugin, no agents. Why skills for build systems. (2) Skill progression — fundamentals → advanced → multi-directory → modular → debugging. The pedagogical structure. Cross-ref: developer-tools for broader build workflow

**Step 1:** Create directories and all files with the Write tool.

**Step 2:** Commit.

```bash
git add gnu-make/docs/
git commit -m "add gnu-make plugin documentation"
```

---

### Task 20: game-development docs (5 agents)

**Files:**

- Create: `game-development/docs/_index.md`
- Create: `game-development/docs/tutorials/_index.md`
- Create: `game-development/docs/tutorials/getting-started.md`
- Create: `game-development/docs/howto/_index.md`
- Create: `game-development/docs/howto/prototype-game-mechanic.md`
- Create: `game-development/docs/howto/set-up-unity-project.md`
- Create: `game-development/docs/reference/_index.md`
- Create: `game-development/docs/reference/agents.md`
- Create: `game-development/docs/explanation/_index.md`
- Create: `game-development/docs/explanation/architecture.md`
- Create: `game-development/docs/explanation/choosing-the-right-agent.md`

**Content focus:**

- Agents: 3d-artist, game-designer, game-developer, unity-game-developer, unreal-engine-developer
- Tutorial: ~800 words. Design a game mechanic with game-designer, implement in Unity with unity-game-developer
- How-tos: (1) Prototype a game mechanic, (2) Set up a Unity project
- Reference: All 5 agents
- Explanation: (1) Architecture — design vs implementation vs engine-specific agents. (2) Choosing the right agent — game-designer (mechanics) vs game-developer (general) vs engine-specific

**Step 1:** Create directories and all files with the Write tool.

**Step 2:** Commit.

```bash
git add game-development/docs/
git commit -m "add game-development plugin documentation"
```

---

## Group D: Tier 3 (Comprehensive) Plugin Documentation

Each Tier 3 plugin gets ~16 files: _index.md, 4 section indexes, 1 tutorial, 3-5 how-tos, reference pages per component type, 2-3 explanation pages.

### Task 21: hugo-repo docs (7 skills, 2 agents, 4 commands, 4 templates)

**Files:**

- Create: `hugo-repo/docs/_index.md`
- Create: `hugo-repo/docs/tutorials/_index.md`
- Create: `hugo-repo/docs/tutorials/getting-started.md`
- Create: `hugo-repo/docs/howto/_index.md`
- Create: `hugo-repo/docs/howto/set-up-hugo-site.md`
- Create: `hugo-repo/docs/howto/add-content-section.md`
- Create: `hugo-repo/docs/howto/deploy-to-github-pages.md`
- Create: `hugo-repo/docs/howto/deploy-to-s3.md`
- Create: `hugo-repo/docs/howto/customize-theme.md`
- Create: `hugo-repo/docs/reference/_index.md`
- Create: `hugo-repo/docs/reference/skills.md`
- Create: `hugo-repo/docs/reference/agents.md`
- Create: `hugo-repo/docs/reference/commands.md`
- Create: `hugo-repo/docs/reference/templates.md`
- Create: `hugo-repo/docs/explanation/_index.md`
- Create: `hugo-repo/docs/explanation/architecture.md`
- Create: `hugo-repo/docs/explanation/design-decisions.md`
- Create: `hugo-repo/docs/explanation/component-interaction.md`

**Content focus:**

- Tutorial: ~1200 words. End-to-end: /hugo-init a site, add a content section with /hugo-add-section, customize theme, deploy with /hugo-deploy. Exercise skills, agents, and commands together.
- How-tos: (1) Set up Hugo site with module mounts, (2) Add a content section from a new directory, (3) Deploy to GitHub Pages, (4) Deploy to AWS S3, (5) Customize a Hugo theme
- Reference: Separate page per component type. Skills: all 7 with trigger phrases, coverage areas. Agents: both with model, tools, triggers. Commands: all 4 with arguments. Templates: all 4 with placeholders.
- Explanation: (1) Architecture — how the 17 components fit together. (2) Design decisions — why 7 skills, why separate deployment skills, module mounts over symlinks. (3) Component interaction — how skills inform agents, how commands use templates. Cross-ref: diataxis-docs for documentation structure

**Step 1:** Create directories and all files with the Write tool.

**Step 2:** Commit.

```bash
git add hugo-repo/docs/
git commit -m "add hugo-repo plugin documentation"
```

---

### Task 22: programming-languages docs (8 agents)

**Files:**

- Create: `programming-languages/docs/_index.md`
- Create: `programming-languages/docs/tutorials/_index.md`
- Create: `programming-languages/docs/tutorials/getting-started.md`
- Create: `programming-languages/docs/howto/_index.md`
- Create: `programming-languages/docs/howto/modernize-legacy-code.md`
- Create: `programming-languages/docs/howto/optimize-performance.md`
- Create: `programming-languages/docs/howto/migrate-to-typescript.md`
- Create: `programming-languages/docs/reference/_index.md`
- Create: `programming-languages/docs/reference/agents.md`
- Create: `programming-languages/docs/explanation/_index.md`
- Create: `programming-languages/docs/explanation/architecture.md`
- Create: `programming-languages/docs/explanation/language-specialization.md`

**Content focus:**

- Agents: cpp-pro, csharp-expert, go-simplifier, javascript-pro, rust-pro, shell-scripting-pro, typescript-pro, zsh-expert
- Tutorial: ~1200 words. Write code in one language, use the appropriate language agent to review and improve it, then use go-simplifier to clean up Go code
- How-tos: (1) Modernize legacy code to modern language idioms, (2) Optimize performance in any supported language, (3) Migrate JavaScript to TypeScript
- Reference: All 8 agents with languages, tools, model, capabilities
- Explanation: (1) Architecture — one agent per language vs a general "coding" agent. (2) Language specialization — what each agent knows that a general agent doesn't. Cross-ref: code-quality for language-agnostic review

**Step 1:** Create directories and all files with the Write tool.

**Step 2:** Commit.

```bash
git add programming-languages/docs/
git commit -m "add programming-languages plugin documentation"
```

---

### Task 23: diataxis-docs docs (7 agents)

**Files:**

- Create: `diataxis-docs/docs/_index.md`
- Create: `diataxis-docs/docs/tutorials/_index.md`
- Create: `diataxis-docs/docs/tutorials/getting-started.md`
- Create: `diataxis-docs/docs/howto/_index.md`
- Create: `diataxis-docs/docs/howto/restructure-docs-to-diataxis.md`
- Create: `diataxis-docs/docs/howto/write-a-tutorial.md`
- Create: `diataxis-docs/docs/howto/write-reference-docs.md`
- Create: `diataxis-docs/docs/howto/validate-doc-quality.md`
- Create: `diataxis-docs/docs/reference/_index.md`
- Create: `diataxis-docs/docs/reference/agents.md`
- Create: `diataxis-docs/docs/explanation/_index.md`
- Create: `diataxis-docs/docs/explanation/architecture.md`
- Create: `diataxis-docs/docs/explanation/diataxis-in-practice.md`
- Create: `diataxis-docs/docs/explanation/orchestration-model.md`

**Content focus:**

- Agents: diataxis-orchestrator, doc-crosslink-validator, doc-explanation-writer, doc-howto-writer, doc-inventory, doc-reference-gen, doc-tutorial-writer
- Tutorial: ~1200 words. Restructure a documentation directory using the full Diataxis pipeline: inventory → orchestrate → write tutorials/howtos/reference/explanation → validate
- How-tos: (1) Restructure existing docs to Diataxis, (2) Write a tutorial with doc-tutorial-writer, (3) Write reference docs with doc-reference-gen, (4) Validate doc quality with doc-crosslink-validator
- Reference: All 7 agents with roles, models, tools, triggers
- Explanation: (1) Architecture — the orchestrator pattern, why 7 specialized agents. (2) Diataxis in practice — how the four types work together, common pitfalls. (3) Orchestration model — how diataxis-orchestrator coordinates the specialist agents. Cross-ref: developer-tools/documentation-engineer for non-Diataxis docs, hugo-repo for site generation

**Step 1:** Create directories and all files with the Write tool.

**Step 2:** Commit.

```bash
git add diataxis-docs/docs/
git commit -m "add diataxis-docs plugin documentation"
```

---

## Summary of Commits

1. `add marketplace root docs scaffolding and overview` (Task 1)
2. `add marketplace getting-started tutorial` (Task 2)
3. `add marketplace how-to guides` (Task 3)
4. `add marketplace reference documentation` (Task 4)
5. `add marketplace explanation documentation` (Task 5)
6. `add desktop-development plugin documentation` (Task 6)
7. `add mcp-development plugin documentation` (Task 7)
8. `add ai-development plugin documentation` (Task 8)
9. `add utility-skills plugin documentation` (Task 9)
10. `add code-quality plugin documentation` (Task 10)
11. `add web-development plugin documentation` (Task 11)
12. `add backend-development plugin documentation` (Task 12)
13. `add cli-development plugin documentation` (Task 13)
14. `add developer-tools plugin documentation` (Task 14)
15. `add utility-agents plugin documentation` (Task 15)
16. `add mobile-development plugin documentation` (Task 16)
17. `add devops-and-infra plugin documentation` (Task 17)
18. `add research plugin documentation` (Task 18)
19. `add gnu-make plugin documentation` (Task 19)
20. `add game-development plugin documentation` (Task 20)
21. `add hugo-repo plugin documentation` (Task 21)
22. `add programming-languages plugin documentation` (Task 22)
23. `add diataxis-docs plugin documentation` (Task 23)
