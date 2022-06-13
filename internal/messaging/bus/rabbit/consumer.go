package rabbit

import (
	"context"
	"encoding/json"
	cloud "github.com/cloudevents/sdk-go/v2"
	"github.com/rs/zerolog/log"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/logging"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/requestid"
	"strconv"
)

// Consume the event consumer
func (mb *MessageBus) Consume(globalCtx context.Context, quit <-chan struct{}, process func(ctx context.Context, event cloud.Event) (ack bool, err error)) error {
	channel, err := mb.connection.Channel()
	if err != nil {
		return mb.err("failed to open rabbit channel", err)
	}
	defer channel.Close()
	deliveries, err := channel.Consume(mb.queue, "", false, false, false, false, nil)
	if err != nil {
		return mb.err("failed to consume rabbit queue", err)
	}
	for {
		select {
		case <-quit:
			return nil
		case delivery, ok := <-deliveries:
			if !ok {
				return nil
			}
			var event cloud.Event
			deliveryErr := json.Unmarshal(delivery.Body, &event)
			if deliveryErr != nil {
				rejectedErr := delivery.Reject(false)
				if rejectedErr != nil {
					log.Err(rejectedErr).Msg("failed to send Reject message to the server")
				}
				log.Err(deliveryErr).Msg("failed to decode(unmarshal) rabbit delivery")
				continue
			}
			requestID, storeID := requestid.GetRequestIDAndStoreIDFromCloudEvent(event)

			ctx := context.WithValue(globalCtx, "x_request_id", requestID)                  //nolint
			ctx = context.WithValue(ctx, "logger", logging.BuildLogger(storeID, requestID)) //nolint

			ack, processErr := process(ctx, event)
			if processErr != nil {
				logging.GetLoggerFromCtx(ctx).Err(processErr).Msg("failed to process rabbit delivery")
			}
			if ack && processErr == nil {
				err = delivery.Ack(false)
			} else {
				err = delivery.Reject(false)
			}
			if err != nil {
				logging.GetLoggerFromCtx(ctx).Err(err).Msgf("failed to send message %s for rabbit delivery", strconv.FormatBool(ack))
			}
		}
	}
}
