package youtrack

type Issue struct {
	Type            string             `json:"$type,omitempty"`
	ID              string             `json:"id,omitempty"`
	IDReadable      string             `json:"idReadable,omitempty"`
	Summary         string             `json:"summary,omitempty"`
	Description     string             `json:"description,omitempty"`
	Created         *int64             `json:"created,omitempty"`
	Updated         *int64             `json:"updated,omitempty"`
	Resolved        *int64             `json:"resolved,omitempty"`
	NumberInProject *int64             `json:"numberInProject,omitempty"`
	Project         *Project           `json:"project,omitempty"`
	Reporter        *User              `json:"reporter,omitempty"`
	Updater         *User              `json:"updater,omitempty"`
	Tags            []IssueTag         `json:"tags,omitempty"`
	CustomFields    []IssueCustomField `json:"customFields,omitempty"`
	Links           []IssueLink        `json:"links,omitempty"`
	Comments        []IssueComment     `json:"comments,omitempty"`
	Attachments     []IssueAttachment  `json:"attachments,omitempty"`
	Votes           *int32             `json:"votes,omitempty"`
	Visibility      *Visibility        `json:"visibility,omitempty"`
}

type Project struct {
	Type        string `json:"$type,omitempty"`
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	ShortName   string `json:"shortName,omitempty"`
	Description string `json:"description,omitempty"`
	Archived    *bool  `json:"archived,omitempty"`
	Leader      *User  `json:"leader,omitempty"`
}

type User struct {
	Type     string `json:"$type,omitempty"`
	ID       string `json:"id,omitempty"`
	Login    string `json:"login,omitempty"`
	FullName string `json:"fullName,omitempty"`
	Email    string `json:"email,omitempty"`
	RingID   string `json:"ringId,omitempty"`
	Guest    *bool  `json:"guest,omitempty"`
	Online   *bool  `json:"online,omitempty"`
	Banned   *bool  `json:"banned,omitempty"`
}

type IssueComment struct {
	Type         string      `json:"$type,omitempty"`
	ID           string      `json:"id,omitempty"`
	Text         string      `json:"text,omitempty"`
	TextPreview  string      `json:"textPreview,omitempty"`
	Created      *int64      `json:"created,omitempty"`
	Updated      *int64      `json:"updated,omitempty"`
	Author       *User       `json:"author,omitempty"`
	Deleted      *bool       `json:"deleted,omitempty"`
	UsesMarkdown *bool       `json:"usesMarkdown,omitempty"`
	Visibility   *Visibility `json:"visibility,omitempty"`
}

type IssueTag struct {
	Type           string      `json:"$type,omitempty"`
	ID             string      `json:"id,omitempty"`
	Name           string      `json:"name,omitempty"`
	Color          *FieldStyle `json:"color,omitempty"`
	UntagOnResolve *bool       `json:"untagOnResolve,omitempty"`
	Owner          *User       `json:"owner,omitempty"`
}

type FieldStyle struct {
	Type       string `json:"$type,omitempty"`
	ID         string `json:"id,omitempty"`
	Background string `json:"background,omitempty"`
	Foreground string `json:"foreground,omitempty"`
}

type IssueCustomField struct {
	Type               string              `json:"$type,omitempty"`
	ID                 string              `json:"id,omitempty"`
	Name               string              `json:"name,omitempty"`
	ProjectCustomField *ProjectCustomField `json:"projectCustomField,omitempty"`
	Value              any                 `json:"value,omitempty"`
}

type ProjectCustomField struct {
	Type  string       `json:"$type,omitempty"`
	ID    string       `json:"id,omitempty"`
	Field *CustomField `json:"field,omitempty"`
}

type CustomField struct {
	Type          string     `json:"$type,omitempty"`
	ID            string     `json:"id,omitempty"`
	Name          string     `json:"name,omitempty"`
	LocalizedName string     `json:"localizedName,omitempty"`
	FieldType     *FieldType `json:"fieldType,omitempty"`
}

type FieldType struct {
	Type string `json:"$type,omitempty"`
	ID   string `json:"id,omitempty"`
}

type IssueLink struct {
	Type      string         `json:"$type,omitempty"`
	ID        string         `json:"id,omitempty"`
	Direction string         `json:"direction,omitempty"`
	LinkType  *IssueLinkType `json:"linkType,omitempty"`
	Issues    []Issue        `json:"issues,omitempty"`
}

