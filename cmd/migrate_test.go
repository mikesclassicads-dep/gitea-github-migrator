package cmd

import (
	"context"
	"testing"

	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/migrations"
	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
)

// Test_migrate is an integration tests for migrate command
// using repo JonasFranzDEV/test
func Test_migrate(t *testing.T) {
	assert.NoError(t, migrate(
		context.Background(),
		github.NewClient(nil),
		migrations.DemoMigratory,
		"JonasFranzDEV",
		"test",
		false,
	))
}
