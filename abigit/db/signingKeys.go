package db

import (
	"crypto/rand"
	"database/sql"
	"github.com/codemicro/abigit/abigit/db/models"
	"github.com/pkg/errors"
)

func (db *DB) FetchSigningKey() (*models.SigningKey, error) {
	ctx, cancel := db.newContext()
	defer cancel()
	o := new(models.SigningKey)
	if err := db.bun.NewSelect().Model(o).Scan(ctx); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			randomData := make([]byte, 64)
			_, _ = rand.Read(randomData)
			o.Key = randomData
			o.ID = db.snowflake.Generate()
			if _, err := db.bun.NewInsert().Model(o).Exec(ctx); err != nil {
				return nil, errors.WithStack(err)
			}
			return o, nil
		}
		return nil, errors.WithStack(err)
	}
	return o, nil
}
