package tools

import (
	"context"
	"encoding/json"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// InitializeHooks sets up lifecycle hooks for the MCP server.
//
// When verbose is true, every JSON-RPC request, successful response, and
// error is logged to stderr — useful for debugging protocol-level traffic.
//
// Per the MCP spec (2024-11-05), initialization involves:
//   - Protocol version negotiation
//   - Capability exchange between client and server
//   - Implementation info sharing
func InitializeHooks(verbose bool) *server.Hooks {
	hooks := &server.Hooks{}

	// Initialization lifecycle hooks — always active.
	hooks.AddBeforeInitialize(func(ctx context.Context, id any, message *mcp.InitializeRequest) {
		log.Printf("[initialize] Client connecting: %s v%s (protocol: %s)",
			message.Params.ClientInfo.Name,
			message.Params.ClientInfo.Version,
			message.Params.ProtocolVersion,
		)
	})

	hooks.AddAfterInitialize(func(ctx context.Context, id any, message *mcp.InitializeRequest, result *mcp.InitializeResult) {
		log.Printf("[initialized] Session established with %s v%s — server: %s v%s (protocol: %s)",
			message.Params.ClientInfo.Name,
			message.Params.ClientInfo.Version,
			result.ServerInfo.Name,
			result.ServerInfo.Version,
			result.ProtocolVersion,
		)
		if result.Capabilities.Tools != nil {
			log.Printf("[initialized] Tools capability enabled (listChanged: %v)", result.Capabilities.Tools.ListChanged)
		}
		if result.Capabilities.Resources != nil {
			log.Printf("[initialized] Resources capability enabled")
		}
		if result.Capabilities.Prompts != nil {
			log.Printf("[initialized] Prompts capability enabled")
		}
		if result.Capabilities.Logging != nil {
			log.Printf("[initialized] Logging capability enabled")
		}
	})

	if !verbose {
		return hooks
	}

	// Verbose: log every incoming request.
	hooks.AddBeforeAny(func(ctx context.Context, id any, method mcp.MCPMethod, message any) {
		b, _ := json.Marshal(message)
		log.Printf("[req] id=%v method=%s body=%s", id, method, b)
	})

	// Verbose: log every successful response.
	hooks.AddOnSuccess(func(ctx context.Context, id any, method mcp.MCPMethod, message any, result any) {
		b, _ := json.Marshal(result)
		log.Printf("[res] id=%v method=%s result=%s", id, method, b)
	})

	// Verbose: log every error.
	hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
		log.Printf("[err] id=%v method=%s error=%v", id, method, err)
	})

	return hooks
}
