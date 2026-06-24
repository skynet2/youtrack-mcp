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

func TestLinkTypesHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().ListIssueLinkTypes(gomock.Any()).
		Return([]youtrack.IssueLinkType{{Name: "Subtask", TargetToSource: "subtask of"}}, nil)

	handler := issueLinkTypesHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "subtask of")
}

func TestIssueLinkHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().IssueLink(gomock.Any(), "IT-9", "subtask of", "IT-8").
		DoAndReturn(func(_ context.Context, source, phrase, target string) error {
			assert.Equal(t, "IT-9", source)
			assert.Equal(t, "subtask of", phrase)
			assert.Equal(t, "IT-8", target)
			return nil
		})

	handler := issueLinkHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{
		"source": "IT-9", "link_type": "subtask of", "target": "IT-8",
	}))

	require.NoError(t, err)
	require.False(t, result.IsError)
}

func TestIssueUnlinkHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().IssueUnlink(gomock.Any(), "IT-9", "subtask of", "IT-8").Return(nil)

	handler := issueUnlinkHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{
		"source": "IT-9", "link_type": "subtask of", "target": "IT-8",
	}))

	require.NoError(t, err)
	require.False(t, result.IsError)
}

func TestIssueLinkHandler_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().IssueLink(gomock.Any(), "IT-9", "subtask of", "IT-8").
		Return(errors.New("boom"))

	handler := issueLinkHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{
		"source": "IT-9", "link_type": "subtask of", "target": "IT-8",
	}))

	require.NoError(t, err)
	require.True(t, result.IsError)
}
