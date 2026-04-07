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

func TestActivityList_Success(t *testing.T) {
	tests := []struct {
		name       string
		args       map[string]any
		categories string
		issueQuery string
		start      string
		end        string
		skip       int
		top        int
		response   []youtrack.ActivityItem
	}{
		{
			name: "with all parameters",
			args: map[string]any{
				"categories":  "IssueCreatedCategory",
				"issue_query": "project: TST",
				"start":       "1700000000000",
				"end":         "1700100000000",
				"skip":        float64(0),
				"top":         float64(10),
			},
			categories: "IssueCreatedCategory",
			issueQuery: "project: TST",
			start:      "1700000000000",
			end:        "1700100000000",
			skip:       0,
			top:        10,
			response: []youtrack.ActivityItem{
				{ID: "act-1", TargetMember: "summary"},
			},
		},
		{
			name:     "no parameters",
			args:     map[string]any{},
			response: []youtrack.ActivityItem{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mocks.NewMockYouTrackClient(ctrl)
			log := zerolog.Nop()

			mock.EXPECT().ActivityList(gomock.Any(), tt.categories, tt.issueQuery, tt.start, tt.end, tt.skip, tt.top).
				Return(tt.response, nil)

			handler := activityListHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.False(t, result.IsError)
			require.Len(t, result.Content, 1)

			text := result.Content[0].(mcplib.TextContent).Text
			for _, item := range tt.response {
				assert.Contains(t, text, item.ID)
			}
		})
	}
}

func TestActivityList_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().ActivityList(gomock.Any(), "", "", "", "", 0, 0).
		Return(nil, errors.New("gateway timeout"))

	handler := activityListHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{}))

	require.NoError(t, err)
	require.True(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "gateway timeout")
}

func TestIssueActivityList_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().IssueActivityList(gomock.Any(), "TST-1", "CommentsCategory", 0, 5).
		Return([]youtrack.ActivityItem{
			{ID: "act-2", TargetMember: "comments"},
		}, nil)

	handler := issueActivityListHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{
		"issue_id":   "TST-1",
		"categories": "CommentsCategory",
		"top":        float64(5),
	}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "act-2")
}

func TestIssueActivityList_Failure(t *testing.T) {
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
				m.EXPECT().IssueActivityList(gomock.Any(), "TST-1", "", 0, 0).
					Return(nil, errors.New("server error"))
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

			handler := issueActivityListHandler(mock, log)
			result, err := handler(context.Background(), newReq(tt.args))

			require.NoError(t, err)
			require.True(t, result.IsError)
			text := result.Content[0].(mcplib.TextContent).Text
			assert.Contains(t, text, tt.errContains)
		})
	}
}
