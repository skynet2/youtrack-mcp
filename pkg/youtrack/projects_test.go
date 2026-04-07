package youtrack

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectList_Success(t *testing.T) {
	tests := []struct {
		name string
		skip int
		top  int
	}{
		{
			name: "with pagination",
			skip: 5,
			top:  10,
		},
		{
			name: "without pagination",
			skip: 0,
			top:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "/api/admin/projects", r.URL.Path)
				assert.Equal(t, ProjectFields, r.URL.Query().Get("fields"))
				assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`[{"$type":"Project","id":"0-1","name":"Demo","shortName":"DEMO","description":"Demo project"}]`))
			})

			result, err := client.ProjectList(context.Background(), tt.skip, tt.top)
			require.NoError(t, err)
			require.Len(t, result, 1)
			assert.Equal(t, "0-1", result[0].ID)
			assert.Equal(t, "Demo", result[0].Name)
			assert.Equal(t, "DEMO", result[0].ShortName)
			assert.Equal(t, "Demo project", result[0].Description)
		})
	}
}

func TestProjectList_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	result, err := client.ProjectList(context.Background(), 0, 0)
	require.ErrorIs(t, err, ErrUnauthorized)
	require.Nil(t, result)
}

func TestProjectGet_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/admin/projects/0-1", r.URL.Path)
		assert.Equal(t, ProjectFields, r.URL.Query().Get("fields"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"$type":"Project","id":"0-1","name":"Demo","shortName":"DEMO","description":"A demo project","leader":{"$type":"User","id":"1-1","login":"admin","fullName":"Admin User"}}`))
	})

	result, err := client.ProjectGet(context.Background(), "0-1")
	require.NoError(t, err)
	assert.Equal(t, "0-1", result.ID)
	assert.Equal(t, "Demo", result.Name)
	assert.Equal(t, "DEMO", result.ShortName)
	require.NotNil(t, result.Leader)
	assert.Equal(t, "admin", result.Leader.Login)
}

func TestProjectGet_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := client.ProjectGet(context.Background(), "nonexistent")
	require.ErrorIs(t, err, ErrNotFound)
}
