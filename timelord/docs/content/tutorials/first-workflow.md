---
title: "Your First Temporal Workflow"
weight: 1
---

# Your First Temporal Workflow

Build a complete order processing workflow from scratch using Go and Temporal.

## What You'll Learn

- Set up a Temporal development environment
- Create a workflow that orchestrates multiple activities
- Implement activities for external operations
- Configure a worker to execute workflows
- Test your workflow with the test framework

## Prerequisites

- Go 1.22 or later
- Docker and Docker Compose (for local Temporal)
- Basic Go knowledge

## Step 1: Start Temporal Locally

Start a local Temporal server using Docker Compose:

```bash
# Clone the Temporal docker-compose repository
git clone https://github.com/temporalio/docker-compose.git
cd docker-compose

# Start Temporal
docker-compose up -d
```

Access the Temporal Web UI at http://localhost:8080.

## Step 2: Create the Project

Use timelord to scaffold a new project:

```bash
timelord scaffold project orderapp
cd orderapp
go mod tidy
```

This creates:

```
orderapp/
├── go.mod
├── Makefile
├── Dockerfile
├── workflows/
│   └── sample_workflow.go
├── activities/
│   └── sample_activities.go
├── worker/
│   └── main.go
└── internal/
```

## Step 3: Define the Workflow

Edit `workflows/order_workflow.go`:

```go
package workflows

import (
    "time"

    "go.temporal.io/sdk/temporal"
    "go.temporal.io/sdk/workflow"

    "orderapp/activities"
)

type OrderInput struct {
    OrderID    string
    CustomerID string
    Items      []string
    Amount     float64
}

type OrderOutput struct {
    OrderID       string
    PaymentID     string
    ShipmentID    string
    Status        string
}

func OrderWorkflow(ctx workflow.Context, input OrderInput) (*OrderOutput, error) {
    logger := workflow.GetLogger(ctx)
    logger.Info("OrderWorkflow started", "orderID", input.OrderID)

    // Configure activity options
    ao := workflow.ActivityOptions{
        StartToCloseTimeout: 5 * time.Minute,
        RetryPolicy: &temporal.RetryPolicy{
            InitialInterval:    time.Second,
            BackoffCoefficient: 2.0,
            MaximumInterval:    time.Minute,
            MaximumAttempts:    3,
        },
    }
    ctx = workflow.WithActivityOptions(ctx, ao)

    // Create activities instance for type-safe execution
    var a *activities.OrderActivities

    // Step 1: Validate the order
    var validationResult activities.ValidationResult
    err := workflow.ExecuteActivity(ctx, a.ValidateOrder, activities.ValidateOrderInput{
        OrderID:    input.OrderID,
        CustomerID: input.CustomerID,
        Items:      input.Items,
        Amount:     input.Amount,
    }).Get(ctx, &validationResult)
    if err != nil {
        return nil, err
    }
    logger.Info("Order validated", "orderID", input.OrderID)

    // Step 2: Process payment
    var paymentResult activities.PaymentResult
    err = workflow.ExecuteActivity(ctx, a.ProcessPayment, activities.ProcessPaymentInput{
        OrderID:    input.OrderID,
        CustomerID: input.CustomerID,
        Amount:     input.Amount,
    }).Get(ctx, &paymentResult)
    if err != nil {
        return nil, err
    }
    logger.Info("Payment processed", "paymentID", paymentResult.PaymentID)

    // Step 3: Ship the order
    var shipmentResult activities.ShipmentResult
    err = workflow.ExecuteActivity(ctx, a.ShipOrder, activities.ShipOrderInput{
        OrderID:    input.OrderID,
        CustomerID: input.CustomerID,
        Items:      input.Items,
    }).Get(ctx, &shipmentResult)
    if err != nil {
        // Payment succeeded but shipping failed
        // In a real app, you might compensate (refund)
        return nil, err
    }
    logger.Info("Order shipped", "shipmentID", shipmentResult.ShipmentID)

    return &OrderOutput{
        OrderID:    input.OrderID,
        PaymentID:  paymentResult.PaymentID,
        ShipmentID: shipmentResult.ShipmentID,
        Status:     "completed",
    }, nil
}
```

