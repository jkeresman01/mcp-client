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

		if debugMode {
			fmt.Println("Debug: Sending request:", req.Method)
		}

		resp, err := t.Send(req)
		if err != nil {
			return transport.WrapError("list-tools", err)
		}

		if resp.Error != nil {
			return &transport.MCPError{
				Operation: "list-tools",
				Err:       fmt.Errorf("server returned error: %v", resp.Error),
				Hints: []string{
					"The server may not have initialized properly",
					"Try running 'mcp-client init' first",
					"Check server logs for more details",
				},
			}
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
Arguments should be provided as a JSON string.

Examples:
  mcp-client call-tool --name calculator --args '{"op":"add","a":5,"b":3}'
  mcp-client call-tool --name search --args '{"query":"golang"}' --server prod`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if toolName == "" {
			return &transport.MCPError{
				Operation: "call-tool",
				Err:       fmt.Errorf("tool name is required"),
				Hints: []string{
					"Specify tool name: --name <tool-name>",
					"List available tools: mcp-client list-tools",
				},
			}
		}

		t, err := getTransport()
		if err != nil {
			return err
		}
		defer t.Close()

		var parsedArgs map[string]interface{}
		if err := json.Unmarshal([]byte(toolArgs), &parsedArgs); err != nil {
			return transport.NewInvalidArgumentsError(err.Error())
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

		if debugMode {
			fmt.Printf("Debug: Calling tool '%s' with args: %s\n", toolName, toolArgs)
		}

		resp, err := t.Send(req)
		if err != nil {
			return transport.WrapError("call-tool", err)
		}

		if resp.Error != nil {
			errMap, ok := resp.Error.(map[string]interface{})
			if ok {
				if code, exists := errMap["code"]; exists && code == -32601 {
					return transport.NewToolNotFoundError(toolName)
				}
			}
			return &transport.MCPError{
				Operation: "call-tool",
				Err:       fmt.Errorf("server error: %v", resp.Error),
				Hints: []string{
					"Verify the tool name is correct (case-sensitive)",
					"Check that all required arguments are provided",
					"List available tools: mcp-client list-tools",
				},
			}
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
