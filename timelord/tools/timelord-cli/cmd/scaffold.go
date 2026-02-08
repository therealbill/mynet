package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var scaffoldCmd = &cobra.Command{
	Use:   "scaffold",
	Short: "Scaffold Temporal application components",
	Long:  `Generate boilerplate code for Temporal projects, workflows, activities, workers, and Nexus services.`,
}

var scaffoldProjectCmd = &cobra.Command{
	Use:   "project [name]",
	Short: "Scaffold a new Temporal project",
	Long: `Generate a complete Temporal project structure with:
  - Go module setup
  - Workflow and activity directories
  - Worker configuration
  - Makefile and Dockerfile
  - Test setup`,
	Args: cobra.ExactArgs(1),
	RunE: runScaffoldProject,
}

var scaffoldWorkflowCmd = &cobra.Command{
	Use:   "workflow [name]",
	Short: "Scaffold a new workflow",
	Long:  `Generate a workflow file with standard structure and patterns.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runScaffoldWorkflow,
}

var scaffoldActivityCmd = &cobra.Command{
	Use:   "activity [name]",
	Short: "Scaffold a new activity",
	Long:  `Generate an activity file with timeout and retry configuration.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runScaffoldActivity,
}

var scaffoldWorkerCmd = &cobra.Command{
	Use:   "worker [name]",
	Short: "Scaffold a new worker",
	Long:  `Generate a worker file with registration and configuration.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runScaffoldWorker,
}

var scaffoldNexusServiceCmd = &cobra.Command{
	Use:   "nexus-service [name]",
	Short: "Scaffold a Nexus service handler",
	Long: `Generate a Nexus service handler file with:
  - Input/output types
  - A sync operation (nexus.NewSyncOperation)
  - An async operation (temporalnexus.NewWorkflowRunOperation)
  - A backing workflow function
  - A NewXxxService() constructor`,
	Args: cobra.ExactArgs(1),
	RunE: runScaffoldNexusService,
}

var scaffoldNexusCallerCmd = &cobra.Command{
	Use:   "nexus-caller [name]",
	Short: "Scaffold a Nexus caller workflow",
	Long: `Generate a Nexus caller workflow file with:
  - Input/output types
  - A caller workflow using workflow.NewNexusClient() and ExecuteOperation()
  - Error handling for NexusOperationFailure
  - scheduleToCloseTimeout configured`,
	Args: cobra.ExactArgs(1),
	RunE: runScaffoldNexusCaller,
}

var (
	outputDir   string
	taskQueue   string
	packageName string
	endpoint    string
	serviceName string
)

func init() {
	rootCmd.AddCommand(scaffoldCmd)
	scaffoldCmd.AddCommand(scaffoldProjectCmd)
	scaffoldCmd.AddCommand(scaffoldWorkflowCmd)
	scaffoldCmd.AddCommand(scaffoldActivityCmd)
	scaffoldCmd.AddCommand(scaffoldWorkerCmd)
	scaffoldCmd.AddCommand(scaffoldNexusServiceCmd)
	scaffoldCmd.AddCommand(scaffoldNexusCallerCmd)

	scaffoldProjectCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory")
	scaffoldWorkflowCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory")
	scaffoldWorkflowCmd.Flags().StringVarP(&taskQueue, "task-queue", "t", "default", "Task queue name")
	scaffoldActivityCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory")
	scaffoldWorkerCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory")
	scaffoldWorkerCmd.Flags().StringVarP(&taskQueue, "task-queue", "t", "default", "Task queue name")
	scaffoldWorkerCmd.Flags().StringVarP(&packageName, "package", "p", "", "Package name for imports")
	scaffoldNexusServiceCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory")
	scaffoldNexusCallerCmd.Flags().StringVarP(&outputDir, "output", "o", ".", "Output directory")
	scaffoldNexusCallerCmd.Flags().StringVar(&endpoint, "endpoint", "my-endpoint", "Nexus endpoint name")
	scaffoldNexusCallerCmd.Flags().StringVar(&serviceName, "service-name", "", "Nexus service name (defaults to lowercase of scaffold name)")
}

type ScaffoldResult struct {
	Type    string   `json:"type"`
	Name    string   `json:"name"`
	Files   []string `json:"files"`
	Message string   `json:"message"`
}

func outputResult(result ScaffoldResult) {
	if jsonOutput {
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Printf("Created %s: %s\n", result.Type, result.Name)
		if verbose {
			for _, f := range result.Files {
				fmt.Printf("  - %s\n", f)
			}
		}
		if result.Message != "" {
			fmt.Println(result.Message)
		}
	}
}

func runScaffoldProject(cmd *cobra.Command, args []string) error {
	name := args[0]
	projectDir := filepath.Join(outputDir, name)

	// Create directory structure
	dirs := []string{
		"workflows",
		"activities",
		"worker",
		"internal",
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(projectDir, dir), 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	files := []string{}

	// Create go.mod
	goModContent := fmt.Sprintf(`module %s

go 1.22

require (
	go.temporal.io/sdk v1.29.1
)
`, name)
	goModPath := filepath.Join(projectDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte(goModContent), 0644); err != nil {
		return fmt.Errorf("failed to write go.mod: %w", err)
	}
	files = append(files, "go.mod")

	// Create Makefile
	makefileContent := fmt.Sprintf(`PROJECT = %s
TASK_QUEUE = %s

.PHONY: build run-worker test lint

build:
	go build -o bin/worker ./worker

run-worker:
	go run ./worker

test:
	go test -v ./...

lint:
	golangci-lint run

docker-build:
	docker build -t $(PROJECT):latest .

.DEFAULT_GOAL := build
`, name, name)
	makefilePath := filepath.Join(projectDir, "Makefile")
	if err := os.WriteFile(makefilePath, []byte(makefileContent), 0644); err != nil {
		return fmt.Errorf("failed to write Makefile: %w", err)
	}
	files = append(files, "Makefile")

	// Create Dockerfile
	dockerfileContent := `FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /worker ./worker

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /worker .

CMD ["./worker"]
`
	dockerfilePath := filepath.Join(projectDir, "Dockerfile")
	if err := os.WriteFile(dockerfilePath, []byte(dockerfileContent), 0644); err != nil {
		return fmt.Errorf("failed to write Dockerfile: %w", err)
	}
	files = append(files, "Dockerfile")

	// Create sample workflow
	workflowData := TemplateData{
		Name:        "Sample",
		PackageName: name,
		TaskQueue:   name,
	}
	workflowContent, err := executeTemplate("workflow", workflowData)
	if err != nil {
		return err
	}
	workflowPath := filepath.Join(projectDir, "workflows", "sample_workflow.go")
	if err := os.WriteFile(workflowPath, []byte(workflowContent), 0644); err != nil {
		return fmt.Errorf("failed to write workflow: %w", err)
	}
	files = append(files, "workflows/sample_workflow.go")

	// Create sample activity
	activityData := TemplateData{
		Name:        "Sample",
		PackageName: name,
	}
	activityContent, err := executeTemplate("activity", activityData)
	if err != nil {
		return err
	}
	activityPath := filepath.Join(projectDir, "activities", "sample_activities.go")
	if err := os.WriteFile(activityPath, []byte(activityContent), 0644); err != nil {
		return fmt.Errorf("failed to write activity: %w", err)
	}
	files = append(files, "activities/sample_activities.go")

	// Create worker
	workerData := TemplateData{
		Name:        name,
		PackageName: name,
		TaskQueue:   name,
	}
	workerContent, err := executeTemplate("worker", workerData)
	if err != nil {
		return err
	}
	workerPath := filepath.Join(projectDir, "worker", "main.go")
	if err := os.WriteFile(workerPath, []byte(workerContent), 0644); err != nil {
		return fmt.Errorf("failed to write worker: %w", err)
	}
	files = append(files, "worker/main.go")

	outputResult(ScaffoldResult{
		Type:    "project",
		Name:    name,
		Files:   files,
		Message: fmt.Sprintf("\nNext steps:\n  cd %s\n  go mod tidy\n  make run-worker", name),
	})

	return nil
}

func runScaffoldWorkflow(cmd *cobra.Command, args []string) error {
	name := args[0]
	data := TemplateData{
		Name:      toExportedName(name),
		TaskQueue: taskQueue,
	}

	content, err := executeTemplate("workflow", data)
	if err != nil {
		return err
	}

	filename := toSnakeCase(name) + "_workflow.go"
	filepath := filepath.Join(outputDir, filename)

	if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write workflow: %w", err)
	}

	outputResult(ScaffoldResult{
		Type:  "workflow",
		Name:  name,
		Files: []string{filename},
	})

	return nil
}

func runScaffoldActivity(cmd *cobra.Command, args []string) error {
	name := args[0]
	data := TemplateData{
		Name: toExportedName(name),
	}

	content, err := executeTemplate("activity", data)
	if err != nil {
		return err
	}

	filename := toSnakeCase(name) + "_activities.go"
	filepath := filepath.Join(outputDir, filename)

	if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write activity: %w", err)
	}

	outputResult(ScaffoldResult{
		Type:  "activity",
		Name:  name,
		Files: []string{filename},
	})

	return nil
}

func runScaffoldWorker(cmd *cobra.Command, args []string) error {
	name := args[0]
	data := TemplateData{
		Name:        toExportedName(name),
		TaskQueue:   taskQueue,
		PackageName: packageName,
	}

	content, err := executeTemplate("worker", data)
	if err != nil {
		return err
	}

	filename := toSnakeCase(name) + "_worker.go"
	filepath := filepath.Join(outputDir, filename)

	if err := os.WriteFile(filepath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write worker: %w", err)
	}

	outputResult(ScaffoldResult{
		Type:  "worker",
		Name:  name,
		Files: []string{filename},
	})

	return nil
}

func runScaffoldNexusService(cmd *cobra.Command, args []string) error {
	name := args[0]
	exportedName := toExportedName(name)
	svcName := strings.ToLower(exportedName[:1]) + exportedName[1:]

	data := TemplateData{
		Name:        exportedName,
		ServiceName: svcName,
	}

	content, err := executeTemplate("nexus-service", data)
	if err != nil {
		return err
	}

	filename := toSnakeCase(name) + "_nexus_service.go"
	outPath := filepath.Join(outputDir, filename)

	if err := os.WriteFile(outPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write nexus service: %w", err)
	}

	outputResult(ScaffoldResult{
		Type:  "nexus-service",
		Name:  name,
		Files: []string{filename},
		Message: fmt.Sprintf("\nNext steps:\n"+
			"  1. Customize the input/output types\n"+
			"  2. Implement sync and async operation logic\n"+
			"  3. Register the service on your handler worker:\n"+
			"     service := New%sService()\n"+
			"     w.RegisterNexusService(service)\n"+
			"  4. Create a Nexus endpoint:\n"+
			"     temporal operator nexus endpoint create \\\n"+
			"       --name my-endpoint \\\n"+
			"       --target-namespace handler-ns \\\n"+
			"       --target-task-queue handler-tq", exportedName),
	})

	return nil
}

func runScaffoldNexusCaller(cmd *cobra.Command, args []string) error {
	name := args[0]
	exportedName := toExportedName(name)

	svcName := serviceName
	if svcName == "" {
		svcName = strings.ToLower(exportedName[:1]) + exportedName[1:]
	}

	data := TemplateData{
		Name:        exportedName,
		Endpoint:    endpoint,
		ServiceName: svcName,
	}

	content, err := executeTemplate("nexus-caller", data)
	if err != nil {
		return err
	}

	filename := toSnakeCase(name) + "_nexus_caller.go"
	outPath := filepath.Join(outputDir, filename)

	if err := os.WriteFile(outPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write nexus caller: %w", err)
	}

	outputResult(ScaffoldResult{
		Type:  "nexus-caller",
		Name:  name,
		Files: []string{filename},
		Message: fmt.Sprintf("\nNext steps:\n"+
			"  1. Customize the input/output types\n"+
			"  2. Ensure the Nexus endpoint '%s' exists:\n"+
			"     temporal operator nexus endpoint list\n"+
			"  3. Register the caller workflow on your worker:\n"+
			"     w.RegisterWorkflow(%sCallerWorkflow)\n"+
			"  4. Adjust scheduleToCloseTimeout for your use case", endpoint, exportedName),
	})

	return nil
}

// TemplateData holds the data passed to scaffold templates.
type TemplateData struct {
	Name        string
	PackageName string
	TaskQueue   string
	Endpoint    string
	ServiceName string
}

func executeTemplate(name string, data TemplateData) (string, error) {
	templates := map[string]string{
		"workflow":      workflowTemplate,
		"activity":      activityTemplate,
		"worker":        workerTemplate,
		"nexus-service": nexusServiceTemplate,
		"nexus-caller":  nexusCallerTemplate,
	}

	tmplContent, ok := templates[name]
	if !ok {
		return "", fmt.Errorf("unknown template: %s", name)
	}

	tmpl, err := template.New(name).Parse(tmplContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

func toExportedName(s string) string {
	if s == "" {
		return s
	}
	// Convert first character to uppercase
	return strings.ToUpper(s[:1]) + s[1:]
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

const workflowTemplate = `package workflows

import (
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// {{.Name}}WorkflowInput is the input for {{.Name}}Workflow
type {{.Name}}WorkflowInput struct {
	// Add your input fields here
	ID string
}

// {{.Name}}WorkflowOutput is the output for {{.Name}}Workflow
type {{.Name}}WorkflowOutput struct {
	// Add your output fields here
	Result string
}

// {{.Name}}Workflow orchestrates the {{.Name}} process
func {{.Name}}Workflow(ctx workflow.Context, input {{.Name}}WorkflowInput) (*{{.Name}}WorkflowOutput, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("{{.Name}}Workflow started", "input", input)

	// Configure activity options
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    5,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Execute activities
	// var result ActivityResult
	// err := workflow.ExecuteActivity(ctx, YourActivity, activityInput).Get(ctx, &result)
	// if err != nil {
	//     return nil, err
	// }

	return &{{.Name}}WorkflowOutput{
		Result: "completed",
	}, nil
}
`

const activityTemplate = `package activities

import (
	"context"

	"go.temporal.io/sdk/activity"
)

// {{.Name}}Activities contains activities for {{.Name}}
type {{.Name}}Activities struct {
	// Add dependencies here (database, HTTP client, etc.)
}

// New{{.Name}}Activities creates a new {{.Name}}Activities instance
func New{{.Name}}Activities() *{{.Name}}Activities {
	return &{{.Name}}Activities{}
}

// {{.Name}}ActivityInput is the input for {{.Name}}Activity
type {{.Name}}ActivityInput struct {
	// Add your input fields here
	ID string
}

// {{.Name}}ActivityOutput is the output for {{.Name}}Activity
type {{.Name}}ActivityOutput struct {
	// Add your output fields here
	Result string
}

// {{.Name}}Activity performs the {{.Name}} operation
func (a *{{.Name}}Activities) {{.Name}}Activity(ctx context.Context, input {{.Name}}ActivityInput) (*{{.Name}}ActivityOutput, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("{{.Name}}Activity started", "input", input)

	// For long-running activities, record heartbeat
	// activity.RecordHeartbeat(ctx, "processing")

	// Check for cancellation
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Implement your activity logic here

	return &{{.Name}}ActivityOutput{
		Result: "completed",
	}, nil
}
`

const workerTemplate = `package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
{{- if .PackageName }}
	"{{.PackageName}}/activities"
	"{{.PackageName}}/workflows"
{{- end }}
)

const taskQueue = "{{.TaskQueue}}"

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
	w := worker.New(c, taskQueue, worker.Options{
		// Configure worker options as needed
		// MaxConcurrentActivityExecutionSize: 100,
		// MaxConcurrentWorkflowTaskExecutionSize: 100,
	})

	// Register workflows
	// w.RegisterWorkflow(workflows.SampleWorkflow)

	// Register activities
	// activities := activities.NewSampleActivities()
	// w.RegisterActivity(activities)

	// Start worker
	log.Printf("Starting worker on task queue: %s", taskQueue)
	if err := w.Run(worker.InterruptCh()); err != nil {
		log.Fatalf("Worker failed: %v", err)
	}
}
`

const nexusServiceTemplate = `package nexus