## Step 4: Implement Activities

Edit `activities/order_activities.go`:

```go
package activities

import (
    "context"
    "fmt"

    "go.temporal.io/sdk/activity"
)

type OrderActivities struct {
    // Add dependencies like database, payment gateway, etc.
}

func NewOrderActivities() *OrderActivities {
    return &OrderActivities{}
}

// Validation

type ValidateOrderInput struct {
    OrderID    string
    CustomerID string
    Items      []string
    Amount     float64
}

type ValidationResult struct {
    Valid   bool
    Message string
}

func (a *OrderActivities) ValidateOrder(ctx context.Context, input ValidateOrderInput) (*ValidationResult, error) {
    logger := activity.GetLogger(ctx)
    logger.Info("Validating order", "orderID", input.OrderID)

    // Validate order data
    if input.Amount <= 0 {
        return &ValidationResult{
            Valid:   false,
            Message: "Invalid order amount",
        }, fmt.Errorf("invalid order amount: %f", input.Amount)
    }

    if len(input.Items) == 0 {
        return &ValidationResult{
            Valid:   false,
            Message: "No items in order",
        }, fmt.Errorf("no items in order")
    }

    return &ValidationResult{
        Valid:   true,
        Message: "Order validated",
    }, nil
}

// Payment

type ProcessPaymentInput struct {
    OrderID    string
    CustomerID string
    Amount     float64
}

type PaymentResult struct {
    PaymentID string
    Status    string
}

func (a *OrderActivities) ProcessPayment(ctx context.Context, input ProcessPaymentInput) (*PaymentResult, error) {
    logger := activity.GetLogger(ctx)
    logger.Info("Processing payment", "orderID", input.OrderID, "amount", input.Amount)

    // In a real app, call payment gateway here
    // For this tutorial, simulate success
    paymentID := fmt.Sprintf("PAY-%s", input.OrderID)

    return &PaymentResult{
        PaymentID: paymentID,
        Status:    "captured",
    }, nil
}

// Shipping

type ShipOrderInput struct {
    OrderID    string
    CustomerID string
    Items      []string
}

type ShipmentResult struct {
    ShipmentID     string
    TrackingNumber string
    Status         string
}

func (a *OrderActivities) ShipOrder(ctx context.Context, input ShipOrderInput) (*ShipmentResult, error) {
    logger := activity.GetLogger(ctx)
    logger.Info("Shipping order", "orderID", input.OrderID)

    // Heartbeat for longer operations
    activity.RecordHeartbeat(ctx, "preparing shipment")

    // In a real app, call shipping API here
    shipmentID := fmt.Sprintf("SHIP-%s", input.OrderID)

    return &ShipmentResult{
        ShipmentID:     shipmentID,
        TrackingNumber: fmt.Sprintf("TRACK-%s", input.OrderID),
        Status:         "shipped",
    }, nil
}
```

## Step 5: Configure the Worker

Edit `worker/main.go`:

```go
package main

import (
    "log"

    "go.temporal.io/sdk/client"
    "go.temporal.io/sdk/worker"

    "orderapp/activities"
    "orderapp/workflows"
)

const taskQueue = "orders"

func main() {
    // Create Temporal client
    c, err := client.Dial(client.Options{
        HostPort: client.DefaultHostPort,
    })
    if err != nil {
        log.Fatalf("Failed to create Temporal client: %v", err)
    }
    defer c.Close()

    // Create worker
    w := worker.New(c, taskQueue, worker.Options{})

    // Register workflows
    w.RegisterWorkflow(workflows.OrderWorkflow)

    // Register activities
    orderActivities := activities.NewOrderActivities()
    w.RegisterActivity(orderActivities)

    // Start worker
    log.Printf("Starting worker on task queue: %s", taskQueue)
    if err := w.Run(worker.InterruptCh()); err != nil {
        log.Fatalf("Worker failed: %v", err)
    }
}
```

