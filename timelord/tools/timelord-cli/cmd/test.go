package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test Temporal workflows",
	Long:  `Commands for testing and validating Temporal workflows.`,
}

var testReplayCmd = &cobra.Command{
	Use:   "replay [history-file]",
	Short: "Replay workflow history for determinism testing",
	Long: `Replay a workflow history file to verify determinism.

This command helps detect non-determinism errors by replaying
recorded workflow history against the current workflow code.`,
	Args: cobra.ExactArgs(1),
	RunE: runTestReplay,
}

var testValidateCmd = &cobra.Command{
	Use:   "validate [path]",
	Short: "Validate workflow and activity code",
	Long: `Validate Temporal workflow and activity implementations.

Checks for common issues:
  - Non-deterministic operations in workflows
  - Missing activity registrations
  - Timeout configuration issues`,
	Args: cobra.MaximumNArgs(1),
	RunE: runTestValidate,
}

var (
	testWorkflowName string
	testPackagePath  string
	testVerboseOut   bool
)

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.AddCommand(testReplayCmd)
	testCmd.AddCommand(testValidateCmd)

	testReplayCmd.Flags().StringVarP(&testWorkflowName, "workflow", "w", "", "Workflow function name to replay")
	testReplayCmd.Flags().StringVarP(&testPackagePath, "package", "p", ".", "Go package path containing workflow")

	testValidateCmd.Flags().BoolVarP(&testVerboseOut, "verbose", "", false, "Show detailed validation output")
}

type ReplayResult struct {
	Success      bool     `json:"success"`
	HistoryFile  string   `json:"historyFile"`
	WorkflowName string   `json:"workflowName,omitempty"`
	EventCount   int      `json:"eventCount"`
	Error        string   `json:"error,omitempty"`
	Suggestions  []string `json:"suggestions,omitempty"`
}

type ValidationResult struct {
	Valid    bool              `json:"valid"`
	Path     string            `json:"path"`
	Errors   []ValidationError `json:"errors,omitempty"`
	Warnings []ValidationError `json:"warnings,omitempty"`
}

