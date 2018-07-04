package migrations

import (
	"context"
	"fmt"
	"sync"

	"code.gitea.io/sdk/gitea"
	"github.com/google/go-github/github"
)

// FetchMigratory adds GitHub fetching functions to migratory
type FetchMigratory struct {
	Migratory
	GHClient  *github.Client
	RepoOwner string
	RepoName  string
}

func (fm *FetchMigratory) ctx() context.Context {
	return context.Background()
}

// MigrateFromGitHub migrates RepoOwner/RepoName from GitHub to Gitea
func (fm *FetchMigratory) MigrateFromGitHub() error {
	fm.Status = &MigratoryStatus{
		Stage: Importing,
	}
	ghRepo, _, err := fm.GHClient.Repositories.Get(fm.ctx(), fm.RepoOwner, fm.RepoName)
	if err != nil {
		fm.Status.Stage = Failed
		fm.Status.FatalError = err
		return fmt.Errorf("GHClient Repostiories Get: %v", err)
	}
	fm.repository, err = fm.Repository(ghRepo)
	if err != nil {
		fm.Status.Stage = Failed
		fm.Status.FatalError = err
		return fmt.Errorf("Repository migration: %v", err)
	}
	var wg sync.WaitGroup
	if fm.Options.Issues || fm.Options.PullRequests {
		issues, err := fm.FetchIssues()
		if err != nil {
			fm.Status.Stage = Failed
			fm.Status.FatalError = err
			return err
		}
		fm.Status.Stage = Migrating
		for _, issue := range issues {
			if (!issue.IsPullRequest() || fm.Options.PullRequests) &&
				(issue.IsPullRequest() || fm.Options.Issues) {
				fm.Status.Issues++
				giteaIssue, err := fm.Issue(issue)
				if err != nil {
					fm.Status.IssuesError++
					// TODO log errors
					continue
				}
				wg.Add(1)
				go func() {
					fm.FetchAndMigrateComments(issue, giteaIssue)
					wg.Done()
				}()
				fm.Status.IssuesMigrated++
			}
		}
	}
	wg.Wait()
	if fm.Status.FatalError != nil {
		fm.Status.Stage = Failed
		return nil
	}
	fm.Status.Stage = Finished
	return nil
}

// FetchAndMigrateComments loads all comments from GitHub and migrates them to Gitea
func (fm *FetchMigratory) FetchAndMigrateComments(issue *github.Issue, giteaIssue *gitea.Issue) {
	comments, _, err := fm.GHClient.Issues.ListComments(fm.ctx(), fm.RepoOwner, fm.RepoName, issue.GetNumber(), nil)
	if err != nil {
		// TODO log errors
		return
	}
	fm.Status.Comments += int64(len(comments))
	for _, gc := range comments {
		if _, err := fm.IssueComment(giteaIssue, gc); err != nil {
			fm.Status.CommentsError++
			// TODO log errors
			return
		}
		fm.Status.CommentsMigrated++
	}
}

// FetchIssues fetches all issues from GitHub
func (fm *FetchMigratory) FetchIssues() ([]*github.Issue, error) {
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
		issues, resp, err := fm.GHClient.Issues.ListByRepo(fm.ctx(), fm.RepoOwner, fm.RepoName, opt)
		if err != nil {
			return nil, fmt.Errorf("error while listing repos: %v", err)
		}
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allIssues, nil
}
