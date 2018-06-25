package migrations

import (
	"fmt"

	"code.gitea.io/sdk/gitea"
	"github.com/google/go-github/github"
)

// Issue migrates a GitHub Issue to a Gitea Issue
func (m *Migratory) Issue(gi *github.Issue) (*gitea.Issue, error) {
	if m.migratedMilestones == nil {
		m.migratedMilestones = make(map[int64]int64)
	}
	if m.migratedLabels == nil {
		m.migratedLabels = make(map[int64]int64)
	}

	// Migrate milestone if it is not already migrated
	milestone := int64(0)
	if gi.Milestone != nil {
		// Lookup if milestone is already migrated
		if migratedMilestone, ok := m.migratedMilestones[*gi.Milestone.ID]; ok {
			milestone = migratedMilestone
		} else if ms, err := m.Milestone(gi.Milestone); err != nil {
			return nil, err
		} else {
			milestone = ms.ID
		}
	}
	// Migrate labels
	labels, err := m.labels(gi.Labels)
	if err != nil {
		return nil, err
	}

	return m.Client.CreateIssue(m.repository.Owner.UserName, m.repository.Name,
		gitea.CreateIssueOption{
			Title:     gi.GetTitle(),
			Body:      fmt.Sprintf("Author: @%s Posted at: %s\n\n\n%s", *gi.User.Login, gi.GetCreatedAt().Format("02.01.2006 15:04"), gi.GetBody()),
			Closed:    *gi.State == "closed",
			Milestone: milestone,
			Labels:    labels,
		})
}

func (m *Migratory) labels(gls []github.Label) (results []int64, err error) {
	for _, gl := range gls {
		if migratedLabel, ok := m.migratedLabels[*gl.ID]; ok {
			results = append(results, migratedLabel)
		} else {
			var newLabel *gitea.Label
			if newLabel, err = m.Label(&gl); err != nil {
				return nil, err
			}
			m.migratedLabels[*gl.ID] = newLabel.ID
			results = append(results, newLabel.ID)
		}
	}
	return
}

// Label migrates a GitHub Label to a Gitea Label without caching its id
func (m *Migratory) Label(gl *github.Label) (*gitea.Label, error) {
	return m.Client.CreateLabel(m.repository.Owner.UserName, m.repository.Name,
		gitea.CreateLabelOption{
			Name:  gl.GetName(),
			Color: fmt.Sprintf("#%s", gl.GetColor()),
		})
}

// Milestone migrates a GitHub Milesteon to a Gitea Milestone and caches its id
func (m *Migratory) Milestone(gm *github.Milestone) (*gitea.Milestone, error) {
	ms, err := m.Client.CreateMilestone(m.repository.Owner.UserName, m.repository.Name,
		gitea.CreateMilestoneOption{
			Title:       gm.GetTitle(),
			Description: gm.GetDescription(),
			Deadline:    gm.DueOn,
		})
	if err != nil {
		return nil, err
	}
	m.migratedMilestones[*gm.ID] = ms.ID
	if gm.State != nil && *gm.State != "open" {
		return m.Client.EditMilestone(m.repository.Owner.UserName, m.repository.Name,
			ms.ID, gitea.EditMilestoneOption{
				State: githubStateToGiteaState(gm.State),
			})
	}
	return ms, err
}

func githubStateToGiteaState(ghstate *string) *string {
	if ghstate == nil {
		return ghstate
	}
	switch *ghstate {
	case "open":
		fallthrough
	case "closed":
		return ghstate
	case "all":
		open := "open"
		return &open
	}
	return nil
}

// IssueComment migrates a GitHub IssueComment to a Gitea Comment
func (m *Migratory) IssueComment(issue *gitea.Issue, gic *github.IssueComment) (*gitea.Comment, error) {
	return m.Client.CreateIssueComment(m.repository.Owner.UserName,
		m.repository.Name,
		issue.Index,
		gitea.CreateIssueCommentOption{
			Body: fmt.Sprintf("Author: @%s Posted at: %s\n\n\n%s", *gic.User.Login, gic.GetCreatedAt().Format("02.01.2006 15:04"), gic.GetBody()),
		})
}