type IssueLinkType struct {
	Type          string `json:"$type,omitempty"`
	ID            string `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	LocalizedName string `json:"localizedName,omitempty"`
}

type IssueAttachment struct {
	Type      string `json:"$type,omitempty"`
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	URL       string `json:"url,omitempty"`
	Size      *int64 `json:"size,omitempty"`
	MimeType  string `json:"mimeType,omitempty"`
	Extension string `json:"extension,omitempty"`
}

type Visibility struct {
	Type            string      `json:"$type,omitempty"`
	ID              string      `json:"id,omitempty"`
	PermittedGroups []UserGroup `json:"permittedGroups,omitempty"`
	PermittedUsers  []User      `json:"permittedUsers,omitempty"`
}

type UserGroup struct {
	Type   string `json:"$type,omitempty"`
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	RingID string `json:"ringId,omitempty"`
}

type IssueWorkItem struct {
	Type         string         `json:"$type,omitempty"`
	ID           string         `json:"id,omitempty"`
	Author       *User          `json:"author,omitempty"`
	Creator      *User          `json:"creator,omitempty"`
	Text         string         `json:"text,omitempty"`
	TextPreview  string         `json:"textPreview,omitempty"`
	WorkType     *WorkItemType  `json:"type,omitempty"`
	Created      *int64         `json:"created,omitempty"`
	Updated      *int64         `json:"updated,omitempty"`
	Duration     *DurationValue `json:"duration,omitempty"`
	Date         *int64         `json:"date,omitempty"`
	UsesMarkdown *bool          `json:"usesMarkdown,omitempty"`
}

type WorkItemType struct {
	Type string `json:"$type,omitempty"`
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type DurationValue struct {
	Type         string `json:"$type,omitempty"`
	ID           string `json:"id,omitempty"`
	Minutes      *int32 `json:"minutes,omitempty"`
	Presentation string `json:"presentation,omitempty"`
}

type Agile struct {
	Type          string    `json:"$type,omitempty"`
	ID            string    `json:"id,omitempty"`
	Name          string    `json:"name,omitempty"`
	Owner         *User     `json:"owner,omitempty"`
	Projects      []Project `json:"projects,omitempty"`
	Sprints       []Sprint  `json:"sprints,omitempty"`
	CurrentSprint *Sprint   `json:"currentSprint,omitempty"`
}

type Sprint struct {
	Type                  string  `json:"$type,omitempty"`
	ID                    string  `json:"id,omitempty"`
	Name                  string  `json:"name,omitempty"`
	Goal                  string  `json:"goal,omitempty"`
	Start                 *int64  `json:"start,omitempty"`
	Finish                *int64  `json:"finish,omitempty"`
	Archived              *bool   `json:"archived,omitempty"`
	IsDefault             *bool   `json:"isDefault,omitempty"`
	Issues                []Issue `json:"issues,omitempty"`
	UnresolvedIssuesCount *int32  `json:"unresolvedIssuesCount,omitempty"`
}

type Article struct {
	Type          string              `json:"$type,omitempty"`
	ID            string              `json:"id,omitempty"`
	IDReadable    string              `json:"idReadable,omitempty"`
	Summary       string              `json:"summary,omitempty"`
	Content       string              `json:"content,omitempty"`
	Reporter      *User               `json:"reporter,omitempty"`
	UpdatedBy     *User               `json:"updatedBy,omitempty"`
	Project       *Project            `json:"project,omitempty"`
	ParentArticle *Article            `json:"parentArticle,omitempty"`
	ChildArticles []Article           `json:"childArticles,omitempty"`
	HasChildren   *bool               `json:"hasChildren,omitempty"`
	Created       *int64              `json:"created,omitempty"`
	Updated       *int64              `json:"updated,omitempty"`
	Visibility    *Visibility         `json:"visibility,omitempty"`
	Attachments   []ArticleAttachment `json:"attachments,omitempty"`
}

type ArticleAttachment struct {
	Type     string `json:"$type,omitempty"`
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	URL      string `json:"url,omitempty"`
	Size     *int64 `json:"size,omitempty"`
	MimeType string `json:"mimeType,omitempty"`
}

type ActivityItem struct {
	Type         string            `json:"$type,omitempty"`
	ID           string            `json:"id,omitempty"`
	Author       *User             `json:"author,omitempty"`
	Timestamp    *int64            `json:"timestamp,omitempty"`
	Removed      any               `json:"removed,omitempty"`
	Added        any               `json:"added,omitempty"`
	Target       any               `json:"target,omitempty"`
	TargetMember string            `json:"targetMember,omitempty"`
	Field        *ActivityField    `json:"field,omitempty"`
	Category     *ActivityCategory `json:"category,omitempty"`
}

type ActivityField struct {
	Type        string       `json:"$type,omitempty"`
	ID          string       `json:"id,omitempty"`
	Name        string       `json:"name,omitempty"`
	CustomField *CustomField `json:"customField,omitempty"`
}

type ActivityCategory struct {
	Type string `json:"$type,omitempty"`
	ID   string `json:"id,omitempty"`
}
