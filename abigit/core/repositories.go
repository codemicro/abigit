package core

import (
	"github.com/codemicro/abigit/abigit/config"
	"github.com/codemicro/abigit/abigit/util"
	"github.com/go-git/go-git/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/pkg/errors"
	"os"
	"os/exec"
	"path/filepath"
)

type RepoOnDisk struct {
	Slug        string
	Description string

	Path string
	Size int64
}

// ListRepositories returns the list of repositories in a given directory.
//
// It assumes that any directory is a repository.
func ListRepositories() ([]*RepoOnDisk, error) {
	dirEntries, err := os.ReadDir(config.Git.RepositoriesPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var o []*RepoOnDisk

	for _, entry := range dirEntries {
		if !entry.IsDir() {
			continue
		}

		fi, err := entry.Info()
		if err != nil {
			return nil, errors.WithStack(err)
		}

		o = append(o, &RepoOnDisk{
			Slug:        entry.Name(),
			Description: "Lorem ipsum dolor sit amet.",

			Path: filepath.Join(config.Git.RepositoriesPath, entry.Name()),
			Size: fi.Size(),
		})
	}

	return o, nil
}

func doesRepositoryExist(slug string) (bool, error) {
	dirEntries, err := os.ReadDir(config.Git.RepositoriesPath)
	if err != nil {
		return false, errors.WithStack(err)
	}

	for _, entry := range dirEntries {
		if entry.Name() == slug {
			return true, nil
		}
	}

	return false, nil
}

func GetRepository(slug string) (*RepoOnDisk, error) {
	if found, err := doesRepositoryExist(slug); err != nil {
		return nil, errors.WithStack(err)
	} else if !found {
		return nil, util.NewRichErrorFromFiberError(fiber.ErrNotFound, nil)
	}

	fp := filepath.Join(config.Git.RepositoriesPath, slug)

	stat, err := os.Stat(fp)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// TODO: populate with meaningful information
	return &RepoOnDisk{
		Slug:        slug,
		Description: "Lorem ipsum dolor sit amet.",

		Path: fp,
		Size: stat.Size(),
	}, nil
}

func ValidateRepositoryName(repoNameSlug string) error {
	if repoNameSlug == "" {
		return util.NewRichError(fiber.StatusBadRequest, "Repository name cannot be empty", nil)
	} else if len(repoNameSlug) >= 128 {
		return util.NewRichError(fiber.StatusBadRequest, "Repository name too long (max: 128 chars)", nil)
	}

	repoExists, err := doesRepositoryExist(repoNameSlug)
	if err != nil {
		return errors.WithStack(err)
	}
	if repoExists {
		return util.NewRichError(fiber.StatusBadRequest, "Name already in use", nil)
	}

	return nil
}

const postUpdateHookContents = `#!/bin/sh
exec git update-server-info
`

func CreateRepository(name string) (*RepoOnDisk, error) {
	rod := &RepoOnDisk{
		Slug: slug.Make(name),
	}

	if err := ValidateRepositoryName(rod.Slug); err != nil {
		return nil, err
	}

	rod.Path = filepath.Join(config.Git.RepositoriesPath, rod.Slug)

	_, err := git.PlainInit(rod.Path, true)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Create `post-update` hook
	if err := os.MkdirAll(filepath.Join(rod.Path, "hooks"), os.ModeDir); err != nil {
		return nil, errors.WithStack(err)
	}

	if err := os.WriteFile(
		filepath.Join(rod.Path, "hooks", "post-update"),
		[]byte(postUpdateHookContents),
		775,
	); err != nil {
		return nil, errors.WithStack(err)
	}

	// We need to execute it so it can be cloned
	cmd := exec.Command(filepath.Join("hooks", "post-update"))
	cmd.Dir = rod.Path
	if err := cmd.Run(); err != nil {
		return nil, errors.WithStack(err)
	}

	return rod, nil
}

func IsRepositoryEmpty(repo *git.Repository) (bool, error) {
	ref, err := repo.Reference("HEAD", false)
	if err != nil {
		return false, errors.WithStack(err)
	}

	return ref.Hash().IsZero(), nil
}
