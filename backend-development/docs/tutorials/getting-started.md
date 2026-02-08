---
title: "Getting Started with Backend Development"
description: "Design an API with backend-architect and implement it in Go with go-architect"
weight: 1
---

# Getting Started with Backend Development

Design a backend API service using the architecture-first workflow: backend-architect defines the system, go-architect implements it in Go, and sql-pro handles complex database queries.

## What You'll Build

By the end of this tutorial, you will have:

- Triggered backend-architect with a service description and received a complete architecture
- Reviewed service boundaries, Goa API contracts, and an observability strategy
- Handed off the architecture to go-architect for Go implementation
- Used sql-pro for a complex analytical query
- Understood the design-then-implement workflow that the backend-development plugin follows

This takes about 30 minutes.

## Prerequisites

- Claude Code CLI installed and authenticated (run `claude --version` to verify)
- The backend-development plugin installed in your project's `.claude/settings.json`
- An idea for a backend service (this tutorial uses a customer billing API as the example)

## Step 1: Describe Your Service Requirements

Open Claude Code in your project directory and describe what you need:

```
We need an API for customer billing. It handles subscription plans, usage tracking,
invoice generation, and payment status. We expect about 10,000 active subscriptions
and need month-end invoice runs that process all accounts.
```

Claude Code matches your request to the backend-architect agent based on the API design context. The agent activates and begins gathering requirements.

## Step 2: Answer the Architect's Questions

backend-architect asks clarifying questions before designing the system. Expect questions about:

- Scale: How many requests per second? How large is the dataset?
- Constraints: Must it integrate with an existing system? Any compliance requirements?
- Access patterns: Are reads or writes dominant? Do you need real-time updates?
- Existing infrastructure: What databases and cloud services are already in use?

Answer these questions. backend-architect uses your responses to make decisions about service boundaries, database choice, and API shape.

## Step 3: Review the Architecture Output

backend-architect produces a complete architecture proposal. For a billing API, expect:

- **Service boundaries:** A single service (monolith) handling subscriptions, usage, invoices, and payments -- not microservices, because the domain is tightly coupled at this scale
- **API contracts:** Goa DSL definitions for endpoints like `POST /subscriptions`, `GET /invoices/{id}`, `POST /invoices/generate` -- design-first with automatic OpenAPI generation
- **Database choice:** PostgreSQL for relational billing data with strong consistency guarantees
- **Data model:** Tables for customers, subscriptions, usage_events, invoices, and payments with foreign key relationships
- **Observability strategy:** OpenTelemetry spans on every API request and database query, key span attributes like `customer_id`, `invoice_amount`, and `subscription_plan`, traces exported to Honeycomb

Review the proposal carefully. The observability strategy is not an afterthought -- backend-architect treats it as a required output. Every architecture includes what to trace, which span attributes matter, and where alerts should fire.

### Checkpoint

At this point you should have:

- Triggered backend-architect with your service description
- Answered the agent's clarifying questions
- Received a complete architecture covering API contracts, database design, and observability

If backend-architect did not activate, verify the backend-development plugin is listed in your `.claude/settings.json` and restart Claude Code.

## Step 4: Hand Off to go-architect

With the architecture defined, ask for implementation:

```
Implement the billing service based on this architecture
```

Claude Code routes this to go-architect, which reads the architecture output and implements it in Go. go-architect produces:

- **Goa service design:** DSL files defining the API types, endpoints, and error responses
- **Generated code:** Goa-generated HTTP transport, OpenAPI spec, and client code
- **Business logic:** Handler implementations for subscription CRUD, usage recording, and invoice generation
- **OpenTelemetry instrumentation:** Spans on every handler and database call, with the span attributes that backend-architect specified
- **Tests:** Table-driven tests with subtests for each endpoint
- **Linting:** Code passes `golangci-lint run` and `go vet` with no issues

Review the generated Go code. Verify that Goa endpoints match the API contracts from the architecture, that error handling follows Go conventions (explicit error returns, no panics), and that OpenTelemetry spans include the business-context attributes.

### Checkpoint

At this point you should have:

- A Go service implementing the billing API architecture
- Goa-generated HTTP transport and OpenAPI specification
- OpenTelemetry instrumentation matching the observability strategy
- Passing tests and clean linting output

## Step 5: Use sql-pro for Complex Queries

The invoice generation endpoint needs a query that calculates total usage per customer for a billing period, applies tiered pricing, and returns the top accounts by revenue. Ask:

```
Write a PostgreSQL query to calculate monthly invoice amounts for all active
subscriptions, applying tiered pricing based on usage, and show the top 10
customers by revenue with month-over-month growth
```

Claude Code routes this to sql-pro. The agent:

1. Confirms the dialect (PostgreSQL)
2. Writes the query using CTEs for readability -- one CTE for usage aggregation, one for tier calculation, one for revenue comparison
3. Runs EXPLAIN ANALYZE on the query
4. Recommends indexes on `usage_events(customer_id, period)` and `subscriptions(status)` with a clear explanation of the read/write tradeoff

Review the query. sql-pro produces readable SQL with CTEs rather than deeply nested subqueries, and every index recommendation includes justification.

## What You Learned

In this tutorial, you:

- **Triggered backend-architect** by describing an API requirement -- the agent activated and produced a complete architecture before any code was written
- **Reviewed architecture-first output** including service boundaries, Goa API contracts, database design, and observability strategy
- **Handed off to go-architect** for implementation -- the agent read the architecture and produced idiomatic Go with Goa endpoints, tests, and instrumentation
- **Used sql-pro** for a complex analytical query -- the agent wrote readable SQL with CTEs and justified its index recommendations

The core workflow is: **backend-architect designs, go-architect implements, sql-pro optimizes queries.** Architecture decisions happen before code. Observability is included from the start, not bolted on later.

## Next Steps

- [Design API Architecture]({{< ref "howto/design-api-architecture" >}}) -- detailed guide for working with backend-architect
- [Optimize SQL Queries]({{< ref "howto/optimize-sql-queries" >}}) -- use sql-pro to fix slow queries
- [Architecture]({{< ref "explanation/architecture" >}}) -- understand the design-implement separation and when to use each agent
- [Agent Reference]({{< ref "reference/agents" >}}) -- full specification of all three agents

For deployment, see the [devops-and-infra](/devops-and-infra/) plugin. For CLI tooling around your service, see the [cli-development](/cli-development/) plugin.
