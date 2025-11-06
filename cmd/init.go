package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jkeresman01/mcp-client/transport"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Send initialize request to MCP server",
	Long: `Initialize the connection with an MCP server.
This sends the initialize request with protocol version and client info.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := getTransport()
		if err != nil {
			return err
		}
		defer t.Close()

		req := transport.RPCRequest{
			JSONRPC: "2.0",
			ID:      1,
			Method:  "initialize",
			Params: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"clientInfo": map[string]string{
					"name":    "mcp-client",
					"version": "0.2.0",
				},
				"capabilities": map[string]interface{}{},
			},
		}

		if debugMode {
			reqJSON, _ := json.MarshalIndent(req, "", "  ")
			fmt.Println("Debug: Request:")
			fmt.Println(string(reqJSON))
		}

		resp, err := t.Send(req)
		if err != nil {
			return transport.WrapError("initialize", err)
		}

		if resp.Error != nil {
			return &transport.MCPError{
				Operation: "initialize",
				Err:       fmt.Errorf("server returned error: %v", resp.Error),
				Hints: []string{
					"The server might not support this MCP protocol version",
					"Check if the server is properly configured",
					"Verify the transport type matches the server's configuration",
				},
			}
		}

		if debugMode {
			respJSON, _ := json.MarshalIndent(resp, "", "  ")
			fmt.Println("Debug: Response:")
			fmt.Println(string(respJSON))
		}

		fmt.Printf("Successfully initialized connection!\n\n")
		output, _ := json.MarshalIndent(resp.Result, "", "  ")
		fmt.Println("Server capabilities:")
		fmt.Println(string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
