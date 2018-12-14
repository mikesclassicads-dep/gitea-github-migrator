package migrations

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"code.gitea.io/sdk/gitea"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

// FetchMigratory adds GitHub fetching functions to migratory
type FetchMigratory struct {
	Migratory
	GHClient  *github.Client
	RepoOwner string
	RepoName  string
	Logger    *logrus.Logger
	LogOutput *bytes.Buffer
}

func (fm *FetchMigratory) ctx() context.Context {
	return context.Background()
}

// MigrateFromGitHub migrates RepoOwner/RepoName from GitHub to Gitea
func (fm *FetchMigratory) MigrateFromGitHub() error {
	fm.Status = &MigratoryStatus{
		Stage: Importing,
	}

	fm.Logger.WithFields(logrus.Fields{
		"repo": fmt.Sprintf("%s/%s", fm.RepoOwner, fm.RepoName),
	}).Info("migrating git repository")
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
	fm.Logger.WithFields(logrus.Fields{
		"repo": fmt.Sprintf("%s/%s", fm.RepoOwner, fm.RepoName),
	}).Info("git repository migrated")
	if fm.Options.Issues || fm.Options.PullRequests {
		if err := fm.MigrateIssuesFromGitHub(); err != nil {
			return err
		}
	}
	if fm.Status.FatalError != nil {
		fm.Status.Stage = Failed
		fm.Logger.WithFields(logrus.Fields{
			"repo": fmt.Sprintf("%s/%s", fm.RepoOwner, fm.RepoName),
		}).Errorf("migration failed: %v", fm.Status.FatalError)
		return nil
	}
	fm.Status.Stage = Finished
	fm.Logger.WithFields(logrus.Fields{
		"repo": fmt.Sprintf("%s/%s", fm.RepoOwner, fm.RepoName),
	}).Info("migration successful")
	return nil
}

// MigrateIssuesFromGitHub migrates all issues from GitHub to Gitea
func (fm *FetchMigratory) MigrateIssuesFromGitHub() error {
	var commentsChan chan *[]*github.IssueComment
	if fm.Options.Comments {
		commentsChan = fm.fetchCommentsAsync()
	}
	issues, err := fm.FetchIssues()
	if err != nil {
		fm.Status.Stage = Failed
		fm.Status.FatalError = err
		fm.Logger.WithFields(logrus.Fields{
			"repo": fmt.Sprintf("%s/%s", fm.RepoOwner, fm.RepoName),
		}).Errorf("migration failed: %v", fm.Status.FatalError)
		return err
	}
	fm.Status.Stage = Migrating
	fm.Status.Issues = int64(len(issues))
	migratedIssues := make(map[int]*gitea.Issue)
	for _, issue := range issues {
		if (!issue.IsPullRequest() || fm.Options.PullRequests) &&
			(issue.IsPullRequest() || fm.Options.Issues) {
			migratedIssues[issue.GetNumber()], err = fm.Issue(issue)
			if err != nil {
				fm.Status.IssuesError++
				fm.Logger.WithFields(logrus.Fields{
					"repo":  fmt.Sprintf("%s/%s", fm.RepoOwner, fm.RepoName),
					"issue": issue.GetNumber(),
				}).Warnf("error while migrating: %v", err)
				continue
			}
			fm.Status.IssuesMigrated++
			fm.Logger.WithFields(logrus.Fields{
				"repo":  fmt.Sprintf("%s/%s", fm.RepoOwner, fm.RepoName),
				"issue": issue.GetNumber(),
			}).Info("issue migrated")
		} else {
			fm.Status.Issues--
		}
	}
	if fm.Options.Comments {
		var comments []*github.IssueComment
		var cmts *[]*github.IssueComment
		if cmts = <-commentsChan; cmts == nil {
			fm.Status.Stage = Failed
			err := fmt.Errorf("error while fetching issue comments")
			fm.Logger.WithFields(logrus.Fields{
				"repo": fmt.Sprintf("%s/%s", fm.RepoOwner, fm.RepoName),
			}).Errorf("migration failed: %v", fm.Status.FatalError)
			return err
		}
		return fm.migrateCommentsFromGitHub(comments, migratedIssues)
	}
	return nil
}

