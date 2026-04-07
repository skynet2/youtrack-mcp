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

func TestArticleList_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().ArticleList(gomock.Any(), 0, 10).
		Return([]youtrack.Article{
			{ID: "art-1", Summary: "Getting Started"},
			{ID: "art-2", Summary: "Advanced Guide"},
		}, nil)

	handler := articleListHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{"top": float64(10)}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "Getting Started")
	assert.Contains(t, text, "Advanced Guide")
}

func TestArticleList_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().ArticleList(gomock.Any(), 0, 0).
		Return(nil, errors.New("unavailable"))

	handler := articleListHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{}))

	require.NoError(t, err)
	require.True(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "unavailable")
}

func TestArticleGet_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().ArticleGet(gomock.Any(), "art-1").
		Return(youtrack.Article{ID: "art-1", Summary: "Getting Started", Content: "Welcome"}, nil)

	handler := articleGetHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{"id": "art-1"}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "Getting Started")
	assert.Contains(t, text, "Welcome")
}

func TestArticleGet_Failure(t *testing.T) {
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
			args: map[string]any{"id": "art-1"},
			setupMock: func(m *mocks.MockYouTrackClient) {
				m.EXPECT().ArticleGet(gomock.Any(), "art-1").
					Return(youtrack.Article{}, errors.New("not found"))
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

			handler := articleGetHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.True(t, result.IsError)
			text := result.Content[0].(mcplib.TextContent).Text
			assert.Contains(t, text, tt.errContains)
		})
	}
}

func TestArticleCreate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	expectedArticle := youtrack.Article{
		Project: &youtrack.Project{ID: "proj-1"},
		Summary: "New Article",
		Content: "Article body",
	}
	mock.EXPECT().ArticleCreate(gomock.Any(), expectedArticle).
		Return(youtrack.Article{ID: "art-new", Summary: "New Article"}, nil)

	handler := articleCreateHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{
		"project_id": "proj-1",
		"summary":    "New Article",
		"content":    "Article body",
	}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "art-new")
}

func TestArticleCreate_SuccessWithParent(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	expectedArticle := youtrack.Article{
		Project:       &youtrack.Project{ID: "proj-1"},
		Summary:       "Child Article",
		ParentArticle: &youtrack.Article{ID: "art-parent"},
	}
	mock.EXPECT().ArticleCreate(gomock.Any(), expectedArticle).
		Return(youtrack.Article{ID: "art-child", Summary: "Child Article"}, nil)

	handler := articleCreateHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{
		"project_id": "proj-1",
		"summary":    "Child Article",
		"parent_id":  "art-parent",
	}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "art-child")
}

func TestArticleCreate_Failure(t *testing.T) {
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
				m.EXPECT().ArticleCreate(gomock.Any(), gomock.Any()).
					Return(youtrack.Article{}, errors.New("quota exceeded"))
			},
			errContains: "quota exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockYouTrackClient(ctrl)
			log := zerolog.Nop()

			tt.setupMock(mock)

			handler := articleCreateHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.True(t, result.IsError)
			text := result.Content[0].(mcplib.TextContent).Text
			assert.Contains(t, text, tt.errContains)
		})
	}
}

func TestArticleUpdate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().ArticleUpdate(gomock.Any(), "art-1", youtrack.Article{Summary: "Updated Title"}).
		Return(youtrack.Article{ID: "art-1", Summary: "Updated Title"}, nil)

	handler := articleUpdateHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{
		"id":      "art-1",
		"summary": "Updated Title",
	}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "Updated Title")
}

func TestArticleUpdate_Failure(t *testing.T) {
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
			args: map[string]any{"id": "art-1", "summary": "test"},
			setupMock: func(m *mocks.MockYouTrackClient) {
				m.EXPECT().ArticleUpdate(gomock.Any(), "art-1", gomock.Any()).
					Return(youtrack.Article{}, errors.New("conflict"))
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

			handler := articleUpdateHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.True(t, result.IsError)
			text := result.Content[0].(mcplib.TextContent).Text
			assert.Contains(t, text, tt.errContains)
		})
	}
}

func TestArticleDelete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().ArticleDelete(gomock.Any(), "art-1").Return(nil)

	handler := articleDeleteHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{"id": "art-1"}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Equal(t, "ok", text)
}

func TestArticleDelete_Failure(t *testing.T) {
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
			args: map[string]any{"id": "art-1"},
			setupMock: func(m *mocks.MockYouTrackClient) {
				m.EXPECT().ArticleDelete(gomock.Any(), "art-1").
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

			handler := articleDeleteHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.True(t, result.IsError)
			text := result.Content[0].(mcplib.TextContent).Text
			assert.Contains(t, text, tt.errContains)
		})
	}
}
