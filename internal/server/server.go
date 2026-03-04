package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/mark3labs/mcp-go/server"
)

const (
	serverName    = "my-mcp-server"
	serverVersion = "0.1.0"
)

// New creates and configures a new MCP server instance.
// An optional *server.Hooks can be provided to hook into lifecycle events
// such as initialization (see MCP spec 2024-11-05 lifecycle).
func New(hooks *server.Hooks) *server.MCPServer {
	opts := []server.ServerOption{
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	}

	if hooks != nil {
		opts = append(opts, server.WithHooks(hooks))
	}

	s := server.NewMCPServer(
		serverName,
		serverVersion,
		opts...,
	)

	return s
}

// uuidSessionIDGenerator generates a new UUID-based session ID for each SSE connection.
func uuidSessionIDGenerator(_ context.Context, _ *http.Request) (string, error) {
	return uuid.New().String(), nil
}

// Serve starts the MCP server using the specified transport.
// Supported transports: "stdio", "sse", "http".
func Serve(s *server.MCPServer, transport string) error {
	switch transport {
	case "stdio":
		return server.ServeStdio(s)

	case "sse":
		sseServer := server.NewSSEServer(s,
			server.WithBaseURL("http://localhost:8080"),
			server.WithSessionIDGenerator(uuidSessionIDGenerator),
		)
		fmt.Println("Starting SSE server on :8080")
		return sseServer.Start(":8080")

	case "http":
		httpServer := server.NewStreamableHTTPServer(s)
		fmt.Println("Starting streamable HTTP server on :8080 — endpoint: http://localhost:8080/mcp")
		return httpServer.Start(":8080")

	default:
		return fmt.Errorf("unsupported transport: %s (use stdio, sse, or http)", transport)
	}
}
