package transport_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/skynet2/youtrack-mcp/internal/transport"
)

const testKey = "s3cr3t-key"

func newProtected(t *testing.T) http.Handler {
	t.Helper()
	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	return transport.BearerAuth(next, testKey)
}

func TestBearerAuth_Allows(t *testing.T) {
	cases := map[string]string{
		"exact bearer": "Bearer " + testKey,
	}

	for name, header := range cases {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Set("Authorization", header)
			rec := httptest.NewRecorder()

			newProtected(t).ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, "ok", rec.Body.String())
		})
	}
}

func TestBearerAuth_DisabledWhenKeyEmpty(t *testing.T) {
	cases := map[string]string{
		"no header":     "",
		"random bearer": "Bearer whatever",
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := transport.BearerAuth(next, "")

	for name, header := range cases {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Set("Authorization", header)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			require.Equal(t, http.StatusOK, rec.Code)
		})
	}
}

func TestBearerAuth_Rejects(t *testing.T) {
	cases := map[string]string{
		"missing header": "",
		"wrong key":      "Bearer wrong",
		"no prefix":      testKey,
		"empty bearer":   "Bearer ",
	}

	for name, header := range cases {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			req.Header.Set("Authorization", header)
			rec := httptest.NewRecorder()

			newProtected(t).ServeHTTP(rec, req)

			require.Equal(t, http.StatusUnauthorized, rec.Code)
			assert.Equal(t, "Bearer", rec.Header().Get("WWW-Authenticate"))
		})
	}
}
