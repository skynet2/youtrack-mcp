package mcp

import (
	"context"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"
)

func registerAgileTools(s *server.MCPServer, c YouTrackClient, log zerolog.Logger) {
	s.AddTool(mcplib.NewTool("agile_list",
		mcplib.WithDescription("List YouTrack agile boards."),
		mcplib.WithNumber("skip", mcplib.Description("pagination offset")),
		mcplib.WithNumber("top", mcplib.Description("max results")),
	), server.ToolHandlerFunc(agileListHandler(c, log)))

	s.AddTool(mcplib.NewTool("sprint_list",
		mcplib.WithDescription("List sprints for a YouTrack agile board."),
		mcplib.WithString("agile_id", mcplib.Required(), mcplib.Description("agile board ID")),
		mcplib.WithNumber("skip", mcplib.Description("pagination offset")),
		mcplib.WithNumber("top", mcplib.Description("max results")),
	), server.ToolHandlerFunc(sprintListHandler(c, log)))

	s.AddTool(mcplib.NewTool("sprint_get",
		mcplib.WithDescription("Get details of a specific sprint including its issues."),
		mcplib.WithString("agile_id", mcplib.Required(), mcplib.Description("agile board ID")),
		mcplib.WithString("sprint_id", mcplib.Required(), mcplib.Description("sprint ID")),
	), server.ToolHandlerFunc(sprintGetHandler(c, log)))
}

func agileListHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		out, err := c.AgileList(ctx, optInt(req, "skip"), optInt(req, "top"))
		if err != nil {
			return toolErr(log, "agile_list", err)
		}
		return toolJSON(out)
	}
}

func sprintListHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		agileID, err := reqString(req, "agile_id")
		if err != nil {
			return toolErr(log, "sprint_list", err)
		}
		out, err := c.SprintList(ctx, agileID, optInt(req, "skip"), optInt(req, "top"))
		if err != nil {
			return toolErr(log, "sprint_list", err)
		}
		return toolJSON(out)
	}
}

func sprintGetHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		agileID, err := reqString(req, "agile_id")
		if err != nil {
			return toolErr(log, "sprint_get", err)
		}
		sprintID, err := reqString(req, "sprint_id")
		if err != nil {
			return toolErr(log, "sprint_get", err)
		}
		out, err := c.SprintGet(ctx, agileID, sprintID)
		if err != nil {
			return toolErr(log, "sprint_get", err)
		}
		return toolJSON(out)
	}
}
