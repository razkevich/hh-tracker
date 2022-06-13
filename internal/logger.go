package internal

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

// LoggerConfig is the collection of middlewares used in system
type LoggerConfig struct {
	GlobalLogLevel string
}

// ProvideLogger configure logger
func ProvideLogger(conf LoggerConfig) zerolog.Logger {
	level, err := zerolog.ParseLevel(conf.GlobalLogLevel)

	// Datadog prefers the keyword 'status' for our log level
	// https://docs.datadoghq.com/logs/log_configuration/attributes_naming_convention/#reserved-attributes
	zerolog.LevelFieldName = "status"
	zerolog.TimestampFieldName = "timestamp"

	if err != nil {
		level = zerolog.InfoLevel
		defer log.Warn().Msg("Could not parse logging level: `" + conf.GlobalLogLevel + "`, defaulting to `Info`")
	}
	zerolog.SetGlobalLevel(level)

	defer log.Warn().Msg("Setting logging level to " + level.String())

	return zerolog.New(os.Stdout).With().Timestamp().Logger()

}
