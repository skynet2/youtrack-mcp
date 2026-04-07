package youtrack

import (
	"context"
	"fmt"
	"net/url"
)

func (c *Client) WorkItemList(ctx context.Context, issueID string, skip, top int) ([]IssueWorkItem, error) {
	q := url.Values{}
	q.Set("fields", WorkItemFields)
	if skip > 0 {
		q.Set("$skip", fmt.Sprintf("%d", skip))
	}
	if top > 0 {
		q.Set("$top", fmt.Sprintf("%d", top))
	}

	var out []IssueWorkItem
	if err := c.get(ctx, apiPath("api", "issues", issueID, "timeTracking", "workItems"), q, &out); err != nil {
		return nil, fmt.Errorf("workitem list: %w", err)
	}
	return out, nil
}

func (c *Client) WorkItemCreate(ctx context.Context, issueID string, item IssueWorkItem) (IssueWorkItem, error) {
	q := url.Values{}
	q.Set("fields", WorkItemFields)

	var out IssueWorkItem
	if err := c.post(ctx, apiPath("api", "issues", issueID, "timeTracking", "workItems"), q, item, &out); err != nil {
		return IssueWorkItem{}, fmt.Errorf("workitem create: %w", err)
	}
	return out, nil
}
