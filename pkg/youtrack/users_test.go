package youtrack

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserList_Success(t *testing.T) {
	tests := []struct {
		name  string
		query string
		skip  int
		top   int
	}{
		{
			name:  "with query",
			query: "admin",
			skip:  0,
			top:   10,
		},
		{
			name:  "without query",
			query: "",
			skip:  0,
			top:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "/api/users", r.URL.Path)
				assert.Equal(t, UserFields, r.URL.Query().Get("fields"))

				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`[{"$type":"User","id":"1-1","login":"admin","fullName":"Admin User","email":"admin@example.com","ringId":"abc-123"}]`))
			})

			result, err := client.UserList(context.Background(), tt.query, tt.skip, tt.top)
			require.NoError(t, err)
			require.Len(t, result, 1)
			assert.Equal(t, "1-1", result[0].ID)
			assert.Equal(t, "admin", result[0].Login)
			assert.Equal(t, "Admin User", result[0].FullName)
			assert.Equal(t, "admin@example.com", result[0].Email)
		})
	}
}

func TestUserList_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	result, err := client.UserList(context.Background(), "", 0, 0)
	require.ErrorIs(t, err, ErrUnauthorized)
	require.Nil(t, result)
}

func TestUserCurrent_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/users/me", r.URL.Path)
		assert.Equal(t, UserFields, r.URL.Query().Get("fields"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"$type":"User","id":"1-1","login":"currentuser","fullName":"Current User","email":"user@example.com","guest":false,"online":true}`))
	})

	result, err := client.UserCurrent(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "1-1", result.ID)
	assert.Equal(t, "currentuser", result.Login)
	assert.Equal(t, "Current User", result.FullName)
}

func TestUserCurrent_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	_, err := client.UserCurrent(context.Background())
	require.ErrorIs(t, err, ErrUnauthorized)
}
