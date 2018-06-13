package cmd

import (
	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/api"
	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/config"
	"github.com/jinzhu/configor"
	"github.com/urfave/cli"
)

var CmdWeb = cli.Command{
	Name:   "web",
	Usage:  "Starts the web interface",
	Action: runWeb,
}

func runWeb(_ *cli.Context) error {
	if err := configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(&config.Config, "config.yml"); err != nil {
		return err
	}
	r := api.InitRoutes()
	return r.Run(":8081")
}
