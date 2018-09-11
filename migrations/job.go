package migrations

import (
	"bytes"
	"fmt"
	"strings"

	"code.gitea.io/sdk/gitea"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

// Job manages all migrations of a "migartion job"
type Job struct {
	Repositories []string
	Options      *Options
	Client       *gitea.Client
	GHClient     *github.Client
	UseStdErr    bool

	migratories map[string]*FetchMigratory
}

// JobReport represents the current status of a Job
type JobReport struct {
	Pending  []string                    `json:"pending"`
	Running  map[string]*MigratoryStatus `json:"running"`
	Finished map[string]*MigratoryStatus `json:"finished"`
	Failed   map[string]string           `json:"failed"`
}

// NewJob returns an instance of initialized instance of Job
func NewJob(options *Options, client *gitea.Client, githubClient *github.Client, repos ...string) *Job {
	return &Job{Repositories: repos, Options: options, Client: client, GHClient: githubClient}
}

// StatusReport generates a JobReport indicating which state the job is
func (job *Job) StatusReport() *JobReport {
	report := &JobReport{
		Pending:  make([]string, 0),
		Finished: make(map[string]*MigratoryStatus),
		Running:  make(map[string]*MigratoryStatus),
		Failed:   make(map[string]string),
	}
	for _, repo := range job.Repositories {
		if migratory, ok := job.migratories[repo]; ok {
			migratory.Status.Log = migratory.LogOutput.String()
			switch migratory.Status.Stage {
			case Finished:
				report.Finished[repo] = migratory.Status
			case Importing:
			case Migrating:
				report.Running[repo] = migratory.Status
			case Failed:
				report.Failed[repo] = migratory.Status.FatalError.Error()
			default:
				report.Pending = append(report.Pending, repo)
				fmt.Printf("unknown status %d\n", migratory.Status.Stage)
			}
		} else {
			report.Pending = append(report.Pending, repo)
		}
	}
	return report
}

// StartMigration migrates all repos from Repositories
func (job *Job) StartMigration() chan error {
	errs := make(chan error, len(job.Repositories))
	var pendingRepos = len(job.Repositories)
	autoclose := func() {
		pendingRepos--
		if pendingRepos <= 0 {
			close(errs)
		}
	}
	job.migratories = make(map[string]*FetchMigratory, pendingRepos)
	for _, repo := range job.Repositories {
		mig, err := job.initFetchMigratory(repo)
		job.migratories[repo] = mig
		if err != nil {
			mig.Status = &MigratoryStatus{
				Stage:      Failed,
				FatalError: err,
			}
			errs <- err
			autoclose()
			continue
		}
		go func() {
			err := mig.MigrateFromGitHub()
			errs <- err
			autoclose()
		}()
	}
	return errs
}

func (job *Job) initFetchMigratory(repo string) (*FetchMigratory, error) {
	res := strings.Split(repo, "/")
	if len(res) != 2 {
		return nil, fmt.Errorf("invalid repo name: %s", repo)
	}
	fm := &FetchMigratory{
		Migratory: Migratory{
			Client:  job.Client,
			Options: *job.Options,
		},
		RepoName:  res[1],
		RepoOwner: res[0],
		GHClient:  job.GHClient,
		Logger:    logrus.New(),
		LogOutput: new(bytes.Buffer),
	}
	if !job.UseStdErr {
		fm.Logger.Formatter = &logrus.TextFormatter{
			DisableColors:          true,
			DisableLevelTruncation: true,
			DisableTimestamp:       true,
		}
		fm.Logger.SetOutput(fm.LogOutput)
	} else {
		fm.LogOutput = nil
	}
	return fm, nil
}

// Finished indicates if the job is finished or not
func (job *Job) Finished() bool {
	return (len(job.StatusReport().Failed) + len(job.StatusReport().Finished)) >= len(job.Repositories)
}
