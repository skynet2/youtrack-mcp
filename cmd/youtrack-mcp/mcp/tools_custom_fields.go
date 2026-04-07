package mcp

import (
	"context"

	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"

	"github.com/skynet2/youtrack-mcp/pkg/youtrack"
)

func registerCustomFieldTools(s *server.MCPServer, c YouTrackClient, log zerolog.Logger) {
	s.AddTool(mcplib.NewTool("issue_custom_fields",
		mcplib.WithDescription("Get custom field values for a YouTrack issue."),
		mcplib.WithString("issue_id", mcplib.Required(), mcplib.Description("issue readable ID")),
	), server.ToolHandlerFunc(issueCustomFieldsHandler(c, log)))

	s.AddTool(mcplib.NewTool("issue_update_field",
		mcplib.WithDescription("Update a custom field value on a YouTrack issue. The value format depends on field type: for enum/state fields use {\"name\":\"value\"}, for simple fields use the direct value."),
		mcplib.WithString("issue_id", mcplib.Required(), mcplib.Description("issue readable ID")),
		mcplib.WithString("field_id", mcplib.Required(), mcplib.Description("custom field ID")),
		mcplib.WithObject("value", mcplib.Required(), mcplib.Description("new field value")),
	), server.ToolHandlerFunc(issueUpdateFieldHandler(c, log)))
}

func issueCustomFieldsHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		issueID, err := reqString(req, "issue_id")
		if err != nil {
			return toolErr(log, "issue_custom_fields", err)
		}
		out, err := c.IssueCustomFields(ctx, issueID)
		if err != nil {
			return toolErr(log, "issue_custom_fields", err)
		}
		return toolJSON(out)
	}
}

func issueUpdateFieldHandler(c YouTrackClient, log zerolog.Logger) toolHandler {
	return func(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
		issueID, err := reqString(req, "issue_id")
		if err != nil {
			return toolErr(log, "issue_update_field", err)
		}
		fieldID, err := reqString(req, "field_id")
		if err != nil {
			return toolErr(log, "issue_update_field", err)
		}
		value, err := reqMap(req, "value")
		if err != nil {
			return toolErr(log, "issue_update_field", err)
		}
		field := youtrack.IssueCustomField{Value: value}
		out, err := c.IssueUpdateField(ctx, issueID, fieldID, field)
		if err != nil {
			return toolErr(log, "issue_update_field", err)
		}
		return toolJSON(out)
	}
}
