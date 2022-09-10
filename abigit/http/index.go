package http

import (
	"github.com/codemicro/abigit/abigit/core"
	views2 "github.com/codemicro/abigit/abigit/http/views"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
)

func (e *Endpoints) Index(ctx *fiber.Ctx) error {
	repos, err := core.ListRepositories()
	if err != nil {
		return errors.WithStack(err)
	}

	rctx, err := e.newRenderContext(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	return views2.SendPage(
		ctx,
		views2.Index(
			rctx,
			&views2.IndexProps{
				AllRepos: repos,
			},
		),
	)
}
