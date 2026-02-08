---
title: "Agents"
description: "MCP development agent specifications"
weight: 1
---

# Agents

Technical specifications for the mcp-development plugin agents.

## mcp-registry-navigator

| Field | Value |
|-------|-------|
| Name | mcp-registry-navigator |
| Model | sonnet |
| Color | green |
| Tools | Read, Write, Edit, WebSearch |

**Trigger conditions:**

- Finding MCP servers for a named service
- Searching MCP registries (mcp.so, GitHub, npm, PyPI)
- Evaluating server capabilities, maintenance, and security
- Generating MCP server configurations
- Publishing servers to registries

**Capabilities:**

| Capability | Description |
|------------|-------------|
| Registry search | Search across mcp.so, modelcontextprotocol GitHub org, npm, and PyPI |
| Cross-reference | Verify servers appear in multiple sources before recommending |
| Maintenance assessment | Evaluate last commit date, release cadence, open/closed issue ratio |
| Security evaluation | Check auth support, transport security, input validation, dependency freshness |
| Capability analysis | Report which MCP features a server implements (tools, resources, completions) |
| Trust rating | Rate servers as Recommended, Acceptable, or Avoid with justification |
| Configuration generation | Produce `.mcp.json` blocks with transport, auth, and capability declarations |
| Publishing guidance | Validate `mcp.json` schema, tool annotations, version compatibility, README completeness |

**Constraints:**

- Does not recommend servers with no commits in 6+ months without flagging the risk
- Does not fabricate download counts, match rates, or performance metrics
- Searches relevant registries per query rather than listing all registries exhaustively

## mcp-integration-engineer

| Field | Value |
|-------|-------|
| Name | mcp-integration-engineer |
| Model | sonnet |
| Color | blue |
| Tools | Read, Write, Edit, Bash |

**Trigger conditions:**

- Connecting MCP servers to clients
- Multi-server orchestration and workflow design
- Authentication and transport security setup
- Fault tolerance and reconnection configuration
- Connection debugging and diagnostics

**Capabilities:**

| Capability | Description |
|------------|-------------|
| Client-server integration | Configure connections between MCP clients and servers |
| Transport configuration | Set up stdio, streamable HTTP, or SSE transport per server requirements |
| Auth configuration | Configure API tokens, OAuth, and environment-variable-based secret management |
| Multi-server orchestration | Design workflows where multiple MCP servers interact (e.g., Slack triggers GitHub actions) |
| Fault tolerance | Add circuit breakers for remote servers, reconnection logic for stdio |
| Error handling | Configure retry with exponential backoff, fail-open for non-critical servers |
| Connection diagnostics | Test end-to-end connectivity, identify authentication failures, debug transport issues |

**Defaults:**

| Setting | Default |
|---------|---------|
| Transport | Streamable HTTP (unless server only supports stdio) |
| Retry strategy | Exponential backoff for all remote server calls |
| Non-critical server failure | Fail-open (degrade gracefully) |
| Secret storage | Environment variables (never inline in config files) |

**Constraints:**

- Does not assume all servers support the same transport -- verifies capabilities first
- Does not hardcode server URLs or credentials in configuration files
- Avoids synchronous round-trips through 3+ servers -- prefers event-driven patterns
