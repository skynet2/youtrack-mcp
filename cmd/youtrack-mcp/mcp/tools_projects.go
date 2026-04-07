package mcp

import (
	"context"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"
)

func registerProjectTools(s *server.MCPServer, c YouTrackClient, log zerolog.Logger) {
	s.AddTool(mcplib.NewTool("project_list",
		mcplib.WithDescription("List YouTrack projects."),
		mcplib.WithNumber("skip", mcplib.Description("pagination offset")),
		mcplib.WithNumber("top", mcplib.Description("max results")),
	), server.ToolHandlerFunc(projectListHandler(c, log)))

	s.AddTool(mcplib.NewTool("project_get",
		mcplib.WithDescription("Get a YouTrack project by ID or short name."),
		mcplib.WithString("id", mcplib.Required(), mcplib.Description("project ID or short name")),
	), server.ToolHandlerFunc(projectGetHandler(c, log)))
}

func projectListHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		out, err := c.ProjectList(ctx, optInt(req, "skip"), optInt(req, "top"))
		if err != nil {
			return toolErr(log, "project_list", err)
		}
		return toolJSON(out)
	}
}

func projectGetHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		id, err := reqString(req, "id")
		if err != nil {
			return toolErr(log, "project_get", err)
		}
		out, err := c.ProjectGet(ctx, id)
		if err != nil {
			return toolErr(log, "project_get", err)
		}
		return toolJSON(out)
	}
}
