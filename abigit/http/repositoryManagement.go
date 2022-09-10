package http

import (
	"fmt"
	"github.com/codemicro/abigit/abigit/core"
	"github.com/codemicro/abigit/abigit/http/views"
	"github.com/codemicro/abigit/abigit/urls"
	"github.com/codemicro/abigit/abigit/util"
	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"
)

func (e *Endpoints) createRepository(ctx *fiber.Ctx) error {
	si := e.getSessionInformation(ctx)
	if si == nil {
		return util.NewRichErrorFromFiberError(fiber.ErrUnauthorized, "You must be logged in to do that.")
	}

	rctx, err := e.newRenderContext(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	switch ctx.Method() {
	case fiber.MethodGet:
		return views.SendPage(
			ctx,
			views.CreateRepository(rctx, new(views.CreateRepositoryProps)),
		)
	case fiber.MethodPost:
		props := new(views.CreateRepositoryProps)

		repoName := ctx.FormValue("name")

		rod, err := core.CreateRepository(repoName)
		if err != nil {
			if util.IsRichError(err) {
				re := err.(*util.RichError)
				ctx.Status(re.Status)
				props.Problem = re.Reason
				goto respondError
			}
			return errors.WithStack(err)
		}

		return ctx.Redirect(urls.Make(urls.RepositoryByName, rod.Slug))
	respondError:
		return views.SendPage(
			ctx,
			views.CreateRepository(rctx, props),
		)
	default:
		return util.NewRichErrorFromFiberError(fiber.ErrMethodNotAllowed, nil)
	}
}

func (e *Endpoints) createRepositoryValidation(ctx *fiber.Ctx) error {
	ctx.Type("html")

	repoName := ctx.FormValue("name")
	slug := slug.Make(repoName)
	if err := core.ValidateRepositoryName(slug); err != nil {
		if util.IsRichError(err) {
			return ctx.SendString(fmt.Sprintf(`<span style="color: red;">%s</span>`, err.(*util.RichError).Reason))
		}
	}

	if slug != repoName {
		return ctx.SendString(fmt.Sprintf("Will be created with the name <b>%s</b>", slug))
	}

	return nil
}
