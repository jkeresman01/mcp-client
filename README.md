# mcp-client

Simple MCP Client supporting SSE, stdio and StreamableHTTP

Idea is to be able to test MCP server faster from CLI instead of lunching MCP inspector.

The MCP Inspector already covers a lot (UI + proxy); a CLI should focus on scripting, CI, and quick local checks, not replicate the GUI.

## How MCP works

An MCP client (for example, ChatGPT, or another application embedding a model) opens a session with a server.

An MCP server exposes:

- **Tools** (functions the model can call, like "run SQL query" or "create GitHub issue")
- **Resources** (structured data endpoints, like "all my calendar events" or "config files")
- **Prompts** (reusable templates that the model can fill in)

Communication flows over a transport: stdio (for local processes), HTTP with Server-Sent Events (SSE), or Streamable HTTP.

More info: https://modelcontextprotocol.io/specification

## Installation

```bash
go build -o mcp-client
```

## Usage

### Initialize Connection

```bash
mcp-client init --transport streamable-http --url http://localhost:8765
```

### List Tools

```bash
mcp-client list-tools --transport streamable-http --url http://localhost:8765
```

### Call a Tool

```bash
mcp-client call-tool --name "my-tool" --args '{"param": "value"}' --transport streamable-http --url http://localhost:8765
```

### List Resources

```bash
mcp-client list-resources --transport streamable-http --url http://localhost:8765
```

### Get Resource

```bash
mcp-client get-resource --id "resource-1" --transport streamable-http --url http://localhost:8765
```

### List Prompts

```bash
mcp-client list-prompts --transport streamable-http --url http://localhost:8765
```

### Run a Prompt

```bash
mcp-client run-prompt --name "my-prompt" --input '{"var": "value"}' --transport streamable-http --url http://localhost:8765
```

## Transport Types

- `streamable-http`: Standard HTTP request/response
- `sse`: Server-Sent Events for streaming
- `stdio`: Communicate with local process via stdin/stdout

### Using stdio

```bash
mcp-client list-tools --transport stdio --command "/path/to/mcp-server" --args "arg1,arg2"
```
