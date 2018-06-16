package web

import (
	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/web/auth"
	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/web/context"
	"github.com/go-macaron/session"
	"github.com/gobuffalo/packr"
	"gopkg.in/macaron.v1"
)

// InitRoutes initiates the gin routes and loads values from config
func InitRoutes() *macaron.Macaron {
	m := macaron.Classic()
	auth.InitGitHubOAuthConfig()
	tmplBox := packr.NewBox("templates")
	publicBox := packr.NewBox("public")
	m.Use(macaron.Recovery())
	m.Use(session.Sessioner())
	m.Use(macaron.Renderer(macaron.RenderOptions{
		TemplateFileSystem: &BundledFS{tmplBox},
	}))
	m.Use(macaron.Statics(macaron.StaticOptions{
		Prefix:     "static",
		FileSystem: publicBox,
	}, ""))
	m.Use(context.Contexter())
	m.Get("/", func(ctx *context.Context) {
		if ctx.User != nil {
			ctx.HTML(200, "migrate")
			return
		}
		ctx.HTML(200, "login_github") // 200 is the response code.
	})
	m.Group("/github", func() {
		m.Get("/", auth.RedirectToGitHub)
		m.Get("/callback", auth.CallbackFromGitHub)
	})
	return m
}
