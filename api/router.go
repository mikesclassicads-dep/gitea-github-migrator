// Package classification Gitea GitHub Migrator API
//
// the purpose of this application is to provide access to the migrator backend
//
//     Schemes: http, https
//     Host: localhost
//     BasePath: /api
//     Version: 0.0.9
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Jonas Franz<info@jonasfranz.de> https://jonasfranz.de
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package api

import (
	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/api/auth"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr"
)

func InitRoutes() *gin.Engine {
	r := gin.Default()
	auth.InitOAuthConfig()

	box := packr.NewBox("../frontend/dist/")
	r.Use(static.Serve("/", &BundledFS{box}))
	api := r.Group("/api")
	{
		au := api.Group("/auth")
		{
			au.GET("/github", auth.RedirectToGitHub)
		}
	}

	return r
}
