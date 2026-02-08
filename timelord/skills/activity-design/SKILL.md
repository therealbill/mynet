---
name: activity-design
description: This skill should be used when the user asks about "activity implementation", "activity timeouts", "activity retries", "heartbeat", "idempotent activity", "StartToCloseTimeout", "ScheduleToCloseTimeout", "HeartbeatTimeout", or needs guidance on designing and configuring Temporal activities.
version: 1.0.0
---

# Temporal Activity Design

Guidance for implementing robust, idempotent activities with proper timeout and retry configurations.

## Activity Purpose

Activities are the only way workflows can interact with the external world:

- HTTP/gRPC calls to external services
- Database read/write operations
- File system access
- Message queue publishing
- Third-party API integrations

Activities execute in workers, separate from workflows, enabling:

- Non-deterministic operations
- Long-running external calls
- Retries with different workers

## Timeout Configuration

Three timeout types control activity execution:

### ScheduleToCloseTimeout

Maximum time from scheduling to completion, including all retries.

```go
ao := workflow.ActivityOptions{
    ScheduleToCloseTimeout: 30 * time.Minute, // Total time budget
}
```

**Use when:** Total time budget matters more than individual attempt duration.

### StartToCloseTimeout

Maximum time for a single execution attempt.

```go
ao := workflow.ActivityOptions{
    StartToCloseTimeout: 5 * time.Minute, // Per-attempt limit
}
```

**Use when:** Need to limit individual attempt duration, allowing retries within total budget.

### HeartbeatTimeout

Maximum time between heartbeats before marking as failed.

```go
ao := workflow.ActivityOptions{
    HeartbeatTimeout: 30 * time.Second, // Must heartbeat every 30s
}
```

**Use when:** Long-running activities need progress monitoring.

### Timeout Selection Guide

| Activity Type | Recommended Configuration |
|--------------|---------------------------|
| Quick API call (<10s) | StartToCloseTimeout: 30s |
| Moderate operation (10s-2m) | StartToCloseTimeout: 5m |
| Long processing (>2m) | StartToCloseTimeout: 10m + HeartbeatTimeout: 30s |
| Batch job | ScheduleToCloseTimeout: 1h, HeartbeatTimeout: 1m |

### Complete Configuration Example

```go
ao := workflow.ActivityOptions{
    StartToCloseTimeout:    10 * time.Minute,
    HeartbeatTimeout:       30 * time.Second,
    ScheduleToCloseTimeout: 1 * time.Hour,
    RetryPolicy: &temporal.RetryPolicy{
        InitialInterval:    time.Second,
        BackoffCoefficient: 2.0,
        MaximumInterval:    time.Minute,
        MaximumAttempts:    10,
    },
}
ctx = workflow.WithActivityOptions(ctx, ao)
```

## Retry Policy Configuration

Configure intelligent retry behavior for transient failures.

### Retry Policy Fields

```go
RetryPolicy: &temporal.RetryPolicy{
    InitialInterval:        time.Second,      // First retry delay
    BackoffCoefficient:     2.0,              // Multiplier per retry
    MaximumInterval:        time.Minute,      // Max delay between retries
    MaximumAttempts:        5,                // 0 = unlimited retries
    NonRetryableErrorTypes: []string{"InvalidInputError"},
}
```

### Retry Timing Example

With InitialInterval=1s, BackoffCoefficient=2.0, MaximumInterval=60s:

| Attempt | Delay Before |
|---------|--------------|
| 1 | 0s (immediate) |
| 2 | 1s |
| 3 | 2s |
| 4 | 4s |
| 5 | 8s |
| 6 | 16s |
| 7 | 32s |
| 8+ | 60s (capped) |

### Non-Retryable Errors

Mark errors as non-retryable for invalid input or business logic failures:

```go
// Go: Return application error with non-retryable flag
func ProcessOrder(ctx context.Context, order Order) error {
    if order.Amount <= 0 {
        return temporal.NewApplicationError(
            "invalid order amount",
            "InvalidInputError",  // Matches NonRetryableErrorTypes
            nil,
        )
    }
    // Process order...
    return nil
}
```

## Heartbeats

Report activity progress for long-running operations.

### When to Use Heartbeats

- Activity duration > 30 seconds
- Processing items in a loop
- Need to detect worker crashes quickly
- Want to report progress

### Heartbeat Implementation

```go
func ProcessLargeFile(ctx context.Context, fileID string) error {
    file := openFile(fileID)
    defer file.Close()

    totalLines := countLines(file)
    processed := 0

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        // Process line
        processLine(scanner.Text())
        processed++

        // Heartbeat with progress details
        activity.RecordHeartbeat(ctx, HeartbeatProgress{
            Processed: processed,
            Total:     totalLines,
            Percent:   float64(processed) / float64(totalLines) * 100,
        })

        // Check for cancellation
        if ctx.Err() != nil {
            return ctx.Err()
        }
    }

    return scanner.Err()
}
```

### Heartbeat Details

Pass progress information for observability:

```go
type HeartbeatProgress struct {
    Processed int     `json:"processed"`
    Total     int     `json:"total"`
    Percent   float64 `json:"percent"`
    LastItem  string  `json:"lastItem"`
}

activity.RecordHeartbeat(ctx, HeartbeatProgress{
    Processed: 500,
    Total:     1000,
    Percent:   50.0,
    LastItem:  "item-xyz",
})
```

## Idempotency

Design activities to produce the same result when called multiple times.

### Why Idempotency Matters

Activities may execute multiple times due to:

