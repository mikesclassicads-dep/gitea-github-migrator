package cmd

import "github.com/urfave/cli"

var defaultMigrateFlags = []cli.Flag{
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
		Name:   "gh-token",
		Usage:  "GitHub Token (optional)",
		EnvVar: "GITHUB_TOKEN",
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
	cli.BoolFlag{
		Name:   "only-repo",
		Usage:  "skip issues etc. and only migrate repo",
		EnvVar: "ONLY_REPO",
	},
}
