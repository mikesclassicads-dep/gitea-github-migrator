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

// Deprecated: Please use Job or FetchMigratory instead
func migrate(c context.Context, gc *github.Client, m *migrations.Migratory, username, repo string, onlyRepo bool) error {
	fmt.Printf("Fetching repository %s/%s...\n", username, repo)
	gr, _, err := gc.Repositories.Get(c, username, repo)
	if err != nil {
		return fmt.Errorf("error while fetching repo[%s/%s]: %v", username, repo, err)
	}

	fmt.Printf("Migrating repository %s/%s...\n", username, repo)
	var mr *gitea.Repository
	if mr, err = m.Repository(gr); err != nil {
		return fmt.Errorf("error while migrating repo[%s/%s]: %v", username, repo, err)
	}
	fmt.Printf("Repository migrated to %s/%s\n", mr.Owner.UserName, mr.Name)

	if onlyRepo {
		return nil
	}

	fmt.Println("Fetching issues...")
	opt := &github.IssueListByRepoOptions{
		Sort:      "created",
		Direction: "asc",
		State:     "all",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	var allIssues = make([]*github.Issue, 0)
	for {
		issues, resp, err := gc.Issues.ListByRepo(c, username, repo, opt)
		if err != nil {
			return fmt.Errorf("error while listing repos: %v", err)
		}
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	fmt.Println("Migrating issues...")
	for _, gi := range allIssues {
		fmt.Printf("Migrating #%d...\n", *gi.Number)
		issue, err := m.Issue(gi)
		if err != nil {
			return fmt.Errorf("migrating issue[id: %d]: %v", *gi.ID, err)
		}
		comments, _, err := gc.Issues.ListComments(c, username, repo, gi.GetNumber(), nil)
		if err != nil {
			return fmt.Errorf("fetching issue[id: %d] comments: %v", *gi.ID, err)
		}
		for _, gc := range comments {
			fmt.Printf("-> %d...", gc.ID)
			if _, err := m.IssueComment(issue, gc); err != nil {
				return fmt.Errorf("migrating issue comment [issue: %d, comment: %d]: %v", *gi.ID, gc.ID, err)
			}
			fmt.Print("Done!\n")
		}
		fmt.Printf("Migrated #%d...\n", *gi.Number)

	}
	return nil
}
