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

		resp, err := t.Send(req)
		if err != nil {
			return err
		}

		output, _ := json.MarshalIndent(resp.Result, "", "  ")
		fmt.Println("Prompts:\n", string(output))
		return nil
	},
}

var getPromptCmd = &cobra.Command{
	Use:   "get-prompt",
	Short: "Get details of a specific prompt",
	Long:  `Retrieve details of a specific prompt from the MCP server using its name.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if promptName == "" {
			return fmt.Errorf("--name is required")
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
				return fmt.Errorf("invalid --arguments JSON: %v", err)
			}
			params["arguments"] = parsedArgs
		}

		req := transport.RPCRequest{
			JSONRPC: "2.0",
			ID:      21,
			Method:  "prompts/get",
			Params:  params,
		}

		resp, err := t.Send(req)
		if err != nil {
			return err
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
