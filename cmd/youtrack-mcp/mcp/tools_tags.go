package mcp

import (
	"context"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"
)

func registerTagTools(s *server.MCPServer, c YouTrackClient, log zerolog.Logger) {
	s.AddTool(mcplib.NewTool("tag_list",
		mcplib.WithDescription("List all YouTrack issue tags."),
		mcplib.WithNumber("skip", mcplib.Description("pagination offset")),
		mcplib.WithNumber("top", mcplib.Description("max results")),
	), server.ToolHandlerFunc(tagListHandler(c, log)))
}

func tagListHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		out, err := c.TagList(ctx, optInt(req, "skip"), optInt(req, "top"))
		if err != nil {
			return toolErr(log, "tag_list", err)
		}
		return toolJSON(out)
	}
}
