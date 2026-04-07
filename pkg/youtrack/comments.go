package youtrack

import (
	"context"
	"fmt"
	"net/url"
)

func (c *Client) CommentList(ctx context.Context, issueID string, skip, top int) ([]IssueComment, error) {
	q := url.Values{}
	q.Set("fields", CommentFields)
	if skip > 0 {
		q.Set("$skip", fmt.Sprintf("%d", skip))
	}
	if top > 0 {
		q.Set("$top", fmt.Sprintf("%d", top))
	}

	var out []IssueComment
	if err := c.get(ctx, apiPath("api", "issues", issueID, "comments"), q, &out); err != nil {
		return nil, fmt.Errorf("comment list: %w", err)
	}
	return out, nil
}

func (c *Client) CommentCreate(ctx context.Context, issueID string, comment IssueComment) (IssueComment, error) {
	q := url.Values{}
	q.Set("fields", CommentFields)

	var out IssueComment
	if err := c.post(ctx, apiPath("api", "issues", issueID, "comments"), q, comment, &out); err != nil {
		return IssueComment{}, fmt.Errorf("comment create: %w", err)
	}
	return out, nil
}

func (c *Client) CommentUpdate(ctx context.Context, issueID, commentID string, comment IssueComment) (IssueComment, error) {
	q := url.Values{}
	q.Set("fields", CommentFields)

	var out IssueComment
	if err := c.post(ctx, apiPath("api", "issues", issueID, "comments", commentID), q, comment, &out); err != nil {
		return IssueComment{}, fmt.Errorf("comment update: %w", err)
	}
	return out, nil
}

func (c *Client) CommentDelete(ctx context.Context, issueID, commentID string) error {
	if err := c.delete(ctx, apiPath("api", "issues", issueID, "comments", commentID), nil); err != nil {
		return fmt.Errorf("comment delete: %w", err)
	}
	return nil
}
