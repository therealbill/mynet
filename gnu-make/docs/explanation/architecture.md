---
title: "Architecture of a Skills-Only Plugin"
description: "Why the gnu-make plugin uses skills instead of agents, and how knowledge injection fits build system tooling"
weight: 1
---

# Architecture of a Skills-Only Plugin

The gnu-make plugin is a skills-only plugin. It contains five skills and zero agents. This is a deliberate architectural choice that reflects how build system expertise is best delivered.

## Skills vs Agents for Build Systems

Agents perform work. They execute multi-step workflows, make decisions, and produce artifacts. A deployment agent might provision infrastructure, run tests, and push to production. An analysis agent might scan a codebase, identify issues, and generate reports.

Skills inject knowledge. They provide patterns, best practices, decision frameworks, and anti-pattern awareness to an AI assistant that is already performing work on behalf of the user. A skill does not act independently. It augments the assistant's responses with domain expertise.

For build systems, skills are the right fit because the user (or the assistant on the user's behalf) is already writing or modifying a Makefile. The need is not for an autonomous workflow but for expert guidance applied to a specific context. When someone asks "create a Makefile for my C project," the assistant needs to know about the `##` help pattern, `.PHONY` declarations, pattern rules, and proper variable usage. That knowledge comes from skills.

An agent would be appropriate if the goal were something like "analyze my entire build system and produce a migration plan." The gnu-make plugin does not address that use case. It addresses the much more common case of writing correct, professional Makefiles in the flow of normal development work.

## Knowledge Injection as a Pattern

Each skill in the gnu-make plugin works by injecting knowledge at the point where the assistant is formulating a response. When the user asks to create a Makefile, the makefile-fundamentals skill activates and provides:

- Non-negotiable practices that should always be present (TAB characters, `.PHONY`, `.DELETE_ON_ERROR`)
- The `##` self-documenting help pattern with its implementation
- Standard variable names and their conventional meanings
- Common mistakes and their fixes

The assistant then applies this knowledge to the user's specific situation. The skill does not generate a generic template; it ensures the assistant has the expertise to create a Makefile tailored to the user's project.

This approach has several advantages over template-based or agent-based alternatives:

- **Context-sensitive**: The assistant adapts the knowledge to the specific project, language, and build requirements.
- **Corrective**: Skills include anti-patterns, so the assistant can identify and fix problems in existing Makefiles rather than only creating new ones.
- **Educational**: Skills carry explanations of why each practice matters, enabling the assistant to teach as it builds.
- **Composable**: Multiple skills can activate together when a request spans several areas.

## The Pedagogical Progression

The five skills are ordered in a deliberate progression that mirrors how Make expertise develops in practice.

**makefile-fundamentals** establishes the foundation. Every Makefile needs TAB characters, `.PHONY` declarations, and a help target. These are non-negotiable regardless of project size. The skill also introduces pattern rules and automatic variables at a basic level, providing enough to write correct Makefiles for simple projects.

**makefile-advanced-features** builds on the fundamentals by going deeper into pattern rules, automatic variables, and DRY practices. Where fundamentals says "use `%.o: %.c`," advanced features covers `$(wildcard)` for auto-discovery, `$(patsubst)` for transformation, conditional compilation with `ifeq`, and target-specific variables. The core lesson is that pattern rules are simpler than explicit rules, not more complex.

**makefile-recursive-multi-directory** extends the single-Makefile knowledge to projects spanning multiple directories. It teaches the phony target pattern for recursive builds, dependency declarations between subdirectories, variable export, and the critical `$(MAKE)` variable. The key insight is that phony targets enable 8x speedup through parallelization.

**makefile-includes-modularity** addresses organization as Makefiles grow. At 150+ lines, a monolithic Makefile becomes hard to navigate and causes merge conflicts. This skill teaches the `include` directive, standard module names (`config.mk`, `rules.mk`, `targets.mk`), environment-specific configurations, and cross-project shared configuration.

**makefile-debugging-optimization** caps the progression with diagnostic expertise. Rather than jumping to solutions, this skill teaches systematic debugging using Make's built-in flags: `make -n` for dry-runs, `make -d` for decision tracing, `make -p` for database inspection, and `make --trace` for execution tracing. Optimization comes after diagnosis, not before.

## How Skills Complement Each Other

The skills are designed to work together. A user creating a large multi-directory project might trigger three or four skills simultaneously:

- **fundamentals** for the `##` help pattern and `.PHONY` in each subdirectory Makefile
- **advanced-features** for pattern rules and auto-discovery within each subdirectory
- **recursive-multi-directory** for the root Makefile coordinating subdirectories
- **includes-modularity** if any individual Makefile exceeds 150 lines

Each skill focuses on its domain without duplicating material from others. Fundamentals introduces pattern rules; advanced features covers them in depth. Fundamentals establishes `.PHONY`; recursive-multi-directory applies it to subdirectory targets.

## Relationship to the Broader Ecosystem

The gnu-make plugin focuses specifically on GNU Make best practices. For broader build workflow concerns -- CI/CD integration, multi-tool build pipelines, containerized builds -- the developer-tools plugin provides complementary coverage. The gnu-make plugin does not attempt to cover these areas. Its scope is the Makefile itself: writing it correctly, organizing it well, and debugging it systematically.

This narrow scope is intentional. A skill that tries to cover everything dilutes its expertise. By staying focused on Make, each skill can provide deep, opinionated guidance rather than surface-level advice.

## See Also

- [Skill Reference]({{< ref "reference/skills" >}}) -- technical specifications for all five skills
- [Skill Progression]({{< ref "explanation/skill-progression" >}}) -- the learning path through fundamentals to debugging
- [Getting Started]({{< ref "tutorials/getting-started" >}}) -- hands-on introduction to the skills
