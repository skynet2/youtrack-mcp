package youtrack

import (
	"context"
	"fmt"
	"net/url"
)

func (c *Client) IssueSearch(ctx context.Context, query string, skip, top int) ([]Issue, error) {
	q := url.Values{}
	q.Set("fields", IssueListFields)
	if query != "" {
		q.Set("query", query)
	}
	if skip > 0 {
		q.Set("$skip", fmt.Sprintf("%d", skip))
	}
	if top > 0 {
		q.Set("$top", fmt.Sprintf("%d", top))
	}

	var out []Issue
	if err := c.get(ctx, "/api/issues", q, &out); err != nil {
		return nil, fmt.Errorf("issue search: %w", err)
	}
	return out, nil
}

func (c *Client) IssueGet(ctx context.Context, id string) (Issue, error) {
	q := url.Values{}
	q.Set("fields", IssueDetailFields)

	var out Issue
	if err := c.get(ctx, apiPath("api", "issues", id), q, &out); err != nil {
		return Issue{}, fmt.Errorf("issue get: %w", err)
	}
	return out, nil
}

func (c *Client) IssueCreate(ctx context.Context, issue Issue) (Issue, error) {
	q := url.Values{}
	q.Set("fields", IssueDetailFields)

	var out Issue
	if err := c.post(ctx, "/api/issues", q, issue, &out); err != nil {
		return Issue{}, fmt.Errorf("issue create: %w", err)
	}
	return out, nil
}

func (c *Client) IssueUpdate(ctx context.Context, id string, issue Issue) (Issue, error) {
	q := url.Values{}
	q.Set("fields", IssueDetailFields)

	var out Issue
	if err := c.post(ctx, apiPath("api", "issues", id), q, issue, &out); err != nil {
		return Issue{}, fmt.Errorf("issue update: %w", err)
	}
	return out, nil
}

func (c *Client) IssueDelete(ctx context.Context, id string) error {
	if err := c.delete(ctx, apiPath("api", "issues", id), nil); err != nil {
		return fmt.Errorf("issue delete: %w", err)
	}
	return nil
}

func (c *Client) IssueAddTag(ctx context.Context, issueID, tagID string) error {
	body := IssueTag{ID: tagID}
	if err := c.post(ctx, apiPath("api", "issues", issueID, "tags"), nil, body, nil); err != nil {
		return fmt.Errorf("issue add tag: %w", err)
	}
	return nil
}

func (c *Client) IssueRemoveTag(ctx context.Context, issueID, tagID string) error {
	if err := c.delete(ctx, apiPath("api", "issues", issueID, "tags", tagID), nil); err != nil {
		return fmt.Errorf("issue remove tag: %w", err)
	}
	return nil
}

func (c *Client) IssueGetLinks(ctx context.Context, issueID string) ([]IssueLink, error) {
	q := url.Values{}
	q.Set("fields", IssueLinkFields)

	var out []IssueLink
	if err := c.get(ctx, apiPath("api", "issues", issueID, "links"), q, &out); err != nil {
		return nil, fmt.Errorf("issue get links: %w", err)
	}
	return out, nil
}
