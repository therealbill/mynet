---
title: "Find and Evaluate MCP Servers"
description: "Assess MCP server trustworthiness, maintenance, and capabilities"
weight: 1
---

# Find and Evaluate MCP Servers

Find the right MCP server for a service and verify it is production-ready before committing to it.

## Problem

You need an MCP server for a specific service but do not know which one to use, whether it is maintained, or whether it meets your security requirements.

## Solution

Use mcp-registry-navigator to search, evaluate, and generate configuration for a vetted server.

## Prerequisites

- The mcp-development plugin installed
- The name of the service you want to connect to

## Steps

### 1. Search by Service Name

Ask the registry navigator to find servers:

```
Find MCP servers for <service-name>.
```

The agent searches mcp.so, the modelcontextprotocol GitHub organization, and relevant package registries. It cross-references results across sources to filter out abandoned forks and duplicates.

### 2. Review Maintenance Activity

For each result, check:

- **Last commit date** -- servers with no commits in 6+ months are flagged as risky
- **Release cadence** -- regular releases indicate active maintenance
- **Open vs. closed issue ratio** -- a high ratio of unresolved issues signals neglect

Ask for more detail on a specific server if the initial summary is insufficient:

```
Show me the maintenance history for the <server-name> MCP server.
```

### 3. Assess Security

Ask the agent for a security evaluation:

```
Is the <server-name> MCP server secure enough for production use?
```

The agent checks:

- **Authentication support** -- does the server require and validate auth tokens?
- **Transport security** -- does it support HTTPS/TLS for remote connections?
- **Input validation** -- does it sanitize inputs to tools and resources?
- **Dependency freshness** -- are dependencies up to date, or are there known vulnerabilities?

### 4. Evaluate Community Signals

Look at broader adoption indicators:

- **Stars and forks** -- not definitive on their own, but low counts on an old repo suggest limited adoption
- **Documentation quality** -- a server with clear setup instructions and API documentation is safer to adopt
- **Issue responsiveness** -- maintainers who respond to bug reports within days are preferable to those with months-old unanswered issues

The agent weights these signals relative to each other. A well-documented server with moderate stars outranks a popular server with no documentation.

### 5. Generate Configuration

Once you have selected a server, ask for configuration:

```
Generate the MCP configuration for <server-name>.
```

The agent produces a `.mcp.json` block with:

- Correct transport type (stdio, streamable HTTP, or SSE)
- Capability declarations matching the server's actual features
- Environment variable placeholders for all secrets
- Retry and timeout defaults appropriate for the transport type

### 6. Verify the Configuration

Confirm the generated configuration is valid:

- The server package name and version resolve in the declared registry
- The transport type matches your deployment model (local process vs. remote service)
- No credentials are hardcoded -- all secrets reference environment variables
- Declared capabilities match what the server actually implements

Hand off to mcp-integration-engineer to test the connection end-to-end if needed:

```
Test the connection to <server-name>.
```

## Verification

After completing these steps you should have:

- [ ] A server selected with a clear Recommended or Acceptable rating
- [ ] A maintenance and security assessment on record
- [ ] A valid `.mcp.json` configuration with no hardcoded secrets
- [ ] A successful connection test confirming the server responds
