package migrations

import "code.gitea.io/sdk/gitea"

// Migratory is the context for migrating things from GitHub to Gitea
type Migratory struct {
	Client       *gitea.Client
	AuthUsername string
	AuthPassword string

	Private    bool
	NewOwnerID int64

	repository *gitea.Repository
	// key: github milestone id | value: gitea milestone id
	migratedMilestones map[int64]int64
	// key: github label id | value: gitea label id
	migratedLabels map[int64]int64
}
