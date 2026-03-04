package tools

import (
	"github.com/mark3labs/mcp-go/server"
)

// Register registers all tool handlers with the MCP server.
func Register(s *server.MCPServer) {
	registerHelloTool(s)
	registerFilesystemTools(s)
	registerShellTools(s)
	registerSysInfoTools(s)
	registerHTTPTools(s)
}
