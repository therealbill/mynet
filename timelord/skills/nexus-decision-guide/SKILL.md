---
name: nexus-decision-guide
description: Architecture decision framework for evaluating whether Temporal Nexus is right for your use case. Use when user asks "should I use nexus", "nexus vs child workflows", "nexus vs activities", "when to use nexus", "cross-namespace pattern", "multi-namespace architecture", "nexus tradeoffs", or "nexus benefits". Covers tradeoff analysis, complexity scoring, and migration paths. Do NOT use for Nexus implementation details or code examples — use nexus-operations instead.
version: 1.0.0
---

# Nexus Decision Guide

Guidance for deciding whether Temporal Nexus is the right communication pattern for your architecture.

## Decision Framework

### Step 1: Do You Need Cross-Namespace Communication?

If all your workflows run in a single namespace, **stop here — you don't need Nexus.**

- Same-namespace decomposition → **Child Workflows**
- Same-namespace messaging → **Signals, Queries, or Updates**
- Same-namespace external calls → **Activities**

### Step 2: What Drives the Namespace Separation?

| Reason | Nexus Fit |
|--------|-----------|
| Team ownership boundaries | Strong — Nexus provides clean service contracts |
| Security/compliance isolation | Strong — namespaces + Nexus maintain isolation |
| Independent scaling | Strong — separate workers, separate scaling |
| Environment separation (dev/staging/prod) | Not applicable — use same namespace structure per env |
| Just organizational preference | Weak — consider if the operational overhead is worth it |

### Step 3: What Are the Communication Requirements?

| Need | Best Pattern |
|------|-------------|
| Fire-and-forget notification | Signals (if same NS) or Activities calling Signal API |
| Quick data lookup (< 10s) | Nexus sync operation |
| Long-running cross-team process | Nexus async (workflow-backed) operation |
| Cross-namespace, no durability needed | Activity with Temporal client (simpler, less durable) |
| Cross-namespace with full durability | **Nexus** (the primary use case) |

### Step 4: Evaluate the Tradeoffs

**Nexus adds:**

- Nexus Endpoints (cluster-level routing resources to create and manage)
- Handler workers (separate workers for Nexus service registration)
- Service contracts (typed interfaces to define and maintain)
- Cross-namespace debugging complexity (tracing across caller + handler)

**Nexus provides:**

- At-least-once execution with automatic retries and circuit breaking
- Clean API boundaries between teams
- Independent deployment and scaling per namespace
- Full durable execution across namespace boundaries

**The question:** Does the isolation and durability benefit outweigh the operational complexity for your case?

## Architecture Patterns

### Pattern: Service Mesh via Nexus

Multiple teams each owning a namespace, connected via Nexus endpoints:

```
┌──────────────┐   Nexus: payments-ep   ┌──────────────┐
│  orders-ns   │──────────────────────>│ payments-ns  │
│  (Team A)    │                        │ (Team B)     │
│              │   Nexus: inventory-ep  ┌──────────────┐
│              │──────────────────────>│ inventory-ns │
└──────────────┘                       │ (Team C)     │
                                       └──────────────┘
```

**Good when:** Teams need independence, clear ownership, separate deployment cycles.

### Pattern: Shared Platform Service

Common platform service (e.g., notifications, payments) exposed via Nexus to all consuming namespaces:

```
┌────────────┐
│ service-a  │──┐
└────────────┘  │   Nexus: notify-ep   ┌────────────────┐
                ├────────────────────>│ notifications  │
┌────────────┐  │                     │ (platform team)│
│ service-b  │──┘                     └────────────────┘
└────────────┘
```

**Good when:** A central capability is consumed by many teams.

### Anti-Pattern: Nexus Within Single Namespace

Don't use Nexus if caller and handler are in the same namespace. Child workflows give you the same decomposition with less overhead.

### Anti-Pattern: Nexus for Simple External Calls

If you just need to call an HTTP API and don't need cross-namespace durability, use a regular Activity. Nexus is for Temporal-to-Temporal communication.

## Complexity Assessment

Rate your scenario:

| Factor | Low (1) | Medium (2) | High (3) |
|--------|---------|------------|----------|
| Number of consuming namespaces | 1 | 2-3 | 4+ |
| Handler operation complexity | Simple sync | Mix of sync/async | Complex chains |
| Team independence requirement | Nice to have | Important | Critical |
| Existing namespace topology | Single NS | Few NS, ad-hoc | Multi-NS established |

**Score 4-6:** Consider if simpler patterns (child workflows, activities) suffice.
**Score 7-9:** Nexus is likely a good fit.
**Score 10-12:** Nexus is strongly recommended.

## Migration Path

If you're currently using ad-hoc cross-namespace communication (activities calling Temporal client, etc.):

1. Define Nexus service contracts for existing cross-namespace calls
2. Create Nexus endpoints pointing to handler namespaces
3. Implement handler services wrapping existing workflows
4. Migrate callers one at a time to use NexusClient
5. Remove old activity-based cross-namespace glue code

## Additional Resources

- **nexus-operations** skill — detailed implementation guidance
- **namespace-management** skill — namespace setup and Nexus endpoint management
