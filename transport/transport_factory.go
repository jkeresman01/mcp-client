package cmd

import (
	"fmt"
	"github.com/jkeresman01/mcp-client/transport"
)

func getTransport() (transport.Transport, error) {
	switch transportType {
	case "streamable-http":
		return transport.NewStreamableHttp(serverURL), nil
	// case "sse":
	//     return transport.NewSSE(serverURL), nil
	// case "stdio":
	//     if commandPath == "" {
	//         return nil, fmt.Errorf("--command is required for stdio transport")
	//     }
	//     return transport.NewSTDIO(commandPath, commandArgs), nil
	default:
		return nil, fmt.Errorf("unknown transport type: %s", transportType)
	}
}
