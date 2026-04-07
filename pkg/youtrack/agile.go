package youtrack

import (
	"context"
	"fmt"
	"net/url"
)

func (c *Client) AgileList(ctx context.Context, skip, top int) ([]Agile, error) {
	q := url.Values{}
	q.Set("fields", AgileFields)
	if skip > 0 {
		q.Set("$skip", fmt.Sprintf("%d", skip))
	}
	if top > 0 {
		q.Set("$top", fmt.Sprintf("%d", top))
	}

	var out []Agile
	if err := c.get(ctx, "/api/agiles", q, &out); err != nil {
		return nil, fmt.Errorf("agile list: %w", err)
	}
	return out, nil
}

func (c *Client) SprintList(ctx context.Context, agileID string, skip, top int) ([]Sprint, error) {
	q := url.Values{}
	q.Set("fields", SprintFields)
	if skip > 0 {
		q.Set("$skip", fmt.Sprintf("%d", skip))
	}
	if top > 0 {
		q.Set("$top", fmt.Sprintf("%d", top))
	}

	var out []Sprint
	if err := c.get(ctx, apiPath("api", "agiles", agileID, "sprints"), q, &out); err != nil {
		return nil, fmt.Errorf("sprint list: %w", err)
	}
	return out, nil
}

func (c *Client) SprintGet(ctx context.Context, agileID, sprintID string) (Sprint, error) {
	q := url.Values{}
	q.Set("fields", SprintDetailFields)

	var out Sprint
	if err := c.get(ctx, apiPath("api", "agiles", agileID, "sprints", sprintID), q, &out); err != nil {
		return Sprint{}, fmt.Errorf("sprint get: %w", err)
	}
	return out, nil
}
