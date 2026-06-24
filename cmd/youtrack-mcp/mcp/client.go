package mcp

import (
	"context"

	"github.com/skynet2/youtrack-mcp/pkg/youtrack"
)

//go:generate mockgen -source=client.go -destination=mocks/mock_client.go -package=mocks

type YouTrackClient interface {
	// Issues
	IssueSearch(ctx context.Context, query string, skip, top int) ([]youtrack.Issue, error)
	IssueGet(ctx context.Context, id string) (youtrack.Issue, error)
	IssueCreate(ctx context.Context, issue youtrack.Issue) (youtrack.Issue, error)
	IssueUpdate(ctx context.Context, id string, issue youtrack.Issue) (youtrack.Issue, error)
	IssueDelete(ctx context.Context, id string) error
	IssueAddTag(ctx context.Context, issueID, tagID string) error
	IssueRemoveTag(ctx context.Context, issueID, tagID string) error
	IssueGetLinks(ctx context.Context, issueID string) ([]youtrack.IssueLink, error)
	ListIssueLinkTypes(ctx context.Context) ([]youtrack.IssueLinkType, error)
	IssueLink(ctx context.Context, sourceID, linkPhrase, targetID string) error
	IssueUnlink(ctx context.Context, sourceID, linkPhrase, targetID string) error
	// Projects
	ProjectList(ctx context.Context, skip, top int) ([]youtrack.Project, error)
	ProjectGet(ctx context.Context, id string) (youtrack.Project, error)
	// Comments
	CommentList(ctx context.Context, issueID string, skip, top int) ([]youtrack.IssueComment, error)
	CommentCreate(ctx context.Context, issueID string, comment youtrack.IssueComment) (youtrack.IssueComment, error)
	CommentUpdate(ctx context.Context, issueID, commentID string, comment youtrack.IssueComment) (youtrack.IssueComment, error)
	CommentDelete(ctx context.Context, issueID, commentID string) error
	// Tags
	TagList(ctx context.Context, skip, top int) ([]youtrack.IssueTag, error)
	// Time Tracking
	WorkItemList(ctx context.Context, issueID string, skip, top int) ([]youtrack.IssueWorkItem, error)
	WorkItemCreate(ctx context.Context, issueID string, item youtrack.IssueWorkItem) (youtrack.IssueWorkItem, error)
	// Agile
	AgileList(ctx context.Context, skip, top int) ([]youtrack.Agile, error)
	SprintList(ctx context.Context, agileID string, skip, top int) ([]youtrack.Sprint, error)
	SprintGet(ctx context.Context, agileID, sprintID string) (youtrack.Sprint, error)
	// Users
	UserList(ctx context.Context, query string, skip, top int) ([]youtrack.User, error)
	UserCurrent(ctx context.Context) (youtrack.User, error)
	// Custom Fields
	IssueCustomFields(ctx context.Context, issueID string) ([]youtrack.IssueCustomField, error)
	IssueUpdateField(ctx context.Context, issueID, fieldID string, field youtrack.IssueCustomField) (youtrack.IssueCustomField, error)
	// Knowledge Base
	ArticleList(ctx context.Context, skip, top int) ([]youtrack.Article, error)
	ArticleGet(ctx context.Context, id string) (youtrack.Article, error)
	ArticleCreate(ctx context.Context, article youtrack.Article) (youtrack.Article, error)
	ArticleUpdate(ctx context.Context, id string, article youtrack.Article) (youtrack.Article, error)
	ArticleDelete(ctx context.Context, id string) error
	// Activities
	ActivityList(ctx context.Context, categories string, issueQuery string, start, end string, skip, top int) ([]youtrack.ActivityItem, error)
	IssueActivityList(ctx context.Context, issueID string, categories string, skip, top int) ([]youtrack.ActivityItem, error)
}

var _ YouTrackClient = (*youtrack.Client)(nil)
