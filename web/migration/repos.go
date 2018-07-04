package migration

import (
	"regexp"
	"strings"

	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/migrations"
	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/web/context"
	"github.com/google/go-github/github"

	bgctx "context"
)

const repoRegex = "^[A-Za-z0-9-.]+/[A-Za-z0-9-.]+$"

// ListRepos shows all available repos of the signed in user
func ListRepos(ctx *context.Context) {
	repos, _, err := ctx.Client.Repositories.List(bgctx.Background(), ctx.User.Username, &github.RepositoryListOptions{})
	if err != nil {
		ctx.Handle(500, "list repositories", err)
		return
	}
	ctx.Data["Repos"] = repos
	ctx.HTML(200, "repos")
}

// ListReposPost handles the form submission of ListRepos
func ListReposPost(ctx *context.Context) {
	if err := ctx.Req.ParseForm(); err != nil {
		ctx.Handle(500, "parse form", err)
		return
	}
	// TODO implement migration options
	job := migrations.NewJob(&migrations.Options{
		Labels:       true,
		Comments:     true,
		Issues:       true,
		Milestones:   true,
		PullRequests: true,
		Strategy:     migrations.Classic,
		NewOwnerID:   int(ctx.GiteaUser.ID), // TODO implement user/org selection
	}, ctx.GiteaClient, ctx.Client)
	for repo, val := range ctx.Req.Form {
		activated := strings.Join(val, "")
		if activated != "on" {
			continue
		}
		// Validate repo format (reponame/owner)
		if matched, err := regexp.MatchString(repoRegex, repo); err != nil || !matched {
			continue
		}
		job.Repositories = append(job.Repositories, repo)
	}
	go job.StartMigration()
	ctx.SetCurrentJob(job)
	ctx.Data["StatusReport"] = job.StatusReport()
	ctx.HTML(200, "migration")
}
