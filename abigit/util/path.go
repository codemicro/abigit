package util

import (
	"github.com/pkg/errors"
	"io/fs"
	"path/filepath"
)

func DirSize(path string) (int64, error) {
	var size int64

	err := filepath.WalkDir(path, func(_ string, info fs.DirEntry, err error) error {
		if err != nil {
			return errors.WithStack(err)
		}
		if !info.IsDir() {
			fi, err := info.Info()
			if err != nil {
				return errors.WithStack(err)
			}
			size += fi.Size()
		}
		return nil
	})

	if err != nil {
		return 0, errors.WithStack(err)
	}

	return size, nil
}