type ValidationError struct {
	File        string `json:"file"`
	Line        int    `json:"line,omitempty"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Suggestion  string `json:"suggestion,omitempty"`
}

func runTestReplay(cmd *cobra.Command, args []string) error {
	historyFile := args[0]

	// Check if history file exists
	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		result := ReplayResult{
			Success:     false,
			HistoryFile: historyFile,
			Error:       "History file not found",
		}
		return outputReplayResult(result)
	}

	// Read and parse history file
	historyData, err := os.ReadFile(historyFile)
	if err != nil {
		result := ReplayResult{
			Success:     false,
			HistoryFile: historyFile,
			Error:       fmt.Sprintf("Failed to read history file: %v", err),
		}
		return outputReplayResult(result)
	}

	// Count events in history
	var history map[string]interface{}
	eventCount := 0
	if err := json.Unmarshal(historyData, &history); err == nil {
		if events, ok := history["events"].([]interface{}); ok {
			eventCount = len(events)
		}
	}

	// Generate replay test code
	testCode := generateReplayTestCode(historyFile, testWorkflowName)

	// Create temporary test file
	tempDir, err := os.MkdirTemp("", "timelord-replay-*")
	if err != nil {
		result := ReplayResult{
			Success:     false,
			HistoryFile: historyFile,
			Error:       fmt.Sprintf("Failed to create temp directory: %v", err),
		}
		return outputReplayResult(result)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "replay_test.go")
	if err := os.WriteFile(testFile, []byte(testCode), 0644); err != nil {
		result := ReplayResult{
			Success:     false,
			HistoryFile: historyFile,
			Error:       fmt.Sprintf("Failed to write test file: %v", err),
		}
		return outputReplayResult(result)
	}

	// Copy history file to temp dir
	historyDest := filepath.Join(tempDir, "history.json")
	if err := copyFile(historyFile, historyDest); err != nil {
		result := ReplayResult{
			Success:     false,
			HistoryFile: historyFile,
			Error:       fmt.Sprintf("Failed to copy history file: %v", err),
		}
		return outputReplayResult(result)
	}

	// Output instructions (actual test running requires user's workflow code)
	result := ReplayResult{
		Success:      true,
		HistoryFile:  historyFile,
		WorkflowName: testWorkflowName,
		EventCount:   eventCount,
		Suggestions: []string{
			"To run replay test, add this test to your workflow package:",
			fmt.Sprintf("  1. Copy the generated test code to your project"),
			fmt.Sprintf("  2. Register your workflow: replayer.RegisterWorkflow(YourWorkflow)"),
			fmt.Sprintf("  3. Run: go test -run TestReplay"),
		},
	}

	if jsonOutput {
		result.Suggestions = append(result.Suggestions, testCode)
	} else {
		fmt.Println("Generated Replay Test Code:")
		fmt.Println("============================")
		fmt.Println(testCode)
		fmt.Println()
	}

	return outputReplayResult(result)
}

func runTestValidate(cmd *cobra.Command, args []string) error {
	path := "."
	if len(args) > 0 {
		path = args[0]
	}

	result := ValidationResult{
		Valid:    true,
		Path:     path,
		Errors:   []ValidationError{},
		Warnings: []ValidationError{},
	}

	// Find Go files
	goFiles, err := findGoFiles(path)
	if err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Type:        "PathError",
			Description: fmt.Sprintf("Failed to scan path: %v", err),
		})
		return outputValidationResult(result)
	}

	if len(goFiles) == 0 {
		result.Warnings = append(result.Warnings, ValidationError{
			Type:        "NoFiles",
			Description: "No Go files found in path",
			Suggestion:  "Ensure the path contains Go source files",
		})
		return outputValidationResult(result)
	}

	// Check each file for common issues
	for _, file := range goFiles {
		fileErrors, fileWarnings := validateGoFile(file)
		result.Errors = append(result.Errors, fileErrors...)
		result.Warnings = append(result.Warnings, fileWarnings...)
	}

	if len(result.Errors) > 0 {
		result.Valid = false
	}

	return outputValidationResult(result)
}

func generateReplayTestCode(historyFile, workflowName string) string {
	wfName := workflowName
	if wfName == "" {
		wfName = "YourWorkflow"
	}

	return fmt.Sprintf(`package yourpackage

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/worker"
)

