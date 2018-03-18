package migrations

import (
	"code.gitea.io/sdk/gitea"
	"github.com/google/go-github/github"
)

func (m *Migratory) Repository(gr *github.Repository) (*gitea.Repository, error) {
	var err error
	m.repository, err = m.Client.MigrateRepo(gitea.MigrateRepoOption{
		Description:  *gr.Description,
		AuthPassword: m.AuthPassword,
		AuthUsername: m.AuthUsername,
		CloneAddr:    gr.GetCloneURL(),
		RepoName:     gr.GetName(),
		UID:          m.NewOwnerID,
		Private:      m.Private,
	})
	return m.repository, err
}
