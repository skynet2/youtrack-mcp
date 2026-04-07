package mcp

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	mcplib "github.com/mark3labs/mcp-go/mcp"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/skynet2/youtrack-mcp/cmd/youtrack-mcp/mcp/mocks"
	"github.com/skynet2/youtrack-mcp/pkg/youtrack"
)

func newReq(args map[string]any) mcplib.CallToolRequest {
	return mcplib.CallToolRequest{
		Params: mcplib.CallToolParams{
			Arguments: args,
		},
	}
}

func TestIssueSearch_Success(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]any
		query    string
		skip     int
		top      int
		response []youtrack.Issue
	}{
		{
			name:  "with all parameters",
			args:  map[string]any{"query": "project: TST", "skip": float64(5), "top": float64(10)},
			query: "project: TST",
			skip:  5,
			top:   10,
			response: []youtrack.Issue{
				{IDReadable: "TST-1", Summary: "First issue"},
				{IDReadable: "TST-2", Summary: "Second issue"},
			},
		},
		{
			name:     "with no parameters",
			args:     map[string]any{},
			query:    "",
			skip:     0,
			top:      0,
			response: []youtrack.Issue{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockYouTrackClient(ctrl)
			log := zerolog.Nop()

			mock.EXPECT().IssueSearch(gomock.Any(), tt.query, tt.skip, tt.top).
				Return(tt.response, nil)

			handler := issueSearchHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.False(t, result.IsError)
			require.Len(t, result.Content, 1)
			text := result.Content[0].(mcplib.TextContent).Text

			for _, issue := range tt.response {
				assert.Contains(t, text, issue.IDReadable)
			}
		})
	}
}

func TestIssueSearch_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().IssueSearch(gomock.Any(), "", 0, 0).
		Return(nil, errors.New("connection refused"))

	handler := issueSearchHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{}))

	require.NoError(t, err)
	require.True(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "connection refused")
}

func TestIssueGet_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().IssueGet(gomock.Any(), "TST-42").
		Return(youtrack.Issue{IDReadable: "TST-42", Summary: "Test issue"}, nil)

	handler := issueGetHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{"id": "TST-42"}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "TST-42")
	assert.Contains(t, text, "Test issue")
}

func TestIssueGet_Failure(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]any
		setupMock   func(m *mocks.MockYouTrackClient)
		errContains string
	}{
		{
			name: "missing id",
			args: map[string]any{},
			setupMock: func(m *mocks.MockYouTrackClient) {
			},
			errContains: "id",
		},
		{
			name: "client error",
			args: map[string]any{"id": "TST-1"},
			setupMock: func(m *mocks.MockYouTrackClient) {
				m.EXPECT().IssueGet(gomock.Any(), "TST-1").
					Return(youtrack.Issue{}, errors.New("not found"))
			},
			errContains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockYouTrackClient(ctrl)
			log := zerolog.Nop()

			tt.setupMock(mock)

			handler := issueGetHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.True(t, result.IsError)
			text := result.Content[0].(mcplib.TextContent).Text
			assert.Contains(t, text, tt.errContains)
		})
	}
}

func TestIssueCreate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	expectedIssue := youtrack.Issue{
		Project: &youtrack.Project{ID: "proj-1"},
		Summary: "New issue",
	}
	mock.EXPECT().IssueCreate(gomock.Any(), expectedIssue).
		Return(youtrack.Issue{IDReadable: "TST-99", Summary: "New issue"}, nil)

	handler := issueCreateHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{
		"project_id": "proj-1",
		"summary":    "New issue",
	}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "TST-99")
}

func TestIssueCreate_SuccessWithDescription(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	expectedIssue := youtrack.Issue{
		Project:     &youtrack.Project{ID: "proj-1"},
		Summary:     "New issue",
		Description: "Some details",
	}
	mock.EXPECT().IssueCreate(gomock.Any(), expectedIssue).
		Return(youtrack.Issue{IDReadable: "TST-100", Summary: "New issue"}, nil)

	handler := issueCreateHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{
		"project_id":  "proj-1",
		"summary":     "New issue",
		"description": "Some details",
	}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "TST-100")
}

func TestIssueCreate_Failure(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]any
		setupMock   func(m *mocks.MockYouTrackClient)
		errContains string
	}{
		{
			name:        "missing project_id",
			args:        map[string]any{"summary": "test"},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "project_id",
		},
		{
			name:        "missing summary",
			args:        map[string]any{"project_id": "proj-1"},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "summary",
		},
		{
			name: "client error",
			args: map[string]any{"project_id": "proj-1", "summary": "test"},
			setupMock: func(m *mocks.MockYouTrackClient) {
				m.EXPECT().IssueCreate(gomock.Any(), gomock.Any()).
					Return(youtrack.Issue{}, errors.New("permission denied"))
			},
			errContains: "permission denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockYouTrackClient(ctrl)
			log := zerolog.Nop()

			tt.setupMock(mock)

			handler := issueCreateHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.True(t, result.IsError)
			text := result.Content[0].(mcplib.TextContent).Text
			assert.Contains(t, text, tt.errContains)
		})
	}
}

func TestIssueUpdate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().IssueUpdate(gomock.Any(), "TST-1", youtrack.Issue{Summary: "Updated"}).
		Return(youtrack.Issue{IDReadable: "TST-1", Summary: "Updated"}, nil)

	handler := issueUpdateHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{
		"id":      "TST-1",
		"summary": "Updated",
	}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "Updated")
}

