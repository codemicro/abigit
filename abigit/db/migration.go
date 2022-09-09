package db

import (
	"context"
	_ "embed"
	"github.com/codemicro/abigit/abigit/db/migrations"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun/migrate"
	"time"
)

func (db *DB) Migrate() error {
	log.Info().Msg("running migrations")

	migs, err := migrations.GetMigrations()
	if err != nil {
		return errors.WithStack(err)
	}

	mig := migrate.NewMigrator(db.bun, migs)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := mig.Init(ctx); err != nil {
		return errors.WithStack(err)
	}

	group, err := mig.Migrate(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	if group.IsZero() {
		log.Info().Msg("database up to date")
	} else {
		log.Info().Msg("migrations applied")
	}

	return nil
}
