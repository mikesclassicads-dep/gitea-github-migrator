package cmd

import (
	"code.gitea.io/sdk/gitea"
	"context"
	"fmt"
	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/migrations"
	"github.com/google/go-github/github"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
	"strings"
)

var CmdMigrate = cli.Command{
	Name:   "migrate",
	Usage:  "migrates a github to a gitea repository",
	Action: runMigrate,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "gh-repo",
			Usage:  "GitHub Repository",
			Value:  "username/reponame",
			EnvVar: "GH_REPOSITORY",
		},
		cli.IntFlag{
			Name:   "owner",
			Usage:  "Owner ID",
			EnvVar: "OWNER_ID",
			Value:  0,
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "Gitea Token",
			EnvVar: "GITEA_TOKEN",
		},
		cli.StringFlag{
			Name:   "url",
			Usage:  "Gitea URL",
			EnvVar: "GITEA_URL",
		},
		cli.BoolFlag{
			Name:   "private",
			Usage:  "should new repository be private",
			EnvVar: "GITEA_PRIVATE",
		},
	},
}

func runMigrate(ctx *cli.Context) error {
	m := migrations.Migratory{
		Client:     gitea.NewClient(ctx.String("url"), ctx.String("token")),
		Private:    ctx.Bool("private"),
		NewOwnerID: ctx.Int("owner"),
	}
	c := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "246947e59029a71cb568e290cc6b28adeda514df"},
	)
	tc := oauth2.NewClient(c, ts)
	gc := github.NewClient(tc)

	username := strings.Split(ctx.String("gh-repo"), "/")[0]
	repo := strings.Split(ctx.String("gh-repo"), "/")[1]

	//p := mpb.New()

	fmt.Printf("Fetching repository %s/%s...\n", username, repo)
	gr, _, err := gc.Repositories.Get(c, username, repo)
	if err != nil {
		return err
	}
	fmt.Printf("Migrating repository %s/%s...\n", username, repo)
	if mr, err := m.Repository(gr); err != nil {
		return err
	} else {
		fmt.Printf("Repository migrated to %s/%s\n", mr.Owner.UserName, mr.Name)
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
			return err
		}
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	fmt.Println("Migrating issues...")
	//bar := p.AddBar(int64(len(issues)))
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
		//bar.Increment()

	}
	return nil
}
