package main

import (
	"github.com/rs/zerolog/log"
)

func main() {
	app, err := setupApp()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not setup server")
	} else {
		log.Info().Msg("Starting server")
		app.Start()
	}
}
