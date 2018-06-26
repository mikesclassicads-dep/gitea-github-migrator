package migrations

import (
	"time"

	"code.gitea.io/sdk/gitea"
)

var DemoMigratory = &Migratory{
	AuthUsername: "demo",
	AuthPassword: "demo",
	Client:       gitea.NewClient("http://gitea:3000", "8bffa364d5a4b2f18421426da0baf6ccddd16d6b"),
	repository: &gitea.Repository{
		Name: "demo",
		Owner: &gitea.User{
			UserName: "demo",
		},
	},
	NewOwnerID:         1,
	migratedMilestones: make(map[int64]int64),
	migratedLabels:     make(map[int64]int64),
}

var demoTime = time.Date(2018, 01, 01, 01, 01, 01, 01, time.UTC)
