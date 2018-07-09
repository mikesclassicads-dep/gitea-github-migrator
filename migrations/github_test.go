package migrations

import (
	"testing"
)

func BenchmarkGetIssueIndexFromHTMLURLAlt(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		getIssueIndexFromHTMLURLAlt("https://github.com/octocat/Hello-World/issues/1347#issuecomment-1")
	}
}

func BenchmarkGetIssueIndexFromHTMLURL(b *testing.B) {
	for i := 0; i <= b.N; i++ {
		getIssueIndexFromHTMLURL("https://github.com/octocat/Hello-World/issues/1347#issuecomment-1")
	}
}
