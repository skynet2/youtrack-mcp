package youtrack

const (
	IssueListFields     = "$type,id,idReadable,summary,created,updated,resolved,project($type,id,name,shortName),reporter($type,id,login,fullName),customFields($type,id,name,value($type,id,name,text,login,fullName,minutes,presentation))"
	IssueDetailFields   = "$type,id,idReadable,summary,description,created,updated,resolved,numberInProject,project($type,id,name,shortName),reporter($type,id,login,fullName),updater($type,id,login,fullName),customFields($type,id,name,value($type,id,name,text,login,fullName,minutes,presentation)),tags($type,id,name),links($type,id,direction,linkType($type,id,name,localizedName),issues($type,id,idReadable,summary)),votes"
	ProjectFields       = "$type,id,name,shortName,description,archived,leader($type,id,login,fullName)"
	CommentFields       = "$type,id,text,textPreview,created,updated,author($type,id,login,fullName),deleted,usesMarkdown"
	TagFields           = "$type,id,name,owner($type,id,login,fullName),color($type,id,background,foreground),untagOnResolve"
	WorkItemFields      = "$type,id,author($type,id,login,fullName),creator($type,id,login,fullName),text,type($type,id,name),created,updated,duration($type,id,minutes,presentation),date,usesMarkdown"
	AgileFields         = "$type,id,name,owner($type,id,login,fullName),projects($type,id,name,shortName),currentSprint($type,id,name)"
	SprintFields        = "$type,id,name,goal,start,finish,archived,isDefault,unresolvedIssuesCount"
	SprintDetailFields  = "$type,id,name,goal,start,finish,archived,isDefault,issues($type,id,idReadable,summary),unresolvedIssuesCount"
	UserFields          = "$type,id,login,fullName,email,ringId,guest,online,banned"
	CustomFieldFields   = "$type,id,name,projectCustomField($type,id,field($type,id,name,localizedName,fieldType($type,id))),value($type,id,name,text,login,fullName,minutes,presentation,color($type,id,background,foreground))"
	ArticleListFields   = "$type,id,idReadable,summary,created,updated,project($type,id,name,shortName),reporter($type,id,login,fullName),hasChildren"
	ArticleDetailFields = "$type,id,idReadable,summary,content,created,updated,project($type,id,name,shortName),reporter($type,id,login,fullName),updatedBy($type,id,login,fullName),parentArticle($type,id,idReadable,summary),childArticles($type,id,idReadable,summary),hasChildren,visibility($type,id)"
	ActivityFields      = "$type,id,author($type,id,login,fullName),timestamp,added,removed,target,targetMember,field($type,id,name,customField($type,id,name)),category($type,id)"
	IssueLinkFields     = "$type,id,direction,linkType($type,id,name,localizedName),issues($type,id,idReadable,summary)"
	IssueLinkTypeFields = "$type,id,name,localizedName,sourceToTarget,targetToSource,directed"
)
