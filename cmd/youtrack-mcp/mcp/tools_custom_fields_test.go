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

func TestIssueCustomFields_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().IssueCustomFields(gomock.Any(), "TST-1").
		Return([]youtrack.IssueCustomField{
			{ID: "cf-1", Name: "Priority"},
			{ID: "cf-2", Name: "State"},
		}, nil)

	handler := issueCustomFieldsHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{"issue_id": "TST-1"}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "Priority")
	assert.Contains(t, text, "State")
}

func TestIssueCustomFields_Failure(t *testing.T) {
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
				m.EXPECT().IssueCustomFields(gomock.Any(), "TST-1").
					Return(nil, errors.New("internal error"))
			},
			errContains: "internal error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockYouTrackClient(ctrl)
			log := zerolog.Nop()

			tt.setupMock(mock)

			handler := issueCustomFieldsHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.True(t, result.IsError)
			text := result.Content[0].(mcplib.TextContent).Text
			assert.Contains(t, text, tt.errContains)
		})
	}
}

func TestIssueUpdateField_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	value := map[string]any{"name": "Critical"}
	mock.EXPECT().IssueUpdateField(gomock.Any(), "TST-1", "cf-1", youtrack.IssueCustomField{Value: value}).
		Return(youtrack.IssueCustomField{ID: "cf-1", Name: "Priority", Value: value}, nil)

	handler := issueUpdateFieldHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{
		"issue_id": "TST-1",
		"field_id": "cf-1",
		"value":    value,
	}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "Priority")
}

func TestIssueUpdateField_Failure(t *testing.T) {
	tests := []struct {
		name        string
		args        map[string]any
		setupMock   func(m *mocks.MockYouTrackClient)
		errContains string
	}{
		{
			name:        "missing issue_id",
			args:        map[string]any{"field_id": "cf-1", "value": map[string]any{"name": "x"}},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "issue_id",
		},
		{
			name:        "missing field_id",
			args:        map[string]any{"issue_id": "TST-1", "value": map[string]any{"name": "x"}},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "field_id",
		},
		{
			name:        "missing value",
			args:        map[string]any{"issue_id": "TST-1", "field_id": "cf-1"},
			setupMock:   func(m *mocks.MockYouTrackClient) {},
			errContains: "value",
		},
		{
			name: "client error",
			args: map[string]any{"issue_id": "TST-1", "field_id": "cf-1", "value": map[string]any{"name": "x"}},
			setupMock: func(m *mocks.MockYouTrackClient) {
				m.EXPECT().IssueUpdateField(gomock.Any(), "TST-1", "cf-1", gomock.Any()).
					Return(youtrack.IssueCustomField{}, errors.New("invalid value"))
			},
			errContains: "invalid value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockYouTrackClient(ctrl)
			log := zerolog.Nop()

			tt.setupMock(mock)

			handler := issueUpdateFieldHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.True(t, result.IsError)
			text := result.Content[0].(mcplib.TextContent).Text
			assert.Contains(t, text, tt.errContains)
		})
	}
}