import (
	"context"
	"fmt"

	"github.com/nexus-rpc/sdk-go/nexus"
	"go.temporal.io/sdk/client"
	temporalnexus "go.temporal.io/sdk/temporalnexus"
	"go.temporal.io/sdk/workflow"
)

// --- Input/Output Types ---

// {{.Name}}SyncInput is the input for the synchronous {{.Name}} operation
type {{.Name}}SyncInput struct {
	// Add your input fields here
	ID      string
	Message string
}

// {{.Name}}SyncOutput is the output for the synchronous {{.Name}} operation
type {{.Name}}SyncOutput struct {
	// Add your output fields here
	Result  string
	Message string
}

// {{.Name}}AsyncInput is the input for the asynchronous {{.Name}} operation
type {{.Name}}AsyncInput struct {
	// Add your input fields here
	ID string
}

// {{.Name}}AsyncOutput is the output for the asynchronous {{.Name}} operation
type {{.Name}}AsyncOutput struct {
	// Add your output fields here
	Result string
	Status string
}

// --- Synchronous Operation ---

// {{.Name}}EchoOp is a synchronous Nexus operation that completes immediately (< 10s).
var {{.Name}}EchoOp = nexus.NewSyncOperation("{{.ServiceName}}-echo",
	func(ctx context.Context, input {{.Name}}SyncInput, opts nexus.StartOperationOptions) ({{.Name}}SyncOutput, error) {
		if input.Message == "" {
			return {{.Name}}SyncOutput{}, nexus.HandlerErrorf(nexus.HandlerErrorTypeBadRequest, "message is required")
		}
		return {{.Name}}SyncOutput{
			Result:  "ok",
			Message: input.Message,
		}, nil
	})

