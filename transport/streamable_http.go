package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type streamableHttpTransport struct {
	url string
}

func NewStreamableHttp(url string) Transport {
	return &streamableHttpTransport{url: url}
}

func (t *streamableHttpTransport) Send(req RPCRequest) (*RPCResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(t.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, NewConnectionError("streamable-http", t.url, err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rpcResp RPCResponse
	if err := json.Unmarshal(data, &rpcResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON-RPC response: %v\nRaw: %s", err, string(data))
	}

	return &rpcResp, nil
}

func (t *streamableHttpTransport) Listen(handler func(RPCResponse)) error {
	// StreamableHttp is request/response — no continuous listen
	return fmt.Errorf("Listen is not supported for StreamableHttp transport")
}

func (t *streamableHttpTransport) Close() error {
	// Nothing to close for simple HTTP transport
	return nil
}
