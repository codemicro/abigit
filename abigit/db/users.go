package db

import (
	"database/sql"
	"github.com/codemicro/abigit/abigit/db/models"
	"github.com/codemicro/abigit/abigit/util"
	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"strings"
)

var (
	ErrEmailInUse      = util.NewRichError(409, "email address in use", nil)
	ErrExternalIDInUse = util.NewRichError(409, "external ID in use", nil)
)

func (db *DB) RegisterUser(user *models.User) error {
	user.ID = db.snowflake.Generate()

	ctx, cancel := db.newContext()
	defer cancel()

	if _, err := db.bun.NewInsert().Model(user).Exec(ctx); err != nil {
		if e, ok := err.(sqlite3.Error); ok {
			if e.ExtendedCode == sqlite3.ErrConstraintUnique {
				if strings.Contains(e.Error(), "extern_id") {
					return ErrExternalIDInUse
				} else if strings.Contains(e.Error(), "email") {
					return ErrEmailInUse
				}
			}
		}
		return errors.WithStack(err)
	}

	return nil
}

func (db *DB) GetUserByExternalID(extern string) (*models.User, error) {
	ctx, cancel := db.newContext()
	defer cancel()

	o := new(models.User)
	if err := db.bun.NewSelect().Model(o).Where("extern_id = ?", extern).Scan(ctx, o); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, errors.WithStack(err)
	}

	return o, nil
}
