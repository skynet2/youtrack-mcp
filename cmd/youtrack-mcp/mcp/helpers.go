package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog"
)

type toolHandler func(context.Context, mcplib.CallToolRequest) (*mcplib.CallToolResult, error)

func toolErr(log zerolog.Logger, tool string, err error) (*mcplib.CallToolResult, error) {
	log.Error().Err(err).Str("tool", tool).Msg("tool call failed")
	return mcplib.NewToolResultError(err.Error()), nil
}

func toolJSON(v any) (*mcplib.CallToolResult, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return mcplib.NewToolResultError(err.Error()), nil
	}
	return mcplib.NewToolResultText(string(b)), nil
}

func reqString(req mcplib.CallToolRequest, key string) (string, error) {
	return req.RequireString(key)
}

func reqInt(req mcplib.CallToolRequest, key string) (int, error) {
	return req.RequireInt(key)
}

func optString(req mcplib.CallToolRequest, key string) string {
	return req.GetString(key, "")
}

func optInt(req mcplib.CallToolRequest, key string) int {
	return req.GetInt(key, 0)
}

func reqMap(req mcplib.CallToolRequest, key string) (map[string]any, error) {
	args := req.GetArguments()
	v, ok := args[key]
	if !ok {
		return nil, fmt.Errorf("required argument %q not found", key)
	}
	m, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("argument %q is not an object", key)
	}
	return m, nil
}
