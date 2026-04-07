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

func TestCommentList_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().CommentList(gomock.Any(), "TST-1", 0, 0).
		Return([]youtrack.IssueComment{
			{ID: "c-1", Text: "First comment"},
			{ID: "c-2", Text: "Second comment"},
		}, nil)

	handler := commentListHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{"issue_id": "TST-1"}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "c-1")
	assert.Contains(t, text, "c-2")
}

func TestCommentList_Failure(t *testing.T) {
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
				m.EXPECT().CommentList(gomock.Any(), "TST-1", 0, 0).
					Return(nil, errors.New("service unavailable"))
			},
			errContains: "service unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockYouTrackClient(ctrl)
			log := zerolog.Nop()

			tt.setupMock(mock)

			handler := commentListHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.True(t, result.IsError)
			text := result.Content[0].(mcplib.TextContent).Text
			assert.Contains(t, text, tt.errContains)
		})
	}
}

func TestCommentCreate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().CommentCreate(gomock.Any(), "TST-1", youtrack.IssueComment{Text: "Hello"}).
		Return(youtrack.IssueComment{ID: "c-new", Text: "Hello"}, nil)

	handler := commentCreateHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{
		"issue_id": "TST-1",
		"text":     "Hello",
	}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "c-new")
}

func TestCommentCreate_Failure(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]any
		setupMock   func(m *mocks.MockYouTrackClient)
		errContains string
	}{
		{
			name:        "missing issue_id",
			args:        map[string]any{"text": "Hello"},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "issue_id",
		},
		{
			name:        "missing text",
			args:        map[string]any{"issue_id": "TST-1"},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "text",
		},
		{
			name: "client error",
			args: map[string]any{"issue_id": "TST-1", "text": "Hello"},
			setupMock: func(m *mocks.MockYouTrackClient) {
				m.EXPECT().CommentCreate(gomock.Any(), "TST-1", gomock.Any()).
					Return(youtrack.IssueComment{}, errors.New("bad request"))
			},
			errContains: "bad request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockYouTrackClient(ctrl)
			log := zerolog.Nop()

			tt.setupMock(mock)

			handler := commentCreateHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.True(t, result.IsError)
			text := result.Content[0].(mcplib.TextContent).Text
			assert.Contains(t, text, tt.errContains)
		})
	}
}

func TestCommentUpdate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().CommentUpdate(gomock.Any(), "TST-1", "c-1", youtrack.IssueComment{Text: "Updated"}).
		Return(youtrack.IssueComment{ID: "c-1", Text: "Updated"}, nil)

	handler := commentUpdateHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{
		"issue_id":   "TST-1",
		"comment_id": "c-1",
		"text":       "Updated",
	}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "c-1")
}

func TestCommentUpdate_Failure(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]any
		setupMock   func(m *mocks.MockYouTrackClient)
		errContains string
	}{
		{
			name:        "missing issue_id",
			args:        map[string]any{"comment_id": "c-1", "text": "x"},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "issue_id",
		},
		{
			name:        "missing comment_id",
			args:        map[string]any{"issue_id": "TST-1", "text": "x"},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "comment_id",
		},
		{
			name:        "missing text",
			args:        map[string]any{"issue_id": "TST-1", "comment_id": "c-1"},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "text",
		},
		{
			name: "client error",
			args: map[string]any{"issue_id": "TST-1", "comment_id": "c-1", "text": "x"},
			setupMock: func(m *mocks.MockYouTrackClient) {
				m.EXPECT().CommentUpdate(gomock.Any(), "TST-1", "c-1", gomock.Any()).
					Return(youtrack.IssueComment{}, errors.New("not found"))
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

			handler := commentUpdateHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.True(t, result.IsError)
			text := result.Content[0].(mcplib.TextContent).Text
			assert.Contains(t, text, tt.errContains)
		})
	}
}

func TestCommentDelete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().CommentDelete(gomock.Any(), "TST-1", "c-1").Return(nil)

	handler := commentDeleteHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{
		"issue_id":   "TST-1",
		"comment_id": "c-1",
	}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Equal(t, "ok", text)
}

func TestCommentDelete_Failure(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]any
		setupMock   func(m *mocks.MockYouTrackClient)
		errContains string
	}{
		{
			name:        "missing issue_id",
			args:        map[string]any{"comment_id": "c-1"},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "issue_id",
		},
		{
			name:        "missing comment_id",
			args:        map[string]any{"issue_id": "TST-1"},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "comment_id",
		},
		{
			name: "client error",
			args: map[string]any{"issue_id": "TST-1", "comment_id": "c-1"},
			setupMock: func(m *mocks.MockYouTrackClient) {
				m.EXPECT().CommentDelete(gomock.Any(), "TST-1", "c-1").
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

			handler := commentDeleteHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.True(t, result.IsError)
			text := result.Content[0].(mcplib.TextContent).Text
			assert.Contains(t, text, tt.errContains)
		})
	}
}
