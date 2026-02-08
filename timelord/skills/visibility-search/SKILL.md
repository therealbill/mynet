---
name: visibility-search
description: This skill should be used when the user asks about "search attributes", "temporal search", "list workflows", "visibility query", "elasticsearch temporal", "find workflow", "workflow filter", or needs guidance on searching and filtering Temporal workflows.
version: 1.0.0
---

# Visibility and Search

Guidance for searching and filtering Temporal workflows using visibility features.

## Visibility Overview

Temporal visibility allows searching workflows by:

- **Standard attributes** - Built-in fields (status, type, time)
- **Custom search attributes** - User-defined indexed fields

### Visibility Stores

| Store | Features | Use Case |
|-------|----------|----------|
| Standard (SQL) | Basic queries | Development, small scale |
| Advanced (Elasticsearch) | Full queries, aggregations | Production, complex queries |

## Standard Attributes

Built-in searchable attributes:

| Attribute | Type | Description |
|-----------|------|-------------|
| `WorkflowId` | Keyword | Workflow identifier |
| `RunId` | Keyword | Run identifier |
| `WorkflowType` | Keyword | Workflow function name |
| `ExecutionStatus` | Keyword | Running, Completed, Failed, etc. |
| `StartTime` | Datetime | When workflow started |
| `CloseTime` | Datetime | When workflow completed |
| `ExecutionDuration` | Int | Duration in milliseconds |
| `TaskQueue` | Keyword | Task queue name |

## Query Syntax

### Basic Queries

```sql
-- By status
ExecutionStatus = 'Running'
ExecutionStatus = 'Failed'
ExecutionStatus = 'Completed'
ExecutionStatus = 'Terminated'
ExecutionStatus = 'TimedOut'

-- By workflow type
WorkflowType = 'OrderWorkflow'

-- By workflow ID
WorkflowId = 'order-12345'

-- By task queue
TaskQueue = 'order-processing'
```

### Time-Based Queries

```sql
-- Started after date
StartTime > '2024-01-01T00:00:00Z'

-- Started within range
StartTime BETWEEN '2024-01-01' AND '2024-01-31'

-- Completed recently
CloseTime > '2024-01-15T00:00:00Z'

-- Long-running workflows
ExecutionDuration > 3600000  -- > 1 hour in ms
```

### Compound Queries

```sql
-- AND conditions
WorkflowType = 'OrderWorkflow' AND ExecutionStatus = 'Failed'

-- OR conditions
ExecutionStatus = 'Failed' OR ExecutionStatus = 'TimedOut'

-- Complex combinations
(WorkflowType = 'OrderWorkflow' OR WorkflowType = 'PaymentWorkflow')
AND ExecutionStatus = 'Running'
AND StartTime > '2024-01-01'
```

### ORDER BY

```sql
-- Order by start time (descending)
ORDER BY StartTime DESC

-- Multiple columns
ORDER BY WorkflowType ASC, StartTime DESC
```

## CLI Usage

### List with Query

```bash
# Running workflows
temporal workflow list --query "ExecutionStatus='Running'"

# Failed order workflows
temporal workflow list \
  --query "WorkflowType='OrderWorkflow' AND ExecutionStatus='Failed'"

# Recent failures
temporal workflow list \
  --query "ExecutionStatus='Failed' AND CloseTime > '2024-01-15'" \
  --limit 20

# With output format
temporal workflow list \
  --query "ExecutionStatus='Running'" \
  --output json
```

### timelord CLI

```bash
timelord workflow list --query "ExecutionStatus='Failed'" --limit 10 --json
```

## Custom Search Attributes

### Creating Attributes

```bash
# Keyword (exact match)
temporal operator search-attribute create \
  --namespace default \
  --name CustomerId \
  --type Keyword

# Text (full-text search)
temporal operator search-attribute create \
  --namespace default \
  --name Description \
  --type Text

# Numeric types
temporal operator search-attribute create \
  --namespace default \
  --name OrderTotal \
  --type Double

temporal operator search-attribute create \
  --namespace default \
  --name ItemCount \
  --type Int

# Boolean
temporal operator search-attribute create \
  --namespace default \
  --name IsPriority \
  --type Bool

# Datetime
temporal operator search-attribute create \
  --namespace default \
  --name DueDate \
  --type Datetime

# Multiple values
temporal operator search-attribute create \
  --namespace default \
  --name Tags \
  --type KeywordList
```

### Setting Attributes in Workflow

**At workflow start:**

