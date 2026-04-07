package youtrack

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActivityList_Success(t *testing.T) {
	tests := []struct {
		name       string
		categories string
		issueQuery string
		start      string
		end        string
		skip       int
		top        int
	}{
		{
			name:       "with all parameters",
			categories: "IssueCreatedCategory",
			issueQuery: "project: DEMO",
			start:      "1700000000000",
			end:        "1700100000000",
			skip:       0,
			top:        50,
		},
		{
			name:       "with no optional parameters",
			categories: "",
			issueQuery: "",
			start:      "",
			end:        "",
			skip:       0,
			top:        0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "/api/activities", r.URL.Path)
				assert.Equal(t, ActivityFields, r.URL.Query().Get("fields"))

				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`[{"$type":"IssueCreatedActivityItem","id":"act-1","author":{"$type":"User","id":"1-1","login":"admin","fullName":"Admin User"},"timestamp":1700000000000,"field":{"$type":"ActivityField","id":"f-1","name":"created"},"category":{"$type":"ActivityCategory","id":"IssueCreatedCategory"}}]`))
			})

			result, err := client.ActivityList(context.Background(), tt.categories, tt.issueQuery, tt.start, tt.end, tt.skip, tt.top)
			require.NoError(t, err)
			require.Len(t, result, 1)
			assert.Equal(t, "act-1", result[0].ID)
			assert.Equal(t, "admin", result[0].Author.Login)
			require.NotNil(t, result[0].Field)
			assert.Equal(t, "created", result[0].Field.Name)
			require.NotNil(t, result[0].Category)
			assert.Equal(t, "IssueCreatedCategory", result[0].Category.ID)
		})
	}
}

func TestActivityList_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	result, err := client.ActivityList(context.Background(), "", "", "", "", 0, 0)
	require.ErrorIs(t, err, ErrUnauthorized)
	require.Nil(t, result)
}

func TestIssueActivityList_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/issues/DEMO-1/activities", r.URL.Path)
		assert.Equal(t, ActivityFields, r.URL.Query().Get("fields"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"$type":"CustomFieldActivityItem","id":"act-2","author":{"$type":"User","id":"1-1","login":"dev"},"timestamp":1700050000000,"targetMember":"Priority","field":{"$type":"ActivityField","id":"f-2","name":"Priority"}}]`))
	})

	result, err := client.IssueActivityList(context.Background(), "DEMO-1", "", 0, 0)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "act-2", result[0].ID)
	assert.Equal(t, "dev", result[0].Author.Login)
	assert.Equal(t, "Priority", result[0].TargetMember)
}

func TestIssueActivityList_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	result, err := client.IssueActivityList(context.Background(), "NONEXIST-1", "", 0, 0)
	require.ErrorIs(t, err, ErrNotFound)
	require.Nil(t, result)
}
