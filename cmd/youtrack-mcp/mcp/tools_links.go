package mcp

import (
	"context"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"
)

func registerLinkTools(s *server.MCPServer, c YouTrackClient, log zerolog.Logger) {
	s.AddTool(mcplib.NewTool("issue_link_types",
		mcplib.WithDescription("List available YouTrack issue link types and their directed phrases (use these phrases as link_type in issue_link/issue_unlink)."),
	), server.ToolHandlerFunc(issueLinkTypesHandler(c, log)))

	s.AddTool(mcplib.NewTool("issue_link",
		mcplib.WithDescription("Create a link between two issues. link_type is a directed phrase from issue_link_types (the sourceToTarget/targetToSource phrase, NOT the type name): e.g. 'subtask of' makes source a child of target, 'parent for' makes source the parent, 'relates to' for a plain relation. Example: issue_link(source='PROJ-2', link_type='subtask of', target='PROJ-1')."),
		mcplib.WithString("source", mcplib.Required(), mcplib.Description("readable ID the link is anchored on (e.g. 'PROJ-2')")),
		mcplib.WithString("link_type", mcplib.Required(), mcplib.Description("directed phrase, e.g. 'subtask of', 'parent for', 'relates to', 'depends on'")),
		mcplib.WithString("target", mcplib.Required(), mcplib.Description("readable ID of the other issue (e.g. 'PROJ-1')")),
	), server.ToolHandlerFunc(issueLinkHandler(c, log)))

	s.AddTool(mcplib.NewTool("issue_unlink",
		mcplib.WithDescription("Remove a link between two issues. Same arguments as issue_link."),
		mcplib.WithString("source", mcplib.Required(), mcplib.Description("readable ID the link is anchored on")),
		mcplib.WithString("link_type", mcplib.Required(), mcplib.Description("directed phrase, e.g. 'subtask of'")),
		mcplib.WithString("target", mcplib.Required(), mcplib.Description("readable ID of the other issue")),
	), server.ToolHandlerFunc(issueUnlinkHandler(c, log)))
}

func issueLinkTypesHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		out, err := c.ListIssueLinkTypes(ctx)
		if err != nil {
			return toolErr(log, "issue_link_types", err)
		}
		return toolJSON(out)
	}
}

func issueLinkHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		source, err := reqString(req, "source")
		if err != nil {
			return toolErr(log, "issue_link", err)
		}
		linkType, err := reqString(req, "link_type")
		if err != nil {
			return toolErr(log, "issue_link", err)
		}
		target, err := reqString(req, "target")
		if err != nil {
			return toolErr(log, "issue_link", err)
		}
		if err := c.IssueLink(ctx, source, linkType, target); err != nil {
			return toolErr(log, "issue_link", err)
		}
		return toolJSON(map[string]string{"status": "linked", "source": source, "link_type": linkType, "target": target})
	}
}

func issueUnlinkHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		source, err := reqString(req, "source")
		if err != nil {
			return toolErr(log, "issue_unlink", err)
		}
		linkType, err := reqString(req, "link_type")
		if err != nil {
			return toolErr(log, "issue_unlink", err)
		}
		target, err := reqString(req, "target")
		if err != nil {
			return toolErr(log, "issue_unlink", err)
		}
		if err := c.IssueUnlink(ctx, source, linkType, target); err != nil {
			return toolErr(log, "issue_unlink", err)
		}
		return toolJSON(map[string]string{"status": "unlinked", "source": source, "link_type": linkType, "target": target})
	}
}
