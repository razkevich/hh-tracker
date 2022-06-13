package domainservices

import (
	"context"
	cloud "github.com/cloudevents/sdk-go/v2"
	mongo "gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/driver"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/logging"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/messaging/consumer/handler"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/messaging/consumer/service"
)

// Handler takes a service interface and return a handler
func Handler(service service.Service, mongoClient *mongo.Client) handler.Function {
	return func(ctx context.Context, event cloud.Event, mongoClient *mongo.Client) (bool, error) {
		logging.GetLoggerFromCtx(ctx).Warn().Msg("hello world")
		return true, nil
	}
}
