package http

import (
	"bytes"
	"fmt"
	"github.com/codemicro/abigit/abigit/config"
	"github.com/codemicro/abigit/abigit/core"
	"github.com/codemicro/abigit/abigit/http/views"
	"github.com/codemicro/abigit/abigit/urls"
	"github.com/codemicro/abigit/abigit/util"
	"github.com/codemicro/htmlRenderer"
	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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

	if !core.ValidateSlug(repositorySlug) {
		// this means something problematic is going on, like an attempted path traversal
		logger.Debug().Msgf("slugs don't match")
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

func (e *Endpoints) displayRepository(ctx *fiber.Ctx) error {
	repoSlug := ctx.Params("slug")
	if repoSlug == "" || !core.ValidateSlug(repoSlug) {
		return fiber.ErrBadRequest
	}

	repoInfo, err := core.GetRepository(repoSlug)
	if err != nil {
		return errors.WithStack(err)
	}

	repo, err := gogit.PlainOpen(repoInfo.Path)
	if err != nil {
		return errors.WithStack(err)
	}

	//refIter, err := repo.References()
	//if err != nil {
	//	return errors.WithStack(err)
	//}
	//defer refIter.Close()

	//var refs []*plumbing.Reference
	//
	//if err := refIter.ForEach(func(ref *plumbing.Reference) error {
	//	fmt.Println(ref.Name(), ref.Type(), ref.Hash().String(), ref.Hash().IsZero())
	//	refs = append(refs, ref)
	//	return nil
	//}); err != nil {
	//	return errors.WithStack(err)
	//}

	isEmpty, err := core.IsRepositoryEmpty(repo)
	if err != nil {
		return errors.WithStack(err)
	}

	rctx, err := e.newRenderContext(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	return views.SendPage(
		ctx,
		views.ViewRepository(rctx, &views.ViewRepositoryProps{
			Repo:    repoInfo,
			IsEmpty: isEmpty,
		}),
	)
}

func (e *Endpoints) repositorySizeOnDisk(ctx *fiber.Ctx) error {
	repoSlug := ctx.Params("slug")
	if repoSlug == "" || !core.ValidateSlug(repoSlug) {
		return fiber.ErrBadRequest
	}

	repo, err := core.GetRepository(repoSlug)
	if err != nil {
		return errors.WithStack(err)
	}

	repoSize, err := repo.Size()
	if err != nil {
		return errors.WithStack(err)
	}

	return ctx.SendString(views.FormatFileSize(repoSize))
}

func (e *Endpoints) repositoryTabs(ctx *fiber.Ctx) error {
	repoSlug := ctx.Params("slug")
	if repoSlug == "" || !core.ValidateSlug(repoSlug) {
		return fiber.ErrBadRequest
	}

	repoInfo, err := core.GetRepository(repoSlug)
	if err != nil {
		return errors.WithStack(err)
	}

	repo, err := gogit.PlainOpen(repoInfo.Path)
	if err != nil {
		return errors.WithStack(err)
	}

	props := &views.RepositoryTabProps{
		Repo: repoInfo,
	}

	switch strings.ToLower(ctx.Query("tab", "readme")) {
	case "tree":
		props.CurrentTab = views.TabSelectorShowTree
	case "refs":
		props.CurrentTab = views.TabSelectorShowRefs

		ri, err := repo.References()
		if err != nil {
			return errors.WithStack(err)
		}

		if err := ri.ForEach(func(ref *plumbing.Reference) error {

			if ref.Type() == plumbing.HashReference {
				if ref.Name().IsBranch() {
					props.Refs.Branches = append(props.Refs.Branches, ref)
				}
				if ref.Name().IsTag() {
					props.Refs.Tags = append(props.Refs.Tags, ref)
				}
			}

			return nil
		}); err != nil {
			return errors.WithStack(err)
		}

	case "clone":
		props.CurrentTab = views.TabSelectorClone

		props.Clone.SSHUser = config.SSH.User
		props.Clone.SSHHost = config.SSH.Host
		props.Clone.SSHStoragePath = config.Git.RepositoriesPath
	case "commits":
		props.CurrentTab = views.TabSelectorCommits
	case "readme":
		fallthrough
	default:
		props.CurrentTab = views.TabSelectorReadme

		readmeContent, err := core.GetReadmeContent(repo)
		if err != nil {
			if errors.Is(err, core.ErrNoReadme) {
				defaultBranch, err := core.GetDefaultBranch(repo)
				if err != nil {
					log.Warn().Msg("could not fetch default branch")
				}

				readmeContent = []byte("*No README file available*")

				if defaultBranch != "" {
					readmeContent = append(readmeContent,
						[]byte(" - *make sure there's a file called README.md in the root of the repository on the `"+defaultBranch+"` branch*")...,
					)
				}
			} else {
				return errors.WithStack(err)
			}
		}

		buf := new(bytes.Buffer)
		markdownProcessor := htmlRenderer.NewProcessor(htmlRenderer.WithHeaderLinks())
		if err := markdownProcessor.Convert(readmeContent, buf); err != nil {
			return errors.WithStack(err)
		}

		props.Readme.Content = buf.String()
	}

	rctx, err := e.newRenderContext(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	return views.SendPage(
		ctx,
		views.RepositoryTabs(rctx, props),
	)
}
