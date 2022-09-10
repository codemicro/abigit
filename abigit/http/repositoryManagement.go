package http

import (
	"fmt"
	"github.com/codemicro/abigit/abigit/config"
	"github.com/codemicro/abigit/abigit/core"
	"github.com/codemicro/abigit/abigit/http/views"
	"github.com/codemicro/abigit/abigit/urls"
	"github.com/codemicro/abigit/abigit/util"
	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"strings"
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

// serveRepository is designed to be used as a middleware
func (e *Endpoints) serveRepository(ctx *fiber.Ctx) error {
	if method := ctx.Method(); !(method == fiber.MethodGet || method == fiber.MethodHead) {
		return ctx.Next()
	}

	logger := log.With().Str("location", "serveRepository").Logger()

	repositorySlug := ctx.Params("slug")
	if repositorySlug == "" || !strings.HasSuffix(repositorySlug, ".git") {
		logger.Debug().Msg("empty slug or no .git suffix")
		return ctx.Next()
	}

	prefix := urls.MakeRelative(ctx.Route().Path, repositorySlug)
	repositorySlug = repositorySlug[:len(repositorySlug)-4] // remove ".git"

	if calcSlug := slug.Make(repositorySlug); repositorySlug != calcSlug {
		// this means something problematic is going on, like an attempted path traversal
		logger.Debug().Msgf("slugs don't match: %s != %s", repositorySlug, calcSlug)
		return fiber.ErrBadRequest
	}

	requestedFile := strings.TrimPrefix(ctx.Path(), prefix)

	if len(requestedFile) == 0 {
		return ctx.Redirect(urls.MakeRelative(urls.RepositoryByName, repositorySlug))
	}

	if strings.Contains(requestedFile, "..") {
		logger.Debug().Msgf("requested file contains possible path traversal: %s", requestedFile)
		return fiber.ErrBadRequest
	}

	osFilepath := filepath.Join(config.Git.RepositoriesPath, repositorySlug, requestedFile)

	log.Debug().Msg(osFilepath)

	stat, err := os.Stat(osFilepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fiber.ErrNotFound
		}
		return errors.WithStack(err)
	}

	if stat.IsDir() {
		logger.Debug().Msg("directory")
		return fiber.ErrBadRequest
	}

	contentLength := int(stat.Size())

	file, err := os.Open(osFilepath)
	if err != nil {
		return errors.WithStack(err)
	}

	ctx.Set(fiber.HeaderContentType, "application/octet-stream")

	switch ctx.Method() {
	case fiber.MethodGet:
		ctx.Response().SetBodyStream(file, contentLength)
	case fiber.MethodHead:
		ctx.Request().ResetBody()
		ctx.Response().SkipBody = true
		ctx.Response().Header.SetContentLength(contentLength)
		if err := file.Close(); err != nil {
			return err
		}
	}

	return nil
}
