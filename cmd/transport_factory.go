package cmd

import (
	"fmt"
	"github.com/jkeresman01/mcp-client/transport"
)

func getTransport() (transport.Transport, error) {
	if debugMode {
		fmt.Printf("Debug: Creating %s transport\n", transportType)
	}

	switch transportType {
	case "streamable-http":
		if serverURL == "" {
			return nil, fmt.Errorf("--url is required for streamable-http transport")
		}
		return transport.NewStreamableHttp(serverURL), nil

	case "sse":
		if serverURL == "" {
			return nil, fmt.Errorf("--url is required for sse transport")
		}
		return transport.NewSSE(serverURL), nil

	case "stdio":
		if commandPath == "" {
			return nil, &transport.MCPError{
				Operation: "creating stdio transport",
				Err:       fmt.Errorf("--command is required for stdio transport"),
				Hints: []string{
					"Specify the server executable: --command /path/to/server",
					"Example: mcp-client list-tools --transport stdio --command ./my-mcp-server",
					"Or configure a server: mcp-client config add myserver --transport stdio --command ./server",
				},
			}
		}
		return transport.NewSTDIO(commandPath, commandArgs), nil

	default:
		return nil, &transport.MCPError{
			Operation: "selecting transport",
			Err:       fmt.Errorf("unknown transport type: %s", transportType),
			Hints: []string{
				"Valid transport types are: streamable-http, sse, stdio",
				"Example: --transport streamable-http",
			},
		}
	}
}
