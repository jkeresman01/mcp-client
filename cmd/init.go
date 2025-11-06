package cmd

import (
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
					"version": "0.1",
				},
				"capabilities": map[string]interface{}{},
			},
		}
		resp, err := t.Send(req)
		if err != nil {
			return err
		}
		fmt.Printf("Response: %+v\n", resp)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
