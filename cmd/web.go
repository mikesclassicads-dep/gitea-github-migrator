// +build web

package cmd

import (
	"fmt"
	"net/http"

	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/config"
	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/web"
	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// CmdWeb stars the web interface
var CmdWeb = cli.Command{
	Name:   "web",
	Usage:  "Starts the web interface",
	Action: runWeb,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "c,config",
			Usage:  "config file",
			Value:  "config.yml",
			EnvVar: "MIGRATOR_CONFIG",
		},
	},
}

func runWeb(ctx *cli.Context) error {
	if err := configor.New(&configor.Config{ErrorOnUnmatchedKeys: true}).Load(&config.Config, ctx.String("config")); err != nil {
		return err
	}
	r := web.InitRoutes()

	hostname := config.Config.Web.Host
	if len(hostname) == 0 {
		hostname = "0.0.0.0"
	}
	port := config.Config.Web.Port
	if port == 0 {
		port = 4000
	}
	logrus.Infof("Server is running at http://%s:%d", hostname, port)
	logrus.SetLevel(logrus.PanicLevel)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", hostname, port), r)
}
