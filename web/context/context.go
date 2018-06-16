package context

import (
	"fmt"
	"strings"

	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/config"
	"github.com/go-macaron/session"
	"github.com/google/go-github/github"
	"gopkg.in/macaron.v1"
)

// Context represents context of a request.
type Context struct {
	*macaron.Context
	Flash   *session.Flash
	Session session.Store

	Client *github.Client
	User   *User
	Link   string // current request URL
}

type User struct {
	Username  string
	AvatarURL string
	Token     string
}

func (ctx *Context) Handle(status int, title string, err error) {
	if err != nil {
		if macaron.Env != macaron.PROD {
			ctx.Data["ErrorMsg"] = err
		}
	}
	ctx.Data["ErrTitle"] = title

	switch status {
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
		c.Map(ctx)
	}
}
