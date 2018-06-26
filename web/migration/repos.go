package migration

import (
	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/web/context"
	"github.com/google/go-github/github"

	bgctx "context"
)

// ListRepos shows all available repos of the signed in user
func ListRepos(ctx *context.Context) {
	repos, _, err := ctx.Client.Repositories.List(bgctx.Background(), ctx.User.Username, &github.RepositoryListOptions{})
	if err != nil {
		ctx.Handle(500, "list repositories", err)
	}
	ctx.Data["Repos"] = repos
	ctx.HTML(200, "repos")
}
