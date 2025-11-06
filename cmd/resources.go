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

		resp, err := t.Send(req)
		if err != nil {
			return err
		}

		output, _ := json.MarshalIndent(resp.Result, "", "  ")
		fmt.Println("Resources:\n", string(output))
		return nil
	},
}

var getResourceCmd = &cobra.Command{
	Use:   "get-resource",
	Short: "Get content of a specific MCP resource by ID",
	Long:  `Retrieve the content of a specific resource from the MCP server using its ID.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if resourceID == "" {
			return fmt.Errorf("--id is required")
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

		resp, err := t.Send(req)
		if err != nil {
			return err
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
