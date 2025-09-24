package cmd

import (
	"encoding/json"
	"fmt"
	"mcp-client/transport"

	"github.com/spf13/cobra"
)

var (
	promptName  string
	promptInput string
)

var listPromptsCmd = &cobra.Command{
	Use:   "list-prompts",
	Short: "List all available prompts",
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

var runPromptCmd = &cobra.Command{
	Use:   "run-prompt",
	Short: "Run a prompt with input",
	RunE: func(cmd *cobra.Command, args []string) error {
		if promptName == "" {
			return fmt.Errorf("--name is required")
		}

		t, err := getTransport()
		if err != nil {
			return err
		}
		defer t.Close()

		var parsedInput map[string]interface{}
		if err := json.Unmarshal([]byte(promptInput), &parsedInput); err != nil {
			return fmt.Errorf("invalid --input JSON: %v", err)
		}

		req := transport.RPCRequest{
			JSONRPC: "2.0",
			ID:      21,
			Method:  "prompts/run",
			Params: map[string]interface{}{
				"name":  promptName,
				"input": parsedInput,
			},
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
	runPromptCmd.Flags().StringVar(&promptName, "name", "", "Name of the prompt to run (required)")
	runPromptCmd.Flags().StringVar(&promptInput, "input", "{}", "JSON-encoded input to pass to the prompt")
	runPromptCmd.MarkFlagRequired("name")

	rootCmd.AddCommand(listPromptsCmd)
	rootCmd.AddCommand(runPromptCmd)
}
