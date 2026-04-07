package youtrack

import (
	"context"
	"fmt"
	"net/url"
)

func (c *Client) TagList(ctx context.Context, skip, top int) ([]IssueTag, error) {
	q := url.Values{}
	q.Set("fields", TagFields)
	if skip > 0 {
		q.Set("$skip", fmt.Sprintf("%d", skip))
	}
	if top > 0 {
		q.Set("$top", fmt.Sprintf("%d", top))
	}

	var out []IssueTag
	if err := c.get(ctx, "/api/issueTags", q, &out); err != nil {
		return nil, fmt.Errorf("tag list: %w", err)
	}
	return out, nil
}
