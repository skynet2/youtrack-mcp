package mcp

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog"
)

func NewServer(client YouTrackClient, log zerolog.Logger, version string) *server.MCPServer {
	s := server.NewMCPServer("youtrack-mcp", version)
	registerIssueTools(s, client, log)
	registerProjectTools(s, client, log)
	registerCommentTools(s, client, log)
	registerTagTools(s, client, log)
	registerTimeTrackingTools(s, client, log)
	registerAgileTools(s, client, log)
	registerUserTools(s, client, log)
	registerCustomFieldTools(s, client, log)
	registerKnowledgeBaseTools(s, client, log)
	registerActivityTools(s, client, log)
	registerLinkTools(s, client, log)
	return s
}
