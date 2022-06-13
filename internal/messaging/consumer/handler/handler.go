package handler

import (
	"context"
	cloud "github.com/cloudevents/sdk-go/v2"
	mongo "gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/driver"
)

// Function is a handler function called with a specific event is consumed.
// It returns a boolean indicating whether the event should be an acknowledged
// or an error which, if not nil, will cause event consumption to cease.
type Function func(ctx context.Context, event cloud.Event, mongoClient *mongo.Client) (bool, error)
