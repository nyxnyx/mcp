package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
)

var httpClient = &http.Client{Timeout: 30 * time.Second}

func registerHTTPTools(s *mcpserver.MCPServer) {
	s.AddTool(mcp.NewTool("http_request",
		mcp.WithDescription("Make an HTTP request and return the response body. Useful for calling external APIs or webhooks."),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("The URL to request"),
		),
		mcp.WithString("method",
			mcp.Description("HTTP method: GET, POST, PUT, DELETE (default: GET)"),
		),
		mcp.WithString("body",
			mcp.Description("Request body (for POST/PUT)"),
		),
		mcp.WithString("headers",
			mcp.Description(`JSON object of request headers, e.g. {"Authorization":"Bearer token","Content-Type":"application/json"}`),
		),
	), handleHTTPRequest)
}

func handleHTTPRequest(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	url, err := req.RequireString("url")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	method := "GET"
	if m, ok := req.GetArguments()["method"].(string); ok && m != "" {
		method = strings.ToUpper(m)
	}

	var bodyReader io.Reader
	if body, ok := req.GetArguments()["body"].(string); ok && body != "" {
		bodyReader = strings.NewReader(body)
	}

	httpReq, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("http_request: invalid request: %v", err)), nil
	}

	// Parse optional headers
	if headersStr, ok := req.GetArguments()["headers"].(string); ok && headersStr != "" {
		var headers map[string]string
		if err := json.Unmarshal([]byte(headersStr), &headers); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("http_request: invalid headers JSON: %v", err)), nil
		}
		for k, v := range headers {
			httpReq.Header.Set(k, v)
		}
	}

	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("http_request: %v", err)), nil
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20)) // 1MB limit
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("http_request: reading body: %v", err)), nil
	}

	result := fmt.Sprintf("HTTP %d %s\n\n%s", resp.StatusCode, resp.Status, string(respBody))
	return mcp.NewToolResultText(result), nil
}
