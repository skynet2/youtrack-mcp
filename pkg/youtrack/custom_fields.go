package youtrack

import (
	"context"
	"fmt"
	"net/url"
)

func (c *Client) IssueCustomFields(ctx context.Context, issueID string) ([]IssueCustomField, error) {
	q := url.Values{}
	q.Set("fields", CustomFieldFields)

	var out []IssueCustomField
	if err := c.get(ctx, apiPath("api", "issues", issueID, "customFields"), q, &out); err != nil {
		return nil, fmt.Errorf("issue custom fields: %w", err)
	}
	return out, nil
}

func (c *Client) IssueUpdateField(ctx context.Context, issueID, fieldID string, field IssueCustomField) (IssueCustomField, error) {
	q := url.Values{}
	q.Set("fields", CustomFieldFields)

	var out IssueCustomField
	if err := c.post(ctx, apiPath("api", "issues", issueID, "customFields", fieldID), q, field, &out); err != nil {
		return IssueCustomField{}, fmt.Errorf("issue update field: %w", err)
	}
	return out, nil
}
