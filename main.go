package main

import (
	"os"

	"github.com/urfave/cli"

	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/cmd"
)

func main() {
	app := cli.NewApp()
	app.Name = "gg-migrator"
	app.Usage = "GitHub to Gitea migrator for repositories"
	app.Description = `Migrate your GitHub repositories including issues to Gitea`
	app.Commands = cli.Commands{
		cmd.CmdMigrate,
		cmd.CmdMigrateAll,
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
