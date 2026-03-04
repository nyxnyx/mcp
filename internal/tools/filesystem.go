package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerFilesystemTools(s *server.MCPServer) {
	// read_file
	s.AddTool(mcp.NewTool("read_file",
		mcp.WithDescription("Read the contents of a file from the local filesystem"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Absolute or relative path to the file"),
		),
	), handleReadFile)

	// write_file
	s.AddTool(mcp.NewTool("write_file",
		mcp.WithDescription("Write content to a file on the local filesystem (creates or overwrites)"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Absolute or relative path to the file"),
		),
		mcp.WithString("content",
			mcp.Required(),
			mcp.Description("Content to write to the file"),
		),
	), handleWriteFile)

	// list_dir
	s.AddTool(mcp.NewTool("list_dir",
		mcp.WithDescription("List files and directories at a given path"),
		mcp.WithString("path",
			mcp.Required(),
			mcp.Description("Absolute or relative path to the directory"),
		),
	), handleListDir)
}

func handleReadFile(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, err := req.RequireString("path")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("read_file: %v", err)), nil
	}
	return mcp.NewToolResultText(string(data)), nil
}

func handleWriteFile(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, err := req.RequireString("path")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	content, err := req.RequireString("content")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("write_file mkdir: %v", err)), nil
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("write_file: %v", err)), nil
	}
	return mcp.NewToolResultText(fmt.Sprintf("Written %d bytes to %s", len(content), path)), nil
}

func handleListDir(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	path, err := req.RequireString("path")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("list_dir: %v", err)), nil
	}
	var sb strings.Builder
	for _, e := range entries {
		info, _ := e.Info()
		if e.IsDir() {
			fmt.Fprintf(&sb, "[DIR]  %s\n", e.Name())
		} else {
			fmt.Fprintf(&sb, "[FILE] %s (%d bytes)\n", e.Name(), info.Size())
		}
	}
	return mcp.NewToolResultText(sb.String()), nil
}
