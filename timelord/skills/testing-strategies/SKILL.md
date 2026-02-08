---
name: testing-strategies
description: This skill should be used when the user asks about "testing workflows", "TestWorkflowEnvironment", "workflow tests", "activity mocking", "replay testing", "temporal testing", "unit test workflow", or needs guidance on testing Temporal applications.
version: 1.0.0
---

# Temporal Testing Strategies

Guidance for testing Temporal workflows and activities using the Go SDK testing framework.

## Testing Overview

Temporal provides a comprehensive testing framework:

- **Unit Testing**: Test workflows and activities in isolation
- **Integration Testing**: Test with a real Temporal server
- **Replay Testing**: Validate workflow determinism against history

## TestWorkflowEnvironment

The test environment simulates Temporal without a server.

### Basic Setup

```go
import (
    "testing"

    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"
    "go.temporal.io/sdk/testsuite"
)

type WorkflowTestSuite struct {
    suite.Suite
    testsuite.WorkflowTestSuite
}

func (s *WorkflowTestSuite) SetupTest() {
    s.SetLogger(log.NewNopLogger())
}

func TestWorkflowTestSuite(t *testing.T) {
    suite.Run(t, new(WorkflowTestSuite))
}
```

### Testing Workflows

```go
func (s *WorkflowTestSuite) TestOrderWorkflow_Success() {
    env := s.NewTestWorkflowEnvironment()

    // Mock activities
    env.OnActivity(ValidateOrder, mock.Anything, mock.Anything).Return(nil)
    env.OnActivity(ProcessPayment, mock.Anything, mock.Anything).Return(
        &PaymentResult{TransactionID: "txn-123"}, nil,
    )
    env.OnActivity(ShipOrder, mock.Anything, mock.Anything).Return(
        &ShipmentResult{TrackingNumber: "track-456"}, nil,
    )

    // Execute workflow
    order := Order{ID: "order-1", Amount: 100}
    env.ExecuteWorkflow(OrderWorkflow, order)

    // Assert completion
    s.True(env.IsWorkflowCompleted())
    s.NoError(env.GetWorkflowError())

    // Get result
    var result OrderResult
    s.NoError(env.GetWorkflowResult(&result))
    s.Equal("txn-123", result.TransactionID)
    s.Equal("track-456", result.TrackingNumber)
}
```

### Testing Failure Scenarios

```go
func (s *WorkflowTestSuite) TestOrderWorkflow_PaymentFails() {
    env := s.NewTestWorkflowEnvironment()

    // Mock successful validation
    env.OnActivity(ValidateOrder, mock.Anything, mock.Anything).Return(nil)

    // Mock payment failure
    env.OnActivity(ProcessPayment, mock.Anything, mock.Anything).Return(
        nil, errors.New("payment declined"),
    )

    // Execute workflow
    order := Order{ID: "order-1", Amount: 100}
    env.ExecuteWorkflow(OrderWorkflow, order)

    // Assert failure
    s.True(env.IsWorkflowCompleted())
    s.Error(env.GetWorkflowError())
    s.Contains(env.GetWorkflowError().Error(), "payment declined")
}
```

## Activity Mocking

Mock activities to isolate workflow logic.

### Basic Mocking

```go
// Mock with exact arguments
env.OnActivity(SendEmail, mock.Anything, EmailRequest{
    To:      "user@example.com",
    Subject: "Order Confirmed",
}).Return(nil)

// Mock with any arguments
env.OnActivity(SendEmail, mock.Anything, mock.Anything).Return(nil)
```

### Conditional Mocking

```go
// Different returns based on input
env.OnActivity(GetUser, mock.Anything, "user-1").Return(
    &User{ID: "user-1", Name: "Alice"}, nil,
)
env.OnActivity(GetUser, mock.Anything, "user-2").Return(
    &User{ID: "user-2", Name: "Bob"}, nil,
)
env.OnActivity(GetUser, mock.Anything, "user-unknown").Return(
    nil, errors.New("user not found"),
)
```

### Mock with Function

```go
// Dynamic mock behavior
callCount := 0
env.OnActivity(ProcessItem, mock.Anything, mock.Anything).Return(
    func(ctx context.Context, item Item) (*Result, error) {
        callCount++
        if item.Priority == "high" {
            return &Result{Status: "expedited"}, nil
        }
        return &Result{Status: "normal"}, nil
    },
)
```

## Testing Signals

Test workflows that receive signals.