```go
options := client.StartWorkflowOptions{
    ID:        "order-12345",
    TaskQueue: "orders",
    SearchAttributes: map[string]interface{}{
        "CustomerId":  "cust-789",
        "OrderTotal":  99.99,
        "IsPriority":  true,
        "Tags":        []string{"electronics", "express"},
    },
}

we, err := c.ExecuteWorkflow(ctx, options, OrderWorkflow, order)
```

**During workflow execution:**

```go
func OrderWorkflow(ctx workflow.Context, order Order) error {
    // Update search attributes
    err := workflow.UpsertSearchAttributes(ctx, map[string]interface{}{
        "OrderStatus": "processing",
        "ItemCount":   len(order.Items),
    })
    if err != nil {
        return err
    }

    // Continue workflow...
    return nil
}
```

### Querying Custom Attributes

```sql
-- By customer
CustomerId = 'cust-789'

-- By price range
OrderTotal >= 100.00 AND OrderTotal < 500.00

-- By priority
IsPriority = true

-- By tags (KeywordList)
Tags = 'electronics'

-- Combining standard and custom
WorkflowType = 'OrderWorkflow'
AND CustomerId = 'cust-789'
AND ExecutionStatus = 'Running'
```

## Common Query Patterns

### Debugging Queries

```sql
-- All failed workflows today
ExecutionStatus = 'Failed'
AND CloseTime > '2024-01-15T00:00:00Z'

-- Stuck workflows (running for > 1 hour)
ExecutionStatus = 'Running'
AND StartTime < '2024-01-15T10:00:00Z'

-- Terminated workflows
ExecutionStatus = 'Terminated'
AND CloseTime > '2024-01-14'
```

### Business Queries

```sql
-- High-value orders
OrderTotal > 1000 AND ExecutionStatus = 'Running'

-- Customer's workflows
CustomerId = 'cust-123'

-- Priority items due soon
IsPriority = true AND DueDate < '2024-01-16'

-- Orders by region
Region = 'US-WEST' AND WorkflowType = 'OrderWorkflow'
```

### Operations Queries

```sql
-- Workflows on specific task queue
TaskQueue = 'critical-processing'

-- Long-running workflows
ExecutionDuration > 86400000  -- > 1 day

-- Workflows started by specific service
InitiatingService = 'api-gateway'
```

## Elasticsearch Setup

### Why Elasticsearch?

| Feature | Standard SQL | Elasticsearch |
|---------|-------------|---------------|
| Basic queries | ✓ | ✓ |
| ORDER BY | Limited | ✓ |
| COUNT queries | Limited | ✓ |
| Full-text search | ✗ | ✓ |
| Aggregations | ✗ | ✓ |
| Performance at scale | Limited | Excellent |

### Configuration

**Helm values for Elasticsearch:**

```yaml
elasticsearch:
  enabled: true
  replicas: 3

server:
  config:
    persistence:
      advancedVisibilityStore: es-visibility
```

### Verification

```bash
# Check visibility store
temporal operator cluster describe

# Verify search attributes
temporal operator search-attribute list
```

## Best Practices

### Search Attribute Design

| Do | Don't |
|----|-------|
| Index frequently queried fields | Index everything |
| Use appropriate types | Store JSON as Keyword |
| Plan attributes upfront | Add attributes ad-hoc |
| Keep names consistent | Use inconsistent naming |

### Query Optimization

| Do | Don't |
|----|-------|
| Use indexed fields | Scan all workflows |
| Add time bounds | Open-ended time queries |
| Limit results | Fetch unlimited rows |
| Use specific conditions | Use OR excessively |

### Naming Conventions

```
CustomerId       # Keyword - entity reference
OrderTotal       # Double - numeric value
IsPriority       # Bool - flag
ProcessingStage  # Keyword - enum value
Tags             # KeywordList - multiple values
Description      # Text - searchable content
DueDate          # Datetime - temporal value
```

## Troubleshooting

### Query Returns No Results

- Check attribute exists: `temporal operator search-attribute list`
- Verify attribute type matches query
- Check time zone in datetime queries
- Verify workflow actually set the attribute

### Slow Queries

- Add time bounds to queries
- Reduce result limits
- Use more specific conditions
- Consider Elasticsearch for complex queries

### Attribute Not Searchable

- Confirm attribute was created before workflow started
- Verify attribute name spelling
- Check namespace context

## Additional Resources

### Reference Files

For detailed search patterns, consult:

- **`references/query-examples.md`** - Common query patterns
- **`references/elasticsearch-setup.md`** - Advanced ES configuration
