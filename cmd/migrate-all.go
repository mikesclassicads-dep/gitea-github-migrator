package cmd

import (
	"context"
	"fmt"
	"sync"

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
	m := &migrations.Migratory{
		Client:     gitea.NewClient(ctx.String("url"), ctx.String("token")),
		Private:    ctx.Bool("private"),
		NewOwnerID: ctx.Int64("owner"),
	}
	if m.NewOwnerID == 0 {
		usr, err := m.Client.GetMyUserInfo()
		if err != nil {
			return fmt.Errorf("cannot fetch user info about current user: %v", err)
		}
		m.NewOwnerID = usr.ID
	}
	c := context.Background()

	var gc *github.Client
	if ctx.IsSet("gh-token") {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: ctx.String("gh-token")},
		)
		tc := oauth2.NewClient(c, ts)
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
		repos, resp, err := gc.Repositories.List(c, ctx.String("gh-user"), opt)
		if err != nil {
			return err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	errs := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(len(allRepos))
	for _, repo := range allRepos {
		go func(r *github.Repository) {
			defer wg.Done()
			errs <- migrate(c, gc, m, r.Owner.GetLogin(), r.GetName(), ctx.Bool("only-repo"))
		}(repo)
	}

	go func() {
		for i := range errs {
			if i != nil {
				fmt.Printf("error: %v", i)
			}
		}
	}()

	wg.Wait()
	return nil
}
