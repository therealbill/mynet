package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Manage Temporal cluster",
	Long:  `Commands for checking cluster status, configuration, and metrics.`,
}

var clusterStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check cluster health status",
	Long:  `Check the health and availability of the Temporal cluster.`,
	RunE:  runClusterStatus,
}

var clusterInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show cluster configuration",
	Long:  `Display cluster configuration details including version and capabilities.`,
	RunE:  runClusterInfo,
}

var clusterMetricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Show key cluster metrics",
	Long:  `Display summary of key cluster metrics.`,
	RunE:  runClusterMetrics,
}

var (
	temporalAddress string
)

func init() {
	rootCmd.AddCommand(clusterCmd)
	clusterCmd.AddCommand(clusterStatusCmd)
	clusterCmd.AddCommand(clusterInfoCmd)
	clusterCmd.AddCommand(clusterMetricsCmd)

	clusterCmd.PersistentFlags().StringVar(&temporalAddress, "address", "", "Temporal server address (default: localhost:7233)")
}

type ClusterStatus struct {
	Healthy     bool   `json:"healthy"`
	Status      string `json:"status"`
	Address     string `json:"address"`
	Message     string `json:"message,omitempty"`
	Error       string `json:"error,omitempty"`
}

type ClusterInfo struct {
	ServerVersion     string            `json:"serverVersion"`
	ClusterID         string            `json:"clusterId"`
	ClusterName       string            `json:"clusterName"`
	HistoryShardCount int               `json:"historyShardCount"`
	PersistenceStore  string            `json:"persistenceStore"`
	VisibilityStore   string            `json:"visibilityStore"`
	Features          map[string]bool   `json:"features"`
}

type ClusterMetrics struct {
	ActiveNamespaces  int     `json:"activeNamespaces"`
	WorkflowsRunning  int     `json:"workflowsRunning"`
	TaskQueuesActive  int     `json:"taskQueuesActive"`
	FrontendLatencyP99 float64 `json:"frontendLatencyP99Ms"`
	HistoryLatencyP99  float64 `json:"historyLatencyP99Ms"`
}

func getTemporalAddress() string {
	if temporalAddress != "" {
		return temporalAddress
	}
	if addr := os.Getenv("TEMPORAL_ADDRESS"); addr != "" {
		return addr
	}
	return "localhost:7233"
}

func runTemporalCLI(args ...string) (string, error) {
	addr := getTemporalAddress()
	fullArgs := append([]string{"--address", addr}, args...)
	cmd := exec.Command("temporal", fullArgs...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func runClusterStatus(cmd *cobra.Command, args []string) error {
	addr := getTemporalAddress()
	status := ClusterStatus{
		Address: addr,
	}

	// Try to get cluster health
	output, err := runTemporalCLI("operator", "cluster", "health")
	if err != nil {
		status.Healthy = false
		status.Status = "UNAVAILABLE"
		status.Error = strings.TrimSpace(output)
		if status.Error == "" {
			status.Error = err.Error()
		}
	} else {
		// Parse output for health status
		if strings.Contains(output, "SERVING") {
			status.Healthy = true
			status.Status = "SERVING"
			status.Message = "Cluster is healthy and accepting requests"
		} else {
			status.Healthy = false
			status.Status = "NOT_SERVING"
			status.Message = strings.TrimSpace(output)
		}
	}

	if jsonOutput {
		data, _ := json.MarshalIndent(status, "", "  ")
		fmt.Println(string(data))
	} else {
		if status.Healthy {
			fmt.Printf("Cluster Status: %s\n", status.Status)
			fmt.Printf("Address: %s\n", status.Address)
			fmt.Printf("Message: %s\n", status.Message)
		} else {
			fmt.Printf("Cluster Status: %s\n", status.Status)
			fmt.Printf("Address: %s\n", status.Address)
			if status.Error != "" {
				fmt.Printf("Error: %s\n", status.Error)
			}
		}
	}

	if !status.Healthy {
		os.Exit(1)
	}
	return nil
}

func runClusterInfo(cmd *cobra.Command, args []string) error {
	output, err := runTemporalCLI("operator", "cluster", "describe", "--output", "json")
	if err != nil {
		if jsonOutput {
			errResp := map[string]string{"error": strings.TrimSpace(output)}
			data, _ := json.MarshalIndent(errResp, "", "  ")
			fmt.Println(string(data))
		} else {
			fmt.Fprintf(os.Stderr, "Error getting cluster info: %s\n", strings.TrimSpace(output))
		}
		return fmt.Errorf("failed to get cluster info")
	}

	if jsonOutput {
		// Pass through JSON output
		fmt.Println(output)
	} else {
		// Parse and display formatted output
		var info map[string]interface{}
		if err := json.Unmarshal([]byte(output), &info); err != nil {
			fmt.Println(output)
			return nil
		}

		fmt.Println("Cluster Information")
		fmt.Println("-------------------")
		if v, ok := info["serverVersion"]; ok {
			fmt.Printf("Server Version: %v\n", v)
		}
		if v, ok := info["clusterId"]; ok {
			fmt.Printf("Cluster ID: %v\n", v)
		}
		if v, ok := info["clusterName"]; ok {
			fmt.Printf("Cluster Name: %v\n", v)
		}
		if v, ok := info["historyShardCount"]; ok {
			fmt.Printf("History Shards: %v\n", v)
		}
		if v, ok := info["persistenceStore"]; ok {
			fmt.Printf("Persistence Store: %v\n", v)
		}
		if v, ok := info["visibilityStore"]; ok {
			fmt.Printf("Visibility Store: %v\n", v)
		}
	}

	return nil
}

func runClusterMetrics(cmd *cobra.Command, args []string) error {
	// Get namespace count
	nsOutput, err := runTemporalCLI("operator", "namespace", "list", "--output", "json")
	namespaceCount := 0
	if err == nil {
		var namespaces []interface{}
		if json.Unmarshal([]byte(nsOutput), &namespaces) == nil {
			namespaceCount = len(namespaces)
		}
	}

	metrics := ClusterMetrics{
		ActiveNamespaces: namespaceCount,
		// Note: Other metrics would require Prometheus queries
		// This is a simplified version
	}

	if jsonOutput {
		data, _ := json.MarshalIndent(metrics, "", "  ")
		fmt.Println(string(data))
	} else {
		fmt.Println("Cluster Metrics Summary")
		fmt.Println("-----------------------")
		fmt.Printf("Active Namespaces: %d\n", metrics.ActiveNamespaces)
		fmt.Println("\nNote: For detailed metrics, query Prometheus directly or use Grafana dashboards.")
		fmt.Println("Common metric queries:")
		fmt.Println("  - Workflow rate: sum(rate(temporal_workflow_started_total[5m]))")
		fmt.Println("  - Task latency: histogram_quantile(0.99, rate(temporal_schedule_to_start_latency_bucket[5m]))")
	}

	return nil
}