// --- Asynchronous (Workflow-Backed) Operation ---

// {{.Name}}ProcessOp is an asynchronous Nexus operation backed by a workflow.
// It starts a workflow and returns when the workflow completes.
var {{.Name}}ProcessOp = temporalnexus.NewWorkflowRunOperation("{{.ServiceName}}-process",
	{{.Name}}ProcessWorkflow,
	func(ctx context.Context, input {{.Name}}AsyncInput, opts nexus.StartOperationOptions) (client.StartWorkflowOptions, error) {
		return client.StartWorkflowOptions{
			// Use RequestID as workflow ID to deduplicate retries
			ID: opts.RequestID,
		}, nil
	})

// {{.Name}}ProcessWorkflow is the backing workflow for the async Nexus operation.
func {{.Name}}ProcessWorkflow(ctx workflow.Context, input {{.Name}}AsyncInput) (*{{.Name}}AsyncOutput, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("{{.Name}}ProcessWorkflow started", "id", input.ID)

	// Implement your long-running workflow logic here.
	// You can use activities, child workflows, timers, etc.
	//
	// Example:
	// var result ActivityOutput
	// err := workflow.ExecuteActivity(ctx, MyActivity, input).Get(ctx, &result)
	// if err != nil {
	//     return nil, fmt.Errorf("activity failed: %w", err)
	// }

	return &{{.Name}}AsyncOutput{
		Result: fmt.Sprintf("processed-%s", input.ID),
		Status: "completed",
	}, nil
}

