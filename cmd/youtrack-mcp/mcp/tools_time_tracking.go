package mcp

import (
	"context"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"

	"github.com/skynet2/youtrack-mcp/pkg/youtrack"
)

func registerTimeTrackingTools(s *server.MCPServer, c YouTrackClient, log zerolog.Logger) {
	s.AddTool(mcplib.NewTool("workitem_list",
		mcplib.WithDescription("List time tracking work items for a YouTrack issue."),
		mcplib.WithString("issue_id", mcplib.Required(), mcplib.Description("issue readable ID")),
		mcplib.WithNumber("skip", mcplib.Description("pagination offset")),
		mcplib.WithNumber("top", mcplib.Description("max results")),
	), server.ToolHandlerFunc(workItemListHandler(c, log)))

	s.AddTool(mcplib.NewTool("workitem_create",
		mcplib.WithDescription("Add a time tracking work item to a YouTrack issue."),
		mcplib.WithString("issue_id", mcplib.Required(), mcplib.Description("issue readable ID")),
		mcplib.WithNumber("minutes", mcplib.Required(), mcplib.Description("duration in minutes")),
		mcplib.WithString("text", mcplib.Description("work item description")),
		mcplib.WithString("type_id", mcplib.Description("work item type ID")),
		mcplib.WithNumber("date", mcplib.Description("date as unix timestamp in ms")),
	), server.ToolHandlerFunc(workItemCreateHandler(c, log)))
}

func workItemListHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		issueID, err := reqString(req, "issue_id")
		if err != nil {
			return toolErr(log, "workitem_list", err)
		}
		out, err := c.WorkItemList(ctx, issueID, optInt(req, "skip"), optInt(req, "top"))
		if err != nil {
			return toolErr(log, "workitem_list", err)
		}
		return toolJSON(out)
	}
}

func workItemCreateHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		issueID, err := reqString(req, "issue_id")
		if err != nil {
			return toolErr(log, "workitem_create", err)
		}
		minutes, err := reqInt(req, "minutes")
		if err != nil {
			return toolErr(log, "workitem_create", err)
		}
		m := int32(minutes)
		item := youtrack.IssueWorkItem{
			Duration: &youtrack.DurationValue{Minutes: &m},
		}
		if t := optString(req, "text"); t != "" {
			item.Text = t
		}
		if tid := optString(req, "type_id"); tid != "" {
			item.WorkType = &youtrack.WorkItemType{ID: tid}
		}
		if d := optInt(req, "date"); d > 0 {
			d64 := int64(d)
			item.Date = &d64
		}
		out, err := c.WorkItemCreate(ctx, issueID, item)
		if err != nil {
			return toolErr(log, "workitem_create", err)
		}
		return toolJSON(out)
	}
}
