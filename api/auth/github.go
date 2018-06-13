package auth

import (
	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/api/responses"
	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/config"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	githubauth "golang.org/x/oauth2/github"
)

var (
	oauthConfig *oauth2.Config
)

func InitOAuthConfig() {
	oauthConfig = &oauth2.Config{
		ClientID:     config.Config.GitHub.ClientID,
		ClientSecret: config.Config.GitHub.ClientSecret,
		Scopes:       []string{"repo"},
		Endpoint:     githubauth.Endpoint,
	}
}

// RedirectToGitHub returns the redirect URL for github
// swagger:route GET /auth/github getGitHubRedirect
//     Produces:
//       - application/json
//     Responses:
//       200: Redirect
func RedirectToGitHub(c *gin.Context) {
	c.JSON(200, responses.RedirectResponse{
		URL: oauthConfig.AuthCodeURL("bla", oauth2.AccessTypeOnline),
	})
}
