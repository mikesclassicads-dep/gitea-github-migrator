package migrations

import "code.gitea.io/sdk/gitea"

// MigratoryStage represents the actual step in the process
type MigratoryStage int

const (
	// Importing imports the repo to Gitea
	Importing MigratoryStage = iota
	// Migrating migrates issues, etc to Gitea
	Migrating
	// Finished means that everything is migrated successfully
	Finished
	// Failed is only entered if a fatal error occurs
	Failed
)

// Migratory is the context for migrating things from GitHub to Gitea
type Migratory struct {
	Options
	Client *gitea.Client

	Status *MigratoryStatus

	repository *gitea.Repository
	// key: github milestone id | value: gitea milestone id
	migratedMilestones map[int64]int64
	// key: github label id | value: gitea label id
	migratedLabels map[int64]int64
}

// MigratoryStatus represents the actual state of a migratory
type MigratoryStatus struct {
	Stage MigratoryStage `json:"stage"`

	Issues         int64 `json:"total_issues"`
	IssuesMigrated int64 `json:"migrated_issues"`
	IssuesError    int64 `json:"failed_issues"`

	Comments         int64 `json:"total_comments"`
	CommentsError    int64 `json:"failed_comments"`
	CommentsMigrated int64 `json:"migrated_comments"`

	// FatalError should only be used if stage == failed; indicates which fatal error occurred
	FatalError error
	Log        string `json:"log"`
}