## Step 6: Start the Workflow

Create `starter/main.go` to start workflows:

```go
package main

import (
    "context"
    "log"

    "go.temporal.io/sdk/client"

    "orderapp/workflows"
)

func main() {
    c, err := client.Dial(client.Options{})
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    defer c.Close()

    input := workflows.OrderInput{
        OrderID:    "ORDER-001",
        CustomerID: "CUST-123",
        Items:      []string{"item-1", "item-2"},
        Amount:     99.99,
    }

    options := client.StartWorkflowOptions{
        ID:        "order-" + input.OrderID,
        TaskQueue: "orders",
    }

    we, err := c.ExecuteWorkflow(context.Background(), options, workflows.OrderWorkflow, input)
    if err != nil {
        log.Fatalf("Failed to start workflow: %v", err)
    }

    log.Printf("Started workflow: WorkflowID=%s, RunID=%s", we.GetID(), we.GetRunID())

    // Wait for result
    var result workflows.OrderOutput
    if err := we.Get(context.Background(), &result); err != nil {
        log.Fatalf("Workflow failed: %v", err)
    }

    log.Printf("Workflow completed: %+v", result)
}
```

## Step 7: Run It

```bash
# Terminal 1: Start the worker
go run ./worker

# Terminal 2: Start a workflow
go run ./starter
```

Check the Temporal Web UI at http://localhost:8080 to see your workflow execution.

## Step 8: Write Tests

Create `workflows/order_workflow_test.go`:

```go
package workflows

import (
    "testing"

    "github.com/stretchr/testify/mock"
    "github.com/stretchr/testify/suite"
    "go.temporal.io/sdk/testsuite"

    "orderapp/activities"
)

type OrderWorkflowTestSuite struct {
    suite.Suite
    testsuite.WorkflowTestSuite
}

func (s *OrderWorkflowTestSuite) TestOrderWorkflow_Success() {
    env := s.NewTestWorkflowEnvironment()

    // Mock activities
    env.OnActivity((*activities.OrderActivities).ValidateOrder, mock.Anything, mock.Anything).Return(
        &activities.ValidationResult{Valid: true}, nil,
    )
    env.OnActivity((*activities.OrderActivities).ProcessPayment, mock.Anything, mock.Anything).Return(
        &activities.PaymentResult{PaymentID: "PAY-001"}, nil,
    )
    env.OnActivity((*activities.OrderActivities).ShipOrder, mock.Anything, mock.Anything).Return(
        &activities.ShipmentResult{ShipmentID: "SHIP-001", TrackingNumber: "TRACK-001"}, nil,
    )

    // Execute workflow
    input := OrderInput{
        OrderID:    "ORDER-001",
        CustomerID: "CUST-123",
        Items:      []string{"item-1"},
        Amount:     50.00,
    }
    env.ExecuteWorkflow(OrderWorkflow, input)

    // Assert
    s.True(env.IsWorkflowCompleted())
    s.NoError(env.GetWorkflowError())

    var result OrderOutput
    s.NoError(env.GetWorkflowResult(&result))
    s.Equal("ORDER-001", result.OrderID)
    s.Equal("PAY-001", result.PaymentID)
    s.Equal("completed", result.Status)
}

func TestOrderWorkflowTestSuite(t *testing.T) {
    suite.Run(t, new(OrderWorkflowTestSuite))
}
```

Run tests:

```bash
go test ./workflows/...
```

## Checkpoint

You've successfully:

- [x] Set up a local Temporal environment
- [x] Created a workflow orchestrating multiple activities
- [x] Implemented activities with proper patterns
- [x] Configured a worker
- [x] Written unit tests

## Next Steps

- Add error handling and saga compensation
- Implement workflow versioning
- Set up monitoring with Prometheus
- Deploy to Kubernetes

See the How-To guides for these advanced topics.
