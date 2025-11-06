package transport

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

type sseTransport struct {
	url       string
	client    *http.Client
	closeOnce sync.Once
	closeCh   chan struct{}
}

func NewSSE(url string) Transport {
	return &sseTransport{
		url:     url,
		client:  &http.Client{},
		closeCh: make(chan struct{}),
	}
}

func (t *sseTransport) Send(req RPCRequest) (*RPCResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", t.url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := t.client.Do(httpReq)
	if err != nil {
		return nil, NewConnectionError("sse", t.url, err)
	}
	defer resp.Body.Close()

	if resp.Header.Get("Content-Type") == "text/event-stream" {
		return t.readSSEResponse(resp.Body)
	}

	// Fallback to regular JSON response
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rpcResp RPCResponse
	if err := json.Unmarshal(data, &rpcResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON-RPC response: %v", err)
	}

	return &rpcResp, nil
}

func (t *sseTransport) readSSEResponse(body io.Reader) (*RPCResponse, error) {
	scanner := bufio.NewScanner(body)
	var dataLines []string

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "data: ") {
			dataLines = append(dataLines, strings.TrimPrefix(line, "data: "))
		} else if line == "" && len(dataLines) > 0 {
			// End of event
			data := strings.Join(dataLines, "\n")
			var rpcResp RPCResponse
			if err := json.Unmarshal([]byte(data), &rpcResp); err != nil {
				return nil, fmt.Errorf("failed to parse SSE JSON-RPC response: %v", err)
			}
			return &rpcResp, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("no valid SSE response received")
}

func (t *sseTransport) Listen(handler func(RPCResponse)) error {
	req, err := http.NewRequest("GET", t.url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	resp, err := t.client.Do(req)
	if err != nil {
		return NewConnectionError("sse", t.url, err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	var dataLines []string

	for {
		select {
		case <-t.closeCh:
			return nil
		default:
			if !scanner.Scan() {
				if err := scanner.Err(); err != nil {
					return err
				}
				return nil
			}

			line := scanner.Text()

			if strings.HasPrefix(line, "data: ") {
				dataLines = append(dataLines, strings.TrimPrefix(line, "data: "))
			} else if line == "" && len(dataLines) > 0 {
				// End of event
				data := strings.Join(dataLines, "\n")
				var rpcResp RPCResponse
				if err := json.Unmarshal([]byte(data), &rpcResp); err == nil {
					handler(rpcResp)
				}
				dataLines = nil
			}
		}
	}
}

func (t *sseTransport) Close() error {
	t.closeOnce.Do(func() {
		close(t.closeCh)
	})
	return nil
}
