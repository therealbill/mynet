---
title: "Optimize SQL Queries"
description: "Use sql-pro to analyze execution plans and improve query performance"
weight: 2
---

# Optimize SQL Queries

Identify and fix slow database queries using sql-pro's execution plan analysis, index recommendations, and query rewrites.

## Prerequisites

- Claude Code with the backend-development plugin installed
- A slow query or a query you suspect could be faster
- Access to the database (or a representative copy) for running EXPLAIN ANALYZE

## Steps

### 1. Identify the Slow Query

Gather the query and its context before engaging sql-pro:

- The full SQL query text
- The database dialect (PostgreSQL, MySQL, SQL Server -- sql-pro defaults to PostgreSQL if not specified)
- How slow it is now and how fast it needs to be
- Approximate table sizes and row counts
- Any existing indexes on the involved tables

If you do not have the execution plan yet, sql-pro will generate one. If you already have EXPLAIN ANALYZE output, include it.

### 2. Trigger sql-pro

Present the slow query with context. Example prompts that activate sql-pro:

- "This query takes 30 seconds, can you optimize it?"
- "Write a query to find the top 10 customers by revenue"
- "The billing report query is timing out in production"

Be specific about what "slow" means. "Takes 30 seconds but needs to be under 1 second" gives sql-pro a clear target. "It's slow" requires a follow-up question.

```
This PostgreSQL query takes 12 seconds on a table with 50 million rows.
The users table has indexes on (id) and (email). Here's the query:

SELECT u.id, u.email, COUNT(o.id) as order_count, SUM(o.total) as revenue
FROM users u
JOIN orders o ON o.user_id = u.id
WHERE o.created_at > NOW() - INTERVAL '30 days'
GROUP BY u.id, u.email
ORDER BY revenue DESC
LIMIT 100;
```

### 3. Review the EXPLAIN ANALYZE Output

sql-pro runs EXPLAIN ANALYZE on the query (or interprets the plan you provided) and identifies bottlenecks:

- **Sequential scans** on large tables where an index scan would be faster
- **Hash joins** that spill to disk due to insufficient `work_mem`
- **Sort operations** that materialize large intermediate results
- **Nested loops** with high row multipliers
- **Missing index** on join or filter columns

The agent explains each bottleneck in plain language, not just the raw plan output. It connects the plan nodes to specific parts of your query so you understand which clause causes which cost.

### 4. Apply Optimizations

sql-pro recommends specific changes, which may include:

- **Index additions:** Each recommendation includes the exact CREATE INDEX statement and an explanation of the read/write tradeoff. sql-pro never recommends an index without explaining its cost to write operations.
- **Query rewrites:** Restructuring joins, replacing correlated subqueries with CTEs, or reordering WHERE clauses to leverage indexes better.
- **CTE refactoring:** Breaking complex queries into named CTEs for readability and sometimes for performance (materializing intermediate results).
- **Configuration changes:** Adjusting `work_mem`, `effective_cache_size`, or other server parameters when the bottleneck is resource configuration rather than query structure.

Apply the recommended changes. If sql-pro suggests multiple optimizations, apply them incrementally so you can measure the impact of each one.

### 5. Verify the Improvement

Run EXPLAIN ANALYZE again on the optimized query and compare:

- **Execution time:** The primary metric. Did it meet the target?
- **Plan changes:** Verify that sequential scans became index scans, that sorts use indexes, and that joins use appropriate strategies
- **Row estimates vs. actuals:** Large discrepancies between estimated and actual rows indicate stale statistics (`ANALYZE` the table) or a need for extended statistics

Ask sql-pro to compare the before and after plans if the results are not clear:

```
Here's the EXPLAIN ANALYZE output after adding the index. Is this better?
```

The agent interprets both plans side by side and confirms whether the optimization succeeded or suggests further changes.

## Verification

A successful optimization shows:

- [ ] Execution time meets the target
- [ ] EXPLAIN ANALYZE confirms the expected plan changes (index scans, efficient joins)
- [ ] Row estimates are close to actual row counts
- [ ] Any new indexes are justified with read/write tradeoff analysis

## Troubleshooting

**Indexes do not help:**

Not all slow queries are fixed by indexes. If the query processes a large fraction of the table (more than 10-20% of rows), PostgreSQL intentionally chooses a sequential scan because it is faster than thousands of random index lookups. In this case, sql-pro may recommend:

- Partitioning the table by the filter column (e.g., `created_at` range partitioning)
- Materialized views for expensive aggregations that do not need real-time data
- Query restructuring to reduce the working set before aggregation

**The problem is schema design, not the query:**

Sometimes a slow query reveals a schema problem. If sql-pro determines that no amount of indexing or query rewriting will meet the performance target, it escalates to schema-level recommendations:

- Denormalizing specific columns to eliminate expensive joins
- Adding summary tables maintained by triggers or batch jobs
- Restructuring the table layout for the actual access pattern

For schema redesign, backend-architect may be more appropriate than sql-pro, as schema decisions affect the entire service architecture.

**EXPLAIN ANALYZE shows good performance but the app is still slow:**

The bottleneck may not be the query itself. Connection pool exhaustion, ORM overhead, N+1 query patterns, or network latency between the application and database are common causes. sql-pro focuses on individual query performance; backend-architect can help diagnose application-level database performance issues.

## See Also

- [Getting Started]({{< ref "tutorials/getting-started" >}}) -- tutorial showing sql-pro in the context of a full backend project
- [Agent Reference]({{< ref "reference/agents" >}}) -- sql-pro specification and capabilities
- [Design API Architecture]({{< ref "howto/design-api-architecture" >}}) -- when the performance problem requires architectural changes
