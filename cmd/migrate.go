package cmd

import (
	"context"
	"fmt"

	"code.gitea.io/sdk/gitea"
	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/migrations"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
)

// CmdMigrate migrates a given repository to gitea
var CmdMigrate = cli.Command{
	Name:   "migrate",
	Usage:  "migrates a github to a gitea repository",
	Action: runMigrate,
	Flags: append(defaultMigrateFlags,
		cli.StringFlag{
			Name:   "gh-repo",
			Usage:  "GitHub Repository",
			Value:  "username/reponame",
			EnvVar: "GH_REPOSITORY",
		},
	),
}

func runMigrate(ctx *cli.Context) error {
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
	logrus.SetLevel(logrus.InfoLevel)
	job := migrations.NewJob(&migrations.Options{
		Private:    ctx.Bool("private"),
		NewOwnerID: ctx.Int("owner"),

		Comments:     !onlyRepos,
		Issues:       !onlyRepos,
		Labels:       !onlyRepos,
		Milestones:   !onlyRepos,
		PullRequests: !onlyRepos,
		Strategy:     migrations.Classic,
	}, gitea.NewClient(ctx.String("url"), ctx.String("token")), gc, ctx.String("gh-repo"))
	if job.Options.NewOwnerID == 0 {
		usr, err := job.Client.GetMyUserInfo()
		if err != nil {
			return fmt.Errorf("cannot fetch user info about current user: %v", err)
		}
		job.Options.NewOwnerID = int(usr.ID)
	}
	errs := job.StartMigration()
	for err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
