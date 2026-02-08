package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var workflowCmd = &cobra.Command{
	Use:   "workflow",
	Short: "Manage and inspect workflows",
	Long:  `Commands for listing, describing, and diagnosing workflow executions.`,
}

var workflowListCmd = &cobra.Command{
	Use:   "list",
	Short: "List workflow executions",
	Long:  `List workflow executions with optional filtering.`,
	RunE:  runWorkflowList,
}

var workflowDescribeCmd = &cobra.Command{
	Use:   "describe [workflow-id]",
	Short: "Describe a workflow execution",
	Long:  `Show detailed information about a workflow execution.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runWorkflowDescribe,
}

var workflowHistoryCmd = &cobra.Command{
	Use:   "history [workflow-id]",
	Short: "Show workflow event history",
	Long:  `Display the event history for a workflow execution.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runWorkflowHistory,
}

var workflowDiagnoseCmd = &cobra.Command{
	Use:   "diagnose [workflow-id]",
	Short: "Diagnose workflow issues",
	Long:  `Analyze a workflow execution to identify potential issues.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runWorkflowDiagnose,
}

var (
	workflowNamespace string
	workflowQuery     string
	workflowLimit     int
	workflowRunID     string
)

func init() {
	rootCmd.AddCommand(workflowCmd)
	workflowCmd.AddCommand(workflowListCmd)
	workflowCmd.AddCommand(workflowDescribeCmd)
	workflowCmd.AddCommand(workflowHistoryCmd)
	workflowCmd.AddCommand(workflowDiagnoseCmd)

	workflowCmd.PersistentFlags().StringVarP(&workflowNamespace, "namespace", "n", "default", "Temporal namespace")

	workflowListCmd.Flags().StringVarP(&workflowQuery, "query", "q", "", "List filter query")
	workflowListCmd.Flags().IntVarP(&workflowLimit, "limit", "l", 10, "Maximum workflows to return")

	workflowDescribeCmd.Flags().StringVar(&workflowRunID, "run-id", "", "Run ID (optional)")
	workflowHistoryCmd.Flags().StringVar(&workflowRunID, "run-id", "", "Run ID (optional)")
	workflowDiagnoseCmd.Flags().StringVar(&workflowRunID, "run-id", "", "Run ID (optional)")
}

type WorkflowListResult struct {
	Workflows []WorkflowSummary `json:"workflows"`
	Count     int               `json:"count"`
}

type WorkflowSummary struct {
	WorkflowID   string `json:"workflowId"`
	RunID        string `json:"runId"`
	Type         string `json:"type"`
	Status       string `json:"status"`
	StartTime    string `json:"startTime"`
	CloseTime    string `json:"closeTime,omitempty"`
}

type WorkflowDescription struct {
	WorkflowID      string                 `json:"workflowId"`
	RunID           string                 `json:"runId"`
	Type            string                 `json:"type"`
	Status          string                 `json:"status"`
	StartTime       string                 `json:"startTime"`
	CloseTime       string                 `json:"closeTime,omitempty"`
	HistoryLength   int                    `json:"historyLength"`
	TaskQueue       string                 `json:"taskQueue"`
	Memo            map[string]interface{} `json:"memo,omitempty"`
	SearchAttributes map[string]interface{} `json:"searchAttributes,omitempty"`
}

type DiagnosisResult struct {
	WorkflowID   string           `json:"workflowId"`
	Status       string           `json:"status"`
	Issues       []DiagnosisIssue `json:"issues"`
	Suggestions  []string         `json:"suggestions"`
	LastEvent    string           `json:"lastEvent"`
	PendingItems []string         `json:"pendingItems,omitempty"`
}

type DiagnosisIssue struct {
	Severity    string `json:"severity"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

func runWorkflowList(cmd *cobra.Command, args []string) error {
	listArgs := []string{
		"workflow", "list",
		"--namespace", workflowNamespace,
		"--limit", fmt.Sprintf("%d", workflowLimit),
		"--output", "json",
	}

	if workflowQuery != "" {
		listArgs = append(listArgs, "--query", workflowQuery)
	}

	output, err := runTemporalCLI(listArgs...)
	if err != nil {
		if jsonOutput {
			errResp := map[string]string{"error": strings.TrimSpace(output)}
			data, _ := json.MarshalIndent(errResp, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Fprintf(os.Stderr, "Error listing workflows: %s\n", strings.TrimSpace(output))
		}
		return fmt.Errorf("failed to list workflows")
	}

	if jsonOutput {
		fmt.Println(output)
	} else {
		var workflows []map[string]interface{}
		if err := json.Unmarshal([]byte(output), &workflows); err != nil {
			fmt.Println(output)
			return nil
		}

		fmt.Printf("Workflows (namespace: %s)\n", workflowNamespace)
		fmt.Println("----------------------------")

		for _, wf := range workflows {
			workflowID := wf["execution"].(map[string]interface{})["workflowId"]
			wfType := ""
			if t, ok := wf["type"].(map[string]interface{}); ok {
				wfType = t["name"].(string)
			}
			status := wf["status"]

			fmt.Printf("  %v [%v] - %v\n", workflowID, status, wfType)
		}

		fmt.Printf("\nShowing %d workflows\n", len(workflows))
	}

	return nil
}

func runWorkflowDescribe(cmd *cobra.Command, args []string) error {
	workflowID := args[0]

	describeArgs := []string{
		"workflow", "describe",
		"--namespace", workflowNamespace,
		"--workflow-id", workflowID,
		"--output", "json",
	}

	if workflowRunID != "" {
		describeArgs = append(describeArgs, "--run-id", workflowRunID)
	}

	output, err := runTemporalCLI(describeArgs...)
	if err != nil {
		if jsonOutput {
			errResp := map[string]string{
				"error":      strings.TrimSpace(output),
				"workflowId": workflowID,
			}
			data, _ := json.MarshalIndent(errResp, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Fprintf(os.Stderr, "Error describing workflow: %s\n", strings.TrimSpace(output))
		}
		return fmt.Errorf("failed to describe workflow")
	}

	if jsonOutput {
		fmt.Println(output)
	} else {
		var desc map[string]interface{}
		if err := json.Unmarshal([]byte(output), &desc); err != nil {
			fmt.Println(output)
			return nil
		}

		fmt.Printf("Workflow: %s\n", workflowID)
		fmt.Println("----------------------------")

		if exec, ok := desc["workflowExecutionInfo"].(map[string]interface{}); ok {
			if status, ok := exec["status"]; ok {
				fmt.Printf("Status: %v\n", status)
			}
			if t, ok := exec["type"].(map[string]interface{}); ok {
				fmt.Printf("Type: %v\n", t["name"])
			}
			if start, ok := exec["startTime"]; ok {
				fmt.Printf("Start Time: %v\n", start)
			}
			if close, ok := exec["closeTime"]; ok && close != nil {
				fmt.Printf("Close Time: %v\n", close)
			}
			if tq, ok := exec["taskQueue"]; ok {
				fmt.Printf("Task Queue: %v\n", tq)
			}
			if hl, ok := exec["historyLength"]; ok {
				fmt.Printf("History Length: %v events\n", hl)
			}
		}
	}

	return nil
}

func runWorkflowHistory(cmd *cobra.Command, args []string) error {
	workflowID := args[0]

	historyArgs := []string{
		"workflow", "show",
		"--namespace", workflowNamespace,
		"--workflow-id", workflowID,
	}

	if workflowRunID != "" {
		historyArgs = append(historyArgs, "--run-id", workflowRunID)
	}

	if jsonOutput {
		historyArgs = append(historyArgs, "--output", "json")
	}

	output, err := runTemporalCLI(historyArgs...)
	if err != nil {
		if jsonOutput {
			errResp := map[string]string{
				"error":      strings.TrimSpace(output),
				"workflowId": workflowID,
			}
			data, _ := json.MarshalIndent(errResp, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Fprintf(os.Stderr, "Error getting workflow history: %s\n", strings.TrimSpace(output))
		}
		return fmt.Errorf("failed to get workflow history")
	}

	fmt.Println(output)
	return nil
}

func runWorkflowDiagnose(cmd *cobra.Command, args []string) error {
	workflowID := args[0]

	// Get workflow description
	describeArgs := []string{
		"workflow", "describe",
		"--namespace", workflowNamespace,
		"--workflow-id", workflowID,
		"--output", "json",
	}

	if workflowRunID != "" {
		describeArgs = append(describeArgs, "--run-id", workflowRunID)
	}

	descOutput, err := runTemporalCLI(describeArgs...)
	if err != nil {
		if jsonOutput {
			errResp := map[string]string{
				"error":      strings.TrimSpace(descOutput),
				"workflowId": workflowID,
			}
			data, _ := json.MarshalIndent(errResp, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Fprintf(os.Stderr, "Error getting workflow info: %s\n", strings.TrimSpace(descOutput))
		}
		return fmt.Errorf("failed to get workflow info")
	}

	// Get workflow history
	historyArgs := []string{
		"workflow", "show",
		"--namespace", workflowNamespace,
		"--workflow-id", workflowID,
		"--output", "json",
	}

	if workflowRunID != "" {
		historyArgs = append(historyArgs, "--run-id", workflowRunID)
	}

	historyOutput, err := runTemporalCLI(historyArgs...)
	if err != nil {
		historyOutput = ""
	}

	// Analyze
	diagnosis := analyzeWorkflow(workflowID, descOutput, historyOutput)

	if jsonOutput {
		data, _ := json.MarshalIndent(diagnosis, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Printf("Diagnosis for Workflow: %s\n", workflowID)
		fmt.Println("================================")
		fmt.Printf("Status: %s\n", diagnosis.Status)
		fmt.Printf("Last Event: %s\n", diagnosis.LastEvent)

		if len(diagnosis.PendingItems) > 0 {
			fmt.Println("\nPending Items:")
			for _, item := range diagnosis.PendingItems {
				fmt.Printf("  - %s\n", item)
			}
		}

		if len(diagnosis.Issues) > 0 {
			fmt.Println("\nIssues Found:")
			for _, issue := range diagnosis.Issues {
				fmt.Printf("  [%s] %s: %s\n", issue.Severity, issue.Type, issue.Description)
			}
		} else {
			fmt.Println("\nNo issues detected.")
		}

		if len(diagnosis.Suggestions) > 0 {
			fmt.Println("\nSuggestions:")
			for _, suggestion := range diagnosis.Suggestions {
				fmt.Printf("  - %s\n", suggestion)
			}
		}
	}

	return nil
}

func analyzeWorkflow(workflowID, descOutput, historyOutput string) DiagnosisResult {
	result := DiagnosisResult{
		WorkflowID:  workflowID,
		Status:      "Unknown",
		Issues:      []DiagnosisIssue{},
		Suggestions: []string{},
	}

	// Parse description
	var desc map[string]interface{}
	if err := json.Unmarshal([]byte(descOutput), &desc); err == nil {
		if exec, ok := desc["workflowExecutionInfo"].(map[string]interface{}); ok {
			if status, ok := exec["status"].(string); ok {
				result.Status = status
			}
		}
	}

	// Parse history for analysis
	var history map[string]interface{}
	if err := json.Unmarshal([]byte(historyOutput), &history); err == nil {
		events, ok := history["events"].([]interface{})
		if ok && len(events) > 0 {
			// Get last event
			lastEvent := events[len(events)-1].(map[string]interface{})
			if eventType, ok := lastEvent["eventType"].(string); ok {
				result.LastEvent = eventType
			}

			// Analyze events for issues
			pendingActivities := 0
			failedActivities := 0
			timedOutActivities := 0

			for _, e := range events {
				event := e.(map[string]interface{})
				eventType, _ := event["eventType"].(string)

				switch eventType {
				case "ActivityTaskScheduled":
					pendingActivities++
				case "ActivityTaskCompleted":
					pendingActivities--
				case "ActivityTaskFailed":
					pendingActivities--
					failedActivities++
				case "ActivityTaskTimedOut":
					pendingActivities--
					timedOutActivities++
				}
			}

			// Check for pending activities
			if pendingActivities > 0 {
				result.PendingItems = append(result.PendingItems,
					fmt.Sprintf("%d pending activity(ies)", pendingActivities))
			}

			// Check for failures
			if failedActivities > 0 {
				result.Issues = append(result.Issues, DiagnosisIssue{
					Severity:    "WARNING",
					Type:        "ActivityFailure",
					Description: fmt.Sprintf("%d activity(ies) failed", failedActivities),
				})
				result.Suggestions = append(result.Suggestions,
					"Check activity error messages in history",
					"Review retry policy configuration")
			}

			// Check for timeouts
			if timedOutActivities > 0 {
				result.Issues = append(result.Issues, DiagnosisIssue{
					Severity:    "WARNING",
					Type:        "ActivityTimeout",
					Description: fmt.Sprintf("%d activity(ies) timed out", timedOutActivities),
				})
				result.Suggestions = append(result.Suggestions,
					"Consider increasing activity timeout",
					"Add heartbeats for long-running activities",
					"Check worker connectivity")
			}

			// Check workflow status
			if result.Status == "Running" && pendingActivities > 0 {
				result.Suggestions = append(result.Suggestions,
					"Workflow is waiting for activities to complete",
					"Verify workers are running on the correct task queue")
			}

			if result.Status == "Failed" {
				result.Issues = append(result.Issues, DiagnosisIssue{
					Severity:    "ERROR",
					Type:        "WorkflowFailed",
					Description: "Workflow execution failed",
				})
				result.Suggestions = append(result.Suggestions,
					"Check workflow error in the last WorkflowTaskFailed event",
					"Review activity error handling")
			}
		}
	}

	// Default suggestions for running workflows
	if result.Status == "Running" && len(result.Issues) == 0 {
		result.Suggestions = append(result.Suggestions,
			"Workflow appears healthy and running normally")
	}

	return result
}
