package mcp

import (
	"context"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"

	"github.com/skynet2/youtrack-mcp/pkg/youtrack"
)

func registerKnowledgeBaseTools(s *server.MCPServer, c YouTrackClient, log zerolog.Logger) {
	s.AddTool(mcplib.NewTool("article_list",
		mcplib.WithDescription("List YouTrack knowledge base articles."),
		mcplib.WithNumber("skip", mcplib.Description("pagination offset")),
		mcplib.WithNumber("top", mcplib.Description("max results")),
	), server.ToolHandlerFunc(articleListHandler(c, log)))

	s.AddTool(mcplib.NewTool("article_get",
		mcplib.WithDescription("Get a YouTrack knowledge base article by ID."),
		mcplib.WithString("id", mcplib.Required(), mcplib.Description("article ID")),
	), server.ToolHandlerFunc(articleGetHandler(c, log)))

	s.AddTool(mcplib.NewTool("article_create",
		mcplib.WithDescription("Create a YouTrack knowledge base article."),
		mcplib.WithString("project_id", mcplib.Required(), mcplib.Description("project ID")),
		mcplib.WithString("summary", mcplib.Required(), mcplib.Description("article title")),
		mcplib.WithString("content", mcplib.Description("article body (markdown)")),
		mcplib.WithString("parent_id", mcplib.Description("parent article ID for nesting")),
	), server.ToolHandlerFunc(articleCreateHandler(c, log)))

	s.AddTool(mcplib.NewTool("article_update",
		mcplib.WithDescription("Update a YouTrack knowledge base article."),
		mcplib.WithString("id", mcplib.Required(), mcplib.Description("article ID")),
		mcplib.WithString("summary", mcplib.Description("new title")),
		mcplib.WithString("content", mcplib.Description("new body")),
	), server.ToolHandlerFunc(articleUpdateHandler(c, log)))

	s.AddTool(mcplib.NewTool("article_delete",
		mcplib.WithDescription("Delete a YouTrack knowledge base article."),
		mcplib.WithString("id", mcplib.Required(), mcplib.Description("article ID")),
	), server.ToolHandlerFunc(articleDeleteHandler(c, log)))
}

func articleListHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		out, err := c.ArticleList(ctx, optInt(req, "skip"), optInt(req, "top"))
		if err != nil {
			return toolErr(log, "article_list", err)
		}
		return toolJSON(out)
	}
}

func articleGetHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		id, err := reqString(req, "id")
		if err != nil {
			return toolErr(log, "article_get", err)
		}
		out, err := c.ArticleGet(ctx, id)
		if err != nil {
			return toolErr(log, "article_get", err)
		}
		return toolJSON(out)
	}
}

func articleCreateHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		projectID, err := reqString(req, "project_id")
		if err != nil {
			return toolErr(log, "article_create", err)
		}
		summary, err := reqString(req, "summary")
		if err != nil {
			return toolErr(log, "article_create", err)
		}
		article := youtrack.Article{
			Project: &youtrack.Project{ID: projectID},
			Summary: summary,
		}
		if c := optString(req, "content"); c != "" {
			article.Content = c
		}
		if pid := optString(req, "parent_id"); pid != "" {
			article.ParentArticle = &youtrack.Article{ID: pid}
		}
		out, err := c.ArticleCreate(ctx, article)
		if err != nil {
			return toolErr(log, "article_create", err)
		}
		return toolJSON(out)
	}
}

func articleUpdateHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		id, err := reqString(req, "id")
		if err != nil {
			return toolErr(log, "article_update", err)
		}
		var article youtrack.Article
		if s := optString(req, "summary"); s != "" {
			article.Summary = s
		}
		if ct := optString(req, "content"); ct != "" {
			article.Content = ct
		}
		out, err := c.ArticleUpdate(ctx, id, article)
		if err != nil {
			return toolErr(log, "article_update", err)
		}
		return toolJSON(out)
	}
}

func articleDeleteHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		id, err := reqString(req, "id")
		if err != nil {
			return toolErr(log, "article_delete", err)
		}
		if err := c.ArticleDelete(ctx, id); err != nil {
			return toolErr(log, "article_delete", err)
		}
		return mcplib.NewToolResultText("ok"), nil
	}
}
