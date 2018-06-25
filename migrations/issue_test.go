package migrations

import (
	"testing"

	"code.gitea.io/sdk/gitea"
	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
)

func Test_githubStateToGiteaState(t *testing.T) {
	open := "open"
	all := "all"
	closed := "closed"
	tests := map[*string]string{&all: "open", &open: "open", &closed: "closed"}
	for input, exceptedResult := range tests {
		actualResult := githubStateToGiteaState(input)
		assert.NotNil(t, actualResult)
		assert.NotEmpty(t, *actualResult)
		assert.Equal(t, exceptedResult, *actualResult)
	}
	nilInput := "teoafweogwoe"
	assert.Nil(t, githubStateToGiteaState(&nilInput))
}

func TestMigratory_Label(t *testing.T) {
	res, err := DemoMigratory.Label(&github.Label{
		Name:  github.String("testlabel"),
		Color: github.String("123456"),
	})
	assert.NoError(t, err)
	assert.Equal(t, "123456", res.Color)
	assert.Equal(t, "testlabel", res.Name)
}

func TestMigratory_Milestone(t *testing.T) {
	res, err := DemoMigratory.Milestone(&github.Milestone{
		ID:          github.Int64(1),
		State:       github.String("open"),
		Description: github.String("test milestone"),
		Title:       github.String("TEST"),
		DueOn:       &demoTime,
	})
	assert.NoError(t, err)
	assert.Equal(t, "TEST", res.Title)
	assert.Equal(t, "test milestone", res.Description)
	assert.Equal(t, demoTime.Unix(), res.Deadline.Unix())
	assert.Equal(t, gitea.StateOpen, res.State)
}