func (fm *FetchMigratory) migrateCommentsFromGitHub(comments []*github.IssueComment, migratedIssues map[int]*gitea.Issue) error {
	fm.Status.Comments = int64(len(comments))
	commentsByIssue := make(map[*gitea.Issue][]*github.IssueComment, len(migratedIssues))
	for _, comment := range comments {
		issueIndex, err := getIssueIndexFromHTMLURL(comment.GetHTMLURL())
		if err != nil {
			fm.Status.CommentsError++
			fm.Logger.WithFields(logrus.Fields{
				"repo":    fmt.Sprintf("%s/%s", fm.RepoOwner, fm.RepoName),
				"issue":   issueIndex,
				"comment": comment.GetID(),
			}).Warnf("error while migrating comment: %v", err)
			continue
		}
		if issue, ok := migratedIssues[issueIndex]; ok && issue != nil {
			if list, ok := commentsByIssue[issue]; !ok && list != nil {
				commentsByIssue[issue] = []*github.IssueComment{comment}
			} else {
				commentsByIssue[issue] = append(list, comment)
			}
		} else {
			fm.Status.CommentsError++
			continue
		}
	}
	wg := sync.WaitGroup{}
	for issue, comms := range commentsByIssue {
		wg.Add(1)
		go func(i *gitea.Issue, cs []*github.IssueComment) {
			for _, comm := range cs {
				if _, err := fm.IssueComment(i, comm); err != nil {
					fm.Status.CommentsError++
					fm.Logger.WithFields(logrus.Fields{
						"repo":    fmt.Sprintf("%s/%s", fm.RepoOwner, fm.RepoName),
						"comment": comm.GetID(),
					}).Warnf("error while migrating comment: %v", err)
					continue
				}
				fm.Status.CommentsMigrated++
				fm.Logger.WithFields(logrus.Fields{
					"repo":    fmt.Sprintf("%s/%s", fm.RepoOwner, fm.RepoName),
					"comment": comm.GetID(),
				}).Info("comment migrated")
			}
			wg.Done()
		}(issue, comms)
	}
	wg.Wait()
	return nil
}

var issueIndexRegex = regexp.MustCompile(`/(issues|pull)/([0-9]+)#`)

func getIssueIndexFromHTMLURL(htmlURL string) (int, error) {
	// Alt is 4 times faster but more error prune
	if res, err := getIssueIndexFromHTMLURLAlt(htmlURL); err == nil {
		return res, nil
	}
	matches := issueIndexRegex.FindStringSubmatch(htmlURL)
	if len(matches) < 3 {
		return 0, fmt.Errorf("cannot parse issue id from HTML URL: %s", htmlURL)
	}
	return strconv.Atoi(matches[2])
}
func getIssueIndexFromHTMLURLAlt(htmlURL string) (int, error) {
	res := strings.Split(htmlURL, "/issues/")
	if len(res) != 2 {
		res = strings.Split(htmlURL, "/pull/")
	}
	if len(res) != 2 {
		return 0, fmt.Errorf("invalid HTMLURL: %s", htmlURL)
	}
	number := res[1]
	number = strings.Split(number, "#")[0]
	return strconv.Atoi(number)
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

// FetchComments fetches all comments from GitHub
func (fm *FetchMigratory) FetchComments() ([]*github.IssueComment, error) {
	var allComments = make([]*github.IssueComment, 0)
	opt := &github.IssueListCommentsOptions{
		Sort:      "created",
		Direction: "asc",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	for {
		comments, resp, err := fm.GHClient.Issues.ListComments(fm.ctx(), fm.RepoOwner, fm.RepoName, 0, opt)
		if err != nil {
			return nil, fmt.Errorf("error while listing repos: %v", err)
		}
		allComments = append(allComments, comments...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return allComments, nil
}

func (fm *FetchMigratory) fetchCommentsAsync() chan *[]*github.IssueComment {
	ret := make(chan *[]*github.IssueComment, 1)
	go func(f *FetchMigratory) {
		comments, err := f.FetchComments()
		if err != nil {
			f.Status.FatalError = err
			ret <- nil
			fm.Logger.WithFields(logrus.Fields{
				"repo": fmt.Sprintf("%s/%s", fm.RepoOwner, fm.RepoName),
			}).Errorf("fetching comments failed: %v", fm.Status.FatalError)
			return
		}
		f.Status.Comments = int64(len(comments))
		ret <- &comments
	}(fm)
	return ret
}
