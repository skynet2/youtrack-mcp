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

func TestUserList_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().UserList(gomock.Any(), "admin", 0, 10).
		Return([]youtrack.User{
			{ID: "u-1", Login: "admin", FullName: "Admin User"},
		}, nil)

	handler := userListHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{
		"query": "admin",
		"top":   float64(10),
	}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "admin")
	assert.Contains(t, text, "Admin User")
}

func TestUserList_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().UserList(gomock.Any(), "", 0, 0).
		Return(nil, errors.New("timeout"))

	handler := userListHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{}))

	require.NoError(t, err)
	require.True(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "timeout")
}

func TestUserCurrent_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().UserCurrent(gomock.Any()).
		Return(youtrack.User{ID: "u-me", Login: "testuser", FullName: "Test User"}, nil)

	handler := userCurrentHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{}))

	require.NoError(t, err)
	require.False(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "testuser")
	assert.Contains(t, text, "Test User")
}

func TestUserCurrent_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	mock := mocks.NewMockYouTrackClient(ctrl)
	log := zerolog.Nop()

	mock.EXPECT().UserCurrent(gomock.Any()).
		Return(youtrack.User{}, errors.New("unauthorized"))

	handler := userCurrentHandler(mock, log)
	result, err := handler(context.Background(), newReq(map[string]any{}))

	require.NoError(t, err)
	require.True(t, result.IsError)
	text := result.Content[0].(mcplib.TextContent).Text
	assert.Contains(t, text, "unauthorized")
}