```go
func (s *WorkflowTestSuite) TestApprovalWorkflow_Approved() {
    env := s.NewTestWorkflowEnvironment()

    // Register callback to send signal
    env.RegisterDelayedCallback(func() {
        env.SignalWorkflow("approval", ApprovalSignal{
            Approved: true,
            ApproverID: "manager-1",
        })
    }, time.Millisecond*100)

    // Execute workflow
    env.ExecuteWorkflow(ApprovalWorkflow, "doc-123")

    s.True(env.IsWorkflowCompleted())
    s.NoError(env.GetWorkflowError())

    var result ApprovalResult
    s.NoError(env.GetWorkflowResult(&result))
    s.True(result.Approved)
}

func (s *WorkflowTestSuite) TestApprovalWorkflow_Rejected() {
    env := s.NewTestWorkflowEnvironment()

    env.RegisterDelayedCallback(func() {
        env.SignalWorkflow("approval", ApprovalSignal{
            Approved: false,
            Reason:   "Budget exceeded",
        })
    }, time.Millisecond*100)

    env.ExecuteWorkflow(ApprovalWorkflow, "doc-123")

    s.True(env.IsWorkflowCompleted())

    var result ApprovalResult
    s.NoError(env.GetWorkflowResult(&result))
    s.False(result.Approved)
    s.Equal("Budget exceeded", result.RejectionReason)
}
```

## Testing Queries

Test query handlers in workflows.

```go
func (s *WorkflowTestSuite) TestStateMachineWorkflow_Query() {
    env := s.NewTestWorkflowEnvironment()

    // Start workflow (don't complete it)
    env.RegisterDelayedCallback(func() {
        // Query current state
        result, err := env.QueryWorkflow("getState")
        s.NoError(err)

        var state string
        s.NoError(result.Get(&state))
        s.Equal("pending", state)

        // Send signal to change state
        env.SignalWorkflow("approve", nil)
    }, time.Millisecond*50)

    env.RegisterDelayedCallback(func() {
        // Query state after signal
        result, err := env.QueryWorkflow("getState")
        s.NoError(err)

        var state string
        s.NoError(result.Get(&state))
        s.Equal("approved", state)

        // Complete the workflow
        env.SignalWorkflow("complete", nil)
    }, time.Millisecond*100)

    env.ExecuteWorkflow(StateMachineWorkflow, "item-123")
    s.True(env.IsWorkflowCompleted())
}
```

## Testing Timers

Test workflows with time-based logic.

```go
func (s *WorkflowTestSuite) TestReminderWorkflow_SendsReminders() {
    env := s.NewTestWorkflowEnvironment()

    remindersSent := 0
    env.OnActivity(SendReminder, mock.Anything, mock.Anything).Return(
        func(ctx context.Context, msg string) error {
            remindersSent++
            return nil
        },
    )

    // Execute workflow that sends 3 daily reminders
    env.ExecuteWorkflow(ReminderWorkflow, "task-123")

    s.True(env.IsWorkflowCompleted())
    s.Equal(3, remindersSent)
}

func (s *WorkflowTestSuite) TestTimeoutWorkflow_TimesOut() {
    env := s.NewTestWorkflowEnvironment()

    // No signal sent - should timeout
    env.ExecuteWorkflow(TimeoutWorkflow, time.Hour)

    s.True(env.IsWorkflowCompleted())
    s.Error(env.GetWorkflowError())
    s.Contains(env.GetWorkflowError().Error(), "timeout")
}
```

## Testing Child Workflows

Test parent-child workflow relationships.

```go
func (s *WorkflowTestSuite) TestParentWorkflow_ExecutesChildren() {
    env := s.NewTestWorkflowEnvironment()

    // Mock child workflow
    env.OnWorkflow(ChildWorkflow, mock.Anything).Return(
        &ChildResult{Success: true}, nil,
    )

    // Execute parent
    env.ExecuteWorkflow(ParentWorkflow, ParentInput{ItemCount: 3})

    s.True(env.IsWorkflowCompleted())
    s.NoError(env.GetWorkflowError())

    var result ParentResult
    s.NoError(env.GetWorkflowResult(&result))
    s.Equal(3, result.ProcessedCount)
}
```

## Testing Nexus Operations

### Testing Sync Operation Handlers

Sync operations are plain functions — test them directly:

```go
func TestEchoOp(t *testing.T) {
    result, err := echoHandler(context.Background(), EchoInput{Message: "hello"}, nexus.StartOperationOptions{})
    require.NoError(t, err)
    require.Equal(t, "hello", result.Message)
}
```

### Testing Async Operation Handlers

Async operations start workflows — test the backing workflow as a normal workflow test, and test the operation handler setup separately.

### Testing Caller Workflows with Nexus Mocks

Mock Nexus operations in the test environment:

```go
func (s *WorkflowTestSuite) TestCallerWorkflow_NexusSuccess() {
    env := s.NewTestWorkflowEnvironment()

    // Register Nexus operation handler for testing
    env.RegisterNexusService(testNexusService)

    env.ExecuteWorkflow(CallerWorkflow, CallerInput{OrderID: "order-1"})

    s.True(env.IsWorkflowCompleted())
    s.NoError(env.GetWorkflowError())
}
```

