---
title: "How to Use Plugins Together"
description: "Combine agents from multiple Mynet plugins in a single Claude Code session to build multi-step workflows: research into reports, code generation with review, and development with documentation."
weight: 3
---

# How to Use Plugins Together

Combine agents from multiple Mynet plugins in a single session to handle workflows that span research, coding, review, and documentation.

## Prerequisites

- Two or more plugins installed in your Claude Code project -- see [How to Install a Plugin](../install-a-plugin/)
- Understanding of what each plugin provides -- see [How to Find the Right Plugin](../find-right-plugin/)

## Steps

### 1. Install Multiple Plugins

Add all required plugins to your project's `.claude/settings.json`:

```json
{
  "plugins": [
    {
      "name": "research",
      "source": "/path/to/claude-plugins/research"
    },
    {
      "name": "code-quality",
      "source": "/path/to/claude-plugins/code-quality"
    },
    {
      "name": "programming-languages",
      "source": "/path/to/claude-plugins/programming-languages"
    },
    {
      "name": "diataxis-docs",
      "source": "/path/to/claude-plugins/diataxis-docs"
    },
    {
      "name": "developer-tools",
      "source": "/path/to/claude-plugins/developer-tools"
    }
  ]
}
```

Restart your Claude Code session to load all plugins.

### 2. Workflow: Research into Reports

Combine the `research` plugin's investigation agents with its report generator to produce structured deliverables.

**Step A -- Gather findings with the comprehensive researcher:**

```
@comprehensive-researcher Investigate current best practices for WebSocket
connection pooling in Go microservices. Cover performance benchmarks,
production patterns, and common pitfalls.
```

The comprehensive-researcher agent performs multi-source investigation and returns structured findings.

**Step B -- Generate a report from the findings:**

```
@report-generator Create a technical report from the research above.
Include an executive summary, key findings, and recommendations.
```

The report-generator agent takes the research output and produces a formatted report with sections, citations, and actionable recommendations.

### 3. Workflow: Code Generation with Quality Review

Pair `programming-languages` agents with `code-quality` agents for a write-then-review cycle.

**Step A -- Generate code with a language expert:**

```
@go-simplifier Implement a connection pool manager for WebSocket connections
with configurable max connections, health checking, and graceful shutdown.
```

The go-simplifier agent writes idiomatic Go code following language best practices.

**Step B -- Review the generated code:**

```
@code-reviewer Review the connection pool implementation above. Check for
race conditions, resource leaks, and error handling gaps.
```

The code-reviewer agent examines the code for architectural issues, bugs, and improvement opportunities.

**Step C -- Write tests for the reviewed code:**

```
@test-writer-fixer Write table-driven tests for the connection pool manager.
Cover normal operation, pool exhaustion, and connection timeout scenarios.
```

The test-writer-fixer agent generates comprehensive tests targeting the code that was just written and reviewed.

### 4. Workflow: Development with Documentation

Combine `developer-tools` or `programming-languages` agents with `diataxis-docs` agents to produce code and its documentation together.

**Step A -- Build the feature:**

```
@typescript-pro Implement a plugin loader module that reads plugin.json
manifests from a directory and validates their schema.
```

**Step B -- Generate a how-to guide for the feature:**

```
@doc-howto-writer Write a how-to guide for loading and validating plugins
using the module implemented above. Target developers integrating with
the plugin system.
```

The doc-howto-writer agent creates a Diataxis-compliant how-to guide with numbered steps, prerequisites, and verification.

**Step C -- Generate API reference documentation:**

```
@doc-reference-gen Generate reference documentation for the plugin loader
module's public API, covering all exported functions, types, and error codes.
```

### 5. Chain Outputs Across Steps

Each agent in the session can see the full conversation history. Outputs from earlier agents serve as inputs to later ones without any manual copy-paste. To get the best results:

- Run agents in sequence within the same session, not in parallel
- Reference earlier outputs explicitly: "Review the code above" or "Document the implementation from Step A"
- Keep each request focused on one action so the receiving agent has clear input

## Verify It Works

After running a multi-plugin workflow:

- [ ] Each agent responded with output appropriate to its role
- [ ] Later agents referenced and built upon earlier outputs
- [ ] The final deliverable (report, reviewed code, documented feature) is coherent and complete
- [ ] No agent reported missing context or unresolved references

## Troubleshooting

**Agent from one plugin cannot see output from another:**

- Confirm both plugins are listed in `.claude/settings.json`
- Run both agents in the same session -- agents share conversation context within a session but not across sessions

**Agent produces generic output, ignoring earlier context:**

- Reference the earlier output explicitly in your prompt: "Using the Go implementation above, write tests for..."
- Keep the conversation focused -- very long sessions may cause earlier context to become less prominent

**Conflicting advice from two agents:**

- Prefer the domain-specific agent's recommendation (e.g., trust `go-simplifier` on Go idioms over `code-reviewer` on Go-specific style)
- Use the more specialized agent for the final pass

**Too many plugins slow down session startup:**

- Install only the plugins needed for your current task
- Remove unused plugins from `.claude/settings.json` when switching projects

**Agents repeat work already done by a previous agent:**

- Be explicit about the division of labor: "Do not re-implement the code. Only review what was written above."
- Assign clear, non-overlapping roles in each prompt

## Next Steps

- [How to Install a Plugin]({{< ref "howto/install-a-plugin" >}}) -- add or remove plugins as your workflow changes
- [How to Find the Right Plugin]({{< ref "howto/find-right-plugin" >}}) -- discover additional plugins to extend your workflow
- [Marketplace Architecture]({{< ref "explanation/marketplace-architecture" >}}) -- understand how plugins, agents, and skills interact
- [Plugin Design Philosophy]({{< ref "explanation/plugin-design-philosophy" >}}) -- why plugins are structured this way
