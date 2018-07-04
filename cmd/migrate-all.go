package cmd

import (
	"context"
	"fmt"

	"code.gitea.io/sdk/gitea"
	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/migrations"
	"github.com/google/go-github/github"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
)

// CmdMigrateAll is a command to migrate all repositories of an user
var CmdMigrateAll = cli.Command{
	Name:   "migrate-all",
	Usage:  "migrates all repositories of an user from github to a gitea repository",
	Action: runMigrateAll,
	Flags: append(defaultMigrateFlags,
		cli.StringFlag{
			Name:   "gh-user",
			EnvVar: "GH_USER",
			Usage:  "GitHub Username",
		},
	),
}

func runMigrateAll(ctx *cli.Context) error {
	onlyRepos := ctx.Bool("only-repo")
	var gc *github.Client
	if ctx.IsSet("gh-token") {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: ctx.String("gh-token")},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		gc = github.NewClient(tc)
	} else {
		gc = github.NewClient(nil)
	}

	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := gc.Repositories.List(context.Background(), ctx.String("gh-user"), opt)
		if err != nil {
			return err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	job := migrations.NewJob(&migrations.Options{
		Private:    ctx.Bool("private"),
		NewOwnerID: ctx.Int("owner"),

		Comments:     !onlyRepos,
		Issues:       !onlyRepos,
		Labels:       !onlyRepos,
		Milestones:   !onlyRepos,
		PullRequests: !onlyRepos,
		Strategy:     migrations.Classic,
	}, gitea.NewClient(ctx.String("url"), ctx.String("token")), gc)
	if job.Options.NewOwnerID == 0 {
		usr, err := job.Client.GetMyUserInfo()
		if err != nil {
			return fmt.Errorf("cannot fetch user info about current user: %v", err)
		}
		job.Options.NewOwnerID = int(usr.ID)
	}
	for _, repo := range allRepos {
		job.Repositories = append(job.Repositories, repo.GetFullName())
	}
	errs := job.StartMigration()
	for i := range errs {
		if i != nil {
			fmt.Printf("error: %v\n", i)
		}
	}
	return nil
}
