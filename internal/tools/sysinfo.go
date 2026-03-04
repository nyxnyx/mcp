package tools

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerSysInfoTools(s *server.MCPServer) {
	s.AddTool(mcp.NewTool("system_info",
		mcp.WithDescription("Get local system information: hostname, uptime, memory usage, disk usage, and CPU load"),
	), handleSystemInfo)
}

func handleSystemInfo(_ context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var sb strings.Builder

	// Hostname
	if h, err := os.Hostname(); err == nil {
		fmt.Fprintf(&sb, "Hostname: %s\n", h)
	}

	// Uptime
	if out := runCmd("uptime", "-p"); out != "" {
		fmt.Fprintf(&sb, "Uptime:   %s\n", strings.TrimSpace(out))
	}

	// Memory (free -h)
	if out := runCmd("free", "-h"); out != "" {
		fmt.Fprintf(&sb, "\nMemory:\n%s\n", out)
	}

	// Disk (df -h /)
	if out := runCmd("df", "-h", "/"); out != "" {
		fmt.Fprintf(&sb, "Disk (/):\n%s\n", out)
	}

	// CPU load (/proc/loadavg)
	if data, err := os.ReadFile("/proc/loadavg"); err == nil {
		fields := strings.Fields(string(data))
		if len(fields) >= 3 {
			fmt.Fprintf(&sb, "CPU load: %s %s %s (1m 5m 15m)\n", fields[0], fields[1], fields[2])
		}
	}

	return mcp.NewToolResultText(sb.String()), nil
}

func runCmd(name string, args ...string) string {
	var out bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return ""
	}
	return out.String()
}
