package youtrack

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgileList_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/agiles", r.URL.Path)
		assert.Equal(t, AgileFields, r.URL.Query().Get("fields"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"$type":"Agile","id":"104-1","name":"Scrum Board","owner":{"$type":"User","id":"1-1","login":"admin"},"projects":[{"$type":"Project","id":"0-1","name":"Demo","shortName":"DEMO"}],"currentSprint":{"$type":"Sprint","id":"105-1","name":"Sprint 1"}}]`))
	})

	result, err := client.AgileList(context.Background(), 0, 0)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "104-1", result[0].ID)
	assert.Equal(t, "Scrum Board", result[0].Name)
	assert.Equal(t, "admin", result[0].Owner.Login)
	require.Len(t, result[0].Projects, 1)
	assert.Equal(t, "DEMO", result[0].Projects[0].ShortName)
	require.NotNil(t, result[0].CurrentSprint)
	assert.Equal(t, "Sprint 1", result[0].CurrentSprint.Name)
}

func TestAgileList_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	result, err := client.AgileList(context.Background(), 0, 0)
	require.ErrorIs(t, err, ErrUnauthorized)
	require.Nil(t, result)
}

func TestSprintList_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/agiles/104-1/sprints", r.URL.Path)
		assert.Equal(t, SprintFields, r.URL.Query().Get("fields"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"$type":"Sprint","id":"105-1","name":"Sprint 1","goal":"Deliver MVP"},{"$type":"Sprint","id":"105-2","name":"Sprint 2"}]`))
	})

	result, err := client.SprintList(context.Background(), "104-1", 0, 0)
	require.NoError(t, err)
	require.Len(t, result, 2)
	assert.Equal(t, "105-1", result[0].ID)
	assert.Equal(t, "Sprint 1", result[0].Name)
	assert.Equal(t, "Deliver MVP", result[0].Goal)
	assert.Equal(t, "105-2", result[1].ID)
}

func TestSprintList_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	result, err := client.SprintList(context.Background(), "nonexistent", 0, 0)
	require.ErrorIs(t, err, ErrNotFound)
	require.Nil(t, result)
}

func TestSprintGet_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/agiles/104-1/sprints/105-1", r.URL.Path)
		assert.Equal(t, SprintDetailFields, r.URL.Query().Get("fields"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"$type":"Sprint","id":"105-1","name":"Sprint 1","goal":"Deliver MVP","issues":[{"$type":"Issue","id":"2-123","idReadable":"DEMO-1","summary":"Test issue"}],"unresolvedIssuesCount":3}`))
	})

	result, err := client.SprintGet(context.Background(), "104-1", "105-1")
	require.NoError(t, err)
	assert.Equal(t, "105-1", result.ID)
	assert.Equal(t, "Sprint 1", result.Name)
	assert.Equal(t, "Deliver MVP", result.Goal)
	require.Len(t, result.Issues, 1)
	assert.Equal(t, "DEMO-1", result.Issues[0].IDReadable)
	require.NotNil(t, result.UnresolvedIssuesCount)
	assert.Equal(t, int32(3), *result.UnresolvedIssuesCount)
}

func TestSprintGet_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := client.SprintGet(context.Background(), "104-1", "nonexistent")
	require.ErrorIs(t, err, ErrNotFound)
}