// --- Service Constructor ---

// New{{.Name}}Service creates a new Nexus service with all {{.Name}} operations registered.
// Register this on your handler worker:
//
//	service := New{{.Name}}Service()
//	w.RegisterNexusService(service)
func New{{.Name}}Service() *nexus.Service {
	service := nexus.NewService("{{.ServiceName}}")
	service.Register({{.Name}}EchoOp)
	service.Register({{.Name}}ProcessOp)
	return service
}
`

const nexusCallerTemplate = `package workflows

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

// --- Input/Output Types ---

// {{.Name}}CallerInput is the input for {{.Name}}CallerWorkflow
type {{.Name}}CallerInput struct {
	// Add your input fields here
	ID string
}

// {{.Name}}CallerOutput is the output for {{.Name}}CallerWorkflow
type {{.Name}}CallerOutput struct {
	// Add your output fields here
	Result string
	Status string
}

// {{.Name}}CallerWorkflow calls a Nexus operation in another namespace.
// It uses workflow.NewNexusClient to create a client bound to a specific
// Nexus endpoint and service, then executes an operation and waits for the result.
//
// Prerequisites:
//   - Nexus endpoint "{{.Endpoint}}" must exist and route to the handler namespace
//   - Handler workers must be running with the "{{.ServiceName}}" Nexus service registered
func {{.Name}}CallerWorkflow(ctx workflow.Context, input {{.Name}}CallerInput) (*{{.Name}}CallerOutput, error) {
	logger := workflow.GetLogger(ctx)
	logger.Info("{{.Name}}CallerWorkflow started", "id", input.ID)

	// Create Nexus client bound to the endpoint and service
	nexusClient := workflow.NewNexusClient("{{.Endpoint}}", "{{.ServiceName}}")

	// Execute the Nexus operation
	// This blocks until the operation completes (sync or async).
	// For async operations, Temporal tracks the backing workflow automatically.
	future := nexusClient.ExecuteOperation(ctx, "{{.ServiceName}}-process", input, workflow.NexusOperationOptions{
		// ScheduleToCloseTimeout is required - set it based on expected operation duration.
		// Without this, the operation may hang indefinitely.
		ScheduleToCloseTimeout: 10 * time.Minute,
	})

	var result {{.Name}}CallerOutput
	if err := future.Get(ctx, &result); err != nil {
		// Errors from Nexus operations include:
		// - NexusOperationFailure: application-level error from the handler
		// - Timeout: scheduleToCloseTimeout exceeded
		// - Cancellation: workflow or operation was canceled
		return nil, fmt.Errorf("nexus operation failed: %w", err)
	}

	logger.Info("{{.Name}}CallerWorkflow completed", "result", result)
	return &result, nil
}
`
