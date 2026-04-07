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

func TestIssueSearch_Success(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		skip     int
		top      int
		wantPath string
	}{
		{
			name:     "with all parameters",
			query:    "project: DEMO",
			skip:     10,
			top:      25,
			wantPath: "/api/issues",
		},
		{
			name:     "with no optional parameters",
			query:    "",
			skip:     0,
			top:      0,
			wantPath: "/api/issues",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, tt.wantPath, r.URL.Path)
				assert.Equal(t, IssueListFields, r.URL.Query().Get("fields"))
				assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`[{"$type":"Issue","id":"2-123","idReadable":"DEMO-1","summary":"Test issue","project":{"$type":"Project","id":"0-1","name":"Demo","shortName":"DEMO"}}]`))
			})

			result, err := client.IssueSearch(context.Background(), tt.query, tt.skip, tt.top)
			require.NoError(t, err)
			require.Len(t, result, 1)
			assert.Equal(t, "2-123", result[0].ID)
			assert.Equal(t, "DEMO-1", result[0].IDReadable)
			assert.Equal(t, "Test issue", result[0].Summary)
			assert.Equal(t, "DEMO", result[0].Project.ShortName)
		})
	}
}

func TestIssueSearch_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	result, err := client.IssueSearch(context.Background(), "", 0, 0)
	require.ErrorIs(t, err, ErrUnauthorized)
	require.Nil(t, result)
}

func TestIssueGet_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/issues/DEMO-1", r.URL.Path)
		assert.Equal(t, IssueDetailFields, r.URL.Query().Get("fields"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"$type":"Issue","id":"2-123","idReadable":"DEMO-1","summary":"Test issue","description":"Some description"}`))
	})

	result, err := client.IssueGet(context.Background(), "DEMO-1")
	require.NoError(t, err)
	assert.Equal(t, "2-123", result.ID)
	assert.Equal(t, "DEMO-1", result.IDReadable)
	assert.Equal(t, "Some description", result.Description)
}

func TestIssueGet_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := client.IssueGet(context.Background(), "NONEXIST-1")
	require.ErrorIs(t, err, ErrNotFound)
}

func TestIssueCreate_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/issues", r.URL.Path)
		assert.Equal(t, IssueDetailFields, r.URL.Query().Get("fields"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var reqBody Issue
		require.NoError(t, json.Unmarshal(body, &reqBody))
		assert.Equal(t, "New issue", reqBody.Summary)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"$type":"Issue","id":"2-124","idReadable":"DEMO-2","summary":"New issue"}`))
	})

	result, err := client.IssueCreate(context.Background(), Issue{
		Summary: "New issue",
		Project: &Project{ID: "0-1"},
	})
	require.NoError(t, err)
	assert.Equal(t, "2-124", result.ID)
	assert.Equal(t, "DEMO-2", result.IDReadable)
}

func TestIssueCreate_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})

	_, err := client.IssueCreate(context.Background(), Issue{Summary: "New issue"})
	require.ErrorIs(t, err, ErrForbidden)
}

func TestIssueUpdate_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/issues/DEMO-1", r.URL.Path)
		assert.Equal(t, IssueDetailFields, r.URL.Query().Get("fields"))

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var reqBody Issue
		require.NoError(t, json.Unmarshal(body, &reqBody))
		assert.Equal(t, "Updated summary", reqBody.Summary)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"$type":"Issue","id":"2-123","idReadable":"DEMO-1","summary":"Updated summary"}`))
	})

	result, err := client.IssueUpdate(context.Background(), "DEMO-1", Issue{Summary: "Updated summary"})
	require.NoError(t, err)
	assert.Equal(t, "Updated summary", result.Summary)
}

func TestIssueUpdate_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := client.IssueUpdate(context.Background(), "NONEXIST-1", Issue{Summary: "Updated"})
	require.ErrorIs(t, err, ErrNotFound)
}

func TestIssueDelete_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/issues/DEMO-1", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	})

	err := client.IssueDelete(context.Background(), "DEMO-1")
	require.NoError(t, err)
}

func TestIssueDelete_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	err := client.IssueDelete(context.Background(), "NONEXIST-1")
	require.ErrorIs(t, err, ErrNotFound)
}

func TestIssueAddTag_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/issues/DEMO-1/tags", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var reqBody IssueTag
		require.NoError(t, json.Unmarshal(body, &reqBody))
		assert.Equal(t, "6-1", reqBody.ID)

		w.WriteHeader(http.StatusOK)
	})

	err := client.IssueAddTag(context.Background(), "DEMO-1", "6-1")
	require.NoError(t, err)
}

func TestIssueAddTag_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	err := client.IssueAddTag(context.Background(), "NONEXIST-1", "6-1")
	require.ErrorIs(t, err, ErrNotFound)
}

func TestIssueRemoveTag_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/api/issues/DEMO-1/tags/6-1", r.URL.Path)

		w.WriteHeader(http.StatusOK)
	})

	err := client.IssueRemoveTag(context.Background(), "DEMO-1", "6-1")
	require.NoError(t, err)
}

func TestIssueRemoveTag_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})

	err := client.IssueRemoveTag(context.Background(), "DEMO-1", "6-1")
	require.ErrorIs(t, err, ErrForbidden)
}

func TestIssueGetLinks_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/issues/DEMO-1/links", r.URL.Path)
		assert.Equal(t, IssueLinkFields, r.URL.Query().Get("fields"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"$type":"IssueLink","id":"3-1","direction":"OUTWARD","linkType":{"$type":"IssueLinkType","id":"4-1","name":"Depends on"}}]`))
	})

	result, err := client.IssueGetLinks(context.Background(), "DEMO-1")
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "3-1", result[0].ID)
	assert.Equal(t, "OUTWARD", result[0].Direction)
	assert.Equal(t, "Depends on", result[0].LinkType.Name)
}

func TestIssueGetLinks_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := client.IssueGetLinks(context.Background(), "NONEXIST-1")
	require.ErrorIs(t, err, ErrNotFound)
}
