# MCP Server (Go)

A Model Context Protocol (MCP) server scaffold built with [mcp-go](https://github.com/mark3labs/mcp-go).

## Project Structure

```
.
├── main.go                         # Entry point
├── internal/
│   ├── server/
│   │   └── server.go               # MCP server creation & transport setup
│   └── tools/
│       ├── filesystem.go           # Example filesystem info tool
│       ├── hello.go                # Example "hello" tool
│       ├── http.go                 # HTTP utility tool
│       ├── initialize.go           # Initialization helper
│       ├── register.go             # Central tool registration
│       ├── shell.go                # Shell command tool
│       └── sysinfo.go              # System info tool
├── go.mod
└── go.sum
```

## Quick Start

### Build

```bash
go build -o mcp-server .
```

### Run (stdio — default)

```bash
./mcp-server
```

### Run (SSE transport)

```bash
MCP_TRANSPORT=sse ./mcp-server
```

### Run (Streamable HTTP transport)

```bash
MCP_TRANSPORT=http ./mcp-server
```

## Adding a New Tool

1. Create a new file in `internal/tools/`, e.g. `internal/tools/mytool.go`
2. Define the tool and its handler following the pattern in `hello.go`
3. Register it in `internal/tools/register.go`

### Example

```go
// internal/tools/mytool.go
package tools

import (
    "context"
    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)

func registerMyTool(s *server.MCPServer) {
    tool := mcp.NewTool("my_tool",
        mcp.WithDescription("Does something useful"),
        mcp.WithString("input",
            mcp.Required(),
            mcp.Description("The input to process"),
        ),
    )

    s.AddTool(tool, handleMyTool)
}

func handleMyTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    input, err := request.RequireString("input")
    if err != nil {
        return mcp.NewToolResultError(err.Error()), nil
    }

    return mcp.NewToolResultText("Processed: " + input), nil
}
```

Then add `registerMyTool(s)` to `Register()` in `register.go`.

## Configuration

| Env Variable      | Default | Description                          |
|--------------------|---------|--------------------------------------|
| `MCP_TRANSPORT`    | `stdio` | Transport mode: `stdio`, `sse`, `http` |

## Dependencies

- [mcp-go](https://github.com/mark3labs/mcp-go) — Go SDK for the Model Context Protocol
