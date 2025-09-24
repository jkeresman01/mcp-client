package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	serverURL     string
	transportType string
	commandPath   string
	commandArgs   []string
)

var rootCmd = &cobra.Command{
	Use:   "mcp-client",
	Short: "A general-purpose MCP CLI client",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Transport: %s\n", transportType)
		fmt.Printf("URL: %s\n", serverURL)
		fmt.Printf("Command: %s %v\n", commandPath, commandArgs)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVar(&serverURL, "url", "http://127.0.0.1:8765", "MCP server URL (SSE or StreamableHttp)")
	rootCmd.PersistentFlags().StringVar(&transportType, "transport", "streamable-http", "Transport type: sse | streamable-http | stdio")
	rootCmd.PersistentFlags().StringVar(&commandPath, "command", "", "Command for stdio transport")
	rootCmd.PersistentFlags().StringSliceVar(&commandArgs, "args", []string{}, "Arguments for stdio transport")
}