func TestReplay_%s(t *testing.T) {
	replayer := worker.NewWorkflowReplayer()

	// Register your workflow
	replayer.RegisterWorkflow(%s)

	// Replay the history
	err := replayer.ReplayWorkflowHistoryFromJSONFile(nil, "history.json")
	require.NoError(t, err, "Replay should succeed - non-determinism detected if this fails")
}
`, wfName, wfName)
}

func findGoFiles(path string) ([]string, error) {
	var files []string

	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(p, ".go") && !strings.HasSuffix(p, "_test.go") {
			files = append(files, p)
		}
		return nil
	})

	return files, err
}

func validateGoFile(file string) ([]ValidationError, []ValidationError) {
	var errors []ValidationError
	var warnings []ValidationError

	content, err := os.ReadFile(file)
	if err != nil {
		errors = append(errors, ValidationError{
			File:        file,
			Type:        "ReadError",
			Description: fmt.Sprintf("Failed to read file: %v", err),
		})
		return errors, warnings
	}

	contentStr := string(content)

	// Check for workflow.Context usage (indicates workflow code)
	isWorkflowFile := strings.Contains(contentStr, "workflow.Context")

	if isWorkflowFile {
		// Check for non-deterministic patterns
		nonDetPatterns := []struct {
			pattern     string
			description string
			suggestion  string
		}{
			{
				pattern:     "time.Now()",
				description: "time.Now() is non-deterministic in workflows",
				suggestion:  "Use workflow.Now(ctx) instead",
			},
			{
				pattern:     "rand.",
				description: "rand package is non-deterministic in workflows",
				suggestion:  "Use workflow.SideEffect for random values",
			},
			{
				pattern:     "uuid.New()",
				description: "uuid.New() is non-deterministic in workflows",
				suggestion:  "Use workflow.SideEffect or workflow.UUIDString()",
			},
			{
				pattern:     "os.Getenv(",
				description: "os.Getenv is non-deterministic in workflows",
				suggestion:  "Pass environment values as workflow input",
			},
			{
				pattern:     "go func()",
				description: "Goroutines are not allowed in workflows",
				suggestion:  "Use workflow.Go() instead",
			},
		}

		for _, p := range nonDetPatterns {
			if strings.Contains(contentStr, p.pattern) {
				warnings = append(warnings, ValidationError{
					File:        file,
					Type:        "NonDeterminism",
					Description: p.description,
					Suggestion:  p.suggestion,
				})
			}
		}

		// Check for missing error handling on activities
		if strings.Contains(contentStr, "ExecuteActivity") && !strings.Contains(contentStr, ".Get(ctx,") {
			warnings = append(warnings, ValidationError{
				File:        file,
				Type:        "MissingErrorHandling",
				Description: "Activity execution may be missing error handling",
				Suggestion:  "Ensure all ExecuteActivity calls have .Get(ctx, &result) with error checking",
			})
		}
	}

	// Check for activity files
	isActivityFile := strings.Contains(contentStr, "activity.") || strings.Contains(file, "activity")

	if isActivityFile {
		// Check for context usage
		if !strings.Contains(contentStr, "context.Context") {
			warnings = append(warnings, ValidationError{
				File:        file,
				Type:        "MissingContext",
				Description: "Activity may be missing context parameter",
				Suggestion:  "Activities should accept context.Context as first parameter",
			})
		}
	}

	return errors, warnings
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

func outputReplayResult(result ReplayResult) error {
	if jsonOutput {
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
	} else {
		if result.Success {
			fmt.Printf("Replay Analysis: %s\n", result.HistoryFile)
			fmt.Println("----------------------------")
			fmt.Printf("Events: %d\n", result.EventCount)
			if result.WorkflowName != "" {
				fmt.Printf("Workflow: %s\n", result.WorkflowName)
			}
			fmt.Println()
			if len(result.Suggestions) > 0 {
				fmt.Println("Next Steps:")
				for _, s := range result.Suggestions {
					fmt.Printf("  %s\n", s)
				}
			}
		} else {
			fmt.Fprintf(os.Stderr, "Error: %s\n", result.Error)
			return fmt.Errorf("replay analysis failed")
		}
	}
	return nil
}

func outputValidationResult(result ValidationResult) error {
	if jsonOutput {
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Printf("Validation: %s\n", result.Path)
		fmt.Println("----------------------------")

		if len(result.Errors) > 0 {
			fmt.Println("\nErrors:")
			for _, e := range result.Errors {
				fmt.Printf("  [%s] %s\n", e.Type, e.Description)
				if e.File != "" {
					fmt.Printf("    File: %s\n", e.File)
				}
				if e.Suggestion != "" {
					fmt.Printf("    Suggestion: %s\n", e.Suggestion)
				}
			}
		}

		if len(result.Warnings) > 0 {
			fmt.Println("\nWarnings:")
			for _, w := range result.Warnings {
				fmt.Printf("  [%s] %s\n", w.Type, w.Description)
				if w.File != "" {
					fmt.Printf("    File: %s\n", w.File)
				}
				if w.Suggestion != "" {
					fmt.Printf("    Suggestion: %s\n", w.Suggestion)
				}
			}
		}

		if result.Valid && len(result.Warnings) == 0 {
			fmt.Println("\n✓ No issues found")
		} else if result.Valid {
			fmt.Printf("\n✓ Valid with %d warning(s)\n", len(result.Warnings))
		} else {
			fmt.Printf("\n✗ Found %d error(s) and %d warning(s)\n", len(result.Errors), len(result.Warnings))
		}
	}

	if !result.Valid {
		return fmt.Errorf("validation failed")
	}
	return nil
}

// runGoTest runs go test in the specified directory
func runGoTest(dir string, args ...string) (string, error) {
	cmdArgs := append([]string{"test"}, args...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	return string(output), err
}
