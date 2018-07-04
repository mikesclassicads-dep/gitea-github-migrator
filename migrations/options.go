package migrations

// Options defines the way a repository gets migrated
type Options struct {
	Issues       bool
	Milestones   bool
	Labels       bool
	Comments     bool
	PullRequests bool

	AuthUsername string
	AuthPassword string

	Private    bool
	NewOwnerID int

	Strategy Strategy
}

// Strategy represents the procedure of migration.
type Strategy int

const (
	// Classic works for all Gitea versions and creates comments by the user migrating the repository. This does not require
	// admin permissions. The issue "number" is also assinged by Gitea and could be different to the GitHub issue "number".
	// Creation date of comments, issues, milestones, etc. will be the date of creation.
	Classic Strategy = iota
	// Advanced works for all Gitea versions 1.6+ and utilizes the Gitea Migration API which allows the tool to create comments
	// with Ghost Users. Creation date and issue numbers will be the same like GitHub. It requires admin permissions for repo
	// (creation date, issue number) and/or
	Advanced
)
