package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jkeresman01/mcp-client/transport"
	"github.com/spf13/cobra"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Start an interactive MCP session",
	Long: `Start an interactive session where you can execute multiple commands
without having to type 'mcp-client' each time. Type 'help' for available commands
or 'exit' to quit.`,
	RunE: runInteractive,
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
}

func runInteractive(cmd *cobra.Command, args []string) error {
	// Get transport once for the session
	t, err := getTransport()
	if err != nil {
		return err
	}
	defer t.Close()

	// Try to initialize the connection
	if err := initializeConnection(t); err != nil {
		fmt.Printf("Warning: Failed to initialize connection: %v\n", err)
		fmt.Println("You can still try commands, but the server may not be ready.\n")
	} else {
		fmt.Println("Connected successfully!\n")
	}

	printWelcome()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("mcp> ")

		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if shouldExit(line) {
			fmt.Println("Goodbye!")
			break
		}

		if err := handleInteractiveCommand(t, line); err != nil {
			fmt.Printf("%v\n", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %v", err)
	}

	return nil
}

func printWelcome() {
	fmt.Println("-----------    MCP Client - Interactive Mode  ----------------")
	fmt.Println()
	fmt.Println("Available commands:")
	fmt.Println("  help                          - Show this help message")
	fmt.Println("  list-tools                    - List all available tools")
	fmt.Println("  list-resources                - List all available resources")
	fmt.Println("  list-prompts                  - List all available prompts")
	fmt.Println("  call <tool> <args>            - Call a tool with JSON args")
	fmt.Println("  get-resource <uri>            - Get resource content")
	fmt.Println("  get-prompt <name> [args]      - Get prompt details")
	fmt.Println("  exit, quit, q                 - Exit interactive mode")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  call calculator {\"op\":\"add\",\"a\":5,\"b\":3}")
	fmt.Println("  get-resource file:///path/to/file")
	fmt.Println()
}

func shouldExit(line string) bool {
	lower := strings.ToLower(line)
	return lower == "exit" || lower == "quit" || lower == "q"
}

func initializeConnection(t transport.Transport) error {
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

	_, err := t.Send(req)
	return err
}

func handleInteractiveCommand(t transport.Transport, line string) error {
	parts := parseCommandLine(line)
	if len(parts) == 0 {
		return nil
	}

	command := strings.ToLower(parts[0])

	switch command {
	case "help", "h", "?":
		printWelcome()
		return nil

	case "list-tools", "lt":
		return listToolsInteractive(t)

	case "list-resources", "lr":
		return listResourcesInteractive(t)

	case "list-prompts", "lp":
		return listPromptsInteractive(t)

	case "call", "c":
		if len(parts) < 2 {
			return fmt.Errorf("usage: call <tool-name> <json-args>\nExample: call calculator {\"op\":\"add\",\"a\":5,\"b\":3}")
		}
		toolName := parts[1]
		args := "{}"
		if len(parts) >= 3 {
			args = strings.Join(parts[2:], " ")
		}
		return callToolInteractive(t, toolName, args)

	case "get-resource", "gr":
		if len(parts) < 2 {
			return fmt.Errorf("usage: get-resource <resource-uri>")
		}
		return getResourceInteractive(t, parts[1])

	case "get-prompt", "gp":
		if len(parts) < 2 {
			return fmt.Errorf("usage: get-prompt <prompt-name> [json-args]")
		}
		promptName := parts[1]
		args := "{}"
		if len(parts) >= 3 {
			args = strings.Join(parts[2:], " ")
		}
		return getPromptInteractive(t, promptName, args)

	default:
		return fmt.Errorf("unknown command: %s\nType 'help' for available commands", command)
	}
}

func parseCommandLine(line string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false
	inJSON := false
	braceCount := 0

	for _, char := range line {
		switch char {
		case '"':
			inQuotes = !inQuotes
			current.WriteRune(char)
		case '{':
			braceCount++
			inJSON = true
			current.WriteRune(char)
		case '}':
			braceCount--
			if braceCount == 0 {
				inJSON = false
			}
			current.WriteRune(char)
		case ' ':
			if inQuotes || inJSON {
				current.WriteRune(char)
			} else if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(char)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

func listToolsInteractive(t transport.Transport) error {
	req := transport.RPCRequest{
		JSONRPC: "2.0",
		ID:      2,
		Method:  "tools/list",
		Params:  map[string]interface{}{},
	}

	resp, err := t.Send(req)
	if err != nil {
		return transport.WrapError("list-tools", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("server error: %v", resp.Error)
	}

	output, _ := json.MarshalIndent(resp.Result, "", "  ")
	fmt.Printf("Tools:\n%s\n", string(output))
	return nil
}

func listResourcesInteractive(t transport.Transport) error {
	req := transport.RPCRequest{
		JSONRPC: "2.0",
		ID:      10,
		Method:  "resources/list",
		Params:  map[string]interface{}{},
	}

	resp, err := t.Send(req)
	if err != nil {
		return transport.WrapError("list-resources", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("server error: %v", resp.Error)
	}

	output, _ := json.MarshalIndent(resp.Result, "", "  ")
	fmt.Printf("Resources:\n%s\n", string(output))
	return nil
}

func listPromptsInteractive(t transport.Transport) error {
	req := transport.RPCRequest{
		JSONRPC: "2.0",
		ID:      20,
		Method:  "prompts/list",
		Params:  map[string]interface{}{},
	}

	resp, err := t.Send(req)
	if err != nil {
		return transport.WrapError("list-prompts", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("server error: %v", resp.Error)
	}

	output, _ := json.MarshalIndent(resp.Result, "", "  ")
	fmt.Printf("Prompts:\n%s\n", string(output))
	return nil
}

func callToolInteractive(t transport.Transport, toolName, argsJSON string) error {
	var parsedArgs map[string]interface{}
	if err := json.Unmarshal([]byte(argsJSON), &parsedArgs); err != nil {
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

	resp, err := t.Send(req)
	if err != nil {
		return transport.WrapError("call-tool", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("server error: %v", resp.Error)
	}

	output, _ := json.MarshalIndent(resp.Result, "", "  ")
	fmt.Printf(" Tool Result:\n%s\n", string(output))
	return nil
}

func getResourceInteractive(t transport.Transport, uri string) error {
	req := transport.RPCRequest{
		JSONRPC: "2.0",
		ID:      11,
		Method:  "resources/read",
		Params: map[string]interface{}{
			"uri": uri,
		},
	}

	resp, err := t.Send(req)
	if err != nil {
		return transport.WrapError("get-resource", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("server error: %v", resp.Error)
	}

	output, _ := json.MarshalIndent(resp.Result, "", "  ")
	fmt.Printf("Resource:\n%s\n", string(output))
	return nil
}

func getPromptInteractive(t transport.Transport, name, argsJSON string) error {
	params := map[string]interface{}{
		"name": name,
	}

	if argsJSON != "{}" {
		var parsedArgs map[string]interface{}
		if err := json.Unmarshal([]byte(argsJSON), &parsedArgs); err != nil {
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

	resp, err := t.Send(req)
	if err != nil {
		return transport.WrapError("get-prompt", err)
	}

	if resp.Error != nil {
		return fmt.Errorf("server error: %v", resp.Error)
	}

	output, _ := json.MarshalIndent(resp.Result, "", "  ")
	fmt.Printf("💬 Prompt:\n%s\n", string(output))
	return nil
}
