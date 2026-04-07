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

func TestWorkItemList_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/issues/DEMO-1/timeTracking/workItems", r.URL.Path)
		assert.Equal(t, WorkItemFields, r.URL.Query().Get("fields"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"$type":"IssueWorkItem","id":"7-1","text":"Coding","author":{"$type":"User","id":"1-1","login":"dev"},"duration":{"$type":"DurationValue","minutes":60,"presentation":"1h"}}]`))
	})

	result, err := client.WorkItemList(context.Background(), "DEMO-1", 0, 0)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "7-1", result[0].ID)
	assert.Equal(t, "Coding", result[0].Text)
	assert.Equal(t, "dev", result[0].Author.Login)
	require.NotNil(t, result[0].Duration)
	assert.Equal(t, int32(60), *result[0].Duration.Minutes)
	assert.Equal(t, "1h", result[0].Duration.Presentation)
}

func TestWorkItemList_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	result, err := client.WorkItemList(context.Background(), "NONEXIST-1", 0, 0)
	require.ErrorIs(t, err, ErrNotFound)
	require.Nil(t, result)
}

func TestWorkItemCreate_Success(t *testing.T) {
	minutes := int32(120)

	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/issues/DEMO-1/timeTracking/workItems", r.URL.Path)
		assert.Equal(t, WorkItemFields, r.URL.Query().Get("fields"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var reqBody IssueWorkItem
		require.NoError(t, json.Unmarshal(body, &reqBody))
		assert.Equal(t, "Review", reqBody.Text)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"$type":"IssueWorkItem","id":"7-2","text":"Review","duration":{"$type":"DurationValue","minutes":120,"presentation":"2h"}}`))
	})

	result, err := client.WorkItemCreate(context.Background(), "DEMO-1", IssueWorkItem{
		Text:     "Review",
		Duration: &DurationValue{Minutes: &minutes},
	})
	require.NoError(t, err)
	assert.Equal(t, "7-2", result.ID)
	assert.Equal(t, "Review", result.Text)
	require.NotNil(t, result.Duration)
	assert.Equal(t, int32(120), *result.Duration.Minutes)
}

func TestWorkItemCreate_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	})

	_, err := client.WorkItemCreate(context.Background(), "DEMO-1", IssueWorkItem{Text: "test"})
	require.ErrorIs(t, err, ErrForbidden)
}
