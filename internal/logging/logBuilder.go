package logging

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// BuildLogger Create a logger that has the store id and requestid as part of it's context
func BuildLogger(storeUUID, requestID string) zerolog.Logger {
	return log.With().Str("store_uuid", storeUUID).Str("x_request_id", requestID).Logger()
}

// GetLoggerFromCtx Retrieve the logger object from the context
func GetLoggerFromCtx(ctx context.Context) *zerolog.Logger {
	if ctx.Value("logger") == nil {
		logger := log.Logger

		logger.Error().Msg("No logger found in context")
		return &logger
	}
	val := ctx.Value("logger").(zerolog.Logger)
	return &val
}
