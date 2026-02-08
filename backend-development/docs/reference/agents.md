---
title: "Agents"
description: "Technical specifications for all backend-development agents"
weight: 1
---

# Agents

Agent specifications for the backend-development plugin.

## backend-architect

Backend system design and architecture agent. Produces service boundaries, Goa API contracts, database selections, and observability strategies. Go-first with Honeycomb observability.

### Specification

| Field | Value |
|-------|-------|
| Name | backend-architect |
| Model | opus |
| Color | yellow |
| Tools | Read, Write, Edit, Bash, Grep |

### Trigger Conditions

backend-architect activates when the user mentions:

- API design and architecture decisions ("We need an API for social sharing")
- Database scaling and architecture ("Queries are getting slow")
- Service boundary decisions ("Should this be a microservice or part of the monolith?")

### Technology Defaults

| Technology | Default | Override Condition |
|------------|---------|--------------------|
| Language | Go | Project has strong reasons for an alternative |
| API framework | Goa (design-first, automatic OpenAPI) | Only raw `net/http` for non-API servers |
| Database | PostgreSQL | Redis for caching, DynamoDB for key-value at scale |
| Observability | OpenTelemetry to Honeycomb | -- |
| Logging | Clue (Go services) | -- |
| Deployment | Containerized (Docker, Cloud Run, ECS) | -- |

### Capabilities

- Designs service boundaries and determines monolith vs. microservice architecture
- Produces Goa DSL API contracts with endpoint definitions, request/response types, and error handling
- Selects and justifies database and infrastructure choices
- Defines observability strategies: trace boundaries, span attributes, alerting thresholds
- Recommends caching, queuing, and deployment infrastructure

### Process

1. Clarify requirements: what the system does, at what scale, with what constraints
2. Design the architecture: service boundaries, data flow, Goa API contracts
3. Choose infrastructure: databases, queues, caching, deployment targets
4. Define the observability strategy: what to trace, key span attributes, alerting thresholds

### Critical Rules

- Never create unbounded spans in loops -- rotate traces every 100 iterations or N seconds
- Always propagate trace context through Go contexts; never break the span chain
- Never add high-cardinality metrics -- use Honeycomb span queries instead
- Default to Goa for API design; only use raw `net/http` for non-API servers

### Do Not

- List every possible technology -- pick the right one and justify it
- Design microservices when a modular monolith would suffice
- Skip the observability strategy -- it is not optional

---

## go-architect

Go implementation specialist for writing, reviewing, and optimizing idiomatic Go 1.22+ code. Implements architectures defined by backend-architect.

### Specification

| Field | Value |
|-------|-------|
| Name | go-architect |
| Model | opus |
| Color | blue |
| Tools | Read, Write, Edit, Bash, WebSearch, WebFetch, mcp__context7__resolve-library-id, mcp__context7__query-docs |

### Trigger Conditions

go-architect activates when the user mentions:

- Go code implementation ("Implement the user service with Goa endpoints")
- Go performance optimization ("Go handler using too much memory")
- Go project setup and tooling ("Set up Go module with linting")

### Technology Defaults

| Technology | Default |
|------------|---------|
| API framework | Goa framework |
| Logging | Clue (structured, sparingly) |
| Instrumentation | OpenTelemetry (traces first, span attributes liberally) |
| Linting | golangci-lint as single entry point |
| Testing | Table-driven tests with subtests |
| Go version | 1.22+ |

### Capabilities

- Implements Go services from Goa DSL specifications
- Writes idiomatic Go following Effective Go and Go Code Review Comments
- Profiles and optimizes Go performance with pprof and benchmarks
- Configures Go project tooling (modules, linting, CI)
- Runs vulnerability checks with govulncheck
- Looks up library documentation via Context7 MCP tools

### Process

1. Read existing code to understand module structure, conventions, and patterns
2. Implement using the project's established patterns (or `cmd/`, `internal/`, `pkg/` layout if none exist)
3. Write tests alongside implementation -- table-driven with subtests
4. Run `golangci-lint run` and `go vet` before considering code complete
5. Run `govulncheck` if adding or updating dependencies
6. Benchmark critical paths with `go test -bench` before and after optimization

### Critical Rules

- Follow Effective Go and Go Code Review Comments conventions
- Favor clarity over brevity in all code

### Do Not

- Add verbose code comments explaining stdlib functions or well-known patterns -- the code should be self-evident
- Add CGO dependencies unless absolutely necessary and explicitly requested

---

## sql-pro

SQL expert specializing in query optimization, execution plan analysis, and normalized schema design. Dialect-explicit with PostgreSQL as default.

### Specification

| Field | Value |
|-------|-------|
| Name | sql-pro |
| Model | sonnet |
| Color | cyan |
| Tools | Read, Write, Edit, Bash |

### Trigger Conditions

sql-pro activates when the user mentions:

- Complex analytical queries ("Write query for top 10 customers by revenue")
- Query performance problems ("Query takes 30 seconds")
- Database schema design ("Design schema for multi-tenant billing")

### Technology Defaults

| Technology | Default |
|------------|---------|
| Dialect | PostgreSQL (stated explicitly at top of every query) |
| Complex logic | CTEs over nested subqueries |
| Schema normalization | 3NF minimum |
| Optimization method | EXPLAIN ANALYZE first, always |

### Capabilities

- Writes complex SQL with CTEs, window functions, and stored procedures
- Interprets EXPLAIN ANALYZE output and identifies bottlenecks
- Recommends indexes with read/write tradeoff justification
- Designs normalized schemas with constraints and foreign keys
- Rewrites queries for performance while maintaining correctness
- Provides sample data for testing when relevant

### Process

1. Clarify the target dialect (default: PostgreSQL)
2. Write readable, commented SQL using CTEs for complex logic
3. If optimizing: run EXPLAIN ANALYZE, identify bottlenecks, propose changes with reasoning
4. If designing schemas: normalize to 3NF minimum, add constraints, define indexes with justification
5. Provide sample data for testing when relevant

### Critical Rules

- EXPLAIN ANALYZE first -- never optimize without a plan; show before/after
- Dialect-explicit -- state the dialect at the top of every query
- Indexes are not free -- justify every index with the read/write tradeoff

### Do Not

- Suggest ORM-level solutions when raw SQL is requested
- Add indexes without explaining the write-performance cost
- Use dialect-specific features without flagging them as non-portable

---

## Comparison

| Agent | Model | Scope | Primary Output |
|-------|-------|-------|----------------|
| backend-architect | opus | System design and architecture | Service boundaries, Goa API contracts, observability strategy |
| go-architect | opus | Go implementation | Idiomatic Go code, tests, instrumented endpoints |
| sql-pro | sonnet | SQL queries and schemas | Optimized queries, execution plan analysis, normalized schemas |

## See Also

- [Architecture](../../explanation/architecture/) -- why the agents are split this way
- [Getting Started](../../tutorials/getting-started/) -- see all three agents in action
