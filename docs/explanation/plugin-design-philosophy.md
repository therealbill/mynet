---
title: "Plugin Design Philosophy"
description: "Design decisions behind plugin component types and specialization model"
weight: 2
---

# Plugin Design Philosophy

The Mynet marketplace contains 64 agents, 30 skills, 10 commands, and 9 templates spread across 19 plugins. These numbers are not arbitrary. They reflect a core design principle: narrow specialization produces better results than broad generalization, and different types of assistance require different component types.

## The specialization thesis

A single "backend agent" could theoretically handle SQL queries, Go architecture, and API design. In practice, it would do all three poorly. Its system prompt would need to cover so much ground that it could not go deep on any topic. It would lack room for the EXPLAIN ANALYZE interpretation hints that make `sql-pro` effective, or the Goa framework conventions that `go-architect` carries.

The marketplace takes the opposite approach: the `backend-development` plugin contains three agents (`backend-architect`, `go-architect`, `sql-pro`), each with a system prompt entirely focused on one domain. When a user asks about query optimization, `sql-pro` activates with its full context of window functions, CTEs, indexing strategies, and execution plan analysis. No prompt space is wasted on unrelated Go or API concerns.

This is not just a theoretical preference. The description field that drives agent dispatch works better with narrow agents. A broad description like "handles backend development" competes with many other agents for any backend-related query. A narrow description like "write complex SQL queries, optimize execution plans, and design normalized schemas" wins the dispatch match cleanly when the user's intent involves SQL.

## Four component types, four purposes

The marketplace defines four component types, each serving a distinct role in how Claude Code provides assistance.

**Agents** are autonomous specialists. They receive control of the conversation, have access to tools (file reading, code execution, web search), and can perform multi-step tasks. An agent definition is a Markdown file with YAML frontmatter that includes a name, description, and system prompt. The system prompt shapes the agent's personality, expertise, and approach. Agents are the right choice when a task requires autonomy -- writing code, analyzing architectures, debugging problems.

**Skills** are knowledge injections. Unlike agents, skills do not take over the conversation. They activate based on topic detection and provide domain knowledge to whichever agent is currently active. A skill like `makefile-fundamentals` does not write Makefiles itself -- it provides the active agent with knowledge about `.PHONY` declarations, the `##` help pattern, and tab-vs-space rules. Skills are the right choice when any agent might benefit from domain-specific knowledge. The `gnu-make` plugin contains zero agents and five skills because Makefile expertise is useful regardless of which agent is doing the work.

**Commands** are user-initiated actions. They register as `/command-name` in Claude Code and activate only when explicitly invoked. Commands have defined inputs and outputs -- `/tl-deploy` deploys a Temporal cluster, `/hugo-init` scaffolds a Hugo site. Commands are the right choice for repeatable workflows where the user knows exactly what they want and the action has a predictable structure.

**Templates** are file scaffolding. They provide starter files with placeholder values -- a `workflow.go.tmpl` template in the `timelord` plugin generates a properly structured Temporal workflow file. Templates are the right choice when the output is a well-defined file format with known conventions that should not be reinvented each time.

The relationship between these types is complementary, not hierarchical:

| Component | Activation      | Autonomy | Purpose                        |
|-----------|-----------------|----------|--------------------------------|
| Agent     | Claude dispatch  | Full     | Multi-step tasks with tools    |
| Skill     | Topic detection  | None     | Knowledge injection            |
| Command   | User invocation  | Scripted | Repeatable workflows           |
| Template  | Agent/command use| None     | File generation                |

## Plugin domain boundaries

Each plugin covers one coherent domain, but "coherent" does not mean "small." The `programming-languages` plugin contains 8 agents (one per language), while `desktop-development` contains just 1. Both are correctly scoped.

The guiding principle for domain boundaries is discoverability: what expertise groups together from a user's perspective? All language experts go in `programming-languages` even though a C++ expert and a Rust expert share no code, because a user searching for "language specialist" expects to find them in one place. Conversely, Go appears in three different plugins -- `programming-languages` has `go-simplifier` for general Go idioms, `backend-development` has `go-architect` for Go service design, and `cli-development` has `go-tui-developer` for Go terminal UIs. The same language spans multiple plugins because the domains (language idioms vs. service architecture vs. terminal UIs) serve different user needs.

Some plugins are agent-heavy while others are skill-heavy or command-heavy. This variation is not inconsistent design -- it reflects the nature of the domain:

- **`diataxis-docs`** has 7 agents because documentation writing requires autonomous, multi-step work across different document types
- **`gnu-make`** has 5 skills and zero agents because Makefile expertise is best delivered as knowledge injection to whatever agent is already working
- **`timelord`** has 3 agents, 16 skills, 6 commands, and 5 templates because Temporal.io operations span autonomous work, deep knowledge, repeatable workflows, and boilerplate generation

## Cross-plugin collaboration

Plugins never import from each other. There are no shared utilities, no common base classes, no dependency declarations between plugins. This is a deliberate trade-off.

The benefit is complete independence. Any plugin can be added or removed without affecting any other. The `hugo-repo` plugin could be deleted tomorrow and `timelord` would not notice. This independence also means plugins can be installed individually -- a user who only needs documentation tools can install `diataxis-docs` alone.

The cost is potential duplication. If two plugins need similar knowledge, they each carry their own copy. In practice, this has been minimal because plugins cover genuinely different domains.

Collaboration between plugins happens entirely through Claude's dispatch layer. When a user writes TypeScript code and then asks for a code review, `typescript-pro` from `programming-languages` handles the writing and `code-reviewer` from `code-quality` handles the review. Neither plugin knows the other exists. Claude orchestrates the handoff based on intent matching against descriptions. This loose coupling is not a limitation -- it is the mechanism that makes the flat, independent plugin model work.

## Related

- [Component Conventions]({{< ref "reference/component-conventions" >}}) -- field-by-field specification for agents, skills, commands, and templates
