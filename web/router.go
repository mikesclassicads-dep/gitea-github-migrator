package web

import (
	"encoding/gob"

	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/web/auth"
	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/web/context"
	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/web/migration"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/session"
	"github.com/gobuffalo/packr"
	"gopkg.in/macaron.v1"
)

// InitRoutes initiates the gin routes and loads values from config
func InitRoutes() *macaron.Macaron {
	gob.Register(&context.User{})
	m := macaron.Classic()
	auth.InitGitHubOAuthConfig()
	tmplBox := packr.NewBox("templates")
	publicBox := packr.NewBox("public")
	m.Use(macaron.Recovery())
	m.Use(session.Sessioner(session.Options{
		Provider:       "file",
		ProviderConfig: "data/sessions",
	}))
	m.Use(macaron.Renderer(macaron.RenderOptions{
		TemplateFileSystem: &BundledFS{tmplBox},
	}))
	m.Use(macaron.Statics(macaron.StaticOptions{
		Prefix:     "static",
		FileSystem: publicBox,
	}, ""))
	m.Use(context.Contexter())

	// BEGIN: Router
	m.Get("/", func(ctx *context.Context) {
		if ctx.User != nil {
			if ctx.GiteaUser == nil {
				ctx.HTML(200, "login_gitea")
				return
			}
			ctx.HTML(200, "dashboard")
			return
		}
		ctx.HTML(200, "login_github") // 200 is the response code.
	})
	m.Get("/logout", func(c *macaron.Context, sess session.Store) {
		sess.Destory(c)
		c.Redirect("/")
	})
	m.Group("/github", func() {
		m.Get("/", auth.RedirectToGitHub)
		m.Get("/callback", auth.CallbackFromGitHub)
	})
	m.Group("/gitea", func() {
		m.Post("/", binding.BindIgnErr(auth.GiteaLoginForm{}), auth.LoginToGitea)
	})
	m.Combo("/repos", reqSignIn).Get(migration.ListRepos).Post(migration.ListReposPost)
	m.Get("/status", reqSignIn, migration.StatusReport)
	return m
}

func reqSignIn(ctx *context.Context) {
	if ctx.User == nil || ctx.GiteaUser == nil {
		ctx.Redirect("/")
	}
}
