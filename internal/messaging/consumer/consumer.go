package consumer

import (
	"context"
	cloud "github.com/cloudevents/sdk-go/v2"
	"github.com/rs/zerolog/log"
	mongo "gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/driver"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/messaging/bus"
	domain_services "gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/messaging/consumer/domain-services"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/messaging/consumer/handler"
	"gitlab.elasticpath.com/commerce-cloud/personal-data.svc/internal/service"
	"time"
)

// Consumer consumes messages from a bus and runs handler.go to process them
type Consumer struct {
	mongoClient *mongo.Client
	client      bus.IMessageBusClient
	handler     handler.Function
	run         *runContext
}

type runContext struct {
	quit   chan struct{}
	closed chan struct{}
}

// Run processes events from the message bus, projecting those events into the database.
// This function does not return until signalled and would normally be run in a goroutine.
func (c *Consumer) Run() {
	log.Print("starting Message bus consumer in Goroutine")
	if c.run != nil {
		panic("Run() invoked twice")
	}
	c.run = &runContext{
		quit:   make(chan struct{}),
		closed: make(chan struct{}),
	}
	go c.consume()
}

func (c *Consumer) consume() {
	for {
		select {
		case <-c.run.quit:
			close(c.run.closed)
			return
		case <-time.After(time.Second * 30):
			err := c.client.ReconnectIfDisconnected(context.Background())
			if err != nil {
				log.Error().Err(err).Msg("error connecting to RabbitMQ")
				continue
			}
			err = c.client.Consume(context.Background(), c.run.quit, c.process)
			if err != nil {
				log.Error().Err(err).Msg("error in message bus consumer")
			}
			log.Error().Msg("message bus consumer stopped")
		}
	}
}

func (c *Consumer) process(ctx context.Context, event cloud.Event) (bool, error) {
	return c.handler(ctx, event, c.mongoClient)
}

// Stop stops a running consumer, waiting until it has stopped before returning
func (c *Consumer) Stop() {
	if c.run == nil {
		return
	}
	close(c.run.quit)
	<-c.run.closed
	c.run = nil
}

// ProvideConsumer is a wire provider function for the consumer
func ProvideConsumer(mongoClient *mongo.Client, bus bus.IMessageBusClient, logEntriesService service.LogEntriesService) *Consumer {
	c := &Consumer{
		mongoClient: mongoClient,
		client:      bus,
		handler:     domain_services.Handler(logEntriesService, mongoClient),
	}
	return c
}
