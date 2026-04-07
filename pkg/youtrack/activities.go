package youtrack

import (
	"context"
	"fmt"
	"net/url"
)

func (c *Client) ActivityList(ctx context.Context, categories string, issueQuery string, start, end string, skip, top int) ([]ActivityItem, error) {
	q := url.Values{}
	q.Set("fields", ActivityFields)
	if categories != "" {
		q.Set("categories", categories)
	}
	if issueQuery != "" {
		q.Set("issueQuery", issueQuery)
	}
	if start != "" {
		q.Set("start", start)
	}
	if end != "" {
		q.Set("end", end)
	}
	if skip > 0 {
		q.Set("$skip", fmt.Sprintf("%d", skip))
	}
	if top > 0 {
		q.Set("$top", fmt.Sprintf("%d", top))
	}

	var out []ActivityItem
	if err := c.get(ctx, "/api/activities", q, &out); err != nil {
		return nil, fmt.Errorf("activity list: %w", err)
	}
	return out, nil
}

func (c *Client) IssueActivityList(ctx context.Context, issueID string, categories string, skip, top int) ([]ActivityItem, error) {
	q := url.Values{}
	q.Set("fields", ActivityFields)
	if categories != "" {
		q.Set("categories", categories)
	}
	if skip > 0 {
		q.Set("$skip", fmt.Sprintf("%d", skip))
	}
	if top > 0 {
		q.Set("$top", fmt.Sprintf("%d", top))
	}

	var out []ActivityItem
	if err := c.get(ctx, apiPath("api", "issues", issueID, "activities"), q, &out); err != nil {
		return nil, fmt.Errorf("issue activity list: %w", err)
	}
	return out, nil
}
