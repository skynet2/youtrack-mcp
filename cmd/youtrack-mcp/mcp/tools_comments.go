package mcp

import (
	"context"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"

	"github.com/skynet2/youtrack-mcp/pkg/youtrack"
)

func registerCommentTools(s *server.MCPServer, c YouTrackClient, log zerolog.Logger) {
	s.AddTool(mcplib.NewTool("comment_list",
		mcplib.WithDescription("List comments on a YouTrack issue."),
		mcplib.WithString("issue_id", mcplib.Required(), mcplib.Description("issue readable ID")),
		mcplib.WithNumber("skip", mcplib.Description("pagination offset")),
		mcplib.WithNumber("top", mcplib.Description("max results")),
	), server.ToolHandlerFunc(commentListHandler(c, log)))

	s.AddTool(mcplib.NewTool("comment_create",
		mcplib.WithDescription("Add a comment to a YouTrack issue."),
		mcplib.WithString("issue_id", mcplib.Required(), mcplib.Description("issue readable ID")),
		mcplib.WithString("text", mcplib.Required(), mcplib.Description("comment text (markdown)")),
	), server.ToolHandlerFunc(commentCreateHandler(c, log)))

	s.AddTool(mcplib.NewTool("comment_update",
		mcplib.WithDescription("Update a comment on a YouTrack issue."),
		mcplib.WithString("issue_id", mcplib.Required(), mcplib.Description("issue readable ID")),
		mcplib.WithString("comment_id", mcplib.Required(), mcplib.Description("comment ID")),
		mcplib.WithString("text", mcplib.Required(), mcplib.Description("updated comment text")),
	), server.ToolHandlerFunc(commentUpdateHandler(c, log)))

	s.AddTool(mcplib.NewTool("comment_delete",
		mcplib.WithDescription("Delete a comment from a YouTrack issue."),
		mcplib.WithString("issue_id", mcplib.Required(), mcplib.Description("issue readable ID")),
		mcplib.WithString("comment_id", mcplib.Required(), mcplib.Description("comment ID")),
	), server.ToolHandlerFunc(commentDeleteHandler(c, log)))
}

func commentListHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		issueID, err := reqString(req, "issue_id")
		if err != nil {
			return toolErr(log, "comment_list", err)
		}
		out, err := c.CommentList(ctx, issueID, optInt(req, "skip"), optInt(req, "top"))
		if err != nil {
			return toolErr(log, "comment_list", err)
		}
		return toolJSON(out)
	}
}

func commentCreateHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		issueID, err := reqString(req, "issue_id")
		if err != nil {
			return toolErr(log, "comment_create", err)
		}
		text, err := reqString(req, "text")
		if err != nil {
			return toolErr(log, "comment_create", err)
		}
		out, err := c.CommentCreate(ctx, issueID, youtrack.IssueComment{Text: text})
		if err != nil {
			return toolErr(log, "comment_create", err)
		}
		return toolJSON(out)
	}
}

func commentUpdateHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		issueID, err := reqString(req, "issue_id")
		if err != nil {
			return toolErr(log, "comment_update", err)
		}
		commentID, err := reqString(req, "comment_id")
		if err != nil {
			return toolErr(log, "comment_update", err)
		}
		text, err := reqString(req, "text")
		if err != nil {
			return toolErr(log, "comment_update", err)
		}
		out, err := c.CommentUpdate(ctx, issueID, commentID, youtrack.IssueComment{Text: text})
		if err != nil {
			return toolErr(log, "comment_update", err)
		}
		return toolJSON(out)
	}
}

func commentDeleteHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		issueID, err := reqString(req, "issue_id")
		if err != nil {
			return toolErr(log, "comment_delete", err)
		}
		commentID, err := reqString(req, "comment_id")
		if err != nil {
			return toolErr(log, "comment_delete", err)
		}
		if err := c.CommentDelete(ctx, issueID, commentID); err != nil {
			return toolErr(log, "comment_delete", err)
		}
		return mcplib.NewToolResultText("ok"), nil
	}
}
