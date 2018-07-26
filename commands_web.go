// +build web

package main

import (
	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/cmd"
	"github.com/urfave/cli"
)

var cmds = cli.Commands{
	cmd.CmdMigrate,
	cmd.CmdMigrateAll,
	cmd.CmdWeb,
}
