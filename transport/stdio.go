package transport

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

type stdioTransport struct {
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	scanner *bufio.Scanner
	mu      sync.Mutex
}

func NewSTDIO(command string, args []string) Transport {
	return &stdioTransport{
		cmd: exec.Command(command, args...),
	}
}

func (t *stdioTransport) start() error {
	var err error

	t.stdin, err = t.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdin pipe: %v", err)
	}

	t.stdout, err = t.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %v", err)
	}

	if err := t.cmd.Start(); err != nil {
		return WrapError("stdio start", err)
	}

	t.scanner = bufio.NewScanner(t.stdout)
	return nil
}

func (t *stdioTransport) Send(req RPCRequest) (*RPCResponse, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Start the process if not already started
	if t.stdin == nil {
		if err := t.start(); err != nil {
			return nil, err
		}
	}

	// Send the request
	data, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	if _, err := t.stdin.Write(append(data, '\n')); err != nil {
		return nil, fmt.Errorf("failed to write request: %v", err)
	}

	// Read the response
	if !t.scanner.Scan() {
		if err := t.scanner.Err(); err != nil {
			return nil, fmt.Errorf("failed to read response: %v", err)
		}
		return nil, fmt.Errorf("no response received")
	}

	var rpcResp RPCResponse
	if err := json.Unmarshal(t.scanner.Bytes(), &rpcResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	return &rpcResp, nil
}

func (t *stdioTransport) Listen(handler func(RPCResponse)) error {
	t.mu.Lock()
	if t.stdin == nil {
		if err := t.start(); err != nil {
			t.mu.Unlock()
			return err
		}
	}
	t.mu.Unlock()

	for t.scanner.Scan() {
		var rpcResp RPCResponse
		if err := json.Unmarshal(t.scanner.Bytes(), &rpcResp); err == nil {
			handler(rpcResp)
		}
	}

	if err := t.scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %v", err)
	}

	return nil
}

func (t *stdioTransport) Close() error {
	if t.stdin != nil {
		t.stdin.Close()
	}
	if t.stdout != nil {
		t.stdout.Close()
	}
	if t.cmd != nil && t.cmd.Process != nil {
		return t.cmd.Process.Kill()
	}
	return nil
}
