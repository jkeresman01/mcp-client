package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/jkeresman01/mcp-client/transport"

	"github.com/spf13/cobra"
)

var (
	promptName      string
	promptInput     string
	promptArguments string
)

var listPromptsCmd = &cobra.Command{
	Use:   "list-prompts",
	Short: "List all available prompts",
	Long:  `Retrieve and display all prompts available on the MCP server.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		t, err := getTransport()
		if err != nil {
			return err
		}
		defer t.Close()

		req := transport.RPCRequest{
			JSONRPC: "2.0",
			ID:      20,
			Method:  "prompts/list",
			Params:  map[string]interface{}{},
		}

		if debugMode {
			fmt.Println("Debug: Sending request:", req.Method)
		}

		resp, err := t.Send(req)
		if err != nil {
			return transport.WrapError("list-prompts", err)
		}

		if resp.Error != nil {
			return &transport.MCPError{
				Operation: "list-prompts",
				Err:       fmt.Errorf("server returned error: %v", resp.Error),
				Hints: []string{
					"The server may not support prompts",
					"Check server logs for more details",
				},
			}
		}

		output, _ := json.MarshalIndent(resp.Result, "", "  ")
		fmt.Println("Prompts:\n", string(output))
		return nil
	},
}

var getPromptCmd = &cobra.Command{
	Use:   "get-prompt",
	Short: "Get details of a specific prompt",
	Long: `Retrieve details of a specific prompt from the MCP server using its name.

Examples:
  mcp-client get-prompt --name greeting
  mcp-client get-prompt --name template --arguments '{"var":"value"}'`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if promptName == "" {
			return &transport.MCPError{
				Operation: "get-prompt",
				Err:       fmt.Errorf("prompt name is required"),
				Hints: []string{
					"Specify prompt name: --name <prompt-name>",
					"List available prompts: mcp-client list-prompts",
				},
			}
		}

		t, err := getTransport()
		if err != nil {
			return err
		}
		defer t.Close()

		params := map[string]interface{}{
			"name": promptName,
		}

		if promptArguments != "" && promptArguments != "{}" {
			var parsedArgs map[string]interface{}
			if err := json.Unmarshal([]byte(promptArguments), &parsedArgs); err != nil {
				return transport.NewInvalidArgumentsError(err.Error())
			}
			params["arguments"] = parsedArgs
		}

		req := transport.RPCRequest{
			JSONRPC: "2.0",
			ID:      21,
			Method:  "prompts/get",
			Params:  params,
		}

		if debugMode {
			fmt.Printf("Debug: Getting prompt '%s'\n", promptName)
		}

		resp, err := t.Send(req)
		if err != nil {
			return transport.WrapError("get-prompt", err)
		}

		if resp.Error != nil {
			errMap, ok := resp.Error.(map[string]interface{})
			if ok {
				if code, exists := errMap["code"]; exists && code == -32602 {
					return &transport.MCPError{
						Operation: "get-prompt",
						Err:       fmt.Errorf("prompt '%s' not found", promptName),
						Hints: []string{
							"List available prompts: mcp-client list-prompts",
							"Check if the prompt name is correct (case-sensitive)",
						},
					}
				}
			}
			return &transport.MCPError{
				Operation: "get-prompt",
				Err:       fmt.Errorf("server error: %v", resp.Error),
				Hints: []string{
					"Verify the prompt name is correct",
					"Check that all required arguments are provided",
					"Check server logs for more details",
				},
			}
		}

		output, _ := json.MarshalIndent(resp.Result, "", "  ")
		fmt.Println("Prompt Result:\n", string(output))
		return nil
	},
}

func init() {
	getPromptCmd.Flags().StringVar(&promptName, "name", "", "Name of the prompt to get (required)")
	getPromptCmd.Flags().StringVar(&promptArguments, "arguments", "{}", "JSON-encoded arguments to pass to the prompt")
	getPromptCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(listPromptsCmd)
	rootCmd.AddCommand(getPromptCmd)
}
