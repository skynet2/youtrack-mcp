package mcp

import (
	"context"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"
)

func registerUserTools(s *server.MCPServer, c YouTrackClient, log zerolog.Logger) {
	s.AddTool(mcplib.NewTool("user_list",
		mcplib.WithDescription("List or search YouTrack users."),
		mcplib.WithString("query", mcplib.Description("search query")),
		mcplib.WithNumber("skip", mcplib.Description("pagination offset")),
		mcplib.WithNumber("top", mcplib.Description("max results")),
	), server.ToolHandlerFunc(userListHandler(c, log)))

	s.AddTool(mcplib.NewTool("user_current",
		mcplib.WithDescription("Get the currently authenticated YouTrack user."),
	), server.ToolHandlerFunc(userCurrentHandler(c, log)))
}

func userListHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		out, err := c.UserList(ctx, optString(req, "query"), optInt(req, "skip"), optInt(req, "top"))
		if err != nil {
			return toolErr(log, "user_list", err)
		}
		return toolJSON(out)
	}
}

func userCurrentHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		out, err := c.UserCurrent(ctx)
		if err != nil {
			return toolErr(log, "user_current", err)
		}
		return toolJSON(out)
	}
}
