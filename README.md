# mcp-client

Simple MCP Client supporting SSE, stdio and StreambleHTTP

Idea is to be able to test MCP server faster from CLI instead of lunching MCP inspector.

The MCP Inspector already covers a lot (UI + proxy); a CLI should focus on scripting, CI, and quick local checks, not replicate the GUI.


How MCP works:

An MCP client (for example, ChatGPT, or another application embedding a model) opens a session with a server.

An MCP server exposes:

Tools (functions the model can call, like “run SQL query” or “create GitHub issue”).

Resources (structured data endpoints, like “all my calendar events” or “config files”).

Prompts (reusable templates that the model can fill in).


Communication flows over a transport: stdio (for local processes), HTTP with Server-Sent Events (SSE), or Streamable HTTP.

Just started this, work in progess.
