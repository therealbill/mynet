package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var nexusCmd = &cobra.Command{
	Use:   "nexus",
	Short: "Manage Nexus endpoints",
	Long:  `Commands for listing, creating, describing, and deleting Temporal Nexus endpoints.`,
}

var nexusListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Nexus endpoints",
	Long:  `List all Nexus endpoints in the Temporal cluster.`,
	RunE:  runNexusList,
}

var nexusCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new Nexus endpoint",
	Long:  `Create a new Nexus endpoint that routes to a target namespace and task queue.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runNexusCreate,
}

var nexusDescribeCmd = &cobra.Command{
	Use:   "describe [name]",
	Short: "Describe a Nexus endpoint",
	Long:  `Show detailed information about a Nexus endpoint.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runNexusDescribe,
}

var nexusDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete a Nexus endpoint",
	Long:  `Delete a Nexus endpoint from the Temporal cluster.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runNexusDelete,
}

var (
	targetNamespace string
	targetTaskQueue string
)

func init() {
	rootCmd.AddCommand(nexusCmd)
	nexusCmd.AddCommand(nexusListCmd)
	nexusCmd.AddCommand(nexusCreateCmd)
	nexusCmd.AddCommand(nexusDescribeCmd)
	nexusCmd.AddCommand(nexusDeleteCmd)

	nexusCreateCmd.Flags().StringVar(&targetNamespace, "target-namespace", "", "Target namespace for the Nexus endpoint")
	nexusCreateCmd.Flags().StringVar(&targetTaskQueue, "target-task-queue", "", "Target task queue for the Nexus endpoint")
	_ = nexusCreateCmd.MarkFlagRequired("target-namespace")
	_ = nexusCreateCmd.MarkFlagRequired("target-task-queue")
}

type NexusEndpointInfo struct {
	Name            string `json:"name"`
	TargetNamespace string `json:"targetNamespace"`
	TargetTaskQueue string `json:"targetTaskQueue"`
}

type NexusEndpointListResult struct {
	Endpoints []NexusEndpointInfo `json:"endpoints"`
	Count     int                 `json:"count"`
}

type NexusEndpointCreateResult struct {
	Name            string `json:"name"`
	TargetNamespace string `json:"targetNamespace"`
	TargetTaskQueue string `json:"targetTaskQueue"`
	Message         string `json:"message"`
}

type NexusEndpointDeleteResult struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func runNexusList(cmd *cobra.Command, args []string) error {
	output, err := runTemporalCLI("operator", "nexus", "endpoint", "list", "--output", "json")
	if err != nil {
		if jsonOutput {
			errResp := map[string]string{"error": strings.TrimSpace(output)}
			data, _ := json.MarshalIndent(errResp, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Fprintf(os.Stderr, "Error listing Nexus endpoints: %s\n", strings.TrimSpace(output))
		}
		return fmt.Errorf("failed to list Nexus endpoints")
	}

	if jsonOutput {
		// Parse and restructure the output
		var rawEndpoints []map[string]any
		if err := json.Unmarshal([]byte(output), &rawEndpoints); err != nil {
			fmt.Println(output)
			return nil
		}

		endpoints := make([]NexusEndpointInfo, 0, len(rawEndpoints))
		for _, ep := range rawEndpoints {
			info := NexusEndpointInfo{}
			if name, ok := ep["name"].(string); ok {
				info.Name = name
			}
			if targetNs, ok := ep["targetNamespace"].(string); ok {
				info.TargetNamespace = targetNs
			}
			if targetTq, ok := ep["targetTaskQueue"].(string); ok {
				info.TargetTaskQueue = targetTq
			}
			endpoints = append(endpoints, info)
		}

		result := NexusEndpointListResult{
			Endpoints: endpoints,
			Count:     len(endpoints),
		}
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
	} else {
		var rawEndpoints []map[string]any
		if err := json.Unmarshal([]byte(output), &rawEndpoints); err != nil {
			fmt.Println(output)
			return nil
		}

		fmt.Println("Nexus Endpoints")
		fmt.Println("---------------")
		for _, ep := range rawEndpoints {
			name := ep["name"]
			targetNs := ep["targetNamespace"]
			targetTq := ep["targetTaskQueue"]
			fmt.Printf("  %v -> %v/%v\n", name, targetNs, targetTq)
		}
		fmt.Printf("\nTotal: %d endpoints\n", len(rawEndpoints))
	}

	return nil
}

func runNexusCreate(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Build command arguments
	createArgs := []string{
		"operator", "nexus", "endpoint", "create",
		"--name", name,
		"--target-namespace", targetNamespace,
		"--target-task-queue", targetTaskQueue,
	}

	output, err := runTemporalCLI(createArgs...)
	if err != nil {
		if jsonOutput {
			errResp := map[string]string{
				"error": strings.TrimSpace(output),
				"name":  name,
			}
			data, _ := json.MarshalIndent(errResp, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Fprintf(os.Stderr, "Error creating Nexus endpoint: %s\n", strings.TrimSpace(output))
		}
		return fmt.Errorf("failed to create Nexus endpoint")
	}

	result := NexusEndpointCreateResult{
		Name:            name,
		TargetNamespace: targetNamespace,
		TargetTaskQueue: targetTaskQueue,
		Message:         fmt.Sprintf("Nexus endpoint '%s' created successfully", name),
	}

	if jsonOutput {
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Printf("Created Nexus endpoint: %s\n", name)
		fmt.Printf("Target Namespace: %s\n", targetNamespace)
		fmt.Printf("Target Task Queue: %s\n", targetTaskQueue)
	}

	return nil
}

func runNexusDescribe(cmd *cobra.Command, args []string) error {
	name := args[0]

	output, err := runTemporalCLI("operator", "nexus", "endpoint", "describe", "--name", name, "--output", "json")
	if err != nil {
		if jsonOutput {
			errResp := map[string]string{
				"error": strings.TrimSpace(output),
				"name":  name,
			}
			data, _ := json.MarshalIndent(errResp, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Fprintf(os.Stderr, "Error describing Nexus endpoint: %s\n", strings.TrimSpace(output))
		}
		return fmt.Errorf("failed to describe Nexus endpoint")
	}

	if jsonOutput {
		fmt.Println(output)
	} else {
		var info map[string]any
		if err := json.Unmarshal([]byte(output), &info); err != nil {
			fmt.Println(output)
			return nil
		}

		fmt.Printf("Nexus Endpoint: %s\n", name)
		fmt.Println("-----------------------")

		if targetNs, ok := info["targetNamespace"]; ok {
			fmt.Printf("Target Namespace: %v\n", targetNs)
		}
		if targetTq, ok := info["targetTaskQueue"]; ok {
			fmt.Printf("Target Task Queue: %v\n", targetTq)
		}
		if id, ok := info["id"]; ok {
			fmt.Printf("ID: %v\n", id)
		}
		if version, ok := info["version"]; ok {
			fmt.Printf("Version: %v\n", version)
		}
		if desc, ok := info["description"]; ok && desc != "" {
			fmt.Printf("Description: %v\n", desc)
		}
	}

	return nil
}

func runNexusDelete(cmd *cobra.Command, args []string) error {
	name := args[0]

	output, err := runTemporalCLI("operator", "nexus", "endpoint", "delete", "--name", name)
	if err != nil {
		if jsonOutput {
			errResp := map[string]string{
				"error": strings.TrimSpace(output),
				"name":  name,
			}
			data, _ := json.MarshalIndent(errResp, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Fprintf(os.Stderr, "Error deleting Nexus endpoint: %s\n", strings.TrimSpace(output))
		}
		return fmt.Errorf("failed to delete Nexus endpoint")
	}

	result := NexusEndpointDeleteResult{
		Name:    name,
		Message: fmt.Sprintf("Nexus endpoint '%s' deleted successfully", name),
	}

	if jsonOutput {
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Printf("Deleted Nexus endpoint: %s\n", name)
	}

	return nil
}
