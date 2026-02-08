---
title: "Design API Architecture"
description: "Use backend-architect to design service boundaries, API contracts, and observability"
weight: 1
---

# Design API Architecture

Get a complete architecture design -- service boundaries, Goa API contracts, database choice, and observability strategy -- before writing any code.

## Prerequisites

- Claude Code with the backend-development plugin installed
- A clear idea of what the system needs to do, even if the details are rough

## Steps

### 1. Describe Your Requirements

Tell backend-architect what the system needs to do, at what scale, and with what constraints. Include:

- The core problem the service solves
- Expected scale (users, requests per second, data volume)
- Integration points (existing services, third-party APIs, databases already in use)
- Any hard constraints (compliance, latency targets, deployment environment)

Example prompts that trigger backend-architect:

- "We need an API for social sharing features"
- "Our queries are getting slow as the dataset grows"
- "Should this be a microservice or part of the monolith?"

The more context you provide up front, the fewer clarifying questions the agent needs to ask.

### 2. Trigger backend-architect

backend-architect activates when it detects API design, system architecture, or database scaling decisions in your request. If it does not activate automatically, be explicit:

```
Design the architecture for a notification service that handles email, SMS,
and push notifications for 500,000 users with delivery confirmation tracking
```

The agent starts with clarifying questions. Answer them -- these questions determine whether the design ends up as a monolith or microservice, PostgreSQL or DynamoDB, synchronous or event-driven.

### 3. Review Service Boundaries and Data Flow

backend-architect produces a service boundary diagram describing:

- Which components belong in the same service and which should be separate
- How data flows between components (synchronous API calls, message queues, shared databases)
- Where the API surface lives and what it exposes

For most systems at moderate scale, backend-architect defaults to a modular monolith rather than microservices. It separates into microservices only when there is a concrete reason: independent scaling requirements, different team ownership, or fundamentally different runtime characteristics.

Challenge the boundaries if they do not feel right. Ask "why is this one service and not two?" or "what happens when the notification volume spikes 10x?"

### 4. Review API Contracts

backend-architect designs APIs using Goa DSL. The output includes:

- Endpoint definitions with request/response types
- Error types and status codes
- Authentication and authorization requirements
- Automatic OpenAPI specification generation

Verify that the API contracts cover all the use cases from your requirements. Missing endpoints at this stage are easier to add than after implementation has started.

### 5. Review the Observability Strategy

Every architecture from backend-architect includes an observability section. This is not optional. The strategy covers:

- **What to trace:** Which operations get spans, which span attributes carry business context
- **Key span attributes:** Identifiers and values that enable debugging and analysis in Honeycomb (e.g., `customer_id`, `notification_type`, `delivery_status`)
- **Alerting thresholds:** What latency or error rate should trigger alerts
- **Logging policy:** What goes to logs (critical errors only) vs. what goes to traces (everything else)

If the observability section is missing or thin, ask backend-architect to expand it. Instrumenting after deployment is significantly harder than designing instrumentation into the architecture.

### 6. Iterate If Needed

Architecture is iterative. Common adjustments:

- **Monolith vs. microservice:** If backend-architect chose a monolith but you need independent scaling, explain the scaling requirement and ask for a revised design
- **Database choice:** If the agent chose PostgreSQL but your access pattern is pure key-value at high throughput, push back with the specific pattern
- **API framework:** backend-architect defaults to Goa. If your team uses a different Go framework, state that constraint before the design starts

Each iteration refines the architecture. When you are satisfied, the design is ready for go-architect to implement.

## Verification

A complete architecture from backend-architect includes:

- [ ] Service boundaries with clear justification
- [ ] Goa API contracts covering all use cases
- [ ] Database choice with reasoning for the access pattern
- [ ] Observability strategy with span attributes, trace boundaries, and alerting
- [ ] Infrastructure recommendations (deployment target, caching, queues if needed)

## Troubleshooting

**backend-architect suggests Go but you need a different language:**

backend-architect defaults to Go because the plugin's implementation agent (go-architect) is a Go specialist. If you need a different language, state that constraint explicitly. backend-architect can still design the architecture -- service boundaries, API contracts, and observability apply regardless of language -- but the implementation handoff to go-architect will not apply.

**The design feels too abstract to implement:**

Ask backend-architect to produce Goa DSL for the API contracts. Concrete DSL definitions bridge the gap between architecture and implementation and give go-architect a precise starting point.

**You need a design review, not a new architecture:**

Paste the existing architecture or code and ask backend-architect to review it. The agent evaluates existing designs against the same criteria it uses for new ones: service boundaries, observability, API contract quality.

## See Also

- [Getting Started](../../tutorials/getting-started/) -- full tutorial walking through the design-then-implement workflow
- [Architecture](../../explanation/architecture/) -- why architecture-first and how the agents divide responsibilities
- [Agent Reference](../../reference/agents/) -- backend-architect specification