### Integration Testing Nexus

Test against a real dev server with multiple namespaces:

```bash
# Start dev server
temporal server start-dev

# Create test namespaces and endpoint
temporal operator namespace create --namespace test-caller
temporal operator namespace create --namespace test-handler
temporal operator nexus endpoint create \
  --name test-endpoint \
  --target-namespace test-handler \
  --target-task-queue test-handler-tq
```

Then run both handler and caller workers, execute the caller workflow, and verify the end-to-end result.

## Testing Activities Directly

Test activity functions in isolation.

```go
func TestSendEmail_Success(t *testing.T) {
    // Create mock SMTP client
    mockClient := &MockSMTPClient{}
    mockClient.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(nil)

    // Create activity with dependencies
    activities := &EmailActivities{client: mockClient}

    // Test activity directly
    err := activities.SendEmail(context.Background(), EmailRequest{
        To:      "user@example.com",
        Subject: "Test",
        Body:    "Hello",
    })

    require.NoError(t, err)
    mockClient.AssertExpectations(t)
}

func TestProcessPayment_InvalidAmount(t *testing.T) {
    activities := &PaymentActivities{}

    _, err := activities.ProcessPayment(context.Background(), PaymentRequest{
        Amount: -100, // Invalid
    })

    require.Error(t, err)
    require.Contains(t, err.Error(), "invalid amount")
}
```

## Replay Testing

Validate workflow determinism using recorded history.

### Recording History

Export history from running workflows:

```bash
temporal workflow show --workflow-id my-workflow-id --output json > history.json
```

### Replay Test

```go
func TestWorkflow_Replay(t *testing.T) {
    replayer := worker.NewWorkflowReplayer()

    // Register workflow
    replayer.RegisterWorkflow(OrderWorkflow)

    // Replay against recorded history
    err := replayer.ReplayWorkflowHistoryFromJSONFile(
        nil,
        "testdata/order_workflow_history.json",
    )

    require.NoError(t, err)
}
```

### Detecting Non-Determinism

Replay tests catch determinism issues:

```go
func TestWorkflow_ReplayDetectsNonDeterminism(t *testing.T) {
    replayer := worker.NewWorkflowReplayer()
    replayer.RegisterWorkflow(ModifiedOrderWorkflow) // Changed workflow

    err := replayer.ReplayWorkflowHistoryFromJSONFile(
        nil,
        "testdata/original_workflow_history.json",
    )

    // Will fail if workflow structure changed incompatibly
    if err != nil {
        t.Logf("Non-determinism detected: %v", err)
    }
}
```

## Integration Testing

Test against a real Temporal server.

```go
func TestIntegration_OrderWorkflow(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    // Connect to test server
    c, err := client.Dial(client.Options{
        HostPort: "localhost:7233",
    })
    require.NoError(t, err)
    defer c.Close()

    // Start worker
    w := worker.New(c, "test-queue", worker.Options{})
    w.RegisterWorkflow(OrderWorkflow)
    w.RegisterActivity(&RealActivities{})
    go w.Run(worker.InterruptCh())

    // Execute workflow
    run, err := c.ExecuteWorkflow(context.Background(), client.StartWorkflowOptions{
        ID:        "test-order-" + uuid.New().String(),
        TaskQueue: "test-queue",
    }, OrderWorkflow, Order{ID: "test-1", Amount: 50})
    require.NoError(t, err)

    // Wait for completion
    var result OrderResult
    err = run.Get(context.Background(), &result)
    require.NoError(t, err)
    require.NotEmpty(t, result.TransactionID)
}
```

## Test Organization

Structure tests by category:

```
myapp/
├── workflows/
│   ├── order_workflow.go
│   └── order_workflow_test.go      # Unit tests
├── activities/
│   ├── payment_activities.go
│   └── payment_activities_test.go  # Activity tests
├── testdata/
│   └── order_workflow_history.json # Replay histories
└── integration/
    └── workflow_integration_test.go # Integration tests
```

## Best Practices

**Unit Testing:**

- Mock all activities in workflow tests
- Test happy path and failure scenarios
- Test signal handling and queries
- Use deterministic time in tests

**Activity Testing:**

- Test activities independently with mocked dependencies
- Test error conditions and edge cases
- Verify idempotency where applicable

**Replay Testing:**

- Store histories in version control
- Run replay tests in CI for critical workflows
- Update histories when intentionally changing workflows

**Integration Testing:**

- Use test containers or local Temporal for CI
- Clean up test workflows after execution
- Use unique workflow IDs to avoid conflicts

## Additional Resources

### Reference Files

For detailed testing patterns, consult:

- **`references/test-patterns.md`** - Advanced testing patterns
- **`references/replay-testing.md`** - Comprehensive replay testing guide
