package youtrack

import (
	"context"
	"fmt"
	"net/url"
)

func (c *Client) ArticleList(ctx context.Context, skip, top int) ([]Article, error) {
	q := url.Values{}
	q.Set("fields", ArticleListFields)
	if skip > 0 {
		q.Set("$skip", fmt.Sprintf("%d", skip))
	}
	if top > 0 {
		q.Set("$top", fmt.Sprintf("%d", top))
	}

	var out []Article
	if err := c.get(ctx, "/api/articles", q, &out); err != nil {
		return nil, fmt.Errorf("article list: %w", err)
	}
	return out, nil
}

func (c *Client) ArticleGet(ctx context.Context, id string) (Article, error) {
	q := url.Values{}
	q.Set("fields", ArticleDetailFields)

	var out Article
	if err := c.get(ctx, apiPath("api", "articles", id), q, &out); err != nil {
		return Article{}, fmt.Errorf("article get: %w", err)
	}
	return out, nil
}

func (c *Client) ArticleCreate(ctx context.Context, article Article) (Article, error) {
	q := url.Values{}
	q.Set("fields", ArticleDetailFields)

	var out Article
	if err := c.post(ctx, "/api/articles", q, article, &out); err != nil {
		return Article{}, fmt.Errorf("article create: %w", err)
	}
	return out, nil
}

func (c *Client) ArticleUpdate(ctx context.Context, id string, article Article) (Article, error) {
	q := url.Values{}
	q.Set("fields", ArticleDetailFields)

	var out Article
	if err := c.post(ctx, apiPath("api", "articles", id), q, article, &out); err != nil {
		return Article{}, fmt.Errorf("article update: %w", err)
	}
	return out, nil
}

func (c *Client) ArticleDelete(ctx context.Context, id string) error {
	if err := c.delete(ctx, apiPath("api", "articles", id), nil); err != nil {
		return fmt.Errorf("article delete: %w", err)
	}
	return nil
}
