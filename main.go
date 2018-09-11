package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	version = "0.0.0"
	build   = "0"
)

func main() {
	app := cli.NewApp()
	app.Name = "gitea-github-migrator"
	app.Version = fmt.Sprintf("%s+%s", version, build)
	app.Usage = "GitHub to Gitea migrator for repositories"
	app.Description = `Migrate your GitHub repositories including issues to Gitea`
	app.Commands = cmds
	if err := app.Run(os.Args); err != nil {
		logrus.Panic(err)
	}
}
