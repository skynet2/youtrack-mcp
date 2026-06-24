package youtrack

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func linkTypeFixtures() []IssueLinkType {
	return []IssueLinkType{
		{ID: "109-0", Name: "Relates", SourceToTarget: "relates to", Directed: false},
		{ID: "109-3", Name: "Subtask", SourceToTarget: "parent for", TargetToSource: "subtask of", Directed: true},
	}
}

func TestResolveLink_Success(t *testing.T) {
	tests := []struct {
		name   string
		phrase string
		want   string
	}{
		{name: "outward directed phrase", phrase: "parent for", want: "109-3s"},
		{name: "inward directed phrase", phrase: "subtask of", want: "109-3t"},
		{name: "undirected name", phrase: "Relates", want: "109-0"},
		{name: "undirected sourceToTarget phrase", phrase: "relates to", want: "109-0"},
		{name: "case insensitive", phrase: "SubTask Of", want: "109-3t"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveLink(linkTypeFixtures(), tt.phrase)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestResolveLink_Failure(t *testing.T) {
	tests := []struct {
		name   string
		phrase string
	}{
		{name: "unknown phrase", phrase: "no such phrase"},
		{name: "directed type name", phrase: "Subtask"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := resolveLink(linkTypeFixtures(), tt.phrase)
			require.ErrorIs(t, err, ErrLinkTypeNotFound)
		})
	}
}

func TestListIssueLinkTypes_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/issueLinkTypes", r.URL.Path)
		assert.Equal(t, IssueLinkTypeFields, r.URL.Query().Get("fields"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"$type":"IssueLinkType","id":"109-3","name":"Subtask","sourceToTarget":"parent for","targetToSource":"subtask of","directed":true}]`))
	})

	out, err := client.ListIssueLinkTypes(context.Background())
	require.NoError(t, err)
	require.Len(t, out, 1)
	assert.Equal(t, "Subtask", out[0].Name)
	assert.Equal(t, "subtask of", out[0].TargetToSource)
	assert.True(t, out[0].Directed)
}

func TestListIssueLinkTypes_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})
	out, err := client.ListIssueLinkTypes(context.Background())
	require.ErrorIs(t, err, ErrUnauthorized)
	require.Nil(t, out)
}

func TestIssueLink_Success(t *testing.T) {
	var postHit bool
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/issueLinkTypes":
			_, _ = w.Write([]byte(`[{"id":"109-3","name":"Subtask","sourceToTarget":"parent for","targetToSource":"subtask of","directed":true}]`))
		case "/api/issues/IT-8":
			_, _ = w.Write([]byte(`{"id":"2-2517","idReadable":"IT-8"}`))
		case "/api/issues/IT-9/links/109-3t/issues":
			postHit = true
			assert.Equal(t, http.MethodPost, r.Method)
			body, _ := io.ReadAll(r.Body)
			assert.JSONEq(t, `{"id":"2-2517"}`, string(body))
			_, _ = w.Write([]byte(`{"id":"2-2517","idReadable":"IT-8"}`))
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	})

	err := client.IssueLink(context.Background(), "IT-9", "subtask of", "IT-8")
	require.NoError(t, err)
	assert.True(t, postHit)
}

func TestIssueUnlink_Success(t *testing.T) {
	var delHit bool
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/issueLinkTypes":
			_, _ = w.Write([]byte(`[{"id":"109-3","name":"Subtask","sourceToTarget":"parent for","targetToSource":"subtask of","directed":true}]`))
		case "/api/issues/IT-8":
			_, _ = w.Write([]byte(`{"id":"2-2517","idReadable":"IT-8"}`))
		case "/api/issues/IT-9/links/109-3t/issues/2-2517":
			delHit = true
			assert.Equal(t, http.MethodDelete, r.Method)
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	})

	err := client.IssueUnlink(context.Background(), "IT-9", "subtask of", "IT-8")
	require.NoError(t, err)
	assert.True(t, delHit)
}

func TestIssueLink_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`[{"id":"109-3","name":"Subtask","sourceToTarget":"parent for","targetToSource":"subtask of","directed":true}]`))
	})
	err := client.IssueLink(context.Background(), "IT-9", "no such phrase", "IT-8")
	require.ErrorIs(t, err, ErrLinkTypeNotFound)
}

func TestIssueLink_TargetResolveFailure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/issueLinkTypes":
			_, _ = w.Write([]byte(`[{"id":"109-3","name":"Subtask","sourceToTarget":"parent for","targetToSource":"subtask of","directed":true}]`))
		case "/api/issues/IT-8":
			w.WriteHeader(http.StatusNotFound)
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	})

	err := client.IssueLink(context.Background(), "IT-9", "subtask of", "IT-8")
	require.ErrorIs(t, err, ErrNotFound)
}
