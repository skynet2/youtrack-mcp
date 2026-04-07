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

func TestIssueCustomFields_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/issues/DEMO-1/customFields", r.URL.Path)
		assert.Equal(t, CustomFieldFields, r.URL.Query().Get("fields"))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"$type":"SingleEnumIssueCustomField","id":"110-1","name":"Priority","projectCustomField":{"$type":"ProjectCustomField","id":"111-1","field":{"$type":"CustomField","id":"58-1","name":"Priority","fieldType":{"$type":"FieldType","id":"enum[1]"}}},"value":{"$type":"EnumBundleElement","id":"67-1","name":"Critical"}}]`))
	})

	result, err := client.IssueCustomFields(context.Background(), "DEMO-1")
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, "110-1", result[0].ID)
	assert.Equal(t, "Priority", result[0].Name)
	require.NotNil(t, result[0].ProjectCustomField)
	require.NotNil(t, result[0].ProjectCustomField.Field)
	assert.Equal(t, "Priority", result[0].ProjectCustomField.Field.Name)
}

func TestIssueCustomFields_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	result, err := client.IssueCustomFields(context.Background(), "NONEXIST-1")
	require.ErrorIs(t, err, ErrNotFound)
	require.Nil(t, result)
}

func TestIssueUpdateField_Success(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/api/issues/DEMO-1/customFields/110-1", r.URL.Path)
		assert.Equal(t, CustomFieldFields, r.URL.Query().Get("fields"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var reqBody map[string]any
		require.NoError(t, json.Unmarshal(body, &reqBody))

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"$type":"SingleEnumIssueCustomField","id":"110-1","name":"Priority","value":{"$type":"EnumBundleElement","id":"67-2","name":"Normal"}}`))
	})

	result, err := client.IssueUpdateField(context.Background(), "DEMO-1", "110-1", IssueCustomField{
		Value: map[string]any{"name": "Normal"},
	})
	require.NoError(t, err)
	assert.Equal(t, "110-1", result.ID)
	assert.Equal(t, "Priority", result.Name)
}

func TestIssueUpdateField_Failure(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := client.IssueUpdateField(context.Background(), "NONEXIST-1", "110-1", IssueCustomField{})
	require.ErrorIs(t, err, ErrNotFound)
}
