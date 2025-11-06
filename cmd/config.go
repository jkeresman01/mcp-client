package cmd

import (
	"fmt"
	"os"

	"github.com/jkeresman01/mcp-client/config"
	"github.com/spf13/cobra"
)

var (
	configFile string
	serverName string
	setDefault bool
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage MCP client configuration",
	Long:  `Manage server configurations, set defaults, and view current settings.`,
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured servers",
	Long:  `Display all servers configured in the MCP client configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(configFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %v", err)
		}

		if len(cfg.Servers) == 0 {
			fmt.Println("No servers configured.")
			fmt.Println("\nAdd a server with:")
			fmt.Println("  mcp-client config add <name> --url <url> --transport <type>")
			return nil
		}

		fmt.Println("Configured servers:")
		fmt.Println()

		for name, server := range cfg.Servers {
			marker := "  "
			if cfg.DefaultServer == name {
				marker = "* "
			}

			fmt.Printf("%s%s:\n", marker, name)
			fmt.Printf("    Transport: %s\n", server.Transport)

			if server.URL != "" {
				fmt.Printf("    URL:       %s\n", server.URL)
			}
			if server.Command != "" {
				fmt.Printf("    Command:   %s\n", server.Command)
				if len(server.Args) > 0 {
					fmt.Printf("    Args:      %v\n", server.Args)
				}
			}
			fmt.Println()
		}

		if cfg.DefaultServer != "" {
			fmt.Printf("Default server: %s (marked with *)\n", cfg.DefaultServer)
		}

		return nil
	},
}

var configAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add or update a server configuration",
	Long:  `Add a new server or update an existing server configuration.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cfg, err := config.Load(configFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %v", err)
		}

		// Validate required fields
		if transportType == "" {
			return fmt.Errorf("--transport is required")
		}

		if transportType == "stdio" {
			if commandPath == "" {
				return fmt.Errorf("--command is required for stdio transport")
			}
		} else {
			if serverURL == "" {
				return fmt.Errorf("--url is required for %s transport", transportType)
			}
		}

		server := config.ServerConfig{
			Transport: transportType,
			URL:       serverURL,
			Command:   commandPath,
			Args:      commandArgs,
		}

		cfg.AddServer(name, server)

		if setDefault {
			cfg.DefaultServer = name
		}

		if err := cfg.Save(configFile); err != nil {
			return fmt.Errorf("failed to save config: %v", err)
		}

		fmt.Printf("Server '%s' configured successfully!\n", name)
		if setDefault {
			fmt.Printf("Set as default server\n")
		}

		return nil
	},
}

var configRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a server configuration",
	Long:  `Remove a server from the configuration file.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cfg, err := config.Load(configFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %v", err)
		}

		if _, exists := cfg.Servers[name]; !exists {
			return fmt.Errorf("server '%s' not found in configuration", name)
		}

		delete(cfg.Servers, name)

		if cfg.DefaultServer == name {
			cfg.DefaultServer = ""
		}

		if err := cfg.Save(configFile); err != nil {
			return fmt.Errorf("failed to save config: %v", err)
		}

		fmt.Printf("Server '%s' removed successfully!\n", name)

		return nil
	},
}

var configSetDefaultCmd = &cobra.Command{
	Use:   "set-default <name>",
	Short: "Set the default server",
	Long:  `Set which server configuration should be used by default.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		cfg, err := config.Load(configFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %v", err)
		}

		if _, exists := cfg.Servers[name]; !exists {
			return fmt.Errorf("server '%s' not found in configuration", name)
		}

		cfg.DefaultServer = name

		if err := cfg.Save(configFile); err != nil {
			return fmt.Errorf("failed to save config: %v", err)
		}

		fmt.Printf("Default server set to '%s'\n", name)

		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show [name]",
	Short: "Show configuration details",
	Long:  `Show details of a specific server or the current configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load(configFile)
		if err != nil {
			return fmt.Errorf("failed to load config: %v", err)
		}

		name := ""
		if len(args) > 0 {
			name = args[0]
		}

		if name != "" {
			server, err := cfg.GetServer(name)
			if err != nil {
				return err
			}

			fmt.Printf("Server: %s\n", name)
			fmt.Printf("  Transport: %s\n", server.Transport)
			if server.URL != "" {
				fmt.Printf("  URL:       %s\n", server.URL)
			}
			if server.Command != "" {
				fmt.Printf("  Command:   %s\n", server.Command)
				if len(server.Args) > 0 {
					fmt.Printf("  Args:      %v\n", server.Args)
				}
			}
		} else {
			// Show current effective configuration
			fmt.Println("Current Configuration:")
			fmt.Printf("  Config file: %s\n", getConfigFilePath())
			fmt.Printf("  Default server: %s\n", cfg.DefaultServer)
			fmt.Printf("  Total servers: %d\n", len(cfg.Servers))
		}

		return nil
	},
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new configuration file",
	Long:  `Create a new configuration file with example servers.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		path := getConfigFilePath()

		// Check if config already exists
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("Configuration file already exists at: %s\n", path)
			fmt.Println("Use 'mcp-client config list' to view current configuration")
			return nil
		}

		cfg := &config.Config{
			DefaultServer: "local",
			Servers: map[string]config.ServerConfig{
				"local": {
					Transport: "streamable-http",
					URL:       "http://localhost:8765",
				},
			},
		}

		if err := cfg.Save(path); err != nil {
			return fmt.Errorf("failed to create config: %v", err)
		}

		fmt.Printf("Configuration file created at: %s\n", path)
		fmt.Println("\nExample server 'local' has been configured.")
		fmt.Println("Edit the file or use 'mcp-client config add' to add more servers.")

		return nil
	},
}

func init() {
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configAddCmd)
	configCmd.AddCommand(configRemoveCmd)
	configCmd.AddCommand(configSetDefaultCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configInitCmd)

	// Flags for add command
	configAddCmd.Flags().BoolVar(&setDefault, "default", false, "Set as default server")

	rootCmd.AddCommand(configCmd)
}

func getConfigFilePath() string {
	if configFile != "" {
		return configFile
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".mcp-config.json"
	}
	return home + "/.mcp-config.json"
}
