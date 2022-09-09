package migrations

import (
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun/migrate"
)

// //go:embed *.sql
// var files embed.FS

var mig = migrate.NewMigrations()

var logger = log.Logger.With().Str("location", "migrations").Logger()

func GetMigrations() (*migrate.Migrations, error) {
	// if err := mig.Discover(files); err != nil {
	// 	return nil, errors.WithStack(err)
	// }
	return mig, nil
}
