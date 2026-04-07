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

func TestAgileList_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().AgileList(gomock.Any(), 0, 5).
		Return([]youtrack.Agile{
			{ID: "agile-1", Name: "Board Alpha"},
		}, nil)

	handler := agileListHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{"top": float64(5)}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "Board Alpha")
}

func TestAgileList_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().AgileList(gomock.Any(), 0, 0).
		Return(nil, errors.New("forbidden"))

	handler := agileListHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{}))

	require.NoError(t, err)
	require.True(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "forbidden")
}

func TestSprintList_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().SprintList(gomock.Any(), "agile-1", 0, 10).
		Return([]youtrack.Sprint{
			{ID: "sprint-1", Name: "Sprint 1"},
			{ID: "sprint-2", Name: "Sprint 2"},
		}, nil)

	handler := sprintListHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{
		"agile_id": "agile-1",
		"top":      float64(10),
	}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "Sprint 1")
	assert.Contains(t, text, "Sprint 2")
}

func TestSprintList_Failure(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]any
		setupMock   func(m *mocks.MockYouTrackClient)
		errContains string
	}{
		{
			name:        "missing agile_id",
			args:        map[string]any{},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "agile_id",
		},
		{
			name: "client error",
			args: map[string]any{"agile_id": "agile-1"},
			setupMock: func(m *mocks.MockYouTrackClient) {
				m.EXPECT().SprintList(gomock.Any(), "agile-1", 0, 0).
					Return(nil, errors.New("not found"))
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

			handler := sprintListHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.True(t, result.IsError)
			text := result.Content[0].(mcplib.TextContent).Text
			assert.Contains(t, text, tt.errContains)
		})
	}
}

func TestSprintGet_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().SprintGet(gomock.Any(), "agile-1", "sprint-1").
		Return(youtrack.Sprint{ID: "sprint-1", Name: "Sprint 1", Goal: "Ship feature"}, nil)

	handler := sprintGetHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{
		"agile_id":  "agile-1",
		"sprint_id": "sprint-1",
	}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "Sprint 1")
	assert.Contains(t, text, "Ship feature")
}

func TestSprintGet_Failure(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]any
		setupMock   func(m *mocks.MockYouTrackClient)
		errContains string
	}{
		{
			name:        "missing agile_id",
			args:        map[string]any{"sprint_id": "s-1"},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "agile_id",
		},
		{
			name:        "missing sprint_id",
			args:        map[string]any{"agile_id": "a-1"},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "sprint_id",
		},
		{
			name: "client error",
			args: map[string]any{"agile_id": "a-1", "sprint_id": "s-1"},
			setupMock: func(m *mocks.MockYouTrackClient) {
				m.EXPECT().SprintGet(gomock.Any(), "a-1", "s-1").
					Return(youtrack.Sprint{}, errors.New("access denied"))
			},
			errContains: "access denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockYouTrackClient(ctrl)
			log := zerolog.Nop()

			tt.setupMock(mock)

			handler := sprintGetHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.True(t, result.IsError)
			text := result.Content[0].(mcplib.TextContent).Text
			assert.Contains(t, text, tt.errContains)
		})
	}
}
