package requestid

import (
	"context"
	cloud "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/logging"
)

// StoreIDOnlyEvent json type used to read store_id
type StoreIDOnlyEvent struct {
	StoreID string `json:"store_id"`
}

// GetRequestIDFromCtx Determine the request id from context
func GetRequestIDFromCtx(ctx context.Context) string {
	if ctx.Value("x_request_id") == nil {

		requestID := "pds-auto-" + uuid.New().String()
		ctx = context.WithValue(ctx, "x_request_id", requestID) //nolint
		logging.GetLoggerFromCtx(ctx).Warn().Msg("Could not find a request id in the context, generating a new one")

		return requestID
	}
	return ctx.Value("x_request_id").(string)
}

// GetRequestIDAndStoreIDFromCloudEvent Determine the request id and store id from a cloud event
func GetRequestIDAndStoreIDFromCloudEvent(event cloud.Event) (string, string) {
	var requestID string
	var storeID string
	var err error

	if requestID, err = types.ToString(event.Extensions()["requestid"]); err != nil {
		requestID = "pds-gen-" + uuid.New().String()
		log.Warn().Str("x_request_id", requestID).Err(err).Msgf("Unknown requestid on event type: %s id: %s  received value[%s], using %s", event.Type(), event.ID(), event.Extensions()["requestid"], requestID)
	}
	var e StoreIDOnlyEvent
	err = event.DataAs(&e)

	if err != nil {
		log.Warn().Str("x_request_id", requestID).Err(err).Msg("failed unmarshal object to retrieve store id")
		storeID = "unknown"
	}

	return requestID, storeID
}
