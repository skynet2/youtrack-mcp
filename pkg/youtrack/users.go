package youtrack

import (
	"context"
	"fmt"
	"net/url"
)

func (c *Client) UserList(ctx context.Context, query string, skip, top int) ([]User, error) {
	q := url.Values{}
	q.Set("fields", UserFields)
	if query != "" {
		q.Set("query", query)
	}
	if skip > 0 {
		q.Set("$skip", fmt.Sprintf("%d", skip))
	}
	if top > 0 {
		q.Set("$top", fmt.Sprintf("%d", top))
	}

	var out []User
	if err := c.get(ctx, "/api/users", q, &out); err != nil {
		return nil, fmt.Errorf("user list: %w", err)
	}
	return out, nil
}

func (c *Client) UserCurrent(ctx context.Context) (User, error) {
	q := url.Values{}
	q.Set("fields", UserFields)

	var out User
	if err := c.get(ctx, "/api/users/me", q, &out); err != nil {
		return User{}, fmt.Errorf("user current: %w", err)
	}
	return out, nil
}
