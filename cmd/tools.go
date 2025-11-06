package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jkeresman01/mcp-client/transport"

	"github.com/spf13/cobra"
)

var (
	toolName string
	toolArgs string
)

var listToolsCmd = &cobra.Command{
	Use:   "list-tools",
	Short: "List all registered MCP tools",
	Long:  `Retrieve and display all tools available on the MCP server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := getTransport()
		if err != nil {
			return err
		}
		defer t.Close()

		req := transport.RPCRequest{
			JSONRPC: "2.0",
			ID:      2,
			Method:  "tools/list",
			Params:  map[string]interface{}{},
		}

		resp, err := t.Send(req)
		if err != nil {
			return err
		}

		out, _ := json.MarshalIndent(resp.Result, "", "  ")
		fmt.Println("Tools:\n", string(out))
		return nil
	},
}

var callToolCmd = &cobra.Command{
	Use:   "call-tool",
	Short: "Call an MCP tool by name with JSON args",
	Long: `Execute a specific tool on the MCP server with provided arguments.
Arguments should be provided as a JSON string.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := getTransport()
		if err != nil {
			return err
		}
		defer t.Close()

		var parsedArgs map[string]interface{}
		if err := json.Unmarshal([]byte(toolArgs), &parsedArgs); err != nil {
			return fmt.Errorf("invalid --args JSON: %v", err)
		}

		req := transport.RPCRequest{
			JSONRPC: "2.0",
			ID:      3,
			Method:  "tools/call",
			Params: map[string]interface{}{
				"name":      toolName,
				"arguments": parsedArgs,
			},
		}

		resp, err := t.Send(req)
		if err != nil {
			return err
		}

		out, _ := json.MarshalIndent(resp.Result, "", "  ")
		fmt.Println("Tool Result:\n", string(out))
		return nil
	},
}

func init() {
	callToolCmd.Flags().StringVar(&toolName, "name", "", "Name of the tool to call (required)")
	callToolCmd.Flags().StringVar(&toolArgs, "args", "{}", "JSON-encoded arguments to pass to the tool")
	callToolCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(listToolsCmd)
	rootCmd.AddCommand(callToolCmd)
}
