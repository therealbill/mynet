---
title: "Architecture"
description: "Why dedicated agents for MCP server management"
weight: 1
---

# Architecture

Why the mcp-development plugin uses two specialized agents instead of one, and how their responsibilities divide.

## What MCP Is

Model Context Protocol (MCP) is a standard for connecting AI assistants to external tools and services. An MCP server exposes tools, resources, and completions through a defined protocol. An MCP client -- such as Claude Code -- connects to one or more servers and makes their capabilities available during a session. The protocol supports multiple transport types: stdio for local processes, streamable HTTP for remote services, and SSE for server-push scenarios.

## Why Two Agents

Discovery and integration are distinct concerns with different knowledge requirements and different tool needs.

**mcp-registry-navigator** needs to search the web, read registry metadata, and evaluate community signals. Its work is research-oriented: given a service name, find every available MCP server, assess each one's quality, and produce a recommendation. It uses WebSearch to query registries and Read/Write/Edit to generate configuration files.

**mcp-integration-engineer** needs to run commands, test connections, and configure systems. Its work is implementation-oriented: given a chosen server and configuration, make the connection work, handle authentication, add fault tolerance, and verify end-to-end. It uses Bash to execute processes and test connectivity, and Read/Write/Edit to modify configuration files.

Combining these into a single agent would create a tool that is good at neither task. The registry navigator would carry Bash access it does not need. The integration engineer would carry WebSearch access it does not need. Worse, the system prompt would need to cover both registry evaluation criteria and connection debugging procedures, diluting focus in both areas.

## Registry Navigator's Role

The registry navigator activates proactively when you mention a service name in the context of MCP. It does not wait for you to ask "search the registry" -- a question like "Is there an MCP server for Notion?" is enough.

Its evaluation is opinionated. It rates servers as Recommended, Acceptable, or Avoid based on maintenance activity, security practices, and community health. It will not recommend a server with no commits in six months without explicitly flagging the risk. This prevents you from adopting abandoned software without informed consent.

It also handles the publishing side: if you have built an MCP server and want to submit it to a registry, the navigator validates your metadata, tool annotations, and version declarations before submission.

## Integration Engineer's Role

The integration engineer activates when the task shifts from "what should I use?" to "how do I connect it?" It takes configuration output from the registry navigator and turns it into a working integration.

Its defaults are production-oriented. It uses streamable HTTP transport unless the server only supports stdio. It adds exponential backoff retries to all remote calls. For non-critical servers, it fails open -- if a server is down, the system degrades gracefully rather than blocking entirely.

For multi-server scenarios -- where Slack triggers GitHub actions, or a database server feeds into an analytics server -- the integration engineer designs the orchestration architecture. It avoids synchronous chains through three or more servers, preferring event-driven patterns that are more resilient to partial failures.

## When to Use Which

The division is straightforward:

- **"Does X exist?"** -- registry navigator. Any question about what servers are available, whether they are trustworthy, or how to publish one.
- **"How do I connect X?"** -- integration engineer. Any question about making a server work, debugging a connection, or orchestrating multiple servers.

In practice, a typical session flows from navigator to engineer: find a server, evaluate it, generate configuration, then hand off to the engineer to connect, authenticate, and test. The agents do not need to coordinate explicitly -- the configuration file the navigator generates is the contract the engineer consumes.
