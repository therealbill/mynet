---
title: "MCP Development"
description: "MCP server discovery, integration, and multi-server orchestration"
weight: 14
bookCollapseSection: true
---

# MCP Development

Model Context Protocol (MCP) server discovery, integration, and multi-server orchestration.

## Components

| Type | Name | Description |
|------|------|-------------|
| Agent | mcp-registry-navigator | MCP registry discovery — find servers, evaluate capabilities, generate configs |
| Agent | mcp-integration-engineer | MCP server integration — client-server setup, multi-server orchestration, fault tolerance |

## Documentation

- [Getting Started](tutorials/getting-started/) — Find and configure an MCP server
- [Find and Evaluate MCP Servers](howto/find-and-evaluate-mcp-servers/) — Assess trustworthiness and capabilities
- [Agent Reference](reference/agents/) — Both agents' specifications
- [Architecture](explanation/architecture/) — Why dedicated agents for MCP work
