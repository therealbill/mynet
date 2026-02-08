---
title: "Architecture"
description: "Architecture-first vs code-first and when to use each agent"
weight: 1
---

# Architecture

The backend-development plugin separates system design from code implementation. This page explains why that separation exists, how the three agents divide responsibility, and when each one is the right choice.

## The Design-Implement Separation

backend-architect and go-architect operate on different levels of abstraction. backend-architect produces architecture: service boundaries, API contracts, database selections, and observability strategies. go-architect produces code: Go implementations of those contracts with tests, instrumentation, and linting.

This split mirrors how backend systems are actually built. The decision to use a monolith vs. microservices is fundamentally different from the decision to use a sync.Pool for allocation reuse. The first requires understanding business domains, scaling requirements, and organizational constraints. The second requires understanding Go runtime internals and profiling data. Combining both into a single agent would mean either the architecture suffers from implementation-level thinking or the implementation suffers from architecture-level abstraction.

backend-architect answers "what should we build and why." go-architect answers "how do we build it in Go."

## Why Both Use Opus

backend-architect runs on opus because architecture decisions require broad reasoning across business domains, technical tradeoffs, and long-term system evolution. These are high-stakes decisions that affect everything downstream.

go-architect also runs on opus, but with a different focus. Go implementation at production quality requires deep understanding of the language's concurrency model, memory management, standard library conventions, and the Goa framework's design-first approach. The reasoning demands are different in kind -- not broader, but deeper in a specific technical domain.

sql-pro runs on sonnet. SQL query optimization and schema design are well-constrained problems with clear inputs (a query, an execution plan, a set of tables) and clear outputs (an optimized query, an index recommendation, a normalized schema). The reasoning pattern is more analytical than architectural, and sonnet handles it efficiently.

## The Goa-Centric Workflow

backend-architect defaults to Goa for API design. This is not an arbitrary framework preference -- it reflects the plugin's architecture-first philosophy.

Goa is a design-first API framework. You define your API in a Go DSL -- endpoints, types, errors, authentication -- and Goa generates the HTTP transport, OpenAPI specification, and client code from that definition. The DSL serves as both the specification and the source of truth.

This aligns with how backend-architect works. The agent produces API contracts as part of the architecture. With Goa, those contracts are not documentation that drifts from implementation -- they are the implementation's foundation. When go-architect picks up the architecture and implements it, the Goa DSL is the bridge between design and code.

The generated OpenAPI specification also means the API is documented from day one, not after someone remembers to update a swagger file.

## Observability as a First-Class Concern

Every architecture from backend-architect includes an observability strategy. This is a deliberate design decision, not a checklist item.

The observability strategy answers specific questions:

- Which operations get spans? (Every API request, database query, and external call.)
- What business context goes into span attributes? (Customer IDs, operation types, amounts -- the data that makes traces useful for debugging production issues.)
- What alerting thresholds apply? (Latency percentiles, error rates, queue depths.)
- What goes into logs vs. traces? (Critical errors go to logs via Clue. Everything else goes to traces via OpenTelemetry to Honeycomb.)

backend-architect includes this because adding observability after a system is deployed is significantly harder than designing it in. Retrofitting spans into existing code requires understanding every code path, which is exactly the knowledge that is freshest during architecture design and fades rapidly after implementation.

The technology stack is opinionated: OpenTelemetry for instrumentation, Honeycomb for analysis. Traces are the primary signal. Metrics are derived from traces via Honeycomb's aggregation, not pre-computed. Logs are for the narrow set of problems that traces cannot capture.

## sql-pro as the Data Specialist

sql-pro fills a different role than backend-architect and go-architect. It is not part of the design-implement pipeline -- it is a specialist called in when the database is the bottleneck.

Three situations trigger sql-pro:

1. **A specific query is slow.** sql-pro runs EXPLAIN ANALYZE, interprets the execution plan, and recommends specific fixes: indexes, query rewrites, or configuration changes.
2. **A complex query needs writing.** Analytical queries with window functions, CTEs, and multi-table aggregations benefit from SQL-specific expertise rather than general-purpose Go implementation.
3. **A schema needs designing.** Normalization, constraints, foreign keys, and index strategies for a new data model.

sql-pro operates at the SQL level, not the application level. It does not know or care whether the query runs inside a Goa handler or a cron job. It cares about the execution plan, the data distribution, and the index strategy. When the problem is bigger than a single query -- when it requires rethinking the data model or service architecture -- backend-architect is the right agent.

## Cross-Plugin Relationships

The backend-development plugin produces backend services. Other plugins handle what happens before, after, and around those services:

- **programming-languages** provides go-simplifier for Go code cleanup and simplification. If go-architect produces code that works but feels overly complex, go-simplifier can reduce it to its essential form.
- **devops-and-infra** handles deployment, CI/CD pipelines, and infrastructure provisioning. backend-architect designs the system; devops-and-infra deploys it.
- **cli-development** handles command-line interfaces. If the backend service needs a CLI tool for administration, health checks, or data migration, cli-development provides agents for CLI argument parsing and terminal UI.
- **code-quality** provides testing agents. While go-architect writes unit tests inline, code-quality agents handle integration testing, end-to-end testing, and test strategy for larger systems.

These plugins do not import each other. They collaborate through Claude's dispatch layer: backend-architect may recommend deploying to Cloud Run, but the devops-and-infra plugin handles the actual deployment configuration. The boundary is clean -- each plugin owns its domain.

## When to Use Each Agent

**Use backend-architect when:**

- Starting a new service or feature that requires architectural decisions
- Deciding between monolith and microservices
- Choosing databases, caches, or message queues
- Designing API contracts before implementation
- Reviewing an existing architecture for improvements

**Use go-architect when:**

- Implementing a designed architecture in Go
- Writing or modifying Go code for existing services
- Profiling and optimizing Go performance
- Setting up Go project tooling and CI

**Use sql-pro when:**

- A specific query is slow and needs optimization
- Writing complex analytical queries
- Designing a new database schema
- Reviewing execution plans

The agents are not sequential -- you do not have to use backend-architect before go-architect on every task. If you are adding a new endpoint to an existing service with clear patterns, go-architect can implement it directly. backend-architect is for when the "what" and "why" are unclear, not for rubber-stamping obvious implementation work.

## See Also

- [Agent Reference]({{< ref "reference/agents" >}}) -- full specification of all three agents
- [Getting Started]({{< ref "tutorials/getting-started" >}}) -- see the design-implement workflow in practice
- [Design API Architecture]({{< ref "howto/design-api-architecture" >}}) -- practical guide for working with backend-architect
