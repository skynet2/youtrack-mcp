package youtrack

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

func resolveLink(types []IssueLinkType, phrase string) (string, error) {
	p := strings.ToLower(strings.TrimSpace(phrase))
	for _, lt := range types {
		if lt.Directed {
			if strings.ToLower(lt.SourceToTarget) == p {
				return lt.ID + "s", nil
			}
			if strings.ToLower(lt.TargetToSource) == p {
				return lt.ID + "t", nil
			}
			if strings.ToLower(lt.Name) == p {
				return "", fmt.Errorf("%w: %q is a directed link type; use %q or %q", ErrLinkTypeNotFound, lt.Name, lt.SourceToTarget, lt.TargetToSource)
			}
			continue
		}
		if strings.ToLower(lt.Name) == p || strings.ToLower(lt.SourceToTarget) == p {
			return lt.ID, nil
		}
	}
	return "", fmt.Errorf("%w: %q; valid phrases: %s", ErrLinkTypeNotFound, phrase, validPhrases(types))
}

func validPhrases(types []IssueLinkType) string {
	var out []string
	for _, lt := range types {
		if lt.Directed {
			out = append(out, fmt.Sprintf("%q", lt.SourceToTarget), fmt.Sprintf("%q", lt.TargetToSource))
			continue
		}
		out = append(out, fmt.Sprintf("%q", lt.Name))
	}
	return strings.Join(out, ", ")
}

func (c *Client) ListIssueLinkTypes(ctx context.Context) ([]IssueLinkType, error) {
	q := url.Values{}
	q.Set("fields", IssueLinkTypeFields)

	var out []IssueLinkType
	if err := c.get(ctx, "/api/issueLinkTypes", q, &out); err != nil {
		return nil, fmt.Errorf("list issue link types: %w", err)
	}
	return out, nil
}

func (c *Client) IssueLink(ctx context.Context, sourceID, linkPhrase, targetID string) error {
	linkID, targetDBID, err := c.resolveLinkArgs(ctx, linkPhrase, targetID)
	if err != nil {
		return err
	}

	body := map[string]string{"id": targetDBID}
	path := apiPath("api", "issues", sourceID, "links", linkID, "issues")
	if err := c.post(ctx, path, nil, body, nil); err != nil {
		return fmt.Errorf("issue link: %w", err)
	}
	return nil
}

func (c *Client) IssueUnlink(ctx context.Context, sourceID, linkPhrase, targetID string) error {
	linkID, targetDBID, err := c.resolveLinkArgs(ctx, linkPhrase, targetID)
	if err != nil {
		return err
	}

	path := apiPath("api", "issues", sourceID, "links", linkID, "issues", targetDBID)
	if err := c.delete(ctx, path, nil); err != nil {
		return fmt.Errorf("issue unlink: %w", err)
	}
	return nil
}

func (c *Client) resolveLinkArgs(ctx context.Context, linkPhrase, targetID string) (string, string, error) {
	types, err := c.ListIssueLinkTypes(ctx)
	if err != nil {
		return "", "", fmt.Errorf("resolve link: %w", err)
	}
	linkID, err := resolveLink(types, linkPhrase)
	if err != nil {
		return "", "", err
	}
	target, err := c.IssueGet(ctx, targetID)
	if err != nil {
		return "", "", fmt.Errorf("resolve target: %w", err)
	}
	return linkID, target.ID, nil
}
