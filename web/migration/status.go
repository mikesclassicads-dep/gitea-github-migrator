package migration

import (
	"git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/web/context"
)

// StatusReport returns the status of the current Job in JSON
func StatusReport(ctx *context.Context) {
	if job := ctx.GetCurrentJob(); job != nil {
		ctx.JSON(200, job.StatusReport())
		return
	}
	ctx.Status(404)
}
