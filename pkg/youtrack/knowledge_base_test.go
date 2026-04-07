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

func TestArticleList_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/articles", r.URL.Path)
		assert.Equal(t, ArticleListFields, r.URL.Query().Get("fields"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"$type":"Article","id":"150-1","idReadable":"DEMO-A-1","summary":"Getting Started","project":{"$type":"Project","id":"0-1","name":"Demo","shortName":"DEMO"},"reporter":{"$type":"User","id":"1-1","login":"admin"}}]`))
	})

	result, err := client.ArticleList(context.Background(), 0, 0)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "150-1", result[0].ID)
	assert.Equal(t, "DEMO-A-1", result[0].IDReadable)
	assert.Equal(t, "Getting Started", result[0].Summary)
	assert.Equal(t, "DEMO", result[0].Project.ShortName)
}

func TestArticleList_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	result, err := client.ArticleList(context.Background(), 0, 0)
	require.ErrorIs(t, err, ErrUnauthorized)
	require.Nil(t, result)
}

func TestArticleGet_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/articles/DEMO-A-1", r.URL.Path)
		assert.Equal(t, ArticleDetailFields, r.URL.Query().Get("fields"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"$type":"Article","id":"150-1","idReadable":"DEMO-A-1","summary":"Getting Started","content":"Welcome to the project.","reporter":{"$type":"User","id":"1-1","login":"admin"},"childArticles":[{"$type":"Article","id":"150-2","idReadable":"DEMO-A-2","summary":"Installation"}]}`))
	})

	result, err := client.ArticleGet(context.Background(), "DEMO-A-1")
	require.NoError(t, err)
	assert.Equal(t, "150-1", result.ID)
	assert.Equal(t, "Getting Started", result.Summary)
	assert.Equal(t, "Welcome to the project.", result.Content)
	require.Len(t, result.ChildArticles, 1)
	assert.Equal(t, "Installation", result.ChildArticles[0].Summary)
}

func TestArticleGet_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := client.ArticleGet(context.Background(), "nonexistent")
	require.ErrorIs(t, err, ErrNotFound)
}

func TestArticleCreate_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/articles", r.URL.Path)
		assert.Equal(t, ArticleDetailFields, r.URL.Query().Get("fields"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var reqBody Article
		require.NoError(t, json.Unmarshal(body, &reqBody))
		assert.Equal(t, "New Article", reqBody.Summary)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"$type":"Article","id":"150-3","idReadable":"DEMO-A-3","summary":"New Article","content":"Article content"}`))
	})

	result, err := client.ArticleCreate(context.Background(), Article{
		Summary: "New Article",
		Content: "Article content",
		Project: &Project{ID: "0-1"},
	})
	require.NoError(t, err)
	assert.Equal(t, "150-3", result.ID)
	assert.Equal(t, "New Article", result.Summary)
}

func TestArticleCreate_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})

	_, err := client.ArticleCreate(context.Background(), Article{Summary: "test"})
	require.ErrorIs(t, err, ErrForbidden)
}

func TestArticleUpdate_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/articles/DEMO-A-1", r.URL.Path)
		assert.Equal(t, ArticleDetailFields, r.URL.Query().Get("fields"))

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var reqBody Article
		require.NoError(t, json.Unmarshal(body, &reqBody))
		assert.Equal(t, "Updated Article", reqBody.Summary)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"$type":"Article","id":"150-1","idReadable":"DEMO-A-1","summary":"Updated Article"}`))
	})

	result, err := client.ArticleUpdate(context.Background(), "DEMO-A-1", Article{Summary: "Updated Article"})
	require.NoError(t, err)
	assert.Equal(t, "Updated Article", result.Summary)
}

func TestArticleUpdate_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := client.ArticleUpdate(context.Background(), "nonexistent", Article{Summary: "x"})
	require.ErrorIs(t, err, ErrNotFound)
}

func TestArticleDelete_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/articles/DEMO-A-1", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	})

	err := client.ArticleDelete(context.Background(), "DEMO-A-1")
	require.NoError(t, err)
}

func TestArticleDelete_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	err := client.ArticleDelete(context.Background(), "nonexistent")
	require.ErrorIs(t, err, ErrNotFound)
}
