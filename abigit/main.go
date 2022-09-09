package main

import (
	"fmt"
	"github.com/codemicro/abigit/abigit/config"
	"github.com/codemicro/abigit/abigit/db"
	"github.com/codemicro/abigit/abigit/endpoints"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
)

func run() error {
	database, err := db.New()
	if err != nil {
		return errors.WithStack(err)
	}

	if err := database.Migrate(); err != nil {
		return errors.Wrap(err, "failed migration")
	}

	e, err := endpoints.New(database)
	if err != nil {
		return errors.WithStack(err)
	}

	app := e.SetupApp()

	serveAddr := config.HTTP.Host + ":" + strconv.Itoa(config.HTTP.Port)

	log.Info().Msgf("starting server on %s", serveAddr)

	if err := app.Listen(serveAddr); err != nil {
		return errors.Wrap(err, "fiber server run failed")
	}

	return nil
}

func main() {
	config.InitLogging()
	if err := run(); err != nil {
		fmt.Printf("%+v\n", err)
		log.Error().Stack().Err(err).Msg("failed to run coordinator")
	}
}
