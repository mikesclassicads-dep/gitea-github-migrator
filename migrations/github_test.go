package migrations

import (
	"testing"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
)

func BenchmarkGetIssueIndexFromHTMLURLAlt(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		getIssueIndexFromHTMLURLAlt("https://github.com/octocat/Hello-World/issues/1347#issuecomment-1")
	}
}

func TestGetIssueIndexFromHTMLURLAlt(t *testing.T) {
	res, err := getIssueIndexFromHTMLURLAlt("https://github.com/octocat/Hello-World/issues/1347#issuecomment-1")
	assert.NoError(t, err)
	assert.Equal(t, 1347, res)
	res, err = getIssueIndexFromHTMLURLAlt("https://github.com/oment-1")
	assert.Error(t, err)
}

func BenchmarkGetIssueIndexFromHTMLURL(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		getIssueIndexFromHTMLURL("https://github.com/octocat/Hello-World/issues/1347#issuecomment-1")
	}
}

func TestGetIssueIndexFromHTMLURL(t *testing.T) {
	res, err := getIssueIndexFromHTMLURL("https://github.com/octocat/Hello-World/issues/1347#issuecomment-1")
	assert.NoError(t, err)
	assert.Equal(t, 1347, res)
	res, err = getIssueIndexFromHTMLURL("https://github.com/oment-1")
	assert.Error(t, err)
}

var testFMig = &FetchMigratory{
	Migratory: *DemoMigratory,
	GHClient:  github.NewClient(nil),
	RepoOwner: "JonasFranzDEV",
	RepoName:  "test",
}

func TestFetchMigratory_FetchIssues(t *testing.T) {
	issues, err := testFMig.FetchIssues()
	assert.NoError(t, err)
	assert.True(t, len(issues) > 0, "at least one issue found")
}

func TestFetchMigratory_FetchComments(t *testing.T) {
	comments, err := testFMig.FetchIssues()
	assert.NoError(t, err)
	assert.True(t, len(comments) > 0, "at least one comment found")
}
