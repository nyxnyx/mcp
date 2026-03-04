package main

import (
	"fmt"
	"log"
	"os"

	"github.com/grzecho/mcp-server/internal/server"
	"github.com/grzecho/mcp-server/internal/tools"
)

func main() {
	// Configure logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Determine transport mode from environment variable
	// Supported: "stdio" (default), "sse", "http"
	transport := os.Getenv("MCP_TRANSPORT")
	if transport == "" {
		transport = "stdio"
	}

	// Create and configure the MCP server with initialization hooks.
	// Set MCP_VERBOSE=1 to enable protocol-level request/response logging.
	verbose := os.Getenv("MCP_VERBOSE") != ""
	hooks := tools.InitializeHooks(verbose)
	srv := server.New(hooks)

	// Register all tools
	tools.Register(srv)

	// Start serving
	if err := server.Serve(srv, transport); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