- Worker crashes and restarts
- Retry policies after timeouts
- Network partitions

### Idempotency Strategies

**1. Idempotency Keys**

```go
func CreatePayment(ctx context.Context, req PaymentRequest) (*Payment, error) {
    // Use workflow-provided idempotency key
    idempotencyKey := req.IdempotencyKey

    // Check for existing payment
    existing, err := db.GetPaymentByKey(idempotencyKey)
    if err == nil {
        return existing, nil // Already processed
    }

    // Create new payment
    payment := &Payment{
        ID:             uuid.New().String(),
        IdempotencyKey: idempotencyKey,
        Amount:         req.Amount,
    }

    err = db.CreatePayment(payment)
    if err != nil {
        return nil, err
    }

    return payment, nil
}
```

**2. Database Transactions**

```go
func TransferFunds(ctx context.Context, transfer Transfer) error {
    return db.Transaction(func(tx *sql.Tx) error {
        // Check if already processed
        var exists bool
        tx.QueryRow("SELECT EXISTS(SELECT 1 FROM transfers WHERE id = $1)",
            transfer.ID).Scan(&exists)
        if exists {
            return nil // Already processed
        }

        // Perform transfer atomically
        _, err := tx.Exec("UPDATE accounts SET balance = balance - $1 WHERE id = $2",
            transfer.Amount, transfer.FromAccount)
        if err != nil {
            return err
        }

        _, err = tx.Exec("UPDATE accounts SET balance = balance + $1 WHERE id = $2",
            transfer.Amount, transfer.ToAccount)
        if err != nil {
            return err
        }

        // Record transfer
        _, err = tx.Exec("INSERT INTO transfers (id, from_account, to_account, amount) VALUES ($1, $2, $3, $4)",
            transfer.ID, transfer.FromAccount, transfer.ToAccount, transfer.Amount)
        return err
    })
}
```

**3. Conditional Writes**

```go
func UpdateInventory(ctx context.Context, req InventoryUpdate) error {
    // Use conditional update with version check
    result, err := db.Exec(`
        UPDATE inventory
        SET quantity = quantity - $1, version = version + 1
        WHERE product_id = $2 AND version = $3 AND quantity >= $1`,
        req.Quantity, req.ProductID, req.ExpectedVersion)

    if err != nil {
        return err
    }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        return temporal.NewApplicationError(
            "inventory update conflict",
            "ConflictError",
            nil,
        )
    }

    return nil
}
```

## Activity Best Practices

### Input/Output Design

Keep activity inputs and outputs serializable and reasonably sized:

```go
// Good: Structured, serializable input
type SendEmailInput struct {
    To      string   `json:"to"`
    Subject string   `json:"subject"`
    Body    string   `json:"body"`
    ReplyTo string   `json:"replyTo,omitempty"`
}

// Good: Structured output with relevant details
type SendEmailOutput struct {
    MessageID   string    `json:"messageId"`
    SentAt      time.Time `json:"sentAt"`
    Recipient   string    `json:"recipient"`
}

func SendEmail(ctx context.Context, input SendEmailInput) (*SendEmailOutput, error) {
    // Implementation
}
```

### Error Handling

Return meaningful errors with appropriate retry behavior:

```go
func CallExternalAPI(ctx context.Context, req APIRequest) (*APIResponse, error) {
    resp, err := httpClient.Do(req.ToHTTPRequest())
    if err != nil {
        // Network errors are retryable by default
        return nil, fmt.Errorf("network error: %w", err)
    }
    defer resp.Body.Close()

    switch resp.StatusCode {
    case 200:
        var result APIResponse
        json.NewDecoder(resp.Body).Decode(&result)
        return &result, nil

    case 400, 422:
        // Client error - don't retry
        return nil, temporal.NewApplicationError(
            "invalid request",
            "ValidationError",
            nil,
        )

    case 429:
        // Rate limited - retry with delay
        return nil, temporal.NewApplicationError(
            "rate limited",
            "RateLimitError",
            nil,
        )

    case 500, 502, 503:
        // Server error - retry
        return nil, fmt.Errorf("server error: %d", resp.StatusCode)

    default:
        return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
    }
}
```

### Cancellation Handling

Respect context cancellation for graceful shutdown:

```go
func LongRunningActivity(ctx context.Context, items []Item) error {
    for _, item := range items {
        // Check cancellation before each item
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        if err := processItem(item); err != nil {
            return err
        }

        activity.RecordHeartbeat(ctx, item.ID)
    }
    return nil
}
```

## Activity Registration

Register activities with the worker:

```go
func main() {
    c, _ := client.Dial(client.Options{})
    defer c.Close()

    w := worker.New(c, "task-queue", worker.Options{})

    // Register activities
    w.RegisterActivity(SendEmail)
    w.RegisterActivity(ProcessPayment)
    w.RegisterActivity(UpdateInventory)

    // Or register struct with methods
    w.RegisterActivity(&EmailActivities{
        client: smtpClient,
    })

    w.Run(worker.InterruptCh())
}

// Activity struct for dependency injection
type EmailActivities struct {
    client *smtp.Client
}

func (a *EmailActivities) SendEmail(ctx context.Context, input SendEmailInput) error {
    return a.client.Send(input.To, input.Subject, input.Body)
}
```

## Additional Resources

### Reference Files

For detailed patterns and examples, consult:

- **`references/timeout-patterns.md`** - Advanced timeout configurations
- **`references/idempotency-examples.md`** - Idempotency implementation patterns
