package cmd

import (
	"fmt"
	"net/http"

	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/config"
	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/web"
	"github.com/jinzhu/configor"
	"github.com/urfave/cli"
)

// CmdWeb stars the web interface
var CmdWeb = cli.Command{
	Name:   "web",
	Usage:  "Starts the web interface",
	Action: runWeb,
}

func runWeb(_ *cli.Context) error {
	if err := configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(&config.Config, "config.yml"); err != nil {
		return err
	}
	r := web.InitRoutes()

	fmt.Println("Server is running...")
	// TODO add port / host to config
	return http.ListenAndServe("0.0.0.0:4000", r)
}
