package transport

import (
	"fmt"
	"strings"
)

// MCPError represents an enhanced error with troubleshooting information
type MCPError struct {
	Operation string
	Err       error
	Hints     []string
}

func (e *MCPError) Error() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Error during %s: %v\n", e.Operation, e.Err))

	if len(e.Hints) > 0 {
		sb.WriteString("\nTroubleshooting hints:\n")
		for i, hint := range e.Hints {
			sb.WriteString(fmt.Sprintf("   %d. %s\n", i+1, hint))
		}
	}

	return sb.String()
}

// WrapError creates an MCPError with contextual hints
func WrapError(operation string, err error) error {
	if err == nil {
		return nil
	}

	mcpErr := &MCPError{
		Operation: operation,
		Err:       err,
		Hints:     generateHints(operation, err),
	}

	return mcpErr
}

// generateHints provides context-specific troubleshooting hints
func generateHints(operation string, err error) []string {
	errMsg := strings.ToLower(err.Error())
	var hints []string

	// Connection errors
	if strings.Contains(errMsg, "connection refused") {
		hints = append(hints, "Make sure the MCP server is running")
		hints = append(hints, "Check if the URL/port is correct")
		hints = append(hints, "Verify firewall settings aren't blocking the connection")
	}

	if strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "deadline exceeded") {
		hints = append(hints, "The server might be slow to respond - try increasing timeout")
		hints = append(hints, "Check your network connection")
		hints = append(hints, "Verify the server is not overloaded")
	}

	// Transport-specific errors
	if strings.Contains(operation, "stdio") {
		if strings.Contains(errMsg, "no such file") || strings.Contains(errMsg, "executable file not found") {
			hints = append(hints, "Verify the command path is correct")
			hints = append(hints, "Make sure the executable has proper permissions")
			hints = append(hints, "Try using an absolute path instead of relative")
		}
		if strings.Contains(errMsg, "permission denied") {
			hints = append(hints, "The executable might not have execute permissions")
			hints = append(hints, "Try: chmod +x /path/to/executable")
		}
	}

	// JSON/Protocol errors
	if strings.Contains(errMsg, "json") || strings.Contains(errMsg, "unmarshal") {
		hints = append(hints, "The server response might not be valid JSON-RPC")
		hints = append(hints, "Check if you're using the correct MCP protocol version")
		hints = append(hints, "Enable debug mode to see raw responses: --debug")
	}

	// Tool/Resource errors
	if strings.Contains(operation, "tool") && strings.Contains(errMsg, "not found") {
		hints = append(hints, "Run 'mcp-client list-tools' to see available tools")
		hints = append(hints, "Check if the tool name is spelled correctly (case-sensitive)")
	}

	if strings.Contains(operation, "resource") && strings.Contains(errMsg, "not found") {
		hints = append(hints, "Run 'mcp-client list-resources' to see available resources")
		hints = append(hints, "Verify the resource URI/ID is correct")
	}

	// Authentication errors
	if strings.Contains(errMsg, "unauthorized") || strings.Contains(errMsg, "401") {
		hints = append(hints, "Check if authentication is required")
		hints = append(hints, "Verify API keys or credentials are configured correctly")
	}

	if strings.Contains(errMsg, "forbidden") || strings.Contains(errMsg, "403") {
		hints = append(hints, "You might not have permission to access this resource")
		hints = append(hints, "Check if your credentials have the necessary scope")
	}

	// Server errors
	if strings.Contains(errMsg, "500") || strings.Contains(errMsg, "internal server error") {
		hints = append(hints, "This is a server-side error")
		hints = append(hints, "Check the MCP server logs for more details")
		hints = append(hints, "Try the request again - it might be a temporary issue")
	}

	// Configuration errors
	if strings.Contains(errMsg, "server") && strings.Contains(errMsg, "not found in configuration") {
		hints = append(hints, "Run 'mcp-client config list' to see configured servers")
		hints = append(hints, "Add a server with: mcp-client config add <name> --url <url>")
	}

	// Generic fallback hints
	if len(hints) == 0 {
		hints = append(hints, "Enable debug mode for more details: --debug")
		hints = append(hints, "Check the server logs for error details")
		hints = append(hints, "Verify your MCP client is up to date")
	}

	return hints
}

// Common error constructors for better UX
func NewConnectionError(transport, address string, err error) error {
	return &MCPError{
		Operation: fmt.Sprintf("%s connection to %s", transport, address),
		Err:       err,
		Hints: []string{
			"Verify the server is running and accessible",
			"Check network connectivity",
			"Ensure the transport type matches the server configuration",
		},
	}
}

func NewToolNotFoundError(toolName string) error {
	return &MCPError{
		Operation: "calling tool",
		Err:       fmt.Errorf("tool '%s' not found", toolName),
		Hints: []string{
			"Run 'mcp-client list-tools' to see available tools",
			"Check if the tool name is spelled correctly (case-sensitive)",
			"Verify the server has initialized properly",
		},
	}
}

func NewInvalidArgumentsError(details string) error {
	return &MCPError{
		Operation: "parsing arguments",
		Err:       fmt.Errorf("invalid arguments: %s", details),
		Hints: []string{
			"Arguments must be valid JSON format",
			"Example: --args '{\"param1\": \"value1\", \"param2\": 123}'",
			"Use single quotes around the JSON string",
		},
	}
}
