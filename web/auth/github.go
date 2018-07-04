package auth

import (
	"context"
	"net/http"

	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/config"
	webcontext "git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/web/context"
	"github.com/adam-hanna/randomstrings"
	"github.com/go-macaron/session"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	githubauth "golang.org/x/oauth2/github"
)

var (
	githubOAuthConfig *oauth2.Config
)

// InitGitHubOAuthConfig loads values from config into githubOAuthConfig
func InitGitHubOAuthConfig() {
	githubOAuthConfig = &oauth2.Config{
		ClientID:     config.Config.GitHub.ClientID,
		ClientSecret: config.Config.GitHub.ClientSecret,
		Scopes:       []string{"repo"},
		Endpoint:     githubauth.Endpoint,
	}
}

// RedirectToGitHub returns the redirect URL for github
func RedirectToGitHub(ctx *webcontext.Context, session session.Store) {
	state, err := randomstrings.GenerateRandomString(64)
	if err != nil {
		return
	}
	session.Set("state", state)
	ctx.Redirect(githubOAuthConfig.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

// CallbackFromGitHub handles the callback from the GitHub OAuth provider
func CallbackFromGitHub(ctx *webcontext.Context, session session.Store) {
	bg := context.Background()
	var state string
	var ok bool
	if state, ok = session.Get("state").(string); state == "" || !ok || state != ctx.Query("state") {
		ctx.Handle(400, "invalid session", nil)
		return
	}
	token, err := githubOAuthConfig.Exchange(bg, ctx.Query("code"))
	if err != nil {
		ctx.Handle(403, "access denied", err)
		return
	}
	tc := oauth2.NewClient(bg, oauth2.StaticTokenSource(token))
	client := github.NewClient(tc)
	user, _, err := client.Users.Get(bg, "")
	if err != nil {
		ctx.Handle(403, "access denied", err)
		return
	}
	session.Set("user", &webcontext.User{
		ID:        user.GetID(),
		AvatarURL: *user.AvatarURL,
		Username:  user.GetLogin(),
		Token:     token.AccessToken,
	})
	ctx.Redirect("/")
}
