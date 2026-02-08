package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	jsonOutput bool
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "timelord",
	Short: "Timelord CLI for Temporal.io development and operations",
	Long: `Timelord is a CLI tool for working with Temporal.io applications.

It provides commands for:
  - Scaffolding new projects, workflows, activities, and workers
  - Checking cluster status and configuration
  - Managing namespaces
  - Inspecting and diagnosing workflow executions
  - Running tests`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}
