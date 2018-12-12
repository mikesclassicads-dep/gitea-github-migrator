package migrations

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJob_StatusReport(t *testing.T) {
	jobWithStatus := func(status *MigratoryStatus) *Job {
		return &Job{
			migratories: map[string]*FetchMigratory{
				"test/test": {
					Migratory: Migratory{
						Status: status,
					},
				},
			},
			Repositories: []string{
				"test/test",
			},
		}
	}
	// Pending
	pendingJob := &Job{
		Repositories: []string{
			"test/test",
		},
	}
	report := pendingJob.StatusReport()
	assert.Len(t, report.Pending, 1)
	assert.Equal(t, report.Pending[0], "test/test")
	assert.Len(t, report.Failed, 0)
	assert.Len(t, report.Running, 0)
	assert.Len(t, report.Finished, 0)

	// Finished
	report = jobWithStatus(&MigratoryStatus{
		Stage: Finished,
	}).StatusReport()
	assert.Len(t, report.Pending, 0)
	assert.Len(t, report.Failed, 0)
	assert.Len(t, report.Running, 0)
	assert.Len(t, report.Finished, 1)
	assert.Equal(t, Finished, report.Finished["test/test"].Stage)

	// Failed
	report = jobWithStatus(&MigratoryStatus{
		Stage:      Failed,
		FatalError: fmt.Errorf("test"),
	}).StatusReport()
	assert.Len(t, report.Failed, 1)
	assert.Equal(t, "test", report.Failed["test/test"])
	assert.Len(t, report.Pending, 0)
	assert.Len(t, report.Running, 0)
	assert.Len(t, report.Finished, 0)

	// Running
	report = jobWithStatus(&MigratoryStatus{
		Stage: Migrating,
	}).StatusReport()
	assert.Len(t, report.Running, 1)
	assert.Equal(t, Migrating, report.Running["test/test"].Stage)
	assert.Len(t, report.Pending, 0)
	assert.Len(t, report.Failed, 0)
	assert.Len(t, report.Finished, 0)
}
