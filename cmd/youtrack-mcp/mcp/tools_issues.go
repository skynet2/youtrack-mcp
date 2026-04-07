package mcp

import (
	"context"
	"encoding/json"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"

	"github.com/skynet2/youtrack-mcp/pkg/youtrack"
)

func registerIssueTools(s *server.MCPServer, c YouTrackClient, log zerolog.Logger) {
	s.AddTool(mcplib.NewTool("issue_search",
		mcplib.WithDescription("Search YouTrack issues using query syntax (e.g. 'project: DEMO State: Open')."),
		mcplib.WithString("query", mcplib.Description("YouTrack search query")),
		mcplib.WithNumber("skip", mcplib.Description("pagination offset")),
		mcplib.WithNumber("top", mcplib.Description("max results to return")),
	), server.ToolHandlerFunc(issueSearchHandler(c, log)))

	s.AddTool(mcplib.NewTool("issue_get",
		mcplib.WithDescription("Get a single YouTrack issue by readable ID (e.g. 'DEMO-123')."),
		mcplib.WithString("id", mcplib.Required(), mcplib.Description("issue readable ID")),
	), server.ToolHandlerFunc(issueGetHandler(c, log)))

	s.AddTool(mcplib.NewTool("issue_create",
		mcplib.WithDescription("Create a new YouTrack issue. Use issue_custom_fields on an existing issue in the same project to discover field names and $type values. Custom fields must include $type discriminator, e.g. [{\"$type\":\"SingleEnumIssueCustomField\",\"name\":\"Type\",\"value\":{\"name\":\"Task\"}}]. Common $types: SingleEnumIssueCustomField, StateIssueCustomField, MultiVersionIssueCustomField, SingleUserIssueCustomField, PeriodIssueCustomField, SimpleIssueCustomField, DateIssueCustomField, TextIssueCustomField."),
		mcplib.WithString("project_id", mcplib.Required(), mcplib.Description("project internal ID (e.g. '0-2'). Use project_list to find IDs.")),
		mcplib.WithString("summary", mcplib.Required(), mcplib.Description("issue title")),
		mcplib.WithString("description", mcplib.Description("issue body (markdown)")),
		mcplib.WithArray("custom_fields", mcplib.Description("array of custom field objects with $type, name, and value")),
	), server.ToolHandlerFunc(issueCreateHandler(c, log)))

	s.AddTool(mcplib.NewTool("issue_update",
		mcplib.WithDescription("Update an existing YouTrack issue."),
		mcplib.WithString("id", mcplib.Required(), mcplib.Description("issue readable ID")),
		mcplib.WithString("summary", mcplib.Description("new summary")),
		mcplib.WithString("description", mcplib.Description("new description")),
	), server.ToolHandlerFunc(issueUpdateHandler(c, log)))

	s.AddTool(mcplib.NewTool("issue_delete",
		mcplib.WithDescription("Delete a YouTrack issue."),
		mcplib.WithString("id", mcplib.Required(), mcplib.Description("issue readable ID")),
	), server.ToolHandlerFunc(issueDeleteHandler(c, log)))

	s.AddTool(mcplib.NewTool("issue_add_tag",
		mcplib.WithDescription("Add a tag to a YouTrack issue."),
		mcplib.WithString("issue_id", mcplib.Required(), mcplib.Description("issue readable ID")),
		mcplib.WithString("tag_id", mcplib.Required(), mcplib.Description("tag ID")),
	), server.ToolHandlerFunc(issueAddTagHandler(c, log)))

	s.AddTool(mcplib.NewTool("issue_remove_tag",
		mcplib.WithDescription("Remove a tag from a YouTrack issue."),
		mcplib.WithString("issue_id", mcplib.Required(), mcplib.Description("issue readable ID")),
		mcplib.WithString("tag_id", mcplib.Required(), mcplib.Description("tag ID")),
	), server.ToolHandlerFunc(issueRemoveTagHandler(c, log)))

	s.AddTool(mcplib.NewTool("issue_get_links",
		mcplib.WithDescription("Get links for a YouTrack issue (related, blocks, duplicates, etc.)."),
		mcplib.WithString("issue_id", mcplib.Required(), mcplib.Description("issue readable ID")),
	), server.ToolHandlerFunc(issueGetLinksHandler(c, log)))
}

func issueSearchHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		query := optString(req, "query")
		skip := optInt(req, "skip")
		top := optInt(req, "top")
		out, err := c.IssueSearch(ctx, query, skip, top)
		if err != nil {
			return toolErr(log, "issue_search", err)
		}
		return toolJSON(out)
	}
}

func issueGetHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		id, err := reqString(req, "id")
		if err != nil {
			return toolErr(log, "issue_get", err)
		}
		out, err := c.IssueGet(ctx, id)
		if err != nil {
			return toolErr(log, "issue_get", err)
		}
		return toolJSON(out)
	}
}

func issueCreateHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		projectID, err := reqString(req, "project_id")
		if err != nil {
			return toolErr(log, "issue_create", err)
		}
		summary, err := reqString(req, "summary")
		if err != nil {
			return toolErr(log, "issue_create", err)
		}
		issue := youtrack.Issue{
			Project: &youtrack.Project{ID: projectID},
			Summary: summary,
		}
		if desc := optString(req, "description"); desc != "" {
			issue.Description = desc
		}
		if raw, ok := req.GetArguments()["custom_fields"]; ok {
			b, jErr := json.Marshal(raw)
			if jErr != nil {
				return toolErr(log, "issue_create", jErr)
			}
			if jErr = json.Unmarshal(b, &issue.CustomFields); jErr != nil {
				return toolErr(log, "issue_create", jErr)
			}
		}
		out, err := c.IssueCreate(ctx, issue)
		if err != nil {
			return toolErr(log, "issue_create", err)
		}
		return toolJSON(out)
	}
}

func issueUpdateHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		id, err := reqString(req, "id")
		if err != nil {
			return toolErr(log, "issue_update", err)
		}
		var issue youtrack.Issue
		if s := optString(req, "summary"); s != "" {
			issue.Summary = s
		}
		if d := optString(req, "description"); d != "" {
			issue.Description = d
		}
		out, err := c.IssueUpdate(ctx, id, issue)
		if err != nil {
			return toolErr(log, "issue_update", err)
		}
		return toolJSON(out)
	}
}

func issueDeleteHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		id, err := reqString(req, "id")
		if err != nil {
			return toolErr(log, "issue_delete", err)
		}
		if err := c.IssueDelete(ctx, id); err != nil {
			return toolErr(log, "issue_delete", err)
		}
		return mcplib.NewToolResultText("ok"), nil
	}
}

func issueAddTagHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		issueID, err := reqString(req, "issue_id")
		if err != nil {
			return toolErr(log, "issue_add_tag", err)
		}
		tagID, err := reqString(req, "tag_id")
		if err != nil {
			return toolErr(log, "issue_add_tag", err)
		}
		if err := c.IssueAddTag(ctx, issueID, tagID); err != nil {
			return toolErr(log, "issue_add_tag", err)
		}
		return mcplib.NewToolResultText("ok"), nil
	}
}

func issueRemoveTagHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		issueID, err := reqString(req, "issue_id")
		if err != nil {
			return toolErr(log, "issue_remove_tag", err)
		}
		tagID, err := reqString(req, "tag_id")
		if err != nil {
			return toolErr(log, "issue_remove_tag", err)
		}
		if err := c.IssueRemoveTag(ctx, issueID, tagID); err != nil {
			return toolErr(log, "issue_remove_tag", err)
		}
		return mcplib.NewToolResultText("ok"), nil
	}
}

func issueGetLinksHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		issueID, err := reqString(req, "issue_id")
		if err != nil {
			return toolErr(log, "issue_get_links", err)
		}
		out, err := c.IssueGetLinks(ctx, issueID)
		if err != nil {
			return toolErr(log, "issue_get_links", err)
		}
		return toolJSON(out)
	}
}