func TestIssueUpdate_Failure(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]any
		setupMock   func(m *mocks.MockYouTrackClient)
		errContains string
	}{
		{
			name:        "missing id",
			args:        map[string]any{"summary": "test"},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "id",
		},
		{
			name: "client error",
			args: map[string]any{"id": "TST-1", "summary": "test"},
			setupMock: func(m *mocks.MockYouTrackClient) {
				m.EXPECT().IssueUpdate(gomock.Any(), "TST-1", gomock.Any()).
					Return(youtrack.Issue{}, errors.New("server error"))
			},
			errContains: "server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockYouTrackClient(ctrl)
			log := zerolog.Nop()

			tt.setupMock(mock)

			handler := issueUpdateHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.True(t, result.IsError)
			text := result.Content[0].(mcplib.TextContent).Text
			assert.Contains(t, text, tt.errContains)
		})
	}
}

func TestIssueDelete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().IssueDelete(gomock.Any(), "TST-1").Return(nil)

	handler := issueDeleteHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{"id": "TST-1"}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Equal(t, "ok", text)
}

func TestIssueDelete_Failure(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]any
		setupMock   func(m *mocks.MockYouTrackClient)
		errContains string
	}{
		{
			name:        "missing id",
			args:        map[string]any{},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "id",
		},
		{
			name: "client error",
			args: map[string]any{"id": "TST-1"},
			setupMock: func(m *mocks.MockYouTrackClient) {
				m.EXPECT().IssueDelete(gomock.Any(), "TST-1").
					Return(errors.New("forbidden"))
			},
			errContains: "forbidden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockYouTrackClient(ctrl)
			log := zerolog.Nop()

			tt.setupMock(mock)

			handler := issueDeleteHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.True(t, result.IsError)
			text := result.Content[0].(mcplib.TextContent).Text
			assert.Contains(t, text, tt.errContains)
		})
	}
}

func TestIssueAddTag_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().IssueAddTag(gomock.Any(), "TST-1", "tag-123").Return(nil)

	handler := issueAddTagHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{
		"issue_id": "TST-1",
		"tag_id":   "tag-123",
	}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Equal(t, "ok", text)
}

func TestIssueAddTag_Failure(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]any
		setupMock   func(m *mocks.MockYouTrackClient)
		errContains string
	}{
		{
			name:        "missing issue_id",
			args:        map[string]any{"tag_id": "tag-1"},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "issue_id",
		},
		{
			name:        "missing tag_id",
			args:        map[string]any{"issue_id": "TST-1"},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "tag_id",
		},
		{
			name: "client error",
			args: map[string]any{"issue_id": "TST-1", "tag_id": "tag-1"},
			setupMock: func(m *mocks.MockYouTrackClient) {
				m.EXPECT().IssueAddTag(gomock.Any(), "TST-1", "tag-1").
					Return(errors.New("tag not found"))
			},
			errContains: "tag not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockYouTrackClient(ctrl)
			log := zerolog.Nop()

			tt.setupMock(mock)

			handler := issueAddTagHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.True(t, result.IsError)
			text := result.Content[0].(mcplib.TextContent).Text
			assert.Contains(t, text, tt.errContains)
		})
	}
}

func TestIssueRemoveTag_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().IssueRemoveTag(gomock.Any(), "TST-1", "tag-456").Return(nil)

	handler := issueRemoveTagHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{
		"issue_id": "TST-1",
		"tag_id":   "tag-456",
	}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Equal(t, "ok", text)
}

func TestIssueRemoveTag_Failure(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]any
		setupMock   func(m *mocks.MockYouTrackClient)
		errContains string
	}{
		{
			name:        "missing issue_id",
			args:        map[string]any{"tag_id": "tag-1"},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "issue_id",
		},
		{
			name:        "missing tag_id",
			args:        map[string]any{"issue_id": "TST-1"},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "tag_id",
		},
		{
			name: "client error",
			args: map[string]any{"issue_id": "TST-1", "tag_id": "tag-1"},
			setupMock: func(m *mocks.MockYouTrackClient) {
				m.EXPECT().IssueRemoveTag(gomock.Any(), "TST-1", "tag-1").
					Return(errors.New("conflict"))
			},
			errContains: "conflict",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockYouTrackClient(ctrl)
			log := zerolog.Nop()

			tt.setupMock(mock)

			handler := issueRemoveTagHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.True(t, result.IsError)
			text := result.Content[0].(mcplib.TextContent).Text
			assert.Contains(t, text, tt.errContains)
		})
	}
}

func TestIssueGetLinks_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().IssueGetLinks(gomock.Any(), "TST-1").
		Return([]youtrack.IssueLink{
			{ID: "link-1", Direction: "OUTWARD"},
		}, nil)

	handler := issueGetLinksHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{"issue_id": "TST-1"}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "link-1")
}

func TestIssueGetLinks_Failure(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]any
		setupMock   func(m *mocks.MockYouTrackClient)
		errContains string
	}{
		{
			name:        "missing issue_id",
			args:        map[string]any{},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "issue_id",
		},
		{
			name: "client error",
			args: map[string]any{"issue_id": "TST-1"},
			setupMock: func(m *mocks.MockYouTrackClient) {
				m.EXPECT().IssueGetLinks(gomock.Any(), "TST-1").
					Return(nil, errors.New("timeout"))
			},
			errContains: "timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockYouTrackClient(ctrl)
			log := zerolog.Nop()

			tt.setupMock(mock)

			handler := issueGetLinksHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.True(t, result.IsError)
			text := result.Content[0].(mcplib.TextContent).Text
			assert.Contains(t, text, tt.errContains)
		})
	}
}
