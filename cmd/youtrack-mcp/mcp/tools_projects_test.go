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

func TestProjectList_Success(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]any
		skip     int
		top      int
		response []youtrack.Project
	}{
		{
			name: "with pagination",
			args: map[string]any{"skip": float64(0), "top": float64(5)},
			skip: 0,
			top:  5,
			response: []youtrack.Project{
				{ID: "proj-1", Name: "Alpha", ShortName: "ALP"},
				{ID: "proj-2", Name: "Beta", ShortName: "BET"},
			},
		},
		{
			name:     "no parameters",
			args:     map[string]any{},
			skip:     0,
			top:      0,
			response: []youtrack.Project{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockYouTrackClient(ctrl)
			log := zerolog.Nop()

			mock.EXPECT().ProjectList(gomock.Any(), tt.skip, tt.top).
				Return(tt.response, nil)

			handler := projectListHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.False(t, result.IsError)
			require.Len(t, result.Content, 1)

			text := result.Content[0].(mcplib.TextContent).Text
			for _, p := range tt.response {
				assert.Contains(t, text, p.Name)
			}
		})
	}
}

func TestProjectList_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().ProjectList(gomock.Any(), 0, 0).
		Return(nil, errors.New("unauthorized"))

	handler := projectListHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{}))

	require.NoError(t, err)
	require.True(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "unauthorized")
}

func TestProjectGet_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().ProjectGet(gomock.Any(), "proj-1").
		Return(youtrack.Project{ID: "proj-1", Name: "Alpha", ShortName: "ALP"}, nil)

	handler := projectGetHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{"id": "proj-1"}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "Alpha")
	assert.Contains(t, text, "ALP")
}

func TestProjectGet_Failure(t *testing.T) {
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
			args: map[string]any{"id": "proj-1"},
			setupMock: func(m *mocks.MockYouTrackClient) {
				m.EXPECT().ProjectGet(gomock.Any(), "proj-1").
					Return(youtrack.Project{}, errors.New("not found"))
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

			handler := projectGetHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.True(t, result.IsError)
			text := result.Content[0].(mcplib.TextContent).Text
			assert.Contains(t, text, tt.errContains)
		})
	}
}
