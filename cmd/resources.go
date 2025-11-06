package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jkeresman01/mcp-client/transport"

	"github.com/spf13/cobra"
)

var resourceID string

var listResourcesCmd = &cobra.Command{
	Use:   "list-resources",
	Short: "List all available MCP resources",
	Long:  `Retrieve and display all resources available on the MCP server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := getTransport()
		if err != nil {
			return err
		}
		defer t.Close()

		req := transport.RPCRequest{
			JSONRPC: "2.0",
			ID:      10,
			Method:  "resources/list",
			Params:  map[string]interface{}{},
		}

		if debugMode {
			fmt.Println("Debug: Sending request:", req.Method)
		}

		resp, err := t.Send(req)
		if err != nil {
			return transport.WrapError("list-resources", err)
		}

		if resp.Error != nil {
			return &transport.MCPError{
				Operation: "list-resources",
				Err:       fmt.Errorf("server returned error: %v", resp.Error),
				Hints: []string{
					"The server may not have initialized properly",
					"Try running 'mcp-client init' first",
					"Check server logs for more details",
				},
			}
		}

		output, _ := json.MarshalIndent(resp.Result, "", "  ")
		fmt.Println("Resources:\n", string(output))
		return nil
	},
}

var getResourceCmd = &cobra.Command{
	Use:   "get-resource",
	Short: "Get content of a specific MCP resource by ID",
	Long: `Retrieve the content of a specific resource from the MCP server using its ID.

Examples:
  mcp-client get-resource --id file:///path/to/file
  mcp-client get-resource --id resource://my-resource --server prod`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if resourceID == "" {
			return &transport.MCPError{
				Operation: "get-resource",
				Err:       fmt.Errorf("resource ID is required"),
				Hints: []string{
					"Specify resource URI: --id <resource-uri>",
					"List available resources: mcp-client list-resources",
				},
			}
		}

		t, err := getTransport()
		if err != nil {
			return err
		}
		defer t.Close()

		req := transport.RPCRequest{
			JSONRPC: "2.0",
			ID:      11,
			Method:  "resources/read",
			Params: map[string]interface{}{
				"uri": resourceID,
			},
		}

		if debugMode {
			fmt.Printf("Debug: Reading resource: %s\n", resourceID)
		}

		resp, err := t.Send(req)
		if err != nil {
			return transport.WrapError("get-resource", err)
		}

		if resp.Error != nil {
			errMap, ok := resp.Error.(map[string]interface{})
			if ok {
				if code, exists := errMap["code"]; exists && code == -32602 {
					return &transport.MCPError{
						Operation: "get-resource",
						Err:       fmt.Errorf("resource '%s' not found", resourceID),
						Hints: []string{
							"List available resources: mcp-client list-resources",
							"Check if the resource URI is correct",
							"Verify you have permission to access this resource",
						},
					}
				}
			}
			return &transport.MCPError{
				Operation: "get-resource",
				Err:       fmt.Errorf("server error: %v", resp.Error),
				Hints: []string{
					"Verify the resource URI is correct",
					"Check server logs for more details",
				},
			}
		}

		output, _ := json.MarshalIndent(resp.Result, "", "  ")
		fmt.Println("Resource content:\n", string(output))
		return nil
	},
}

func init() {
	getResourceCmd.Flags().StringVar(&resourceID, "id", "", "ID/URI of the resource to fetch")
	getResourceCmd.MarkFlagRequired("id")

	rootCmd.AddCommand(listResourcesCmd)
	rootCmd.AddCommand(getResourceCmd)
}
