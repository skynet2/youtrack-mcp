package youtrack

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommentList_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/issues/DEMO-1/comments", r.URL.Path)
		assert.Equal(t, CommentFields, r.URL.Query().Get("fields"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"$type":"IssueComment","id":"4-1","text":"First comment","author":{"$type":"User","id":"1-1","login":"admin","fullName":"Admin User"}}]`))
	})

	result, err := client.CommentList(context.Background(), "DEMO-1", 0, 0)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "4-1", result[0].ID)
	assert.Equal(t, "First comment", result[0].Text)
	assert.Equal(t, "admin", result[0].Author.Login)
}

func TestCommentList_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	result, err := client.CommentList(context.Background(), "NONEXIST-1", 0, 0)
	require.ErrorIs(t, err, ErrNotFound)
	require.Nil(t, result)
}

func TestCommentCreate_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/issues/DEMO-1/comments", r.URL.Path)
		assert.Equal(t, CommentFields, r.URL.Query().Get("fields"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var reqBody IssueComment
		require.NoError(t, json.Unmarshal(body, &reqBody))
		assert.Equal(t, "New comment text", reqBody.Text)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"$type":"IssueComment","id":"4-2","text":"New comment text","author":{"$type":"User","id":"1-1","login":"admin"}}`))
	})

	result, err := client.CommentCreate(context.Background(), "DEMO-1", IssueComment{Text: "New comment text"})
	require.NoError(t, err)
	assert.Equal(t, "4-2", result.ID)
	assert.Equal(t, "New comment text", result.Text)
}

func TestCommentCreate_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})

	_, err := client.CommentCreate(context.Background(), "DEMO-1", IssueComment{Text: "test"})
	require.ErrorIs(t, err, ErrForbidden)
}

func TestCommentUpdate_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/issues/DEMO-1/comments/4-1", r.URL.Path)
		assert.Equal(t, CommentFields, r.URL.Query().Get("fields"))

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var reqBody IssueComment
		require.NoError(t, json.Unmarshal(body, &reqBody))
		assert.Equal(t, "Updated text", reqBody.Text)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"$type":"IssueComment","id":"4-1","text":"Updated text"}`))
	})

	result, err := client.CommentUpdate(context.Background(), "DEMO-1", "4-1", IssueComment{Text: "Updated text"})
	require.NoError(t, err)
	assert.Equal(t, "4-1", result.ID)
	assert.Equal(t, "Updated text", result.Text)
}

func TestCommentUpdate_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := client.CommentUpdate(context.Background(), "DEMO-1", "nonexistent", IssueComment{Text: "x"})
	require.ErrorIs(t, err, ErrNotFound)
}

func TestCommentDelete_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/issues/DEMO-1/comments/4-1", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	})

	err := client.CommentDelete(context.Background(), "DEMO-1", "4-1")
	require.NoError(t, err)
}

func TestCommentDelete_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	err := client.CommentDelete(context.Background(), "DEMO-1", "nonexistent")
	require.ErrorIs(t, err, ErrNotFound)
}
