---
title: "Getting Started with MCP Development"
description: "Find and configure an MCP server using the mcp-development agents"
weight: 1
---

# Getting Started with MCP Development

Walk through discovering, evaluating, and connecting an MCP server from scratch using the two mcp-development agents.

## What You'll Learn

- Use mcp-registry-navigator to find MCP servers by service name
- Evaluate search results for trustworthiness and maintenance
- Generate a valid MCP server configuration
- Use mcp-integration-engineer to set up authentication and test the connection

## Prerequisites

- The mcp-development plugin installed in your Claude Code environment
- A service you want to connect to (this tutorial uses Notion as an example)

## Step 1: Search for a Server

Start by asking a natural question:

```
Is there an MCP server for Notion?
```

The mcp-registry-navigator agent activates automatically. It searches across registries including mcp.so, the modelcontextprotocol GitHub organization, and package managers like npm and PyPI. It cross-references results to verify that servers appear in multiple sources before presenting them.

## Step 2: Evaluate the Results

The agent returns a list of candidate servers. For each one, it reports on:

- **Transport type** -- stdio, streamable HTTP, or SSE
- **Maintenance status** -- last commit date, release cadence, open vs. closed issue ratio
- **Security signals** -- authentication method, input validation, dependency freshness
- **Capabilities** -- which MCP features the server implements (tools, resources, completions)

Each server receives a rating: Recommended, Acceptable, or Avoid, with a one-sentence justification.

### Checkpoint

Review the results. You should see at least one server with a clear maintenance and security profile. If all candidates are rated Avoid, the agent will flag this and suggest alternatives.

## Step 3: Generate Configuration

Choose a server and ask:

```
Generate the MCP configuration for the Notion server by <author>.
```

The mcp-registry-navigator produces a `.mcp.json` configuration block with the correct transport settings, capability declarations, and placeholder values for credentials.

### Checkpoint

Verify the generated configuration contains:

- The correct server package name and version
- A transport type that matches your client (stdio for local, streamable HTTP for remote)
- Environment variable references for secrets -- never hardcoded API keys

## Step 4: Set Up Authentication

Now hand off to the integration agent:

```
Set up authentication for this Notion MCP server.
```

The mcp-integration-engineer agent takes over. It examines the configuration, identifies the required authentication method (typically an API token for Notion), and configures the connection:

- Sets environment variable names for credentials
- Configures transport security settings
- Adds retry logic with exponential backoff for remote servers

## Step 5: Test the Connection

Ask the integration engineer to verify:

```
Test the connection to the Notion MCP server.
```

The agent runs a minimal end-to-end check: it confirms the server starts, authentication succeeds, and at least one tool or resource responds.

### Checkpoint

You should see confirmation that:

- The server process started without errors
- Authentication completed successfully
- The server's advertised capabilities match what the configuration declares

## Summary

You used two agents with distinct responsibilities:

- **mcp-registry-navigator** found servers, evaluated their quality, and generated configuration. It handles the "what exists and is it trustworthy?" question.
- **mcp-integration-engineer** connected the server, configured authentication, and verified the integration. It handles the "how do I make it work?" question.

## Next Steps

- [Find and Evaluate MCP Servers]({{< relref "../howto/find-and-evaluate-mcp-servers" >}}) -- deeper evaluation criteria
- [Architecture]({{< relref "../explanation/architecture" >}}) -- why two agents instead of one
- [Agent Reference]({{< relref "../reference/agents" >}}) -- full specifications for both agents
