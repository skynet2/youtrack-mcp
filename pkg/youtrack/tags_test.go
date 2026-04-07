package youtrack

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTagList_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/issueTags", r.URL.Path)
		assert.Equal(t, TagFields, r.URL.Query().Get("fields"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"$type":"IssueTag","id":"6-1","name":"Bug","owner":{"$type":"User","id":"1-1","login":"admin"},"color":{"$type":"FieldStyle","id":"0","background":"#ff0000","foreground":"#ffffff"}}]`))
	})

	result, err := client.TagList(context.Background(), 0, 0)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "6-1", result[0].ID)
	assert.Equal(t, "Bug", result[0].Name)
	assert.Equal(t, "admin", result[0].Owner.Login)
	require.NotNil(t, result[0].Color)
	assert.Equal(t, "#ff0000", result[0].Color.Background)
}

func TestTagList_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	result, err := client.TagList(context.Background(), 0, 0)
	require.ErrorIs(t, err, ErrUnauthorized)
	require.Nil(t, result)
}
