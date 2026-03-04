package tools

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// registerHelloTool registers a simple "hello" tool as a starting example.
// Replace or remove this once you add your own tools.
func registerHelloTool(s *server.MCPServer) {
	tool := mcp.NewTool("hello",
		mcp.WithDescription("Says hello to the given name"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("The name to greet"),
		),
	)

	s.AddTool(tool, handleHello)
}

func handleHello(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, err := request.RequireString("name")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	greeting := fmt.Sprintf("Hello, %s! 👋", name)
	return mcp.NewToolResultText(greeting), nil
}
