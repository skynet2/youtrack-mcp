package mcp

import (
	"context"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"
)

func registerActivityTools(s *server.MCPServer, c YouTrackClient, log zerolog.Logger) {
	s.AddTool(mcplib.NewTool("activity_list",
		mcplib.WithDescription("List YouTrack activities (global change log). Filter by categories, issue query, or time range."),
		mcplib.WithString("categories", mcplib.Required(), mcplib.Description("comma-separated activity categories. Common: IssueCreatedCategory, IssueResolvedCategory, CommentsCategory, CustomFieldCategory, DescriptionCategory, SummaryCategory, LinksCategory, TagsCategory, AttachmentsCategory, SprintCategory, VcsChangeCategory")),
		mcplib.WithString("issue_query", mcplib.Description("YouTrack issue search query to filter activities")),
		mcplib.WithString("start", mcplib.Description("start timestamp (ms)")),
		mcplib.WithString("end", mcplib.Description("end timestamp (ms)")),
		mcplib.WithNumber("skip", mcplib.Description("pagination offset")),
		mcplib.WithNumber("top", mcplib.Description("max results")),
	), server.ToolHandlerFunc(activityListHandler(c, log)))

	s.AddTool(mcplib.NewTool("issue_activity_list",
		mcplib.WithDescription("List activities (change history) for a specific YouTrack issue."),
		mcplib.WithString("issue_id", mcplib.Required(), mcplib.Description("issue readable ID")),
		mcplib.WithString("categories", mcplib.Description("comma-separated categories")),
		mcplib.WithNumber("skip", mcplib.Description("pagination offset")),
		mcplib.WithNumber("top", mcplib.Description("max results")),
	), server.ToolHandlerFunc(issueActivityListHandler(c, log)))
}

func activityListHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		out, err := c.ActivityList(ctx,
			optString(req, "categories"),
			optString(req, "issue_query"),
			optString(req, "start"),
			optString(req, "end"),
			optInt(req, "skip"),
			optInt(req, "top"),
		)
		if err != nil {
			return toolErr(log, "activity_list", err)
		}
		return toolJSON(out)
	}
}

func issueActivityListHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		issueID, err := reqString(req, "issue_id")
		if err != nil {
			return toolErr(log, "issue_activity_list", err)
		}
		out, err := c.IssueActivityList(ctx, issueID,
			optString(req, "categories"),
			optInt(req, "skip"),
			optInt(req, "top"),
		)
		if err != nil {
			return toolErr(log, "issue_activity_list", err)
		}
		return toolJSON(out)
	}
}
