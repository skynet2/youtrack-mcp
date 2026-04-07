package youtrack

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	client := New(Config{
		BaseURL:    srv.URL,
		Token:      "test-token",
		HTTPClient: srv.Client(),
	})
	return client, srv
}

func TestClient_AuthHeader(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	})

	var out map[string]any
	err := client.get(context.Background(), "/api/test", nil, &out)
	require.NoError(t, err)
}

func TestClient_ErrorHandling_Success(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantErr    error
	}{
		{
			name:       "401 returns ErrUnauthorized",
			statusCode: http.StatusUnauthorized,
			wantErr:    ErrUnauthorized,
		},
		{
			name:       "404 returns ErrNotFound",
			statusCode: http.StatusNotFound,
			wantErr:    ErrNotFound,
		},
		{
			name:       "403 returns ErrForbidden",
			statusCode: http.StatusForbidden,
			wantErr:    ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			})

			var out map[string]any
			err := client.get(context.Background(), "/api/test", nil, &out)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestClient_ErrorHandling_APIError(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`internal server error`))
	})

	var out map[string]any
	err := client.get(context.Background(), "/api/test", nil, &out)
	require.Error(t, err)

	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))
	assert.Equal(t, http.StatusInternalServerError, apiErr.StatusCode)
	assert.Equal(t, "internal server error", apiErr.Message)
}

func TestClient_NilHTTPClient(t *testing.T) {
	c := New(Config{
		BaseURL: "http://localhost",
		Token:   "tok",
	})
	assert.Equal(t, 30*time.Second, c.httpClient.Timeout)
}

func TestClient_BaseURLTrailingSlash(t *testing.T) {
	c := New(Config{
		BaseURL: "http://localhost/",
		Token:   "tok",
	})
	assert.Equal(t, "http://localhost", c.baseURL)
}
