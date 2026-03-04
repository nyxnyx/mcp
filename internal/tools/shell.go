package tools

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerShellTools(s *server.MCPServer) {
	s.AddTool(mcp.NewTool("run_command",
		mcp.WithDescription("Run a shell command on the local machine and return its output. Commands time out after 30 seconds."),
		mcp.WithString("command",
			mcp.Required(),
			mcp.Description("The shell command to execute (runs via /bin/sh -c)"),
		),
	), handleRunCommand)
}

func handleRunCommand(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	command, err := req.RequireString("command")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	runErr := cmd.Run()

	result := out.String()
	if runErr != nil {
		return mcp.NewToolResultText(fmt.Sprintf("exit error: %v\n%s", runErr, result)), nil
	}
	if result == "" {
		result = "(no output)"
	}
	return mcp.NewToolResultText(result), nil
}
