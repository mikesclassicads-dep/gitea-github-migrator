package context

import (
	bgctx "context"
	"fmt"
	"strings"

	"code.gitea.io/sdk/gitea"
	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/config"
	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/migrations"
	"github.com/go-macaron/session"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"gopkg.in/macaron.v1"
)

// Context represents context of a request.
type Context struct {
	*macaron.Context
	Flash   *session.Flash
	Session session.Store

	Client      *github.Client
	GiteaClient *gitea.Client
	User        *User //GitHub user
	GiteaUser   *User
	Link        string // current request URL
}

// User is an abstraction of a Gitea or GitHub user, saving the required information
type User struct {
	ID        int64
	Username  string
	AvatarURL string
	Token     string
}

var runningJobs = make(map[string]*migrations.Job)

// GetCurrentJob returns the current job of the user
// Bug(JonasFranzDEV): prevents scalability (FIXME)
func (ctx *Context) GetCurrentJob() *migrations.Job {
	return runningJobs[ctx.Session.ID()]
}

// SetCurrentJob sets the current job of the user
// Bug(JonasFranzDEV): prevents scalability (FIXME)
func (ctx *Context) SetCurrentJob(job *migrations.Job) {
	runningJobs[ctx.Session.ID()] = job
}

// Handle displays the corresponding error message
func (ctx *Context) Handle(status int, title string, err error) {
	if err != nil {
		if macaron.Env != macaron.PROD {
			ctx.Data["ErrorMsg"] = err
		}
	}
	logrus.Warnf("Handle: %v", err)
	ctx.Data["ErrTitle"] = title

	switch status {
	case 403:
		ctx.Data["Title"] = "Access denied"
	case 404:
		ctx.Data["Title"] = "Page not found"
	case 500:
		ctx.Data["Title"] = "Internal Server Error"
	default:
		ctx.Context.HTML(status, "status/unknown_error")
		return
	}
	ctx.Context.HTML(status, fmt.Sprintf("status/%d", status))
}

// Contexter injects context.Context into macaron
func Contexter() macaron.Handler {
	return func(c *macaron.Context, sess session.Store, f *session.Flash) {
		ctx := &Context{
			Context: c,
			Flash:   f,
			Session: sess,
			Link:    c.Req.URL.String(),
		}
		c.Data["Link"] = ctx.Link
		if ctx.Req.Method == "POST" && strings.Contains(ctx.Req.Header.Get("Content-Type"), "multipart/form-data") {
			if err := ctx.Req.ParseMultipartForm(5242880); err != nil &&
				strings.Contains(err.Error(), "EOF") {
				ctx.Handle(500, "ParseMultipartForm", err)
			}
		}
		ctx.Data["Config"] = config.Config
		usr := sess.Get("user")
		if usr != nil {
			ctx.User = usr.(*User)
			ctx.Data["User"] = ctx.User
		}
		giteaUsr := sess.Get("gitea_user")
		if giteaUsr != nil {
			ctx.GiteaUser = giteaUsr.(*User)
			ctx.Data["GiteaUser"] = ctx.GiteaUser
		}
		if ctx.User != nil && ctx.User.Token != "" {
			tc := oauth2.NewClient(bgctx.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ctx.User.Token}))
			ctx.Client = github.NewClient(tc)
		} else {
			ctx.Client = github.NewClient(nil)
		}
		if giteaURL, ok := sess.Get("gitea").(string); ok && giteaURL != "" && ctx.GiteaUser != nil && ctx.GiteaUser.Token != "" {
			ctx.GiteaClient = gitea.NewClient(giteaURL, ctx.GiteaUser.Token)
		}
		c.Map(ctx)
	}
}
