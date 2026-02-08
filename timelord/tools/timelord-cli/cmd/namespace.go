package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var namespaceCmd = &cobra.Command{
	Use:   "namespace",
	Short: "Manage Temporal namespaces",
	Long:  `Commands for listing, creating, and describing Temporal namespaces.`,
}

var namespaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all namespaces",
	Long:  `List all namespaces in the Temporal cluster.`,
	RunE:  runNamespaceList,
}

var namespaceCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new namespace",
	Long:  `Create a new namespace with optional retention period.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runNamespaceCreate,
}

var namespaceDescribeCmd = &cobra.Command{
	Use:   "describe [name]",
	Short: "Describe a namespace",
	Long:  `Show detailed information about a namespace.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runNamespaceDescribe,
}

var (
	retention   string
	description string
)

func init() {
	rootCmd.AddCommand(namespaceCmd)
	namespaceCmd.AddCommand(namespaceListCmd)
	namespaceCmd.AddCommand(namespaceCreateCmd)
	namespaceCmd.AddCommand(namespaceDescribeCmd)

	namespaceCreateCmd.Flags().StringVar(&retention, "retention", "3d", "Workflow execution retention period")
	namespaceCreateCmd.Flags().StringVar(&description, "description", "", "Namespace description")
}

type NamespaceInfo struct {
	Name        string `json:"name"`
	State       string `json:"state"`
	Retention   string `json:"retention"`
	Description string `json:"description,omitempty"`
}

type NamespaceListResult struct {
	Namespaces []NamespaceInfo `json:"namespaces"`
	Count      int             `json:"count"`
}

type NamespaceCreateResult struct {
	Name      string `json:"name"`
	Retention string `json:"retention"`
	Message   string `json:"message"`
}

func runNamespaceList(cmd *cobra.Command, args []string) error {
	output, err := runTemporalCLI("operator", "namespace", "list", "--output", "json")
	if err != nil {
		if jsonOutput {
			errResp := map[string]string{"error": strings.TrimSpace(output)}
			data, _ := json.MarshalIndent(errResp, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Fprintf(os.Stderr, "Error listing namespaces: %s\n", strings.TrimSpace(output))
		}
		return fmt.Errorf("failed to list namespaces")
	}

	if jsonOutput {
		// Parse and restructure the output
		var rawNamespaces []map[string]interface{}
		if err := json.Unmarshal([]byte(output), &rawNamespaces); err != nil {
			fmt.Println(output)
			return nil
		}

		namespaces := make([]NamespaceInfo, 0, len(rawNamespaces))
		for _, ns := range rawNamespaces {
			info := NamespaceInfo{}
			if name, ok := ns["name"].(string); ok {
				info.Name = name
			}
			if state, ok := ns["state"].(string); ok {
				info.State = state
			}
			namespaces = append(namespaces, info)
		}

		result := NamespaceListResult{
			Namespaces: namespaces,
			Count:      len(namespaces),
		}
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
	} else {
		var rawNamespaces []map[string]interface{}
		if err := json.Unmarshal([]byte(output), &rawNamespaces); err != nil {
			fmt.Println(output)
			return nil
		}

		fmt.Println("Namespaces")
		fmt.Println("----------")
		for _, ns := range rawNamespaces {
			name := ns["name"]
			state := ns["state"]
			fmt.Printf("  %v (%v)\n", name, state)
		}
		fmt.Printf("\nTotal: %d namespaces\n", len(rawNamespaces))
	}

	return nil
}

func runNamespaceCreate(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Build command arguments
	createArgs := []string{
		"operator", "namespace", "create",
		"--namespace", name,
		"--retention", retention,
	}
	if description != "" {
		createArgs = append(createArgs, "--description", description)
	}

	output, err := runTemporalCLI(createArgs...)
	if err != nil {
		if jsonOutput {
			errResp := map[string]string{
				"error":     strings.TrimSpace(output),
				"namespace": name,
			}
			data, _ := json.MarshalIndent(errResp, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Fprintf(os.Stderr, "Error creating namespace: %s\n", strings.TrimSpace(output))
		}
		return fmt.Errorf("failed to create namespace")
	}

	result := NamespaceCreateResult{
		Name:      name,
		Retention: retention,
		Message:   fmt.Sprintf("Namespace '%s' created successfully", name),
	}

	if jsonOutput {
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Printf("Created namespace: %s\n", name)
		fmt.Printf("Retention: %s\n", retention)
		if description != "" {
			fmt.Printf("Description: %s\n", description)
		}
	}

	return nil
}

func runNamespaceDescribe(cmd *cobra.Command, args []string) error {
	name := args[0]

	output, err := runTemporalCLI("operator", "namespace", "describe", "--namespace", name, "--output", "json")
	if err != nil {
		if jsonOutput {
			errResp := map[string]string{
				"error":     strings.TrimSpace(output),
				"namespace": name,
			}
			data, _ := json.MarshalIndent(errResp, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Fprintf(os.Stderr, "Error describing namespace: %s\n", strings.TrimSpace(output))
		}
		return fmt.Errorf("failed to describe namespace")
	}

	if jsonOutput {
		fmt.Println(output)
	} else {
		var info map[string]interface{}
		if err := json.Unmarshal([]byte(output), &info); err != nil {
			fmt.Println(output)
			return nil
		}

		fmt.Printf("Namespace: %s\n", name)
		fmt.Println("-----------------------")

		if state, ok := info["state"]; ok {
			fmt.Printf("State: %v\n", state)
		}

		if nsInfo, ok := info["namespaceInfo"].(map[string]interface{}); ok {
			if desc, ok := nsInfo["description"]; ok && desc != "" {
				fmt.Printf("Description: %v\n", desc)
			}
		}

		if config, ok := info["config"].(map[string]interface{}); ok {
			if retention, ok := config["workflowExecutionRetentionTtl"]; ok {
				fmt.Printf("Retention: %v\n", retention)
			}
			if historyArchival, ok := config["historyArchivalState"]; ok {
				fmt.Printf("History Archival: %v\n", historyArchival)
			}
			if visibilityArchival, ok := config["visibilityArchivalState"]; ok {
				fmt.Printf("Visibility Archival: %v\n", visibilityArchival)
			}
		}

		if replicationConfig, ok := info["replicationConfig"].(map[string]interface{}); ok {
			if activeCluster, ok := replicationConfig["activeClusterName"]; ok {
				fmt.Printf("Active Cluster: %v\n", activeCluster)
			}
		}
	}

	return nil
}

// Helper to run temporal CLI
func runTemporalCLINs(args ...string) (string, error) {
	addr := getTemporalAddress()
	fullArgs := append([]string{"--address", addr}, args...)
	execCmd := exec.Command("temporal", fullArgs...)
	output, err := execCmd.CombinedOutput()
	return string(output), err
}
