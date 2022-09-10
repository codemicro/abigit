package migrations

import (
	"context"
	"github.com/codemicro/abigit/abigit/models"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func init() {
	tps := []any{
		(*models.User)(nil),
		(*models.SigningKey)(nil),
	}

	mig.MustRegister(func(ctx context.Context, db *bun.DB) error {
		logger.Debug().Msg("1 up")

		for _, t := range tps {
			if _, err := db.NewCreateTable().Model(t).Exec(ctx); err != nil {
				return errors.WithStack(err)
			}
		}

		return nil
	},
		func(ctx context.Context, db *bun.DB) error {
			logger.Debug().Msg("1 down")

			for _, t := range tps {
				if _, err := db.NewDropTable().Model(t).Exec(ctx); err != nil {
					return errors.WithStack(err)
				}
			}

			return nil
		})
}
