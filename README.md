# mcp-client

Simple MCP Client supporting SSE, stdio and StreamableHTTP with **interactive mode**, **configuration files**, and **enhanced error messages**.

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

## Quick Start

### 1. Initialize Configuration

```bash
# Create config file with example server
mcp-client config init

# Add your own server
mcp-client config add myserver \
  --url http://localhost:8765 \
  --transport streamable-http \
  --default
```

### 2. Start Interactive Mode

```bash
mcp-client interactive --server myserver
```

Inside interactive mode:
```
mcp> list-tools
mcp> call calculator {"op":"add","a":5,"b":3}
mcp> list-resources
mcp> help
mcp> exit
```

### 3. Use CLI Commands

```bash
# List available tools
mcp-client list-tools --server myserver

# Call a tool
mcp-client call-tool --name calculator \
  --args '{"op":"add","a":5,"b":3}' \
  --server myserver
```

## Configuration Management

### List Configured Servers

```bash
mcp-client config list
```

Output:
```
Configured servers:

* local:
    Transport: streamable-http
    URL:       http://localhost:8765

  production:
    Transport: sse
    URL:       https://api.example.com/mcp

Default server: local (marked with *)
```

### Add a Server

**Streamable HTTP:**
```bash
mcp-client config add prod \
  --url https://api.example.com/mcp \
  --transport streamable-http \
  --default
```

**SSE:**
```bash
mcp-client config add sse-server \
  --url http://localhost:3000/events \
  --transport sse
```

**Stdio:**
```bash
mcp-client config add local-stdio \
  --transport stdio \
  --command /path/to/mcp-server \
  --args "arg1,arg2"
```

### Remove a Server

```bash
mcp-client config remove prod
```

### Set Default Server

```bash
mcp-client config set-default local
```

### Show Configuration

```bash
# Show specific server
mcp-client config show prod

# Show overall config info
mcp-client config show
```

## Interactive Mode

Interactive mode provides a REPL-style interface for faster testing:

```bash
# Start with configured server
mcp-client interactive --server myserver

# Start with direct connection
mcp-client interactive \
  --transport streamable-http \
  --url http://localhost:8765
```

### Interactive Commands

| Command | Alias | Description | Example |
|---------|-------|-------------|---------|
| `help` | `h`, `?` | Show help | `help` |
| `list-tools` | `lt` | List all tools | `list-tools` |
| `list-resources` | `lr` | List resources | `list-resources` |
| `list-prompts` | `lp` | List prompts | `list-prompts` |
| `call` | `c` | Call a tool | `call calculator {"op":"add","a":5,"b":3}` |
| `get-resource` | `gr` | Get resource | `get-resource file:///path/to/file` |
| `get-prompt` | `gp` | Get prompt | `get-prompt greeting {"name":"Alice"}` |
| `exit` | `quit`, `q` | Exit | `exit` |

## CLI Commands

### Initialize Connection

```bash
mcp-client init --server myserver

# Or without config
mcp-client init \
  --transport streamable-http \
  --url http://localhost:8765
```

### List Tools

```bash
mcp-client list-tools --server myserver
```

### Call a Tool

```bash
mcp-client call-tool \
  --name "calculator" \
  --args '{"op":"add","a":5,"b":3}' \
  --server myserver
```

### List Resources

```bash
mcp-client list-resources --server myserver
```

### Get Resource

```bash
mcp-client get-resource \
  --id "file:///path/to/file" \
  --server myserver
```

### List Prompts

```bash
mcp-client list-prompts --server myserver
```

### Get Prompt

```bash
mcp-client get-prompt \
  --name "greeting" \
  --arguments '{"name":"Alice"}' \
  --server myserver
```

## Transport Types

### Streamable HTTP
Standard HTTP request/response:
```bash
mcp-client list-tools \
  --transport streamable-http \
  --url http://localhost:8765
```

### SSE (Server-Sent Events)
Streaming responses:
```bash
mcp-client list-tools \
  --transport sse \
  --url http://localhost:8765
```

### Stdio
Communicate with local process:
```bash
mcp-client list-tools \
  --transport stdio \
  --command "/path/to/mcp-server" \
  --args "arg1,arg2"
```

## Debug Mode

Enable debug mode for detailed request/response information:

```bash
mcp-client list-tools --server myserver --debug
```

Output includes:
- Full request JSON
- Full response JSON
- Connection details
- Timing information
## Examples

### Example 1: Quick Test with Config

```bash
# Setup
mcp-client config add local \
  --url http://localhost:8765 \
  --transport streamable-http \
  --default

# Test
mcp-client init --server local
mcp-client list-tools --server local
mcp-client call-tool --name echo --args '{"text":"Hello"}' --server local
```

### Example 2: Interactive Session

```bash
mcp-client interactive --server local
```

```
Connected successfully!

mcp> list-tools
Tools:
{
  "tools": [
    {
      "name": "calculator",
      "description": "Perform basic calculations"
    }
  ]
}

mcp> call calculator {"op":"add","a":5,"b":3}
Tool Result:
{
  "content": [
    {
      "type": "text",
      "text": "Result: 8"
    }
  ]
}

mcp> exit
Goodbye!
```

### Example 3: Stdio with Local Server

```bash
# Add stdio server
mcp-client config add local-fs \
  --transport stdio \
  --command "./filesystem-server" \
  --default

# Use it
mcp-client interactive --server local-fs
```

### Example 4: Multiple Servers

```bash
# Configure multiple servers
mcp-client config add dev --url http://localhost:8765 --transport streamable-http
mcp-client config add staging --url https://staging.example.com/mcp --transport sse
mcp-client config add prod --url https://api.example.com/mcp --transport streamable-http

# Switch between them
mcp-client list-tools --server dev
mcp-client list-tools --server staging
mcp-client list-tools --server prod
```

## Configuration File Format

The configuration file (`~/.mcp-config.json`) uses this format:

```json
{
  "default_server": "local",
  "servers": {
    "local": {
      "url": "http://localhost:8765",
      "transport": "streamable-http"
    },
    "production": {
      "url": "https://api.example.com/mcp",
      "transport": "sse"
    },
    "local-stdio": {
      "transport": "stdio",
      "command": "/path/to/mcp-server",
      "args": ["--verbose", "--port", "8080"]
    }
  }
}
```

## Global Flags

| Flag | Description | Example |
|------|-------------|---------|
| `--server` | Use named server from config | `--server prod` |
| `--config` | Custom config file path | `--config ./my-config.json` |
| `--debug` | Enable debug output | `--debug` |
| `--transport` | Transport type | `--transport streamable-http` |
| `--url` | Server URL | `--url http://localhost:8765` |

