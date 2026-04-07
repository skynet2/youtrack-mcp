# YouTrack MCP Server Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a Go MCP server exposing YouTrack REST API as MCP tools for Claude Code.

**Architecture:** Three-layer design mirroring kiwi-tcms-mcp: CLI (cobra) → MCP tool handlers → YouTrack HTTP client. Config via viper (YAML + env vars). Bearer token auth. All API methods return typed Go structs.

**Tech Stack:** Go 1.24+, mcp-go, cobra, viper, zerolog, testify, gomock, httptest

---

## Tool Inventory (31 tools)

### Issues (8)
- `issue_search` — search with YouTrack query syntax (GET /api/issues?query=...)
- `issue_get` — get single issue by readable ID (GET /api/issues/{id})
- `issue_create` — create issue (POST /api/issues)
- `issue_update` — update issue fields (POST /api/issues/{id})
- `issue_delete` — delete issue (DELETE /api/issues/{id})
- `issue_add_tag` — add tag to issue (POST /api/issues/{id}/tags, body: {id: tagId})
- `issue_remove_tag` — remove tag (DELETE /api/issues/{id}/tags/{tagId})
- `issue_get_links` — get issue links (GET /api/issues/{id}/links)

### Projects (2)
- `project_list` — list projects (GET /api/admin/projects)
- `project_get` — get project (GET /api/admin/projects/{id})

### Comments (4)
- `comment_list` — list comments (GET /api/issues/{id}/comments)
- `comment_create` — add comment (POST /api/issues/{id}/comments)
- `comment_update` — update comment (POST /api/issues/{id}/comments/{commentId})
- `comment_delete` — delete comment (DELETE /api/issues/{id}/comments/{commentId})

### Tags (1)
- `tag_list` — list all tags (GET /api/issueTags)

### Time Tracking (2)
- `workitem_list` — list work items (GET /api/issues/{id}/timeTracking/workItems)
- `workitem_create` — add work item (POST /api/issues/{id}/timeTracking/workItems)

### Agile (3)
- `agile_list` — list boards (GET /api/agiles)
- `sprint_list` — list sprints (GET /api/agiles/{id}/sprints)
- `sprint_get` — get sprint (GET /api/agiles/{id}/sprints/{sprintId})

### Users (2)
- `user_list` — list/search users (GET /api/users)
- `user_current` — current user (GET /api/users/me)

### Custom Fields (2)
- `issue_custom_fields` — get issue custom fields (GET /api/issues/{id}/customFields)
- `issue_update_field` — update a custom field (POST /api/issues/{id}/customFields/{fieldId})

### Knowledge Base (5)
- `article_list` — list articles (GET /api/articles)
- `article_get` — get article (GET /api/articles/{id})
- `article_create` — create article (POST /api/articles)
- `article_update` — update article (POST /api/articles/{id})
- `article_delete` — delete article (DELETE /api/articles/{id})

### Activities (2)
- `activity_list` — global activities (GET /api/activities)
- `issue_activity_list` — issue activities (GET /api/issues/{id}/activities)

---

## Project Structure

```
youtrack-mcp/
├── cmd/youtrack-mcp/
│   ├── main.go
│   ├── cli/
│   │   ├── root.go
│   │   └── serve.go
│   └── mcp/
│       ├── server.go
│       ├── client.go              # YouTrackClient interface
│       ├── helpers.go             # toolErr, toolJSON, reqString, etc.
│       ├── tools_issues.go
│       ├── tools_projects.go
│       ├── tools_comments.go
│       ├── tools_tags.go
│       ├── tools_time_tracking.go
│       ├── tools_agile.go
│       ├── tools_users.go
│       ├── tools_custom_fields.go
│       ├── tools_knowledge_base.go
│       ├── tools_activities.go
│       └── mocks/                 # generated
├── pkg/youtrack/
│   ├── client.go
│   ├── types.go
│   ├── errors.go
│   ├── issues.go
│   ├── projects.go
│   ├── comments.go
│   ├── tags.go
│   ├── time_tracking.go
│   ├── agile.go
│   ├── users.go
│   ├── custom_fields.go
│   ├── knowledge_base.go
│   └── activities.go
├── internal/config/
│   └── config.go
├── configs/
│   └── config.example.yaml
├── Makefile
├── .golangci.yml
└── .goreleaser.yml
```

## Config

```yaml
url: https://youtrack.example.com    # YouTrack base URL
token: perm:xxx                       # Permanent token
timeout: 30s
log_level: info
```

Env prefix: `YOUTRACK_` (e.g. YOUTRACK_URL, YOUTRACK_TOKEN)

## Testing Strategy

- **pkg/youtrack/*_test.go**: httptest.NewServer mocking YouTrack REST responses
- **cmd/mcp/tools_*_test.go**: gomock mocking YouTrackClient interface
- **internal/config/config_test.go**: file + env var config loading
- All tests: table-driven, separate success/failure, no branching
