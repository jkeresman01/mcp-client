package cmd

import (
	"fmt"
	"github.com/jkeresman01/mcp-client/config"
	"github.com/spf13/cobra"
)

var (
	serverURL     string
	transportType string
	commandPath   string
	commandArgs   []string
	serverName    string
	debugMode     bool
)

var rootCmd = &cobra.Command{
	Use:   "mcp-client",
	Short: "A general-purpose MCP CLI client",
	Long: `MCP Client - A command-line interface for testing Model Context Protocol servers.
Supports stdio, SSE, and StreamableHTTP transports.

Configuration:
  Use 'mcp-client config init' to create a configuration file.
  Then use 'mcp-client config add' to add servers.
  
Examples:
  mcp-client config add local --url http://localhost:8765 --transport streamable-http
  mcp-client list-tools --server local
  mcp-client interactive --server local`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Don't print connection info for config commands or root help
		if cmd.Parent() != nil && cmd.Parent().Name() == "config" {
			return
		}
		if cmd.Name() == "mcp-client" || cmd.Name() == "config" {
			return
		}

		// Load configuration if --server is specified
		if serverName != "" {
			cfg, err := config.Load(configFile)
			if err != nil {
				fmt.Printf("Warning: Could not load config: %v\n", err)
				return
			}

			server, err := cfg.GetServer(serverName)
			if err != nil {
				fmt.Printf("Warning: %v\n", err)
				return
			}

			// Override flags with config values
			transportType = server.Transport
			serverURL = server.URL
			commandPath = server.Command
			commandArgs = server.Args

			if debugMode {
				fmt.Printf("Debug: Using server '%s' from config\n", serverName)
			}
		}

		// Print connection info
		if debugMode {
			fmt.Println("--------------------------------------------------------")
		}
		fmt.Printf("Transport: %s\n", transportType)
		if transportType != "stdio" {
			fmt.Printf("URL: %s\n", serverURL)
		} else {
			fmt.Printf("Command: %s %v\n", commandPath, commandArgs)
		}
		if debugMode {
			fmt.Println("--------------------------------------------------------")
		}
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
	rootCmd.PersistentFlags().StringVar(&serverName, "server", "", "Use a named server from config file")
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "Config file path (default: ~/.mcp-config.json)")
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "Enable debug mode with verbose output")
}
