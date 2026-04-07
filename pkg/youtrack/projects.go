package youtrack

import (
	"context"
	"fmt"
	"net/url"
)

func (c *Client) ProjectList(ctx context.Context, skip, top int) ([]Project, error) {
	q := url.Values{}
	q.Set("fields", ProjectFields)
	if skip > 0 {
		q.Set("$skip", fmt.Sprintf("%d", skip))
	}
	if top > 0 {
		q.Set("$top", fmt.Sprintf("%d", top))
	}

	var out []Project
	if err := c.get(ctx, "/api/admin/projects", q, &out); err != nil {
		return nil, fmt.Errorf("project list: %w", err)
	}
	return out, nil
}

func (c *Client) ProjectGet(ctx context.Context, id string) (Project, error) {
	q := url.Values{}
	q.Set("fields", ProjectFields)

	var out Project
	if err := c.get(ctx, apiPath("api", "admin", "projects", id), q, &out); err != nil {
		return Project{}, fmt.Errorf("project get: %w", err)
	}
	return out, nil
}
