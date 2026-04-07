# youtrack-mcp

[![CI](https://github.com/skynet2/youtrack-mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/skynet2/youtrack-mcp/actions/workflows/ci.yml)
[![Release](https://github.com/skynet2/youtrack-mcp/actions/workflows/release.yml/badge.svg)](https://github.com/skynet2/youtrack-mcp/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/skynet2/youtrack-mcp)](https://goreportcard.com/report/github.com/skynet2/youtrack-mcp)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

MCP stdio server for [JetBrains YouTrack](https://www.jetbrains.com/youtrack/) issue tracking.

## Install

Download a pre-built binary from [Releases](https://github.com/skynet2/youtrack-mcp/releases), or build from source:

    go install github.com/skynet2/youtrack-mcp/cmd/youtrack-mcp@latest

## Configure

Copy `configs/config.example.yaml` to `config.yaml` and fill in credentials,
or set env vars (prefix `YOUTRACK_`):

| Variable | Description |
|----------|-------------|
| `YOUTRACK_URL` | YouTrack instance URL |
| `YOUTRACK_TOKEN` | Permanent token (Hub → Profile → Authentication) |
| `YOUTRACK_TIMEOUT` | HTTP timeout (default `30s`) |
| `YOUTRACK_LOG_LEVEL` | Log level: `debug`, `info`, `warn`, `error` (default `info`) |

## Usage

### Claude Code

Add to `~/.claude.json`:

```json
{
  "mcpServers": {
    "youtrack": {
      "type": "stdio",
      "command": "/path/to/youtrack-mcp",
      "args": [],
      "env": {
        "YOUTRACK_URL": "https://youtrack.example.com",
        "YOUTRACK_TOKEN": "perm:your-permanent-token"
      }
    }
  }
}
```

### Other MCP clients

Any MCP-compatible client can launch the binary as a stdio server:

    /path/to/youtrack-mcp

## Tools

### Issues

| Tool | Description |
|------|-------------|
| `issue_search` | Search issues using YouTrack query syntax |
| `issue_get` | Get a single issue by readable ID (e.g. `DEMO-123`) |
| `issue_create` | Create a new issue (supports custom fields with `$type` discriminator) |
| `issue_update` | Update an existing issue |
| `issue_delete` | Delete an issue |
| `issue_add_tag` | Add a tag to an issue |
| `issue_remove_tag` | Remove a tag from an issue |
| `issue_get_links` | Get issue links (related, blocks, duplicates, etc.) |

### Projects

| Tool | Description |
|------|-------------|
| `project_list` | List all projects |
| `project_get` | Get a project by ID or short name |

### Comments

| Tool | Description |
|------|-------------|
| `comment_list` | List comments on an issue |
| `comment_create` | Add a comment to an issue |
| `comment_update` | Update a comment |
| `comment_delete` | Delete a comment |

### Tags

| Tool | Description |
|------|-------------|
| `tag_list` | List all issue tags |

### Time Tracking

| Tool | Description |
|------|-------------|
| `workitem_list` | List time tracking work items for an issue |
| `workitem_create` | Add a work item to an issue |

### Agile Boards

| Tool | Description |
|------|-------------|
| `agile_list` | List agile boards |
| `sprint_list` | List sprints for a board |
| `sprint_get` | Get sprint details including issues |

### Users

| Tool | Description |
|------|-------------|
| `user_list` | List or search users |
| `user_current` | Get the currently authenticated user |

### Custom Fields

| Tool | Description |
|------|-------------|
| `issue_custom_fields` | Get custom field values for an issue |
| `issue_update_field` | Update a custom field on an issue |

### Knowledge Base

| Tool | Description |
|------|-------------|
| `article_list` | List knowledge base articles |
| `article_get` | Get an article by ID |
| `article_create` | Create a new article |
| `article_update` | Update an article |
| `article_delete` | Delete an article |

### Activities

| Tool | Description |
|------|-------------|
| `activity_list` | List global activities filtered by category |
| `issue_activity_list` | List activity history for a specific issue |

## Dev

    make build     # compile binary
    make test      # run unit tests
    make lint      # golangci-lint
    make generate  # regenerate mocks

## License

[MIT](LICENSE)
