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

func TestWorkItemList_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	minutes := int32(60)
	mock.EXPECT().WorkItemList(gomock.Any(), "TST-1", 0, 0).
		Return([]youtrack.IssueWorkItem{
			{ID: "wi-1", Duration: &youtrack.DurationValue{Minutes: &minutes}},
		}, nil)

	handler := workItemListHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{"issue_id": "TST-1"}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "wi-1")
}

func TestWorkItemList_Failure(t *testing.T) {
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
				m.EXPECT().WorkItemList(gomock.Any(), "TST-1", 0, 0).
					Return(nil, errors.New("network error"))
			},
			errContains: "network error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockYouTrackClient(ctrl)
			log := zerolog.Nop()

			tt.setupMock(mock)

			handler := workItemListHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.True(t, result.IsError)
			text := result.Content[0].(mcplib.TextContent).Text
			assert.Contains(t, text, tt.errContains)
		})
	}
}

func TestWorkItemCreate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	m := int32(30)
	expectedItem := youtrack.IssueWorkItem{
		Duration: &youtrack.DurationValue{Minutes: &m},
	}
	mock.EXPECT().WorkItemCreate(gomock.Any(), "TST-1", expectedItem).
		Return(youtrack.IssueWorkItem{ID: "wi-new"}, nil)

	handler := workItemCreateHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{
		"issue_id": "TST-1",
		"minutes":  float64(30),
	}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "wi-new")
}

func TestWorkItemCreate_SuccessWithOptionalFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	m := int32(45)
	d := int64(1700000000000)
	expectedItem := youtrack.IssueWorkItem{
		Duration: &youtrack.DurationValue{Minutes: &m},
		Text:     "Working on feature",
		WorkType: &youtrack.WorkItemType{ID: "type-1"},
		Date:     &d,
	}
	mock.EXPECT().WorkItemCreate(gomock.Any(), "TST-2", expectedItem).
		Return(youtrack.IssueWorkItem{ID: "wi-full"}, nil)

	handler := workItemCreateHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{
		"issue_id": "TST-2",
		"minutes":  float64(45),
		"text":     "Working on feature",
		"type_id":  "type-1",
		"date":     float64(1700000000000),
	}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "wi-full")
}

func TestWorkItemCreate_Failure(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]any
		setupMock   func(m *mocks.MockYouTrackClient)
		errContains string
	}{
		{
			name:        "missing issue_id",
			args:        map[string]any{"minutes": float64(30)},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "issue_id",
		},
		{
			name:        "missing minutes",
			args:        map[string]any{"issue_id": "TST-1"},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "minutes",
		},
		{
			name: "client error",
			args: map[string]any{"issue_id": "TST-1", "minutes": float64(30)},
			setupMock: func(m *mocks.MockYouTrackClient) {
				m.EXPECT().WorkItemCreate(gomock.Any(), "TST-1", gomock.Any()).
					Return(youtrack.IssueWorkItem{}, errors.New("validation error"))
			},
			errContains: "validation error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockYouTrackClient(ctrl)
			log := zerolog.Nop()

			tt.setupMock(mock)

			handler := workItemCreateHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.True(t, result.IsError)
			text := result.Content[0].(mcplib.TextContent).Text
			assert.Contains(t, text, tt.errContains)
		})
	}
}
