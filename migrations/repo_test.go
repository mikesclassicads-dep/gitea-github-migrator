package migrations

import (
	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMigratory_Repository(t *testing.T) {
	fm := &FetchMigratory{
		Migratory: *DemoMigratory,
		GHClient:  github.NewClient(nil),
		RepoOwner: "JonasFranzDEV",
		RepoName:  "migrate",
	}
	ghRepo, _, err := fm.GHClient.Repositories.Get(fm.ctx(), fm.RepoOwner, fm.RepoName)
	if err != nil {
		t.Skipf("Skipped due to repo is not accessable: %v", err)
		return
	}
	repo, err := fm.Repository(ghRepo)
	assertNoError(t, err)
	assert.Equal(t, "migrate", repo.Name)
	assertNoError(t, fm.Client.DeleteRepo("demo", "migrate"))
}
